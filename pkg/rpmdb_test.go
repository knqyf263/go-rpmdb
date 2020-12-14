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
			file:    "testdata/centos6-plain/Packages",
			pkgList: CentOS6Plain,
		},
		{
			file:    "testdata/centos6-devtools/Packages",
			pkgList: CentOS6DevTools,
		},
		{
			file:    "testdata/centos6-many/Packages",
			pkgList: CentOS6Many,
		},
		{
			file:    "testdata/centos7-plain/Packages",
			pkgList: CentOS7Plain,
		},
		{
			file:    "testdata/centos7-devtools/Packages",
			pkgList: CentOS7DevTools,
		},
		{
			file:    "testdata/centos7-many/Packages",
			pkgList: CentOS7Many,
		},
		{
			file:    "testdata/centos7-python35/Packages",
			pkgList: CentOS7Python35,
		},
		{
			file:    "testdata/centos7-httpd24/Packages",
			pkgList: CentOS7Httpd24,
		},
		{
			file:    "testdata/centos8-modularitylabel/Packages",
			pkgList: CentOS8Modularitylabel,
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
				if want.SourceRpm != got.SourceRpm {
					t.Errorf("%d: SourceRpm: got %s, want %s", i, got.SourceRpm, want.SourceRpm)
				}
				if want.Vendor != got.Vendor {
					t.Errorf("%d: Vendor: got %s, want %s", i, got.Vendor, want.Vendor)
				}
				if want.Size != got.Size {
					t.Errorf("%d: Size: got %d, want %d", i, got.Size, want.Size)
				}
				if want.License != got.License {
					t.Errorf("%d: License: got %s, want %s", i, got.License, want.License)
				}
				if want.Modularitylabel != got.Modularitylabel {
					t.Errorf("%d: Modularitylabel: got %s, want %s", i, got.Modularitylabel, want.Modularitylabel)
				}
			}
		})
	}
}
