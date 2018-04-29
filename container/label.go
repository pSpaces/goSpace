package container

import (
	"fmt"
)

// Interlabel is an internal interface for manipulating a label.
type Interlabel interface {
	Id() (id string)
}

// Label is a structure used to label arbitrary data.
type Label Tuple

// NewLabel creates a new label with identifier id and an optional value v.
func NewLabel(id string) (l Label) {
	params := []interface{}{id}
	l = Label(NewTuple(params...))
	return l
}

// DeepCopy returns a deep copy of the label {
func (l *Label) DeepCopy() (lc Label) {
	lc = NewLabel(l.ID())
	return lc
}

// Equal returns true if both labels l and m are equivalent, and false otherwise.
func (l *Label) Equal(m *Label) (e bool) {
	e = false

	if l == nil && l == m {
		e = true
	} else if l != nil && m != nil {
		idl := l.ID()
		idm := m.ID()
		e = len(idl) == len(idm)
		if e {
			e = idl == idm
		}
	}

	return e
}

// ID returns label l's identifier id.
func (l *Label) ID() (id string) {
	t := Tuple(*l)
	id = (&t).GetFieldAt(0).(string)
	return id
}

// ParenthesisType returns a pair of strings that encapsulates the tuple.
// ParenthesisType is used in the String() method.
func (l Label) ParenthesisType() (ls string, rs string) {
	ls = "|"
	rs = "|"
	return ls, rs
}

// Delimiter returns the delimiter used to seperated the values in label l.
// Delimiter is used in the String() method.
func (l Label) Delimiter() string {
	return ", "
}

// String returns a print friendly representation of label l.
func (l Label) String() (s string) {
	s = fmt.Sprintf("|%s|", l.ID())
	return s
}
