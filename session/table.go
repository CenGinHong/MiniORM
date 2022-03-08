package session

import (
	"MiniORM/log"
	"MiniORM/schema"
	"fmt"
	"reflect"
	"strings"
)

// Model 将s的refTable进行赋值
func (s *Session) Model(value interface{}) *Session {
	// 这里的value是指针
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		// 解析出schema
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

// CreateTable 创建表格
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	// 组装字段
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	// 创建字段
	if _, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec(); err != nil {
		return err
	}
	return nil
}

// DropTable 销毁table
func (s *Session) DropTable() error {
	//if _, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec(); err != nil {
	//	return err
	//}
	//return nil
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	// 返回查询表格是否存在sql语句
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	// 执行
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
