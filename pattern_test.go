package bytepattern

import (
	"testing"
)

func TestFromHexString_Valid(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		length   int
	}{
		{"AA", "AA", 1},
		{"AA[3]", "AA[3]", 3},
		{"AA ?? BB", "AA ?? BB", 3},
		{"??[2] 01", "??[2] 01", 3},
		{"12 ??[3] FF", "12 ??[3] FF", 5},
	}

	for _, test := range tests {
		var p Pattern
		err := p.FromHexString(test.input)
		if err != nil {
			t.Errorf("unexpected error for input %q: %v", test.input, err)
			continue
		}
		if p.String() != test.expected {
			t.Errorf("got %q, want %q", p.String(), test.expected)
		}
		if p.Length() != test.length {
			t.Errorf("got length %d, want %d", p.Length(), test.length)
		}
	}
}

func TestFromHexString_Invalid(t *testing.T) {
	invalid := []string{
		"[3]",     // no element before [3]
		"??[x]",   // non-numeric repeat
		"AA[0]",   // zero repeat
		"12[-1]",  // negative repeat
		"AA BB[2", // missing closing bracket
		"A",       // incomplete byte
		"GG",      // invalid hex
		"AA<",     // invalid character
	}

	for _, s := range invalid {
		var p Pattern
		if err := p.FromHexString(s); err == nil {
			t.Errorf("expected error for input %q", s)
		}
	}
}

func TestFind(t *testing.T) {
	var p Pattern
	_ = p.FromHexString("01 02 ??[2] 05")
	buffer := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	pos := p.Find(buffer)
	if pos != 1 {
		t.Errorf("Find returned %d, want 1", pos)
	}
}

func TestFind_NotFound(t *testing.T) {
	var p Pattern
	_ = p.FromHexString("AA BB CC")
	buf := []byte{0x00, 0x01, 0x02}
	if p.Find(buf) != -1 {
		t.Error("expected -1 for non-matching buffer")
	}
}

func TestPatch(t *testing.T) {
	var p Pattern
	_ = p.FromHexString("AA ??[2] BB")
	buf := make([]byte, p.Length())
	p.Patch(buf, 0)

	expected := []byte{0xAA, 0x00, 0x00, 0xBB}
	for i, b := range expected {
		if buf[i] != b {
			t.Errorf("Patch mismatch at %d: got %02X, want %02X", i, buf[i], b)
		}
	}
}

func TestFromAnsiString(t *testing.T) {
	var p Pattern
	p.FromAnsiString("ABC")
	want := "41 42 43"
	if p.String() != want {
		t.Errorf("FromAnsiString: got %q, want %q", p.String(), want)
	}
}

func TestFromWideString(t *testing.T) {
	var p Pattern
	p.FromWideString("AB")
	want := "41 00 42 00"
	if p.String() != want {
		t.Errorf("FromWideString: got %q, want %q", p.String(), want)
	}
}

func TestParsePattern(t *testing.T) {
	p, err := ParsePattern("01 ??[2] FF")
	if err != nil {
		t.Fatalf("ParsePattern failed: %v", err)
	}
	if p.String() != "01 ??[2] FF" {
		t.Errorf("ParsePattern: got %q", p.String())
	}
}

func TestPatternLength(t *testing.T) {
	p := Pattern{
		elements: []PatternEl{
			{value: 0x01, times: 2},
			{value: -1, times: 3},
		},
	}
	if p.Length() != 5 {
		t.Errorf("Length: got %d, want 5", p.Length())
	}
}
