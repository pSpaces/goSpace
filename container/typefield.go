package container

import (
	"math"
	"reflect"
	"sync"

	"github.com/pspaces/gospace/function"
)

// TypeRegistry represents a registry of types.
type TypeRegistry *sync.Map

// gtr maintains a global type registry.
var gtr = new(sync.Map)

// TypeField encapsulate a type.
type TypeField struct {
	TypeStr string `bson:"type" json:"type" xml:"type"`
}

// CreateTypeField creates an encapsulation of a type.
func CreateTypeField(t interface{}) TypeField {
	ts := reflect.ValueOf(t).Type().String()
	tf := TypeField{ts}

	_, exists := gtr.Load(ts)
	if !exists {
		registerType := func(ts string, t reflect.Type) {
			gtr.Store(ts, t)
		}

		recursiveTypeRegister(math.MaxUint16, t, registerType)
	}

	return tf
}

// recursiveTypeRegister traverses any arbitrary structure and registers all of its subtypes.
// recursiveTypeRegister will traverse any pointers it will encounter and will not register the pointer type itself.
// recursiveTypeRegister will not traverse unsafe pointers or invalid types.
// recursiveTypeRegister does not use any memoization and will terminate at depth of 2^16.
func recursiveTypeRegister(rd uint64, value interface{}, fun interface{}) {
	reg := fun.(func(string, reflect.Type))

	var ts string
	switch val := reflect.ValueOf(value); val.Kind() {
	case reflect.Invalid, reflect.UnsafePointer:
	case reflect.Func:
		ts = val.Type().String()
		reg(ts, val.Type())
	case reflect.Array, reflect.Slice:
		if rd >= 0 {
			len := val.Len()
			for i := 0; i < len; i++ {
				recursiveTypeRegister(rd-1, val.Index(i).Interface(), fun)
			}
		}

		ts = val.Type().String()
		reg(ts, val.Type())
	case reflect.Map:
		if rd >= 0 {
			keys := val.MapKeys()

			for _, key := range keys {
				recursiveTypeRegister(rd-1, key, fun)
				index := val.MapIndex(key)
				recursiveTypeRegister(rd-1, index, fun)
			}
		}

		ts = val.Type().String()
		reg(ts, val.Type())
	case reflect.Struct:
		if rd >= 0 {
			cnt := val.NumField()

			for i := 0; i < cnt; i++ {
				field := val.Field(i)
				recursiveTypeRegister(rd-1, field.Interface(), fun)
			}
		}

		ts = val.Type().String()
		reg(ts, val.Type())
	case reflect.Ptr:
		if function.IsFunc(value) {
			ts = reflect.TypeOf(val).String()
			reg(ts, val.Type())
		} else {
			rval := val.Elem()
			recursiveTypeRegister(rd-1, rval, fun)
		}
	default:
		ts = val.Type().String()
		reg(ts, val.Type())
	}

	return
}

// Equal returns true if both type field a and b are quivalent, and false otherwise.
func (tf TypeField) Equal(tfo TypeField) (e bool) {
	rta := tf.GetType()
	rtb := tfo.GetType()
	e = rta == rtb
	return e
}

// GetType returns a type associated to this Typefield.
func (tf TypeField) GetType() (t reflect.Type) {
	ti, _ := gtr.Load(tf.TypeStr)
	t = ti.(reflect.Type)
	return t
}

// String returns the type string of this TypeField.
func (tf TypeField) String() string {
	return tf.TypeStr
}
