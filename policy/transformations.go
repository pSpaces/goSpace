package policy

// Transformations defines a structure for a collection of transformations to be applied.
type Transformations struct {
	Tmpl Transformation
	Mtch Transformation
	Rslt Transformation
}

// NewTransformations creates a collection of transformations to be applied.
// NewTransformations returns a pointer to a collection if exactly 3 types of transformations are specified, otherwise nil is returned.
func NewTransformations(tr ...*Transformation) (trs *Transformations) {
	if len(tr) != 3 {
		trs = nil
	} else {
		trc := make([]Transformation, len(tr))

		// Make a copy of the transformations.
		for i := range tr {
			if tr[i] != nil {
				trans := tr[i]
				function := trans.Function()
				params := trans.Parameters()
				ntr := NewTransformation(function, params...)
				trc[i] = ntr
			} else {
				trc[i] = Transformation{}
			}
		}

		trs = &Transformations{Tmpl: trc[0], Mtch: trc[1], Rslt: trc[2]}
	}

	return trs
}

// Template returns a transformation that can be applied to template entities.
func (trs *Transformations) Template() (trans *Transformation) {
	trans = nil

	if trs != nil {
		trans = &trs.Tmpl
	}

	return trans
}

// Match returns an transformation that can be applied to matched entities.
func (trs *Transformations) Match() (match *Transformation) {
	match = nil

	if trs != nil {
		match = &trs.Mtch
	}

	return match
}

// Result returns an transformation that can be applied to result entities.
func (trs *Transformations) Result() (rslt *Transformation) {
	rslt = nil

	if trs != nil {
		rslt = &trs.Rslt
	}

	return rslt
}
