package lqd

import (
	"fmt"
	"unicode/utf8"
)

type ErrInvalidCondition string

func (err ErrInvalidCondition) Error() string {
	return fmt.Sprintf("invalid condition: %s", string(err))
}

func ParseCondition(val string) (Condition, error) {
	lex := &lexer{
		val:    val,
		result: Condition{},
	}

	var err error
	state := initState
	for {
		if state, err = state(lex); err != nil {
			return Condition{}, err
		} else if state == nil {
			return lex.result, nil
		}
	}
}

type stateFn func(*lexer) (stateFn, error)

type lexer struct {
	val    string
	offset int
	result Condition
}

func (l *lexer) NextToken(delims ...rune) (string, bool) {
	if l.offset == len(l.val) {
		return "", false
	}
	var tok string
	if len(delims) == 0 {
		tok = l.val[l.offset:]
		l.offset = len(l.val)
		return tok, true
	}
	subval := l.val[l.offset:]
	index := l.offset
	for _, sym := range subval {
		index = index + utf8.RuneLen(sym)
		for _, r := range delims {
			if sym == r {
				tok = l.val[l.offset : index-1]
				l.offset = index
				return tok, true
			}
		}
	}
	tok = l.val[l.offset:]
	l.offset = len(l.val)
	return tok, true
}

func initState(lex *lexer) (stateFn, error) {
	tok, ok := lex.NextToken('.')
	if !ok {
		return nil, ErrInvalidCondition(lex.val)
	}
	switch tok {
	case eq.str:
		lex.result.Operator = EQ
		return hasOperator, nil
	case gt.str:
		lex.result.Operator = GT
		return hasOperator, nil
	case lt.str:
		lex.result.Operator = LT
		return hasOperator, nil
	case gte.str:
		lex.result.Operator = GTE
		return hasOperator, nil
	case lte.str:
		lex.result.Operator = LTE
		return hasOperator, nil
	case is.str:
		lex.result.Operator = IS
		return isIs, nil
	case in.str:
		lex.result.Operator = IN
		return isIn, nil
	case not.str:
		return hasNot, nil
	default:
		return nil, ErrInvalidCondition(lex.val)
	}
}

func hasNot(lex *lexer) (stateFn, error) {
	tok, ok := lex.NextToken('.')
	if !ok {
		return nil, ErrInvalidCondition(lex.val)
	}
	switch tok {
	case eq.str:
		lex.result.Operator = NEQ
		return hasOperator, nil
	case gt.str:
		lex.result.Operator = LTE
		return hasOperator, nil
	case lt.str:
		lex.result.Operator = GTE
		return hasOperator, nil
	case gte.str:
		lex.result.Operator = LT
		return hasOperator, nil
	case lte.str:
		lex.result.Operator = GT
		return hasOperator, nil
	case is.str:
		lex.result.Operator = ISNOT
		return isIs, nil
	case in.str:
		lex.result.Operator = NOTIN
		return isIn, nil
	case tru.str:
		lex.result.Operator = ISNOT
		lex.result.Values = []interface{}{true}
		return nil, nil
	case fals.str:
		lex.result.Operator = ISNOT
		lex.result.Values = []interface{}{false}
		return nil, nil
	case null.str:
		lex.result.Operator = ISNOT
		lex.result.Values = []interface{}{nil}
		return nil, nil
	default:
		return nil, ErrInvalidCondition(lex.val)
	}
}

func hasOperator(lex *lexer) (stateFn, error) {
	tok, ok := lex.NextToken()
	if !ok {
		return nil, ErrInvalidCondition(lex.val)
	}
	lex.result.Values = append(lex.result.Values, tok)
	return nil, nil
}

func isIn(lex *lexer) (stateFn, error) {
	for tok, ok := lex.NextToken(','); ok; tok, ok = lex.NextToken(',') {
		lex.result.Values = append(lex.result.Values, tok)
	}
	return nil, nil
}

func isIs(lex *lexer) (stateFn, error) {
	tok, ok := lex.NextToken()
	if !ok {
		return nil, ErrInvalidCondition(lex.val)
	}
	switch tok {
	case tru.str:
		lex.result.Values = []interface{}{true}
	case fals.str:
		lex.result.Values = []interface{}{false}
	case null.str:
		lex.result.Values = []interface{}{nil}
	default:
		return nil, ErrInvalidCondition(lex.val)
	}
	return nil, nil
}
