package shared

// TypeField encapsulate a string that specifies the type.
type TypeField struct {
	IsType string // Field of the tuple.
}

// CreateTypeField will create the TypeField and return it.
func CreateTypeField(isType string) TypeField {
	typeField := TypeField{isType}
	return typeField
}

// getType will return the type of the TypeField.
func (t TypeField) getType() string {
	return t.IsType
}
