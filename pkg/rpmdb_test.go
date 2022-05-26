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
		{
			name:    "SLE15 with NDB style rpm database",
			file:    "testdata/sle15-bci/Packages.db",
			pkgList: SLE15WithNDB,
		},
		{
			name:    "Fedora35 with SQLite3 style rpm database",
			file:    "testdata/fedora35/rpmdb.sqlite",
			pkgList: Fedora35WithSQLite3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Open(tt.file)
			require.NoError(t, err)

			got, err := db.ListPackages()
			require.NoError(t, err)

			// They are tested in another function.
			for _, g := range got {
				g.BaseNames = nil
				g.DirIndexes = nil
				g.DirNames = nil
				g.FileSizes = nil
				g.FileDigests = nil
				g.FileModes = nil
				g.FileFlags = nil
				g.UserNames = nil
				g.GroupNames = nil
			}

			for i, p := range tt.pkgList {
				assert.Equal(t, p, got[i])
			}
		})
	}
}

func TestRpmDB_Package(t *testing.T) {
	tests := []struct {
		name                   string
		pkgName                string
		file                   string // Test input file
		want                   *PackageInfo
		wantInstalledFiles     []FileInfo
		wantInstalledFileNames []string
		wantErr                string
	}{
		{
			name:    "centos5 python",
			pkgName: "python",
			file:    "testdata/centos5-plain/Packages",
			want: &PackageInfo{
				Name:      "python",
				Version:   "2.4.3",
				Release:   "56.el5",
				Arch:      "x86_64",
				Size:      74377,
				SourceRpm: "python-2.4.3-56.el5.src.rpm",
				License:   "PSF - see LICENSE",
				Vendor:    "CentOS",
			},
			wantInstalledFiles:     CentOS5PythonInstalledFiles,
			wantInstalledFileNames: CentOS5PythonInstalledFileNames,
		},
		{
			name:    "centos6 glibc",
			pkgName: "glibc",
			file:    "testdata/centos6-plain/Packages",
			want: &PackageInfo{
				Name:            "glibc",
				Version:         "2.12",
				Release:         "1.212.el6",
				Arch:            "x86_64",
				Size:            13117447,
				SourceRpm:       "glibc-2.12-1.212.el6.src.rpm",
				License:         "LGPLv2+ and LGPLv2+ with exceptions and GPLv2+",
				Vendor:          "CentOS",
				DigestAlgorithm: PGPHASHALGO_SHA256,
			},
			wantInstalledFiles:     CentOS6GlibcInstalledFiles,
			wantInstalledFileNames: CentOS6GlibcInstalledFileNames,
		},
		{
			name:    "centos8 nodejs",
			pkgName: "nodejs",
			file:    "testdata/centos8-modularitylabel/Packages",
			want: &PackageInfo{
				Epoch:           intRef(1),
				Name:            "nodejs",
				Version:         "10.21.0",
				Release:         "3.module_el8.2.0+391+8da3adc6",
				Arch:            "x86_64",
				Size:            31483781,
				SourceRpm:       "nodejs-10.21.0-3.module_el8.2.0+391+8da3adc6.src.rpm",
				License:         "MIT and ASL 2.0 and ISC and BSD",
				Vendor:          "CentOS",
				Modularitylabel: "nodejs:10:8020020200707141642:6a468ee4",
				DigestAlgorithm: PGPHASHALGO_SHA256,
			},
			wantInstalledFiles:     CentOS8NodejsInstalledFiles,
			wantInstalledFileNames: CentOS8NodejsInstalledFileNames,
		},
		{
			name:    "CBL-Mariner 2.0 curl",
			pkgName: "curl",
			file:    "testdata/cbl-mariner-2.0/rpmdb.sqlite",
			want: &PackageInfo{
				Name:            "curl",
				Version:         "7.76.0",
				Release:         "6.cm2",
				Arch:            "x86_64",
				Size:            326023,
				SourceRpm:       "curl-7.76.0-6.cm2.src.rpm",
				License:         "MIT",
				Vendor:          "Microsoft Corporation",
				DigestAlgorithm: PGPHASHALGO_SHA256,
			},
			wantInstalledFiles:     Mariner2CurlInstalledFiles,
			wantInstalledFileNames: Mariner2CurlInstalledFileNames,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Open(tt.file)
			require.NoError(t, err)

			got, err := db.Package(tt.pkgName)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			assert.NoError(t, err)

			gotInstalledFiles, err := got.InstalledFiles()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantInstalledFiles, gotInstalledFiles)

			gotInstalledFileNames, err := got.InstalledFileNames()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantInstalledFileNames, gotInstalledFileNames)

			// These fields are tested through InstalledFiles() above
			got.BaseNames = nil
			got.DirIndexes = nil
			got.DirNames = nil
			got.FileSizes = nil
			got.FileDigests = nil
			got.FileModes = nil
			got.FileFlags = nil
			got.UserNames = nil
			got.GroupNames = nil

			assert.Equal(t, tt.want, got)
		})
	}
}
