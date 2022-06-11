package core

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
)

type ReaderI interface {
	// Return next token. If 'onlySyntax' is set, skip
	// whitespace related tokens
	Next(ctxt context.Context, onlySyntax bool) (term TermI, err error)
}

// Converts the top of the stream to a term. Returns nil if the content doesn't
// match.
type MapperF func(context.Context, StreamI, ReaderI) (TermI, error)
type reader struct {
	stream        StreamI
	mapper        map[rune]MapperF
	defaultMapper MapperF
}

func Reader(stream StreamI) ReaderI {
	mapper := map[rune]MapperF{
		'(':  listMapper,
		'[':  vectorMapper,
		'"':  stringMapper,
		'\'': quoteMapper,
	}
	for _, r := range "+-0123456789" {
		mapper[r] = numberMapper
	}

	return &reader{
		stream:        stream,
		mapper:        mapper,
		defaultMapper: symbolMapper,
	}
}

const TAB = rune(0x9)
const LF = rune(0xa)
const CR = rune(0xd)
const SPACE = rune(0x20)

func (r *reader) Next(ctxt context.Context, onlySyntax bool) (term TermI, err error) {
	var p rune
	for {
		p, err = r.stream.Peek(ctxt)
		if err != nil {
			return
		}
		if p != SPACE && p != TAB && p != LF && p != CR {
			break
		} else {
			// consume whitespace
			r.stream.Next(ctxt)
		}
	}
	if m, ok := r.mapper[p]; ok {
		term, err = m(ctxt, r.stream, r)
		if term != nil || err != nil {
			return
		}
	}
	term, err = r.defaultMapper(ctxt, r.stream, r)
	if term == nil && err == nil {
		err = fmt.Errorf("cannot parse term starting with '%v'", p)
	}
	return
}

var (
	stringRE, numberRE, symbolRE *regexp.Regexp
)

func init() {
	stringRE = regexp.MustCompile(`^"((?:[^"\\]|\\.)*)"`)
	numberRE = regexp.MustCompile(`^([+-]?[0-9]+)(\.[0-9]*)?(E[+-]?[0-9]+)?`)
	symbolRE = regexp.MustCompile(`^([^[:space:])\]]+)`)
}

func listMapper(ctxt context.Context, s StreamI, reader ReaderI) (list TermI, err error) {
	terms, err := getList(ctxt, ')', s, reader)
	if err != nil {
		return
	}
	list = List(terms)
	return
}

func vectorMapper(ctxt context.Context, s StreamI, reader ReaderI) (vector TermI, err error) {
	terms, err := getList(ctxt, ']', s, reader)
	if err != nil {
		return
	}
	vector = &VectorT{elements: terms}
	return
}

func getList(ctxt context.Context, closingRune rune, s StreamI, reader ReaderI) (terms []TermI, err error) {
	s.Next(ctxt) // consume opening rune
	var p rune
	for {
		if p, err = s.Peek(ctxt); err != nil {
			return
		}
		if p == closingRune {
			s.Next(ctxt) // consume closing rune
			break
		}
		var t TermI
		if t, err = reader.Next(ctxt, false); err != nil {
			return
		}
		terms = append(terms, t)
	}
	return
}

func stringMapper(ctxt context.Context, s StreamI, _ ReaderI) (term TermI, err error) {
	match, ok, err := s.NextMatch(ctxt, stringRE)
	if ok {
		term = String(string(match[1]))
		return
	}
	if err == nil {
		// a non match is always an error as no symbol can start with '"'
		err = fmt.Errorf("what looked like a string, isn't one")
		return
	}
	return
}

func numberMapper(ctxt context.Context, s StreamI, _ ReaderI) (term TermI, err error) {
	match, ok, err := s.NextMatch(ctxt, numberRE)
	if ok {
		if len(match) == 4 {
			if len(match[2])+len(match[3]) > 0 {
				var n float64
				if n, err = strconv.ParseFloat(string(match[0]), 64); err == nil {
					term = Float64(n)
				}
			} else {
				var n int64
				if n, err = strconv.ParseInt(string(match[0]), 10, 64); err == nil {
					term = Int64(n)
				}
			}
			return
		}
	}
	// ok to return nil as an '+' could be the beginning of a number of a standalone symbol
	return
}

func quoteMapper(ctxt context.Context, s StreamI, reader ReaderI) (list TermI, err error) {
	s.Next(ctxt) // consume quote symbol
	var t TermI
	if t, err = reader.Next(ctxt, false); err != nil {
		return
	}
	terms := []TermI{Symbol("quote"), t}
	list = List(terms)
	return
}

func symbolMapper(ctxt context.Context, s StreamI, _ ReaderI) (term TermI, err error) {
	match, ok, err := s.NextMatch(ctxt, symbolRE)
	if ok {
		term = Symbol(string(match[1]))
		return
	}
	if err == nil {
		// symbols are the most generic atoms, so we should actually never get here
		// except maybe for an empty stream
		err = fmt.Errorf("what looked like a symbol, isn't one")
		return
	}
	return
}
