package lqdsql

import (
	"clypd/lqd"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/lib/pq"
)

var (
	ErrIdentifier = errors.New("invalid identifer")
	ErrOneValue   = errors.New("exactly one value required")
	ErrInvalidIS  = errors.New("invalid IS/ISNOT expression")

	validIdentifier = regexp.MustCompile("^[a-zA-Z][_a-zA-Z0-9]*$")
)

type SQL struct {
	Query  string
	Params []interface{}
}

// BuildWhere constructs a WHERE clause from a list of lqd Expressions
// The paramOffset argument provides a way to indicate how many SQL parameters will precede the WHERE clause
// An instance of this would be if one were constructing a WHERE clause on top of a sub-select that takes parameters
func BuildWhere(paramOffset int, exprs ...lqd.Expr) (SQL, error) {
	var where []string
	sql := SQL{}
	for _, e := range exprs {
		if err := validateIdentifier(e.Identifier); err != nil {
			return SQL{}, err
		}
		ident := pq.QuoteIdentifier(e.Identifier)

		var s string
		switch e.Operator {
		case lqd.IN, lqd.NOTIN:
			ordinals := []string{}
			for _, v := range e.Values {
				sql.Params = append(sql.Params, v)
				ordinals = append(ordinals, fmt.Sprintf("$%d", len(sql.Params)+paramOffset))
			}
			s = fmt.Sprintf("%s %s (%s)", ident, e.Op(), strings.Join(ordinals, ","))
		case lqd.IS, lqd.ISNOT:
			if len(e.Values) != 1 {
				return SQL{}, ErrOneValue
			}
			v := e.Values[0]
			switch t := v.(type) {
			case nil:
				s = fmt.Sprintf("%s %s NULL", ident, e.Op())
			case bool:
				if t {
					s = fmt.Sprintf("%s %s TRUE", ident, e.Op())
				} else {
					s = fmt.Sprintf("%s %s FALSE", ident, e.Op())
				}
			default:
				return SQL{}, ErrInvalidIS
			}

		default:
			if len(e.Values) != 1 {
				return SQL{}, ErrOneValue
			}

			sql.Params = append(sql.Params, e.Values...)
			s = fmt.Sprintf("%s %s $%d", ident, e.Op(), len(sql.Params)+paramOffset)
		}
		where = append(where, s)
	}
	sql.Query = fmt.Sprintf("%s", strings.Join(where, "\nAND "))
	return sql, nil
}

func ConcatSQL(sql ...SQL) SQL {
	switch len(sql) {
	case 0:
		return SQL{}
	case 1:
		return sql[0]
	case 2:
		return SQL{
			Query:  fmt.Sprintf("%s\n%s", sql[0].Query, sql[1].Query),
			Params: append(sql[0].Params, sql[1].Params...),
		}
	default:
		return ConcatSQL(sql[0], ConcatSQL(sql[1:]...))
	}
}

func validateIdentifier(ident string) error {
	if !validIdentifier.MatchString(ident) {
		return ErrIdentifier
	}
	return nil
}
