package rpmdb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileDigest(t *testing.T) {
	tests := []struct {
		algorithm DigestAlgorithm
		expected  string
	}{
		{
			algorithm: PGPHASHALGO_MD5,
			expected:  "md5",
		},
		{
			algorithm: PGPHASHALGO_SHA1,
			expected:  "sha1",
		},
		{
			algorithm: PGPHASHALGO_RIPEMD160,
			expected:  "ripemd160",
		},
		{
			algorithm: 4,
			expected:  "unknown-digest-algorithm",
		},
		{
			algorithm: PGPHASHALGO_MD2,
			expected:  "md2",
		},
		{
			algorithm: PGPHASHALGO_TIGER192,
			expected:  "tiger192",
		},
		{
			algorithm: PGPHASHALGO_HAVAL_5_160,
			expected:  "haval-5-160",
		},
		{
			algorithm: PGPHASHALGO_SHA256,
			expected:  "sha256",
		},
		{
			algorithm: PGPHASHALGO_SHA384,
			expected:  "sha384",
		},
		{
			algorithm: PGPHASHALGO_SHA512,
			expected:  "sha512",
		},
		{
			algorithm: PGPHASHALGO_SHA224,
			expected:  "sha224",
		},
		{
			algorithm: 12,
			expected:  "unknown-digest-algorithm",
		},
		// assert against known good values
		{
			algorithm: 1,
			expected:  "md5",
		},
		{
			algorithm: 2,
			expected:  "sha1",
		},
		{
			algorithm: 8,
			expected:  "sha256",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			actual := test.algorithm.String()
			assert.Equal(t, test.expected, actual)
		})
	}
}
