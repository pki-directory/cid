package multibase

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"unicode/utf8"
)

var (
	ErrEmptyString                = errors.New("empty string")
	ErrUnsupportedMultibasePrefix = errors.New("unsupported multibase prefix")
)

// Encoding represents a multibase encoding. It knows how to encode and decode
// data and the corresponding multibase prefix code.
type Encoding struct {
	// code is the multibase prefix identifying this encoding.
	code rune

	encode func([]byte) string
	decode func(string) ([]byte, error)
}

// Encode encodes the given byte slice and prepends the multibase prefix.
func (e Encoding) Encode(src []byte) string {
	return string(e.code) + e.encode(src)
}

// Decode decodes the given string using the encoding. The input string should
// not include the multibase prefix.
func (e Encoding) Decode(s string) ([]byte, error) {
	return e.decode(s)
}

// -----------------------------------------------------------------------------
// Base specific implementations

// helper functions for numeric bases using big.Int
func encodeBigInt(src []byte, base int) string {
	if len(src) == 0 {
		return ""
	}
	i := new(big.Int).SetBytes(src)
	return i.Text(base)
}

func decodeBigInt(s string, base int) ([]byte, error) {
	if len(s) == 0 {
		return []byte{}, nil
	}
	i := new(big.Int)
	if _, ok := i.SetString(s, base); !ok {
		return nil, fmt.Errorf("invalid base%d string", base)
	}
	return i.Bytes(), nil
}

// Binary (base2) implementation encodes each byte as 8 bits.
func encodeBase2(src []byte) string {
	if len(src) == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(len(src) * 8)
	for _, by := range src {
		for i := 7; i >= 0; i-- {
			if (by>>uint(i))&1 == 1 {
				b.WriteByte('1')
			} else {
				b.WriteByte('0')
			}
		}
	}
	return b.String()
}

func decodeBase2(s string) ([]byte, error) {
	if len(s)%8 != 0 {
		return nil, errors.New("invalid binary string length")
	}
	out := make([]byte, len(s)/8)
	for i := 0; i < len(out); i++ {
		var v byte
		for j := 0; j < 8; j++ {
			c := s[i*8+j]
			switch c {
			case '0':
			case '1':
				v |= 1 << uint(7-j)
			default:
				return nil, fmt.Errorf("invalid binary digit %q", c)
			}
		}
		out[i] = v
	}
	return out, nil
}

// Base16 lower/upper helpers
func encodeBase16Lower(src []byte) string   { return hex.EncodeToString(src) }
func encodeBase16Upper(src []byte) string   { return strings.ToUpper(hex.EncodeToString(src)) }
func decodeBase16(s string) ([]byte, error) { return hex.DecodeString(s) }

// Base32 encoders
var (
	b32LowerNoPad    = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567").WithPadding(base32.NoPadding)
	b32UpperNoPad    = base32.StdEncoding.WithPadding(base32.NoPadding)
	b32HexLowerNoPad = base32.NewEncoding("0123456789abcdefghijklmnopqrstuv").WithPadding(base32.NoPadding)
	b32HexUpperNoPad = base32.HexEncoding.WithPadding(base32.NoPadding)
)

func encodeBase32Lower(src []byte) string { return b32LowerNoPad.EncodeToString(src) }
func decodeBase32Lower(s string) ([]byte, error) {
	return b32LowerNoPad.DecodeString(strings.ToLower(s))
}

func encodeBase32Upper(src []byte) string { return b32UpperNoPad.EncodeToString(src) }
func decodeBase32Upper(s string) ([]byte, error) {
	return b32UpperNoPad.DecodeString(strings.ToUpper(s))
}

func encodeBase32HexLower(src []byte) string { return b32HexLowerNoPad.EncodeToString(src) }
func decodeBase32HexLower(s string) ([]byte, error) {
	return b32HexLowerNoPad.DecodeString(strings.ToLower(s))
}

func encodeBase32HexUpper(src []byte) string { return b32HexUpperNoPad.EncodeToString(src) }
func decodeBase32HexUpper(s string) ([]byte, error) {
	return b32HexUpperNoPad.DecodeString(strings.ToUpper(s))
}

// Base64 URL without padding
func encodeBase64URL(src []byte) string        { return base64.RawURLEncoding.EncodeToString(src) }
func decodeBase64URL(s string) ([]byte, error) { return base64.RawURLEncoding.DecodeString(s) }

