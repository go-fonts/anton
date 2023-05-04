// Copyright ©2020 The go-fonts Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anton_test

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
