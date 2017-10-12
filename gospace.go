package gospace

import (
	shr "github.com/pspaces/gospace/shared"
	spc "github.com/pspaces/gospace/space"
)

type Space = spc.Space
type Tuple = shr.Tuple
type Template = shr.Template

func NewSpace(name string) Space {
	return spc.NewSpace(name)
}

func NewRemoteSpace(name string) Space {
	return spc.NewRemoteSpace(name)
}

// SpaceInterface contains all interfaces that can operate on a space.
type SpaceInterface interface {
	spc.Interspace
}

func CreateTuple(fields []interface{}) Tuple {
	return shr.CreateTuple(fields)
}

// TupleInterface contains all interfaces that can operate on a tuple.
type TupleInterface interface {
	shr.Intertuple
}

func CreateTemplate(fields []interface{}) Template {
	return shr.CreateTemplate(fields)
}

// TemplateInterface contains all interfaces that can operate on a template.
type TemplateInterface interface {
	shr.Intertemplate
}
