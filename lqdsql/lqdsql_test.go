package lqdsql

import (
	"clypd/lqd"
	"reflect"
	"testing"
)

func Test_BuildWhere(t *testing.T) {
	exprs := []lqd.Expr{
		{"start_date", lqd.Condition{lqd.GTE, []interface{}{"2015-01-01"}}},
		{"end_date", lqd.Condition{lqd.LT, []interface{}{"2017-03-01"}}},
		{"status", lqd.Condition{lqd.IN, []interface{}{"PLANNING", "PENDING"}}},
	}
	sql, err := BuildWhere(0, exprs...)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expect := "\"start_date\" >= $1\nAND \"end_date\" < $2\nAND \"status\" IN ($3,$4)"
	if sql.Query != expect {
		t.Errorf("Expected query \n %s \n got \n %s", expect, sql.Query)
	}
	params := []interface{}{"2015-01-01", "2017-03-01", "PLANNING", "PENDING"}
	if !reflect.DeepEqual(sql.Params, params) {
		t.Errorf("Expected params `%+v`, got %+v", params, sql.Params)
	}
}

func Test_ConcatSQL(t *testing.T) {
	a := SQL{
		Query:  "A",
		Params: []interface{}{1},
	}
	b := SQL{
		Query:  "B",
		Params: []interface{}{2},
	}
	c := SQL{
		Query:  "C",
		Params: []interface{}{3},
	}

	tests := []struct {
		sql    []SQL
		expect SQL
	}{
		{sql: []SQL{}, expect: SQL{}},
		{sql: []SQL{a}, expect: a},
		{sql: []SQL{a, b}, expect: SQL{
			Query:  "A\nB",
			Params: []interface{}{1, 2},
		}},
		{sql: []SQL{a, b, c}, expect: SQL{
			Query:  "A\nB\nC",
			Params: []interface{}{1, 2, 3},
		}},
	}

	for i, x := range tests {
		result := ConcatSQL(x.sql...)
		if !reflect.DeepEqual(result, x.expect) {
			t.Errorf("expected result at index %d to be %+v, got %+v", i, x.expect, result)
		}
	}
}

func Test_validateIdentifier(t *testing.T) {
	valid := []string{
		"name",
		"table_id",
		"why_so_seriously_long_name",
		"column2",
		"Caps_are_ok",
	}
	for _, ident := range valid {
		if err := validateIdentifier(ident); err != nil {
			t.Errorf("Expected %s to be valid", ident)
		}
	}

	invalid := []string{
		"",
		"2bad4you",
		"µ†ƒ",
		"_not_ok",
		"-why-",
		";DROP table users",
		"`escape",
		"valid_till_the_end™",
	}
	for _, ident := range invalid {
		if err := validateIdentifier(ident); err == nil {
			t.Errorf("Expected %s to be invalid", ident)
		}
	}
}
