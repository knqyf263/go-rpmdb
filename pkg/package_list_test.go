package rpmdb

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/glebarez/go-sqlite"
)

// To update all test fixtures, run:
//
//	go test ./... -update-testcases=true
var update = flag.Bool("update-testcases", false, "update testcases")

type testCase struct {
	name         string
	databaseFile string // Test input file
	testDir      string
	pkgList      []*PackageInfo
	image        string
}

func TestPackageList(t *testing.T) {
	tests := []testCase{
		{
			name:         "amazonlinux2 plain",
			databaseFile: "testdata/amazonlinux2-plain/Packages",
			testDir:      "testdata/amazonlinux2-plain",
			image:        "public.ecr.aws/amazonlinux/amazonlinux:2",
		},
		{
			name:         "amazonlinux2023 plain",
			databaseFile: "testdata/amazonlinux2023-plain/rpmdb.sqlite",
			testDir:      "testdata/amazonlinux2023-plain",
			image:        "public.ecr.aws/amazonlinux/amazonlinux:2023",
		},
		{
			name:    "amazonlinux2023 devtools",
			testDir: "testdata/amazonlinux2023-devtools",
			image:   "public.ecr.aws/amazonlinux/amazonlinux:2023",
		},
		{
			name:    "oraclelinux9 plain",
			testDir: "testdata/oraclelinux9-plain",
			image:   "oraclelinux:9",
		},
		{
			name:         "CentOS5 plain",
			databaseFile: "testdata/centos5-plain/Packages",
			testDir:      "testdata/centos5-plain",
			image:        "centos:5",
		},
		{
			name:         "CentOS6 Plain",
			databaseFile: "testdata/centos6-plain/Packages",
			image:        "centos:6",
			testDir:      "testdata/centos6-plain",
		},
		{
			name:         "CentOS6 with Development tools",
			databaseFile: "testdata/centos6-devtools/Packages",
			image:        "centos:6",
			testDir:      "testdata/centos6-devtools",
		},
		{
			name:         "CentOS6 with many packages",
			databaseFile: "testdata/centos6-many/Packages",
			image:        "centos:6",
			testDir:      "testdata/centos6-many",
		},
		{
			name:         "CentOS7 Plain",
			databaseFile: "testdata/centos7-plain/Packages",
			image:        "centos:7",
			testDir:      "testdata/centos7-plain",
		},
		{
			name:         "CentOS7 with Development tools",
			databaseFile: "testdata/centos7-devtools/Packages",
			image:        "centos:7",
			testDir:      "testdata/centos7-devtools",
		},
		// TODO: Flakey test, not sure why?
		// {
		// 	name:         "CentOS7 with many packages",
		// 	databaseFile: "testdata/centos7-many/Packages",
		// 	image:        "centos:7",
		// 	testDir:      "testdata/centos7-many",
		// },
		{
			name:         "CentOS7 with Python 3.5",
			databaseFile: "testdata/centos7-python35/Packages",
			image:        "centos/python-35-centos7",
			testDir:      "testdata/centos7-python35",
		},
		{
			name:         "CentOS7 with httpd 2.4",
			databaseFile: "testdata/centos7-httpd24/Packages",
			image:        "centos/httpd-24-centos7",
			testDir:      "testdata/centos7-httpd24",
		},

		/*
			// TODO: remove?
			{
				name:         "RHEL UBI8 from s390x",
				databaseFile: "testdata/ubi8-s390x/Packages",
				pkgList:      UBI8s390x(),
			},
		*/
		// TODO: had to manually add m-dashes in liblzma5 summary
		{
			name:         "SLE15 with NDB style rpm database",
			databaseFile: "testdata/sle15-bci/Packages.db",
			image:        "registry.suse.com/bci/bci-minimal:15.3",
			testDir:      "testdata/sle15-bci",
		},
		{
			name:         "Fedora35 with SQLite3 style rpm database",
			databaseFile: "testdata/fedora35/rpmdb.sqlite",
			image:        "fedora:35",
			testDir:      "testdata/fedora35",
		},
		{
			name:         "Fedora35 plus MongoDB with SQLite3 style rpm database",
			databaseFile: "testdata/fedora35-plus-mongo/rpmdb.sqlite",
			image:        "fedora:35",
			testDir:      "testdata/fedora35-plus-mongo",
		},
		{
			name:         "Fedora38 modules",
			databaseFile: "testdata/fedora38-modules/rpmdb.sqlite",
			image:        "fedora:38",
			testDir:      "testdata/fedora38-modules",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if *update && tt.image != "" && tt.pkgList == nil {
				err := updateTestCase(t, tt)
				require.NoError(t, err)
			}

			if tt.testDir != "" {
				var err error
				tt.pkgList, err = readTestCase(t, tt)
				if err != nil {
					require.NoError(t, err)
				}
			}
			// TODO: remove once ubi8-s390x is gone
			if tt.databaseFile == "" {
				tt.databaseFile = filepath.Join(tt.testDir, "rpmdb.sqlite")
			}

			db, err := Open(tt.databaseFile)
			require.NoError(t, err)

			got, err := db.ListPackages()
			require.NoError(t, err)

			// They are tested in another function.
			for _, g := range got {
				g.PGP = ""
				g.DigestAlgorithm = 0
				g.InstallTime = 0
				g.BaseNames = nil
				g.DirIndexes = nil
				g.DirNames = nil
				g.FileSizes = nil
				g.FileDigests = nil
				g.FileModes = nil
				g.FileFlags = nil
				g.UserNames = nil
				g.GroupNames = nil
				g.Provides = nil
				g.Requires = nil
			}

			if len(tt.pkgList) > len(got) {
				t.Fatalf("Got too few packages. Expected: %d, got: %d", len(tt.pkgList), len(got))
			}
			if len(tt.pkgList) < len(got) {
				t.Fatalf("Got too many packages. Expected: %d, got: %d", len(tt.pkgList), len(got))
			}

			for i, p := range tt.pkgList {
				assert.Equal(t, got[i], p)
			}
		})
	}
}
