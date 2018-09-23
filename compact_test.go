package bprimitives

import (
	"bytes"
	"testing"
)

func TestCompactBytes(t *testing.T) {
	tests := []struct {
		input  uint32
		expect []byte
	}{
		{
			0x01003456,
			[]byte{},
		},

		{
			0x01123456,
			[]byte{0x12},
		},
		{
			0x02008000,
			[]byte{0x80},
		},
		{
			0x05009234,
			[]byte{0x92, 0x34, 0x00, 0x00},
		},
		{
			0x04123456,
			[]byte{0x12, 0x34, 0x56, 0x00},
		},
	}
	for _, test := range tests {
		b := NewCompact(test.input).Bytes()
		if !bytes.Equal(b, test.expect) {
			t.Fatalf("expect %v got %v", test.expect, b)
		}
	}
}

func TestNewCompactFromBytes(t *testing.T) {
	tests := []struct {
		input  []byte
		expect uint32
	}{
		{
			[]byte{0x00},
			0,
		},
		{
			[]byte{1},
			65536,
		},
	}

	for _, test := range tests {
		c := NewCompactFromBytes(test.input)
		if c.Uint32() != test.expect {
			t.Fatalf("expect %v got %v", test.expect, c)
		}
	}
}
