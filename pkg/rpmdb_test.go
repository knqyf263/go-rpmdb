package rpmdb

import (
	"path"
	"testing"
)

func TestPackageList(t *testing.T) {
	vectors := []struct {
		file    string // Test input file
		pkgList []PackageInfo
	}{
		{
			file:    "testdata/Packages_centos6_plain",
			pkgList: CentOS6Plain,
		},
		{
			file:    "testdata/Packages_centos6_dev_tools",
			pkgList: CentOS6DevTools,
		},
		{
			file:    "testdata/Packages_centos6_many",
			pkgList: CentOS6Many,
		},
		{
			file:    "testdata/Packages_centos7_plain",
			pkgList: CentOS7Plain,
		},
		{
			file:    "testdata/Packages_centos7_dev_tools",
			pkgList: CentOS7DevTools,
		},
		{
			file:    "testdata/Packages_centos7_many",
			pkgList: CentOS7Many,
		},
		{
			file:    "testdata/Packages_centos7_python35",
			pkgList: CentOS7Python35,
		},
		{
			file:    "testdata/Packages_centos7_httpd24",
			pkgList: CentOS7Httpd24,
		},
	}

	for _, v := range vectors {
		t.Run(path.Base(v.file), func(t *testing.T) {
			db, err := Open(v.file)
			if err != nil {
				t.Fatalf("Open() error: %v", err)
			}
			pkgList, err := db.ListPackages()
			if err != nil {
				t.Fatalf("ListPackagges() error: %v", err)
			}

			if len(pkgList) != len(v.pkgList) {
				t.Errorf("pkg length: got %v, want %v", len(pkgList), len(v.pkgList))
			}

			for i, got := range pkgList {
				want := v.pkgList[i]
				if want.Epoch != got.Epoch {
					t.Errorf("%d: Epoch: got %d, want %d", i, got.Epoch, want.Epoch)
				}
				if want.Name != got.Name {
					t.Errorf("%d: Name: got %s, want %s", i, got.Name, want.Name)
				}
				if want.Version != got.Version {
					t.Errorf("%d: Version: got %s, want %s", i, got.Version, want.Version)
				}
				if want.Release != got.Release {
					t.Errorf("%d: Release: got %s, want %s", i, got.Release, want.Release)
				}
				if want.Arch != got.Arch {
					t.Errorf("%d: Arch: got %s, want %s", i, got.Arch, want.Arch)
				}
			}
		})
	}
}
