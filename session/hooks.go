package session

import (
	"MiniORM/log"
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

// CallMethod 调用一个已经注册hook
func (s *Session) CallMethod(method string, value interface{}) {
	// 获取结构体的hooks
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	// 如果value存在，就获取当前操作的对象，例如insert操作的对象
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	// 将session作为入参
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		// 调用，并处理error
		if v := fm.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return
}
