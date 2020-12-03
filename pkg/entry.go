package rpmdb

import (
	"bytes"
	"encoding/binary"
	"io"
	"unsafe"

	"golang.org/x/xerrors"
)

const (
	// REGION_TAG_COUNT is sizeof(entryInfo_s)
	REGION_TAG_COUNT = 16

	RPMTAG_HEADERIMAGE      = 61
	RPMTAG_HEADERSIGNATURES = 62
	RPMTAG_HEADERIMMUTABLE  = 63
)

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.11.3-release/lib/header_internal.h#L13-L19
type entryInfo struct {
	Tag    int32  /*!< Tag identifier. */
	Type   uint32 /*!< Tag data type. */
	Offset int32  /*!< Offset into data segment (ondisk only). */
	Count  uint32 /*!< Number of tag elements. */
}

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.11.3-release/lib/header_internal.h#L27-L33
type indexEntry struct {
	Info   entryInfo
	Length int
	Rdlen  int
	Data   []byte
}

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.11.3-release/lib/header.c#L789
func headerImport(data []byte) ([]indexEntry, error) {
	var il, dl int32
	var err error
	reader := bytes.NewReader(data)

	if err = binary.Read(reader, binary.BigEndian, &il); err != nil {
		return nil, xerrors.Errorf("invalid index length: %w", err)
	}
	if err = binary.Read(reader, binary.BigEndian, &dl); err != nil {
		return nil, xerrors.Errorf("invalid data length: %w", err)
	}

	peList := make([]entryInfo, il)
	for i := 0; i < int(il); i++ {
		var pe entryInfo
		err = binary.Read(reader, binary.LittleEndian, &pe)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, xerrors.Errorf("failed to read entry info: %w", err)
		}
		peList[i] = pe
	}

	if len(peList) == 0 {
		return nil, xerrors.New("peList is empty error")
	}
	einfo := ei2h(peList[0])

	dataStart := int32(unsafe.Sizeof(il)) + int32(unsafe.Sizeof(dl)) + il*int32(unsafe.Sizeof(entryInfo{}))
	if !(einfo.Tag == RPMTAG_HEADERIMAGE || einfo.Tag == RPMTAG_HEADERSIGNATURES || einfo.Tag == RPMTAG_HEADERIMMUTABLE) {
		return regionSwab(data, peList[1:], dataStart, int(dl)), nil
	}

	// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header.c#L1842
	var trailer entryInfo
	regionEnd := dataStart + einfo.Offset
	if err := binary.Read(bytes.NewReader(data[regionEnd:regionEnd+REGION_TAG_COUNT]), binary.LittleEndian, &trailer); err != nil {
		return nil, xerrors.Errorf("invalid trailer: %w", err)
	}
	einfo = ei2h(trailer)

	ril := -einfo.Offset / REGION_TAG_COUNT
	if !(ril > 1 && ril < int32(len(peList))) {
		return nil, xerrors.New("invalid region index length error")
	}

	// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header.c#L917
	indexEntries := regionSwab(data, peList[1:ril], dataStart, int(dl))
	if ril < int32(len(peList)-1) {
		dribbleIndexEntries := regionSwab(data, peList[ril:], dataStart, int(dl))

		// Dribble entries replace duplicate region entries.
		for i, indexEntry := range indexEntries {
			for _, newIndexEntry := range dribbleIndexEntries {
				if indexEntry.Info.Tag == newIndexEntry.Info.Tag {
					indexEntries = deleteIndex(indexEntries, i)
					break
				}
			}
		}

		indexEntries = append(indexEntries, dribbleIndexEntries...)
	}

	return indexEntries, nil
}

func deleteIndex(indexEntries []indexEntry, i int) []indexEntry {
	if len(indexEntries)-1 == i {
		return indexEntries[:i]
	}
	return append(indexEntries[:i], indexEntries[i+1:]...)

}

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header_internal.h#L42
func ei2h(pe entryInfo) entryInfo {
	return entryInfo{
		Type:   HtonlU(pe.Type),
		Count:  HtonlU(pe.Count),
		Offset: Htonl(pe.Offset),
		Tag:    Htonl(pe.Tag),
	}
}

// ref. https://github.com/rpm-software-management/rpm/blob/7a2f891d25d78cf797c789ac6859b5f2c589d296/lib/header.c#L498
func regionSwab(data []byte, peList []entryInfo, dataStart int32, dl int) []indexEntry {
	indexEntries := make([]indexEntry, len(peList))
	for i := 0; i < len(peList); i++ {
		pe := peList[i]
		indexEntry := indexEntry{Info: ei2h(pe)}

		start := dataStart + indexEntry.Info.Offset
		if i < len(peList)-1 {
			indexEntry.Length = int(Htonl(peList[i+1].Offset) - indexEntry.Info.Offset)
		} else {

			if indexEntry.Length = dataLength(data[start:], indexEntry); indexEntry.Length < 0 {
				indexEntry.Length = dl - int(indexEntry.Info.Offset)
			}
		}

		end := int(start) + indexEntry.Length
		indexEntry.Data = data[start:end]

		indexEntries[i] = indexEntry
	}
	return indexEntries
}

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header.c#L440
func dataLength(data []byte, entry indexEntry) int {
	switch entry.Info.Type {
	case RPM_STRING_TYPE:
		if entry.Info.Count != 1 {
			return -1

		}
		return strtaglen(data)
	default:
		return -1
	}
}

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header.c#L408
func strtaglen(data []byte) int {
	for i := 0; i < len(data); i++ {
		if data[i] == 0x00 {
			return i
		}
	}
	return -1
}
