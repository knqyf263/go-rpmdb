package rpmdb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"path/filepath"
	"strings"
	"time"

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
	Summary         string
	PGP             string
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
		var err error
		switch ie.Info.Tag {
		case RPMTAG_DIRINDEXES:
			pkgInfo.DirIndexes, err = ie.ParseInt32Array()
		case RPMTAG_DIRNAMES:
			pkgInfo.DirNames, err = ie.ParseStringArray()
		case RPMTAG_BASENAMES:
			pkgInfo.BaseNames, err = ie.ParseStringArray()
		case RPMTAG_MODULARITYLABEL:
			pkgInfo.Modularitylabel, err = ie.ParseString()
		case RPMTAG_NAME:
			pkgInfo.Name, err = ie.ParseString()
		case RPMTAG_EPOCH:
			if ie.Data != nil {
				value, err := ie.ParseInt32()
				if err != nil {
					break
				}
				pkgInfo.Epoch = &value
			}
		case RPMTAG_VERSION:
			pkgInfo.Version, err = ie.ParseString()
		case RPMTAG_RELEASE:
			pkgInfo.Release, err = ie.ParseString()
		case RPMTAG_ARCH:
			pkgInfo.Arch, err = ie.ParseString()
		case RPMTAG_SOURCERPM:
			pkgInfo.SourceRpm, err = ie.ParseString()
		case RPMTAG_LICENSE:
			pkgInfo.License, err = ie.ParseString()
		case RPMTAG_VENDOR:
			pkgInfo.Vendor, err = ie.ParseString()
		case RPMTAG_SIZE:
			pkgInfo.Size, err = ie.ParseInt32()
		case RPMTAG_FILEDIGESTALGO:
			// note: all digests within a package entry only supports a single digest algorithm (there may be future support for
			// algorithm noted for each file entry, but currently unimplemented: https://github.com/rpm-software-management/rpm/blob/0b75075a8d006c8f792d33a57eae7da6b66a4591/lib/rpmtag.h#L256)
			digestAlgorithm, err := ie.ParseInt32()
			if err != nil {
				break
			}

			pkgInfo.DigestAlgorithm = DigestAlgorithm(digestAlgorithm)
		case RPMTAG_FILESIZES:
			// note: there is no distinction between int32, uint32, and []uint32
			pkgInfo.FileSizes, err = ie.ParseInt32Array()
		case RPMTAG_FILEDIGESTS:
			pkgInfo.FileDigests, err = ie.ParseStringArray()
		case RPMTAG_FILEMODES:
			// note: there is no distinction between int16, uint16, and []uint16
			pkgInfo.FileModes, err = ie.ParseUint16Array()
		case RPMTAG_FILEFLAGS:
			// note: there is no distinction between int32, uint32, and []uint32
			pkgInfo.FileFlags, err = ie.ParseInt32Array()
		case RPMTAG_FILEUSERNAME:
			pkgInfo.UserNames, err = ie.ParseStringArray()
		case RPMTAG_FILEGROUPNAME:
			pkgInfo.GroupNames, err = ie.ParseStringArray()
		case RPMTAG_SUMMARY:
			// some libraries have a string value instead of international string, so accounting for both
			if ie.Info.Type != RPM_I18NSTRING_TYPE && ie.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag summary")
			}
			// since this is an international string, getting the first null terminated string
			pkgInfo.Summary = string(bytes.Split(ie.Data, []byte{0})[0])
		case RPMTAG_PGP:
			pkgInfo.PGP, err = parsePGPSignature(ie)
			if err != nil {
				return nil, err
			}
		}

		if err != nil {
			return nil, xerrors.Errorf("error while parsing %v: %w",
						   ie.Info.TagName(), err)
		}
	}

	return pkgInfo, nil
}

const (
	sizeOfInt32  = 4
	sizeOfUInt16 = 2
)

