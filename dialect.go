package lorm

import "github.com/Ai-feier/lorm/internal/errs"

var (
	MySQL   Dialect = &MysqlDialect{}
	SQLite3 Dialect = &Sqlite3Dialect{}
)

type Dialect interface {
	// quoter 返回一个引号，引用列名，表名的引号
	quoter() byte
	// buildUpsert 构造插入冲突部分
	buildUpsert(b *builder, odk *Upsert) error
}

type standardSQL struct {
}

func (s *standardSQL) quoter() byte {
	// TODO implement me
	panic("implement me")
}

func (s *standardSQL) buildUpsert(b *builder,
	odk *Upsert) error {
	panic("implement me")
}

type MysqlDialect struct {
	standardSQL
}

func (m *MysqlDialect) quoter() byte {
	return '`'
}

func (m *MysqlDialect) buildUpsert(b *builder,
	odk *Upsert) error {
	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for idx, a := range odk.assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		switch assign := a.(type) {
		case Column:
			fd, ok := b.model.FieldMap[assign.name]
			if !ok {
				return errs.NewErrUnknownField(assign.name)
			}
			b.quote(fd.ColName)
			b.sb.WriteString("=VALUES(")
			b.quote(fd.ColName)
			b.sb.WriteByte(')')
		case Assignment:
			err := b.buildColumn(nil, assign.column)
			if err != nil {
				return err
			}
			b.sb.WriteString("=")
			return b.buildExpression(assign.val)
		default:
			return errs.NewErrUnsupportedAssignableType(a)
		}
	}
	return nil
}

type Sqlite3Dialect struct {
	standardSQL
}

func (s *Sqlite3Dialect) quoter() byte {
	return '`'
}

func (s *Sqlite3Dialect) buildUpsert(b *builder,
	odk *Upsert) error {
	b.sb.WriteString(" ON CONFLICT")
	if len(odk.conflictColumns) > 0 {
		b.sb.WriteByte('(')
		for i, col := range odk.conflictColumns {
			if i > 0 {
				b.sb.WriteByte(',')
			}
			err := b.buildColumn(nil, col)
			if err != nil {
				return err
			}
		}
		b.sb.WriteByte(')')
	}
	b.sb.WriteString(" DO UPDATE SET ")

	for idx, a := range odk.assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		switch assign := a.(type) {
		case Column:
			fd, ok := b.model.FieldMap[assign.name]
			if !ok {
				return errs.NewErrUnknownField(assign.name)
			}
			b.quote(fd.ColName)
			b.sb.WriteString("=excluded.")
			b.quote(fd.ColName)
		case Assignment:
			err := b.buildColumn(nil, assign.column)
			if err != nil {
				return err
			}
			b.sb.WriteString("=")
			return b.buildExpression(assign.val)
		default:
			return errs.NewErrUnsupportedAssignableType(a)
		}
	}
	return nil
}
