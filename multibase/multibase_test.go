package multibase

import (
	"bytes"
	"testing"
)

type encodingTest struct {
	name string
	enc  Encoding
}

func TestRoundtrip(t *testing.T) {
	data := []byte("hello multibase")
	tests := []encodingTest{
		{"base2", Base2},
		{"base8", Base8},
		{"base10", Base10},
		{"base16", Base16},
		{"base16upper", Base16Upper},
		{"base32", Base32},
		{"base32upper", Base32Upper},
		{"base32hex", Base32Hex},
		{"base32hexupper", Base32HexUpper},
		{"base36", Base36},
		{"base36upper", Base36Upper},
		{"base58btc", Base58BTC},
		{"base64url", Base64URL},
	}

	for _, tc := range tests {
		encoded := tc.enc.Encode(data)
		if rune(encoded[0]) != tc.enc.code {
			t.Fatalf("%s: missing prefix", tc.name)
		}
		decoded, err := Decode(encoded)
		if err != nil {
			t.Fatalf("%s: decode error: %v", tc.name, err)
		}
		if !bytes.Equal(decoded, data) {
			t.Fatalf("%s: roundtrip mismatch", tc.name)
		}
	}
}

func TestDecodeErrors(t *testing.T) {
	if _, err := Decode(""); err == nil {
		t.Fatalf("expected error for empty input")
	}
	if _, err := Decode("x123"); err == nil {
		t.Fatalf("expected error for unknown prefix")
	}
}
