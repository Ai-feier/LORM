//go:build v9
package orm

type RawExpr struct {
	raw string
	args []any 
}

func (r RawExpr) selectable() {}

func (r RawExpr) expr() {}

func (r RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}

// Raw 创建一个 RawExpr
func Raw(expr string, args...interface{}) RawExpr {
	return RawExpr{
		raw:  expr,
		args: args,
	}
}
