package policy

import (
	"fmt"
	"github.com/choleraehyq/gofunctools/functools"
)

// Transformation defines a structure for a transformation.
type Transformation struct {
	Func   interface{}
	Params []interface{}
}

// TransformationError defines an error occuring during a application of a transformation.
type TransformationError struct {
	msg string
	tr  Transformation
}

// NewTransformation creates a new transformation from a function and optional parameter list params.
func NewTransformation(fun interface{}, params ...interface{}) (tr Transformation) {
	tr = Transformation{Func: fun, Params: params}
	return tr
}

// Apply applies attached transformation and additionally passed parameters fparams.
// Apply returns an encapsulated result res and an error err if function application failed.
func (tr *Transformation) Apply(fparams ...interface{}) (res interface{}, err error) {
	err = nil

	if tr != nil {
		fun := tr.Function()

		if fun != nil {
			tparams := tr.Parameters()

			if tparams != nil {
				pfun, e := functools.Partial(fun, tparams...)
				if e == nil {
					pres, pe := functools.Apply(pfun, []interface{}{fparams})
					res = pres.([]interface{})[0]
					err = pe
				} else {
					err = e
				}
			} else {
				pres, e := functools.Apply(fun, []interface{}{fparams})
				res = pres.([]interface{})[0]
				err = e
			}
		} else {
			err = &TransformationError{"Transformation could not be applied.", *tr}
		}
	}

	return res, err
}

// Error returns an error msg associated to the transformation tr.
func (err *TransformationError) Error() (msg string) {
	msg = fmt.Sprintf("%s: %v", err.msg, err.tr)
	return msg
}

// Function returns the function associated to the transformation tr.
func (tr *Transformation) Function() (f interface{}) {
	if tr != nil {
		f = tr.Func
	}

	return f
}

// Parameters returns the parameters associated to the transformation tr.
func (tr *Transformation) Parameters() (p []interface{}) {
	if tr != nil {
		p = tr.Params
	}

	return p
}

// String returns print friendly representation of a transformation tr.
func (tr Transformation) String() (s string) {
	return fmt.Sprintf("{%v, %v}", tr.Func, tr.Params)
}
