// Copyright (c) 2024 Inetum Poland.

package semver

import (
	"testing"

	"github.com/Inetum-Poland/action-bumper/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	tests := []struct {
		level config.BumpLevel
		want  string
	}{
		{config.BumpLevelMajor, "1.0.0"},
		{config.BumpLevelMinor, "0.1.0"},
		{config.BumpLevelPatch, "0.0.1"},
		{config.BumpLevelNone, "0.0.0"},
		{config.BumpLevelEmpty, "0.0.0"},
	}

	for _, tt := range tests {
		t.Run(string(tt.level), func(t *testing.T) {
			got := DefaultVersion(tt.level)
			assert.Equal(t, tt.want, got.String())
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
