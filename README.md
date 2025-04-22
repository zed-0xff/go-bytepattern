# bytepattern

🧬 Efficient byte pattern matcher with wildcard and repeat support — compact representation for binary scanning, patching, and signature parsing.

---

## Features

- Define byte patterns using:
  - Hex bytes (`AA`, `FF`, etc.)
  - Wildcards (`??`)
  - Repeat syntax (`[N]`) to indicate repeated elements
- Compact internal representation (`PatternEl`) — no pattern flattening
- Search (`Find`) and patch (`Patch`) arbitrary byte buffers
- String representations and conversions (ANSI, UTF-16LE/Wide)
- Graceful error handling — no panics on invalid input

---

## Example

```go
import "github.com/zed-0xff/go-bytepattern"

pattern, err := bytepattern.ParsePattern("01 ??[2] FF")
if err != nil {
	log.Fatal(err)
}

buffer := []byte{0x00, 0x01, 0x02, 0x03, 0xFF}
offset := pattern.Find(buffer)
fmt.Println("Pattern found at:", offset)

// Patch a buffer (wildcards will be skipped)
patchBuf := make([]byte, pattern.Length())
pattern.Patch(patchBuf, 0)
```

---

## Pattern Syntax

| Syntax   | Meaning                                |
|----------|----------------------------------------|
| `AA`     | Match byte 0xAA                        |
| `??`     | Match any byte (wildcard)             |
| `[N]`    | Repeat the previous byte/wildcard N times |
| `AA[4]`  | Equivalent to `AA AA AA AA`            |
| `??[3]`  | Matches 3 arbitrary bytes              |

---

## API Overview

```go
type PatternEl struct {
    Value int32 // -1 means wildcard
    Times int32
}

type Pattern struct {
    Elements []PatternEl
}

func (p *Pattern) FromHexString(s string) error
func (p *Pattern) FromArgs(args []string) error
func (p *Pattern) FromAnsiString(s string)
func (p *Pattern) FromWideString(s string)
func (p Pattern) Length() int
func (p Pattern) String() string
func (p Pattern) Find(buffer []byte) int
func (p Pattern) Patch(buffer []byte, offset int)
func ParsePattern(s string) (Pattern, error)
```
