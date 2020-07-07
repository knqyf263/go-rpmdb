#!/bin/bash

getListing() {
	_pkgdbdir="`pwd`/pkg.$$"
	rm -rf "$_pkgdbdir"
	mkdir "$_pkgdbdir"
	cp "$1" "$_pkgdbdir/Packages"
	rpm -qa --dbpath "$_pkgdbdir" --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{VENDOR}\", \`%{SUMMARY}\`, %{SIZE}, %{INSTALLTIME}\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g" | sed "s/^/		/g"
	rm -rf "$_pkgdbdir"
}

cat << EOT
package rpmdb

var (
	// docker run --rm -it centos:6 bash
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{VENDOR}\", \"%{SUMMARY}\", %{SIZE}, %{INSTALLTIME}\},\n" | sed "s/(none)/0/g"
	CentOS6Plain = []PackageInfo{
`getListing "Packages_centos6_plain"`
	}

	// docker run --rm -it centos:6 bash
	// yum groupinstall -y "Development tools"
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{VENDOR}\", \"%{SUMMARY}\", %{SIZE}, %{INSTALLTIME}\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS6DevTools = []PackageInfo{
`getListing "Packages_centos6_dev_tools"`
	}

	// docker run --rm -it centos:6 bash
	// yum groupinstall -y "Development tools"
	// yum install -y rpm-build redhat-rpm-config asciidoc hmaccalc perl-ExtUtils-Embed pesign xmlto
	// yum install -y audit-libs-devel binutils-devel elfutils-devel elfutils-libelf-devel java-devel
	// yum install -y ncurses-devel newt-devel numactl-devel pciutils-devel python-devel zlib-devel
	// yum install -y net-tools bc
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{VENDOR}\", \"%{SUMMARY}\", %{SIZE}, %{INSTALLTIME}\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS6Many = []PackageInfo{
`getListing "Packages_centos6_many"`
	}

	// docker run --rm -it centos:7 bash
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{VENDOR}\", \"%{SUMMARY}\", %{SIZE}, %{INSTALLTIME}\},\n" | sed "s/(none)/0/g"
	CentOS7Plain = []PackageInfo{
`getListing "Packages_centos7_plain"`
	}

	// docker run --rm -it centos:7 bash
	// yum groupinstall -y "Development tools"
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{VENDOR}\", \"%{SUMMARY}\", %{SIZE}, %{INSTALLTIME}\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS7DevTools = []PackageInfo{
`getListing "Packages_centos7_dev_tools"`
	}

	// docker run --rm -it centos:7 bash
	// yum groupinstall -y "Development tools"
	// yum install -y rpm-build redhat-rpm-config asciidoc hmaccalc perl-ExtUtils-Embed pesign xmlto
	// yum install -y audit-libs-devel binutils-devel elfutils-devel elfutils-libelf-devel java-devel
	// yum install -y ncurses-devel newt-devel numactl-devel pciutils-devel python-devel zlib-devel
	// yum install -y net-tools bc
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{VENDOR}\", \"%{SUMMARY}\", %{SIZE}, %{INSTALLTIME}\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS7Many = []PackageInfo{
`getListing "Packages_centos7_many"`
	}

	// docker run --rm -it centos/python-35-centos7 bash
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{VENDOR}\", \"%{SUMMARY}\", %{SIZE}, %{INSTALLTIME}\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS7Python35 = []PackageInfo{
`getListing "Packages_centos7_python35"`
	}

	// docker run --rm -it centos/httpd-24-centos7 bash
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{VENDOR}\", \"%{SUMMARY}\", %{SIZE}, %{INSTALLTIME}\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS7Httpd24 = []PackageInfo{
`getListing "Packages_centos7_httpd24"`
	}
)
EOT

