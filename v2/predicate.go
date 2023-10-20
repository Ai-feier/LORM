//go:build v2
package v2

// op 代表操作符
type op string

// 后面可以每次支持新的操作符就加一个
const (
	opEQ  = "="
	opLT  = "<"
	opGT  = ">"
	opAND = "AND"
	opOR  = "OR"
	opNOT = "NOT"
)

// Expression 代表语句，或者语句的部分
type Expression interface {
	expr()
}


// Predicate 代表一个查询条件
// Predicate 可以通过和 Predicate 组合构成复杂的查询条件
type Predicate struct {
	left Expression
	op op 
	right Expression
}

// 同样继承 Expression
func (p Predicate) expr() {}

func exprOf(e any) Expression {
	switch exp := e.(type) {
	case Expression:
		return exp
	default:
		return valueOf(exp)
	}
}

// And Or 实现链式调用特性

func (p Predicate) And(r Predicate) Predicate {
	return Predicate{
		left: p,
		op: opAND,
		right: r,
	}
}

func (p Predicate) Or(r Predicate) Predicate {
	return Predicate{
		left: p,
		op: opOR,
		right: r,
	}
}

func Not(p Predicate) Predicate {
	return Predicate{
		op: opNOT,
		right: p,
	}
}

func (o op) String() string {
	return string(o)
}