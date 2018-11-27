package lqdurl

import (
	"clypd/lqd"
	"net/http"
	"net/url"
	"sort"
)

func ParseLQDRequest(req *http.Request) ([]lqd.Expr, error) {
	query := req.URL.Query()
	queryColumns := sortedColumns(query)

	var exprs []lqd.Expr
	for _, col := range queryColumns {
		queries := query[col]
		for _, q := range queries {
			condition, err := lqd.ParseCondition(q)
			if err != nil {
				return nil, err
			}
			exprs = append(exprs, lqd.Expr{col, condition})
		}
	}
	return exprs, nil
}

func ParseSimpleRequest(req *http.Request) ([]lqd.Expr, error) {
	query := req.URL.Query()
	queryColumns := sortedColumns(query)

	exprs := []lqd.Expr{}
	for _, col := range queryColumns {
		values := []interface{}{}
		for _, v := range query[col] {
			values = append(values, v)
		}
		op := lqd.IN
		if len(values) == 1 {
			op = lqd.EQ
		}
		e := lqd.Expr{
			col,
			lqd.Condition{
				Operator: op,
				Values:   values,
			},
		}
		exprs = append(exprs, e)
	}
	return exprs, nil
}

func sortedColumns(q url.Values) []string {
	columns := []string{}
	for col := range q {
		columns = append(columns, col)
	}
	sort.Strings(columns)
	return columns
}
