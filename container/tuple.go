package container

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pspaces/gospace/function"
)

// Tuple contains a set of fields, where fields can be any primitive or type.
// A tuple is used to store information which is placed in a tuple space.
type Tuple struct {
	Flds []interface{} `bson:"fields" json:"fields" xml:"fields"`
}

// NewTuple create a tuple according to the values in the fields.
func NewTuple(fields ...interface{}) Tuple {
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

	tp = NewTuple(fields...)

	return tp
}

// Length returns the amount of fields of the tuple.
func (t *Tuple) Length() (sz int) {
	sz = -1

	if t != nil {
		sz = len((*t).Flds)
	}

	return sz
}

// Fields returns the fields of the tuple.
func (t *Tuple) Fields() (flds []interface{}) {
	if t != nil {
		flds = (*t).Flds
	}

	return flds
}

// GetFieldAt returns the i'th field of the tuple.
func (t *Tuple) GetFieldAt(i int) (fld interface{}) {
	if t != nil && i >= 0 && i < t.Length() {
		fld = (*t).Flds[i]
	}

	return fld
}

// SetFieldAt sets the i'th field of the tuple to the value of val.
func (t *Tuple) SetFieldAt(i int, val interface{}) (b bool) {
	if t != nil && i >= 0 && i < t.Length() {
		(*t).Flds[i] = val
		b = true
	}

	return b
}

// Apply iterates through the tuple t and applies the function fun to each field.
// Apply returns true function fun could be applied to all the fields, and false otherwise.
func (t *Tuple) Apply(fun func(field interface{}) interface{}) (b bool) {
	b = false

	if t != nil {
		b = true
		for i := 0; i < t.Length(); i++ {
			t.SetFieldAt(i, fun(t.GetFieldAt(i)))
		}
	}

	return b
}

// Tuple returns a new tuple.
func (t *Tuple) Tuple() (tn Tuple) {
	tn = NewTuple((*t).Flds[:]...)
	return tn
}

// Match pattern matches the tuple against the template tp.
// Match discriminates between encapsulated formal fields and actual fields.
// Match returns true if the template matches the tuple, and false otherwise.
func (t *Tuple) Match(tp Template) (b bool) {
	b = t != nil && t.Length() == (&tp).Length()

	// Run through corresponding fields of tuple and template to see if they are
	// matching.
	for i := 0; i < tp.Length() && b; i++ {
		tf := (*t).GetFieldAt(i)
		tpf := tp.GetFieldAt(i)
		// Check if the field of the template is an encapsulated formal or actual field.
		if reflect.TypeOf(tpf) == reflect.TypeOf(TypeField{}) {
			b = reflect.TypeOf(tf) == tpf.(TypeField).GetType()
		} else if function.IsFunc(tf) && function.IsFunc(tpf) {
			// We can do better. Functions are not being moved or rewritten while one is executing, no?
			// If a function has a static address, and one has two functions with the same static address,
			// then the functionality they provide must be equal. Then we shall match and accept any consequences of
			// inlining and using addresses of anonymous functions.
			b = (function.Name(tf) == function.Name(tpf)) && (function.Signature(tf) == function.Signature(tpf))
		} else {
			b = reflect.DeepEqual(tf, tpf)
		}
	}

	return b
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

	for i := range strs {
		field := t.GetFieldAt(i)
		if field != nil {
			if reflect.TypeOf(field).Kind() == reflect.String {
				strs[i] = fmt.Sprintf("%s%s%s", "\"", field, "\"")
			} else if function.IsFunc(field) {
				strs[i] = fmt.Sprintf("%s %s", function.Name(field), function.Signature(field))
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
func (t *Tuple) WriteToVariables(params ...interface{}) (b bool) {
	b = t != nil

	if b {
		for i, param := range params {
			if reflect.TypeOf(param).Kind() == reflect.Ptr {
				value := reflect.ValueOf(param).Elem()
				if value.CanSet() {
					value.Set(reflect.ValueOf((*t).GetFieldAt(i)))
				}
			}
		}
	}

	return b
}
