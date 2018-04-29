package policy

import (
	"github.com/pspaces/gospace/container"
)

// Aggregation is a structure defining how an aggregation operation should be transformed by an aggregation rule.
type Aggregation struct {
	Lbl     container.Label
	AggRule AggregationRule
}

// NewAggregation constructs a new policy given a label l an aggregation rule r.
func NewAggregation(l container.Label, r AggregationRule) (ap Aggregation) {
	ap = Aggregation{Lbl: l, AggRule: r}
	return ap
}

// Label returns the label l attached to aggregation policy ap.
func (ap *Aggregation) Label() (l container.Label) {
	l = (*ap).Lbl
	return l
}

// Action returns the acction associated to the aggregation rule.
func (ap *Aggregation) Action() (a Action) {
	return ap.AggRule.Object
}

// Apply applies an aggregation policy onto the input action ia.
// Apply returns a modified action oa
func (*Aggregation) Apply(ia Action) (oa Action) {
	return oa
}
