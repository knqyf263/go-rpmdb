#!/bin/bash
set -e

getListing() {
	docker run -it --rm -v $1:/testdata/ $2 rpm --dbpath /testdata/ -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{SOURCERPM}\", %{SIZE}, \"%{LICENSE}\", \"%{VENDOR}\", \"%{SUMMARY}\", %{INSTALLTIME}, \"\", nil, nil, nil\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
}

getListingCentos8() {
	docker run -it --rm -v $1:/testdata/ $2 rpm --dbpath /testdata/ -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{SOURCERPM}\", %{SIZE}, \"%{LICENSE}\", \"%{VENDOR}\", \"%{SUMMARY}\", %{INSTALLTIME}, \"%{MODULARITYLABEL}\", nil, nil, nil\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
}

getFilesCommand() {
	rpm -ql $1 | awk '{printf "\"%s\",\n", $1}'
}

getFilesCommand2() {
	rpm -ql --dbpath $1 $2 | awk '{printf "\"%s\",\n", $1}'
}

getFiles() {
	docker run -it -v $(pwd):/test --rm $1 /bin/bash -c "/test/gen_testcase.sh getFilesCommand2 /test/testdata/$2 $3"
}

getFiles2() {
	docker run -it -v `pwd`:/test --rm $1 /bin/bash -c "/test/gen_testcase.sh getFilesCommand $2"
}

