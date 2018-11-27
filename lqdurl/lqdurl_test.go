package lqdurl

import (
	"clypd/lqd"
	"net/http"
	"testing"
)

func Test_ParseLQDRequest(t *testing.T) {
	shouldWork := []struct {
		url    string
		expect []lqd.Expr
	}{
		{"http://api.net/foo", []lqd.Expr{}},
		{"http://api.net/foo?id=eq.10", []lqd.Expr{{"id", lqd.Condition{lqd.EQ, []interface{}{"10"}}}}},
		{"http://api.net/foo?id=eq.10&status=in.PLANNING,PENDING", []lqd.Expr{
			{"id", lqd.Condition{lqd.EQ, []interface{}{"10"}}},
			{"status", lqd.Condition{lqd.IN, []interface{}{"PLANNING", "PENDING"}}},
		}},
		{"http://api.net/bar?age=gte.21&approved=is.true", []lqd.Expr{
			{"age", lqd.Condition{lqd.GTE, []interface{}{"21"}}},
			{"approved", lqd.Condition{lqd.IS, []interface{}{true}}},
		}},
	}

	for i, x := range shouldWork {
		r, err := http.NewRequest("GET", x.url, nil)
		if err != nil {
			t.Errorf("could not parse url: %s", x.url)
			t.FailNow()
		}
		exprs, err := ParseLQDRequest(r)
		if err != nil {
			t.Errorf("at index %d, unexpected error %+v", i, err)
		}

		for e, expr := range exprs {
			expect := x.expect[e]
			if !expr.Equals(expect) {
				t.Errorf("at test index %d, expected expr %d to be %+v, got %+v", i, e, expect, expr)
			}
		}
	}
}

func Test_ParseSimpleRequest(t *testing.T) {
	shouldWork := []struct {
		url    string
		expect []lqd.Expr
	}{
		{"http://api.net/foo", []lqd.Expr{}},
		{"http://api.net/foo?id=10", []lqd.Expr{{"id", lqd.Condition{lqd.EQ, []interface{}{"10"}}}}},
		{"http://api.net/foo?id=10&id=20&id=30", []lqd.Expr{{"id", lqd.Condition{lqd.IN, []interface{}{"10", "20", "30"}}}}},
		{"http://api.net/foo?id=10&id=20&status=PLANNING&status=PENDING", []lqd.Expr{
			{"id", lqd.Condition{lqd.IN, []interface{}{"10", "20"}}},
			{"status", lqd.Condition{lqd.IN, []interface{}{"PLANNING", "PENDING"}}},
		}},
	}

	for i, x := range shouldWork {
		r, err := http.NewRequest("GET", x.url, nil)
		if err != nil {
			t.Errorf("could not parse url: %s", x.url)
			t.FailNow()
		}
		exprs, err := ParseSimpleRequest(r)
		if err != nil {
			t.Errorf("at index %d, unexpected error %+v", i, err)
		}

		for e, expr := range exprs {
			expect := x.expect[e]
			if !expr.Equals(expect) {
				t.Errorf("at test index %d, expected expr %d to be %+v, got %+v", i, e, expect, expr)
			}
		}
	}
}
