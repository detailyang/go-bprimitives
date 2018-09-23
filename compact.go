package bprimitives

import (
	"math/big"
)

const (
	CompactSize = 4
)

type Compact uint32

func NewCompact(u32 uint32) Compact {
	return Compact(u32)
}

func NewCompactFromBytes(data []byte) Compact {
	bg := big.NewInt(0)
	bg = bg.SetBytes(data)

	if bg.Sign() == 0 {
		return NewCompact(0)
	}

	mantissa := uint32(0)
	exponent := uint8(len(bg.Bytes()))
	if exponent <= 3 {
		mantissa = uint32(bg.Bits()[0])
		mantissa <<= 8 * (3 - exponent)
	} else {
		tmp := new(big.Int).Set(bg)
		mantissa = uint32(tmp.Rsh(tmp, uint(8*(exponent-3))).Bits()[0])
	}

	// When the mantissa already has the sign bit set, the number is too
	// large to fit into the available 23-bits, so divide the number by 256
	// and increment the exponent accordingly.
	if mantissa&0x00800000 != 0 {
		mantissa >>= 8
		exponent++
	}

	c := NewCompact(uint32(exponent<<24) | mantissa)
	if bg.Sign() < 0 {
		c |= 0x00800000
	}

	return c
}

//Bytes
// Like IEEE754 floating point, there are three basic components: the sign,
// the exponent, and the mantissa.  They are broken out as follows:
//
//	* the most significant 8 bits represent the unsigned base 256 exponent
// 	* bit 23 (the 24th bit) represents the sign bit
//	* the least significant 23 bits represent the mantissa
//
//	-------------------------------------------------
//	|   Exponent     |    Sign    |    Mantissa     |
//	-------------------------------------------------
//	| 8 bits [31-24] | 1 bit [23] | 23 bits [22-00] |
//	-------------------------------------------------
//
// The formula to calculate N is:
// 	N = (-1^sign) * mantissa * 256^(exponent-3)
func (c Compact) Bytes() []byte {
	var bg *big.Int

	mantissa := c & 0x007FFFFF
	isNegative := c&0x00800000 != 0
	exponent := c >> 24

	if exponent <= 3 {
		mantissa >>= 8 * (3 - exponent)
		bg = big.NewInt(int64(mantissa))
	} else {
		bg = big.NewInt(int64(mantissa))
		bg.Lsh(bg, uint(8*(exponent-3)))
	}

	if isNegative {
		bg = bg.Neg(bg)
	}

	return bg.Bytes()
}

func (c *Compact) SetUint32(u32 uint32) {
	*c = NewCompact(u32)
}

func (c Compact) Uint32() uint32 {
	return uint32(c)
}
