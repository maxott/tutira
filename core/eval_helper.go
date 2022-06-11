package core

import (
	"fmt"
)

type paramsHelper struct {
	terms []TermI
	// evalFirst bool
	// scope     Scope
	ctxt  EvalContext
	index int
	err   error
}

func ParamsHelper(
	params *ListT,
	evalFirst bool,
	minNumber int,
	maxNumber int, // -1 ... infinity
	scope Scope,
	ctxt EvalContext,
) *paramsHelper {
	pcnt := len(params.Elements())
	if pcnt < minNumber {
		return &paramsHelper{err: ctxt.WrongParameterCount(minNumber, maxNumber)}
	}
	if maxNumber >= 0 && pcnt > maxNumber {
		return &paramsHelper{err: ctxt.WrongParameterCount(minNumber, maxNumber)}
	}
	var terms []TermI
	if evalFirst {
		var err error
		if terms, err = params.EvaluatedElements(scope, ctxt); err != nil {
			return &paramsHelper{err: err}
		}
	} else {
		terms = params.elements
	}

	return &paramsHelper{terms: terms, ctxt: ctxt}
}

func (p *paramsHelper) NextAsList() (v *ListT) {
	if p.ensureNext() != nil {
		return
	}

	term := p.terms[p.index]
	var ok bool
	if v, ok = term.AsList(); ok {
		p.index += 1
	} else {
		p.err = p.ctxt.WrongParameterType("list", term)
	}
	return
}

func (p *paramsHelper) NextAsTerm() (v TermI) {
	if p.ensureNext() != nil {
		return
	}

	v = p.terms[p.index]
	p.index += 1
	return
}

func (p *paramsHelper) Term(index int) (v TermI) {
	if index < 0 || index >= len(p.terms) {
		p.err = p.ctxt.Error(fmt.Sprintf("wrong index '%d'", index), nil)
	}
	return p.terms[index]
}

func (p *paramsHelper) NextAsString() (v string) {
	if p.ensureNext() != nil {
		return
	}

	term := p.terms[p.index]
	if sp, ok := term.AsString(); ok {
		p.index += 1
		v = *sp
	} else {
		p.err = p.ctxt.WrongParameterType("string", term)
	}
	return
}

func (p *paramsHelper) RestAsStrings() []string {
	v := make([]string, 0)
	for p.HasNext() {
		v = append(v, p.NextAsString())
	}
	if p.err == nil {
		return v
	} else {
		return []string{}
	}
}

func (p *paramsHelper) NextAsInt() (v int64) {
	if p.ensureNext() != nil {
		return
	}

	term := p.terms[p.index]
	if num, ok := term.AsNumber(); ok {
		if vp, ok2 := num.AsInt(); ok2 {
			p.index += 1
			v = *vp
			return
		}
	}
	p.err = p.ctxt.WrongParameterType("int64", term)
	return
}

func (p *paramsHelper) NextAsFloat() (v float64) {
	if p.ensureNext() != nil {
		return
	}

	term := p.terms[p.index]
	if num, ok := term.AsNumber(); ok {
		if vp, ok2 := num.AsFloat(); ok2 {
			p.index += 1
			v = *vp
			return
		}
	}
	p.err = p.ctxt.WrongParameterType("float64", term)
	return
}

func (p *paramsHelper) Error() error {
	return p.err
}

func (p *paramsHelper) HasNext() bool {
	return p.err == nil && p.index < len(p.terms)
}

func (p *paramsHelper) ensureNext() error {
	if p.err == nil && !p.HasNext() {
		p.err = p.ctxt.Error("Requesting too many parameters", nil)
	}
	return p.err
}
