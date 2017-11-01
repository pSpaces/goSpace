package container

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pspaces/gospace/function"
)

// Template structure used for matching against tuples.
// Template is a tuple with type information used for pattern matching.
type Template struct {
	Flds []interface{} `bson:"fields" json:"fields" xml:"fields"`
}

// NewTemplate creates a template from the variadic fields provided.
// NewTemplate encapsulates the types of pointer values.
func NewTemplate(fields ...interface{}) (tp Template) {
	tempfields := make([]interface{}, len(fields))
	copy(tempfields, fields)

	// Replace pointers with reflect.Type value used to match type.
	for i, value := range fields {
		if value != nil {
			if reflect.TypeOf(value).Kind() == reflect.Ptr {
				// Encapsulate the parameter with a TypeField.
				tempfields[i] = CreateTypeField(reflect.ValueOf(value).Elem().Interface())
			} else if function.IsFunc(value) && reflect.TypeOf(value).Kind() == reflect.Interface {
				tempfields[i] = CreateTypeField(value)
			}
		}
	}

	tp = Template{tempfields}

	return tp
}

// Equal returns true if both templates tp and tq are strictly equal, and false otherwise.
func (tp *Template) Equal(tq *Template) (e bool) {
	e = false

	if tp == nil && tp == tq {
		e = true
	} else if tp != nil && tq != nil {
		e = tp.Length() == tq.Length()

		length := tp.Length()
		for i := 0; i < length && e; i++ {
			tpf := tp.GetFieldAt(i)
			tqf := tq.GetFieldAt(i)

			if reflect.TypeOf(tpf) == reflect.TypeOf(TypeField{}) && reflect.TypeOf(tqf) == reflect.TypeOf(TypeField{}) {
				atpf := tpf.(TypeField)
				btpf := tqf.(TypeField)
				e = e && atpf.Equal(btpf)
			} else if function.IsFunc(tpf) && function.IsFunc(tqf) {
				e = e && function.Name(tpf) == function.Name(tqf) && function.Signature(tpf) == function.Signature(tqf)
			} else {
				e = e && reflect.DeepEqual(tpf, tqf)
			}
		}
	}

	return e
}

// Match returns true if both templates tp and tq are equivalent up to their types, and false otherwise.
func (tp *Template) Match(tq *Template) (e bool) {
	e = false

	if tp == nil && tp == tq {
		e = true
	} else if tp != nil && tq != nil {
		e = tp.Length() == tq.Length()

		length := tp.Length()
		for i := 0; i < length && e; i++ {
			tpf := tp.GetFieldAt(i)
			tqf := tq.GetFieldAt(i)

			if reflect.TypeOf(tpf) == reflect.TypeOf(TypeField{}) || reflect.TypeOf(tqf) == reflect.TypeOf(TypeField{}) {
				atpf, ais := tpf.(TypeField)
				btpf, bis := tqf.(TypeField)

				if ais && bis {
					e = e && atpf.Equal(btpf)
				} else if ais {
					e = e && reflect.TypeOf(tqf) == atpf.GetType()
				} else if bis {
					e = e && reflect.TypeOf(tpf) == btpf.GetType()
				}
			} else if function.IsFunc(tpf) && function.IsFunc(tqf) {
				e = e && function.Name(tpf) == function.Name(tqf) && function.Signature(tpf) == function.Signature(tqf)
			} else {
				e = e && reflect.DeepEqual(tpf, tqf)
			}
		}
	}

	return e
}

// ExactMatch returns true if both templates tp and tq are equivalent up to their types and the maximum amount of concrete values matched w.r.t. tq given the amount of concrete values in tp.
func (tp *Template) ExactMatch(tq *Template) (e bool, pno, qno uint) {
	pno, qno = 0, 0
	e = false

	if tp == nil && tp == tq {
		pno++
		qno++
		e = true
	} else if tp != nil && tq != nil {
		e = tp.Length() == tq.Length()

		length := tp.Length()
		for i := 0; i < length && e; i++ {
			tpf := tp.GetFieldAt(i)
			tqf := tq.GetFieldAt(i)

			if reflect.TypeOf(tpf) == reflect.TypeOf(TypeField{}) || reflect.TypeOf(tqf) == reflect.TypeOf(TypeField{}) {
				atpf, ais := tpf.(TypeField)
				btpf, bis := tqf.(TypeField)

				if ais && bis {
					e = e && atpf.Equal(btpf)

					if e {
						pno++
					}
				} else if ais {
					e = e && reflect.TypeOf(tqf) == atpf.GetType()
				} else if bis {
					e = e && reflect.TypeOf(tpf) == btpf.GetType()
				}

				if e {
					qno++
				}
			} else if function.IsFunc(tpf) && function.IsFunc(tqf) {
				nm := function.Name(tpf) == function.Name(tqf)
				sm := function.Signature(tpf) == function.Signature(tqf)
				e = e && nm && sm

				if e {
					pno++
				}

				if e && sm {
					qno++
				}
			} else {
				e = e && reflect.DeepEqual(tpf, tqf)

				if e {
					pno++
					qno++
				}
			}
		}
	}

	return e, pno, qno
}

// Length returns the size of the template.
func (tp *Template) Length() (sz int) {
	sz = -1

	if tp != nil {
		sz = len((*tp).Flds)
	}

	return sz
}

// Fields returns the fields of the template.
func (tp *Template) Fields() (flds []interface{}) {
	if tp != nil {
		flds = (*tp).Flds
	}

	return flds
}

// GetFieldAt returns the i'th field of the template.
func (tp *Template) GetFieldAt(i int) (fld interface{}) {
	if tp != nil {
		fld = (*tp).Flds[i]
	}

	return fld
}

// setFieldAt sets the i'th field of the template to the value of val.
// setFieldAt returns true if the field is set, and false otherwise.
func (tp *Template) setFieldAt(i int, val interface{}) (b bool) {
	b = false

	if tp != nil {
		(*tp).Flds[i] = val
		b = true
	}

	return b
}

// Apply iterates through the template tp and applies the function fun to each field.
// Apply returns true function fun could be applied to all the fields, and false otherwise.
func (tp *Template) Apply(fun func(field interface{}) interface{}) (b bool) {
	b = false

	if tp != nil {
		b = true
		for i := 0; i < tp.Length(); i++ {
			b = b && tp.setFieldAt(i, fun(tp.GetFieldAt(i)))
		}
	}

	return b
}

// NewTuple returns a new tuple t from the template.
// NewTuple initializes all tuple fields in t with empty values depending on types in the template.
func (tp *Template) NewTuple() (t Tuple) {
	var element interface{}
	param := make([]interface{}, tp.Length())

	for i := range param {
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

	t = NewTuple(param...)

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

	for i := range strs {
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
