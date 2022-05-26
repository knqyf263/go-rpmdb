# go-rpmdb
Library for enumerating packages in an RPM DB `Packages` file (without bindings).

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
	db, err := rpmdb.Open("./Packages")
	if err != nil {
		return err
	}
	pkgList, err := db.ListPackages()
	if err != nil {
		return err
	}

	fmt.Println("Packages:")
	for _, pkg := range pkgList {
		fmt.Printf("\t%+v\n", *pkg)
		// {Epoch:0 Name:m4 Version:1.4.16 Release:10.el7 Arch:x86_64 Vendor:CentOS Summary:The GNU macro processor Size:525707 InstallTime:1556442601}
		// {Epoch:0 Name:zip Version:3.0 Release:11.el7 Arch:x86_64 Vendor:CentOS Summary:A file compression and packaging utility compatible with PKZIP Size:815173 InstallTime:1556442604}
		// ...
	}
	fmt.Printf("[Total Packages: %d]\n", len(pkgList))
	return nil
}
```