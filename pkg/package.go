package rpmdb

import (
	"bytes"
	"encoding/binary"

	"golang.org/x/xerrors"
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
}

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/tagexts.c#L752
func getNEVRA(indexEntries []indexEntry) (*PackageInfo, error) {
	pkgInfo := &PackageInfo{}
	for _, indexEntry := range indexEntries {
		switch indexEntry.Info.Tag {
		case RPMTAG_MODULARITYLABEL:
			if indexEntry.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag modularitylabel")
			}
			pkgInfo.Modularitylabel = string(bytes.TrimRight(indexEntry.Data, "\x00"))
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
				return nil, xerrors.Errorf("failed to read binary (epoch): %w", err)
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
		case RPMTAG_SOURCERPM:
			if indexEntry.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag sourcerpm")
			}
			pkgInfo.SourceRpm = string(bytes.TrimRight(indexEntry.Data, "\x00"))
			if pkgInfo.SourceRpm == "(none)" {
				pkgInfo.SourceRpm = ""
			}
		case RPMTAG_LICENSE:
			if indexEntry.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag license")
			}
			pkgInfo.License = string(bytes.TrimRight(indexEntry.Data, "\x00"))
			if pkgInfo.License == "(none)" {
				pkgInfo.License = ""
			}
		case RPMTAG_VENDOR:
			if indexEntry.Info.Type != RPM_STRING_TYPE {
				return nil, xerrors.New("invalid tag vendor")
			}
			pkgInfo.Vendor = string(bytes.TrimRight(indexEntry.Data, "\x00"))
			if pkgInfo.Vendor == "(none)" {
				pkgInfo.Vendor = ""
			}
		case RPMTAG_SIZE:
			if indexEntry.Info.Type != RPM_INT32_TYPE {
				return nil, xerrors.New("invalid tag size")
			}

			var size int32
			reader := bytes.NewReader(indexEntry.Data)
			if err := binary.Read(reader, binary.BigEndian, &size); err != nil {
				return nil, xerrors.Errorf("failed to read binary (size): %w", err)
			}
			pkgInfo.Size = int(size)
		}

	}
	return pkgInfo, nil
}
