package core

import (
	"context"
	"fmt"
	_ "regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderString(t *testing.T) {
	ts := []string{`abc`, `ab\"cd`, `Ɐ `, "λ"}
	for _, s := range ts {
		term := parseHelper(t, "\""+s+"\"", StringType)
		if term != nil {
			ts, _ := term.AsString()
			assert.Equal(t, s, *ts)
		}
	}
}

func TestReaderNumber(t *testing.T) {
	for _, n := range []string{"123", "  +1234", "-345"} {
		term := parseHelper(t, n, IntType)
		if term != nil {
			n, _ := strconv.ParseInt(strings.TrimSpace(n), 10, 64)
			ti, _ := term.AsNumber()
			tn, _ := ti.AsInt()
			assert.Equal(t, n, *tn)
		}
	}

	for _, n := range []string{"123.", "-1.234", "-1.234E-9", "-1.E-9", "-345E6"} {
		term := parseHelper(t, n, FloatType)
		if term != nil {
			n, _ := strconv.ParseFloat(n, 64)
			ti, _ := term.AsNumber()
			tn, _ := ti.AsFloat()
			assert.Equal(t, n, *tn)
		}
	}
}

func TestSymbolReader(t *testing.T) {
	for _, s := range []string{"  foo  ", string('\u03BB'), "+ ", "/"} {
		term := parseHelper(t, s, SymbolType)
		if term != nil {
			si, _ := term.AsString()
			assert.Equal(t, strings.TrimSpace(s), *si)
		}
	}
}

func TestReaderList(t *testing.T) {
	if l := parseListHelper(t, "(123)"); l != nil {
		el := l.elements
		assert.Equal(t, 1, len(el))
		assert.Equal(t, IntType, el[0].Type())
	}
	if l := parseListHelper(t, "(\"foo\")"); l != nil {
		el := l.elements
		assert.Equal(t, 1, len(el))
		assert.Equal(t, StringType, el[0].Type())
	}
	if l := parseListHelper(t, "(x 12)"); l != nil {
		el := l.elements
		assert.Equal(t, 2, len(el))
		assert.Equal(t, SymbolType, el[0].Type())
		assert.Equal(t, IntType, el[1].Type())
	}
}

func TestReaderListInList(t *testing.T) {
	if l := parseListHelper(t, "((fn [a] 'a) 4)"); l != nil {
		el := l.elements
		assert.Equal(t, 2, len(el))
		assert.Equal(t, ListType, el[0].Type())
		assert.Equal(t, IntType, el[1].Type())
	}
}

func TestReaderPartialList(t *testing.T) {
	s := "((fn [a]"
	c := context.Background()
	stream := StringStream(s)
	reader := Reader(stream)
	_, err := reader.Next(c, true)
	fmt.Printf("err %v", err)
	var uberr *UnbalancedReaderError
	assert.IsType(t, uberr, err)

}

func TestReaderPartialString(t *testing.T) {
	s := "\"foo"
	c := context.Background()
	stream := StringStream(s)
	reader := Reader(stream)
	_, err := reader.Next(c, true)
	var oserr *OpenStringReaderError
	assert.IsType(t, oserr, err)
	//fmt.Printf("Error: %v", err)

}

func TestReaderQuote(t *testing.T) {
	if l := parseListHelper(t, "'123"); l != nil {
		el := l.elements
		assert.Equal(t, 2, len(el))
		assert.Equal(t, SymbolType, el[0].Type())
		qs, _ := el[0].AsString()
		assert.Equal(t, "quote", *qs)
		assert.Equal(t, IntType, el[1].Type())
	}
	if l := parseListHelper(t, "'(123)"); l != nil {
		el := l.elements
		assert.Equal(t, 2, len(el))
		assert.Equal(t, SymbolType, el[0].Type())
		qs, _ := el[0].AsString()
		assert.Equal(t, "quote", *qs)
		assert.Equal(t, ListType, el[1].Type())
	}
}

func parseHelper(t *testing.T, s string, expType string) (term TermI) {
	c := context.Background()
	stream := StringStream(s)
	reader := Reader(stream)
	var err error
	term, err = reader.Next(c, true)
	assert.NoError(t, err, "reader next")
	assert.NotNil(t, term)
	if term != nil {
		assert.Equal(t, expType, term.Type())
	}
	return
}

func parseListHelper(t *testing.T, s string) (list *ListT) {
	term := parseHelper(t, s, ListType)
	if term != nil {
		list, _ = term.AsList()
	}
	return
}