genTestcases(){
cat << EOT
package rpmdb

var (
	// docker run --rm -it centos:5 bash
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{SOURCERPM}\", %{SIZE}, \"%{LICENSE}\", \"%{VENDOR}\", \'%{SUMMARY}\', %{INSTALLTIME}, \"\", nil, nil, nil\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS5Plain = []*PackageInfo{
		`getListing "$(pwd)/testdata/centos5-plain" centos:5`
	}

	// docker run --rm -it centos:6 bash
	// yum groupinstall -y "Development tools"
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{SOURCERPM}\", %{SIZE}, \"%{LICENSE}\", \"%{VENDOR}\", \'%{SUMMARY}\', %{INSTALLTIME}, \"\", nil, nil, nil\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS6DevTools = []*PackageInfo{
		`getListing "$(pwd)/testdata/centos6-devtools" centos:6`
	}

	// docker run --rm -it centos:6 bash
	// yum groupinstall -y "Development tools"
	// yum install -y rpm-build redhat-rpm-config asciidoc hmaccalc perl-ExtUtils-Embed pesign xmlto
	// yum install -y audit-libs-devel binutils-devel elfutils-devel elfutils-libelf-devel java-devel
	// yum install -y ncurses-devel newt-devel numactl-devel pciutils-devel python-devel zlib-devel
	// yum install -y net-tools bc
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{SOURCERPM}\", %{SIZE}, \"%{LICENSE}\", \"%{VENDOR}\", \'%{SUMMARY}\', %{INSTALLTIME}, \"\", nil, nil, nil\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS6Many = []*PackageInfo{
		`getListing "$(pwd)/testdata/centos6-many" centos:6`
	}

	// docker run --rm -it centos:6 bash
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{SOURCERPM}\", %{SIZE}, \"%{LICENSE}\", \"%{VENDOR}\", \'%{SUMMARY}\', %{INSTALLTIME}, \"\", nil, nil, nil\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS6Plain = []*PackageInfo{
		`getListing "$(pwd)/testdata/centos6-plain" centos:6`
	}

	// docker run --rm -it centos:7 bash
	// yum groupinstall -y "Development tools"
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{SOURCERPM}\", %{SIZE}, \"%{LICENSE}\", \"%{VENDOR}\", \'%{SUMMARY}\', %{INSTALLTIME}, \"\", nil, nil, nil\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS7DevTools = []*PackageInfo{
		`getListing "$(pwd)/testdata/centos7-devtools" centos:7`
	}

	// docker run --rm -it centos/httpd-24-centos7 bash
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{SOURCERPM}\", %{SIZE}, \"%{LICENSE}\", \"%{VENDOR}\", \'%{SUMMARY}\', %{INSTALLTIME}, \"\", nil, nil, nil\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS7Httpd24 = []*PackageInfo{
		`getListing "$(pwd)/testdata/centos7-httpd24" centos:7`
	}

	// docker run --rm -it centos:7 bash
	// yum groupinstall -y "Development tools"
	// yum install -y rpm-build redhat-rpm-config asciidoc hmaccalc perl-ExtUtils-Embed pesign xmlto
	// yum install -y audit-libs-devel binutils-devel elfutils-devel elfutils-libelf-devel java-devel
	// yum install -y ncurses-devel newt-devel numactl-devel pciutils-devel python-devel zlib-devel
	// yum install -y net-tools bc
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{SOURCERPM}\", %{SIZE}, \"%{LICENSE}\", \"%{VENDOR}\", \'%{SUMMARY}\', %{INSTALLTIME}, \"\", nil, nil, nil\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS7Many = []*PackageInfo{
		`getListing "$(pwd)/testdata/centos7-many" centos:7`
	}

	// docker run --rm -it centos:7 bash
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{SOURCERPM}\", %{SIZE}, \"%{LICENSE}\", \"%{VENDOR}\", \'%{SUMMARY}\', %{INSTALLTIME}, \"\", nil, nil, nil\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS7Plain = []*PackageInfo{
		`getListing "$(pwd)/testdata/centos7-plain" centos:7`
	}

	// docker run --rm -it centos/python-35-centos7 bash
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{SOURCERPM}\", %{SIZE}, \"%{LICENSE}\", \"%{VENDOR}\", \'%{SUMMARY}\', %{INSTALLTIME}, \"\", nil, nil, nil\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS7Python35 = []*PackageInfo{
		`getListing "$(pwd)/testdata/centos7-python35" centos:7`
	}

	// docker run --rm -it centos:8 bash
	// yum module install -y container-tools
	// yum groupinstall -y "Development tools"
	// yum -y install nodejs podman-docker
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{SOURCERPM}\", %{SIZE}, \"%{LICENSE}\", \"%{VENDOR}\", \"%{SUMMARY}\", %{INSTALLTIME}, \"%{MODULARITYLABEL}\", nil, nil, nil\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	CentOS8Modularitylabel = []*PackageInfo{
		`getListingCentos8 "$(pwd)/testdata/centos8-modularitylabel" centos:8`
	}

	// docker run --rm -it registry.suse.com/bci/bci-minimal:15.3 bash
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{SOURCERPM}\", %{SIZE}, \"%{LICENSE}\", \"%{VENDOR}\", \'%{SUMMARY}\', %{INSTALLTIME}, \"\", nil, nil, nil\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	SLE15WithNDB = []*PackageInfo{
		`getListing "$(pwd)/testdata/sle15-bci" registry.suse.com/bci/bci-minimal:15.3`
	}

	// docker run --rm -it fedora:35 bash
	// rpm -qa --queryformat "\{%{EPOCH}, \"%{NAME}\", \"%{VERSION}\", \"%{RELEASE}\", \"%{ARCH}\", \"%{SOURCERPM}\", %{SIZE}, \"%{LICENSE}\", \"%{VENDOR}\", \'%{SUMMARY}\', %{INSTALLTIME}, \"\", nil, nil, nil\},\n" | sed "s/^{(none)/{0/g" | sed "s/(none)//g"
	Fedora35WithSQLite3 = []*PackageInfo{
		`getListing "$(pwd)/testdata/fedora35" fedora:35`
	}

	// rpm -ql python --dbpath /path/to/testdata/centos5-plain | awk '{printf "\"%s\",\n", $1}'
	CentOS5PythonInstalledFiles = []string{
		`getFiles2 centos:5 python`
	}

	// rpm -ql glibc --dbpath /path/to/testdata/centos6-plain | awk '{printf "\"%s\",\n", $1}'
	CentOS6GlibcInstalledFiles = []string{
		`getFiles2 centos:6 glibc`
	}

	CentOS8NodejsInstalledFiles = []string{
		`getFiles centos:8 centos8-modularitylabel nodejs`
	}

	Mariner2CurlInstalledFiles = []string{
		`getFiles fedora:35 cbl-mariner-2.0 curl`
	}
)
EOT
}

if [ $# -ge 1 ]; then
	case "$1" in
	"getFilesCommand")
		getFilesCommand $2
		;;
	"getFilesCommand2")
		getFilesCommand2 $2 $3
		;;
	"genTestcases")
		genTestcases
		;;
	*)
		echo "This Operation is not support!"
		;;
	esac

else
	echo "This Operation is not support!"
fi
