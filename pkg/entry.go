package rpmdb

import (
	"bytes"
	"encoding/binary"
	"io"
	"unsafe"

	"golang.org/x/xerrors"
)

const (
	REGION_TAG_COUNT = int32(unsafe.Sizeof(entryInfo{}))
	REGION_TAG_TYPE  = RPM_BIN_TYPE

	// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header.c#L113
	headerMaxbytes = 256 * 1024 * 1024
)

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header.c#L52
var (
	typeSizes = [16]int{
		0,  /*!< RPM_NULL_TYPE */
		1,  /*!< RPM_CHAR_TYPE */
		1,  /*!< RPM_INT8_TYPE */
		2,  /*!< RPM_INT16_TYPE */
		4,  /*!< RPM_INT32_TYPE */
		8,  /*!< RPM_INT64_TYPE */
		-1, /*!< RPM_STRING_TYPE */
		1,  /*!< RPM_BIN_TYPE */
		-1, /*!< RPM_STRING_ARRAY_TYPE */
		-1, /*!< RPM_I18NSTRING_TYPE */
		0,
		0,
		0,
		0,
		0,
		0,
	}
	typeAlign = [16]int{
		1, /*!< RPM_NULL_TYPE */
		1, /*!< RPM_CHAR_TYPE */
		1, /*!< RPM_INT8_TYPE */
		2, /*!< RPM_INT16_TYPE */
		4, /*!< RPM_INT32_TYPE */
		8, /*!< RPM_INT64_TYPE */
		1, /*!< RPM_STRING_TYPE */
		1, /*!< RPM_BIN_TYPE */
		1, /*!< RPM_STRING_ARRAY_TYPE */
		1, /*!< RPM_I18NSTRING_TYPE */
		0,
		0,
		0,
		0,
		0,
		0,
	}
)

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header_internal.h#L14-L20
type entryInfo struct {
	Tag    int32  /*!< Tag identifier. */
	Type   uint32 /*!< Tag data type. */
	Offset int32  /*!< Offset into data segment (ondisk only). */
	Count  uint32 /*!< Number of tag elements. */
}

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header.c#L88-L94
type indexEntry struct {
	Info   entryInfo
	Length int
	Rdlen  int
	Data   []byte
}

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header_internal.h#L23
type hdrblob struct {
	peList    []entryInfo
	il        int32
	dl        int32
	pvlen     int32
	dataStart int32
	dataEnd   int32
	regionTag int32
	ril       int32
	rdl       int32
}

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.11.3-release/lib/header.c#L789
func headerImport(data []byte) ([]indexEntry, error) {
	blob, err := hdrblobInit(data)
	if err != nil {
		return nil, xerrors.Errorf("failed to hdrblobInit error: %w", err)
	}
	indexEntries, err := hdrblobImport(*blob, data)
	if err != nil {
		return nil, xerrors.Errorf("failed to hdrblobImport error: %w", err)
	}
	return indexEntries, nil
}

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header.c#L1974
func hdrblobInit(data []byte) (*hdrblob, error) {
	var blob hdrblob
	var err error
	reader := bytes.NewReader(data)

	if err = binary.Read(reader, binary.BigEndian, &blob.il); err != nil {
		return nil, xerrors.Errorf("invalid index length: %w", err)
	}
	if err = binary.Read(reader, binary.BigEndian, &blob.dl); err != nil {
		return nil, xerrors.Errorf("invalid data length: %w", err)
	}
	blob.dataStart = int32(unsafe.Sizeof(blob.il)) + int32(unsafe.Sizeof(blob.dl)) + blob.il*int32(unsafe.Sizeof(entryInfo{}))
	blob.pvlen = int32(unsafe.Sizeof(blob.il)) + int32(unsafe.Sizeof(blob.dl)) + blob.il*int32(unsafe.Sizeof(entryInfo{})) + blob.dl
	blob.dataEnd = blob.dataStart + blob.dl

	blob.peList = make([]entryInfo, blob.il)
	for i := 0; i < int(blob.il); i++ {
		var pe entryInfo
		err = binary.Read(reader, binary.LittleEndian, &pe)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, xerrors.Errorf("failed to read entry info: %w", err)
		}
		blob.peList[i] = pe
	}
	if blob.pvlen >= headerMaxbytes {
		return nil, xerrors.Errorf("blob size error: size(%d) BAD, 8 + 16 * il(%d) + dl(%d)", blob.pvlen, blob.il, blob.dl)
	}

	if len(blob.peList) == 0 {
		return nil, xerrors.New("peList is empty error")
	}

	if err := hdrblobVerifyRegion(&blob, data); err != nil {
		return nil, xerrors.Errorf("failed to hdrblobVerifyRegion error: %w", err)
	}
	if err := hdrblobVerifyInfo(&blob); err != nil {
		return nil, xerrors.Errorf("failed to hdrblobVerifyInfo error: %w", err)
	}

	return &blob, nil
}

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header.c#L880
func hdrblobImport(blob hdrblob, data []byte) ([]indexEntry, error) {
	ril := blob.ril
	if blob.peList[0].Offset == 0 {
		ril = blob.il
	}

	// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header.c#L917
	indexEntries := regionSwab(data, blob.peList[1:ril], blob.dataStart, int(blob.dl))
	if blob.ril < int32(len(blob.peList)-1) {
		dribbleIndexEntries := regionSwab(data, blob.peList[ril:], blob.dataStart, int(blob.dl))
		uniqTagMap := make(map[int32]indexEntry)

		for _, indexEntry := range append(indexEntries, dribbleIndexEntries...) {
			uniqTagMap[indexEntry.Info.Tag] = indexEntry
		}

		var ies []indexEntry
		for _, indexEntry := range uniqTagMap {
			ies = append(ies, indexEntry)
		}

		return ies, nil
	}
	return indexEntries, nil
}

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header.c#L298-L303
func hdrblobVerifyInfo(blob *hdrblob) error {
	var end int32

	peOffset := 0
	if blob.regionTag != 0 {
		peOffset = 1
	}
	typechk := blob.regionTag == RPMTAG_HEADERIMMUTABLE || blob.regionTag == RPMTAG_HEADERIMAGE

	for _, pe := range blob.peList[peOffset:] {
		info := ei2h(pe)
		if end > info.Offset {
			return xerrors.New("")
		}

		if hdrchkTag(info.Tag) {
			return xerrors.New("")
		}
		if hdrchkType(info.Type) {
			return xerrors.New("")
		}
		if hdrchkAlign(info.Type, info.Offset) {
			return xerrors.New("")
		}
		if hdrchkRange(blob.dl, info.Offset) {
			return xerrors.New("")
		}
		if typechk && hdrchkTagType(info.Tag, info.Type) {
			return xerrors.New("")
		}

		// TODO: verify data length
		/* Verify the data actually fits */
		// len = dataLength(info.Type, ds + info.Offset,
		// 		 info.count, 1, ds + blob.dl);
		// end = info.offset + len;
		// if (hdrchkRange(blob->dl, end) || len <= 0)
		//     goto err;
	}
	return nil
}

