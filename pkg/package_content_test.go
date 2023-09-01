package rpmdb

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/glebarez/go-sqlite"
)

type packageContentTestCase struct {
	name                       string
	pkgName                    string
	testDir                    string
	databaseFile               string
	wantPackageFile            string
	wantInstalledFilesFile     string
	wantInstalledFileNamesFile string
}

func TestPackageContent(t *testing.T) {
	tests := []packageContentTestCase{
		{
			name:                       "fedora35 python3",
			pkgName:                    "python3",
			testDir:                    "testdata/fedora35",
			databaseFile:               "rpmdb.sqlite",
			wantPackageFile:            "python3.json",
			wantInstalledFilesFile:     "python3_files.json",
			wantInstalledFileNamesFile: "python3_file_names.json",
		},
		{
			name:                       "centos5 python",
			pkgName:                    "python",
			testDir:                    "testdata/centos5-plain",
			databaseFile:               "Packages",
			wantPackageFile:            "python.json",
			wantInstalledFilesFile:     "python_files.json",
			wantInstalledFileNamesFile: "python_file_names.json",
		},
		{
			name:                       "centos6 glibc",
			pkgName:                    "glibc",
			testDir:                    "testdata/centos6-plain",
			databaseFile:               "Packages",
			wantPackageFile:            "glibc.json",
			wantInstalledFilesFile:     "glibc_files.json",
			wantInstalledFileNamesFile: "glibc_file_names.json",
		},
		{
			name:                       "centos8 nodejs",
			pkgName:                    "nodejs",
			testDir:                    "testdata/fedora38-modules",
			databaseFile:               "rpmdb.sqlite",
			wantPackageFile:            "nodejs.json",
			wantInstalledFilesFile:     "nodejs_files.json",
			wantInstalledFileNamesFile: "nodejs_file_names.json",
		},
		{
			name:                       "CBL-Mariner 2.0 curl",
			pkgName:                    "curl",
			testDir:                    "testdata/cbl-mariner-2.0",
			databaseFile:               "rpmdb.sqlite",
			wantPackageFile:            "curl.json",
			wantInstalledFilesFile:     "curl_files.json",
			wantInstalledFileNamesFile: "curl_file_names.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := Open(filepath.Join(tt.testDir, tt.databaseFile))
			require.NoError(t, err)

			got, err := db.Package(tt.pkgName)
			assert.NoError(t, err)

			gotInstalledFiles, err := got.InstalledFiles()
			assert.NoError(t, err)

			gotInstalledFileNames, err := got.InstalledFileNames()
			assert.NoError(t, err)

			if *update && tt.wantPackageFile != "" {
				err = updatePackageFile(t, tt, got)
				assert.NoError(t, err)
				err = updateInstalledFiles(t, tt, gotInstalledFiles)
				assert.NoError(t, err)
				err = updateInstalledFileNames(t, tt, gotInstalledFileNames)
				assert.NoError(t, err)
			}

			want, err := readPackageFile(t, tt)
			assert.NoError(t, err)

			wantInstalledFiles, err := readInstalledFiles(t, tt)
			assert.NoError(t, err)

			wantInstalledFileNames, err := readInstalledFileNames(t, tt)
			assert.NoError(t, err)

			assert.Equal(t, wantInstalledFiles, gotInstalledFiles)
			assert.Equal(t, wantInstalledFileNames, gotInstalledFileNames)
			assert.Equal(t, want, got)
		})
	}
}
