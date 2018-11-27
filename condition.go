package lqd

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNoConditionValues = errors.New("no values in condition")
)

type Condition struct {
	Operator
	Values []interface{}
}

func (e Condition) MarshalText() ([]byte, error) {
	var (
		vals string
		err  error
	)
	switch len(e.Values) {
	case 0:
		err = ErrNoConditionValues
	case 1:
		vals = fmt.Sprintf("%v", e.Values[0])
	default:
		vala := make([]string, len(e.Values))
		for i, v := range e.Values {
			vala[i] = fmt.Sprintf("%v", v)
		}
		vals = strings.Join(vala, ",")
	}
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%s.%s", e.symbol().String(), vals)), nil
}

func (e *Condition) UnmarshalText(b []byte) error {
	exp, err := ParseCondition(string(b))
	if err != nil {
		return err
	}
	*e = exp
	return nil
}

func (p Condition) Equals(other Condition) bool {
	if p.Operator.symbol() != other.Operator.symbol() ||
		len(p.Values) != len(other.Values) {
		return false
	}
	for i := range p.Values {
		if p.Values[i] != other.Values[i] {
			return false
		}
	}
	return true
}
