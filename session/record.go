package session

import (
	"MiniORM/clause"
	"errors"
	"reflect"
)

// 在该文件下定义的都是链式调用的末端方法

// Insert 插入
func (s *Session) Insert(values ...interface{}) (int64, error) {
	// 每一列的记录
	recordValues := make([]interface{}, 0)
	if len(values) > 0 {

	}
	// 对于传入的所有记录
	for _, value := range values {
		// 取出表名
		table := s.Model(value).RefTable()
		// 生成insert子句 INSERT INTO table_name(col1, col2, col3, ...)
		// 这里是为了书写方便，实际只需要set一次
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
			// 将所有字段域展开交给scan，注意这里获取的是地址，所以实际dest会被注入数据
			// 注意这里需要结构体的字段和表格字段的顺序是一致的，否则会出现装配失败的情况
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

// Update 更新，可以接受map或者平铺开来的键值对
func (s *Session) Update(kv ...interface{}) (int64, error) {
	// 强转map
	m, ok := kv[0].(map[string]interface{})
	// 转失败的话他可能就按k,v,k,v的形式收集
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	// 构造update子句
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	// 构建sql
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	// 执行
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.refTable.Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

// Limit 链式调用limit
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// Where 链式调用where,desc是 ”name = ? "这种，vars是备填数据
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	vars = append(vars, desc)
	vars = append(vars, args...)
	s.clause.Set(clause.WHERE, vars...)
	return s
}

// OrderBy 链式调用orderby
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

func (s *Session) First(value interface{}) error {
	// 获取类型
	dest := reflect.Indirect(reflect.ValueOf(value))
	// 反射出类型列表
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	// 将值设回value
	dest.Set(destSlice.Index(0))
	return nil
}
