package tuplespace

import "reflect"

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
func (t *Tuple) match(temp Template) bool {
	// Check if the number of fields in the tuple and template should match,
	// i.e. arbitrary number of fields in the tuple.
	// TODO: Get rid of getMatchNumberOfFields
	/*
		if temp.getMatchNumberOfFields() {
			// Check that the tuple and template have matching lengths.
			if t.Length() != temp.length() {
				return false
			}
		} else {
			// If the length of the tuple is less than the length of the template
			// they will never match.
			if t.Length() < temp.length() {
				return false
			}
		}
	*/

	if t.Length() != temp.length() {
		return false
	}

	// Check if fields of the template match the fields of the tuple.
	//return t.matchFieldsOf(temp)

	// Run through corresponding fields of tuple and template to see if they are
	// matching.
	for i := 0; i < temp.length(); i++ {
		// Check if the field of the template is a formal or actual field.
		// Extract corresponding fields from tuple and template.
		tupleField := t.GetFieldAt(i)
		tempField := temp.getFieldAt(i)
		// Should tempField be a type, check if tupleField is of the same type.
		if reflect.TypeOf(tempField) == reflect.TypeOf(TypeField{}) {
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

// matchFieldsOf will match the corresponding fields of the template and tuple
// to see if they match. If all fields match true is returned otherwise false
// is returned.
// NOTE: This method assumes that the length of the template is less than or
// equal to the length of the tuple.

/*
func (t *Tuple) matchFieldsOf(temp Template) bool {
>>>>>>> .r48
	// Run through corresponding fields of tuple and template to see if they are
	// matching.
	for i := 0; i < temp.length(); i++ {
		// Check if the field of the template is a formal or actual field.
		// Extract corresponding fields from tuple and template.
		tupleField := t.GetFieldAt(i)
		tempField := temp.getFieldAt(i)
		// Should tempField be a type, check if tupleField is of the same type.
		if reflect.TypeOf(tempField) == reflect.TypeOf(TypeField{}) {
			//magic
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
*/
