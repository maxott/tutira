package core

import (
	"context"
	"fmt"
	"regexp"
	"unicode/utf8"
)

type StreamI interface {
	Next(ctxt context.Context) (r rune, ok bool, err error)
	NextMatch(ctxt context.Context, re *regexp.Regexp) (match [][]byte, ok bool, err error)
	Peek(ctxt context.Context) (r rune, err error)
	PeekMulti(ctxt context.Context, length int) (rs []rune, err error)
}

type stream struct {
	data   []byte
	offset int
	length int
}

func StringStream(s string) StreamI {
	return ByteStream([]byte(s))
}

func ByteStream(data []byte) StreamI {
	return &stream{
		data:   data,
		offset: 0,
		length: len(data),
	}
}

func (s *stream) Next(ctxt context.Context) (r rune, ok bool, err error) {
	r, size := utf8.DecodeRune(s.data[s.offset:])
	if r == utf8.RuneError {
		err = fmt.Errorf("cannot read rune from data")
	}
	s.offset += size
	ok = true
	return
}

func (s *stream) NextMatch(ctxt context.Context, re *regexp.Regexp) (match [][]byte, ok bool, err error) {
	loc := re.FindSubmatchIndex(s.data[s.offset:])
	// fmt.Printf("Loc len %d - loc %v\n", len(loc), loc)
	if len(loc) == 0 {
		return
	}

	ll := len(loc)
	for i := 0; i < ll; i += 2 {
		fi := loc[i]
		if fi < 0 {
			match = append(match, []byte{})
		} else {
			from := fi + s.offset
			to := loc[i+1] + s.offset
			match = append(match, s.data[from:to])
		}
	}
	ok = true
	s.offset += loc[1]
	return
}

func (s *stream) Peek(ctxt context.Context) (r rune, err error) {
	r, _ = utf8.DecodeRune(s.data[s.offset:])
	if r == utf8.RuneError {
		err = fmt.Errorf("cannot read rune from data")
	}
	return
}

func (s *stream) PeekMulti(ctxt context.Context, length int) (ra []rune, err error) {
	mark := s.offset
	ra = make([]rune, length)
	off := mark
	for i := 0; i < length; i++ {
		r, size := utf8.DecodeRune(s.data[off:])
		if r == utf8.RuneError {
			err = fmt.Errorf("cannot read rune from data")
			ra = make([]rune, 0)
			s.offset = mark
			return
		}
		off += size
	}
	s.offset = off
	return
}
