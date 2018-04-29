package container

// Applier interface describes for map functionality for collections.
type Applier interface {
	// Applier interface iterates over values and applies function fun on each value.
	// Applier returns true function fun could be applied to each element, and false otherwise.
	Apply(fun func(interface{}) interface{}) bool
}
