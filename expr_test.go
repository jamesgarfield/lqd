package lqd

import "testing"

func Test_Expr_UnmarshalText(t *testing.T) {
	shouldWork := []struct {
		expr   string
		expect Expr
	}{
		{"id=eq.4", Expr{"id", Condition{EQ, []interface{}{"4"}}}},
		{"status=in.PLANNING,PENDING,ERROR", Expr{"status", Condition{IN, []interface{}{"PLANNING", "PENDING", "ERROR"}}}},
		{"age=gte.21", Expr{"age", Condition{GTE, []interface{}{"21"}}}},
	}

	for i, x := range shouldWork {
		e := Expr{}
		err := e.UnmarshalText([]byte(x.expr))
		if err != nil {
			t.Errorf("at index %d, unexpected error : %+v", i, err)
		}
		if !e.Equals(x.expect) {
			t.Errorf("at index %d, expected %+v, got %+v", i, x.expect, e)
		}
	}

	shouldFail := []string{
		"id==eq.4",
		"id=eq",
		"eq.4",
		"id=",
		"id",
		"=id=eq.4",
		"id:eq.4",
		"id:eq.4=",
	}

	for i, x := range shouldFail {
		e := Expr{}
		err := e.UnmarshalText([]byte(x))
		if err == nil {
			t.Errorf("at index %d, expected error, got none", i)
		}
	}
}

func Test_Expr_MarshalText(t *testing.T) {
	shouldWork := []struct {
		Expr
		expect string
	}{
		{Expr{"id", Condition{EQ, []interface{}{4}}}, "id=eq.4"},
		{Expr{"status", Condition{IN, []interface{}{"PLANNING", "PENDING"}}}, "status=in.PLANNING,PENDING"},
		{Expr{"age", Condition{GTE, []interface{}{21}}}, "age=gte.21"},
	}

	for i, x := range shouldWork {
		b, err := x.Expr.MarshalText()
		if err != nil {
			t.Errorf("at index %d, unexpected err: %+v", i, err)
		}
		if string(b) != x.expect {
			t.Errorf("expected %s, got %s", x.expect, b)
		}
	}
}
