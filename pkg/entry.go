package rpmdb

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/xerrors"
	"io"
	"unsafe"
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

	dataStart := int32(unsafe.Sizeof(il)) + int32(unsafe.Sizeof(dl)) + il*int32(unsafe.Sizeof(entryInfo{}))

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

	// Ignore negative offset
	indexEntries := regionSwab(data, peList[1:], dataStart, int(dl))
	return indexEntries, nil
}

// ref. https://github.com/rpm-software-management/rpm/blob/7a2f891d25d78cf797c789ac6859b5f2c589d296/lib/header.c#L498
func regionSwab(data []byte, peList []entryInfo, dataStart int32, dl int) []indexEntry {
	indexEntries := make([]indexEntry, len(peList))
	for i := 0; i < len(peList); i++ {
		pe := peList[i]
		indexEntry := indexEntry{
			Info: entryInfo{
				Type:   HtonlU(pe.Type),
				Count:  HtonlU(pe.Count),
				Offset: Htonl(pe.Offset),
				Tag:    Htonl(pe.Tag),
			},
		}
		if i < len(peList)-1 {
			indexEntry.Length = int(Htonl(peList[i+1].Offset) - indexEntry.Info.Offset)
		} else {
			indexEntry.Length = dl - int(indexEntry.Info.Offset)
		}

		start := dataStart + indexEntry.Info.Offset
		end := int(start) + indexEntry.Length
		indexEntry.Data = data[start:end]

		indexEntries[i] = indexEntry
	}
	return indexEntries
}
