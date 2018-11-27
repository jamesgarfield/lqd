package lqd

import "testing"

func Test_ParseCondition(t *testing.T) {

	shouldParse := []struct {
		exp    string
		expect Condition
	}{
		{"eq.100", Condition{EQ, []interface{}{"100"}}},
		{"gt.100", Condition{GT, []interface{}{"100"}}},
		{"lt.100", Condition{LT, []interface{}{"100"}}},
		{"lt.100.05", Condition{LT, []interface{}{"100.05"}}},
		{"gte.100", Condition{GTE, []interface{}{"100"}}},
		{"lte.100", Condition{LTE, []interface{}{"100"}}},
		{"not.eq.100", Condition{NEQ, []interface{}{"100"}}},
		{"eq.not.100", Condition{EQ, []interface{}{"not.100"}}},
		{"in.a,b,c", Condition{IN, []interface{}{"a", "b", "c"}}},
		{"in.1.1,2.2,3.3", Condition{IN, []interface{}{"1.1", "2.2", "3.3"}}},
		{"not.in.1.1,2.2,3.3", Condition{NOTIN, []interface{}{"1.1", "2.2", "3.3"}}},
		{"is.true", Condition{IS, []interface{}{true}}},
		{"is.false", Condition{IS, []interface{}{false}}},
		{"not.is.true", Condition{ISNOT, []interface{}{true}}},
		{"is.null", Condition{IS, []interface{}{nil}}},
		{"not.null", Condition{ISNOT, []interface{}{nil}}},
		{"not.false", Condition{ISNOT, []interface{}{false}}},
		{"not.true", Condition{ISNOT, []interface{}{true}}},
	}

	for _, x := range shouldParse {
		p, err := ParseCondition(x.exp)
		if err != nil {
			t.Errorf("Expected no error parsing Condition %s, got %+v", x.exp, err)
		}
		if !p.Equals(x.expect) {
			t.Errorf("Expected %+v to equal %+v", p, x.expect)
		}
	}

	shouldError := []string{
		"100",
		"is.not.null",
		"not.is.nil",
	}

	for _, exp := range shouldError {
		p, err := ParseCondition(exp)
		if err == nil {
			t.Errorf("Expected error, got value %+v", p)
		}
	}
}
