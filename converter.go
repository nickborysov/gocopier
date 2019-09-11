package copier

import (
	"reflect"
)

var interfaceType = reflect.ValueOf(new(interface{})).Elem().Type()

// Converter represents struct, which used by Copier for convert one type to another
type Converter struct {
	Src     interface{}
	Dst     interface{}
	Convert func(src interface{}) (interface{}, error)
}

// SrcType returns reflect.Type of Src or interfaceType if src is nil
func (p *Converter) SrcType() reflect.Type {
	if p.Src == nil {
		return interfaceType
	}
	return reflect.ValueOf(p.Src).Type()
}

// DstType returns reflect.Type of Dst or interfaceType if dst is nil
func (p *Converter) DstType() reflect.Type {
	if p.Dst == nil {
		return interfaceType
	}
	return reflect.ValueOf(p.Dst).Type()
}
