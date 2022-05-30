package rpmdb

import (
	"bytes"
	"encoding/binary"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"
)

type PackageInfo struct {
	Epoch           *int
	Name            string
	Version         string
	Release         string
	Arch            string
	SourceRpm       string
	Size            int
	License         string
	Vendor          string
	Modularitylabel string
	DigestAlgorithm DigestAlgorithm
	BaseNames       []string
	DirIndexes      []int32
	DirNames        []string
	FileSizes       []int32
	FileDigests     []string
	FileModes       []uint16
	FileFlags       []int32
	UserNames       []string
	GroupNames      []string
}

type FileInfo struct {
	Path      string
	Mode      uint16
	Digest    string
	Size      int32
	Username  string
	Groupname string
	Flags     FileFlags
}

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/tagexts.c#L752
func getNEVRA(indexEntries []indexEntry) (*PackageInfo, error) {
	pkgInfo := &PackageInfo{}
	for _, ie := range indexEntries {
		switch ie.Info.Tag {
		case RPMTAG_DIRINDEXES:
			if ie.Info.Type != RPM_INT32_TYPE {
				return nil, xerrors.New("invalid tag dir indexes")
			}

			dirIndexes, err := parseInt32Array(ie.Data, ie.Length)
			if err != nil {
				return nil, xerrors.Errorf("unable to read dir indexes: %w", err)
			}
			pkgInfo.DirIndexes = dirIndexes
		case RPMTAG_DIRNAMES:
			if ie.Info.Type != RPM_STRING_ARRAY_TYPE {
				return nil, xerrors.New("invalid tag dir names")
			}
			pkgInfo.DirNames = parseStringArray(ie.Data)
		case RPMTAG_BASENAMES:
			if ie.Info.Type != RPM_STRING_ARRAY_TYPE {
				return nil, xerrors.New("invalid tag base names")
			}
			pkgInfo.BaseNames = parseStringArray(ie.Data)
		case RPMTAG_MODULARITYLABEL:
			if ie.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag modularitylabel")
			}
			pkgInfo.Modularitylabel = string(bytes.TrimRight(ie.Data, "\x00"))
		case RPMTAG_NAME:
			if ie.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag name")
			}
			pkgInfo.Name = string(bytes.TrimRight(ie.Data, "\x00"))
		case RPMTAG_EPOCH:
			if ie.Info.Type != RPM_INT32_TYPE {
				return nil, xerrors.New("invalid tag epoch")
			}

			if ie.Data != nil {
				value, err := parseInt32(ie.Data)
				if err != nil {
					return nil, xerrors.Errorf("failed to parse epoch: %w", err)
				}
				pkgInfo.Epoch = &value
			}
		case RPMTAG_VERSION:
			if ie.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag version")
			}
			pkgInfo.Version = string(bytes.TrimRight(ie.Data, "\x00"))
		case RPMTAG_RELEASE:
			if ie.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag release")
			}
			pkgInfo.Release = string(bytes.TrimRight(ie.Data, "\x00"))
		case RPMTAG_ARCH:
			if ie.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag arch")
			}
			pkgInfo.Arch = string(bytes.TrimRight(ie.Data, "\x00"))
		case RPMTAG_SOURCERPM:
			if ie.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag sourcerpm")
			}
			pkgInfo.SourceRpm = string(bytes.TrimRight(ie.Data, "\x00"))
			if pkgInfo.SourceRpm == "(none)" {
				pkgInfo.SourceRpm = ""
			}
		case RPMTAG_LICENSE:
			if ie.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag license")
			}
			pkgInfo.License = string(bytes.TrimRight(ie.Data, "\x00"))
			if pkgInfo.License == "(none)" {
				pkgInfo.License = ""
			}
		case RPMTAG_VENDOR:
			if ie.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag vendor")
			}
			pkgInfo.Vendor = string(bytes.TrimRight(ie.Data, "\x00"))
			if pkgInfo.Vendor == "(none)" {
				pkgInfo.Vendor = ""
			}
		case RPMTAG_SIZE:
			if ie.Info.Type != RPM_INT32_TYPE {
				return nil, xerrors.New("invalid tag size")
			}

			size, err := parseInt32(ie.Data)
			if err != nil {
				return nil, xerrors.Errorf("failed to parse size: %w", err)
			}
			pkgInfo.Size = size
		case RPMTAG_FILEDIGESTALGO:
			// note: all digests within a package entry only supports a single digest algorithm (there may be future support for
			// algorithm noted for each file entry, but currently unimplemented: https://github.com/rpm-software-management/rpm/blob/0b75075a8d006c8f792d33a57eae7da6b66a4591/lib/rpmtag.h#L256)
			if ie.Info.Type != RPM_INT32_TYPE {
				return nil, xerrors.New("invalid tag digest algo")
			}

			digestAlgorithm, err := parseInt32(ie.Data)
			if err != nil {
				return nil, xerrors.Errorf("failed to parse digest algo: %w", err)
			}

			pkgInfo.DigestAlgorithm = DigestAlgorithm(digestAlgorithm)
		case RPMTAG_FILESIZES:
			// note: there is no distinction between int32, uint32, and []uint32
			if ie.Info.Type != RPM_INT32_TYPE {
				return nil, xerrors.New("invalid tag file-sizes")
			}
			fileSizes, err := parseInt32Array(ie.Data, ie.Length)
			if err != nil {
				return nil, xerrors.Errorf("failed to parse file-sizes: %w", err)
			}
			pkgInfo.FileSizes = fileSizes
		case RPMTAG_FILEDIGESTS:
			if ie.Info.Type != RPM_STRING_ARRAY_TYPE {
				return nil, xerrors.New("invalid tag file-digests")
			}
			pkgInfo.FileDigests = parseStringArray(ie.Data)
		case RPMTAG_FILEMODES:
			// note: there is no distinction between int16, uint16, and []uint16
			if ie.Info.Type != RPM_INT16_TYPE {
				return nil, xerrors.New("invalid tag file-modes")
			}
			fileModes, err := uint16Array(ie.Data, ie.Length)
			if err != nil {
				return nil, xerrors.Errorf("failed to parse file-modes: %w", err)
			}
			pkgInfo.FileModes = fileModes
		case RPMTAG_FILEFLAGS:
			// note: there is no distinction between int32, uint32, and []uint32
			if ie.Info.Type != RPM_INT32_TYPE {
				return nil, xerrors.New("invalid tag file-flags")
			}
			fileFlags, err := parseInt32Array(ie.Data, ie.Length)
			if err != nil {
				return nil, xerrors.Errorf("failed to parse file-flags: %w", err)
			}
			pkgInfo.FileFlags = fileFlags
		case RPMTAG_FILEUSERNAME:
			if ie.Info.Type != RPM_STRING_ARRAY_TYPE {
				return nil, xerrors.New("invalid tag usernames")
			}
			pkgInfo.UserNames = parseStringArray(ie.Data)
		case RPMTAG_FILEGROUPNAME:
			if ie.Info.Type != RPM_STRING_ARRAY_TYPE {
				return nil, xerrors.New("invalid tag groupnames")
			}
			pkgInfo.GroupNames = parseStringArray(ie.Data)
		}
	}

	return pkgInfo, nil
}

