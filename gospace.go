package gospace

import (
	"github.com/luhac/gospace/protocol"
	"github.com/luhac/gospace/space"
)

type PointToPoint = protocol.PointToPoint

func NewSpace(name string) PointToPoint {
	return space.NewSpace(name)
}

type SpaceInterface interface {
	space.SpaceInterface
}
