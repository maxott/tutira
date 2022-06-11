package core

import "fmt"

type GenericTerm struct {
	// required
	TermType string
	PrintFn  func(bool) string

	// optional
	EvalF      func(Scope, EvalContext) (TermI, error)
	AsListFn   func() (*ListT, bool)
	AsSymbolFn func() (*Symbol, bool)
	AsNumberFn func() (NumberI, bool)
	AsStringFn func() (*string, bool)
	AsBoolFn   func() (bool, bool)
	AsLambdaFn func() (*LambdaT, bool)
	IsNilFn    func() bool
}

func (g *GenericTerm) Eval(scope Scope, ctxt EvalContext) (term TermI, err error) {
	if g.EvalF != nil {
		return g.EvalF(scope, ctxt)
	} else {
		return g, nil
	}
}

func (g *GenericTerm) AsList() (*ListT, bool) {
	if g.AsListFn != nil {
		return g.AsListFn()
	} else {
		return nil, false
	}
}

func (g *GenericTerm) AsSymbol() (*Symbol, bool) {
	if g.AsSymbolFn != nil {
		return g.AsSymbolFn()
	} else {
		return nil, false
	}
}

func (g *GenericTerm) AsNumber() (NumberI, bool) {
	if g.AsNumberFn != nil {
		return g.AsNumberFn()
	} else {
		return nil, false
	}
}

func (g *GenericTerm) AsString() (*string, bool) {
	if g.AsStringFn != nil {
		return g.AsStringFn()
	} else {
		return nil, false
	}
}

func (g *GenericTerm) AsBool() (bool, bool) {
	if g.AsBoolFn != nil {
		return g.AsBoolFn()
	} else {
		return false, false
	}
}

func (g *GenericTerm) AsLambda() (*LambdaT, bool) {
	if g.AsLambdaFn != nil {
		return g.AsLambdaFn()
	} else {
		return nil, false
	}
}

func (g *GenericTerm) IsNil() bool {
	if g.IsNilFn != nil {
		return g.IsNilFn()
	} else {
		return false
	}
}

func (g *GenericTerm) Type() string {
	return g.TermType
}
func (g *GenericTerm) Print(withParsingInfo bool) string {
	if g.PrintFn != nil {
		return g.PrintFn(withParsingInfo)
	} else {
		return fmt.Sprintf("<<Error: No method provided to print type '%s'>>", g.TermType)
	}
}
