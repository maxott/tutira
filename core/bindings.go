package core

import (
	"fmt"
)

// var globalBindings = map[string]TermI{}
var coreFunction = map[string]LambdaT{}
var namespaceFunctions = map[string]map[string]LambdaT{}

func AddNamespaceFunction(ns string, name string, fn LambdaT, doc string) {
	AddAliasedNamespaceFunction(ns, []string{name}, fn, doc)
}

func AddAliasedNamespaceFunction(ns string, names []string, fn LambdaT, doc string) {
	nsf, ok := namespaceFunctions[ns]
	if !ok {
		nsf = make(map[string]LambdaT)
		namespaceFunctions[ns] = nsf
	}
	for _, name := range names {
		nsf[name] = fn
		coreFunction[fmt.Sprintf("%s/%s", ns, name)] = fn
	}
}

// func addBuiltinBinding(name string, term TermI) {
// 	globalBindings[name] = term
// }

func addCoreFunction(name string, fn LambdaT, doc string) {
	coreFunction[name] = fn
}

func addAliasedCoreFunction(names []string, fn LambdaT, doc string) {
	for _, name := range names {
		coreFunction[name] = fn
	}
}
