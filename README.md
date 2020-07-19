# go-rpmdb
Library for enumerating packages in an RPM DB `Packages` file (without bindings).

```
package main

import (
	"fmt"
	"log"

	rpmdb "github.com/wagoodman/go-rpmdb/pkg"
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
		// {Epoch:0 Name:m4 Version:1.4.16 Release:10.el7 Arch:x86_64}
		// {Epoch:0 Name:zip Version:3.0 Release:11.el7 Arch:x86_64}
		// ...
	}
	fmt.Printf("[Total Packages: %d]\n", len(pkgList))
}
```