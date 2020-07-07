package rpmdb

import (
	"bytes"
	"encoding/binary"

	"golang.org/x/xerrors"
)

type PackageInfo struct {
	Epoch       int
	Name        string
	Version     string
	Release     string
	Arch        string
	Vendor      string
	Summary     string
	Size        uint32
	InstallTime uint32
}

const (
	// rpmTag_e
	// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.11.3-release/lib/rpmtag.h#L28
	RPMTAG_NAME        = 1000
	RPMTAG_VERSION     = 1001
	RPMTAG_RELEASE     = 1002
	RPMTAG_EPOCH       = 1003
	RPMTAG_SUMMARY     = 1004
	RPMTAG_INSTALLTIME = 1008
	RPMTAG_SIZE        = 1009
	RPMTAG_VENDOR      = 1011
	RPMTAG_ARCH        = 1022

	//rpmTagType_e
	// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.11.3-release/lib/rpmtag.h#L362
	RPM_NULL_TYPE         = 0
	RPM_CHAR_TYPE         = 1
	RPM_INT8_TYPE         = 2
	RPM_INT16_TYPE        = 3
	RPM_INT32_TYPE        = 4
	RPM_INT64_TYPE        = 5
	RPM_STRING_TYPE       = 6
	RPM_BIN_TYPE          = 7
	RPM_STRING_ARRAY_TYPE = 8
	RPM_I18NSTRING_TYPE   = 9
)

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.11.3-release/lib/tagexts.c#L649
func getNEVRA(indexEntries []indexEntry) (*PackageInfo, error) {
	pkgInfo := &PackageInfo{}

	for _, indexEntry := range indexEntries {
		switch indexEntry.Info.Tag {
		case RPMTAG_NAME:
			if indexEntry.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag name")
			}
			pkgInfo.Name = string(bytes.TrimRight(indexEntry.Data, "\x00"))
		case RPMTAG_EPOCH:
			if indexEntry.Info.Type != RPM_INT32_TYPE {
				return nil, xerrors.New("invalid tag epoch")
			}

			var epoch int32
			reader := bytes.NewReader(indexEntry.Data)
			if err := binary.Read(reader, binary.BigEndian, &epoch); err != nil {
				return nil, xerrors.Errorf("failed to read binary: %w", err)
			}
			pkgInfo.Epoch = int(epoch)
		case RPMTAG_VERSION:
			if indexEntry.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag version")
			}
			pkgInfo.Version = string(bytes.TrimRight(indexEntry.Data, "\x00"))
		case RPMTAG_RELEASE:
			if indexEntry.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag release")
			}
			pkgInfo.Release = string(bytes.TrimRight(indexEntry.Data, "\x00"))
		case RPMTAG_ARCH:
			if indexEntry.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag arch")
			}
			pkgInfo.Arch = string(bytes.TrimRight(indexEntry.Data, "\x00"))
		case RPMTAG_VENDOR:
			if indexEntry.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag vendor")
			}
			pkgInfo.Vendor = string(bytes.TrimRight(indexEntry.Data, "\x00"))
		case RPMTAG_SUMMARY:
			if indexEntry.Info.Type != RPM_I18NSTRING_TYPE {
				return nil, xerrors.New("invalid tag summary")
			}
			pkgInfo.Summary = string(bytes.TrimRight(indexEntry.Data, "\x00"))
		case RPMTAG_INSTALLTIME:
			if indexEntry.Info.Type != RPM_INT32_TYPE {
				return nil, xerrors.New("invalid tag install time")
			}
			pkgInfo.InstallTime = binary.BigEndian.Uint32(indexEntry.Data)
		case RPMTAG_SIZE:
			if indexEntry.Info.Type != RPM_INT32_TYPE {
				return nil, xerrors.New("invalid tag size")
			}
			pkgInfo.Size = binary.BigEndian.Uint32(indexEntry.Data)
		default:

		}
	}
	return pkgInfo, nil
}
