package core

import (
	_ "fmt"
)

func init() {
	addCoreFunction("list", listFn, "missing docs")
	addCoreFunction("car", carFn, "missing docs")
	addCoreFunction("cdr", cdrFn, "missing docs")
	addAliasedCoreFunction([]string{"nth", "lindex"}, nthFn, "missing docs")
	addCoreFunction("linsert", linsertFn, "missing docs")
}

func listFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	if elt, err := params.EvaluatedElements(scope, ctxt); err == nil {
		return &ListT{elt, -1}, nil
	} else {
		return nil, err
	}
}

func carFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	fn := func(l *ListT) TermI {
		return l.Car()
	}
	return carCdrHelper(fn, scope, params, ctxt)
}

func cdrFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	fn := func(l *ListT) TermI {
		return l.Cdr()
	}
	return carCdrHelper(fn, scope, params, ctxt)
}

// (nth index list)
func nthFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	h := ParamsHelper(params, true, 2, 2, scope, ctxt)
	index := h.NextAsInt()
	list := h.NextAsList()
	if h.Error() != nil {
		return nil, h.Error()
	}
	le := list.elements
	if int(index) >= len(le) {
		return nil, ctxt.Error("Index too large for list", list)
	}
	return le[index], nil
}

// (linsert index term list)
func linsertFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	h := ParamsHelper(params, true, 3, 3, scope, ctxt)
	index := int(h.NextAsInt())
	term := h.NextAsTerm()
	list := h.NextAsList()
	if h.Error() != nil {
		return nil, h.Error()
	}
	return list.Insert(index, term, ctxt)
}

func carCdrHelper(mapper func(*ListT) TermI, scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	if len(params.elements) != 1 {
		return nil, ctxt.Error("Wrong number of arguments", nil)
	}
	if et, err := params.elements[0].Eval(scope, ctxt); err == nil {
		if l, ok := et.AsList(); ok {
			return mapper(l), nil
		} else {
			return nil, ctxt.Error("Wrong argument type", params.elements[0])
		}
	} else {
		return nil, err
	}
}

type ListT struct {
	elements []TermI
	level    int
}

func List(elements []TermI) (l *ListT) {
	return &ListT{elements, -1}
}

func (l *ListT) Eval(scope Scope, ctxt EvalContext) (term TermI, err error) {
	car := l.Car()
	var fnt TermI
	if fnt, err = car.Eval(scope, ctxt); err != nil {
		return
	}
	if lambda, ok := fnt.AsLambda(); ok {
		return (*lambda)(scope, l.Cdr(), &evalContext{l})
	} else {
		err = EvalError{"Car is not, or not bound to a lambda", car}
		return
	}

	// if fns, ok := fnt.AsSymbol(); !ok {
	// 	// maybe already lambda
	// 	if lambda, ok2 := fnt.AsLambda(); !ok2 {
	// 		err = EvalError{"Car is not a symbol or a raw lambda", car}
	// 		return
	// 	} else {
	// 		return (*lambda)(scope, l.Cdr(), &evalContext{})
	// 	}
	// } else {
	// 	if lambda, ok2 := scope.GetLambda(fns); !ok2 {
	// 		err = EvalError{"Not a lambda", fns}
	// 		return
	// 	} else {
	// 		return lambda(scope, l.Cdr(), &evalContext{})
	// 	}
	// }
}

func (l *ListT) Car() TermI {
	return l.elements[0]
}

func (l *ListT) Cdr() *ListT {
	return &ListT{
		elements: l.elements[1:],
		level:    l.level,
	}
}

func (l *ListT) Elements() []TermI {
	return l.elements
}

func (l *ListT) EvaluatedElements(scope Scope, ctxt EvalContext) ([]TermI, error) {
	els := l.elements
	res := make([]TermI, len(els))
	for i, el := range els {
		if elv, err := el.Eval(scope, ctxt); err == nil {
			res[i] = elv
		} else {
			return nil, err
		}
	}
	return res, nil
}

func (l *ListT) EvaluatedElement(scope Scope, index int, ctxt EvalContext) (TermI, error) {
	if index >= len(l.elements) {
		return nil, ctxt.Error("Index too large for list", l)
	}
	term := l.elements[index]
	return term.Eval(scope, ctxt)
}

// Return new list identical to this, but with additional 'term' inserted at 'index'
func (l *ListT) Insert(index int, term TermI, ctxt EvalContext) (*ListT, error) {
	le := l.elements
	if index > len(le) {
		return nil, ctxt.Error("Index too large for list", l)
	}
	resLen := len(le) + 1
	res := make([]TermI, resLen)
	for i := 0; i < index; i++ {
		res[i] = le[i]
	}
	res[index] = term
	for i := index + 1; i < resLen; i++ {
		res[i] = le[i-1]
	}
	return List(res), nil
}

func (l *ListT) Print(withParsingInfo bool) string {
	s := "("
	for i, el := range l.elements {
		if i > 0 {
			s += " "
		}
		s += el.Print(false)
	}
	return s + ")"
}
