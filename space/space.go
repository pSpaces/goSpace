package space

import (
	. "github.com/pspaces/gospace/protocol"
	. "github.com/pspaces/gospace/shared"
	"reflect"
)

// Interspace defines the internal space interface.
// Interspace interface is meant to be used by both external or internal interfaces.
type Interspace interface {
	Put(tuple ...interface{}) (Tuple, error)
	Get(template ...interface{}) (Tuple, error)
	Query(template ...interface{}) (Tuple, error)
	PutP(tuple ...interface{}) (Tuple, error)
	GetP(template ...interface{}) (Tuple, error)
	QueryP(template ...interface{}) (Tuple, error)
	GetAll(template ...interface{}) ([]Tuple, error)
	QueryAll(template ...interface{}) ([]Tuple, error)
}

// Interstellar defines the internal space interface without any error checking.
// Interstellar interface is meant primarily for internal usage.
// This interface can change without any notice.
type Interstellar interface {
	RawPut(tuple ...interface{}) (interface{}, interface{})
	RawGet(template ...interface{}) (interface{}, interface{})
	RawQuery(template ...interface{}) (interface{}, interface{})
	RawPutP(tuple ...interface{}) (interface{}, interface{})
	RawGetP(template ...interface{}) (interface{}, interface{})
	RawQueryP(template ...interface{}) (interface{}, interface{})
	RawGetAll(template ...interface{}) (interface{}, interface{})
	RawQueryAll(template ...interface{}) (interface{}, interface{})
}

// Space is a structure for interacting with a space.
type Space struct {
	id string
	p  PointToPoint
}

// NewSpace creates an empty space with the specified URL.
func NewSpace(url string) Space {
	return Space{url, NewSpaceAlt(url)}
}

// NewRemoteSpace connects to a remote space with the specified URL.
func NewRemoteSpace(url string) Space {
	return Space{url, NewRemoteSpaceAlt(url)}
}

// Put performs a blocking placement a tuple t into space s.
// Put returns the original tuple tp and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) Put(t ...interface{}) (tp Tuple, e error) {
	var result Tuple
	var status interface{} = nil

	if s != nil {
		rawres, rawerr := (*s).RawPut(t...)
		result = rawres.(Tuple)
		status = rawerr
	} else {
		result = CreateTuple(nil)
	}

	tp = result
	e = NewSpaceError(s, status)

	return tp, e
}

// RawPut performs a blocking placement of a tuple t into space s without any error checking.
// RawPut returns the implementation result tp and error state e.
func (s *Space) RawPut(t ...interface{}) (tp interface{}, e interface{}) {
	e = Put((*s).p, t...)
	tp = TupleFromTemplate(t...)
	return tp, e
}

// Get performs a blocking retrieval for a tuple from space s with template t.
// Get returns the matched tuple tp and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) Get(t ...interface{}) (tp Tuple, e error) {
	var result Tuple
	var status interface{} = nil

	if s != nil {
		rawres, rawerr := (*s).RawGet(t...)
		result = rawres.(Tuple)
		status = rawerr
	} else {
		result = CreateTuple(nil)
	}

	tp = result
	e = NewSpaceError(s, status)

	return tp, e
}

// RawGet performs a blocking retrieval a tuple from space s with template t and without any error checking.
// RawGet returns the implementation result tp and error state e.
func (s *Space) RawGet(t ...interface{}) (tp interface{}, e interface{}) {
	e = Get((*s).p, t...)
	tp = TupleFromTemplate(t...)
	return tp, e
}

// Query performs a blocking query for a tuple from space s with template t.
// Query returns the matched tuple tp and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) Query(t ...interface{}) (tp Tuple, e error) {
	var result Tuple
	var status interface{} = nil

	if s != nil {
		rawres, rawerr := (*s).RawQuery(t...)
		result = rawres.(Tuple)
		status = rawerr
	} else {
		result = CreateTuple(nil)
	}

	tp = result
	e = NewSpaceError(s, status)

	return tp, e
}

// RawQuery performs a blocking query for a tuple from space s with template t and without any error checking.
// RawQuery returns the implementation result tp and error state e.
func (s *Space) RawQuery(t ...interface{}) (tp interface{}, e interface{}) {
	e = Query((*s).p, t)
	tp = TupleFromTemplate(t...)
	return tp, e
}

// PutP performs a non-blocking placement a tuple t into space s.
// PutP returns the original tuple tp and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) PutP(t ...interface{}) (tp Tuple, e error) {
	var result Tuple
	var status interface{} = nil

	if s != nil {
		rawres, rawerr := (*s).RawPutP(t...)
		result = rawres.(Tuple)
		status = rawerr
	} else {
		result = CreateTuple(nil)
	}

	tp = result
	e = NewSpaceError(s, status)

	return tp, e
}

// RawPut performs a non-blocking placement of a tuple t into space s without any error checking.
// RawPut returns the implementation result tp and error state e.
func (s *Space) RawPutP(t ...interface{}) (tp interface{}, e interface{}) {
	tp = TupleFromTemplate(t...)
	e = PutP((*s).p, t...)
	return tp, e
}

