package session

import (
	"MiniORM/clause"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	// 每一列的记录
	recordValues := make([]interface{}, 0)
	// 对于传入的所有记录
	for _, value := range values {
		// 取出表名
		table := s.Model(value).RefTable()
		// 生成insert子句 INSERT INTO table_name(col1, col2, col3, ...)
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		// 将要插入的data收集
		recordValues = append(recordValues, table.RecordValues(value))
	}
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Find 传入一个切片，并将查询结果置入切片
func (s *Session) Find(values interface{}) error {
	// 通过指针获取值
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	// 获取切片单个类型的值类型，注意两种Elem的使用场景
	destType := destSlice.Type().Elem()
	// 获取表格
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()
	// 构建子句
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	// 拼凑sql
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	// 执行多行查询
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}
	// 对于每一行查出来的数据
	for rows.Next() {
		// 反射一个相同类型的值出来
		dest := reflect.New(destType).Elem()
		// 收集列表
		var values []interface{}
		// 对于每一个字段
		for _, name := range table.FieldNames {
			// 将所有字段域铺平，注意这里获取的是地址，所以实际dest会被注入数据
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		// 进行scan，将结果映射回字段域
		if err = rows.Scan(values...); err != nil {
			return err
		}
		// 将结果收集
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
