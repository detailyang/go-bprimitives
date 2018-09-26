package bprimitives

import (
	"encoding/binary"
	"errors"
)

var (
	ErrBufferOverflow = errors.New("buffer: try to take more")
)

const (
	BufferDefaultSliceSize = 128
)

// Buffer little-endian
type Buffer struct {
	data []byte
	pos  int
}

func NewReadBuffer(data []byte) *Buffer {
	return &Buffer{
		data: data,
		pos:  0,
	}
}

func NewBuffer() *Buffer {
	return &Buffer{
		data: make([]byte, 0, BufferDefaultSliceSize),
		pos:  0,
	}
}

func (b *Buffer) HasSize(n int) bool {
	return b.pos+n <= len(b.data)
}

func (b *Buffer) PutUint8(u8 uint8) *Buffer {
	b.data = append(b.data, u8)
	b.pos += 1
	return b
}

func (b *Buffer) Uint8(u8 *uint8) *Buffer {
	*u8 = b.data[b.pos]
	b.pos += 1
	return b
}

func (b *Buffer) GetUint8() (uint8, error) {
	if !b.HasSize(1) {
		return 0, ErrBufferOverflow
	}

	u8 := uint8(b.data[b.pos])
	b.pos += 1

	return u8, nil
}

func (b *Buffer) PutUint16(u16 uint16) *Buffer {
	b.data = append(b.data, byte(u16), byte(u16)>>8)
	b.pos += 2
	return b
}

func (b *Buffer) Uint16(u16 *uint16) *Buffer {
	*u16 = binary.LittleEndian.Uint16(b.data[b.pos:])
	b.pos += 2
	return b
}

func (b *Buffer) GetUint16() (uint16, error) {
	if !b.HasSize(2) {
		return 0, ErrBufferOverflow
	}

	u16 := binary.LittleEndian.Uint16(b.data[b.pos:])
	b.pos += 2

	return u16, nil
}

func (b *Buffer) PutCompact(c Compact) *Buffer {
	u32 := c.Uint32()
	b.data = append(b.data, byte(u32), byte(u32>>8), byte(u32>>16), byte(u32>>24))
	b.pos += 4
	return b
}

func (b *Buffer) GetCompact() (Compact, error) {
	if !b.HasSize(CompactSize) {
		return 0, ErrBufferOverflow
	}

	var c Compact
	u32 := binary.LittleEndian.Uint32(b.data[b.pos:])
	b.pos += 4
	c.SetUint32(u32)

	return c, nil
}

func (b *Buffer) Compact(c *Compact) *Buffer {
	u32 := binary.LittleEndian.Uint32(b.data[b.pos:])
	b.pos += 4
	c.SetUint32(u32)
	return b
}

func (b *Buffer) PutUint32(u32 uint32) *Buffer {
	b.data = append(b.data, byte(u32), byte(u32>>8), byte(u32>>16), byte(u32>>24))
	b.pos += 4
	return b
}

func (b *Buffer) Uint32(u32 *uint32) *Buffer {
	*u32 = binary.LittleEndian.Uint32(b.data[b.pos:])
	b.pos += 4
	return b
}

func (b *Buffer) GetUint32() (uint32, error) {
	if !b.HasSize(4) {
		return 0, ErrBufferOverflow
	}

	u32 := binary.LittleEndian.Uint32(b.data[b.pos:])
	b.pos += 4
	return u32, nil
}

func (b *Buffer) PutUint64(u64 uint64) *Buffer {
	b.data = append(b.data, byte(u64), byte(u64>>8), byte(u64>>16), byte(u64>>24),
		byte(u64>>32), byte(u64>>40), byte(u64>>48), byte(u64>>56))
	b.pos += 8
	return b
}

func (b *Buffer) Uint64(u64 *uint64) *Buffer {
	*u64 = binary.LittleEndian.Uint64(b.data[b.pos:])
	b.pos += 8
	return b
}

func (b *Buffer) GetUint64() (uint64, error) {
	if !b.HasSize(8) {
		return 0, ErrBufferOverflow
	}

	u64 := binary.LittleEndian.Uint64(b.data[b.pos:])
	b.pos += 8

	return u64, nil
}

func (b *Buffer) PutHash(hash Hash) *Buffer {
	b.data = append(b.data, hash.Bytes()...)
	b.pos += len(hash)
	return b
}

func (b *Buffer) GetHash() (Hash, error) {
	if !b.HasSize(HashSize) {
		return HashZero, ErrBufferOverflow
	}

	var hash Hash
	data := b.data[b.pos : b.pos+HashSize]
	hash.SetBytes(data)
	b.pos += HashSize

	return hash, nil
}

func (b *Buffer) Hash(hash *Hash) *Buffer {
	data := b.data[b.pos : b.pos+HashSize]
	hash.SetBytes(data)
	b.pos += HashSize
	return b
}

func (b *Buffer) PutVarInt(n uint64) *Buffer {
	switch {
	case n < uint64(0xfd):
		b.PutUint8(uint8(n))
	case n <= 0xffff:
		b.PutUint16(uint16(n))
	case n <= 0xffffffff:
		b.PutUint32(uint32(n))
	default:
		b.PutUint64(n)
	}

	return b
}

func (b *Buffer) GetVarInt() (uint64, error) {
	var u64 uint64

	u8, err := b.GetUint8()
	if err != nil {
		return 0, err
	}

	switch u8 {
	case 0xfd:
		u16, err := b.GetUint16()
		if err != nil {
			return 0, err
		}
		u64 = uint64(u16)
	case 0xfe:
		u32, err := b.GetUint32()
		if err != nil {
			return 0, err
		}
		u64 = uint64(u32)
	case 0xff:
		u64, err = b.GetUint64()
		if err != nil {
			return 0, err
		}
	default:
		u64 = uint64(u8)
	}

	return u64, nil
}

func (b *Buffer) PutBytes(bytes []byte) *Buffer {
	b.data = append(b.data, bytes...)
	b.pos += len(bytes)
	return b
}

func (b *Buffer) GetBytes(n int) ([]byte, error) {
	if !b.HasSize(n) {
		return nil, ErrBufferOverflow
	}

	data := b.data[b.pos : b.pos+n]
	b.pos += n

	return data, nil
}

func (b *Buffer) PutVarBytes(bytes []byte) *Buffer {
	b.PutVarInt(uint64(len(bytes)))
	b.data = append(b.data, bytes...)
	b.pos += len(bytes)
	return b
}

func (b *Buffer) GetVarBytes() ([]byte, error) {
	n, err := b.GetVarInt()
	if err != nil {
		return nil, err
	}

	return b.GetBytes(int(n))
}

func (b *Buffer) PutCheckSum() *Buffer {
	return b.PutBytes(DHash256(b.data).TakeBytes(0, 4))
}

func (b *Buffer) Size() int {
	return b.pos
}

func (b *Buffer) Bytes() []byte {
	return b.data
}