func parsePGPSignature(ie indexEntry) (string, error) {
	type pgpSig struct {
		_          [3]byte
		Date       int32
		KeyID      [8]byte
		PubKeyAlgo uint8
		HashAlgo   uint8
	}

	type textSig struct {
		_          [2]byte
		PubKeyAlgo uint8
		HashAlgo   uint8
		_          [4]byte
		Date       int32
		_          [4]byte
		KeyID      [8]byte
	}

	type pgp4Sig struct {
		_          [2]byte
		PubKeyAlgo uint8
		HashAlgo   uint8
		_          [17]byte
		KeyID      [8]byte
		_          [2]byte
		Date       int32
	}

	pubKeyLookup := map[uint8]string{
		0x01: "RSA",
	}
	hashLookup := map[uint8]string{
		0x02: "SHA1",
		0x08: "SHA256",
	}

	if ie.Info.Type != RPM_BIN_TYPE {
		return "", xerrors.New("invalid PGP signature")
	}

	var tag, signatureType, version uint8
	r := bytes.NewReader(ie.Data)
	err := binary.Read(r, binary.BigEndian, &tag)
	if err != nil {
		return "", err
	}
	err = binary.Read(r, binary.BigEndian, &signatureType)
	if err != nil {
		return "", err
	}
	err = binary.Read(r, binary.BigEndian, &version)
	if err != nil {
		return "", err
	}

	var pubKeyAlgo, hashAlgo, pkgDate string
	var keyId [8]byte

	switch signatureType {
	case 0x01:
		switch version {
		case 0x1c:
			sig := textSig{}
			err = binary.Read(r, binary.BigEndian, &sig)
			if err != nil {
				return "", xerrors.Errorf("invalid PGP signature on decode: %w", err)
			}
			pubKeyAlgo = pubKeyLookup[sig.PubKeyAlgo]
			hashAlgo = hashLookup[sig.HashAlgo]
			pkgDate = time.Unix(int64(sig.Date), 0).UTC().Format("Mon Jan _2 15:04:05 2006")
			keyId = sig.KeyID
		default:
			sig := pgpSig{}
			err = binary.Read(r, binary.BigEndian, &sig)
			if err != nil {
				return "", xerrors.Errorf("invalid PGP signature on decode: %w", err)
			}
			pubKeyAlgo = pubKeyLookup[sig.PubKeyAlgo]
			hashAlgo = hashLookup[sig.HashAlgo]
			pkgDate = time.Unix(int64(sig.Date), 0).UTC().Format("Mon Jan _2 15:04:05 2006")
			keyId = sig.KeyID
		}
	case 0x02:
		switch version {
		case 0x33:
			sig := pgp4Sig{}
			err = binary.Read(r, binary.BigEndian, &sig)
			if err != nil {
				return "", xerrors.Errorf("invalid PGP signature on decode: %w", err)
			}
			pubKeyAlgo = pubKeyLookup[sig.PubKeyAlgo]
			hashAlgo = hashLookup[sig.HashAlgo]
			pkgDate = time.Unix(int64(sig.Date), 0).UTC().Format("Mon Jan _2 15:04:05 2006")
			keyId = sig.KeyID
		default:
			sig := pgpSig{}
			err = binary.Read(r, binary.BigEndian, &sig)
			if err != nil {
				return "", xerrors.Errorf("invalid PGP signature on decode: %w", err)
			}
			pubKeyAlgo = pubKeyLookup[sig.PubKeyAlgo]
			hashAlgo = hashLookup[sig.HashAlgo]
			pkgDate = time.Unix(int64(sig.Date), 0).UTC().Format("Mon Jan _2 15:04:05 2006")
			keyId = sig.KeyID
		}
	}

	result := fmt.Sprintf("%s/%s, %s, Key ID %x",
			     pubKeyAlgo, hashAlgo, pkgDate, keyId)

	return result, nil
}

func parseInt32Array(data []byte, arraySize int) ([]int32, error) {
	length := arraySize / sizeOfInt32
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
	length := arraySize / sizeOfUInt16
	values := make([]uint16, length)
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.BigEndian, &values); err != nil {
		return nil, xerrors.Errorf("failed to read binary: %w", err)
	}
	return values, nil
}

func parseString(data []byte) string {
	return string(bytes.TrimRight(data, "\x00"))
}

func parseStringArray(data []byte) []string {
	return strings.Split(string(bytes.TrimRight(data, "\x00")), "\x00")
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
