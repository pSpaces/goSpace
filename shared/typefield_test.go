package shared

import (
	"reflect"
	"testing"
)

// Test to see if typeField is creating correctly.
func TestCreateTypeField(t *testing.T) {
	actualField := "string"
	testField := reflect.TypeOf("")

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
	actualField := reflect.TypeOf("")
	testTypeField := CreateTypeField(actualField)

	testField := testTypeField.GetType()

	fieldsEqual := reflect.DeepEqual(actualField, testField)

	if !fieldsEqual {
		t.Errorf("GetType() gave %+v, should be %+v", testField, actualField)
	}
}
