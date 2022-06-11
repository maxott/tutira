package core

import (
	"context"
	_ "fmt"
	"regexp"
	_ "strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextRune(t *testing.T) {
	c := context.Background()
	s := StringStream("AZαω")
	for _, exr := range []rune{'A', 'Z', '\u03B1', '\u03C9'} {
		if r, ok, err := s.Next(c); ok {
			assert.Equal(t, exr, r)
		} else {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestNextMatch(t *testing.T) {
	testNextMatch(t, `^\(`, "() and more", "(", 0)

	stringRE := `^(\s*)"((?:[^"\\]|\\.)*)"`
	testNextMatch(t, stringRE, "\"more\"", "more", 2)
	testNextMatch(t, stringRE, "\"mo\\\"re\"", "mo\\\"re", 2)
	testNextMatch(t, stringRE, " \n \"more\" no", "more", 2)
	testNextMatch(t, stringRE, "\"mo\nre\" no", "mo\nre", 2)

}

func testNextMatch(t *testing.T, rx string, in string, exp string, matchIdx int) {
	c := context.Background()
	match, ok, err := StringStream(in).NextMatch(c, regexp.MustCompile(rx))
	if !ok {
		t.Errorf("expected OK to be true but it was false")
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	assert.Equal(t, exp, string(match[matchIdx]))
}
