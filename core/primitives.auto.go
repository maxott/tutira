// Automatically generated code, please do NOT edit,
// but edit the respective template "primitive.go.tmpl"
//

package core

import (
	"fmt"
)

//*************** *ListT *****************

func (s *ListT) AsList() (*ListT, bool) {
	return s, true
}

func (s *ListT) AsSymbol() (*Symbol, bool) {
	return nil, false
}

func (s *ListT) AsNumber() (NumberI, bool) {
	return nil, false
}

func (s *ListT) AsString() (*string, bool) {
	return nil, false
}

func (s *ListT) AsBool() (bool, bool) {
	return false, false
}

func (s *ListT) AsLambda() (*LambdaT, bool) {
	return nil, false
}

func (s *ListT) IsNil() bool {
	return false
}

func (s *ListT) Type() string {
	return ListType
}

//*************** *VectorT *****************
func (s *VectorT) Eval(scope Scope, ctxt EvalContext) (term TermI, err error) {
	return s, nil
}

func (s *VectorT) AsList() (*ListT, bool) {
	return nil, false
}

func (s *VectorT) AsSymbol() (*Symbol, bool) {
	return nil, false
}

func (s *VectorT) AsNumber() (NumberI, bool) {
	return nil, false
}

func (s *VectorT) AsString() (*string, bool) {
	return nil, false
}

func (s *VectorT) AsBool() (bool, bool) {
	return false, false
}

func (s *VectorT) AsLambda() (*LambdaT, bool) {
	return nil, false
}

func (s *VectorT) IsNil() bool {
	return false
}

func (s *VectorT) Type() string {
	return VectorType
}

//*************** Symbol *****************

func (s Symbol) AsList() (*ListT, bool) {
	return nil, false
}

func (s Symbol) AsSymbol() (*Symbol, bool) {
	return &s, true
}

func (s Symbol) AsNumber() (NumberI, bool) {
	return nil, false
}

func (s Symbol) AsString() (*string, bool) {
	v := string(s)
	return &v, true
}

func (s Symbol) AsBool() (bool, bool) {
	return false, false
}

func (s Symbol) AsLambda() (*LambdaT, bool) {
	return nil, false
}

func (s Symbol) IsNil() bool {
	return false
}

func (s Symbol) Type() string {
	return SymbolType
}
func (s Symbol) Print(withParsingInfo bool) string {
	if ps, ok := s.AsString(); ok {
		return *ps
	}
	return fmt.Sprintf("%v", s)
}

//*************** Keyword *****************
func (s Keyword) Eval(scope Scope, ctxt EvalContext) (term TermI, err error) {
	return s, nil
}

func (s Keyword) AsList() (*ListT, bool) {
	return nil, false
}

func (s Keyword) AsSymbol() (*Symbol, bool) {
	return nil, false
}

func (s Keyword) AsNumber() (NumberI, bool) {
	return nil, false
}

func (s Keyword) AsString() (*string, bool) {
	v := string(s)
	return &v, true
}

func (s Keyword) AsBool() (bool, bool) {
	return false, false
}

func (s Keyword) AsLambda() (*LambdaT, bool) {
	return nil, false
}

func (s Keyword) IsNil() bool {
	return false
}

func (s Keyword) Type() string {
	return KeywordType
}
func (s Keyword) Print(withParsingInfo bool) string {
	if ps, ok := s.AsString(); ok {
		return fmt.Sprintf(":%s", *ps)
	}
	return fmt.Sprintf("%v", s)
}

//*************** String *****************
func (s String) Eval(scope Scope, ctxt EvalContext) (term TermI, err error) {
	return s, nil
}

func (s String) AsList() (*ListT, bool) {
	return nil, false
}

func (s String) AsSymbol() (*Symbol, bool) {
	return nil, false
}

func (s String) AsNumber() (NumberI, bool) {
	return nil, false
}

func (s String) AsString() (*string, bool) {
	v := string(s)
	return &v, true
}

func (s String) AsBool() (bool, bool) {
	return false, false
}

func (s String) AsLambda() (*LambdaT, bool) {
	return nil, false
}

func (s String) IsNil() bool {
	return false
}

func (s String) Type() string {
	return StringType
}
func (s String) Print(withParsingInfo bool) string {
	if ps, ok := s.AsString(); ok {
		return fmt.Sprintf("\"%s\"", *ps)
	}
	return fmt.Sprintf("%v", s)
}

//*************** Int64 *****************
func (s Int64) Eval(scope Scope, ctxt EvalContext) (term TermI, err error) {
	return s, nil
}

func (s Int64) AsList() (*ListT, bool) {
	return nil, false
}

func (s Int64) AsSymbol() (*Symbol, bool) {
	return nil, false
}

