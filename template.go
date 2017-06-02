// TODO: Description of package.
package tuplespace

import "reflect"

// Template is used for finding tuples.
// The template struct contains information if the number of fields of the
// template should match the number of fields of the tuple.
// TODO: Get rid of getMatchNumberOfFields
type Template struct {
	MatchNumberOfFields bool          // Boolean to specify if number of fields should match.
	Fields              []interface{} // Field of the template.
}

// CreateTemplate will create the template and return it with the fields
// specified by the user and set the matchNumberOfFields to true as default.
func CreateTemplate(fields []interface{}) Template {
	//creates copy of fields
	tempfields := make([]interface{}, len(fields))
	copy(tempfields, fields)
	//replace pointers with string from reflect.type value (used to match type)
	for i, value := range fields {
		//if value is a pointer
		if reflect.TypeOf(value).Kind() == reflect.Ptr {
			//replace with typefield
			tempfields[i] = CreateTypeField(reflect.ValueOf(value).Elem().Type().String())
		}
	}
	template := Template{true, tempfields}
	return template
}

// DontMatchNumberOfFields changes the value of matchNumberOfFields to false.
func (temp *Template) DontMatchNumberOfFields() {
	temp.MatchNumberOfFields = false
}

// DoMatchNumberOfFields changes the value of matchNumberOfFields to true.
func (temp *Template) DoMatchNumberOfFields() {
	temp.MatchNumberOfFields = true
}

// length returns the amount of fields of the template.
func (temp *Template) length() int {
	return len(temp.Fields)
}

// getFieldAt will return the field at position i in fields of the template.
func (temp *Template) getFieldAt(i int) interface{} {
	return temp.Fields[i]
}

// getMatchNumberOfFields returns the value of matchNumberOfFields.
func (temp *Template) getMatchNumberOfFields() bool {
	return temp.MatchNumberOfFields
}
