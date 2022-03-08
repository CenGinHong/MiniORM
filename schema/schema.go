package schema

import (
	"MiniORM/dialect"
	"go/ast"
	"reflect"
)

// Field 字段域
type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema Schema结构
type Schema struct {
	Model      interface{}       // 映射结构
	Name       string            // 表名
	Fields     []*Field          // 字段结构
	FieldNames []string          // 所有字段名
	fieldMap   map[string]*Field // 字段名与字段结构的映射
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// 获取反射
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}
	// 遍历所有字段
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		// 该字段不能是嵌入或者未导出字段
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				// 获取指针指向的实例
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))), // TypeOf返回变量类型，Indirect能返回了指针所指向的变量的类型
			}
			// 遍历字段标签
			if v, ok := p.Tag.Lookup("miniorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

// RecordValues 将结构体的字段铺平，类似 (A1, A2, A3, ...)
func (s *Schema) RecordValues(dest interface{}) []interface{} {
	// 反射获取值
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	// 字段域值
	var fieldValues []interface{}
	// 对于每一个字段的值
	for _, field := range s.Fields {
		// 将反射结构体的字段值取出来
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
