package rpmdb

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path"
	"testing"

	"github.com/go-test/deep"
)

func TestPackageList(t *testing.T) {
	tests := []struct {
		name    string
		file    string // Test input file
		pkgList []*PackageInfo
	}{
		{
			name:    "CentOS5 plain",
			file:    "testdata/centos5-plain/Packages",
			pkgList: CentOS5Plain,
		},
		{
			name:    "CentOS6 Plain",
			file:    "testdata/centos6-plain/Packages",
			pkgList: CentOS6Plain,
		},
		{
			name:    "CentOS6 with Development tools",
			file:    "testdata/centos6-devtools/Packages",
			pkgList: CentOS6DevTools,
		},
		{
			name:    "CentOS6 with many packages",
			file:    "testdata/centos6-many/Packages",
			pkgList: CentOS6Many,
		},
		{
			name:    "CentOS7 Plain",
			file:    "testdata/centos7-plain/Packages",
			pkgList: CentOS7Plain,
		},
		{
			name:    "CentOS7 with Development tools",
			file:    "testdata/centos7-devtools/Packages",
			pkgList: CentOS7DevTools,
		},
		{
			name:    "CentOS7 with many packages",
			file:    "testdata/centos7-many/Packages",
			pkgList: CentOS7Many,
		},
		{
			name:    "CentOS7 with Python 3.5",
			file:    "testdata/centos7-python35/Packages",
			pkgList: CentOS7Python35,
		},
		{
			name:    "CentOS7 with httpd 2.4",
			file:    "testdata/centos7-httpd24/Packages",
			pkgList: CentOS7Httpd24,
		},
		{
			name:    "CentOS8 with modules",
			file:    "testdata/centos8-modularitylabel/Packages",
			pkgList: CentOS8Modularitylabel,
		},
		{
			name:    "SLE15 with NDB style rpm database",
			file:    "testdata/sle15-bci/Packages.db",
			pkgList: SLE15WithNDB,
		},
		{
			name:    "Fedora35 with SQLite3 style rpm database",
			file:    "testdata/fedora35/rpmdb.sqlite",
			pkgList: Fedora35WithSQLite3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Open(tt.file)
			require.NoError(t, err)

			got, err := db.ListPackages()
			require.NoError(t, err)

			// They are tested in another function.
			for _, g := range got {
				g.Files = nil
			}

			for i, p := range tt.pkgList {
				assert.Equal(t, p, got[i])
			}
		})
	}
}

