package relalg

import (
	"fmt"
	_ "regexp"
	_ "strconv"
	"testing"

	"github.com/maxott/tutiro/core"
	"github.com/stretchr/testify/assert"
)

func TestRelation(t *testing.T) {
	prog := `(relalg/relation)`
	assert.Equal(t, "(relalg/request \"foo\")", EvalOK(prog, t).Print(false))

	prog = `(progn (import "relalg") (relation "foo"))`
	assert.Equal(t, "(relalg/request \"foo\")", EvalOK(prog, t).Print(false))

	prog = `(progn (progn (import "relalg") (relation "foo")) (relation "foo"))`
	assert.Equal(t, "Can't resolve binding - relation", EvalError(prog, t).Error())
}

//*** HELPERS

func EvalOK(prog string, t *testing.T) core.TermI {
	if res, err := core.Eval(prog); err == nil {
		fmt.Printf("Prog: '%s' res: '%v':%T\n", prog, res, res)
		return res
	} else {
		t.Fatalf("Executing '%s' err: %v", prog, err)
		return nil
	}
}

func EvalError(prog string, t *testing.T) error {
	if res, err := core.Eval(prog); err == nil {
		t.Fatalf("Prog: '%s' expected to fail, but returned: '%v':%T\n", prog, res, res)
		return nil
	} else {
		return err
	}
}
