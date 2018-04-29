package policy

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/pspaces/gospace/container"
	"github.com/pspaces/gospace/function"
)

// Action is a structure defining an operation.
type Action struct {
	Oper   string
	Func   interface{}
	Params []interface{}
	Sign   ActionSignature
}

// ActionSignature is a structure for storing signature information for faster matching.
type ActionSignature struct {
	Oper   container.Signature
	Func   container.Signature
	Params container.Signature
}

// NewAction creates an action given a function and optionally a parameter list params.
func NewAction(fun interface{}, params ...interface{}) (a *Action) {
	operator := function.Name(fun)

	sign := ActionSignature{
		Oper:   container.NewSignature(1, operator),
		Func:   container.NewSignature(1, fun),
		Params: container.NewTypeSignature(math.MaxUint16, params),
	}

	a = &Action{Oper: operator, Func: fun, Params: params, Sign: sign}

	return a
}

// Equal returns true if both action a and action b are equivalent, and false otherwise.
func (a *Action) Equal(b *Action) (e bool) {
	e = false

	if a == nil && a == b {
		e = true
	} else if a != nil && b != nil {
		op := a.Operator() == b.Operator()

		signature := false
		aSign := a.Signature()
		bSign := b.Signature()
		if op {
			signature = aSign.Oper == bSign.Oper && aSign.Func == bSign.Func
		}

		fun := false
		if signature {
			fun = reflect.ValueOf(a.Function()).Pointer() == reflect.ValueOf(b.Function()).Pointer()
		}

		params := len(a.Parameters()) == len(b.Parameters())
		if fun && params {
			if aSign.Params != bSign.Params {
				ta := container.NewTemplate(a.Parameters()...)
				tb := container.NewTemplate(b.Parameters()...)
				params = params && (&ta).Equal(&tb)
			}
		}

		e = op && signature && fun && params
	}

	return e
}

// Operator returns the operator name s of the action a.
func (a *Action) Operator() (s string) {
	if a != nil {
		s = (*a).Oper
	}

	return s
}

// Function returns the function f associated to the action a.
func (a *Action) Function() (f interface{}) {
	if a != nil {
		f = (*a).Func
	}

	return f
}

// Parameters returns the actual paramaters p which optionally can be applied to action a.
func (a *Action) Parameters() (p []interface{}) {
	if a != nil {
		p = (*a).Params
	}

	return p
}

// Signature returns the signature s belonging to an action a.
func (a *Action) Signature() (s ActionSignature) {
	if a != nil {
		s = (*a).Sign
	}

	return s
}

// String returns print friendly representation of an action a.
func (a Action) String() (s string) {
	s = fmt.Sprintf("{\n\t%v,\n\t%v,\n\t%v,\n\t%v,\n}", a.Oper, a.Func, a.Params, strings.Replace(a.Sign.String(), "\n", "\n\t", -1))
	return s
}

// String returns print friendly representation of an action signature as.
func (as ActionSignature) String() (s string) {
	s = fmt.Sprintf("{\n\t%v,\n\t%v,\n\t%v,\n}", as.Oper, as.Func, as.Params)
	return s
}
