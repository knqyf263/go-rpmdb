package rpmdb

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/xerrors"
	"strings"
)

type PackageInfo struct {
	Epoch           int
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
	Files           []FileInfo
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
	var baseNames []string
	var dirIndexes []int32
	var dirNames []string
	var fileSizes []int32
	var fileDigests []string
	var fileModes []uint16
	var fileFlags []int32
	var userNames []string
	var groupNames []string
	var err error

	pkgInfo := &PackageInfo{}
	for _, ie := range indexEntries {
		switch ie.Info.Tag {
		case RPMTAG_DIRINDEXES:
			if ie.Info.Type != RPM_INT32_TYPE {
				return nil, xerrors.New("invalid tag dir indexes")
			}

			dirIndexes, err = parseInt32Array(ie.Data, ie.Length)
			if err != nil {
				return nil, xerrors.Errorf("unable to read dir indexes: %w", err)
			}
		case RPMTAG_DIRNAMES:
			if ie.Info.Type != RPM_STRING_ARRAY_TYPE {
				return nil, xerrors.New("invalid tag dir names")
			}
			dirNames = parseStringArray(ie.Data)
		case RPMTAG_BASENAMES:
			if ie.Info.Type != RPM_STRING_ARRAY_TYPE {
				return nil, xerrors.New("invalid tag base names")
			}
			baseNames = parseStringArray(ie.Data)
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

			epoch, err := parseInt32(ie.Data)
			if err != nil {
				return nil, xerrors.Errorf("failed to parse epoch: %w", err)
			}
			pkgInfo.Epoch = epoch
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
			fileSizes, err = parseInt32Array(ie.Data, ie.Length)
			if err != nil {
				return nil, xerrors.Errorf("failed to parse file-sizes: %w", err)
			}
		case RPMTAG_FILEDIGESTS:
			if ie.Info.Type != RPM_STRING_ARRAY_TYPE {
				return nil, xerrors.New("invalid tag file-digests")
			}
			fileDigests = parseStringArray(ie.Data)
		case RPMTAG_FILEMODES:
			// note: there is no distinction between int16, uint16, and []uint16
			if ie.Info.Type != RPM_INT16_TYPE {
				return nil, xerrors.New("invalid tag file-modes")
			}
			fileModes, err = uint16Array(ie.Data, ie.Length)
			if err != nil {
				return nil, xerrors.Errorf("failed to parse file-modes: %w", err)
			}
		case RPMTAG_FILEFLAGS:
			// note: there is no distinction between int32, uint32, and []uint32
			if ie.Info.Type != RPM_INT32_TYPE {
				return nil, xerrors.New("invalid tag file-flags")
			}
			fileFlags, err = parseInt32Array(ie.Data, ie.Length)
			if err != nil {
				return nil, xerrors.Errorf("failed to parse file-flags: %w", err)
			}
		case RPMTAG_FILEUSERNAME:
			if ie.Info.Type != RPM_STRING_ARRAY_TYPE {
				return nil, xerrors.New("invalid tag usernames")
			}
			userNames = parseStringArray(ie.Data)
		case RPMTAG_FILEGROUPNAME:
			if ie.Info.Type != RPM_STRING_ARRAY_TYPE {
				return nil, xerrors.New("invalid tag groupnames")
			}
			groupNames = parseStringArray(ie.Data)
		}
	}

	// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/tagexts.c#L68-L70
	if len(dirIndexes) != len(baseNames) || len(dirNames) > len(baseNames) {
		return nil, xerrors.Errorf("invalid rpm %s", pkgInfo.Name)
	}

	// piece together a list of files and their metadata
	var files []FileInfo
	if dirNames != nil && dirIndexes != nil {
		for i, file := range baseNames {
			var digest, username, groupname string
			var mode uint16
			var size, flags int32

			if fileDigests != nil && len(fileDigests) > i {
				digest = fileDigests[i]
			}

			if fileModes != nil && len(fileModes) > i {
				mode = fileModes[i]
			}

			if fileSizes != nil && len(fileSizes) > i {
				size = fileSizes[i]
			}

			if userNames != nil && len(userNames) > i {
				username = userNames[i]
			}

			if groupNames != nil && len(groupNames) > i {
				groupname = groupNames[i]
			}

			if fileFlags != nil && len(fileFlags) > i {
				flags = fileFlags[i]
			}

			record := FileInfo{
				Path:      dirNames[dirIndexes[i]] + file,
				Mode:      mode,
				Digest:    digest,
				Size:      size,
				Username:  username,
				Groupname: groupname,
				Flags:     FileFlags(flags),
			}
			files = append(files, record)
		}
	}

	pkgInfo.Files = files

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

func (p *PackageInfo) InstalledFiles() ([]string, error) {
	var filePaths []string
	for _, fileInfo := range p.Files {
		filePaths = append(filePaths, fileInfo.Path)
	}

	return filePaths, nil
}
