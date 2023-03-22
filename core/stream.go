package core

import (
	"context"
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"
)

type StreamI interface {
	Next(ctxt context.Context) (r rune, ok bool, err *StreamError)
	NextMatch(ctxt context.Context, re *regexp.Regexp) (match [][]byte, ok bool, err *StreamError)
	Peek(ctxt context.Context, consumeWhitespace bool) (r rune, err *StreamError)
	PeekMulti(ctxt context.Context, length int) (rs []rune, err *StreamError)
	IsAtEnd(ctxt context.Context) bool
	Offset() int
	Mark()
	Reset(ctxt context.Context) *StreamError
}

type StreamError struct {
	Msg    string
	Offset int
	IsEoS  bool
}

func (e *StreamError) Error() string {
	return fmt.Sprintf("%s - offset: %d", e.Msg, e.Offset)
}

func streamError(msg string, s *ByteStream) *StreamError {
	return &StreamError{Msg: msg, Offset: s.offset, IsEoS: false}
}

func eosError(s *ByteStream) *StreamError {
	return &StreamError{Msg: "end-of-stream", Offset: s.offset, IsEoS: true}
}

type ByteStream struct {
	data   []byte
	offset int
	length int
	mark   int
}

func StringStream(s string) StreamI {
	return NewByteStream([]byte(s))
}

func NewByteStream(data []byte) *ByteStream {
	return (&ByteStream{}).SetData(data)
}

func (s *ByteStream) Next(ctxt context.Context) (r rune, ok bool, err *StreamError) {
	r, size := utf8.DecodeRune(s.data[s.offset:])
	if r == utf8.RuneError {
		err = s.createRuneError(ctxt)
	}
	s.offset += size
	ok = true
	return
}

func (s *ByteStream) NextMatch(ctxt context.Context, re *regexp.Regexp) (match [][]byte, ok bool, err *StreamError) {
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

func (s *ByteStream) Peek(ctxt context.Context, consumeWhitespace bool) (r rune, err *StreamError) {
	b := s.data[s.offset:]
	off := 0
	var size int
	for {
		r, size = utf8.DecodeRune(b[off:])
		if r == utf8.RuneError {
			err = s.createRuneError(ctxt)
			return
		}
		if !consumeWhitespace || !unicode.IsSpace(r) {
			break
		}
		off += size
	}
	s.offset += off
	return

	// r, _ = utf8.DecodeRune(s.data[s.offset:])
	// if r == utf8.RuneError {
	// 	err = s.createRuneError(ctxt)
	// }
	// return
}

func (s *ByteStream) PeekMulti(ctxt context.Context, length int) (ra []rune, err *StreamError) {
	mark := s.offset
	ra = make([]rune, length)
	off := mark
	for i := 0; i < length; i++ {
		r, size := utf8.DecodeRune(s.data[off:])
		if r == utf8.RuneError {
			err = s.createRuneError(ctxt)
			ra = make([]rune, 0)
			s.offset = mark
			return
		}
		off += size
	}
	s.offset = off
	return
}

func (s *ByteStream) IsAtEnd(ctxt context.Context) bool {
	return s.offset >= s.length
}

func (s *ByteStream) Offset() int {
	return s.offset
}

func (s *ByteStream) Mark() {
	s.mark = s.offset
}

func (s *ByteStream) Reset(ctxt context.Context) (err *StreamError) {
	if s.mark < 0 {
		return streamError("no mark set", s)
	}
	s.offset = s.mark
	s.mark = -1
	return
}

func (s *ByteStream) ToString() string {
	return string(s.data[s.offset:])
}

func (s *ByteStream) SetData(data []byte) *ByteStream {
	s.data = data
	s.offset = 0
	s.length = len(data)
	s.mark = -1
	return s
}

func (s *ByteStream) AppendData(data []byte) *ByteStream {
	s.data = append(s.data, data...)
	s.length = len(data)
	return s
}

func (s *ByteStream) createRuneError(ctxt context.Context) *StreamError {
	if s.IsAtEnd(ctxt) {
		return eosError(s)
	} else {
		return streamError("cannot read rune from data", s)
	}
}
