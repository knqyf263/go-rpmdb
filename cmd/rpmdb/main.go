package main

import (
	"fmt"
	"log"

	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
)

func main() {
	db, err := rpmdb.Open("./Packages")
	if err != nil {
		log.Fatal(err)
	}
	pkgList, err := db.ListPackages()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Packages:")
	for _, pkg := range pkgList {
		fmt.Printf("\t%+v\n", *pkg)
	}
	fmt.Printf("[Total Packages: %d]\n", len(pkgList))
}

