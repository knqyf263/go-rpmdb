package bdb

const (
	NoEncryptionAlgorithm = 0

	HashMagicNumber = 0x061561

	// the size (in bytes) of an in-page offset
	HashIndexEntrySize = 2
	// all DB pages have the same sized header (in bytes)
	PageHeaderSize = 26

	// all page types supported
	OverflowPageType     PageType = 7
	HashMetadataPageType PageType = 8
	HashPageType         PageType = 13
	HashOffIndexPageType PageType = 3 // aka HOFFPAGE

	HashOffPageSize = 12 // (in bytes)
)

type PageType = uint8
