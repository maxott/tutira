package relalg

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/maxott/welisp/core"
)

const DB_URL_DEF = "postgres://localhost:5432/blue_growth?sslmode=disable"

func init() {
	core.AddNamespaceFunction(NAMESPACE, "query", queryFn, "missing docs")
}

func queryFn(scope core.Scope, params *core.ListT, ctxt core.EvalContext) (core.TermI, error) {
	h := core.ParamsHelper(params, true, 1, 2, scope, ctxt)
	query := h.NextAsString()
	var url string
	if h.HasNext() {
		url = h.NextAsString()
	}
	if h.Error() != nil {
		return nil, h.Error()
	}
	return EvalQuery(query, url)
}

func EvalQuery(query string, url string) (res core.TermI, err error) {
	if url == "" {
		var ok bool
		if url, ok = os.LookupEnv("DB_URL"); !ok {
			url = DB_URL_DEF
		}
	}
	if url == "" {
		err = fmt.Errorf("missing database url")
		return
	}

	var adapter DbAdapterI
	if strings.HasPrefix(url, "postgres:") {
		adapter, err = PostgresAdapter(url)
	} else {
		err = fmt.Errorf("unknown database '%s'", url)
	}
	if err != nil {
		return
	}
	return adapter.Query(context.Background(), query)
}

type DbAdapterI interface {
	Query(ctx context.Context, query string, args ...interface{}) (res core.TermI, err error)
}
