package multibase

var (
	Base16 = &Encoding{
		prefix:  prefixBase16,
		name:    "base16",
		encoder: base16Encoding,
	}

	Base16Upper = &Encoding{
		prefix:  prefixBase16Upper,
		name:    "base16upper",
		encoder: base16UpperEncoding,
	}

	Base32 = &Encoding{
		prefix:  prefixBase32,
		name:    "base32",
		encoder: base32Encoding,
	}

	Base32Upper = &Encoding{
		prefix:  prefixBase32Upper,
		name:    "base32upper",
		encoder: base32UpperEncoding,
	}

	Base32Hex = &Encoding{
		prefix:  prefixBase32Hex,
		name:    "base32hex",
		encoder: base32HexEncoding,
	}

	Base32HexUpper = &Encoding{
		prefix:  prefixBase32HexUpper,
		name:    "base32hexupper",
		encoder: base32HexUpperEncoding,
	}

	Base64URL = &Encoding{
		prefix:  prefixBase64URL,
		name:    "base64url",
		encoder: base64URLEncoding,
	}
)

type encoder interface {
	Encode(dst, src []byte)
	EncodedLen(n int) int
}

type Encoding struct {
	prefix  byte // change to rune if non-ASCII encodings are added
	name    string
	encoder encoder
}

func (e *Encoding) String() string {
	return e.name
}

func (e *Encoding) Encode(data []byte) string {
	if e.encoder == nil {
		panic("multibase: zero value encoding")
	}

	dst := make([]byte, e.encoder.EncodedLen(len(data))+1)
	dst[0] = e.prefix
	e.encoder.Encode(dst[1:], data)
	return string(dst)
}
