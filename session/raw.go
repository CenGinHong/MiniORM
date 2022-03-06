package session

import (
	"MiniORM/dialect"
	"MiniORM/log"
	"MiniORM/schema"
	"database/sql"
	"strings"
)

// Session 会话
type Session struct {
	db       *sql.DB         // 数据库连接指针
	sql      strings.Builder // sql语句拼接
	sqlVars  []interface{}   // 占位符拼接
	refTable *schema.Schema
	dialect  dialect.Dialect
}

// New 新创建一个session
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect}
}

// Clear 清空sql语句
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return result, nil
}

// QueryRow 查询一行数据
func (s *Session) QueryRow() *sql.Row {
	// 清空所有sql拼接句段
	defer s.Clear()
	// 打印sql语句
	log.Info(s.sql.String(), s.sqlVars)
	// 调用原生DB查询
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRows 查询列表
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return rows, err
}
