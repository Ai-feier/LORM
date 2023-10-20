//go:build v6
package orm

import (
	"LORM/v6/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

// Registry 元数据注册中心的抽象
type Registry interface {
	// Get 查找元数据
	Get(val any) (*Model, error)
	// Register 注册一个模型
	Register(val any, opts ...ModelOpt) (*Model, error)
}
// 这种包变量对测试不友好，缺乏隔离
// var defaultRegistry = &registry{
// 	models: make(map[reflect.Type]*Model, 16),
// }



type registry struct {
	// models key 是类型名
	// 这种定义方式是不行的
	// 1. 类型名冲突，例如都是 User，但是一个映射过去 buyer_t
	// 一个映射过去 seller_t
	// 2. 并发不安全
	// models map[string]*Model

	// lock sync.RWMutex
	// models map[reflect.Type]*Model
	models sync.Map
}

func NewRegistry() Registry {
	return &registry{}
}

func (r *registry) Register(val any, opts ...ModelOpt) (*Model, error) {
	// 解析 val 
	m, err := r.parseModel(val)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		err := opt(m)
		if err != nil {
			return nil, err
		}
	}
	typ := reflect.TypeOf(val)
	r.models.Store(typ, m)
	return m, nil
}

func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if !ok {
		// 加载
		var err error
		if m, err = r.parseModel(val); err != nil {
			return nil, err
		}
	}
	r.models.Store(typ, m)
	return m.(*Model), nil
}

// 接受一级结构体指针
// parseModel 支持从标签中提取自定义设置
// 标签形式 orm:"key1=value1,key2=value2"
func (r *registry) parseModel(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Pointer ||
		typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()


	// 获取字段数量
	numField := typ.NumField()
	fds := make(map[string]*field, numField)
	for i := 0; i < numField; i++ {
		fdTyp := typ.Field(i)
		tags, err := r.parseTag(fdTyp.Tag)
		if err != nil {
			return nil, err
		}
		colName := tags[tagKeyColumn]
		if colName == "" {
			colName = underscoreName(fdTyp.Name)
		}
		
		fds[fdTyp.Name] = &field{
			colName: colName,
		}
	}
	
	var tableName string 
	if tn, ok := val.(TableName); ok {
		tableName = tn.TableName()
	}
	if tableName == "" {
		tableName = underscoreName(typ.Name())
	}
	
	return &Model{
		tableName: tableName, 
		fieldMap:  fds,
	}, nil

}

func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag := tag.Get("orm")
	if ormTag == "" {
		return map[string]string{}, nil
	}
	res := make(map[string]string, 1)
	
	// 切割字符串
	pairs := strings.Split(ormTag, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return nil, errs.NewErrInvalidTagContent(pair)
		}
		res[kv[0]] = kv[1]
	}
	return res, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}
	}
	return string(buf)
} 
