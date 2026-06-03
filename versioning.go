package go_discord_wrapper

import "runtime/debug"

// modulePath is this module's import path, used to locate its own version in
// the build info.
const modulePath = "github.com/streame-gg/go-discord-wrapper"

const (
	// repositoryVersionFallback is reported when the module version cannot be
	// read from the build (e.g. when building this module directly, where the
	// version is "(devel)", or when build info is unavailable). When the library
	// is consumed as a tagged dependency, RepositoryVersion reflects the real
	// module version instead.
	repositoryVersionFallback = "v0.1.0-alpha"

	RepositoryURL = "https://github.com/streame-gg/go-discord-wrapper"
)

// RepositoryVersion is the library version advertised to Discord in the REST
// User-Agent header and the gateway IDENTIFY properties. It is resolved once at
// startup from the build's module version so it cannot drift out of sync with
// the published tag; it falls back to repositoryVersionFallback for local
// (untagged) builds.
var RepositoryVersion = resolveVersion()

func resolveVersion() string {
	info, ok := debug.ReadBuildInfo()
	return versionFromBuildInfo(info, ok)
}

// versionFromBuildInfo resolves the module version from already-read build info,
// kept separate from resolveVersion so every branch is testable without
// depending on the ambient build (debug.ReadBuildInfo cannot be stubbed).
func versionFromBuildInfo(info *debug.BuildInfo, ok bool) string {
	if !ok || info == nil {
		return repositoryVersionFallback
	}
	if v := info.Main.Version; info.Main.Path == modulePath && isReleaseVersion(v) {
		return v
	}
	for _, dep := range info.Deps {
		if dep.Path == modulePath && isReleaseVersion(dep.Version) {
			return dep.Version
		}
	}
	return repositoryVersionFallback
}

// isReleaseVersion reports whether v is a usable tagged version rather than the
// empty string or the "(devel)" placeholder Go uses for the main module of an
// untagged build.
func isReleaseVersion(v string) bool {
	return v != "" && v != "(devel)"
}
