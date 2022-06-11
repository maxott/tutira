package core

import (
	_ "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLambda(t *testing.T) {
	prog := `((fn (a) a) 4)`
	assert.Equal(t, Int64(4), EvalOK(prog, t))

	prog = `(let ((f (fn (a) (+ a b))) (b 3)) (f 4))`
	assert.Equal(t, Int64(7), EvalOK(prog, t))

	prog = `(let ((f (fn (a) (+ a b))) (b 3)) (let ((b 8)) (f 4)))`
	assert.Equal(t, Int64(7), EvalOK(prog, t))

}

func TestAdd(t *testing.T) {
	prog := `(+ 2 3 4)`
	assert.Equal(t, Int64(9), EvalOK(prog, t))

	prog = `(+ 2 (+ 3 4))`
	assert.Equal(t, Int64(9), EvalOK(prog, t))
}

func TestAddWithError(t *testing.T) {
	prog := `(+ 2 "foo")`
	assert.EqualError(t, evalError(prog, t), "paramters need to be all numbers")
}

func TestEqual(t *testing.T) {
	prog := `(= 2 2)`
	assert.Equal(t, Bool(true), EvalOK(prog, t))

	prog = `(= 2 3)`
	assert.Equal(t, Bool(false), EvalOK(prog, t))

	prog = `(eq 1.0 1.0)`
	assert.Equal(t, Bool(true), EvalOK(prog, t))

	prog = `(eq 1 1.0)`
	assert.Equal(t, Bool(true), EvalOK(prog, t))

	prog = `(eq 1 1)`
	assert.Equal(t, Bool(true), EvalOK(prog, t))

	prog = `(eq "a" "a")`
	assert.Equal(t, Bool(true), EvalOK(prog, t))

	prog = `(eq "a" "ab")`
	assert.Equal(t, Bool(false), EvalOK(prog, t))

}
