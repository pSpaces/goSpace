package space

import (
	"reflect"

	"github.com/google/uuid"
	"github.com/pspaces/gospace/container"
	"github.com/pspaces/gospace/policy"
	"github.com/pspaces/gospace/protocol"
)

// Interspace defines the internal space interface.
// Interspace interface is meant to be used by both external or internal interfaces.
type Interspace interface {
	Put(tuple ...interface{}) (container.Tuple, error)
	Get(template ...interface{}) (container.Tuple, error)
	Query(template ...interface{}) (container.Tuple, error)
	PutP(tuple ...interface{}) (container.Tuple, error)
	GetP(template ...interface{}) (container.Tuple, error)
	QueryP(template ...interface{}) (container.Tuple, error)
	GetAll(template ...interface{}) ([]container.Tuple, error)
	QueryAll(template ...interface{}) ([]container.Tuple, error)
}

// Interstar defines the internal space aggregation interface.
// Interstar interface is meant to be used by both external or internal interfaces.
type Interstar interface {
	PutAgg(function interface{}, template ...interface{}) (container.Tuple, error)
	GetAgg(function interface{}, template ...interface{}) (container.Tuple, error)
	QueryAgg(function interface{}, template ...interface{}) (container.Tuple, error)
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

// Intercellestial defines the internal space aggregation interface without any error checking.
// Intercellestial interface is meant primarily for internal usage.
// This interface can change without any notice.
type Intercellestial interface {
	RawPutAgg(function interface{}, template ...interface{}) (interface{}, interface{})
	RawGetAgg(function interface{}, template ...interface{}) (interface{}, interface{})
	RawQueryAgg(function interface{}, template ...interface{}) (interface{}, interface{})
}

// Space is a structure for interacting with a space.
type Space struct {
	id string
	ts *TupleSpace
	p  *protocol.PointToPoint
}

// NewSpace creates an empty space s with the specified URL.
func NewSpace(url string, cp ...*policy.Composable) (s Space) {
	id := uuid.New()
	sid, err := id.MarshalText()

	if err == nil {
		// TODO: This is a workaround. Refactor.
		s = Space{string(sid), nil, nil}
		p, ts := NewSpaceAlt(url, cp...)
		s.ts = ts
		s.p = p
	}

	return s
}

// NewRemoteSpace connects to a remote space rs with the specified URL.
func NewRemoteSpace(url string) (rs Space) {
	id := uuid.New()
	sid, err := id.MarshalText()

	if err == nil {
		p, ts := NewRemoteSpaceAlt(url)
		rs = Space{string(sid), ts, p}
	}

	return rs
}

// ID returns the identifier for space s.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) ID() (id string, e error) {
	if s != nil {
		id = (*s).id
	}

	e = NewSpaceError(s, id, nil)

	return id, e
}

// Size retrieves the size of space s at this instant.
// Size returns the space size or -1 if it was not possible to determine the size at this instant.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) Size() (sz int, e error) {
	var result int
	var status interface{}

	if s != nil {
		rawres, rawerr := (*s).RawSize()
		result = rawres.(int)
		status = rawerr
	} else {
		result = -1
	}

	e = NewSpaceError(s, -1, status)

	if e == nil {
		sz = result
	} else {
		sz = -1
	}

	return sz, e
}

// RawSize retrieves the size of space s at this instant without any error checking.
// RawSize returns the implementation result sz and error state e.
func (s *Space) RawSize() (sz interface{}, e interface{}) {
	sz, e = Size(*s.p)
	return sz, e
}

// Put performs a blocking placement a tuple t into space s.
// Put returns the original tuple tp and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) Put(t ...interface{}) (tp container.Tuple, e error) {
	var result container.Tuple
	var status interface{}

	if s != nil {
		rawres, rawerr := (*s).RawPut(t...)
		result = rawres.(container.Tuple)
		status = rawerr
	} else {
		result = container.NewTuple(nil)
	}

	e = NewSpaceError(s, container.NewTuple(t...), status)

	if e == nil {
		tp = result
	} else {
		tp = container.NewTuple(nil)
	}

	return tp, e
}

