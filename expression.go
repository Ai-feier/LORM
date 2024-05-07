package lorm

// RawExpr 原生 sql 语句
type RawExpr struct {
	raw  string
	args []any
}

func (r RawExpr) selectedAlias() string {
	return ""
}

func (r RawExpr) fieldName() string {
	return ""
}

func (r RawExpr) target() TableReference {
	return nil
}

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

// binaryExpr 带有关系的表达式
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

type SubqueryExpr struct {
	s Subquery
	// 谓词，ALL，ANY 或者 SOME
	pred string
}

func (SubqueryExpr) expr() {}

func Any(sub Subquery) SubqueryExpr {
	return SubqueryExpr{
		s: sub,
		pred: "ANY",
	}
}

func All(sub Subquery) SubqueryExpr {
	return SubqueryExpr{
		s: sub,
		pred: "ALL",
	}
}

func Some(sub Subquery) SubqueryExpr {
	return SubqueryExpr{
		s: sub,
		pred: "SOME",
	}
}
