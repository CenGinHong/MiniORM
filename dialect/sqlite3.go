package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type sqlite3 struct {
}

func (s *sqlite3) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s （%s)", typ.Type().Name(), typ.Kind()))
}

// TableExistSQL 串sql,表格是否存在
func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", args
}

var _ Dialect = (*sqlite3)(nil)

// init 包加载时自动注册
func init() {
	// 将sqlite3注册
	RegisterDialect("sqlite3", &sqlite3{})
}

type User struct {
	Name string `miniorm:"PRIMARY KEY"`
	Age  int
}
