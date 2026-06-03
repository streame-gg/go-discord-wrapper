package go_discord_wrapper

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/suite"
)

// VersioningSuite covers the module version resolution helpers.
type VersioningSuite struct {
	suite.Suite
}

func TestVersioningSuite(t *testing.T) { suite.Run(t, new(VersioningSuite)) }

func (s *VersioningSuite) TestIsReleaseVersion() {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"empty is not a release", "", false},
		{"devel placeholder is not a release", "(devel)", false},
		{"semver tag is a release", "v1.2.3", true},
		{"alpha pre-release tag is a release", "v0.1.0-alpha", true},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.Equal(tc.want, isReleaseVersion(tc.in))
		})
	}
}

func (s *VersioningSuite) TestResolveVersion_FallsBackForUntaggedBuild() {
	// Under `go test` the main module reports the "(devel)" version and no dep
	// matches modulePath, so resolveVersion takes the fallback path.
	s.Equal(repositoryVersionFallback, resolveVersion())
}

func (s *VersioningSuite) TestVersionFromBuildInfo() {
	s.Run("build info unavailable → fallback", func() {
		s.Equal(repositoryVersionFallback, versionFromBuildInfo(nil, false))
	})

	s.Run("nil info with ok=true → fallback", func() {
		s.Equal(repositoryVersionFallback, versionFromBuildInfo(nil, true))
	})

	s.Run("main module is a tagged release → its version", func() {
		info := &debug.BuildInfo{Main: debug.Module{Path: modulePath, Version: "v1.5.0"}}
		s.Equal("v1.5.0", versionFromBuildInfo(info, true))
	})

	s.Run("consumed as a tagged dependency → dep version", func() {
		info := &debug.BuildInfo{
			Main: debug.Module{Path: "example.com/consumer", Version: "v2.0.0"},
			Deps: []*debug.Module{
				{Path: "example.com/other", Version: "v9.9.9"},
				{Path: modulePath, Version: "v0.4.2"},
			},
		}
		s.Equal("v0.4.2", versionFromBuildInfo(info, true))
	})

	s.Run("untagged main and no matching dep → fallback", func() {
		info := &debug.BuildInfo{
			Main: debug.Module{Path: modulePath, Version: "(devel)"},
			Deps: []*debug.Module{{Path: "example.com/other", Version: "v1.0.0"}},
		}
		s.Equal(repositoryVersionFallback, versionFromBuildInfo(info, true))
	})
}

func (s *VersioningSuite) TestRepositoryVersion_IsResolved() {
	// The package-level var is resolved at init; for this untagged test build it
	// equals the fallback. The assertion pins that the var is wired to the
	// resolver rather than left empty.
	s.Equal(repositoryVersionFallback, RepositoryVersion)
	s.NotEmpty(RepositoryURL)
}
