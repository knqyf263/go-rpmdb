package rpmdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackageInfo_InstalledFileNames(t *testing.T) {
	tests := []struct {
		name    string
		pkg     *PackageInfo
		want    []string
		wantErr bool
	}{
		{
			name: "valid indexes",
			pkg: &PackageInfo{
				Name:       "valid",
				BaseNames:  []string{"a", "b"},
				DirIndexes: []int32{0, 1},
				DirNames:   []string{"/usr/bin/", "/usr/lib/"},
			},
			want: []string{"/usr/bin/a", "/usr/lib/b"},
		},
		{
			name: "dir index out of range (positive)",
			pkg: &PackageInfo{
				Name:       "positive-oob",
				BaseNames:  []string{"a"},
				DirIndexes: []int32{5},
				DirNames:   []string{"/usr/bin/"},
			},
			wantErr: true,
		},
		{
			name: "negative dir index from crafted header",
			pkg: &PackageInfo{
				Name:       "negative-oob",
				BaseNames:  []string{"a"},
				DirIndexes: []int32{-1},
				DirNames:   []string{"/usr/bin/"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pkg.InstalledFileNames()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
