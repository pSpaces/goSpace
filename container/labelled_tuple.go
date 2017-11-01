package container

import (
	"fmt"
	"reflect"
	"strings"
)

// LabelledTuple is a labelled tuple containing a list of fields and a label set.
// Fields can be any data type and is used to store data.
// TupleLabels is a label set that is associated to the tuple itself.
type LabelledTuple Tuple

// NewLabelledTuple creates a labelled tuple according to the labels and values present in the fields.
// NewLabelledTuple searches the first field for labels.
func NewLabelledTuple(fields ...interface{}) (lt LabelledTuple) {
	if len(fields) == 0 {
		lt = LabelledTuple(NewTuple(Labels{}))
	} else {
		var lbls Labels

		if len(fields) == 0 {
			lbls = Labels{}
		} else {
			lbls = fields[0].(Labels)
			lblsc := make(Labels)
			for _, v := range lbls.Labelling() {
				lbl := lbls.Retrieve(v)
				lblsc.Add(NewLabel(lbl.ID()))
			}
			lbls = lblsc
		}

		if len(fields) < 1 {
			lt = LabelledTuple(NewTuple(lbls))
		} else {
			fields[0] = lbls
			lt = LabelledTuple(NewTuple(fields...))
		}
	}

	return lt
}

// Length returns the amount of fields of the tuple.
func (lt *LabelledTuple) Length() (sz int) {
	sz = -1

	if lt != nil {
		sz = len((*lt).Flds) - 1
	}

	return
}

// Fields returns the fields of the tuple.
func (lt *LabelledTuple) Fields() (flds []interface{}) {
	if lt != nil {
		if lt.Length() >= 2 {
			flds = (*lt).Flds[1:]
		}
	}

	return flds
}

// GetFieldAt returns the i'th field of the tuple.
func (lt *LabelledTuple) GetFieldAt(i int) (fld interface{}) {
	if lt != nil && i >= 0 && i < lt.Length() {
		fld = (*lt).Flds[i+1]
	}

	return
}

// SetFieldAt sets the i'th field of the tuple to the value of val.
func (lt *LabelledTuple) SetFieldAt(i int, val interface{}) (b bool) {
	if lt != nil && i >= 0 && i < lt.Length() {
		(*lt).Flds[i+1] = val
		b = true
	}

	return b
}

// Apply iterates through the labelled tuple t and applies the function fun to each field.
// Apply returns true function fun could be applied to all the fields, and false otherwise.
func (lt *LabelledTuple) Apply(fun func(field interface{}) interface{}) (b bool) {
	b = false

	if lt != nil {
		b = true
		for i := 0; i < lt.Length(); i++ {
			lt.SetFieldAt(i, fun(lt.GetFieldAt(i)))
		}
	}

	return b
}

// Labels returns the label set belonging to the labelled tuple.
func (lt *LabelledTuple) Labels() (ls Labels) {
	return (*lt).Flds[0].(Labels)
}

// Tuple returns a tuple without the label.
func (lt *LabelledTuple) Tuple() (t Tuple) {
	t = NewTuple((*lt).Flds[1:]...)
	return t
}

// MatchLabels matches a labelled tuples t labels against labels ls.
func (lt *LabelledTuple) MatchLabels(ls Labels) (mls *Labels, b bool) {
	b = lt != nil

	if b {
		ltls := lt.Labels()
		mls, b = ltls.Intersect(&ls)
	}

	return mls, b
}

// MatchTemplate pattern matches labelled tuple t against the template tp.
// MatchTemplate discriminates between encapsulated formal fields and actual fields.
// MatchTemplate returns true if the template matches the labelled tuple and false otherwise.
func (lt *LabelledTuple) MatchTemplate(tp Template) (b bool) {
	b = lt != nil && lt.Length() == (&tp).Length()

	if b {
		t := lt.Tuple()
		b = (&t).Match(tp)
	}

	return b
}

// ParenthesisType returns a pair of strings that encapsulates labelled tuple t.
// ParenthesisType is used in the String() method.
func (lt LabelledTuple) ParenthesisType() (string, string) {
	return "(", ")"
}

// Delimiter returns the delimiter used to seperate a labelled tuple t's fields.
// Delimiter is used in the String() method.
func (lt LabelledTuple) Delimiter() string {
	return ", "
}

// String returns a print friendly representation of the tuple.
func (lt LabelledTuple) String() (s string) {
	ld, rd := lt.ParenthesisType()

	delim := lt.Delimiter()

	strs := make([]string, lt.Length())

	for i := range strs {
		field := lt.GetFieldAt(i)

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

	s = fmt.Sprintf("%s%s%s%s%s", ld, lt.Labels(), " : ", strings.Join(strs, delim), rd)

	return s
}
