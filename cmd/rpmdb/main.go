package main

import (
	"fmt"
	"log"

	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
func run() error {
	db := rpmdb.DB{}
	err := db.Open("./Packages")
	if err != nil {
		return err
	}
	pkgList, err := db.ListPackages()
	if err != nil {
		return err
	}

	for _, pkg := range pkgList {
		fmt.Printf("%+v\n", *pkg)
	}
	return nil
}
