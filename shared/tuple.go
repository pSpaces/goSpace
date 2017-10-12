package shared

import (
	"reflect"
)

// Tuple contains a set of fields, where fields can be any primitive or type.
// A tuple is used to store information which is placed in a tuple space.
type Tuple struct {
	Fields []interface{} // Field of the tuple.
}

// CreateTuple will create the tuple and return it with the fields specified
// by the user.
func CreateTuple(fields []interface{}) Tuple {
	tuple := Tuple{fields}
	return tuple
}

// Length returns the amount of fields of the tuple.
func (t *Tuple) Length() int {
	return len(t.Fields)
}

// GetFieldAt will return the field at position i in fields of the tuple.
func (t *Tuple) GetFieldAt(i int) interface{} {
	return t.Fields[i]
}

// SetFieldAt will set the field at position i in fields of the tuple to the
// value of val specified by the user.
func (t *Tuple) SetFieldAt(i int, val interface{}) {
	t.Fields[i] = val
}

// match will return the boolean value according to if the template temp match
// the tuple or not.
func (t *Tuple) Match(temp Template) bool {
	if t.Length() != temp.length() {
		return false
	}

	// Run through corresponding fields of tuple and template to see if they are
	// matching.
	for i := 0; i < temp.length(); i++ {
		// Check if the field of the template is a formal or actual field.
		// Extract corresponding fields from tuple and template.
		tupleField := t.GetFieldAt(i)
		tempField := temp.getFieldAt(i)
		// Check if tempField is a TypeField
		if reflect.TypeOf(tempField) == reflect.TypeOf(TypeField{}) {
			// Check if the type of tupleField is the same type as specified
			// in tempField
			if tempField.(TypeField).getType() == reflect.TypeOf(tupleField).String() {
				continue
			} else {
				return false
			}
		} else if !reflect.DeepEqual(tupleField, tempField) { // Check if tupleField and tempField are equal.
			return false
		}
	}
	return true
}
