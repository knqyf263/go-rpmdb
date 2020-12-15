package rpmdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackageList(t *testing.T) {
	tests := []struct {
		name    string
		file    string // Test input file
		pkgList []*PackageInfo
	}{
		{
			name:    "CentOS5 plain",
			file:    "testdata/centos5-plain/Packages",
			pkgList: CentOS5Plain,
		},
		{
			name:    "CentOS6 Plain",
			file:    "testdata/centos6-plain/Packages",
			pkgList: CentOS6Plain,
		},
		{
			name:    "CentOS6 with Development tools",
			file:    "testdata/centos6-devtools/Packages",
			pkgList: CentOS6DevTools,
		},
		{
			name:    "CentOS6 with many packages",
			file:    "testdata/centos6-many/Packages",
			pkgList: CentOS6Many,
		},
		{
			name:    "CentOS7 Plain",
			file:    "testdata/centos7-plain/Packages",
			pkgList: CentOS7Plain,
		},
		{
			name:    "CentOS7 with Development tools",
			file:    "testdata/centos7-devtools/Packages",
			pkgList: CentOS7DevTools,
		},
		{
			name:    "CentOS7 with many packages",
			file:    "testdata/centos7-many/Packages",
			pkgList: CentOS7Many,
		},
		{
			name:    "CentOS7 with Python 3.5",
			file:    "testdata/centos7-python35/Packages",
			pkgList: CentOS7Python35,
		},
		{
			name:    "CentOS7 with httpd 2.4",
			file:    "testdata/centos7-httpd24/Packages",
			pkgList: CentOS7Httpd24,
		},
		{
			name:    "CentOS8 with modules",
			file:    "testdata/centos8-modularitylabel/Packages",
			pkgList: CentOS8Modularitylabel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Open(tt.file)
			require.NoError(t, err)

			got, err := db.ListPackages()
			require.NoError(t, err)

			assert.ElementsMatch(t, tt.pkgList, got)
		})
	}
}
