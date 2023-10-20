//go:build v2
package v2

type Column struct {
	name string
}

func C(name string) Column {
	return Column{
		name: name,
	}
}

// 继承表达式接口
func (c Column) expr() {}

type value struct {
	val any
}

func (v value) expr() {}

// 将任意值封装为 Expression
func valueOf(val any) value {
	return value{
		val: val,
	}
}

func(c Column) EQ(ary any) Predicate {
	return Predicate{
		left: c,
		op: opEQ,
		right: exprOf(ary),
	}
}

func (c Column) LT(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opLT,
		right: exprOf(arg),
	}
}

func (c Column) GT(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opGT,
		right: exprOf(arg),
	}
}
