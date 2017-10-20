package shared

import (
	"fmt"
	"reflect"
	"strings"
)

// Intertemplate an interface for manipulating templates.
type Intertemplate interface {
	Length() int
	GetFieldAt(i int) interface{}
	NewTuple() Tuple
}

// Template structure used for matching against tuples.
// Templates is, in princple, a tuple with additional attributes used for pattern matching.
type Template struct {
	Fields []interface{}
}

// CreateTemplate creates a template from the variadic fields provided.
// CreateTemplate encapsulates the types of pointer values and are used to provide variable binding.
// Variable binding is used in Pattern matched values such that they can be writen back through the pointers.
func CreateTemplate(fields ...interface{}) Template {
	tempfields := make([]interface{}, len(fields))
	copy(tempfields, fields)

	// Replace pointers with reflect.Type value used to match type.
	for i, value := range fields {
		if reflect.TypeOf(value).Kind() == reflect.Ptr {
			// Encapsulate the parameter with a TypeField.
			tempfields[i] = CreateTypeField(reflect.ValueOf(value).Elem().Type())
		}
	}

	template := Template{tempfields}

	return template
}

// Length returns the amount of fields of the template.
func (tp *Template) Length() int {
	return len((*tp).Fields)
}

// GetFieldAt returns the i'th field of the template.
func (tp *Template) GetFieldAt(i int) interface{} {
	return (*tp).Fields[i]
}

// NewTuple returns a new tuple t from the template.
// NewTuple initializes all tuple fields in t with empty values depending on types in the template.
func (tp *Template) NewTuple() (t Tuple) {
	param := make([]interface{}, tp.Length())
	var element interface{}

	for i, _ := range param {
		field := tp.GetFieldAt(i)

		if field != nil {
			if reflect.TypeOf(field) == reflect.TypeOf(TypeField{}) {
				tf := reflect.ValueOf(field).Interface().(TypeField)
				rt := (tf.GetType()).(reflect.Type)
				element = reflect.New(rt).Elem().Interface()
			} else {
				ptf := reflect.TypeOf(field)
				element = reflect.New(ptf).Elem().Interface()
			}
		} else {
			element = field
		}

		param[i] = element
	}

	t = CreateTuple(param...)

	return t
}

// GetParenthesisType returns a pair of strings that encapsulates the template.
// GetParenthesisType is used in the String() method.
func (tp Template) GetParenthesisType() (string, string) {
	return "(", ")"
}

// GetDelimiter returns the delimiter used to seperated the template fields.
// GetParenthesisType is used in the String() method.
func (tp Template) GetDelimiter() string {
	return ", "
}

// String returns a print friendly representation of the template.
func (tp Template) String() string {
	lp, rp := tp.GetParenthesisType()

	delim := tp.GetDelimiter()

	strs := make([]string, tp.Length())

	for i, _ := range strs {
		field := tp.GetFieldAt(i)
		if field != nil {
			if reflect.TypeOf(field).Kind() == reflect.String {
				strs[i] = fmt.Sprintf("%s%s%s", "\"", field, "\"")
			} else {
				strs[i] = fmt.Sprintf("%v", field)
			}
		} else {
			strs[i] = "nil"
		}
	}

	return fmt.Sprintf("%s%s%s", lp, strings.Join(strs, delim), rp)
}
