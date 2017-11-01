package container

import (
	"fmt"
	"sort"
	"strings"
)

// Labels is structure used to represent a set of labels.
// Labels purpose is to conveniently manipulate many labels.
type Labels map[string]Label

// NewLabels creates a new collections of labels from parameter list ll.
func NewLabels(ll ...Label) (ls Labels) {
	ls = make(Labels)

	for _, l := range ll {
		ls.Add(l.DeepCopy())
	}

	return ls
}

// Add adds a label l to label set ls.
// Add returns true if the label has been added, and false otherwise.
func (ls *Labels) Add(l Label) (b bool) {
	_, exists := (*ls)[l.ID()]
	if !exists {
		(*ls)[l.ID()] = l
	}

	b = !exists

	return b
}

// Retrieve returns a label l in the label set ls.
// Retrieve returns a reference to the label if it exists and nil otherwise.
func (ls *Labels) Retrieve(id string) (l *Label) {
	lbl, exists := (*ls)[id]
	if exists {
		l = &lbl
	} else {
		l = nil
	}

	return l
}

// Intersect returns true if both label sets ls and lm intersect, and false otherwise.
// Intersect returns an intersected label set li if ls and lm intersect, and nil otherwise.
func (ls *Labels) Intersect(lm *Labels) (li *Labels, e bool) {
	e = ls != nil && lm != nil

	if e {
		labelling := ls.Labelling()
		intersect := make([]Label, 0, len(labelling))
		i := 0
		for _, id := range labelling {
			la := ls.Retrieve(id)
			lb := lm.Retrieve(id)
			if lb != nil {
				e = la.Equal(lb)
				if e {
					intersect = append(intersect, NewLabel(id))
					i++
				}
			}
		}

		e = len(intersect) > 0
		if e {
			lbls := NewLabels(intersect[:i]...)
			li = &lbls
		}
	}

	return li, e
}

// Delete deletes a label l from label set ls.
// Delete returns true if the label has been deleted, and false otherwise.
func (ls *Labels) Delete(id string) (b bool) {
	_, exists := (*ls)[id]
	if exists {
		delete(*ls, id)
	}

	b = exists

	return b
}

// Labelling returns the labels identifiers present in the label set ls.
func (ls *Labels) Labelling() (labelling []string) {
	if ls != nil {
		labelling = make([]string, len(*ls))

		i := 0
		for k := range *ls {
			labelling[i] = k
			i++
		}
	}

	return labelling
}

// Set returns the set of all labels in ls.
func (ls *Labels) Set() (set []Label) {
	if ls != nil {
		set = make([]Label, 0, len(*ls))

		for _, v := range *ls {
			set = append(set, v)
		}
	}

	return set
}

// String returns a print friendly representation of the labels set ls.
func (ls Labels) String() (s string) {
	ms := make([]string, len((&ls).Labelling()))

	for i, lid := range (&ls).Labelling() {
		l := (&ls).Retrieve(lid)
		ms[i] = fmt.Sprintf("%s", *l)
	}

	sort.Strings(ms)

	s = fmt.Sprintf("{%s}", strings.Join(ms, ", "))

	return s
}