// RawPut performs a blocking placement of a tuple t into space s without any error checking.
// RawPut returns the implementation result tp and error state e.
func (s *Space) RawPut(t ...interface{}) (tp interface{}, e interface{}) {
	tp, e = Put(*s.p, t...)
	return tp, e
}

// Get performs a blocking retrieval for a tuple from space s with template t.
// Get returns the matched tuple tp and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) Get(t ...interface{}) (tp container.Tuple, e error) {
	var result container.Tuple
	var status interface{}

	if s != nil {
		rawres, rawerr := (*s).RawGet(t...)
		result = rawres.(container.Tuple)
		status = rawerr
	} else {
		result = container.NewTuple(nil)
	}

	e = NewSpaceError(s, container.NewTemplate(t...), status)

	if e == nil {
		tp = result
	} else {
		tp = container.NewTuple(nil)
	}

	return tp, e
}

// RawGet performs a blocking retrieval a tuple from space s with template t and without any error checking.
// RawGet returns the implementation result tp and error state e.
func (s *Space) RawGet(t ...interface{}) (tp interface{}, e interface{}) {
	tp, e = Get(*s.p, t...)
	return tp, e
}

// Query performs a blocking query for a tuple from space s with template t.
// Query returns the matched tuple tp and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) Query(t ...interface{}) (tp container.Tuple, e error) {
	var result container.Tuple
	var status interface{}

	if s != nil {
		rawres, rawerr := (*s).RawQuery(t...)
		result = rawres.(container.Tuple)
		status = rawerr
	} else {
		result = container.NewTuple(nil)
	}

	e = NewSpaceError(s, container.NewTemplate(t...), status)

	if e == nil {
		tp = result
	} else {
		tp = container.NewTuple(nil)
	}

	return tp, e
}

// RawQuery performs a blocking query for a tuple from space s with template t and without any error checking.
// RawQuery returns the implementation result tp and error state e.
func (s *Space) RawQuery(t ...interface{}) (tp interface{}, e interface{}) {
	tp, e = Query(*s.p, t...)
	return tp, e
}

// PutP performs a non-blocking placement a tuple t into space s.
// PutP returns the original tuple tp and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) PutP(t ...interface{}) (tp container.Tuple, e error) {
	var result container.Tuple
	var status interface{}

	if s != nil {
		rawres, rawerr := (*s).RawPutP(t...)
		result = rawres.(container.Tuple)
		status = rawerr
	} else {
		result = container.NewTuple(nil)
	}

	e = NewSpaceError(s, container.NewTuple(t...), status)

	if e == nil {
		tp = result
	} else {
		tp = container.NewTuple(nil)
	}

	return tp, e
}

// RawPutP performs a non-blocking placement of a tuple t into space s without any error checking.
// RawPutP returns the implementation result tp and error state e.
func (s *Space) RawPutP(t ...interface{}) (tp interface{}, e interface{}) {
	tp, e = PutP(*s.p, t...)
	return tp, e
}

// GetP performs a non-blocking retrieval for a tuple from space s with template t.
// GetP returns the matched tuple tp and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) GetP(t ...interface{}) (tp container.Tuple, e error) {
	var result container.Tuple
	var status interface{}

	if s != nil {
		rawres, rawerr := (*s).RawGetP(t...)
		result = rawres.(container.Tuple)
		status = rawerr
	} else {
		result = container.NewTuple(nil)
	}

	e = NewSpaceError(s, container.NewTemplate(t...), status)

	if e == nil {
		tp = result
	} else {
		tp = container.NewTuple(nil)
	}

	return tp, e
}

