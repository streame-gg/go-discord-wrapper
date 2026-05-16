// cmd/doc-check extracts code examples from README.md and Go doc comments and
// verifies that each one compiles.
//
// README.md: every fenced ```go block whose first non-empty line starts with
// "package" is compiled in a temporary module that replaces the local module.
//
// Doc comments: every run of consecutive //\t lines from a declaration's doc
// comment is wrapped in a helper function (build tag "doccheck") and compiled
// inside its own package. Examples that appear to be written from a caller's
// perspective (i.e. they reference the current package by name, e.g.
// "options.WithFoo") are skipped — those fragments rely on free variables that
// cannot be reconstructed without the full caller context.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	root, err := moduleRoot()
	if err != nil {
		fatalf("cannot find module root: %v", err)
	}

	var failures []string
	if errs := checkReadme(root); len(errs) > 0 {
		failures = append(failures, errs...)
	}
	if errs := checkGodoc(root); len(errs) > 0 {
		failures = append(failures, errs...)
	}

	if len(failures) > 0 {
		for _, f := range failures {
			fmt.Fprintln(os.Stderr, f)
		}
		os.Exit(1)
	}
	fmt.Println("doc-check: all examples OK")
}

// ─── README ──────────────────────────────────────────────────────────────────

func checkReadme(root string) []string {
	path := filepath.Join(root, "README.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{fmt.Sprintf("README.md: cannot open: %v", err)}
	}

	var failures []string
	blockNum := 0
	inBlock := false
	var cur strings.Builder
	sc := bufio.NewScanner(bytes.NewReader(data))

	for sc.Scan() {
		line := sc.Text()
		switch {
		case line == "```go":
			inBlock = true
			cur.Reset()
		case line == "```" && inBlock:
			inBlock = false
			blockNum++
			code := cur.String()
			trimmed := strings.TrimSpace(code)
			if strings.HasPrefix(trimmed, "package ") {
				label := fmt.Sprintf("README.md block %d", blockNum)
				if err := compileStandaloneBlock(root, code, label); err != nil {
					failures = append(failures, err.Error())
				}
			}
		case inBlock:
			cur.WriteString(line)
			cur.WriteByte('\n')
		}
	}
	return failures
}

func compileStandaloneBlock(root, code, label string) error {
	tmp, err := os.MkdirTemp("", "doccheck-*")
	if err != nil {
		return fmt.Errorf("%s: mktemp: %v", label, err)
	}
	defer os.RemoveAll(tmp)

	modPath, goVer, err := readModInfo(root)
	if err != nil {
		return fmt.Errorf("%s: %v", label, err)
	}

	goMod := fmt.Sprintf(
		"module doccheck_example\n\ngo %s\n\nrequire %s v0.0.0\n\nreplace %s => %s\n",
		goVer, modPath, modPath, root)
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte(goMod), 0o644); err != nil {
		return fmt.Errorf("%s: write go.mod: %v", label, err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "main.go"), []byte(code), 0o644); err != nil {
		return fmt.Errorf("%s: write main.go: %v", label, err)
	}

	out, err := runCmd(tmp, "go", "build", ".")
	if err != nil {
		return fmt.Errorf("%s: compile error:\n%s", label, out)
	}
	return nil
}

// ─── GODOC ───────────────────────────────────────────────────────────────────

type godocBlock struct {
	pkgDir     string
	pkgName    string
	code       string
	// stubs: "var name Type\n\t_ = name\n\t" declarations for receiver+params.
	stubs      string
	// namedImports are imports that must be named (used in stubs).
	namedImports []string
	// blankImports are pass-through imports from the source file.
	blankImports []string
	label        string
}

