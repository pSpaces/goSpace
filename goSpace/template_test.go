package goSpace

import (
	"reflect"
	"testing"
)

// Test to see if template is creating correct.
func TestCreateTemplate(t *testing.T) {
	// Setup
	actualFields := make([]interface{}, 5)
	actualFields[0] = "Field 1"
	actualFields[1] = 2
	actualFields[2] = 3.14
	actualFields[3] = false
	actualFields[4] = CreateTypeField("int")

	testTemplate := createTestTemplate()
	actualTemplate := Template{actualFields}

	// Test that the two templates are equal.
	templatesEqual := reflect.DeepEqual(actualTemplate, testTemplate)

	if !templatesEqual {
		t.Errorf("CreateTemplate() gave %+v, should be %+v", testTemplate, actualTemplate)
	}
}

// Test to see if template has correct length.
func TestTemplateLength(t *testing.T) {
	// Setup
	testTemplate := createTestTemplate()

	// Get lengths
	actualTemplateLength := 5
	testTemplateLength := testTemplate.length()

	if testTemplateLength != actualTemplateLength {
		t.Errorf("Length(%+v) == %d, should be %d", testTemplate, testTemplateLength, actualTemplateLength)
	}
}

// Test to see if the correct template field is returned.
func TestTemplateGetFieldAt(t *testing.T) {
	// Setup
	testTemplate := createTestTemplate()

	actualFieldAtBool := false
	testFieldAtBool := testTemplate.getFieldAt(3)

	if actualFieldAtBool != testFieldAtBool {
		t.Errorf("GetFieldAt(%d) on template: %+v == %v, should be %v", 3, testTemplate, testFieldAtBool, actualFieldAtBool)
	}

	actualFieldAtPtr := CreateTypeField("int")
	testFieldAtPtr := testTemplate.getFieldAt(4)

	fieldsEqual := reflect.DeepEqual(actualFieldAtPtr, testFieldAtPtr)

	if !fieldsEqual {
		t.Errorf("GetFieldAt(%d) on template: %+v == %v, should be %v", 3, testTemplate, testFieldAtPtr, actualFieldAtPtr)
	}
}

func createTestTemplate() Template {
	testFields := make([]interface{}, 5)
	testFields[0] = "Field 1"
	testFields[1] = 2
	testFields[2] = 3.14
	testFields[3] = false
	intVal := 2
	intPtr := &intVal
	testFields[4] = intPtr

	return CreateTemplate(testFields)
}