// GetP performs a non-blocking retrieval for a tuple from space s with template t.
// GetP returns the matched tuple tp and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) GetP(t ...interface{}) (tp Tuple, e error) {
	var result Tuple
	var status interface{} = nil

	if s != nil {
		rawres, rawerr := (*s).RawGetP(t...)
		result = rawres.(Tuple)
		status = rawerr
	} else {
		result = CreateTuple(nil)
	}

	tp = result
	e = NewSpaceError(s, status)

	return tp, e
}

// RawGetP performs a non-blocking retrieval a tuple from space s with template t and without any error checking.
// RawGetP returns the implementation result tp and error state e.
func (s *Space) RawGetP(t ...interface{}) (tp interface{}, e interface{}) {
	_, e = GetP((*s).p, t...)
	tp = TupleFromTemplate(t...)
	return tp, e
}

// QueryP performs a non-blocking query for a tuple from space s with template t.
// QueryP returns the matched tuple tp and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) QueryP(t ...interface{}) (tp Tuple, e error) {
	var result Tuple
	var status interface{} = nil

	if s != nil {
		rawres, rawerr := (*s).RawQueryP(t...)
		result = rawres.(Tuple)
		status = rawerr
	} else {
		result = CreateTuple(nil)
	}

	tp = result
	e = NewSpaceError(s, status)

	return tp, e
}

// RawQueryP performs a blocking query for a tuple from space s with template t and without any error checking.
// RawQueryP returns the implementation result tp and error state e.
func (s *Space) RawQueryP(t ...interface{}) (tp interface{}, e interface{}) {
	_, e = QueryP((*s).p, t...)
	tp = TupleFromTemplate(t...)
	return tp, e
}

// GetAll performs a non-blocking retrieval for all tuples from space s with template t.
// GetAll returns the matching tuples ts and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) GetAll(t ...interface{}) (ts []Tuple, e error) {
	var result []Tuple
	var status interface{} = nil

	if s != nil {
		rawres, rawerr := (*s).RawGetAll(t...)
		result = rawres.([]Tuple)
		status = rawerr
	} else {
		result = []Tuple{}
	}

	ts = result
	e = NewSpaceError(s, status)

	return ts, e
}

// RawGetAll performs a non-blocking retrieval for all tuples from space s with template t and without any error checking.
// RawGetAll returns the implementation result ts and error state e.
func (s *Space) RawGetAll(t ...interface{}) (ts interface{}, e interface{}) {
	ts, e = GetAll((*s).p, t...)
	return ts, e
}

// QueryAll performs a non-blocking query for all tuples from space s with template t.
// QueryAll returns the matching tuples ts and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) QueryAll(t ...interface{}) (ts []Tuple, e error) {
	var result []Tuple
	var status interface{} = nil

	if s != nil {
		rawres, rawerr := (*s).RawQueryAll(t...)
		result = rawres.([]Tuple)
		status = rawerr
	} else {
		result = []Tuple{}
	}

	ts = result
	e = NewSpaceError(s, status)

	return ts, e
}

// RawQueryAll performs a non-blocking query for all tuples from space s with template t and without any error checking.
// RawQueryAll returns the implementation result ts and error state e.
func (s *Space) RawQueryAll(t ...interface{}) (ts interface{}, e interface{}) {
	ts, e = QueryAll((*s).p, t...)
	return ts, e
}

// InterpretError returns an error message msg given a return state by an operation.
// The state is given by the implementation and this method maps from the state to sane
// error messages. This is an internal method and may change without notice.
func (s *Space) InterpretError(state interface{}) (msg string) {
	const (
		InterSpaceInvalid = iota
		InterStateInvalid
		InterOperationFailed
		InterOperationSuccess
	)

	var errMsg = map[int]string{
		InterSpaceInvalid:     "trying to interpret error on an invalid space",
		InterStateInvalid:     "trying to interpret error on an invalid state",
		InterOperationFailed:  "operation on this space failed",
		InterOperationSuccess: "operation on this space succeeded",
	}

	if s != nil {
		if state != nil {
			status := state.(bool)

			if status {
				msg = errMsg[InterOperationSuccess]
			} else {
				msg = errMsg[InterOperationFailed]
			}
		} else {
			msg = errMsg[InterStateInvalid]
		}
	} else {
		msg = errMsg[InterSpaceInvalid]
	}

	return msg
}

// InterpretError returns a status for an operation that has been succesful given a return state by an operation.
// The state is given by the implementation and this method maps the returned state to a boolean value.
// This is an internal method and may change without notice.
func (s *Space) InterpretOperation(state interface{}) (status bool) {
	status = false

	if s != nil && state != nil {
		status = state.(bool)
	}

	return status
}

// TupleFromTemplate reads a template and returns a new tuple tp.
// TupleFromTemplate extracts values from any pointers it finds in template t.
func TupleFromTemplate(t ...interface{}) (tp Tuple) {
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