const (
	sizeOfInt32  = 4
	sizeOfUInt16 = 2
)

func parseInt32Array(data []byte, arraySize int) ([]int32, error) {
	var length = arraySize / sizeOfInt32
	values := make([]int32, length)
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.BigEndian, &values); err != nil {
		return nil, xerrors.Errorf("failed to read binary: %w", err)
	}
	return values, nil
}

func parseInt32(data []byte) (int, error) {
	var value int32
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.BigEndian, &value); err != nil {
		return 0, xerrors.Errorf("failed to read binary: %w", err)
	}
	return int(value), nil
}

func uint16Array(data []byte, arraySize int) ([]uint16, error) {
	var length = arraySize / sizeOfUInt16
	values := make([]uint16, length)
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.BigEndian, &values); err != nil {
		return nil, xerrors.Errorf("failed to read binary: %w", err)
	}
	return values, nil
}

func parseStringArray(data []byte) []string {
	elements := strings.Split(string(data), "\x00")
	if len(elements) > 0 && elements[len(elements)-1] == "" {
		return elements[:len(elements)-1]
	}
	return elements
}

func (p *PackageInfo) InstalledFileNames() ([]string, error) {
	if len(p.DirNames) == 0 || len(p.DirIndexes) == 0 || len(p.BaseNames) == 0 {
		return nil, nil
	}

	// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/tagexts.c#L68-L70
	if len(p.DirIndexes) != len(p.BaseNames) || len(p.DirNames) > len(p.BaseNames) {
		return nil, xerrors.Errorf("invalid rpm %s", p.Name)
	}

	var filePaths []string
	for i, baseName := range p.BaseNames {
		dir := p.DirNames[p.DirIndexes[i]]
		filePaths = append(filePaths, filepath.Join(dir, baseName))
	}
	return filePaths, nil
}

func (p *PackageInfo) InstalledFiles() ([]FileInfo, error) {
	fileNames, err := p.InstalledFileNames()
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for i, fileName := range fileNames {
		var digest, username, groupname string
		var mode uint16
		var size, flags int32

		if p.FileDigests != nil && len(p.FileDigests) > i {
			digest = p.FileDigests[i]
		}

		if p.FileModes != nil && len(p.FileModes) > i {
			mode = p.FileModes[i]
		}

		if p.FileSizes != nil && len(p.FileSizes) > i {
			size = p.FileSizes[i]
		}

		if p.UserNames != nil && len(p.UserNames) > i {
			username = p.UserNames[i]
		}

		if p.GroupNames != nil && len(p.GroupNames) > i {
			groupname = p.GroupNames[i]
		}

		if p.FileFlags != nil && len(p.FileFlags) > i {
			flags = p.FileFlags[i]
		}

		record := FileInfo{
			Path:      fileName,
			Mode:      mode,
			Digest:    digest,
			Size:      size,
			Username:  username,
			Groupname: groupname,
			Flags:     FileFlags(flags),
		}
		files = append(files, record)
	}

	return files, nil
}

func (p *PackageInfo) EpochNum() int {
	if p.Epoch == nil {
		return 0
	}
	return *p.Epoch
}
