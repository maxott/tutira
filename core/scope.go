package core

type Scope interface {
	Inner(EvalContext) Scope
	AddBinding(*Symbol, TermI)
	GetBinding(*Symbol) (TermI, bool)
	GetLambda(*Symbol) (LambdaT, bool)
}

type BaseScope struct {
	level    int
	upper    Scope
	bindings map[string]TermI
}

func (s *BaseScope) GetLambda(sym *Symbol) (fn LambdaT, ok bool) {
	if t, ok := s.upper.GetBinding(sym); ok {
		if lp, ok2 := t.AsLambda(); ok2 {
			return *lp, true
		} else {
			return nil, false
		}
	} else {
		return s.upper.GetLambda(sym)
	}
}

func (s *BaseScope) Inner(ctxt EvalContext) Scope {
	return &BaseScope{s.level + 1, s, map[string]TermI{}}
}

func (s *BaseScope) AddBinding(sym *Symbol, term TermI) {
	str, _ := sym.AsString()
	s.bindings[*str] = term
}

func (s *BaseScope) GetBinding(sym *Symbol) (term TermI, ok bool) {
	str, _ := sym.AsString()
	if term, ok = s.bindings[*str]; !ok {
		return s.upper.GetBinding(sym)
	}
	return
}
