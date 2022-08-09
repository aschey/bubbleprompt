package main

import (
	"sync"

	"github.com/dop251/goja"
)

type vm struct {
	mutex   sync.Mutex
	runtime *goja.Runtime
}

func newVm() *vm {
	return &vm{runtime: goja.New()}
}

func (v *vm) GlobalObject() *goja.Object {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	return v.runtime.GlobalObject()
}

func (v *vm) ToObject(val goja.Value) *goja.Object {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	if val == nil {
		return nil
	}
	return val.ToObject(v.runtime)
}

func (v *vm) RunString(str string) (goja.Value, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	return v.runtime.RunString(str)
}
