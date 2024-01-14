## 顶级抽象

```go
// Querier sql 查询语句抽象
type Querier[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)
}

// Executor insert update 语句执行抽象
type Executor interface {
	Exec(ctx context.Context) Result
}

// Query sql 中间结构体
// SQL: sql 语句
// Args: sql 语句中的占位符参数
type Query struct {
	SQL  string
	Args []any
}

// QueryBuilder 构造 sql 语句的抽象
type QueryBuilder interface {
	Build() (*Query, error)
}
```



## 元数据

### Registry 注册中心抽象

```go
// Registry 元数据注册中心的抽象
type Registry interface {
	// Get 查找元数据
	Get(val any) (*Model, error)
	// Register 注册一个模型
	Register(val any, opts ...Option) (*Model, error)
}
```



#### Registry 注册中心底层数据实现

存储 元数据(用户结构体 与 操作数据对象)

```go
type registry struct {
	models sync.Map
}
```



核心方法

```go
func (r *registry) Get(val any) (*Model, error)

func (r *registry) Register(val any, opts ...Option)

// 接受一级结构体指针
// parseModel 支持从标签中提取自定义设置
// 标签形式 orm:"key1=value1,key2=value2"
func (r *registry) parseModel(val any) (*Model, error)

func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error)
```





### Model

用户结构体名 -> Model -> database.sql

```go
type Model struct {
	// TableName 结构体对应的表名
	TableName string
	Fields    []*Field
	FieldMap  map[string]*Field
	ColumnMap map[string]*Field
}

// Field 字段
type Field struct {
	ColName string
	GoName  string
	Type    reflect.Type
	Index   int
	// Offset 相对于对象起始地址的字段偏移量
	Offset uintptr
}
```

#### Model Option

扩展 Model 属性

```go
type Option func(m *Model) error
```
WithColumnName
```go
func WithColumnName(field string, columnName string) Option {
	return func(model *Model) error {
		fd, ok := model.FieldMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		// 注意，这里我们根本没有检测 ColName 会不会是空字符串
		// 因为正常情况下，用户都不会写错
		// 即便写错了，也很容易在测试中发现
		fd.ColName = columnName
		return nil
	}
}
```



### 扩展

#### tag parser

支持的 tag

```go
const (
	tagKeyColumn = "column"
)
```



#### TableName

```go
// TableName 用户实现这个接口来返回自定义的表名
type TableName interface {
	TableName() string
}
```



### 流程

TODO: picture

调用方创建模型注册中心 -> 调用方向数据中心注册元数据

调用方向数据中获取元数据 -> 数据中心查找响应的元数据 (如果没找到就注册该条元数据, 懒加载实现)



#### 注册模型

```go
func (r *registry) Register(val any, opts ...Option) (*Model, error)
	->	func (r *registry) parseModel(val any) (*Model, error)
		-> func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error)
```



#### 获取元数据

```go
func (r *registry) Get(val any) (*Model, error) => 是否已有元数据
	-> func (r *registry) Register(val any, opts ...Option) (*Model, error)
		->	func (r *registry) parseModel(val any) (*Model, error)
			-> func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error)
```



### 核心代码

获取调用方结构体反射类型  (check接受一级结构体指针)

​	-> 逐一解析结构体字段

​		-> 解析其 tag  填充 Field 结构体

```go
// 接受一级结构体指针
// parseModel 支持从标签中提取自定义设置
// 标签形式 orm:"key1=value1,key2=value2"
func (r *registry) parseModel(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Ptr ||
		typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()

	// 获得字段的数量
	numField := typ.NumField()
	fds := make(map[string]*Field, numField)
	fields := make([]*Field, 0, numField)
	colMap := make(map[string]*Field, numField)
	for i := 0; i < numField; i++ {
		fdType := typ.Field(i)
		tags, err := r.parseTag(fdType.Tag)
		if err != nil {
			return nil, err
		}
		colName := tags[tagKeyColumn]
		if colName == "" {
			colName = underscoreName(fdType.Name)
		}
		f := &Field{
			ColName: colName,
			Type:    fdType.Type,
			GoName:  fdType.Name,
			Offset:  fdType.Offset,
			Index:   i,
		}
		fds[fdType.Name] = f
		fields = append(fields, f)
		colMap[colName] = f
	}
	var tableName string
	if tn, ok := val.(TableName); ok {
		tableName = tn.TableName()
	}

	if tableName == "" {
		tableName = underscoreName(typ.Name())
	}

	return &Model{
		TableName: tableName,
		Fields:    fields,
		FieldMap:  fds,
		ColumnMap: colMap,
	}, nil
}
```



