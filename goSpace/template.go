// TODO: Description of package.
package goSpace

import "reflect"

// Template is used for finding tuples.
// The template struct contains information if the number of fields of the
// template should match the number of fields of the tuple.
type Template struct {
	Fields []interface{} // Field of the template.
}

// CreateTemplate will create the template and return it with the fields
// specified by the user and set the matchNumberOfFields to true as default.
func CreateTemplate(fields []interface{}) Template {
	// Creates copy of fields
	tempfields := make([]interface{}, len(fields))
	copy(tempfields, fields)
	// Replace pointers with string from reflect.type value (used to match type)
	for i, value := range fields {
		// Check if value is a pointer
		if reflect.TypeOf(value).Kind() == reflect.Ptr {
			// Replace with typefield
			tempfields[i] = CreateTypeField(reflect.ValueOf(value).Elem().Type().String())
		}
	}
	template := Template{tempfields}
	return template
}

// length returns the amount of fields of the template.
func (temp *Template) length() int {
	return len(temp.Fields)
}

// getFieldAt will return the field at position i in fields of the template.
func (temp *Template) getFieldAt(i int) interface{} {
	return temp.Fields[i]
}
