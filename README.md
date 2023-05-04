# anton

[![GitHub release](https://img.shields.io/github/release/go-fonts/anton.svg)](https://github.com/go-fonts/anton/releases)
[![GoDoc](https://godoc.org/github.com/go-fonts/anton?status.svg)](https://godoc.org/github.com/go-fonts/anton)
[![License](https://img.shields.io/badge/License-BSD--3-blue.svg)](https://github.com/go-fonts/anton/raw/master/LICENSE)

`anton` provides the [anton](https://github.com/googlefonts/AntonFont/) fonts as importable Go packages.

The fonts are released under the [SIL Open Font](https://github.com/go-fonts/anton/raw/master/LICENSE-SIL) license.
The Go packages under the [BSD-3](https://github.com/go-fonts/anton/raw/master/LICENSE) license.

## Example

```go
import (
	"fmt"
	"log"

	"github.com/go-fonts/anton/antonregular"
	"golang.org/x/image/font/sfnt"
)

func Example() {
	ttf, err := sfnt.Parse(antonregular.TTF)
	if err != nil {
		log.Fatalf("could not parse anton font: %+v", err)
	}

	var buf sfnt.Buffer
	v, err := ttf.Name(&buf, sfnt.NameIDVersion)
	if err != nil {
		log.Fatalf("could not retrieve font version: %+v", err)
	}

	fmt.Printf("version:    %s\n", v)
	fmt.Printf("num glyphs: %d\n", ttf.NumGlyphs())

	// Output:
	// version:    Version 2.116; ttfautohint (v1.8.3)
	// num glyphs: 1373
}
```
