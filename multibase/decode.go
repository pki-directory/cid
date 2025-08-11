package multibase

import (
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	ErrEmptyString                  = errors.New("empty string")
	ErrUnsupportedMultibaseEncoding = errors.New("unsupported multibase encoding")
)

// Decode decodes a multibase encoded string.
func Decode(s string) ([]byte, error) {
	if len(s) == 0 {
		return nil, ErrEmptyString
	}

	// Note: if non-ASCII encodings are added, we need to use
	// utf8.DecodeRuneInString instead.
	switch s[0] {
	case prefixBase16, prefixBase16Upper:
		return hex.DecodeString(s[1:])
	case prefixBase32:
		return base32Encoding.DecodeString(s[1:])
	case prefixBase32Upper:
		return base32UpperEncoding.DecodeString(s[1:])
	case prefixBase32Hex:
		return base32HexEncoding.DecodeString(s[1:])
	case prefixBase32HexUpper:
		return base32HexUpperEncoding.DecodeString(s[1:])
	case prefixBase64URL:
		return base64URLEncoding.DecodeString(s[1:])
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedMultibaseEncoding, s[0])
	}
}
