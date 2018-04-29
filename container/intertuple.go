package container

// Intertuple defines an interface for manipulating tuples.
type Intertuple interface {
	Applier
	Tuple() Tuple
	Length() int
	Fields() []interface{}
	GetFieldAt(i int) interface{}
	SetFieldAt(i int, val interface{}) bool
}
