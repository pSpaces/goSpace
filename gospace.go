package gospace

import (
	shr "github.com/pspaces/gospace/shared"
	spc "github.com/pspaces/gospace/space"
)

type Space = spc.Space
type Tuple = shr.Tuple
type Template = shr.Template

// NewSpace creates a structure that represents a space.
func NewSpace(name string) Space {
	return spc.NewSpace(name)
}

// NewSpace creates a structure that represents a remote space.
func NewRemoteSpace(name string) Space {
	return spc.NewRemoteSpace(name)
}

// SpaceFrame contains all interfaces that can operate on a space.
type SpaceFrame interface {
	spc.Interspace
	spc.Interstar
}

// CreateTuple creates a structure that represents a tuple.
func CreateTuple(fields ...interface{}) Tuple {
	return shr.CreateTuple(fields...)
}

// TupleFrame contains all interfaces that can operate on a tuple.
type TupleFrame interface {
	shr.Intertuple
}

// CreateTemplate creates a structure that represents a template.
func CreateTemplate(fields ...interface{}) Template {
	return shr.CreateTemplate(fields...)
}

// TemplateFrame contains all interfaces that can operate on a template.
type TemplateFrame interface {
	shr.Intertemplate
}
