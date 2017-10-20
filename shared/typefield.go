package shared

import (
	"reflect"
)

// typeMap maintains a map of registered types.
var typeMap = make(map[string]reflect.Type)

// TypeField encapsulate a type.
type TypeField struct {
	TypeStr string
}

// CreateTypeField creates an encapsulation of a type.
func CreateTypeField(t reflect.Type) TypeField {
	ts := t.String()
	tf := TypeField{ts}

	_, exists := typeMap[ts]
	if !exists {
		typeMap[ts] = t
	}

	return tf
}

// GetType returns a type associated to this Typefield.
func (tf TypeField) GetType() reflect.Type {
	return typeMap[tf.TypeStr]
}

// String returns the type string of this TypeField.
func (tf TypeField) String() string {
	return tf.TypeStr
}
