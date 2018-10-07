package bprimitives

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"math"
	"math/big"
	"math/rand"
	"time"

	"golang.org/x/crypto/ripemd160"
)

const HashSize = 32

var (
	HashZero = Hash{0x00}
	HashOne  = Hash{0x01}
)

type Hash [HashSize]byte
type H256 [256]byte
type H160 [160]byte

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func NewHash(b []byte) Hash {
	var hash Hash
	hash.SetBytes(b)
	return hash
}

func NewHashFromHexString(hexstring string) (Hash, error) {
	data, err := hex.DecodeString(hexstring)
	if err != nil {
		return HashZero, err
	}

	return NewHash(data), nil
}

func NewHashFromReversedHexString(hexstring string) (Hash, error) {
	return NewHashFromHexString(reverseHexString(hexstring))
}

func fallbackRandomBytes(n int) []byte {
	data := make([]byte, n)
	for i := 0; i < n; i++ {
		data[i] = byte(rand.Intn(int(math.MaxInt8) + 1))
	}

	return data
}

func reverseHexString(s string) string {
	r := []rune(s)
	ns := len(s)
	middle := ns / 2
	if middle%2 == 1 {
		middle = middle - 1
	}

	for i := 0; i < middle; i += 2 {
		r[i], r[ns-i-1-1] = r[ns-i-1-1], r[i]
		r[i+1], r[ns-i-1] = r[ns-i-1], r[i+1]
	}

	return string(r)
}

func NewRandomHash() Hash {
	data := make([]byte, HashSize)
	_, err := rand.Read(data)
	if err != nil {
		data = fallbackRandomBytes(HashSize)
	}

	return NewHash(data)
}

func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashSize:]
	}

	copy(h[HashSize-len(b):], b)
}

func (h Hash) TakeBytes(b, e int) []byte {
	return h[b:e]
}

func (h Hash) Cmp(target Hash) int {
	var a, b big.Int
	a.SetBytes(h[:])
	b.SetBytes(target[:])

	return a.Cmp(&b)
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
	return h.Reverse().Hex()
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

type ByHash []Hash

func (a ByHash) Len() int           { return len(a) }
func (a ByHash) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHash) Less(i, j int) bool { return a[i].Cmp(a[j]) < 0 }

func Hash160(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	data = hash.Sum(nil)
	hash = ripemd160.New()
	hash.Write(data)
	data = hash.Sum(nil)

	return data
}

func Hash256(data []byte) Hash {
	hash := sha256.New()
	hash.Write(data)
	data = hash.Sum(nil)

	return NewHash(data)
}

func DHash256(data []byte) Hash {
	return Hash256(Hash256(data).Bytes())
}
