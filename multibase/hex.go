package multibase

import "encoding/hex"

type hexEncoder struct {
	upper bool
}

func (e *hexEncoder) Encode(dst, src []byte) {
	hex.Encode(dst, src)

	// hex encodes in lower case by default
	if e.upper {
		for i := range dst {
			dst[i] = toUpperASCIIByte(dst[i])
		}
	}
}

func (e *hexEncoder) EncodedLen(n int) int {
	return hex.EncodedLen(n)
}

func toUpperASCIIByte(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b &^ 0x20
	}
	return b
}
