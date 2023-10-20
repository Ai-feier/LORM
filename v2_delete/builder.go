//go:build del
package v2_delete

import (
	"fmt"
	"strings"
)

type Builder struct {
	sb *strings.Builder
	args []any
}

func (b *Builder) buildExpression(e Expression) error {
	if e == nil {
		return nil
	}
	switch exp := e.(type) {
	case Column:
		b.sb.WriteByte('`')
		b.sb.WriteString(exp.name)
		b.sb.WriteByte('`')
	case value:
		b.sb.WriteByte('?')
		b.args = append(b.args, exp.val)
	case Predicate:
		_, lp := exp.left.(Predicate)
		if lp {
			b.sb.WriteByte('(')
		}
		if err := b.buildExpression(exp.left); err != nil {
			return err
		}
		if lp {
			b.sb.WriteByte(')')
		}

		b.sb.WriteByte(' ')
		b.sb.WriteString(exp.op.String())
		b.sb.WriteByte(' ')

		_, rp := exp.right.(Predicate)
		if rp {
			b.sb.WriteByte('(')
		}
		if err := b.buildExpression(exp.right); err != nil {
			return err
		}
		if rp {
			b.sb.WriteByte(')')
		}
	default:
		return fmt.Errorf("orm: 不支持的表达式 %v", exp)
	}
	return nil
}
