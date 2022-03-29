package main

import (
	"fmt"
	"log"

	multierror "github.com/hashicorp/go-multierror"
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	db, err := detectDB()
	if err != nil {
		return err
	}
	pkgList, err := db.ListPackages()
	if err != nil {
		return err
	}

	fmt.Println("Packages:")
	for _, pkg := range pkgList {
		// Suppress output
		pkg.Files = nil

		fmt.Printf("\t%+v\n", *pkg)
	}
	fmt.Printf("[Total Packages: %d]\n", len(pkgList))

	return nil
}

func detectDB() (*rpmdb.RpmDB, error) {
	var result error
	db, err := rpmdb.Open("./rpmdb.sqlite")
	if err == nil {
		return db, nil
	}
	result = multierror.Append(result, err)

	db, err = rpmdb.Open("./Packages.db")
	if err == nil {
		return db, nil
	}
	result = multierror.Append(result, err)

	db, err = rpmdb.Open("./Packages")
	if err == nil {
		return db, nil
	}
	result = multierror.Append(result, err)

	return nil, result
}
