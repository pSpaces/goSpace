package container

// Intertemplate an interface for manipulating templates.
type Intertemplate interface {
	Applier
	Length() int
	Fields() []interface{}
	GetFieldAt(i int) interface{}
	NewTuple() Tuple
}