## DB

### 定义

```go
type DB struct {
	dialect    Dialect
	r          model.Registry
	db         *sql.DB
	valCreator valuer.Creator
}
```



### 创建数据连接

```go
// Open 创建一个 DB 实例。
// 默认情况下，该 DB 将使用 MySQL 作为方言
// 如果你使用了其它数据库，可以使用 DBWithDialect 指定
func Open(driver string, dsn string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
}

// 语法糖
func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error)
func MustNewDB(driver string, dsn string, opts ...DBOption) *DB
```



### DB Option

```go
type DBOption func(*DB)
```



DBWithRegistry 更换数据中心实现

```go
func DBWithRegistry(r model.Registry) DBOption {
	return func(db *DB) {
		db.r = r
	}
}
```



DBUseReflectValuer 更改数据返回结果实现

```go
func DBUseReflectValuer() DBOption {
	return func(db *DB) {
		db.valCreator = valuer.NewReflectValue
	}
}
```



优雅启动, 确保数据正确连接

```go
func (db *DB) Wait() error {
	err := db.db.Ping()
	for errors.Is(err, driver.ErrBadConn) {
		log.Println("数据库启动中")
		err = db.db.Ping()
	}
	return nil
}
```



## Valuer

### 作用

在返回数据库查询查询结果时, 填充数据



### Interface

```go
// Value 是对结构体实例的内部抽象
type Value interface {
	// Field 返回字段对应的值
	Field(name string) (any, error)
	// SetColumns 设置新值
	SetColumns(rows *sql.Rows) error
}
```



### 函数式工厂构造

```go
type Creator func(val any, meta *model.Model) Value
```



### 实现类

#### 基于 reflect 实现

类型

```go
// reflectValue 基于反射的 Value
type reflectValue struct {
	val  reflect.Value
	meta *model.Model
}
```



构造方法

```go
// NewReflectValue 返回一个封装好的，基于反射实现的 Value
// 输入 val 必须是一个指向结构体实例的指针，而不能是任何其它类型
func NewReflectValue(val interface{}, meta *model.Model) Value {
	return reflectValue{
		val:  reflect.ValueOf(val).Elem(),
		meta: meta,
	}
}
```



填充数据库查询结果

```go
func (r reflectValue) SetColumns(rows *sql.Rows) error

反射获取查询返回结果
用反射填充 valuer 对象值
```





#### 基于 unsafe 实现

类型

```go
type unsafeValue struct {
	// 基准地址
	addr unsafe.Pointer
	meta *model.Model
}
```



构造方法

```go
func NewUnsafeValue(val interface{}, meta *model.Model) Value {
	return unsafeValue{
		addr: unsafe.Pointer(reflect.ValueOf(val).Pointer()),
		meta: meta,
	}
}
```



填充数据库查询结果

```go
func (u unsafeValue) SetColumns(rows *sql.Rows) error

根据 unsafe 基准地址 + offset 在对应位置, 利用反射在当前位置开辟一个新对象
把其地址指向查询返回结果
```



#### 获取 field 的值

```go
func (u unsafeValue) Field(name string) (any, error)
```



#### benchmark 测试

```go
// 执行命令: go test -bench=BenchmarkQuerier_Get -benchmem -benchtime=10000x
// goos: windows
// goarch: amd64
// pkg: LORM/v8
// cpu: Intel(R) Core(TM) i5-10210U CPU @ 1.60GHz
// BenchmarkQuerier_Get/unsafe-8              10000            418712 ns/op            3324 B/op        112 allocs/op
// BenchmarkQuerier_Get/reflect-8             10000           1462437 ns/op            3503 B/op        120 allocs/op
// PASS
// ok      LORM/v8 19.787s
// 可以看出 unsafe 的性能远远快于直接使用 reflect
```



## builer

### Interface

Expression 代表语句，或者语句的部分

```go
// Expression 代表语句，或者语句的部分
type Expression interface {
    expr()
}
```



Assignable 标记接口

```go
// Assignable 标记接口，
// 实现该接口意味着可以用于赋值语句，
// 用于在 UPDATE 和 UPSERT 中
type Assignable interface {
    assign()
}
```



Selectable select 语句中 colume, rawexpr, aggregate 的抽象

