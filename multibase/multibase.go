package multibase

import (
	"encoding/base32"
	"encoding/base64"
)

// Multibase prefixes
// Available at https://github.com/multiformats/multibase/blob/master/multibase.csv
const (
	prefixBase16         byte = 'f'
	prefixBase16Upper    byte = 'F'
	prefixBase32         byte = 'b'
	prefixBase32Upper    byte = 'B'
	prefixBase32Hex      byte = 'v'
	prefixBase32HexUpper byte = 'V'
	prefixBase64URL      byte = 'u'
)

var (
	base16Encoding         = &hexEncoder{upper: false}
	base16UpperEncoding    = &hexEncoder{upper: true}
	base32Encoding         = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567").WithPadding(base32.NoPadding)
	base32UpperEncoding    = base32.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567").WithPadding(base32.NoPadding)
	base32HexEncoding      = base32.NewEncoding("0123456789abcdefghijklmnopqrstuv").WithPadding(base32.NoPadding)
	base32HexUpperEncoding = base32.NewEncoding("0123456789ABCDEFGHIJKLMNOPQRSTUV").WithPadding(base32.NoPadding)
	base64URLEncoding      = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_").WithPadding(base64.NoPadding)
)