func TestRpmDB_Package(t *testing.T) {
	tests := []struct {
		name               string
		pkgName            string
		file               string // Test input file
		want               *PackageInfo
		wantInstalledFiles []string
		wantErr            string
	}{
		{
			name:    "centos5 python",
			pkgName: "python",
			file:    "testdata/centos5-plain/Packages",
			want: &PackageInfo{
				Name:      "python",
				Version:   "2.4.3",
				Release:   "56.el5",
				Arch:      "x86_64",
				Size:      74377,
				SourceRpm: "python-2.4.3-56.el5.src.rpm",
				License:   "PSF - see LICENSE",
				Vendor:    "CentOS",
			},
			wantInstalledFiles: CentOS5PythonInstalledFiles,
		},
		{
			name:    "centos6 glibc",
			pkgName: "glibc",
			file:    "testdata/centos6-plain/Packages",
			want: &PackageInfo{
				Name:            "glibc",
				Version:         "2.12",
				Release:         "1.212.el6",
				Arch:            "x86_64",
				Size:            13117447,
				SourceRpm:       "glibc-2.12-1.212.el6.src.rpm",
				License:         "LGPLv2+ and LGPLv2+ with exceptions and GPLv2+",
				Vendor:          "CentOS",
				DigestAlgorithm: PGPHASHALGO_SHA256,
			},
			wantInstalledFiles: CentOS6GlibcInstalledFiles,
		},
		{
			name:    "centos8 nodejs",
			pkgName: "nodejs",
			file:    "testdata/centos8-modularitylabel/Packages",
			want: &PackageInfo{
				Epoch:           1,
				Name:            "nodejs",
				Version:         "10.21.0",
				Release:         "3.module_el8.2.0+391+8da3adc6",
				Arch:            "x86_64",
				Size:            31483781,
				SourceRpm:       "nodejs-10.21.0-3.module_el8.2.0+391+8da3adc6.src.rpm",
				License:         "MIT and ASL 2.0 and ISC and BSD",
				Vendor:          "CentOS",
				Modularitylabel: "nodejs:10:8020020200707141642:6a468ee4",
				DigestAlgorithm: PGPHASHALGO_SHA256,
			},
			wantInstalledFiles: CentOS8NodejsInstalledFiles,
		},
		{
			name:    "CBL-Mariner 2.0 curl",
			pkgName: "curl",
			file:    "testdata/cbl-mariner-2.0/rpmdb.sqlite",
			want: &PackageInfo{
				Epoch:           0,
				Name:            "curl",
				Version:         "7.76.0",
				Release:         "6.cm2",
				Arch:            "x86_64",
				Size:            326023,
				SourceRpm:       "curl-7.76.0-6.cm2.src.rpm",
				License:         "MIT",
				Vendor:          "Microsoft Corporation",
				DigestAlgorithm: PGPHASHALGO_SHA256,
			},
			wantInstalledFiles: Mariner2CurlInstalledFiles,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Open(tt.file)
			require.NoError(t, err)

			got, err := db.Package(tt.pkgName)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			assert.NoError(t, err)

			gotInstalledFiles, err := got.InstalledFiles()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantInstalledFiles, gotInstalledFiles)

			// These fields are tested through InstalledFiles()
			got.Files = nil

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPackageFileList(t *testing.T) {
	vectors := []struct {
		file     string // Test input file
		fileList map[string][]FileInfo
	}{
		{
			file: "testdata/centos6-plain/Packages",
			fileList: map[string][]FileInfo{
				"libffi": {
					{Path: "/usr/lib64/libffi.so.5", Mode: 41471, Digest: "", Size: 15, Username: "root", Groupname: "root", Flags: 0},
					{Path: "/usr/lib64/libffi.so.5.0.6", Mode: 33261, Digest: "2009cab32d65011e653d7c87b49ad74541484467b3dc96be05bb2198b6c7a730", Size: 31720, Username: "root", Groupname: "root", Flags: 0},
					{Path: "/usr/share/doc/libffi-3.0.5", Mode: 16877, Digest: "", Size: 4096, Username: "root", Groupname: "root", Flags: 0},
					{Path: "/usr/share/doc/libffi-3.0.5/LICENSE", Mode: 33188, Digest: "b0421fa2fcb17d5d603cc46c66d69a8d943a03d48edbdfd672f24068bf6b2b65", Size: 1119, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/doc/libffi-3.0.5/README", Mode: 33188, Digest: "d8a1231d9090231272d547f7a7ee922298c20d34d4c79772f5ed4badc3a86f8d", Size: 10042, Username: "root", Groupname: "root", Flags: 2},
				},
			},
		},
		{
			file: "testdata/centos7-plain/Packages",
			fileList: map[string][]FileInfo{
				"ncurses": {
					{Path: "/usr/bin/captoinfo", Mode: 41471, Digest: "", Size: 3, Username: "root", Groupname: "root", Flags: 0},
					{Path: "/usr/bin/clear", Mode: 33261, Digest: "68353b0b989463d9e202362c843ee42c408dd1e08dd5e8e93733753749a96208", Size: 7192, Username: "root", Groupname: "root", Flags: 0},
					{Path: "/usr/bin/infocmp", Mode: 33261, Digest: "469fd67a3bdc7967a4c05b39a1b9a87635448520a619e608e702310480cef153", Size: 57416, Username: "root", Groupname: "root", Flags: 0},
					{Path: "/usr/bin/infotocap", Mode: 41471, Digest: "", Size: 3, Username: "root", Groupname: "root", Flags: 0},
					{Path: "/usr/bin/reset", Mode: 41471, Digest: "", Size: 4, Username: "root", Groupname: "root", Flags: 0},
					{Path: "/usr/bin/tabs", Mode: 33261, Digest: "85a7fb2d93019eb9ff1dd907dc649e9be5a49c704a26d94572418aea77affe46", Size: 15680, Username: "root", Groupname: "root", Flags: 0},
					{Path: "/usr/bin/tic", Mode: 33261, Digest: "df2ea23f0fdcd9a13a846de6d1880197d2fd60afe7b9b2945aa77f8595137a0c", Size: 65800, Username: "root", Groupname: "root", Flags: 0},
					{Path: "/usr/bin/toe", Mode: 33261, Digest: "b6cad57397f83d187c1361daf20d2b6a59982f9aa553a95d659edebe3116d26a", Size: 15800, Username: "root", Groupname: "root", Flags: 0},
					{Path: "/usr/bin/tput", Mode: 33261, Digest: "737da2a672c9ac17f86ebba733d316639365ad8459e16939fa03faea8e7d720f", Size: 15784, Username: "root", Groupname: "root", Flags: 0},
					{Path: "/usr/bin/tset", Mode: 33261, Digest: "50fa6ec48545da72f5c92040a39fbacb61ff1e45e14f9998a281b6c3285564c1", Size: 20072, Username: "root", Groupname: "root", Flags: 0},
					{Path: "/usr/share/doc/ncurses-5.9", Mode: 16877, Digest: "", Size: 75, Username: "root", Groupname: "root", Flags: 0},
					{Path: "/usr/share/doc/ncurses-5.9/ANNOUNCE", Mode: 33188, Digest: "1694388b7f5ce0819e1f8fd1c2b40979e82df58541ceb0c8b60c683f29378b78", Size: 13750, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/doc/ncurses-5.9/AUTHORS", Mode: 33188, Digest: "5e59823796c266525a92a6cd31bf144603a7d1b65362e48aa85e74a2b8093d50", Size: 2529, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/doc/ncurses-5.9/NEWS.bz2", Mode: 33188, Digest: "bb48de080557f81b9626ebd0baf48e559ae241dace93d57b7d618a441f8737fb", Size: 131412, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/doc/ncurses-5.9/README", Mode: 33188, Digest: "37e56186af1edbc4b0c41b85e224295fe2ef114399a488651ebc658f57bf80c7", Size: 10212, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/doc/ncurses-5.9/TO-DO", Mode: 33188, Digest: "9a40247610befa57d2c47d0fcd5d3ff3587edad07287f17a8279b98e4221692a", Size: 9651, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/man/man1/captoinfo.1m.gz", Mode: 33188, Digest: "40940eef25e38baaaa2ceb1cd7edb3508718400846485ed6f5c1e13bba1f1a34", Size: 2904, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/man/man1/clear.1.gz", Mode: 33188, Digest: "1ce7d795bb239d39ca5e11808f0766b456766ad1a914c6097beb7f9c8af638b9", Size: 1262, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/man/man1/infocmp.1m.gz", Mode: 33188, Digest: "2649e8bf304f00eb5624293515c4bd6eb7c7f847f33c3308dd8b76c5e44122dd", Size: 6952, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/man/man1/infotocap.1m.gz", Mode: 33188, Digest: "edd4d4bb4d79044d32f3422d5ba1e15302769b8a9a5e2fe0f8ce13967443bc25", Size: 1579, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/man/man1/reset.1.gz", Mode: 41471, Digest: "", Size: 9, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/man/man1/tabs.1.gz", Mode: 33188, Digest: "d9841dc62123346f2973dafb79874f794690f88725135a4d21805284cb973492", Size: 2253, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/man/man1/tic.1m.gz", Mode: 33188, Digest: "a5f8512a7a0e252225bd18efd0bcdbcee752e9bf5d539aef5948d3ab9230da8e", Size: 5677, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/man/man1/toe.1m.gz", Mode: 33188, Digest: "ca295431aa6b43954409c314bb15687dfc93b95ad8fbd5fcc183bd205008f995", Size: 1874, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/man/man1/tput.1.gz", Mode: 33188, Digest: "2f0d53ffbf8bef6d1a932a9955701ada4842f133ecdfb5b324604a703376bd2f", Size: 4529, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/man/man1/tset.1.gz", Mode: 33188, Digest: "7a2332f6d2305af034eafc9c94ed427f5d63c12087f611c4a499546fa9240a9c", Size: 4907, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/man/man5/term.5.gz", Mode: 33188, Digest: "0d53e8274fcd0c91ec79d1c7911c68d6993025335f0ed688413c38cf80edb04a", Size: 4431, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/man/man5/terminfo.5.gz", Mode: 33188, Digest: "c94c45d9713db4c2380b53fc5130e41ec3034e256a0cfc6f523676a49cf7f02e", Size: 33598, Username: "root", Groupname: "root", Flags: 2},
					{Path: "/usr/share/man/man7/term.7.gz", Mode: 33188, Digest: "29346e334d22d23120a45e692b0dc8f2d8262ef077149dbac3f775fbe0c9125d", Size: 4114, Username: "root", Groupname: "root", Flags: 2},
				},
			},
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
				t.Fatalf("ListPackages() error: %v", err)
			}

			assertedPkgLists := 0
			for _, p := range pkgList {
				if expected, ok := v.fileList[p.Name]; ok {
					assertedPkgLists++
					diffs := deep.Equal(expected, p.Files)
					if len(diffs) > 0 {
						t.Logf("Got files:")
						for _, actual := range p.Files {
							t.Logf("   %+v", actual)
						}
						for _, d := range diffs {
							t.Errorf(d)
						}
					}
				}
			}

			if assertedPkgLists != len(v.fileList) {
				t.Errorf("unexpected number of assertions: %d != %d", assertedPkgLists, len(v.fileList))
			}

		})
	}
}
