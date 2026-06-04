package util

import (
	"encoding/json"
	"reflect"
	"strings"
	"sync"
)

// MergePartial applies non-zero/non-nil fields from partial onto old and
// returns the result. Pointer and interface fields are only applied when
// non-nil; struct fields with exported members are merged recursively, while
// atomic value structs such as time.Time (no exported fields) are replaced
// wholesale when non-zero; all other fields (scalars, slices, maps) are applied
// only when not the zero value.
//
// This is the right tool for Discord partial-update payloads (MESSAGE_UPDATE,
// GUILD_MEMBER_UPDATE, etc.) where Discord omits unchanged fields, leaving them
// as zero values in the unmarshaled struct.
func MergePartial[T any](old, partial T) T {
	mergeValues(reflect.ValueOf(&old).Elem(), reflect.ValueOf(&partial).Elem())
	return old
}

func mergeValues(dst, src reflect.Value) {
	if dst.Kind() != reflect.Struct {
		return
	}
	for i := range dst.NumField() {
		d := dst.Field(i)
		s := src.Field(i)
		if !d.CanSet() {
			continue
		}
		switch d.Kind() {
		case reflect.Ptr, reflect.Interface:
			if !s.IsNil() {
				d.Set(s)
			}
		case reflect.Struct:
			// Value structs with no exported (settable) fields — e.g. time.Time —
			// cannot be merged field-by-field: the recursion would set nothing and a
			// non-zero source value would be silently dropped (GuildMember.JoinedAt,
			// Message.Timestamp, …). Treat them as atomic and replace wholesale when
			// the source is non-zero.
			if hasExportedFields(d.Type()) {
				mergeValues(d, s)
			} else if !s.IsZero() {
				d.Set(s)
			}
		default:
			if !s.IsZero() {
				d.Set(s)
			}
		}
	}
}

// MergePartialJSON applies partial onto old, like [MergePartial], but uses raw —
// the original JSON payload that partial was decoded from — to decide which
// fields are present. A top-level field whose JSON key appears in raw is applied
// even when its value is the zero value (an intentional "", false, 0 or null),
// while fields absent from raw retain old's value.
//
// This is the correct semantics for Discord partial-update payloads that
// deliberately reset a field: MESSAGE_UPDATE clearing content to "", a presence
// going from a custom status back to the empty string, GUILD_MEMBER_UPDATE
// setting deaf to false, and so on. Value-based [MergePartial] cannot tell such
// an intentional zero apart from an omitted field; with the raw payload,
// MergePartialJSON can.
//
// Fields tagged json:"-" are never sourced from JSON and so always retain old's
// value. If raw is not a JSON object (empty, null, or malformed), MergePartialJSON
// falls back to value-based [MergePartial].
func MergePartialJSON[T any](old, partial T, raw json.RawMessage) T {
	present := topLevelKeys(raw)
	if present == nil {
		return MergePartial(old, partial)
	}
	mergePresent(reflect.ValueOf(&old).Elem(), reflect.ValueOf(&partial).Elem(), present)
	return old
}

// topLevelKeys decodes the top-level object keys of raw into a lowercased set.
// It returns nil if raw is not a JSON object, signalling the caller to fall
// back to value-based merging.
func topLevelKeys(raw json.RawMessage) map[string]struct{} {
	if len(raw) == 0 {
		return nil
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil || obj == nil {
		return nil
	}
	keys := make(map[string]struct{}, len(obj))
	for k := range obj {
		keys[strings.ToLower(k)] = struct{}{}
	}
	return keys
}

// mergePresent applies, for each settable field of dst, the corresponding src
// field iff that field's JSON key is present in the set. Present fields are
// taken wholesale (honouring intentional zero values and explicit nulls);
// absent fields are left untouched.
func mergePresent(dst, src reflect.Value, present map[string]struct{}) {
	if dst.Kind() != reflect.Struct {
		return
	}
	t := dst.Type()
	for i := range dst.NumField() {
		d := dst.Field(i)
		if !d.CanSet() {
			continue
		}
		name := jsonKey(t.Field(i))
		if name == "" {
			continue // json:"-" — never sourced from the payload.
		}
		if _, ok := present[name]; ok {
			d.Set(src.Field(i))
		}
	}
}

// jsonKey returns the lowercased JSON object key for f, or "" if f is excluded
// from JSON (json:"-") . A field with no json tag uses its Go field name, which
// encoding/json matches case-insensitively — hence the lowercasing here.
func jsonKey(f reflect.StructField) string {
	tag, ok := f.Tag.Lookup("json")
	if !ok {
		return strings.ToLower(f.Name)
	}
	name, _, _ := strings.Cut(tag, ",")
	if name == "-" {
		return ""
	}
	if name == "" {
		return strings.ToLower(f.Name)
	}
	return strings.ToLower(name)
}

// exportedFieldsCache memoises hasExportedFields per struct type.
var exportedFieldsCache sync.Map // reflect.Type -> bool

// hasExportedFields reports whether t (a struct type) has at least one exported
// field, i.e. whether mergeValues can meaningfully recurse into it.
func hasExportedFields(t reflect.Type) bool {
	if v, ok := exportedFieldsCache.Load(t); ok {
		return v.(bool)
	}
	has := false
	for i := range t.NumField() {
		if t.Field(i).IsExported() {
			has = true
			break
		}
	}
	exportedFieldsCache.Store(t, has)
	return has
}
