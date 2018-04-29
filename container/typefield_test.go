package container

import (
	"reflect"
	"testing"
)

// Test to see if typeField is creating correctly.
func TestCreateTypeField(t *testing.T) {
	actualField := "string"
	testField := ""

	actualTypeField := TypeField{actualField}
	testTypeField := CreateTypeField(reflect.ValueOf(&testField).Elem().Interface())

	// Test that the two type fields are equal.
	typeFieldsEqual := reflect.DeepEqual(actualTypeField, testTypeField)

	if !typeFieldsEqual {
		t.Errorf("CreateTypeField() gave %+v, should be %+v", testTypeField, actualTypeField)
	}
}

func TestGetType(t *testing.T) {
	// Setup
	actualField := ""
	testTypeField := CreateTypeField(reflect.ValueOf(&actualField).Elem().Interface())

	testField := testTypeField.GetType()

	fieldsEqual := reflect.DeepEqual(reflect.TypeOf(actualField), testField)

	if !fieldsEqual {
		t.Errorf("GetType() gave %+v, should be %+v", testField, actualField)
	}
}
