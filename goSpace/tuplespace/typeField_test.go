package tuplespace

import (
	"reflect"
	"testing"
)

// Test to see if typeField is creating correctly.
func TestCreateTypeField(t *testing.T) {
	actualField := "String"
	testField := "String"

	actualTypeField := TypeField{actualField}
	testTypeField := CreateTypeField(testField)

	// Test that the two templates are equal.
	typeFieldsEqual := reflect.DeepEqual(actualTypeField, testTypeField)

	if !typeFieldsEqual {
		t.Errorf("CreateTypeField() gave %+v, should be %+v", testTypeField, actualTypeField)
	}
}

func TestGetType(t *testing.T) {
	// Setup
	actualField := "String"
	testTypeField := CreateTypeField(actualField)

	// Use getType()
	testField := testTypeField.getType()

	fieldsEqual := reflect.DeepEqual(actualField, testField)

	if !fieldsEqual {
		t.Errorf("getType() gave %+v, should be %+v", testField, actualField)
	}
}
