package miniorm

import (
	"MiniORM/session"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "gee.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

type User struct {
	Name string `miniorm:"PRIMARY KEY"`
	Age  int
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}

func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	// 这个s只是用来dropTable，没有使用tx
	s := engine.NewSession()
	err := s.Model(&User{}).DropTable()
	if err != nil {
		t.Fatal(err)
	}
	_, err = engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		// 注意这里执行一定要是用参数里的s进行，这个s是在Transaction 里新建后然后被begin后赋值了tx的
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		// 返回一个错误，引起rollback
		return nil, errors.New("error")
	})
	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	e := OpenDB(t)
	defer e.Close()
	s := e.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := e.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return
	})
	u := &User{}
	_ = s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}
