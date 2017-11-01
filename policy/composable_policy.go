package policy

import (
	"fmt"
	"strings"
	"sync"

	"github.com/deckarep/golang-set"
	"github.com/pspaces/gospace/container"
)

// Composable is a structure for containing composable policies.
type Composable struct {
	ActionMap *sync.Map // [ActionTypeSignature]Set(Template,Label)
	LabelMap  *sync.Map // [Label]AggregationPolicy
}

// NewComposable creates a composable policy cp from any amount of aggregation policies aps.
func NewComposable(aps ...Aggregation) (cp *Composable) {
	cp = &Composable{ActionMap: new(sync.Map), LabelMap: new(sync.Map)}

	for _, ap := range aps {
		cp.Add(ap)
	}

	return cp
}

// Add adds an aggregation policy ap to the composable policy cp.
// Add returns true if the aggregation policy ap has been added to the composable policy cp, and false otherwise.
func (cp *Composable) Add(ap Aggregation) (b bool) {
	b = cp != nil

	var crit = &sync.Mutex{}

	if b {
		a := (&ap).Action()
		l := (&ap).Label()
		atemp := container.NewTemplate(a.Parameters()...)
		lid := l.ID()

		as := (&a).Signature()
		isetp, exists := cp.ActionMap.Load(as)
		set, valid := isetp.(*mapset.Set)

		// Check for the existence of the policy in the set.
		polExist := false
		if exists && valid {
			ml := FindLabel(set, &atemp)
			polExist = ml != nil && (*ml).ID() == lid
		}

		// Create a new set for (Template, Label)-pairs.
		if !valid {
			instance := mapset.NewSet()
			polPair := container.NewTuple(atemp, l)
			instance.Add(&polPair)
			set = &instance
		}

		b = !polExist
		if b {
			_, existsPolicy := cp.LabelMap.Load(lid)

			b = !existsPolicy
			if b {
				// Check that the policy does not exists.
				// Otherwise, perform a union and put it back.
				crit.Lock()

				isetp, exists = cp.ActionMap.Load(as)
				setp, valid := isetp.(*mapset.Set)
				if exists && valid {
					instance := (*set).Union(*setp)
					set = &instance
				}
				cp.ActionMap.Store(as, set)
				_, existsPolicy := cp.LabelMap.LoadOrStore(lid, ap)

				crit.Unlock()

				b = !existsPolicy
			}
		}
	}

	return b
}

// FindLabel checks for a set s, a template tp has an associated a label l.
// FindLabel returns l if it is associated to a template tp, and nil otherwise.
func FindLabel(s *mapset.Set, tp *container.Template) (l *container.Label) {
	b := tp != nil

	if b {
		var mpno, mqno uint

		sit := (*s).Iterator()
		for elem := range sit.C {

			ituple := elem.(container.Intertuple)
			temp, et := ituple.GetFieldAt(0).(container.Template)
			label, el := ituple.GetFieldAt(1).(container.Label)

			if !et || !el {
				continue
			}

			match, pno, qno := tp.ExactMatch(&temp)

			if match {
				if qno > mqno {
					if pno > mpno {
						mpno = pno
						lc := label.DeepCopy()
						l = &lc
					}
					mqno = qno
				}
			}
		}
	}

	return l
}

// Find returns a reference to a label l given an action a.
func (cp *Composable) Find(a *Action) (l *container.Label) {
	b := cp != nil
	l = nil

	if b {
		as := a.Signature()
		atemp := container.NewTemplate(a.Parameters()...)
		iset, exists := cp.ActionMap.Load(as)
		set, valid := iset.(*mapset.Set)

		if exists && valid {
			lbl := FindLabel(set, &atemp)
			l = lbl
		}
	}

	return l
}

// Retrieve returns a reference to the aggregation policy ap with label l from the composable policy cp.
// Retrieve returns a reference if it exists and nil otherwise.
func (cp *Composable) Retrieve(l container.Label) (ap *Aggregation) {
	b := cp != nil
	ap = nil

	if b {
		lid := l.ID()
		val, exists := cp.LabelMap.Load(lid)
		if exists {
			pol := val.(Aggregation)
			ap = &pol
		}
	}

	return ap
}

// Delete removes an aggregation policy with label l from the composable policy cp.
// Delete returns true if an aggregation policy with label l has been deleted from the composable policy cp, and false otherwise.
func (cp *Composable) Delete(l container.Label) (b bool) {
	b = cp != nil

	if b {
		lid := l.ID()
		val, exists := cp.LabelMap.Load(lid)
		if exists {
			ap := val.(Aggregation)
			a := ap.Action()
			cp.LabelMap.Delete(l)
			cp.ActionMap.Delete(a)
		}

		b = exists
	}

	return b
}

// String returns a print friendly representation of a composable policy cp.
func (cp Composable) String() (s string) {
	var actionEntries, labelEntries []string

	entries := []string{}
	entry := make(chan string)

	go func() {
		cp.ActionMap.Range(func(k, v interface{}) bool {
			entry <- fmt.Sprintf("\t%v: %v", k, v)
			return true
		})
		close(entry)
	}()

	for entry := range entry {
		entries = append(entries, entry)
	}

	actionEntries = entries

	entries = []string{}
	entry = make(chan string)

	go func() {
		cp.LabelMap.Range(func(k, v interface{}) bool {
			entry <- fmt.Sprintf("\t%v: %v", k, v)
			return true
		})
		close(entry)
	}()

	for entry := range entry {
		entries = append(entries, entry)
	}

	labelEntries = entries

	refs := strings.Join(actionEntries, ",\n")

	if refs != "" {
		refs = fmt.Sprintf("\n%s\n", refs)
	}

	names := strings.Join(labelEntries, ",\n")

	if names != "" {
		names = fmt.Sprintf("\n%s\n", names)
	}

	s = fmt.Sprintf("%s%s%s%s%s%s", "{", refs, "}", "{", names, "}")

	return s
}
