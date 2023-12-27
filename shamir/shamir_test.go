/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package shamir

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestSplitInvalid(t *testing.T) {
	assert := assert.New(t)

	secret := []byte("test")

	_, err := Split(secret, 0, 0)
	assert.NotNil(err)
	_, err = Split(secret, 2, 3)
	assert.NotNil(err)
	_, err = Split(secret, 2, 3)
	assert.NotNil(err)
	_, err = Split(secret, 1000, 3)
	assert.NotNil(err)
	_, err = Split(secret, 10, 1)
	assert.NotNil(err)
	_, err = Split(nil, 3, 2)
	assert.NotNil(err)
}

func TestSplit(t *testing.T) {
	assert := assert.New(t)

	secret := []byte("test")

	out, err := Split(secret, 5, 3)
	assert.Nil(err)
	assert.Len(out, 5)

	for _, share := range out {
		assert.Equal(len(share), len(secret)+1)
	}
}

func TestCombineInvalid(t *testing.T) {
	assert := assert.New(t)

	_, err := Combine(nil)
	assert.NotNil(err)

	// Mis-match in length
	parts := [][]byte{
		[]byte("foo"),
		[]byte("ba"),
	}
	_, err = Combine(parts)
	assert.NotNil(err)

	//Too short
	parts = [][]byte{
		[]byte("f"),
		[]byte("b"),
	}
	_, err = Combine(parts)
	assert.NotNil(err)

	parts = [][]byte{
		[]byte("foo"),
		[]byte("foo"),
	}
	if _, err := Combine(parts); err == nil {
		t.Fatalf("should err")
	}
}

func TestCombine(t *testing.T) {
	assert := assert.New(t)

	secret := []byte("test")

	out, err := Split(secret, 5, 3)
	assert.Nil(err)

	// There is 5*4*3 possible choices,
	// we will just brute force try them all
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if j == i {
				continue
			}
			for k := 0; k < 5; k++ {
				if k == i || k == j {
					continue
				}
				parts := [][]byte{out[i], out[j], out[k]}
				recomb, err := Combine(parts)
				assert.Nil(err)
				assert.True(bytes.Equal(recomb, secret), fmt.Sprintf("parts: (i:%d, j:%d, k:%d) %v", i, j, k, parts))
			}
		}
	}
}

func TestFieldAdd(t *testing.T) {
	if out := add(16, 16); out != 0 {
		t.Fatalf("Bad: %v 16", out)
	}

	if out := add(3, 4); out != 7 {
		t.Fatalf("Bad: %v 7", out)
	}
}

func TestFieldMult(t *testing.T) {
	if out := mult(3, 7); out != 9 {
		t.Fatalf("Bad: %v 9", out)
	}

	if out := mult(3, 0); out != 0 {
		t.Fatalf("Bad: %v 0", out)
	}

	if out := mult(0, 3); out != 0 {
		t.Fatalf("Bad: %v 0", out)
	}
}

func TestFieldDivide(t *testing.T) {
	if out := div(0, 7); out != 0 {
		t.Fatalf("Bad: %v 0", out)
	}

	if out := div(3, 3); out != 1 {
		t.Fatalf("Bad: %v 1", out)
	}

	if out := div(6, 3); out != 2 {
		t.Fatalf("Bad: %v 2", out)
	}
}

func TestPolynomialRandom(t *testing.T) {
	p, err := makePolynomial(42, 2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if p.coefficients[0] != 42 {
		t.Fatalf("bad: %v", p.coefficients)
	}
}

func TestPolynomialEval(t *testing.T) {
	p, err := makePolynomial(42, 1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if out := p.evaluate(0); out != 42 {
		t.Fatalf("bad: %v", out)
	}

	out := p.evaluate(1)
	exp := add(42, mult(1, p.coefficients[1]))
	if out != exp {
		t.Fatalf("bad: %v %v %v", out, exp, p.coefficients)
	}
}

func TestInterpolateRand(t *testing.T) {
	for i := 0; i < 256; i++ {
		p, err := makePolynomial(uint8(i), 2)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		xVals := []uint8{1, 2, 3}
		yVals := []uint8{p.evaluate(1), p.evaluate(2), p.evaluate(3)}
		out := interpolatePolynomial(xVals, yVals, 0)
		if out != uint8(i) {
			t.Fatalf("Bad: %v %d", out, i)
		}
	}
}
