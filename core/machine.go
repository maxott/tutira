package core

import (
	"context"
	"log"
)

func Eval(prog string) (res TermI, err error) {
	stream := StringStream(prog)
	reader := Reader(stream)
	m := NewMachine(reader)
	c := context.Background()
	return m.Eval(c)
}

type MachineI interface {
	Eval(ctxt context.Context) (res TermI, err error)
	ResetScope() MachineI
}

type machine struct {
	reader ReaderI
	scope  Scope
}

func NewMachine(reader ReaderI) MachineI {
	m := &machine{reader: reader}
	return m.ResetScope()
}

func (m *machine) Eval(ctxt context.Context) (res TermI, err error) {
	term, err := m.reader.Next(ctxt, false)
	if err != nil {
		return
	}
	if list, ok := term.(*ListT); ok {
		res, err = list.Eval(m.scope, NewEvalContext(list))
	} else {
		res, err = term.Eval(m.scope, NewEvalContext(nil))
	}
	return
}

func (m *machine) ResetScope() MachineI {
	rs := RootScope{
		functions: coreFunction,
	}
	m.scope = rs.Inner(NewEvalContext(nil))
	return m
}

type RootScope struct {
	functions map[string]LambdaT
}

func (s *RootScope) GetLambda(sym *Symbol) (fn LambdaT, ok bool) {
	str, _ := sym.AsString()
	fn, ok = s.functions[*str]
	return
}

func (s *RootScope) Inner(ctxt EvalContext) Scope {
	return &BaseScope{1, s, map[string]TermI{}}
}

func (s *RootScope) AddBinding(*Symbol, TermI) {
	log.Fatalf("'AddBinding' should never be called")
}

func (s *RootScope) GetBinding(sym *Symbol) (TermI, bool) {
	return s.GetLambda(sym)
}
