// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	goformat "go/format"
	"go/parser"
	"go/token"
	"testing"
)

// TestRewriteImpossibleSubstitution checks that a rewrite whose replacement
// cannot be assigned into every matched position does not panic. The rule
// a -> p.a matches the declaration name _, but a qualified name cannot
// replace an identifier in that position, so it is left unchanged while the
// type int is rewritten to p.int. See go.dev/issue/52145.
func TestRewriteImpossibleSubstitution(t *testing.T) {
	const src = "package p\n\nvar _ int\n"
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	file = rewriteFile(fset, parseExpr("a", "pattern"), parseExpr("p.a", "replacement"), file)

	var buf bytes.Buffer
	if err := goformat.Node(&buf, fset, file); err != nil {
		t.Fatal(err)
	}
	const want = "package p\n\nvar _ p.int\n"
	if got := buf.String(); got != want {
		t.Errorf("rewrite result = %q, want %q", got, want)
	}
}