// RawGetP performs a non-blocking retrieval a tuple from space s with template t and without any error checking.
// RawGetP returns the implementation result tp and error state e.
func (s *Space) RawGetP(t ...interface{}) (tp interface{}, e interface{}) {
	tp, e, _ = GetP(*s.p, t...)
	return tp, e
}

// QueryP performs a non-blocking query for a tuple from space s with template t.
// QueryP returns the matched tuple tp and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) QueryP(t ...interface{}) (tp container.Tuple, e error) {
	var result container.Tuple
	var status interface{}

	if s != nil {
		rawres, rawerr := (*s).RawQueryP(t...)
		result = rawres.(container.Tuple)
		status = rawerr
	} else {
		result = container.NewTuple(nil)
	}

	e = NewSpaceError(s, container.NewTemplate(t...), status)

	if e == nil {
		tp = result
	} else {
		tp = container.NewTuple(nil)
	}

	return tp, e
}

// RawQueryP performs a blocking query for a tuple from space s with template t and without any error checking.
// RawQueryP returns the implementation result tp and error state e.
func (s *Space) RawQueryP(t ...interface{}) (tp interface{}, e interface{}) {
	tp, e, _ = QueryP(*s.p, t...)
	return tp, e
}

// GetAll performs a non-blocking retrieval for all tuples from space s with template t.
// GetAll returns the matching tuples ts and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) GetAll(t ...interface{}) (ts []container.Tuple, e error) {
	var result []container.Tuple
	var status interface{}

	if s != nil {
		rawres, rawerr := (*s).RawGetAll(t...)
		result = rawres.([]container.Tuple)
		status = rawerr
	} else {
		result = []container.Tuple{}
	}

	e = NewSpaceError(s, container.NewTemplate(t...), status)

	if e == nil {
		ts = result
	} else {
		ts = []container.Tuple{}
	}

	return ts, e
}

// RawGetAll performs a non-blocking retrieval for all tuples from space s with template t and without any error checking.
// RawGetAll returns the implementation result ts and error state e.
func (s *Space) RawGetAll(t ...interface{}) (ts interface{}, e interface{}) {
	ts, e = GetAll(*s.p, t...)
	return ts, e
}

// QueryAll performs a non-blocking query for all tuples from space s with template t.
// QueryAll returns the matching tuples ts and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) QueryAll(t ...interface{}) (ts []container.Tuple, e error) {
	var result []container.Tuple
	var status interface{}

	if s != nil {
		rawres, rawerr := (*s).RawQueryAll(t...)
		result = rawres.([]container.Tuple)
		status = rawerr
	} else {
		result = []container.Tuple{}
	}

	e = NewSpaceError(s, container.NewTemplate(t...), status)

	if e == nil {
		ts = result
	} else {
		ts = []container.Tuple{}
	}

	return ts, e
}

// RawQueryAll performs a non-blocking query for all tuples from space s with template t and without any error checking.
// RawQueryAll returns the implementation result ts and error state e.
func (s *Space) RawQueryAll(t ...interface{}) (ts interface{}, e interface{}) {
	ts, e = QueryAll(*s.p, t...)
	return ts, e
}

// PutAgg performs a non-blocking aggregation placement on all tuples from space s that matches template t.
// PutAgg uses an aggregation function f to aggregate a pair of tuples into one.
// PutAgg places either the aggregate tuple, or the intrinsic tuple belonging to a template t if no matching tuples are returned, back into s.
// PutAgg returns either the aggregate or intrinsic tuple and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) PutAgg(f interface{}, t ...interface{}) (tp container.Tuple, e error) {
	var result container.Tuple
	var status interface{}

	if s != nil {
		rawres, rawerr := (*s).RawPutAgg(f, t...)
		result = rawres.(container.Tuple)
		status = rawerr
	} else {
		result = container.NewTuple(nil)
	}

	e = NewSpaceError(s, container.NewTemplate(t...), status)

	if e == nil {
		tp = result
	} else {
		tp = container.NewTuple(nil)
	}

	return tp, e
}