// Base58 BTC alphabet
const btcAlphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func encodeBase58BTC(src []byte) string {
	if len(src) == 0 {
		return ""
	}
	x := new(big.Int).SetBytes(src)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	var out []byte
	mod := new(big.Int)
	for x.Cmp(zero) > 0 {
		x.DivMod(x, base, mod)
		out = append(out, btcAlphabet[mod.Int64()])
	}
	// handle leading zeros
	for _, b := range src {
		if b == 0 {
			out = append(out, btcAlphabet[0])
		} else {
			break
		}
	}
	// reverse
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return string(out)
}

func decodeBase58BTC(s string) ([]byte, error) {
	x := big.NewInt(0)
	base := big.NewInt(58)
	for i := 0; i < len(s); i++ {
		idx := strings.IndexByte(btcAlphabet, s[i])
		if idx < 0 {
			return nil, fmt.Errorf("invalid base58 character %q", s[i])
		}
		x.Mul(x, base)
		x.Add(x, big.NewInt(int64(idx)))
	}
	bytes := x.Bytes()
	// add leading zeros
	nLeading := 0
	for nLeading < len(s) && s[nLeading] == btcAlphabet[0] {
		nLeading++
	}
	if nLeading > 0 {
		bytes = append(make([]byte, nLeading), bytes...)
	}
	return bytes, nil
}

// Numeric base helpers
func encodeBase8(src []byte) string              { return encodeBigInt(src, 8) }
func decodeBase8(s string) ([]byte, error)       { return decodeBigInt(s, 8) }
func encodeBase10(src []byte) string             { return encodeBigInt(src, 10) }
func decodeBase10(s string) ([]byte, error)      { return decodeBigInt(s, 10) }
func encodeBase36Lower(src []byte) string        { return encodeBigInt(src, 36) }
func decodeBase36Lower(s string) ([]byte, error) { return decodeBigInt(strings.ToLower(s), 36) }
func encodeBase36Upper(src []byte) string        { return strings.ToUpper(encodeBigInt(src, 36)) }
func decodeBase36Upper(s string) ([]byte, error) { return decodeBigInt(strings.ToLower(s), 36) }

// -----------------------------------------------------------------------------

// Exported encodings
var (
	Base2          = Encoding{code: '0', encode: encodeBase2, decode: decodeBase2}
	Base8          = Encoding{code: '7', encode: encodeBase8, decode: decodeBase8}
	Base10         = Encoding{code: '9', encode: encodeBase10, decode: decodeBase10}
	Base16         = Encoding{code: 'f', encode: encodeBase16Lower, decode: decodeBase16}
	Base16Upper    = Encoding{code: 'F', encode: encodeBase16Upper, decode: decodeBase16}
	Base32         = Encoding{code: 'b', encode: encodeBase32Lower, decode: decodeBase32Lower}
	Base32Upper    = Encoding{code: 'B', encode: encodeBase32Upper, decode: decodeBase32Upper}
	Base32Hex      = Encoding{code: 'v', encode: encodeBase32HexLower, decode: decodeBase32HexLower}
	Base32HexUpper = Encoding{code: 'V', encode: encodeBase32HexUpper, decode: decodeBase32HexUpper}
	Base36         = Encoding{code: 'k', encode: encodeBase36Lower, decode: decodeBase36Lower}
	Base36Upper    = Encoding{code: 'K', encode: encodeBase36Upper, decode: decodeBase36Upper}
	Base58BTC      = Encoding{code: 'z', encode: encodeBase58BTC, decode: decodeBase58BTC}
	Base64URL      = Encoding{code: 'u', encode: encodeBase64URL, decode: decodeBase64URL}
)

// Decode decodes a multibase encoded string.
func Decode(s string) ([]byte, error) {
	if len(s) == 0 {
		return nil, ErrEmptyString
	}
	r, size := utf8.DecodeRuneInString(s)
	switch r {
	case Base2.code:
		return Base2.decode(s[size:])
	case Base8.code:
		return Base8.decode(s[size:])
	case Base10.code:
		return Base10.decode(s[size:])
	case Base16.code:
		return Base16.decode(s[size:])
	case Base16Upper.code:
		return Base16Upper.decode(s[size:])
	case Base32.code:
		return Base32.decode(s[size:])
	case Base32Upper.code:
		return Base32Upper.decode(s[size:])
	case Base32Hex.code:
		return Base32Hex.decode(s[size:])
	case Base32HexUpper.code:
		return Base32HexUpper.decode(s[size:])
	case Base36.code:
		return Base36.decode(s[size:])
	case Base36Upper.code:
		return Base36Upper.decode(s[size:])
	case Base58BTC.code:
		return Base58BTC.decode(s[size:])
	case Base64URL.code:
		return Base64URL.decode(s[size:])
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedMultibasePrefix, r)
	}
}
