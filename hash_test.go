package bprimitives

import (
	"testing"
)

func TestDHash256(t *testing.T) {
	hash := DHash256([]byte(string("hello")))
	if hash.String() != "503d8319a48348cdc610a582f7bf754b5833df65038606eb48510790dfc99595" {
		t.Fatalf("expect :%s", "503d8319a48348cdc610a582f7bf754b5833df65038606eb48510790dfc99595")
	}

}

func TestReverseHexString(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"1F", "1F"},
		{"1F2F", "2F1F"},
		{"123456", "563412"},
		{"1F2F3F", "3F2F1F"},
	}

	for _, test := range tests {
		r := reverseHexString(test.input)
		if r != test.expect {
			t.Errorf("expect %s got %s", test.expect, r)
		}
	}
}

func TestHashReverse(t *testing.T) {
	rh := Hash{
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x12,
	}
	h := Hash{0x12}
	if !h.Reverse().Equal(rh) {
		t.Fatalf("expect Hash{0x12}")
	}
}
