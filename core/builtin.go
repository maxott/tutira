package core

import (
	"fmt"
	"math"
	"strings"
)

func init() {
	addCoreFunction("quote", quoteFn, "missing docs")
	addCoreFunction("type", typeFn, "missing docs")

	addAliasedCoreFunction([]string{"lambda", "fn"}, lambdaFn, "missing docs")

	addCoreFunction("progn", prognFn, "missing docs")
	addCoreFunction("setq", setqFn, "missing docs")
	addCoreFunction("let", letFn, "missing docs")
	addAliasedCoreFunction([]string{"pipe", "->"}, pipeFn, "missing docs")
	addCoreFunction("import", importFn, "import into the current context all functions from the respective namespaces")

	addAliasedCoreFunction([]string{"eq", "="}, equalFn, "missing docs")
	addAliasedCoreFunction([]string{"le", "<"}, lessFn, "missing docs")
	addAliasedCoreFunction([]string{"leq", "<="}, lessEqualFn, "missing docs")
	addAliasedCoreFunction([]string{"gt", ">"}, greaterFn, "missing docs")
	addAliasedCoreFunction([]string{"geq", ">="}, greaterEqualFn, "missing docs")

	addAliasedCoreFunction([]string{"add", "+"}, addFn, "missing docs")
	addAliasedCoreFunction([]string{"mul", "*"}, mulFn, "missing docs")
}

func quoteFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	els := params.elements
	if len(els) == 1 {
		return els[0], nil
	} else {
		return &ListT{params.elements, -1}, nil
	}
}

func typeFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	els := make([]TermI, len(params.elements))
	for i, t := range params.elements {
		typeS := t.Type()
		if sym, ok := t.(Symbol); ok {
			if b, ok2 := scope.GetBinding(&sym); ok2 {
				typeS = b.Type()
			} else if l, ok3 := scope.GetLambda(&sym); ok3 {
				typeS = l.Type()
			}
		}
		els[i] = String(typeS)
	}
	return &ListT{els, -1}, nil
}

func lambdaFn(outScope Scope, outParams *ListT, octxt EvalContext) (TermI, error) {
	if len(outParams.elements) < 2 {
		return nil, &EvalError{"fn: Need at least two parameters", outParams}
	}
	if vl, ok := outParams.Car().AsList(); !ok {
		return nil, &EvalError{"fn: First parameter needs to be a list of variables", outParams.Car()}
	} else {
		variables := make([]*Symbol, len(vl.elements))
		for i, e := range vl.elements {
			if sym, ok := e.AsSymbol(); ok {
				variables[i] = sym
			} else {
				return nil, octxt.Error("Parameters need to be symbols", e)
			}
		}
		var fn LambdaT
		fn = func(inScope Scope, inParams *ListT, ictxt EvalContext) (TermI, error) {
			if len(inParams.elements) != len(variables) {
				return nil, ictxt.Error("fn: Mismatched numbers of parameters", inParams)
			}
			if params, err := inParams.EvaluatedElements(inScope, ictxt); err != nil {
				return nil, err
			} else {
				scope := outScope.Inner(ictxt)
				for i, v := range params {
					scope.AddBinding(variables[i], v)
				}
				if els, err := outParams.Cdr().EvaluatedElements(scope, octxt); err != nil {
					return nil, err
				} else {
					return els[len(els)-1], nil
				}
			}
		}
		return fn, nil
	}
}

func prognFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	if els, err := params.EvaluatedElements(scope.Inner(ctxt), ctxt); err != nil {
		return nil, err
	} else {
		return els[len(els)-1], nil
	}
}

func setqFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	els := params.Elements()
	if len(els) != 2 {
		// let's keep it simple
		return nil, ctxt.Error("setq only handles two parameters", nil)
	}
	if name, ok := els[0].AsSymbol(); ok {
		val, err := els[1].Eval(scope, ctxt)
		if err != nil {
			return nil, err
		}
		scope.AddBinding(name, val)
		return els[1], nil
	} else {
		return nil, ctxt.Error("setq first parameter is not a symbol", els[0])
	}
}

// (let (params) form form ..)
func letFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	els := params.Elements()
	if len(els) < 2 {
		return nil, ctxt.Error("let needs at least two parameters", nil)
	}
	inner := scope.Inner(ctxt)
	if pl, ok := els[0].AsList(); !ok {
		return nil, ctxt.Error("let first parameter needs to be a list", els[0])
	} else {
		for _, el := range pl.Elements() {
			if pli, ok2 := el.AsList(); !ok2 {
				return nil, ctxt.Error("let nil'ed binding not implemented", el)
			} else {
				iels := pli.Elements()
				if len(iels) != 2 {
					return nil, ctxt.Error("let binding can only contain 2 elements", pli)
				}
				if bName, ok3 := iels[0].AsSymbol(); !ok3 {
					return nil, ctxt.Error("let binding first pair needs to be a symbol", iels[0])
				} else {
					if bValue, err := iels[1].Eval(inner, ctxt); err != nil {
						return nil, err
					} else {
						inner.AddBinding(bName, bValue)
					}
				}
			}
		}
	}
	if body, err := params.Cdr().EvaluatedElements(inner, ctxt); err != nil {
		return nil, err
	} else {
		return body[len(body)-1], nil
	}
}

func importFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	h := ParamsHelper(params, true, 1, -1, scope, ctxt)
	els := make([]TermI, 0)
	for i, ns := range h.RestAsStrings() {
		fna, ok := namespaceFunctions[ns]
		if !ok {
			return nil, ctxt.Error(fmt.Sprintf("unknown namespace '%s'", h.Term(i)), nil)
		}
		for name, fn := range fna {
			sy := Symbol(name)
			els = append(els, sy)
			scope.AddBinding(&sy, fn)
		}
	}
	return &ListT{els, -1}, nil
}

// (-> form1 form2 form3 ...)
func pipeFn(scope Scope, params *ListT, ctxt EvalContext) (res TermI, err error) {
	h := ParamsHelper(params, false, 1, math.MaxInt, scope, ctxt)
	first := h.NextAsTerm()
	if res, err = first.Eval(scope, ctxt); err != nil {
		return
	}
	for h.HasNext() {
		nextForm := h.NextAsList()
		if err = h.Error(); err != nil {
			return nil, err // res may already been set, but don't return it
		}
		if nextForm, err = nextForm.Insert(1, res, ctxt); err != nil {
			return nil, err
		}
		if res, err = nextForm.Eval(scope, ctxt); err != nil {
			return
		}
	}
	err = h.Error()
	return
}

func equalFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	return _cmpHelper(
		func(a *string, b *string) bool { return strings.Compare(*a, *b) == 0 },
		func(a int64, b int64) bool { return a == b },
		func(a float64, b float64) bool { return a == b },
		scope, params, ctxt)
}

func lessFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	return _cmpHelper(
		func(a *string, b *string) bool { return strings.Compare(*a, *b) == -1 },
		func(a int64, b int64) bool { return a < b },
		func(a float64, b float64) bool { return a < b },
		scope, params, ctxt)
}

func lessEqualFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	return _cmpHelper(
		func(a *string, b *string) bool { return strings.Compare(*a, *b) <= 0 },
		func(a int64, b int64) bool { return a <= b },
		func(a float64, b float64) bool { return a <= b },
		scope, params, ctxt)
}

func greaterFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	return _cmpHelper(
		func(a *string, b *string) bool { return strings.Compare(*a, *b) == 1 },
		func(a int64, b int64) bool { return a > b },
		func(a float64, b float64) bool { return a > b },
		scope, params, ctxt)
}

func greaterEqualFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	return _cmpHelper(
		func(a *string, b *string) bool { return strings.Compare(*a, *b) >= 0 },
		func(a int64, b int64) bool { return a >= b },
		func(a float64, b float64) bool { return a >= b },
		scope, params, ctxt)
}

func _cmpHelper(
	stringCmp func(*string, *string) bool,
	intCmp func(int64, int64) bool,
	floatCmp func(float64, float64) bool,
	scope Scope, params *ListT, ctxt EvalContext,
) (TermI, error) {
	if len(params.elements) != 2 {
		return nil, &EvalError{"Wrong number of arguments", nil}
	}
	if els, err := params.EvaluatedElements(scope, ctxt); err != nil {
		return nil, err
	} else {
		ctype := CommonType(els)
		if ctype == "" {
			return nil, &EvalError{"paramters need to be comparable", nil}
		}
		var res bool
		switch ctype {
		case "string":
			{
				p1, _ := els[0].AsString()
				p2, _ := els[1].AsString()
				res = stringCmp(p1, p2)
			}
		case "float64", "int64":
			{
				n1, _ := els[0].AsNumber()
				n2, _ := els[1].AsNumber()
				if ctype == "float64" {
					p1, _ := n1.AsFloat()
					p2, _ := n2.AsFloat()
					res = floatCmp(*p1, *p2)
				} else {
					p1, _ := n1.AsInt()
					p2, _ := n2.AsInt()
					res = intCmp(*p1, *p2)
				}
			}
		default:
			return nil, &EvalError{"unsupported common type", nil}
		}
		return Bool(res), nil
	}
}

func addFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	return reducer(
		func(p int64, v int64) int64 { return p + v },
		func(p float64, v float64) float64 { return p + v },
		scope, params, ctxt,
	)
}

func mulFn(scope Scope, params *ListT, ctxt EvalContext) (TermI, error) {
	return reducer(
		func(p int64, v int64) int64 { return p * v },
		func(p float64, v float64) float64 { return p * v },
		scope, params, ctxt,
	)
}

func reducer(
	mapInt func(int64, int64) int64,
	mapFloat func(float64, float64) float64,
	scope Scope,
	params *ListT,
	ctxt EvalContext,
) (TermI, error) {
	if els, err := params.EvaluatedElements(scope, ctxt); err != nil {
		return nil, err
	} else {
		ctype := CommonType(els)
		if ctype == "" {
			return nil, &EvalError{"paramters need to be all numbers", nil}
		}
		switch ctype {
		case "float64":
			var v float64
			for i, el := range els {
				if vp, ok := ToFloat(el); ok {
					if i == 0 {
						v = *vp
					} else {
						v = mapFloat(v, *vp)
					}
				} else {
					return nil, &EvalError{"Not a number", el}
				}
			}
			return Float64(v), nil
		case "int64":
			var v int64
			for i, el := range els {
				if vp, ok := ToInt(el); ok {
					if i == 0 {
						v = *vp
					} else {
						v = mapInt(v, *vp)
					}
				} else {
					return nil, &EvalError{"Not a number", el}
				}
			}
			return Int64(v), nil
		default:
			return nil, &EvalError{"unsupported common type", nil}
		}
	}
}
