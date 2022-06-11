package core

import (
	_ "fmt"
)

const (
	ListType    string = "list"
	VectorType  string = "vector"
	SymbolType  string = "symbol"
	KeywordType string = "keyword"
	StringType  string = "string"
	IntType     string = "int64"
	FloatType   string = "float64"
	BoolType    string = "bool"
	NilType     string = "nil"
	LambdaType  string = "lambda"
)

type VectorT struct {
	elements []TermI
}

type Int64 int64
type Float64 float64
type String string

func ToInt(term TermI) (*int64, bool) {
	if n, ok := term.AsNumber(); ok {
		if v, ok2 := n.AsInt(); ok2 {
			return v, true
		}
	}
	return nil, false
}

func ToFloat(term TermI) (*float64, bool) {
	if n, ok := term.AsNumber(); ok {
		if v, ok2 := n.AsFloat(); ok2 {
			return v, true
		}
	}
	return nil, false
}

type Symbol string
type Keyword string
type Nil struct{} // not sure if there is a better way to define a Nil type
type Bool bool

func (s Symbol) Eval(scope Scope, ctxt EvalContext) (term TermI, err error) {
	if v, ok := scope.GetBinding(&s); !ok {
		err = &EvalError{"Can't resolve binding", s}
		// term = s
		return
	} else {
		term = v
		return
	}
}

func (v *VectorT) Print(withParsingInfo bool) string {
	s := "["
	for i, el := range v.elements {
		if i > 0 {
			s += " "
		}
		s += el.Print(false)
	}
	return s + "]"
}
