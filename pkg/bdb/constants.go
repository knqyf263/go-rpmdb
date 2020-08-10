package bdb

const (
	NoEncryptionAlgorithm = 0

	HashMagicNumber = 0x061561

	// the size (in bytes) of an in-page offset
	HashIndexEntrySize = 2
	// all DB pages have the same sized header (in bytes)
	PageHeaderSize = 26

	// all page types supported
	HashMetadataPageType PageType = 8
	HashPageType         PageType = 13
	HashOffIndexPageType PageType = 3 // a.k.a HOFFPAGE

	HashOffPageSize = 12 // (in bytes)
)

type PageType = uint8