func checkGodoc(root string) []string {
	byDir := map[string][]godocBlock{}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			n := d.Name()
			if n == ".git" || n == "vendor" || n == "cmd" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		blocks, err := extractGodocBlocks(path)
		if err != nil {
			return err
		}
		for _, b := range blocks {
			byDir[b.pkgDir] = append(byDir[b.pkgDir], b)
		}
		return nil
	})
	if err != nil {
		return []string{fmt.Sprintf("godoc walk: %v", err)}
	}

	var failures []string
	for dir, blocks := range byDir {
		if errs := compileGodocBlocks(root, dir, blocks); len(errs) > 0 {
			failures = append(failures, errs...)
		}
	}
	return failures
}

func extractGodocBlocks(filePath string) ([]godocBlock, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil
	}

	pkgDir := filepath.Dir(filePath)
	pkgName := f.Name.Name
	blankImports := collectImports(f)

	var blocks []godocBlock

	add := func(cg *ast.CommentGroup, stubs string, namedImports []string, label string) {
		if cg == nil {
			return
		}
		code := extractTabBlock(cg)
		if code == "" {
			return
		}
		// Skip prose that happens to use //\t (not valid Go syntax).
		if !looksLikeGoCode(code) {
			return
		}
		// Skip caller-perspective examples: they reference the package by name
		// (e.g. "options.WithFoo") and rely on free variables we cannot reconstruct.
		if strings.Contains(code, pkgName+".") {
			return
		}
		blocks = append(blocks, godocBlock{
			pkgDir:       pkgDir,
			pkgName:      pkgName,
			code:         code,
			stubs:        stubs,
			namedImports: namedImports,
			blankImports: blankImports,
			label:        label,
		})
	}

	// Package-level doc (no declaration to derive stubs from).
	add(f.Doc, "", nil, filepath.Base(filePath)+":package")

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Doc == nil {
				continue
			}
			stubs, named := funcStubs(fset, d)
			label := filepath.Base(filePath) + ":" + d.Name.Name
			add(d.Doc, stubs, named, label)

		case *ast.GenDecl:
			for _, spec := range d.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					doc := ts.Doc
					if doc == nil {
						doc = d.Doc
					}
					add(doc, "", nil, filepath.Base(filePath)+":"+ts.Name.Name)
				}
			}
		}
	}
	return blocks, nil
}

