package multibase

import (
	"github.com/pki-directory/cid/multibase/base16"
)

var (
	Base16 = &Encoding{
		prefix:    prefixBase16,
		name:      "base16",
		encode:    base16.Encode,
		encodeLen: base16.EncodedLen,
	}

	Base16Upper = &Encoding{
		prefix:    prefixBase16Upper,
		name:      "base16upper",
		encode:    base16.EncodeUpper,
		encodeLen: base16.EncodedLen,
	}

	Base32 = &Encoding{
		prefix:    prefixBase32,
		name:      "base32",
		encode:    base32Encoding.Encode,
		encodeLen: base32Encoding.EncodedLen,
	}

	Base32Upper = &Encoding{
		prefix:    prefixBase32Upper,
		name:      "base32upper",
		encode:    base32UpperEncoding.Encode,
		encodeLen: base32UpperEncoding.EncodedLen,
	}

	Base32Hex = &Encoding{
		prefix:    prefixBase32Hex,
		name:      "base32hex",
		encode:    base32HexEncoding.Encode,
		encodeLen: base32HexEncoding.EncodedLen,
	}

	Base32HexUpper = &Encoding{
		prefix:    prefixBase32HexUpper,
		name:      "base32hexupper",
		encode:    base32HexUpperEncoding.Encode,
		encodeLen: base32HexUpperEncoding.EncodedLen,
	}

	Base64URL = &Encoding{
		prefix:    prefixBase64URL,
		name:      "base64url",
		encode:    base64URLEncoding.Encode,
		encodeLen: base64URLEncoding.EncodedLen,
	}
)

type Encoding struct {
	prefix    byte // change to rune if non-ASCII encodings are added
	name      string
	encode    func(dst, src []byte)
	encodeLen func(n int) int
}

func (e *Encoding) Prefix() byte {
	return e.prefix
}

func (e *Encoding) String() string {
	return e.name
}

func (e *Encoding) Encode(data []byte) string {
	dst := make([]byte, e.encodeLen(len(data))+1)
	dst[0] = e.prefix
	e.encode(dst[1:], data)
	return string(dst)
}
