package space

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// SpaceError represents an internal error type used when printing error messages.
type SpaceError struct {
	msg    string
	pkg    string
	fun    string
	sid    string
	sop    bool
	status interface{}
}

const (
	SpaceInvalid = iota
	SpaceNoErrorMethod
	SpaceOperationFailed
)

// errMsg contains all the generic error messages when operating on spaces.
var errMsg = map[int]string{
	SpaceInvalid:         "operation performed on an invalid space",
	SpaceNoErrorMethod:   "no error method available, can't interpret error",
	SpaceOperationFailed: "could not perform operation on this space",
}

// NewSpaceError creates a new error given space spc and the return state of the implemented operation.
// NewSpaceError returns a structure which fulfils the error interface and if an operation error has occured.
// NewSpaceError returns nil if no operation failure has occured.
func NewSpaceError(spc *Space, state interface{}) error {
	var msg, pkg, fun, sid string
	var err error
	var sop bool
	var status interface{}

	pkg, fun = getCalleInfo(2)

	if spc == nil {
		sid = "nil"
		msg = errMsg[SpaceInvalid]
		status = nil
	} else {
		sid = (*spc).id

		spct := reflect.ValueOf(spc)

		method := spct.MethodByName("InterpretError")
		if method.IsValid() {
			vals := method.Call([]reflect.Value{reflect.ValueOf(state)})
			if len(vals) == 1 {
				msg = vals[0].String()
			} else {
				msg = errMsg[SpaceNoErrorMethod]
			}
		} else {
			msg = errMsg[SpaceNoErrorMethod]
		}

		method = spct.MethodByName("InterpretOperation")
		if method.IsValid() {
			vals := method.Call([]reflect.Value{reflect.ValueOf(state)})
			if len(vals) == 1 {
				sop = vals[0].Interface().(bool)
			} else {
				sop = false
			}
		} else {
			sop = false
		}

		status = state
	}

	if sop == true {
		err = nil
	} else {
		err = SpaceError{msg, pkg, fun, sid, sop, status}
	}

	return err
}

// Operation returns a boolean value if an operation has succeeded.
func (e SpaceError) Operation() bool {
	return e.sop
}

// Error prints the error message represented by SpaceError.
func (e SpaceError) Error() string {
	return fmt.Sprintf("%s: %s(%s).%s: %s.", e.pkg, "Space", e.sid, e.fun, e.msg)
}

// getCalleInfo determines the package and function names associated to a function call.
// getCalleInfo uses the runtime package, and no file or line information is provided,
// since this can not be guaranteed due to compiler optimizations.
func getCalleInfo(depth int) (pkg string, fun string) {
	fpc, _, _, _ := runtime.Caller(depth)

	fname := runtime.FuncForPC(fpc).Name()

	fparts := strings.Split(fname, ".")

	pkg = strings.Join(fparts[:len(fparts)-2], ".")

	fun = fparts[len(fparts)-1]

	return pkg, fun
}