func extractTabBlock(cg *ast.CommentGroup) string {
	var buf strings.Builder
	for _, c := range cg.List {
		if strings.HasPrefix(c.Text, "//\t") {
			buf.WriteString(strings.TrimPrefix(c.Text, "//\t"))
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func looksLikeGoCode(code string) bool {
	src := "package p\nfunc _() {\n" + code + "\n}"
	_, err := parser.ParseFile(token.NewFileSet(), "", src, 0)
	return err == nil
}

// funcStubs returns var declarations for a function's receiver and parameters.
// It also returns the imports that must be named (not blank) to support the types.
func funcStubs(fset *token.FileSet, d *ast.FuncDecl) (stubs string, namedImports []string) {
	type param struct{ name, typ string }
	var params []param

	if d.Recv != nil {
		for _, f := range d.Recv.List {
			typ := nodeStr(fset, f.Type)
			for _, name := range f.Names {
				params = append(params, param{name.Name, typ})
			}
		}
	}
	if d.Type.Params != nil {
		for _, f := range d.Type.Params.List {
			if _, ok := f.Type.(*ast.Ellipsis); ok {
				continue // skip variadic
			}
			typ := nodeStr(fset, f.Type)
			for _, name := range f.Names {
				params = append(params, param{name.Name, typ})
			}
		}
	}

	if len(params) == 0 {
		return "", nil
	}

	importSet := map[string]bool{}
	var sb strings.Builder
	for _, p := range params {
		fmt.Fprintf(&sb, "var %s %s\n\t_ = %s\n\t", p.name, p.typ, p.name)
		if imp := namedImportForType(p.typ); imp != "" {
			importSet[imp] = true
		}
	}

	var named []string
	for imp := range importSet {
		named = append(named, imp)
	}
	return sb.String(), named
}

// namedImportForType maps a type string to the stdlib import it requires.
// Only stdlib packages need this (module-local types are already in scope).
var stdlibTypeImports = map[string]string{
	"context.":  `"context"`,
	"http.":     `"net/http"`,
	"slog.":     `"log/slog"`,
	"time.":     `"time"`,
	"sync.":     `"sync"`,
	"io.":       `"io"`,
	"atomic.":   `"sync/atomic"`,
	"rand.":     `"math/rand"`,
}

func namedImportForType(typ string) string {
	for prefix, imp := range stdlibTypeImports {
		if strings.Contains(typ, prefix) {
			return imp
		}
	}
	return ""
}

func collectImports(f *ast.File) []string {
	var out []string
	for _, imp := range f.Imports {
		out = append(out, imp.Path.Value)
	}
	return out
}

func nodeStr(fset *token.FileSet, n ast.Node) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, n)
	return buf.String()
}

func compileGodocBlocks(root, pkgDir string, blocks []godocBlock) []string {
	if len(blocks) == 0 {
		return nil
	}
	pkgName := blocks[0].pkgName

	// Named imports (required for types in stubs to resolve).
	namedSet := map[string]bool{}
	// Blank imports (pass-through to avoid unused-import errors).
	blankSet := map[string]bool{}

	for _, b := range blocks {
		for _, imp := range b.namedImports {
			namedSet[imp] = true
		}
		for _, imp := range b.blankImports {
			blankSet[imp] = true
		}
	}

	var importLines strings.Builder
	for imp := range namedSet {
		fmt.Fprintf(&importLines, "\t%s\n", imp)
		// Don't also blank-import it.
		delete(blankSet, imp)
	}
	for imp := range blankSet {
		fmt.Fprintf(&importLines, "\t_ %s\n", imp)
	}

	var funcs strings.Builder
	for i, b := range blocks {
		fmt.Fprintf(&funcs, "func _doccheck_%d() {\n", i)
		if b.stubs != "" {
			fmt.Fprintf(&funcs, "\t%s", b.stubs)
		}
		for _, line := range strings.Split(strings.TrimRight(b.code, "\n"), "\n") {
			fmt.Fprintf(&funcs, "\t%s\n", line)
		}
		funcs.WriteString("}\n\n")
	}

	src := fmt.Sprintf("//go:build doccheck\n\npackage %s\n\nimport (\n%s)\n\n%s",
		pkgName, importLines.String(), funcs.String())

	genFile := filepath.Join(pkgDir, "zzz_doccheck_gen.go")
	if err := os.WriteFile(genFile, []byte(src), 0o644); err != nil {
		return []string{fmt.Sprintf("%s: write generated file: %v", pkgDir, err)}
	}
	defer os.Remove(genFile)

	rel, _ := filepath.Rel(root, pkgDir)
	modPath, _, _ := readModInfo(root)
	pkgImport := modPath + "/" + filepath.ToSlash(rel)
	if rel == "." {
		pkgImport = modPath
	}

	out, err := runCmd(root, "go", "build", "-tags", "doccheck", pkgImport)
	if err != nil {
		return []string{fmt.Sprintf("godoc examples in %s: compile error:\n%s", pkgDir, out)}
	}
	return nil
}

// ─── HELPERS ─────────────────────────────────────────────────────────────────

func moduleRoot() (string, error) {
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		return "", err
	}
	mod := strings.TrimSpace(string(out))
	if mod == "" || mod == os.DevNull {
		return "", fmt.Errorf("not in a Go module")
	}
	return filepath.Dir(mod), nil
}

func readModInfo(root string) (modPath, goVer string, err error) {
	data, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", "", err
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "module ") {
			modPath = strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
		if strings.HasPrefix(line, "go ") {
			goVer = strings.TrimSpace(strings.TrimPrefix(line, "go "))
		}
	}
	if modPath == "" {
		return "", "", fmt.Errorf("cannot parse module path from go.mod")
	}
	if goVer == "" {
		goVer = "1.21"
	}
	return modPath, goVer, nil
}

func runCmd(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "doc-check: "+format+"\n", args...)
	os.Exit(1)
}
