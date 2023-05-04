// Copyright ©2023 The go-fonts Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	log.SetPrefix("anton-gen: ")
	log.SetFlags(0)

	var (
		src = flag.String(
			"src",
			"https://github.com/googlefonts/AntonFont/raw/80d0112/fonts/Anton-Regular.ttf",
			"remote file holding TTF files for Anton fonts",
		)
	)

	flag.Parse()

	tmp, err := os.MkdirTemp("", "go-fonts-anton-")
	if err != nil {
		log.Fatalf("could not create tmp dir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	var (
		fname string
	)

	switch {
	case strings.HasPrefix(*src, "http://"),
		strings.HasPrefix(*src, "https://"):
		fname, err = fetch(tmp, *src)
		if err != nil {
			log.Fatalf("could not fetch Anton sources: %+v", err)
		}
	default:
		fname = *src
	}

	f, err := os.Open(fname)
	if err != nil {
		log.Fatalf("could not open ttf file: %+v", err)
	}
	defer f.Close()

	err = gen(path.Base(fname), f)
	if err != nil {
		log.Fatalf("could not generate font: %+v", err)
	}
}

func fetch(tmp, src string) (string, error) {
	resp, err := http.Get(src)
	if err != nil {
		return "", fmt.Errorf("could not GET %q: %w", src, err)
	}
	defer resp.Body.Close()

	f, err := os.Create(path.Join(tmp, "Anton-Regular.ttf"))
	if err != nil {
		return "", fmt.Errorf("could not create ttf file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not copy ttf file: %w", err)
	}

	err = f.Close()
	if err != nil {
		return "", fmt.Errorf("could not save ttf file: %w", err)
	}

	return f.Name(), nil
}

func gen(fname string, r io.Reader) error {
	log.Printf("generating fonts package for %q...", fname)

	raw := new(bytes.Buffer)
	_, err := io.Copy(raw, r)
	if err != nil {
		return fmt.Errorf("could not download TTF file: %w", err)
	}

	err = do(fname, raw.Bytes())
	if err != nil {
		return fmt.Errorf("could not generate package for %q: %w", fname, err)
	}

	return nil
}

func do(ttfName string, src []byte) error {
	fontName := fontName(ttfName)
	pkgName := pkgName(ttfName)
	if err := os.Mkdir(pkgName, 0777); err != nil && !os.IsExist(err) {
		return fmt.Errorf("could not create package dir %q: %w", pkgName, err)
	}

	b := new(bytes.Buffer)
	fmt.Fprintf(b, "// generated by go run gen-fonts.go; DO NOT EDIT\n\n")
	fmt.Fprintf(b, "// Package %s provides the %q TrueType font\n", pkgName, fontName)
	fmt.Fprintf(b, "// from the Anton font family.\n")
	fmt.Fprintf(b, "package %[1]s // import \"github.com/go-fonts/anton/%[1]s\"\n\n", pkgName)
	fmt.Fprintf(b, "import _ \"embed\"\n")
	fmt.Fprintf(b, "// TTF is the data for the %q TrueType font.\n", fontName)
	fmt.Fprintf(b, "//\n//go:embed %s\n", ttfName)
	fmt.Fprintf(b, "var TTF  []byte\n")

	dst, err := format.Source(b.Bytes())
	if err != nil {
		return fmt.Errorf("could not format source: %w", err)
	}

	err = ioutil.WriteFile(filepath.Join(pkgName, "data.go"), dst, 0666)
	if err != nil {
		return fmt.Errorf("could not write package source file: %w", err)
	}

	err = ioutil.WriteFile(filepath.Join(pkgName, ttfName), src, 0666)
	if err != nil {
		return fmt.Errorf("could not write package TTF file: %w", err)
	}

	return nil
}

const suffix = ".ttf"

// fontName maps "Go-Regular.ttf" to "Go Regular".
func fontName(ttfName string) string {
	s := ttfName[:len(ttfName)-len(suffix)]
	s = strings.Replace(s, "-", " ", -1)
	return s
}

// pkgName maps "Go-Regular.ttf" to "goregular".
func pkgName(ttfName string) string {
	s := ttfName[:len(ttfName)-len(suffix)]
	s = strings.Replace(s, "-", "", -1)
	s = strings.ToLower(s)
	return s
}
