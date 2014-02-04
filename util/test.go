package util

import (
	"reflect"
)

// Panic catcher stolen from https://groups.google.com/forum/#!topic/golang-nuts/Hg_u6fdTx0I

// Panics returns true if function f panics with parameters p.
func Panics(f interface{}, p ...interface{}) bool {
	fv := reflect.ValueOf(f)
	ft := reflect.TypeOf(f)

	if ft.NumIn() != len(p) {
		panic("wrong argument count")
	}

	pv := make([]reflect.Value, len(p))
	for i, v := range p {
		if reflect.TypeOf(v) != ft.In(i) {
			panic("wrong argument type")
		}
		pv[i] = reflect.ValueOf(v)
	}

	return call(fv, pv)
}

func call(fv reflect.Value, pv []reflect.Value) (b bool) {
	defer func() {
		if err := recover(); err != nil {
			b = true
		}
	}()

	fv.Call(pv)
	return
}
