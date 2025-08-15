package multihash

import (
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"

	"golang.org/x/crypto/sha3"
)

// Algorithm represents a particular multihash algorithm.
type Algorithm struct {
	code      byte
	digestLen int
	newHash   func() hash.Hash
}

// Code returns the multicodec code of the algorithm.
func (a *Algorithm) Code() byte { return a.code }

// DigestLength returns the length of the digest in bytes.
func (a *Algorithm) DigestLength() int { return a.digestLen }

// New creates a new hasher for the algorithm.
func (a *Algorithm) New() hash.Hash { return a.newHash() }

// Sum computes the multihash of the provided data and returns it prefixed
// with the multicodec code and digest length.
func (a *Algorithm) Sum(data []byte) []byte {
	h := a.newHash()
	_, _ = h.Write(data) // hash.Hash never returns an error
	digest := h.Sum(nil)
	out := make([]byte, 2+len(digest))
	out[0] = a.code
	out[1] = byte(len(digest))
	copy(out[2:], digest)
	return out
}

var (
	// SHA256 implements the sha2-256 multihash algorithm (code 0x12).
	SHA256 = &Algorithm{code: 0x12, digestLen: 32, newHash: sha256.New}

	// SHA512 implements the sha2-512 multihash algorithm (code 0x13).
	SHA512 = &Algorithm{code: 0x13, digestLen: 64, newHash: sha512.New}

	// SHA3_512 implements the sha3-512 multihash algorithm (code 0x14).
	SHA3_512 = &Algorithm{code: 0x14, digestLen: 64, newHash: sha3.New512}

	// SHA3_384 implements the sha3-384 multihash algorithm (code 0x15).
	SHA3_384 = &Algorithm{code: 0x15, digestLen: 48, newHash: sha3.New384}

	// SHA3_256 implements the sha3-256 multihash algorithm (code 0x16).
	SHA3_256 = &Algorithm{code: 0x16, digestLen: 32, newHash: sha3.New256}

	// SHA3_224 implements the sha3-224 multihash algorithm (code 0x17).
	SHA3_224 = &Algorithm{code: 0x17, digestLen: 28, newHash: sha3.New224}
)

var algorithms = map[byte]*Algorithm{
	SHA256.code:   SHA256,
	SHA512.code:   SHA512,
	SHA3_512.code: SHA3_512,
	SHA3_384.code: SHA3_384,
	SHA3_256.code: SHA3_256,
	SHA3_224.code: SHA3_224,
}

// FromBytes inspects a multihash-encoded byte sequence and returns the
// corresponding algorithm constant.
func FromBytes(b []byte) (*Algorithm, error) {
	if len(b) < 2 {
		return nil, errors.New("multihash: data too short")
	}
	code := b[0]
	length := int(b[1])

	alg, ok := algorithms[code]
	if !ok {
		return nil, errors.New("multihash: unknown code")
	}
	if length != alg.digestLen {
		return nil, errors.New("multihash: length mismatch")
	}
	if len(b)-2 < length {
		return nil, errors.New("multihash: digest too short")
	}
	return alg, nil
}
