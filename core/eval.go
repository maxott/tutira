package core

import (
	"fmt"
	"strings"
)

type TermI interface {
	Eval(scope Scope, ctxt EvalContext) (TermI, error)
	AsList() (*ListT, bool)
	AsSymbol() (*Symbol, bool)
	AsNumber() (NumberI, bool)
	AsString() (*string, bool)
	AsBool() (bool, bool)
	AsLambda() (*LambdaT, bool)
	IsNil() bool
	Type() string
	Print(bool) string
}

type NumberI interface {
	AsInt() (*int64, bool)
	AsFloat() (*float64, bool)
}

type EvalContext interface {
	Error(string, TermI) error
	WrongParameterCount(int, int) error
	WrongParameterType(string, TermI) error
}

type LambdaT func(scope Scope, params *ListT, ctxt EvalContext) (TermI, error)

// Return the type all elements can be cast to, returns the empty
// string if no common element can be found
func CommonType(els []TermI) string {
	if len(els) == 0 {
		return ""
	}
	common := els[0].Type()
	for _, el := range els[1:] {
		t := el.Type()
		if t != common {
			var key string
			if strings.Compare(common, t) < 0 {
				key = fmt.Sprintf("%s:%s", common, t)
			} else {
				key = fmt.Sprintf("%s:%s", t, common)
			}
			common = commonType[key]
			if common == "" {
				break
			}
		}
	}
	return common
}

var commonType = map[string]string{
	"float64:int64":  "float64",
	"string:symbol":  "string",
	"keyword:string": "string",
	"keyword:symbol": "string",
}

//****** IMPLEMENTATION ******

type evalContext struct {
	list *ListT
}

func NewEvalContext(list *ListT) *evalContext {
	return &evalContext{list}
}

func (c *evalContext) Error(msg string, cause TermI) error {
	return &EvalError{msg, cause}
}

func (c *evalContext) WrongParameterCount(minCount int, maxCount int) error {
	var msg string
	cnt := len(c.list.elements) - 1
	if minCount <= 0 {
		msg = fmt.Sprintf("expected a max of '%d' parameter(s), but got '%d'", maxCount, cnt)
	} else {
		if maxCount >= 0 {
			msg = fmt.Sprintf("expected between '%d' and '%d' parameters, but got '%d'", minCount, maxCount, cnt)
		} else {
			msg = fmt.Sprintf("expected a minimum of '%d' parameter(s), but got '%d'", minCount, cnt)
		}
	}
	return &EvalError{msg, nil}
}

func (c *evalContext) WrongParameterType(expType string, term TermI) error {
	msg := fmt.Sprintf("expected '%s' type, but got '%s'", expType, term.Type())
	return &EvalError{msg, term}
}

type EvalError struct {
	msg  string
	term TermI
}

func (e EvalError) Error() string {
	if e.term != nil {
		return fmt.Sprintf("%s - %s", e.msg, e.term.Print(true))
	} else {
		return e.msg
	}
}

type FunctionNotFoundError struct {
	name *Symbol
}

func (e FunctionNotFoundError) Error() string {
	name, _ := e.name.AsString()
	return fmt.Sprintf("Function '%s' not found", *name)
}
