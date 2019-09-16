package copier

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
)

var (
	// ErrInvalidDestination represents error invalid destination
	ErrInvalidDestination = errors.New("invalid destination")
	// ErrInvalidSource represents error invalid source
	ErrInvalidSource = errors.New("invalid source")
	// ErrDifferentTypes represents error src and dst fields has different types
	ErrDifferentTypes = errors.New("src and dst fields has different types")
	// ErrCannotSetValue represents error can not set value
	ErrCannotSetValue = errors.New("can not set value")
)

// Copier represents struct of Copier
type Copier struct {
	Converters []Converter
	Logger     *log.Logger
}

// New creates new Copier
func New() *Copier {
	defaultLogger := log.New(os.Stderr, "", log.LstdFlags)
	return &Copier{
		Logger: defaultLogger,
	}
}

// SetConverters set converters to copier
func (c *Copier) SetConverters(cc []Converter) {
	c.Converters = cc
}

// Copy create new copier, set converters and make copy value from source to destination
func Copy(dst, src interface{}, cc ...Converter) error {
	copier := New()
	copier.SetConverters(cc)
	return copier.Copy(dst, src)
}

// Copy make copy value from source to destination
func (c *Copier) Copy(dst, src interface{}) error {
	d := reflect.ValueOf(dst)
	if d.Kind() != reflect.Ptr {
		return ErrInvalidDestination
	}
	s := reflect.ValueOf(src)
	if s.Kind() != reflect.Ptr {
		return ErrInvalidSource
	}
	return c.copyInterface(d, s.Elem())
}

func (c *Copier) copyInterface(dst, src reflect.Value) error {
	if src.Kind() == reflect.Ptr && src.IsNil() {
		return nil
	}
	for _, pattern := range c.Converters {
		dstElem := dst
		if dst.Kind() == reflect.Ptr {
			dstElem = dst.Elem()
		}
		if pattern.SrcType() == src.Type() && pattern.DstType() == dstElem.Type() {
			res, err := pattern.Convert(src.Interface())
			if err != nil {
				return err
			}
			dstElem.Set(reflect.ValueOf(res))
			return nil
		}
	}
	switch src.Kind() {
	case reflect.Slice, reflect.Array:
		return c.copySliceArray(dst, src)
	case reflect.Map:
		return c.copyMap(dst, src)
	case reflect.Struct:
		return c.copyStruct(dst, src)
	case reflect.Ptr:
		return c.copyPtr(dst, src)
	default:
		return c.copyElement(dst, src)
	}
}

func (c *Copier) copyPtr(dst, src reflect.Value) error {
	switch {
	case dst.Kind() != reflect.Ptr && !dst.CanAddr():
		return fmt.Errorf("%s: expected %s, actual %s", ErrDifferentTypes, src.Kind(), dst.Kind())
	case dst.Kind() != reflect.Ptr:
		return c.copyInterface(dst.Addr(), src.Elem())
	case dst.IsNil():
		newElem := reflect.New(reflect.TypeOf(dst.Interface()).Elem())
		dst.Set(newElem)
		return c.copyInterface(newElem, src.Elem())
	default:
		return c.copyInterface(dst.Elem(), src.Elem())
	}
}

func (c *Copier) copyMap(dst, src reflect.Value) error {
	dstElem := dst
	if dst.Kind() == reflect.Ptr {
		dstElem = dst.Elem()
	}
	if src.Kind() != dstElem.Kind() {
		return fmt.Errorf("%s: expected %s, actual %s", ErrDifferentTypes, src.Kind(), dstElem.Kind())
	}
	dstValueType := reflect.TypeOf(dstElem.Interface()).Elem()
	srcValueType := reflect.TypeOf(src.Interface()).Elem()
	srcKeyType := reflect.TypeOf(src.Interface()).Key()
	dstKeyType := reflect.TypeOf(dstElem.Interface()).Key()
	if dstElem.IsNil() {
		dstElem.Set(reflect.MakeMapWithSize(dstElem.Type(), src.Len()))
	}
	for _, key := range src.MapKeys() {
		dstKey := key
		if srcKeyType.Kind() != dstKeyType.Kind() {
			newKey := reflect.New(dstKeyType)
			err := c.copyInterface(newKey, key)
			if err != nil {
				return err
			}
			dstKey = newKey.Elem()
		}
		dstValue := src.MapIndex(key)
		if srcValueType.Kind() != dstValueType.Kind() {
			newValue := reflect.New(dstValueType)
			err := c.copyInterface(newValue, dstValue)
			if err != nil {
				return err
			}
			dstValue = newValue.Elem()
		}
		dstElem.SetMapIndex(dstKey, dstValue)
	}
	return nil
}

func (c *Copier) copyStruct(dst, src reflect.Value) error {
	for i := 0; i < src.NumField(); i++ {
		field := src.Field(i)
		if dst.Kind() == reflect.Ptr {
			dst = dst.Elem()
		}
		structField := reflect.TypeOf(src.Interface()).Field(i)
		dstValue := dst.FieldByName(structField.Name)
		if dstValue == (reflect.Value{}) {
			c.Logger.Printf("Field not found: %s", structField.Name)
			return nil
		}
		err := c.copyInterface(dstValue, field)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Copier) copySliceArray(dst, src reflect.Value) error {
	if src.Len() == 0 {
		return nil
	}
	dstElem := dst
	if dst.Kind() == reflect.Ptr {
		dstElem = dst.Elem()
	}
	var slice reflect.Value
	switch {
	case dstElem.Kind() == reflect.Slice && src.Len() > dstElem.Len():
		slice = reflect.MakeSlice(dstElem.Type(), src.Len(), src.Len())
	case dstElem.Kind() != reflect.Slice && dstElem.Kind() != reflect.Array:
		return fmt.Errorf("%s: expected %s, actual %s", ErrDifferentTypes, src.Kind(), dstElem.Kind())
	default:
		slice = dstElem
	}
	for i := 0; i < src.Len(); i++ {
		err := c.copyInterface(slice.Index(i), src.Index(i))
		if err != nil {
			return err
		}
	}
	if dstElem.Kind() == reflect.Slice {
		dstElem.Set(slice)
	}
	return nil
}

func (c *Copier) copyElement(dst, src reflect.Value) error {
	if dst.Kind() == reflect.Ptr {
		if dst.IsNil() {
			newElem := reflect.New(reflect.TypeOf(dst.Interface()).Elem())
			dst.Set(newElem)
		}
		return c.copyInterface(dst.Elem(), src)
	}
	if src.Kind() != dst.Kind() {
		return fmt.Errorf("%s: expected %s, actual %s", ErrDifferentTypes, src.Kind(), dst.Kind())
	}
	if src.Type() != dst.Type() {
		return fmt.Errorf("%s: expected %s, actual %s", ErrDifferentTypes, src.Type(), dst.Type())
	}
	if !dst.CanSet() {
		return fmt.Errorf("%s: %v to %v", ErrCannotSetValue, src, dst)
	}
	dst.Set(src)
	return nil
}
