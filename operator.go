package lqd

type Operator interface {
	Op() string
	symbol() symbol
}

type symbol struct {
	str, sym string
}

var (
	eq   = symbol{"eq", "="}
	gt   = symbol{"gt", ">"}
	lt   = symbol{"lt", "<"}
	gte  = symbol{"gte", ">="}
	lte  = symbol{"lte", "<="}
	in   = symbol{"in", "IN"}
	is   = symbol{"is", "IS"}
	not  = symbol{"not", "NOT"}
	tru  = symbol{"true", "TRUE"}
	fals = symbol{"false", "FALSE"}
	null = symbol{"null", "NULL"}
)

var (
	EQ    Operator = eq
	NEQ   Operator = symbol{"not.eq", "!="}
	GT    Operator = gt
	LT    Operator = lt
	GTE   Operator = gte
	LTE   Operator = lte
	IN    Operator = in
	NOTIN Operator = symbol{"not.in", "NOT IN"}
	IS    Operator = is
	ISNOT Operator = symbol{"not.is", "IS NOT"}
)

func (s symbol) String() string {
	return s.str
}

func (s symbol) Op() string {
	return string(s.sym)
}

func (s symbol) symbol() symbol {
	return s
}
