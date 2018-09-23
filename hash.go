package bprimitives

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
)

const HashSize = 32

var (
	HashZero = Hash{0x00}
	HashOne  = Hash{0x01}
)

type Hash [HashSize]byte

func NewHash(b []byte) Hash {
	var hash Hash
	hash.SetBytes(b)
	return hash
}

func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashSize:]
	}

	copy(h[HashSize-len(b):], b)
}

func (h Hash) Equal(target Hash) bool {
	return bytes.Equal(h[:], target[:])
}

func (h Hash) IsZero() bool {
	return h.Equal(HashZero)
}

func (h Hash) Clone() Hash {
	var nh Hash
	nh.SetBytes(h.Bytes())
	return nh
}

func (h Hash) Bytes() []byte {
	return h[:]
}

func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) String() string {
	return h.Hex()
}

func (h Hash) Reverse() Hash {
	for i, j := 0, len(h)-1; i < j; i, j = i+1, j-1 {
		h[i], h[j] = h[j], h[i]
	}
	return h
}

func (h Hash) RString() string {
	return h.Reverse().Hex()
}

func Hash256(data []byte) Hash {
	h := sha256.New()
	h.Write(data)
	hash := h.Sum(nil)
	return NewHash(hash)
}

func DHash256(data []byte) Hash {
	return Hash256(Hash256(data).Bytes())
}
