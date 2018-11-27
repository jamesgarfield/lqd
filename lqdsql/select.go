package lqdsql

import (
	"clypd/lqd"
	"fmt"

	"github.com/lib/pq"
)

func SelectFromTable(table string, exprs ...lqd.Expr) (SQL, error) {
	if len(exprs) == 0 {
		return SQL{
			Query: fmt.Sprintf("SELECT *\nFROM %s", pq.QuoteIdentifier(table)),
		}, nil
	}

	where, err := BuildWhere(0, exprs...)
	if err != nil {
		return SQL{}, err
	}
	return SQL{
		Query:  fmt.Sprintf("SELECT *\nFROM %s\nWHERE %s", pq.QuoteIdentifier(table), where.Query),
		Params: where.Params,
	}, nil
}

func SelectFromSubQuery(subquery SQL, exprs ...lqd.Expr) (SQL, error) {
	if len(exprs) == 0 {
		return SQL{
			Query:  fmt.Sprintf("SELECT *\nFROM (%s)", subquery.Query),
			Params: subquery.Params,
		}, nil
	}

	where, err := BuildWhere(len(subquery.Params), exprs...)
	if err != nil {
		return SQL{}, err
	}

	return SQL{
		Query:  fmt.Sprintf("SELECT *\nFROM (%s)\nWHERE %s", subquery.Query, where.Query),
		Params: append(subquery.Params, where.Params...),
	}, nil
}