func hdrchkTag(tag int32) bool {
	return tag < HEADER_I18NTABLE
}

func hdrchkType(t uint32) bool {
	return t < RPM_MIN_TYPE || t > RPM_MAX_TYPE
}

func hdrchkAlign(t uint32, offset int32) bool {
	return offset&int32(typeAlign[t]-1) != 0
}

func hdrchkTagType(tag int32, t uint32) bool {
	// TODO:
	return false
}

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header.c#L1791
func hdrblobVerifyRegion(blob *hdrblob, data []byte) error {
	var einfo entryInfo
	var regionTag int32

	if blob.il < 1 {
		return xerrors.New("region no tags error")
	}

	einfo = ei2h(blob.peList[0])

	if einfo.Tag == RPMTAG_HEADERIMAGE ||
		einfo.Tag == RPMTAG_HEADERSIGNATURES ||
		einfo.Tag == RPMTAG_HEADERIMMUTABLE {

		regionTag = einfo.Tag
	}

	if einfo.Tag != regionTag {
		return nil
	}

	if !(einfo.Type == REGION_TAG_TYPE && einfo.Count == uint32(REGION_TAG_COUNT)) {
		return xerrors.New("invalid region tag")
	}

	if hdrchkRange(blob.dl, einfo.Offset+REGION_TAG_COUNT) {
		return xerrors.New("invalid region offset")
	}

	// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header.c#L1842
	var trailer entryInfo
	regionEnd := blob.dataStart + einfo.Offset
	if err := binary.Read(bytes.NewReader(data[regionEnd:regionEnd+REGION_TAG_COUNT]), binary.LittleEndian, &trailer); err != nil {
		return xerrors.Errorf("failed to parse trailer error: %w", err)
	}
	blob.rdl = regionEnd + REGION_TAG_COUNT - blob.dataStart

	if regionTag == RPMTAG_HEADERSIGNATURES && einfo.Tag == RPMTAG_HEADERIMAGE {
		einfo.Tag = RPMTAG_HEADERSIGNATURES
	}

	if !(einfo.Tag == regionTag && einfo.Type == REGION_TAG_TYPE && einfo.Count == uint32(REGION_TAG_COUNT)) {
		return xerrors.New("invalid region trailer")
	}

	einfo = ei2h(trailer)
	einfo.Offset = -einfo.Offset
	blob.ril = einfo.Offset / int32(unsafe.Sizeof(blob.peList[0]))
	if (einfo.Offset%REGION_TAG_COUNT) != 0 || hdrchkRange(blob.il, blob.ril) || hdrchkRange(blob.dl, blob.rdl) {
		return xerrors.New("invalid region size")
	}

	blob.regionTag = regionTag

	return nil
}

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header.c#L158
func hdrchkRange(dl, offset int32) bool {
	return offset < 0 || offset > dl
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

// ref. https://github.com/rpm-software-management/rpm/blob/rpm-4.14.3-release/lib/header.c#L498
func regionSwab(data []byte, peList []entryInfo, dataStart int32, dl int) []indexEntry {
	// TODO: Make use of dl and make sure dl is calculated properly.
	// TODO: verify dataEnd
	// TODO: refactoring same interface
	indexEntries := make([]indexEntry, len(peList))
	for i := 0; i < len(peList); i++ {
		pe := peList[i]
		indexEntry := indexEntry{Info: ei2h(pe)}

		// Unsupport RPM_STRING_ARRAY_TYPE and RPM_I18NSTRING_TYPE tag type
		if indexEntry.Info.Type == RPM_STRING_ARRAY_TYPE ||
			indexEntry.Info.Type == RPM_I18NSTRING_TYPE {

			continue
		}

		start := dataStart + indexEntry.Info.Offset
		if i < len(peList)-1 && typeSizes[indexEntry.Info.Type] == -1 {
			indexEntry.Length = int(Htonl(peList[i+1].Offset) - indexEntry.Info.Offset)
		} else {
			indexEntry.Length = dataLength(data[start:], indexEntry)
			if indexEntry.Length < 0 {
				continue
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
		return bytes.IndexByte(data, byte(0x00))
	default:
		return typeSizes[entry.Info.Type&0xf] * int(entry.Info.Count)
	}
}
