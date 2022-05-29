package rpmdb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatFileFlags(t *testing.T) {
	tests := []struct {
		flags    FileFlags
		expected string
	}{
		// empty
		{
			flags:    0,
			expected: "",
		},
		// check that the formatting works relative to the configured bits
		{
			flags:    FileFlags(RPMFILE_CONFIG),
			expected: "c",
		},
		{
			flags:    FileFlags(RPMFILE_DOC),
			expected: "d",
		},
		{
			flags:    FileFlags(RPMFILE_MISSINGOK),
			expected: "m",
		},
		{
			flags:    FileFlags(RPMFILE_NOREPLACE),
			expected: "n",
		},
		{
			flags:    FileFlags(RPMFILE_SPECFILE),
			expected: "s",
		},
		{
			flags:    FileFlags(RPMFILE_GHOST),
			expected: "g",
		},
		{
			flags:    FileFlags(RPMFILE_LICENSE),
			expected: "l",
		},
		{
			flags:    FileFlags(RPMFILE_README),
			expected: "r",
		},
		{
			flags:    FileFlags(RPMFILE_ARTIFACT),
			expected: "a",
		},
		{
			flags:    FileFlags(RPMFILE_CONFIG | RPMFILE_DOC | RPMFILE_SPECFILE | RPMFILE_MISSINGOK | RPMFILE_NOREPLACE | RPMFILE_GHOST | RPMFILE_LICENSE | RPMFILE_README | RPMFILE_ARTIFACT),
			expected: "dcsmnglra",
		},
		{
			flags:    FileFlags(RPMFILE_DOC | RPMFILE_ARTIFACT),
			expected: "da",
		},
		// check that the formatting matches relative to verified correct values
		// see helpful examples from: rpm  --dbpath=/var/lib/rpm -qa --queryformat '%{FILEFLAGS:fflags}|%{FILEFLAGS}\n'
		{
			flags:    FileFlags(89),
			expected: "cmng",
		},
		{
			flags:    FileFlags(16),
			expected: "n",
		},
		{
			flags:    FileFlags(64),
			expected: "g",
		},
		{
			flags:    FileFlags(17),
			expected: "cn",
		},
		{
			flags:    FileFlags(4096),
			expected: "a",
		},
	}
	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			assert.Equal(t, test.expected, test.flags.String())
		})
	}
}
