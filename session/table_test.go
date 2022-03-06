package session

import (
	"MiniORM/dialect"
	"MiniORM/log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

type User struct {
	Name string `miniorm:"PRIMARY KEY"`
	Age  int
}

func TestSession_CreateTable(t *testing.T) {
	s := NewSession().Model(&User{})
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table user")
	}
}

func TestMain(m *testing.M) {
	TestDB, _ := sql.Open("sqlite3", "../gee.db")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	TestDB, err := sql.Open("sqlite3", "../gee.db")
	if err != nil {
		log.Error(err)
	}
	TestDial, _ := dialect.GetDialect("sqlite3")
	return New(TestDB, TestDial)
}
