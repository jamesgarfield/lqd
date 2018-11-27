package lqdsql

import (
	"clypd/lqd"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

var (
	ErrNoUpdateValues = errors.New("requires at least one update value")
)

func UpdateOnTable(table string, values map[string]interface{}, exprs ...lqd.Expr) (SQL, error) {
	update, err := buildUpdate(table, values)
	if err != nil {
		return SQL{}, err
	}

	var where SQL
	if len(exprs) > 0 {
		if where, err = BuildWhere(len(update.Params), exprs...); err != nil {
			return SQL{}, err
		}
	}

	return ConcatSQL(update, where, SQL{Query: "RETURNING *"}), nil
}

func buildUpdate(table string, values map[string]interface{}) (SQL, error) {
	if len(values) == 0 {
		return SQL{}, ErrNoUpdateValues
	}

	cols := make([]string, len(values))
	params := make([]interface{}, len(values))
	paramIndex := 0
	for col, val := range values {
		if err := validateIdentifier(col); err != nil {
			return SQL{}, err
		}
		cols[paramIndex] = fmt.Sprintf("%s = $%d", pq.QuoteIdentifier(col), paramIndex+1)
		params[paramIndex] = val
		paramIndex++
	}
	update := fmt.Sprintf("UPDATE %s\nSET\n%s",
		pq.QuoteIdentifier(table),
		strings.Join(cols, ",\n"),
	)
	return SQL{
		Query:  update,
		Params: params,
	}, nil
}
