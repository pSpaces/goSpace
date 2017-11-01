package gospace

import (
	"github.com/pspaces/gospace/container"
	"github.com/pspaces/gospace/policy"
	"github.com/pspaces/gospace/space"
)

// Intertuple is an interface for interacting with a tuple.
type Intertuple = container.Intertuple

// Intertemplate is an interface for interacting with a template.
type Intertemplate = container.Intertemplate

// Space defines a multi-set for tuples.
type Space = space.Space

// Tuple defines a tuple structure.
type Tuple = container.Tuple

// Template defines a template used for pattern matching.
type Template = container.Template

// Label describes a label given some entity.
type Label = container.Label

// Labels defines a set of labels.
type Labels = container.Labels

// LabelledTuple represents a labelled tuple.
type LabelledTuple = container.LabelledTuple

// Action encapsulates an operation with its parameters.
type Action = policy.Action

// AggregationRule defines an aggregation rule.
type AggregationRule = policy.AggregationRule

// AggregationPolicy defines an aggregation policy.
type AggregationPolicy = policy.Aggregation

// ComposablePolicy defines a composable policy.
type ComposablePolicy = policy.Composable

// Transformation defines a transformation.
type Transformation = policy.Transformation

// Transformations defines a transformations applied to an action.
type Transformations = policy.Transformations

// NewSpace creates a structure that represents a space.
func NewSpace(name string, policy ...*ComposablePolicy) Space {
	return space.NewSpace(name, policy...)
}

// NewRemoteSpace creates a structure that represents a remote space.
func NewRemoteSpace(name string) Space {
	return space.NewRemoteSpace(name)
}

// SpaceFrame contains all interfaces that can operate on a space.
type SpaceFrame interface {
	space.Interspace
	space.Interstar
}

// CreateTuple creates a structure that represents a tuple.
func CreateTuple(fields ...interface{}) Tuple {
	return container.NewTuple(fields...)
}

// TupleFrame contains all interfaces that can operate on a tuple.
type TupleFrame interface {
	Intertuple
}

// CreateTemplate creates a structure that represents a template.
func CreateTemplate(fields ...interface{}) Template {
	return container.NewTemplate(fields...)
}

// TemplateFrame contains all interfaces that can operate on a template.
type TemplateFrame interface {
	Intertemplate
}

// NewLabel creates a structure that represents a label.
func NewLabel(id string) Label {
	return container.NewLabel(id)
}

// NewLabels creates a structure that represents a collection of labels.
func NewLabels(ll ...Label) Labels {
	return container.NewLabels(ll...)
}

// NewLabelledTuple creates a structure that represents a labelled tuple.
func NewLabelledTuple(fields ...interface{}) LabelledTuple {
	return container.NewLabelledTuple(fields...)
}

// NewAction creates a structure that represents an action.
func NewAction(function interface{}, params ...interface{}) *Action {
	return policy.NewAction(function, params...)
}

// NewAggregationRule creates a structure that represents an aggregation rule.
func NewAggregationRule(a Action, trs Transformations) AggregationRule {
	return policy.NewAggregationRule(a, trs)
}

// NewAggregationPolicy creates a structure that represents an aggregation policy.
func NewAggregationPolicy(l Label, r AggregationRule) AggregationPolicy {
	return policy.NewAggregation(l, r)
}

// NewComposablePolicy creates a structure that represents a composable policy.
func NewComposablePolicy(ars ...AggregationPolicy) *ComposablePolicy {
	return policy.NewComposable(ars...)
}

// NewTransformation creates a structure for representing a transformation that can be applied to an action.
func NewTransformation(function interface{}, params ...interface{}) Transformation {
	return policy.NewTransformation(function, params...)
}

// NewTransformations creates a structure for representing collection of transformation that can be applied to an action.
func NewTransformations(trs ...*Transformation) *Transformations {
	return policy.NewTransformations(trs...)
}
