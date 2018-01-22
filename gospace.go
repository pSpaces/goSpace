package gospace

import (
	"crypto/tls"

	shr "github.com/pspaces/gospace/shared"
	spc "github.com/pspaces/gospace/space"
)

type Space = spc.Space
type Tuple = shr.Tuple
type Template = shr.Template

// NewSpace creates a structure that represents a space.
func NewSpace(name string, config *tls.Config) Space {
	return spc.NewSpace(name, config)
}

// NewSpace creates a structure that represents a remote space.
func NewRemoteSpace(name string, config *tls.Config) Space {
	return spc.NewRemoteSpace(name, config)
}

// SpaceInterface contains all interfaces that can operate on a space.
type SpaceInterface interface {
	spc.Interspace
}

// CreateTuple creates a structure that represents a tuple.
func CreateTuple(fields ...interface{}) Tuple {
	return shr.CreateTuple(fields...)
}

// TupleInterface contains all interfaces that can operate on a tuple.
type TupleInterface interface {
	shr.Intertuple
}

// CreateTemplate creates a structure that represents a template.
func CreateTemplate(fields ...interface{}) Template {
	return shr.CreateTemplate(fields...)
}

// TemplateInterface contains all interfaces that can operate on a template.
type TemplateInterface interface {
	shr.Intertemplate
}
