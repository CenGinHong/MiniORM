package clause

import "strings"

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

// Set 调用generator, 生成对应子句
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	// 找到对应的generator
	sql, vars := generators[name](vars...)
	// 产生子句
	c.sql[name] = sql
	// values和_where会有vars值
	c.sqlVars[name] = vars
}

// Build 根据传入Type的顺序构造出最终的sql语句
func (c Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	// 拼接所有子句，并且把嵌入参数返回
	for _, order := range orders {
		// 仅当子句存在，即ok是再进行拼接
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