// RawPutAgg performs a non-blocking aggregation retrieval on all tuples from space that match template t and without any error checking.
// RawPutAgg uses an aggregation function f to aggregate a pair of tuples into one.
// RawPutAgg returns the implementation result tp and error state e.
func (s *Space) RawPutAgg(f interface{}, t ...interface{}) (tp interface{}, e interface{}) {
	tp, e = PutAgg(*s.p, f, t...)
	return tp, e
}

// GetAgg performs a non-blocking aggregation retrieval on all tuples from space s that matches template t.
// GetAgg uses an aggregation function f to aggregate a pair of tuples into one.
// GetAgg returns an aggregate tuple tp and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) GetAgg(f interface{}, t ...interface{}) (tp container.Tuple, e error) {
	var result container.Tuple
	var status interface{}

	if s != nil {
		rawres, rawerr := (*s).RawGetAgg(f, t...)
		result = rawres.(container.Tuple)
		status = rawerr
	} else {
		result = container.NewTuple(nil)
	}

	e = NewSpaceError(s, container.NewTemplate(t...), status)

	if e == nil {
		tp = result
	} else {
		tp = container.NewTuple(nil)
	}

	return tp, e
}

// RawGetAgg performs a non-blocking aggregation retrieval on all tuples from space that match template t and without any error checking.
// RawGetAgg uses an aggregation function f to aggregate a pair of tuples into one.
// RawGetAgg returns the implementation result tp and error state e.
func (s *Space) RawGetAgg(f interface{}, t ...interface{}) (tp interface{}, e interface{}) {
	tp, e = GetAgg(*s.p, f, t...)
	return tp, e
}

// QueryAgg performs a non-blocking aggregation query on all tuples from space s that matches template t.
// QueryAgg uses an aggregation function f to aggregate a pair of tuples into one.
// QueryAgg returns an aggregate tuple tp and an error e.
// Error e contains a structure adhering to the error interface if the operation fails, and nil if no error occured.
func (s *Space) QueryAgg(f interface{}, t ...interface{}) (tp container.Tuple, e error) {
	var result container.Tuple
	var status interface{}

	if s != nil {
		rawres, rawerr := (*s).RawQueryAgg(f, t...)
		result = rawres.(container.Tuple)
		status = rawerr
	} else {
		result = container.NewTuple(nil)
	}

	e = NewSpaceError(s, container.NewTemplate(t...), status)

	if e == nil {
		tp = result
	} else {
		tp = container.NewTuple(nil)
	}

	return tp, e
}

// RawQueryAgg performs a non-blocking aggregation query on all tuples from space that match template t and without any error checking.
// RawQueryAgg uses an aggregation function f to aggregate a pair of tuples into one.
// RawQueryAgg returns the implementation result tp and error state e.
func (s *Space) RawQueryAgg(f interface{}, t ...interface{}) (tp interface{}, e interface{}) {
	tp, e = QueryAgg(*s.p, f, t...)
	return tp, e
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

// InterpretOperation returns a status for an operation that has been succesful given a return state by an operation.
// The state is given by the implementation and this method maps the returned state to a boolean value.
// This is an internal method and may change without notice.
func (s *Space) InterpretOperation(state interface{}) (status bool) {
	status = false

	if s != nil && state != nil {
		status = state.(bool)
	}

	return status
}

// InterpretValue returns a representation of the value that was passed to the operation.
// The representation for now is a print friendly string value.
// This is an internal method and may change without notice.
func (s *Space) InterpretValue(value interface{}) (str string) {
	if s != nil {
		if value == nil {
			str = "nil"
		} else if reflect.TypeOf(value) == reflect.TypeOf(container.Tuple{}) {
			str = value.(container.Tuple).String()
		} else if reflect.TypeOf(value) == reflect.TypeOf(container.Template{}) {
			str = value.(container.Template).String()
		}
	}

	return str
}
