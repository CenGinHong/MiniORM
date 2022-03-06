package schema

import (
	"MiniORM/dialect"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

type User struct {
	Name string `miniorm:"PRIMARY KEY"`
	Age  int
}

var TestDial, _ = dialect.GetDialect("sqlite3")

func TestParse(t *testing.T) {
	schema := Parse(&User{}, TestDial)
	if schema.Name != "User" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse User Struct")
	}
	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}
