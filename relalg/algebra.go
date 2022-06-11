package relalg

import (
	"fmt"

	"github.com/maxott/welisp/core"
)

const NAMESPACE = "relalg"

func init() {
	core.AddNamespaceFunction(NAMESPACE, "relation", relationFn, "missing docs")
}

type Request struct {
	core.GenericTerm
}

func relationFn(scope core.Scope, params *core.ListT, ctxt core.EvalContext) (core.TermI, error) {
	// h := core.ParamsHelper(params, true, 1, 2, scope, ctxt)

	// printFn := func(withParsingInfo bool) string {
	// 	return fmt.Sprintf("(relalg/relation \"%s\")", name)
	// }
	return &Request{
		core.GenericTerm{
			TermType: "relalg/request",
			// PrintFn:  printFn,
		},
	}, nil
}

func (r *Request) Print(withParsingInfo bool) string {
	name := "foo"
	return fmt.Sprintf("(relalg/request \"%s\")", name)
}
