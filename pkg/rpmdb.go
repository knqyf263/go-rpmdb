package rpmdb

import (
	"github.com/knqyf263/go-rpmdb/pkg/bdb"
	"golang.org/x/xerrors"
)

type RpmDB struct {
	db *bdb.BerkeleyDB
}

func Open(path string) (*RpmDB, error) {
	db, err := bdb.Open(path)
	if err != nil {
		return nil, err
	}

	return &RpmDB{
		db: db,
	}, nil

}

func (d *RpmDB) Package(name string) (*PackageInfo, error) {
	pkgs, err := d.ListPackages()
	if err != nil {
		return nil, xerrors.Errorf("unable to list packages: %w", err)
	}

	for _, pkg := range pkgs {
		if pkg.Name == name {
			return pkg, nil
		}
	}
	return nil, xerrors.Errorf("%s is not installed", name)
}

func (d *RpmDB) ListPackages() ([]*PackageInfo, error) {
	var pkgList []*PackageInfo

	for entry := range d.db.Read() {
		if entry.Err != nil {
			return nil, entry.Err
		}

		indexEntries, err := headerImport(entry.Value)
		if err != nil {
			return nil, xerrors.Errorf("error during importing header: %w", err)
		}
		pkg, err := getNEVRA(indexEntries)
		if err != nil {
			return nil, xerrors.Errorf("invalid package info: %w", err)
		}
		pkgList = append(pkgList, pkg)
	}

	return pkgList, nil
}