```go
// Selectable select 语句中 colume, rawexpr, aggregate 的抽象
type Selectable interface {
    selectable()
}
```



### 实现类

#### Predicate 代表一个查询条件

- Expression

```go
// Predicate 代表一个查询条件
// Predicate 可以通过和 Predicate 组合构成复杂的查询条件
type Predicate struct {
    left  Expression
    op    op
    right Expression
}
```



#### Aggregate 代表聚合函数，例如 AVG, MAX, MIN 等

- Expression 
- Selectable

```go
// Aggregate 代表聚合函数，例如 AVG, MAX, MIN 等
type Aggregate struct {
    fn    string
    arg   string
    alias string
}
```



#### Column 列名

- Expression 
- Selectable
- Assignable 

```go
type Column struct {
    name  string
    alias string
}
```



#### RawExpr 原生 sql 语句

- Expression
- Selectable
- Assignable

```go
// RawExpr 原生 sql 语句
type RawExpr struct {
    raw  string
    args []any
}
```



#### binaryExpr, MathExpr 带有关系的表达式

- Expression

```go
type binaryExpr struct {
    left  Expression
    op    op
    right Expression
}

type MathExpr binaryExpr
func (m MathExpr) Add(val interface{}) MathExpr
func (m MathExpr) Multi(val interface{}) MathExp
```



#### value 代表单独的值, 可单独作为表达式

```go
type value struct {
    val any
}
```





### builder 实现类

```go
type builder struct {
	// 构造 SQL
	sb strings.Builder
	// 存放 SQL 参数
	args []any
	// 存放当前对象的元数据信息
	model *model.Model
	// 方言抽象
	dialect Dialect
	quoter  byte
}
```



#### Dialect 方言

兼容不同数据库的方言

```go
type Dialect interface {
    // quoter 返回一个引号，引用列名，表名的引号
    quoter() byte
    // buildUpsert 构造插入冲突部分
    buildUpsert(b *builder, odk *Upsert) error
}
```



mysql, sqlite3 方言实现

```go
var (
    MySQL   Dialect = &mysqlDialect{}
    SQLite3 Dialect = &sqlite3Dialect{}
)
```



#### 核心方法

```go
// buildColumn 构造列
func (b *builder) buildColumn(fd string) error

// 构造方言的quote
func (b *builder) quote(name string)

// 构造原生表达式
func (b *builder) raw(r RawExpr)

// 构造 sql 语句参数
func (b *builder) addArgs(args ...any)

// 构造条件表达式
func (b *builder) buildPredicates(ps []Predicate) error

// 构造表达式
func (b *builder) buildExpression(e Expression) error

// 构造二分表达式
func (b *builder) buildBinaryExpr(e binaryExpr) error
流程:
err := b.buildSubExpr(e.left)  => 构造做部分
左右部分有可能又是一个 Predicate, 所以递归构造, 到最后为 raw, value


// 构造二分表达式的右半部分
func (b *builder) buildSubExpr(subExpr Expression) error

// 构造 Aggregate 表达式
func (b *builder) buildAggregate(a Aggregate, useAlias bool) error

// 构造 As 别名
func (b *builder) buildAs(alias string)
```





## Selector

### 定义

构造 select 语句发送给数据库处理

用户需传入要查询的结构体(泛型的使用)

采用构造复杂结构的 Build 模式, 链式调用

```go
// Selector 用于构造 SELECT 语句
type Selector[T any] struct {
	builder
	table   string
	db      *DB
	where   []Predicate
	having  []Predicate
	columns []Selectable
	groupBy []Column
	limit   int
	offset  int
}
```



### 构造方法

```gO
// NewSelector 构造 Selector
func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		builder: builder{
			dialect: db.dialect,
			quoter:  db.dialect.quoter(),
		},
		db: db,
	}
}
```



### Build



依次构造: SELECT 列名... FROM 表名 WHERE 条件 GROUP BY 列名 HAVING 条件 LIMIT limit OFFSET offest

列名, 表名有不同数据库的方言

```go
// Build 构造 sql 查询语句, 底层调用 database.sql 查询数据库
func (s *Selector[T]) Build() (*Query, error)
```



其余方法及 Build() 流程

