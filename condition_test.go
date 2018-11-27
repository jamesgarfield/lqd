package lqd

import "testing"

func Test_Condition_MarshalText(t *testing.T) {
	shouldMarshal := []struct {
		Condition
		expect string
	}{
		{Condition{EQ, []interface{}{"A"}}, "eq.A"},
		{Condition{EQ, []interface{}{"1"}}, "eq.1"},
		{Condition{IN, []interface{}{"A", "B"}}, "in.A,B"},
		{Condition{ISNOT, []interface{}{false}}, "not.is.false"},
	}

	for i, s := range shouldMarshal {
		b, err := s.MarshalText()
		if err != nil {
			t.Errorf("unexpected error at index %d: %+v", i, err)
		}
		if string(b) != s.expect {
			t.Errorf("expected %s, got %s", s.expect, b)
		}
	}
}
