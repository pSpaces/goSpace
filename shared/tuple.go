package shared

import (
	"fmt"
	"reflect"
	"strings"
)

// Intertuple defines an interface for manipulating tuples.
type Intertuple interface {
	Length() int
	GetFieldAt(i int) interface{}
	SetFieldAt(i int, val interface{})
}

// Tuple contains a set of fields, where fields can be any primitive or type.
// A tuple is used to store information which is placed in a tuple space.
type Tuple struct {
	Fields []interface{} // Field of the tuple.
}

// CreateTuple create a tuple according to the values in the fields.
func CreateTuple(fields ...interface{}) Tuple {
	tf := make([]interface{}, len(fields))
	copy(tf, fields)
	tuple := Tuple{tf}
	return tuple
}

// CreateIntrinsicTuple creates a intrinsic tuple belonging to template t.
// CreateIntrinsicTuple initializes according to the template fields in tp by reading of concrete values.
// If CreateIntrinsicTuple finds a field containing a pointer, this field will be dereferenced and the value read.
func CreateIntrinsicTuple(t ...interface{}) (tp Tuple) {
	fields := make([]interface{}, len(t))

	for i, value := range t {
		if reflect.TypeOf(value).Kind() == reflect.Ptr {
			fields[i] = (reflect.ValueOf(value).Elem().Interface()).(interface{})
		} else {
			fields[i] = value
		}
	}

	tp = CreateTuple(fields...)

	return tp
}

// Length returns the amount of fields of the tuple.
func (t *Tuple) Length() int {
	return len((*t).Fields)
}

// GetFieldAt returns the i'th field of the tuple.
func (t *Tuple) GetFieldAt(i int) interface{} {
	return (*t).Fields[i]
}

// SetFieldAt sets the i'th field of the tuple to the value of val.
func (t *Tuple) SetFieldAt(i int, val interface{}) {
	(*t).Fields[i] = val
}

// Match pattern matches the tuple against the template tp.
// Match discriminates between encapsulated formal fields and actual fields.
// Match returns true if the template matches the tuple, and false otherwise.
func (t *Tuple) Match(tp Template) bool {
	if (*t).Length() != tp.Length() {
		return false
	} else if (*t).Length() == 0 && tp.Length() == 0 {
		return true
	}

	// Run through corresponding fields of tuple and template to see if they are
	// matching.
	for i := 0; i < tp.Length(); i++ {
		tf := (*t).GetFieldAt(i)
		tpf := tp.GetFieldAt(i)
		// Check if the field of the template is an encapsulated formal or actual field.
		if reflect.TypeOf(tpf) == reflect.TypeOf(TypeField{}) {
			if reflect.TypeOf(tf) != tpf.(TypeField).GetType() {
				return false
			}
		} else if !reflect.DeepEqual(tf, tpf) {
			return false
		}
	}

	return true
}

// GetParenthesisType returns a pair of strings that encapsulates the tuple.
// GetParenthesisType is used in the String() method.
func (t Tuple) GetParenthesisType() (string, string) {
	return "(", ")"
}

// GetDelimiter returns the delimiter used to seperated the tuple fields.
// GetParenthesisType is used in the String() method.
func (t Tuple) GetDelimiter() string {
	return ", "
}

// String returns a print friendly representation of the tuple.
func (t Tuple) String() string {
	ld, rd := t.GetParenthesisType()

	delim := t.GetDelimiter()

	strs := make([]string, t.Length())

	for i, _ := range strs {
		field := t.GetFieldAt(i)
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

	return fmt.Sprintf("%s%s%s", ld, strings.Join(strs, delim), rd)
}

// WriteToVariables will overwrite the values pointed to by pointers with
// the values contained in the tuple.
// WriteToVariables will ignore unaddressable pointers.
// TODO: There should be placed a lock around the variables that are being
// changed, to ensure that mix of two tuple are written to the variables.
func (t *Tuple) WriteToVariables(params ...interface{}) {
	for i, param := range params {
		if reflect.TypeOf(param).Kind() == reflect.Ptr {
			value := reflect.ValueOf(param).Elem()
			if value.CanSet() {
				value.Set(reflect.ValueOf((*t).GetFieldAt(i)))
			}
		}
	}
}
