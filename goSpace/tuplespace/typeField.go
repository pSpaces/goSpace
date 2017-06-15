package tuplespace

type TypeField struct {
	IsType string // Field of the tuple.
}

func CreateTypeField(isType string) TypeField {
	typeField := TypeField{isType}
	return typeField
}

func (t TypeField) getType() string {
	return t.IsType
}
