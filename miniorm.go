package miniorm

import (
	"MiniORM/dialect"
	"MiniORM/log"
	"MiniORM/session"
	"database/sql"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver string, source string) (e *Engine, err error) {
	// 连接到数据库
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// ping一下去确保数据库连接存活
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	// 获取数据库dialect
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found, driver")
		return
	}
	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return e, nil
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// NewSession 新建一个会话
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

type TxFunc func(*session.Session) (interface{}, error)

func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err = s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			defer func() {
				// 若commit失败返回error，就需要再一次rollback
				if err != nil {
					_ = s.Rollback()
				}
			}()
			err = s.Commit()
		}
	}()
	return f(s)
}
