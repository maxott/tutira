package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProgn1(t *testing.T) {
	prog := `(progn (+ 2 3))`
	assert.Equal(t, Int64(5), EvalOK(prog, t))
}

func TestProgn2(t *testing.T) {
	prog := `(progn (+ 2 3) (+ 3 4))`
	assert.Equal(t, Int64(7), EvalOK(prog, t))
}

func TestSetq(t *testing.T) {
	prog := `(progn (setq a 3) (+ 2 a))`
	assert.Equal(t, Int64(5), EvalOK(prog, t))
}

func TestLet(t *testing.T) {
	prog := `(let ((a 3)) (+ 4 5) (+ 2 a))`
	assert.Equal(t, Int64(5), EvalOK(prog, t))

	prog = `(let ((a 3) (b (+ a 2))) (+ a b))`
	assert.Equal(t, Int64(8), EvalOK(prog, t))

	prog = `(let ((a 3)) (let ((b (+ a 2))) (+ a b)))`
	assert.Equal(t, Int64(8), EvalOK(prog, t))

	prog = `(let ((a 3)) (let ((a 6)) (+ a 1)) (let ((b (+ a 2))) (+ a b)))`
	assert.Equal(t, Int64(8), EvalOK(prog, t))

}

//*** HELPERS

func EvalOK(prog string, t *testing.T) TermI {
	if res, err := eval(prog, t); err == nil {
		fmt.Printf("Prog: '%s' res: '%v':%T\n", prog, res, res)
		return res
	} else {
		t.Fatalf("Executing '%s' err: %v", prog, err)
		return nil
	}
}

func evalError(prog string, t *testing.T) error {
	if res, err := eval(prog, t); err == nil {
		t.Fatalf("Prog: '%s' expected to fail, but returned: '%v':%T\n", prog, res, res)
		return nil
	} else {
		return err
	}
}

func eval(prog string, t *testing.T) (TermI, error) {
	return Eval(prog)
	// list := parse(prog, t)
	// scope := NewScope()
	// return list.Eval(scope)
}
