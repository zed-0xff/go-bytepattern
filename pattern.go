package bytepattern

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type PatternEl struct {
	value int32 // -1 means wildcard
	times int32
}

type Pattern struct {
	elements []PatternEl
}

func (p Pattern) Length() int {
	total := 0
	for _, el := range p.elements {
		total += int(el.times)
	}
	return total
}

func (p Pattern) String() string {
	var sb strings.Builder
	for _, el := range p.elements {
		var val string
		if el.value == -1 {
			val = "??"
		} else {
			val = fmt.Sprintf("%02X", el.value)
		}
		if el.times > 1 {
			val += fmt.Sprintf("[%d]", el.times)
		}
		sb.WriteString(val + " ")
	}
	return strings.TrimSpace(sb.String())
}

func (p Pattern) Find(buffer []byte) int {
	patLen := p.Length()
	if patLen == 0 || len(buffer) < patLen {
		return -1
	}

	for i := 0; i <= len(buffer)-patLen; i++ {
		bufIdx := i
		match := true
		for _, el := range p.elements {
			for k := 0; k < int(el.times); k++ {
				b := buffer[bufIdx]
				if el.value != -1 && el.value != int32(b) {
					match = false
					break
				}
				bufIdx++
			}
			if !match {
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func (p Pattern) Patch(buffer []byte, offset int) {
	bufIdx := offset
	for _, el := range p.elements {
		for i := 0; i < int(el.times); i++ {
			if el.value != -1 {
				buffer[bufIdx] = byte(el.value)
			}
			bufIdx++
		}
	}
}

// no wildcards
func (p *Pattern) FromBytes(input []byte) {
	p.elements = make([]PatternEl, 0, len(input))
	for _, b := range input {
		p.elements = append(p.elements, PatternEl{value: int32(b), times: 1})
	}
}

func (p *Pattern) FromHexString(input string) error {
	p.elements = nil
	tokens, err := tokenize(input)
	if err != nil {
		return err
	}

	lastIdx := -1
	for _, tok := range tokens {
		if strings.HasPrefix(tok, "[") {
			if lastIdx == -1 {
				return errors.New("repeat syntax used without a preceding element")
			}
			n, err := strconv.Atoi(tok[1 : len(tok)-1])
			if err != nil || n < 1 {
				return fmt.Errorf("invalid repeat value: %q", tok)
			}
			p.elements[lastIdx].times = int32(n)
			continue
		}

		var val int32
		if tok == "??" {
			val = -1
		} else {
			b, err := strconv.ParseUint(tok, 16, 8)
			if err != nil {
				return fmt.Errorf("invalid hex token: %q", tok)
			}
			val = int32(b)
		}
		p.elements = append(p.elements, PatternEl{value: val, times: 1})
		lastIdx++
	}
	return nil
}

func (p *Pattern) FromArgs(args []string) error {
	return p.FromHexString(strings.Join(args, ""))
}

func (p *Pattern) FromAnsiString(s string) {
	p.elements = make([]PatternEl, 0, len(s))
	for _, c := range s {
		p.elements = append(p.elements, PatternEl{value: int32(c), times: 1})
	}
}

func (p *Pattern) FromWideString(s string) {
	p.elements = make([]PatternEl, 0, len(s)*2)
	for _, c := range s {
		p.elements = append(p.elements, PatternEl{value: int32(c), times: 1})
		p.elements = append(p.elements, PatternEl{value: 0, times: 1})
	}
}

func Parse(s string) (Pattern, error) {
	var p Pattern
	err := p.FromHexString(s)
	return p, err
}

func tokenize(input string) ([]string, error) {
	input = strings.ReplaceAll(input, " ", "")
	var tokens []string
	for i := 0; i < len(input); {
		if i+1 < len(input) && input[i] == '?' && input[i+1] == '?' {
			tokens = append(tokens, "??")
			i += 2
		} else if input[i] == '[' {
			j := i + 1
			for j < len(input) && input[j] != ']' {
				j++
			}
			if j == len(input) {
				return nil, fmt.Errorf("unclosed [N] repeat block at position %d", i)
			}
			tokens = append(tokens, input[i:j+1])
			i = j + 1
		} else {
			if i+2 > len(input) {
				return nil, fmt.Errorf("incomplete hex byte at position %d", i)
			}
			tokens = append(tokens, input[i:i+2])
			i += 2
		}
	}
	return tokens, nil
}
