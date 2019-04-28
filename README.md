# go-rpmdb
RPM DB bindings for go

## Feature
- Extract installed rpm packages

## Example

Locate `Packages` in the same directory
```
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
            // {Epoch:0 Name:m4 Version:1.4.16 Release:10.el7 Arch:x86_64}
            // {Epoch:0 Name:zip Version:3.0 Release:11.el7 Arch:x86_64}
            // ...

	}
	return nil
}
```