package space

import (
	"fmt"
	"go/build"
	"reflect"
	"strings"

	"github.com/pspaces/gospace/function"
)

// SpaceError represents an internal error type used when printing error messages.
type SpaceError struct {
	Msg     string
	LibInfo function.CallerInfo
	UsrInfo function.CallerInfo
	Sid     string
	Val     string
	Sop     bool
	Status  interface{}
}

// Constants used for enumerating the generic error strings.
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

// Constants used for printing a partial trace.
const (
	libCallDepth = 3
	usrCallDepth = 4
)

// NewSpaceError creates a new error given space spc, a value used in an operation and the return state of the implemented operation.
// NewSpaceError returns a structure which fulfils the error interface and if an operation error has occured.
// NewSpaceError returns nil if no operation failure has occured.
func NewSpaceError(spc *Space, value interface{}, state interface{}) error {
	var msg, sid, val string
	var err error
	var sop bool
	var libInfo, usrInfo function.CallerInfo
	var status interface{}

	libInfo = function.ExtractCallerInfo(libCallDepth)
	usrInfo = function.ExtractCallerInfo(usrCallDepth)

	if spc == nil {
		sid = "nil"
		msg = errMsg[SpaceInvalid]
		status = nil
	} else if state != nil {
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

		method = spct.MethodByName("InterpretValue")
		if method.IsValid() {
			vals := method.Call([]reflect.Value{reflect.ValueOf(value)})
			if len(vals) == 1 {
				val = vals[0].Interface().(string)
			} else {
				val = ""
			}
		} else {
			val = ""
		}

		status = state
	}

	if sop == true {
		err = nil
	} else {
		if state != nil {
			err = SpaceError{Msg: msg, LibInfo: libInfo, UsrInfo: usrInfo, Sid: sid, Val: val, Sop: sop, Status: status}
		}
	}

	return err
}

// Operation returns a boolean value if an operation has succeeded.
func (e SpaceError) Operation() bool {
	return e.Sop
}

// Error prints the error message s represented by SpaceError e.
func (e SpaceError) Error() (s string) {
	separator := strings.Repeat(" ", 2)
	libFile := strings.Replace(e.LibInfo.File, strings.Join([]string{build.Default.GOPATH, "/src/"}, ""), "", 1)
	usrFile := strings.Replace(e.UsrInfo.File, strings.Join([]string{build.Default.GOPATH, "/src/"}, ""), "", 1)
	libInfo := fmt.Sprintf("%s:%d", libFile, e.LibInfo.Line)
	usrInfo := fmt.Sprintf("%s:%d", usrFile, e.UsrInfo.Line)
	call := fmt.Sprintf("%s(%s).%s%s: %s", "Space", e.Sid, e.UsrInfo.Func, e.Val, e.Msg)
	s = fmt.Sprintf("\n%s%s:\n%s%s:\n%s%s%s.", separator, libInfo, separator, usrInfo, separator, separator, call)
	return s
}
