package dialect

import "reflect"

var dialectMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(typ reflect.Value) string                    // 将Go的基本类型转换成该数据库的数据类型
	TableExistSQL(tableName string) (string, []interface{}) // 表是否已经存在
}

func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectMap[name]
	return
}
