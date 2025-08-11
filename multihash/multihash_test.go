package multihash_test

import (
	"bytes"
	"testing"

	"crypto/sha256"
	"crypto/sha512"

	"golang.org/x/crypto/sha3"

	"github.com/pki-directory/cid/multihash"
)

func TestSum(t *testing.T) {
	data := []byte("hello world")

	expected256 := sha256.Sum256(data)
	expected512 := sha512.Sum512(data)
	expected3_512 := sha3.Sum512(data)
	expected3_384 := sha3.Sum384(data)
	expected3_256 := sha3.Sum256(data)
	expected3_224 := sha3.Sum224(data)

	tests := []struct {
		alg    *multihash.Algorithm
		digest []byte
	}{
		{multihash.SHA256, expected256[:]},
		{multihash.SHA512, expected512[:]},
		{multihash.SHA3_512, expected3_512[:]},
		{multihash.SHA3_384, expected3_384[:]},
		{multihash.SHA3_256, expected3_256[:]},
		{multihash.SHA3_224, expected3_224[:]},
	}

	for _, tc := range tests {
		got := tc.alg.Sum(data)
		if got[0] != tc.alg.Code() {
			t.Errorf("code mismatch: got 0x%x want 0x%x", got[0], tc.alg.Code())
		}
		if int(got[1]) != tc.alg.DigestLength() {
			t.Errorf("length mismatch for code 0x%x: got %d want %d", tc.alg.Code(), got[1], tc.alg.DigestLength())
		}
		if !bytes.Equal(got[2:], tc.digest) {
			t.Errorf("digest mismatch for code 0x%x", tc.alg.Code())
		}
	}
}

func TestNew(t *testing.T) {
	data := []byte("a different input")

	expected256 := sha256.Sum256(data)
	expected512 := sha512.Sum512(data)
	expected3_512 := sha3.Sum512(data)
	expected3_384 := sha3.Sum384(data)
	expected3_256 := sha3.Sum256(data)
	expected3_224 := sha3.Sum224(data)

	tests := []struct {
		alg    *multihash.Algorithm
		digest []byte
	}{
		{multihash.SHA256, expected256[:]},
		{multihash.SHA512, expected512[:]},
		{multihash.SHA3_512, expected3_512[:]},
		{multihash.SHA3_384, expected3_384[:]},
		{multihash.SHA3_256, expected3_256[:]},
		{multihash.SHA3_224, expected3_224[:]},
	}

	for _, tc := range tests {
		h := tc.alg.New()
		_, _ = h.Write(data)
		got := h.Sum(nil)
		if !bytes.Equal(got, tc.digest) {
			t.Errorf("digest mismatch for code 0x%x", tc.alg.Code())
		}
	}
}

func TestFromBytes(t *testing.T) {
	data := []byte("data for multihash")

	algs := []*multihash.Algorithm{
		multihash.SHA256,
		multihash.SHA512,
		multihash.SHA3_512,
		multihash.SHA3_384,
		multihash.SHA3_256,
		multihash.SHA3_224,
	}

	for _, alg := range algs {
		digest := alg.Sum(data)
		a, err := multihash.FromBytes(digest)
		if err != nil {
			t.Fatalf("unexpected error for code 0x%x: %v", alg.Code(), err)
		}
		if a != alg {
			t.Fatalf("expected algorithm pointer %p got %p", alg, a)
		}
	}
}

func TestFromBytesErrors(t *testing.T) {
	// Unknown code
	if _, err := multihash.FromBytes([]byte{0xff, 0x20}); err == nil {
		t.Errorf("expected error for unknown code")
	}

	// Length mismatch
	dig := multihash.SHA256.Sum([]byte("test"))
	dig[1] = 0 // corrupt length
	if _, err := multihash.FromBytes(dig); err == nil {
		t.Errorf("expected error for length mismatch")
	}

	// Too short data
	if _, err := multihash.FromBytes([]byte{0x12}); err == nil {
		t.Errorf("expected error for short data")
	}

	// Truncated digest
	dig = multihash.SHA256.Sum([]byte("test"))
	dig = dig[:len(dig)-1]
	if _, err := multihash.FromBytes(dig); err == nil {
		t.Errorf("expected error for truncated digest")
	}
}