```go
// Select 选择要查询的列
func (s *Selector[T]) Select(cols ...Selectable) *Selector[T]

// From 指定表名，如果是空字符串，那么将会使用默认表名
func (s *Selector[T]) From(tbl string) *Selector[T] 

func (s *Selector[T]) buildColumns() error

func (s *Selector[T]) buildColumn(c Column, useAlias bool) error

// Where 用于构造 WHERE 查询条件。如果 ps 长度为 0，那么不会构造 WHERE 部分
func (s *Selector[T]) Where(ps ...Predicate) *Selector[T]

// GroupBy 设置 group by 子句
func (s *Selector[T]) GroupBy(cols ...Column) *Selector[T]

func (s *Selector[T]) Having(ps ...Predicate) *Selector[T]

func (s *Selector[T]) Offset(offset int) *Selector[T]

func (s *Selector[T]) Limit(limit int) *Selector[T]

Build() 流程:
s.model, err = s.db.r.Get(&t) => 从注册中心获取模型
s.sb.WriteString("SELECT ")  => 构造固定语法
if err = s.buildColumns();   => 构造列
	-> for i, c := range s.columns  => 遍历所有列名
		-> switch val := c.(type)  => 判断列类型
			-> case Column:  => 普通列名构造
				if err := s.buildColumn(val, true);
			-> case Aggregate:  => 构造 Aggregate 类型列
					if err := s.buildAggregate(val, true);
			-> case RawExpr:  => 构造原生 sql 列
					s.raw(val)
s.sb.WriteString(" FROM ")  => 构造固定语法
s.quote(s.model.TableName) || s.sb.WriteString(s.table)  => 构造表名
if len(s.where) > 0; if err = s.buildPredicates(s.where); => 构造 where 条件语句
	-> for i := 1; i < len(ps); i++ {  => 不同的 Predicate And 连接
		p = p.And(ps[i])
	-> b.buildExpression(p)  => 构造表达式
if len(s.groupBy) > 0  => 构造 group by
if len(s.having) > 0   => 构造 having
if s.limit > 0         => 构造 limit
if s.offset > 0        => 构造 offset
```







### Get() 与 GetMulti()

```go
func (s *Selector[T]) Get(ctx context.Context) (*T, error)

func (s *Selector[T]) GetMulti(ctx context.Context) (res []*T, err error)

q, err := s.Build()  => 构造 sql
rows, err := db.QueryContext(ctx, q.SQL, q.Args...)  => 执行 sql
val := s.db.valCreator(tp, meta)  => 构造 valuer
err = val.SetColumns(rows)  => 填充数据
```





## Deleter

- Executor
- QueryBuilder

```go
type Deleter[T any] struct {
	builder
	tableName string
	where     []Predicate
	db        *DB
}
```



### Build() (*Query, error)

```go
```





## Inserter & Upsert

- Executor
- QueryBuilder

```go
type Inserter[T any] struct {
	builder
	values  []*T
	db      *DB
	columns []string
	upsert  *Upsert
}
```



### Build() (*Query, error)

```go
m, err := i.db.r.Get(i.values[0])  => 获取模型
i.sb.WriteString("INSERT INTO ")  => 构造固定部分
i.quote(m.TableName)
i.sb.WriteString("(")
fields := m.Fields   => 获取列名
if len(i.columns) != 0
for idx, fd := range fields  => 构造 insert 列名
i.sb.WriteString(") VALUES")
for vIdx, val := range i.values  => 遍历传入参数, 构造 insert 值
	-> refVal := i.db.valCreator(val, i.model)  => 构造 valuer
       for fIdx, field := range fields  => 遍历列名
			-> fdVal, err := refVal.Field(field.GoName)  => 获取列名值, 加入 inserter args, 替换 ?
if i.upsert != nil  => 支持 upsert 语句, 不同数据库格式不同
```





### Exec(ctx context.Context) Result

```go
q, err := i.Build()  => 构造 sql
res, err := i.db.db.ExecContext(ctx, q.SQL, q.Args...)  => 执行 sql
```





## Updater

- Executor
- QueryBuilder



### Build() (*Query, error)

```go
u.model, err = u.db.r.Get(&t)  => 获取模型
u.sb.WriteString("UPDATE ")  => 构造固定部分
u.quote(u.model.TableName)
u.sb.WriteString(" SET ")
val := u.db.valCreator(u.val, u.model)  => 构造 valuer
for i, a := range u.assigns  => 遍历所有要更改的字段
	-> switch assign := a.(type)  => 判断不同的类型
if len(u.where) > 0  => 构造 where
```





### Exec(ctx context.Context) Result

```go
q, err := u.Build()  => 构造 sql
res, err := u.db.db.ExecContext(ctx, q.SQL, q.Args...)  => 执行 sql
```















