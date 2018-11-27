package lqd

import (
	"bytes"
	"errors"
)

var (
	ErrInvalidExpression = errors.New("invalid expression")
)

type Expr struct {
	Identifier string
	Condition
}

func (e Expr) Equals(other Expr) bool {
	if e.Identifier != other.Identifier {
		return false
	}
	return e.Condition.Equals(other.Condition)
}

func (e Expr) MarshalText() ([]byte, error) {
	c, err := e.Condition.MarshalText()
	if err != nil {
		return nil, err
	}
	b := append([]byte(e.Identifier), '=')
	return append(b, c...), nil
}

func (e *Expr) UnmarshalText(b []byte) error {
	pos := bytes.IndexRune(b, '=')
	if pos == -1 {
		return ErrInvalidExpression
	}
	id := b[0:pos]
	cond, err := ParseCondition(string(b[pos+1:]))
	if err != nil {
		return err
	}
	e.Identifier = string(id)
	e.Condition = cond
	return nil
}
