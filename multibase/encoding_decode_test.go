package multibase

import (
	"bytes"
	"testing"
)

var encodings = []*Encoding{
	Base16,
	Base16Upper,
	Base32,
	Base32Upper,
	Base32Hex,
	Base32HexUpper,
	Base64URL,
}

func FuzzRoundtrip(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0})
	f.Add([]byte("multibase"))
	f.Add([]byte("Hello, 世界"))
	f.Add(bytes.Repeat([]byte{0xff}, 1024))

	f.Fuzz(func(t *testing.T, data []byte) {
		t.Logf("data: %x", data)
		for _, enc := range encodings {
			name := enc.name
			encoded := enc.Encode(data)

			if len(encoded) == 0 {
				t.Fatalf("%s: missing prefix and data. data: %x", name, data)
			}

			if encoded[0] != enc.prefix {
				t.Fatalf("%s: incorrect prefix. data: %x", name, data)
			}

			decoded, err := Decode(encoded)
			if err != nil {
				t.Fatalf("%s: decode error: %v. data: %x", name, err, data)
			}

			if !bytes.Equal(decoded, data) {
				t.Fatalf("%s: roundtrip mismatch. got: %x, want: %x.", name, decoded, data)
			}

			t.Logf("[%s] %s", enc.name, encoded)
		}
	})
}