func (s Int64) AsNumber() (NumberI, bool) {
	return s, true
}

func (s Int64) AsString() (*string, bool) {
	return nil, false
}

func (s Int64) AsBool() (bool, bool) {
	return false, false
}

func (s Int64) AsLambda() (*LambdaT, bool) {
	return nil, false
}

func (s Int64) IsNil() bool {
	return false
}

func (s Int64) Type() string {
	return IntType
}
func (s Int64) Print(withParsingInfo bool) string {
	if ps, ok := s.AsString(); ok {
		return *ps
	}
	return fmt.Sprintf("%v", s)
}

func (s Int64) AsInt() (*int64, bool) {
	v := int64(s)
	return &v, true
}

func (s Int64) AsFloat() (*float64, bool) {
	v := float64(s)
	return &v, true
}

//*************** Float64 *****************
func (s Float64) Eval(scope Scope, ctxt EvalContext) (term TermI, err error) {
	return s, nil
}

func (s Float64) AsList() (*ListT, bool) {
	return nil, false
}

func (s Float64) AsSymbol() (*Symbol, bool) {
	return nil, false
}

func (s Float64) AsNumber() (NumberI, bool) {
	return s, true
}

func (s Float64) AsString() (*string, bool) {
	return nil, false
}

func (s Float64) AsBool() (bool, bool) {
	return false, false
}

func (s Float64) AsLambda() (*LambdaT, bool) {
	return nil, false
}

func (s Float64) IsNil() bool {
	return false
}

func (s Float64) Type() string {
	return FloatType
}
func (s Float64) Print(withParsingInfo bool) string {
	if ps, ok := s.AsString(); ok {
		return *ps
	}
	return fmt.Sprintf("%v", s)
}

func (s Float64) AsInt() (*int64, bool) {
	v := int64(s)
	return &v, true
}

func (s Float64) AsFloat() (*float64, bool) {
	v := float64(s)
	return &v, true
}

//*************** Nil *****************
func (s Nil) Eval(scope Scope, ctxt EvalContext) (term TermI, err error) {
	return s, nil
}

func (s Nil) AsList() (*ListT, bool) {
	return nil, false
}

func (s Nil) AsSymbol() (*Symbol, bool) {
	return nil, false
}

func (s Nil) AsNumber() (NumberI, bool) {
	return nil, false
}

func (s Nil) AsString() (*string, bool) {
	return nil, false
}

func (s Nil) AsBool() (bool, bool) {
	return false, false
}

func (s Nil) AsLambda() (*LambdaT, bool) {
	return nil, false
}

func (s Nil) IsNil() bool {
	return true
}

func (s Nil) Type() string {
	return NilType
}
func (s Nil) Print(withParsingInfo bool) string {
	if ps, ok := s.AsString(); ok {
		return *ps
	}
	return fmt.Sprintf("%v", s)
}

//*************** Bool *****************
func (s Bool) Eval(scope Scope, ctxt EvalContext) (term TermI, err error) {
	return s, nil
}

func (s Bool) AsList() (*ListT, bool) {
	return nil, false
}

func (s Bool) AsSymbol() (*Symbol, bool) {
	return nil, false
}

func (s Bool) AsNumber() (NumberI, bool) {
	return nil, false
}

func (s Bool) AsString() (*string, bool) {
	return nil, false
}

func (s Bool) AsBool() (bool, bool) {
	return bool(s), true
}

func (s Bool) AsLambda() (*LambdaT, bool) {
	return nil, false
}

func (s Bool) IsNil() bool {
	return false
}

func (s Bool) Type() string {
	return BoolType
}
func (s Bool) Print(withParsingInfo bool) string {
	if ps, ok := s.AsString(); ok {
		return *ps
	}
	return fmt.Sprintf("%v", s)
}

//*************** LambdaT *****************
func (s LambdaT) Eval(scope Scope, ctxt EvalContext) (term TermI, err error) {
	return s, nil
}

func (s LambdaT) AsList() (*ListT, bool) {
	return nil, false
}

func (s LambdaT) AsSymbol() (*Symbol, bool) {
	return nil, false
}

func (s LambdaT) AsNumber() (NumberI, bool) {
	return nil, false
}

func (s LambdaT) AsString() (*string, bool) {
	return nil, false
}

func (s LambdaT) AsBool() (bool, bool) {
	return false, false
}

func (s LambdaT) AsLambda() (*LambdaT, bool) {
	return &s, true
}

func (s LambdaT) IsNil() bool {
	return false
}

func (s LambdaT) Type() string {
	return LambdaType
}
func (s LambdaT) Print(withParsingInfo bool) string {
	if ps, ok := s.AsString(); ok {
		return *ps
	}
	return fmt.Sprintf("%v", s)
}
