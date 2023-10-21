package LORM

type RawExpr struct {
	raw  string
	args []any
}

func (r RawExpr) selectable() {}

func (r RawExpr) expr() {}

func (r RawExpr) assign() {}

func (r RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}

// Raw 创建一个 RawExpr
func Raw(expr string, args ...interface{}) RawExpr {
	return RawExpr{
		raw:  expr,
		args: args,
	}
}

type binaryExpr struct {
	left  Expression
	op    op
	right Expression
}

func (binaryExpr) expr() {}

type MathExpr binaryExpr

func (m MathExpr) Add(val interface{}) MathExpr {
	return MathExpr{
		left:  m,
		op:    opAdd,
		right: valueOf(val),
	}
}

func (m MathExpr) Multi(val interface{}) MathExpr {
	return MathExpr{
		left:  m,
		op:    opMulti,
		right: valueOf(val),
	}
}

func (m MathExpr) expr() {}
