package h2g

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

type module struct {
	servers map[string]reflect.Method
	rcvr reflect.Value
	typ reflect.Type
}

func newModule(v interface{}) *module {
	m := &module{
		servers: make(map[string]reflect.Method),
		rcvr: reflect.ValueOf(v),
		typ: reflect.TypeOf(v),
	}
	for i := 0; i < m.typ.NumMethod(); i ++ {
		method := m.typ.Method(i)
		m.servers[method.Name] = method
	}
	return m
}

func (m *module) call(name string, buf []byte) ([]byte, error) {
	method, ok := m.servers[name]
	if !ok {
		return nil, fmt.Errorf("method[%s] is not exists", name)
	}
	vType := method.Type.In(2).Elem()
	v := reflect.New(vType).Interface()
	err := json.Unmarshal(buf, &v)
	if err != nil {
		return nil, err
	}
	response := method.Func.Call([]reflect.Value{
		m.rcvr,
		reflect.ValueOf(context.Background()),
		reflect.ValueOf(v),
	})
	if len(response) != 2 {
		return nil, fmt.Errorf("respone length is error")
	}
	r1 := response[1]
	if !r1.IsNil() {
		err, ok = r1.Interface().(error)
		if !ok {
			return nil, fmt.Errorf("response second value's type is not error")
		}
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(response[0].Interface())
}
