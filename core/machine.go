package core

import (
	"context"
	"log"
)

func Eval(prog string) (res TermI, err error) {
	m := NewMachine()
	c := context.Background()
	return m.Eval(c, prog)
}

type machine struct {
	scope Scope
}

func NewMachine() *machine {
	m := &machine{}
	return m.ResetScope()
}

func (m *machine) Eval(ctxt context.Context, prog string) (res TermI, err error) {
	c := context.Background()
	stream := StringStream(prog)
	reader := Reader(stream)
	term, err := reader.Next(c, false)
	if err != nil {
		return
	}
	list, ok := term.(*ListT)
	if !ok {
		err = &EvalError{"Not a list", list}
	}
	return list.Eval(m.scope, NewEvalContext(list))
}

func (m *machine) ResetScope() *machine {
	m.scope = &RootScope{
		functions: coreFunction,
	}
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
