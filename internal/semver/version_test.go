// Copyright (c) 2024-2026 Inetum Poland.

package semver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Inetum-Poland/action-bumper/internal/config"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"with v prefix", "v1.2.3", "1.2.3", false},
		{"without v prefix", "1.2.3", "1.2.3", false},
		{"major only", "v1.0.0", "1.0.0", false},
		{"with prerelease", "v1.2.3-alpha.1", "1.2.3-alpha.1", false},
		{"with build metadata", "v1.2.3+build.123", "1.2.3+build.123", false},
		{"invalid format", "invalid", "", true},
		{"empty string", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func TestDefaultVersion(t *testing.T) {
	// DefaultVersion always returns 0.0.0; Bump() produces the correct first version
	assert.Equal(t, "0.0.0", DefaultVersion().String())

	tests := []struct {
		level      config.BumpLevel
		wantBumped string
	}{
		{config.BumpLevelMajor, "1.0.0"},
		{config.BumpLevelMinor, "0.1.0"},
		{config.BumpLevelPatch, "0.0.1"},
		{config.BumpLevelNone, "0.0.0"},
		{config.BumpLevelEmpty, "0.0.0"},
	}

	for _, tt := range tests {
		t.Run(string(tt.level), func(t *testing.T) {
			bumped := DefaultVersion().Bump(tt.level)
			assert.Equal(t, tt.wantBumped, bumped.String())
		})
	}
}

func TestVersion_Bump(t *testing.T) {
	tests := []struct {
		name    string
		current string
		level   config.BumpLevel
		want    string
	}{
		{"bump major from 1.2.3", "1.2.3", config.BumpLevelMajor, "2.0.0"},
		{"bump major from 0.1.0", "0.1.0", config.BumpLevelMajor, "1.0.0"},
		{"bump minor from 1.2.3", "1.2.3", config.BumpLevelMinor, "1.3.0"},
		{"bump minor from 1.0.0", "1.0.0", config.BumpLevelMinor, "1.1.0"},
		{"bump patch from 1.2.3", "1.2.3", config.BumpLevelPatch, "1.2.4"},
		{"bump patch from 1.0.0", "1.0.0", config.BumpLevelPatch, "1.0.1"},
		{"no bump", "1.2.3", config.BumpLevelNone, "1.2.3"},
		{"empty level", "1.2.3", config.BumpLevelEmpty, "1.2.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := MustParse(tt.current)
			got := v.Bump(tt.level)
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func TestVersion_StringFormats(t *testing.T) {
	v := MustParse("1.2.3")

	assert.Equal(t, "1.2.3", v.String())
	assert.Equal(t, "v1.2.3", v.StringWithV())
}

func TestVersion_Components(t *testing.T) {
	v := MustParse("1.2.3")

	assert.Equal(t, uint64(1), v.Major())
	assert.Equal(t, uint64(2), v.Minor())
	assert.Equal(t, uint64(3), v.Patch())
}

func TestVersion_Tags(t *testing.T) {
	v := MustParse("1.2.3")

	tests := []struct {
		name     string
		method   func(bool) string
		includeV bool
		want     string
	}{
		{"major tag with v", v.MajorTag, true, "v1"},
		{"major tag without v", v.MajorTag, false, "1"},
		{"minor tag with v", v.MinorTag, true, "v1.2"},
		{"minor tag without v", v.MinorTag, false, "1.2"},
		{"full tag with v", v.FullTag, true, "v1.2.3"},
		{"full tag without v", v.FullTag, false, "1.2.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.method(tt.includeV)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		{"equal versions", "1.2.3", "1.2.3", 0},
		{"v1 < v2", "1.2.3", "1.2.4", -1},
		{"v1 > v2", "1.2.4", "1.2.3", 1},
		{"major difference", "2.0.0", "1.9.9", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := MustParse(tt.v1)
			v2 := MustParse(tt.v2)
			got := v1.Compare(v2)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVersion_Equal(t *testing.T) {
	v1 := MustParse("1.2.3")
	v2 := MustParse("1.2.3")
	v3 := MustParse("1.2.4")

	assert.True(t, v1.Equal(v2))
	assert.False(t, v1.Equal(v3))
}

func TestVersion_LessThan(t *testing.T) {
	v1 := MustParse("1.2.3")
	v2 := MustParse("1.2.4")

	assert.True(t, v1.LessThan(v2))
	assert.False(t, v2.LessThan(v1))
	assert.False(t, v1.LessThan(v1))
}

func TestVersion_GreaterThan(t *testing.T) {
	v1 := MustParse("1.2.4")
	v2 := MustParse("1.2.3")

	assert.True(t, v1.GreaterThan(v2))
	assert.False(t, v2.GreaterThan(v1))
	assert.False(t, v1.GreaterThan(v1))
}

func TestMustParse_Panic(t *testing.T) {
	assert.Panics(t, func() {
		MustParse("invalid")
	})
}

func TestVersion_SortManyVersions(t *testing.T) {
	// Test sorting 100+ versions correctly
	versions := make([]*Version, 0, 120)

	// Add versions from v1.0.0 to v10.10.10
	for major := 1; major <= 10; major++ {
		for minor := 0; minor <= 10; minor++ {
			v := MustParse(fmt.Sprintf("v%d.%d.0", major, minor))
			versions = append(versions, v)
		}
	}

	// Find max using comparison
	var maxVersion *Version
	for _, v := range versions {
		if maxVersion == nil || v.GreaterThan(maxVersion) {
			maxVersion = v
		}
	}

	assert.Equal(t, "10.10.0", maxVersion.String())
}

func TestVersion_NumericSorting(t *testing.T) {
	// Ensure v1.10.0 > v1.9.0 (numeric, not lexicographic)
	tests := []struct {
		v1, v2    string
		v1Greater bool
	}{
		{"v1.10.0", "v1.9.0", true},
		{"v1.9.0", "v1.10.0", false},
		{"v2.0.0", "v1.99.99", true},
		{"v10.0.0", "v9.99.99", true},
		{"v1.2.10", "v1.2.9", true},
	}

	for _, tt := range tests {
		t.Run(tt.v1+"_vs_"+tt.v2, func(t *testing.T) {
			v1 := MustParse(tt.v1)
			v2 := MustParse(tt.v2)
			assert.Equal(t, tt.v1Greater, v1.GreaterThan(v2))
		})
	}
}

func TestVersion_PrereleaseHandling(t *testing.T) {
	tests := []struct {
		name    string
		version string
		valid   bool
	}{
		{"alpha prerelease", "v1.0.0-alpha.1", true},
		{"beta prerelease", "v1.0.0-beta.2", true},
		{"rc prerelease", "v1.0.0-rc.1", true},
		{"simple prerelease", "v1.0.0-pre", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := Parse(tt.version)
			if tt.valid {
				assert.NoError(t, err)
				assert.NotNil(t, v)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestVersion_BuildMetadataHandling(t *testing.T) {
	tests := []struct {
		name    string
		version string
		valid   bool
	}{
		{"simple build", "v1.0.0+build.123", true},
		{"sha build", "v1.0.0+20130313144700", true},
		{"complex build", "v1.0.0+build.11.e0f985a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := Parse(tt.version)
			if tt.valid {
				assert.NoError(t, err)
				assert.NotNil(t, v)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestDefaultVersion_ReturnsZero(t *testing.T) {
	assert.Equal(t, "0.0.0", DefaultVersion().String())
}

func TestVersion_FirstBumpProducesCorrectVersion(t *testing.T) {
	// Verify that bumping from 0.0.0 produces correct first versions
	tests := []struct {
		level    config.BumpLevel
		expected string
	}{
		{config.BumpLevelMajor, "1.0.0"},
		{config.BumpLevelMinor, "0.1.0"},
		{config.BumpLevelPatch, "0.0.1"},
	}

	for _, tt := range tests {
		t.Run(string(tt.level), func(t *testing.T) {
			bumped := DefaultVersion().Bump(tt.level)
			assert.Equal(t, tt.expected, bumped.String())
		})
	}
}
