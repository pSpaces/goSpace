package function

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// CallerInfo is a structure for containing function caller information.
type CallerInfo struct {
	File string
	Line int
	Pkg  string
	Func string
}

// ExtractCallerInfo returns caller info associated to a function call given a depth.
// ExtractCallerInfo returns unknown caller info if frame unwinding fails.
func ExtractCallerInfo(depth int) (ci CallerInfo) {
	pc := make([]uintptr, depth)
	n := runtime.Callers(1, pc)

	var file, pkg, fun string
	var line int

	if n != 0 {
		pc = pc[:n]
		frames := runtime.CallersFrames(pc)

		// Unwind the frames.
		frame, available := frames.Next()
		for i := 0; i < depth; i++ {
			if available {
				frame, available = frames.Next()
			} else {
				break
			}
		}

		file = frame.File

		line = frame.Line

		fparts := strings.Split(frame.Function, ".")

		pkg = strings.Join(fparts[:len(fparts)-1], ".")

		fun = fparts[len(fparts)-1]
	} else {
		file = "<unknown file>"
		line = -1
		pkg = "<unknown package>"
		fun = "<unknown function>"
	}

	ci = CallerInfo{File: file, Line: line, Pkg: pkg, Func: fun}

	return ci
}

// IsFunc returns true if interface fun is a function or method and false otherwise.
func IsFunc(fun interface{}) (b bool) {
	var fn reflect.Value
	if reflect.TypeOf(fun) == reflect.TypeOf(reflect.Value{}) {
		fn = fun.(reflect.Value)
	} else if reflect.TypeOf(fun) == reflect.TypeOf(reflect.Method{}) {
		fn = (fun.(reflect.Method)).Func
	} else {
		fn = reflect.ValueOf(fun)
	}

	b = fn.Type().Kind() == reflect.Func

	return b
}

// Name returns name of the function as a string s.
// Name will panic if function is not a function or method.
func Name(fun interface{}) (s string) {
	var fn reflect.Value
	if reflect.TypeOf(fun) == reflect.TypeOf(reflect.Value{}) {
		fn = fun.(reflect.Value)
	} else if reflect.TypeOf(fun) == reflect.TypeOf(reflect.Method{}) {
		fn = (fun.(reflect.Method)).Func
	} else {
		fn = reflect.ValueOf(fun)
	}

	if fn.Type().Kind() != reflect.Func {
		panic("Can not determine function name of non-function value.")
	}

	fp := fn.Pointer()

	s = strings.Replace(runtime.FuncForPC(fp).Name(), "-fm", "", 1)

	return s
}

// Signature returns the function signature of any function or method.
// Signature accepts an optional name parameter used in the function signature.
// Signature will panic if function is not a function or method.
func Signature(fun interface{}, name ...string) (sgn string) {
	var fn reflect.Value
	if reflect.TypeOf(fun) == reflect.TypeOf(reflect.Value{}) {
		fn = fun.(reflect.Value)
	} else {
		fn = reflect.ValueOf(fun)
	}

	if fn.Type().Kind() != reflect.Func {
		panic("Can not determine signature of non-function value.")
	}

	isize := fn.Type().NumIn()
	osize := fn.Type().NumOut()
	vsize := 0

	if fn.Type().IsVariadic() {
		vsize++
	}

	istrs := make([]string, isize)

	for i := 0; i < isize; i++ {
		t := fn.Type().In(i)
		istrs[i] = strings.Replace(t.String(), " ", "", 1)

		if i == (isize - 1) {
			if vsize > 0 {
				istrs[i] = strings.Replace(istrs[i], "[]", "...", 1)
			}
		}
	}

	istr := fmt.Sprintf("%s%s%s", "(", strings.Join(istrs, ", "), ")")

	to := "->"

	ostrs := make([]string, osize)

	for j := 0; j < osize; j++ {
		t := fn.Type().Out(j)
		ostrs[j] = strings.Replace(t.String(), " ", "", 1)
	}

	ostr := fmt.Sprintf("%s%s%s", "(", strings.Join(ostrs, ", "), ")")

	sgn = strings.Join([]string{istr, to, ostr}, " ")

	if len(name) == 1 {
		sgn = fmt.Sprintf("%s: %s", name[0], sgn)
	}

	return sgn
}

// Type returns the type of a function or method.
func Type(fun interface{}) (t reflect.Type) {
	var fn reflect.Value
	if reflect.TypeOf(fun) == reflect.TypeOf(reflect.Value{}) {
		fn = fun.(reflect.Value)
	} else if reflect.TypeOf(fun) == reflect.TypeOf(reflect.Method{}) {
		fn = (fun.(reflect.Method)).Func
	} else {
		fn = reflect.ValueOf(fun)
	}

	t = fn.Type()

	return t
}
