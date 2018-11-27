package lqdsql

import (
	"clypd/lqd"
	"reflect"
	"testing"
)

func Test_SelectFromTable(t *testing.T) {
	sql, err := SelectFromTable("frobs")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expect := "SELECT *\nFROM \"frobs\""
	if sql.Query != expect {
		t.Errorf("Expected query `%s` got %s", expect, sql.Query)
	}

	expr := lqd.Expr{"id", lqd.Condition{lqd.EQ, []interface{}{1}}}
	sql, err = SelectFromTable("frobs", expr)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expect = "SELECT *\nFROM \"frobs\"\nWHERE \"id\" = $1"
	if sql.Query != expect {
		t.Errorf("Expected query `%s` got %s", expect, sql.Query)
	}
	params := []interface{}{1}
	if !reflect.DeepEqual(sql.Params, params) {
		t.Errorf("Expected params `%+v`, got %+v", params, sql.Params)
	}
}

func Test_SelectFromSubquery(t *testing.T) {
	subquery := SQL{
		Query:  "SELECT user.id, user.name, account.name AS account_name\nFROM user JOIN account ON account.id = user.account_id\nWHERE account.marketplace_id = $1",
		Params: []interface{}{3},
	}
	sql, err := SelectFromSubQuery(subquery)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expect := "SELECT *\nFROM (SELECT user.id, user.name, account.name AS account_name\nFROM user JOIN account ON account.id = user.account_id\nWHERE account.marketplace_id = $1)"
	if sql.Query != expect {
		t.Errorf("Expected query `%s` got %s", expect, sql.Query)
	}

	expr := lqd.Expr{"id", lqd.Condition{lqd.EQ, []interface{}{1}}}
	sql, err = SelectFromSubQuery(subquery, expr)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expect = "SELECT *\nFROM (SELECT user.id, user.name, account.name AS account_name\nFROM user JOIN account ON account.id = user.account_id\nWHERE account.marketplace_id = $1)\nWHERE \"id\" = $2"
	if sql.Query != expect {
		t.Errorf("Expected query \n`%s`\ngot\n%s", expect, sql.Query)
	}
	params := []interface{}{3, 1}
	if !reflect.DeepEqual(sql.Params, params) {
		t.Errorf("Expected params `%+v`, got %+v", params, sql.Params)
	}
}
