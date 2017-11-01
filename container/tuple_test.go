package container

import (
	"reflect"
	"testing"
)

// Test to see if tuple is creating correct.
func TestCreateTuple(t *testing.T) {
	// Setup
	testTuple := createTestTuple()

	// Create tuple with different kind of types, manually.
	actualFields := make([]interface{}, 5)
	actualFields[0] = "Field 1"
	actualFields[1] = 2
	actualFields[2] = 3.14
	actualFields[3] = false
	actualFields[4] = 991
	actualTuple := Tuple{actualFields}

	// Test the manually created tuple against the tuple created by the method.
	if !reflect.DeepEqual(testTuple, actualTuple) {
		t.Errorf("CreateTuple() gave %q, should be %q", testTuple, actualTuple)
	}
}

// Test to see if tuple has correct length.
func TestTupleLength(t *testing.T) {
	// Setup
	testTuple := createTestTuple()

	// Create tuple with same length as defined by variable.
	actualTupleLength := 5
	testTupleLength := testTuple.Length()

	if testTupleLength != actualTupleLength {
		t.Errorf("Length(%q) == %q, should be %q", testTuple, testTupleLength, actualTupleLength)
	}
}

// Test to see if the correct tuple field is returned.
func TestTupleGetFieldAt(t *testing.T) {
	// Setup
	testTuple := createTestTuple()

	actualFieldAtString := "Field 1"
	testFieldAtString := testTuple.GetFieldAt(0)

	if testFieldAtString != actualFieldAtString {
		t.Errorf("GetFieldAt(%q) == %q, should be %q", 0, testFieldAtString, actualFieldAtString)
	}
}

// Test to see if the correct tuple field is set.
func TestSetFieldAt(t *testing.T) {
	// Setup
	testTuple := createTestTuple()

	actualSetField := "Field new"

	testTuple.SetFieldAt(0, actualSetField)

	testSetFieldAt := testTuple.GetFieldAt(0)

	if testSetFieldAt != actualSetField {
		t.Errorf("SetFieldAt(%d, %q) == %q, should be %q", 0, actualSetField, testSetFieldAt, actualSetField)
	}
}

// Test to see if a tuple and template with different length wont match
func TestMatchNotSameLength(t *testing.T) {
	// Setup
	testTuple := createTestTuple()

	// Create template of different length
	testFields := make([]interface{}, 4)
	testFields[0] = "Field 1"
	testFields[1] = 2
	testFields[2] = 3.14
	testFields[3] = false
	testTemplateDifferentLength := NewTemplate(testFields...)

	if testTuple.Match(testTemplateDifferentLength) {
		t.Errorf("%q.match(%q) == false as length of %q = %d and length of %q = %d", testTuple, testTemplateDifferentLength, testTuple, testTuple.Length(), testTemplateDifferentLength, testTemplateDifferentLength.Length())
	}
}

// Test to see if a tuple and template with same length, where the template
// only contains typeFields that match.
func TestMatchSameLengthTemplateOnlyMatchingTypeFields(t *testing.T) {
	// Setup
	testTuple := createTestTuple()

	// Create template of same length with matching typeFields
	testFields := make([]interface{}, 5)
	stringVal := "1"
	stringPtr := &stringVal
	testFields[0] = stringPtr
	intVal := 2
	intPtr := &intVal
	testFields[1] = intPtr
	floatVal := 1.1
	floatPtr := &floatVal
	testFields[2] = floatPtr
	boolVal := true
	boolPtr := &boolVal
	testFields[3] = boolPtr
	testFields[4] = intPtr
	testTemplateSameLengthOnlyMatchingTypeFields := NewTemplate(testFields...)

	if !testTuple.Match(testTemplateSameLengthOnlyMatchingTypeFields) {
		t.Errorf("%q.match(%q) was expected to be true as the tuple looks like, %q, and the template looks like, %q.", testTuple, testTemplateSameLengthOnlyMatchingTypeFields, testTuple, testTemplateSameLengthOnlyMatchingTypeFields)
	}
}

// Test to see if a tuple and template with same length, where the template
// contains typeFields but they DON'T match.
func TestMatchSameLengthTemplateNotMatchingTypeFields(t *testing.T) {
	// Setup
	testTuple := createTestTuple()

	// Create template of same length with matching typeFields
	testFields := make([]interface{}, 5)
	stringVal := "1"
	stringPtr := &stringVal
	testFields[0] = stringPtr
	testFields[1] = stringPtr
	testFields[2] = stringPtr
	testFields[3] = stringPtr
	testFields[4] = stringPtr
	testTemplateSameLengthNotMatchingTypeFields := NewTemplate(testFields...)

	if testTuple.Match(testTemplateSameLengthNotMatchingTypeFields) {
		t.Errorf("%q.match(%q) was expected to be true as the tuple looks like, %q, and the template looks like, %q.", testTuple, testTemplateSameLengthNotMatchingTypeFields, testTuple, testTemplateSameLengthNotMatchingTypeFields)
	}
}

// Test to see if a tuple and template with same length, where the template
// only contains fields that match.
func TestMatchSameLengthTemplateOnlyMatchingFields(t *testing.T) {
	// Setup
	testTuple := createTestTuple()

	// Create template of same length with matching typeFields
	testFields := make([]interface{}, 5)
	testFields[0] = "Field 1"
	testFields[1] = 2
	testFields[2] = 3.14
	testFields[3] = false
	testFields[4] = 991
	testTemplateSameLengthOnlyMatchingFields := NewTemplate(testFields...)

	if !testTuple.Match(testTemplateSameLengthOnlyMatchingFields) {
		t.Errorf("%q.match(%q) was expected to be true as the tuple looks like, %q, and the template looks like, %q.", testTuple, testTemplateSameLengthOnlyMatchingFields, testTuple, testTemplateSameLengthOnlyMatchingFields)
	}
}

// Test to see if a tuple and template with same length, where the template
// contains fields but they DON'T match.
func TestMatchSameLengthTemplateNotMatchingFields(t *testing.T) {
	// Setup
	testTuple := createTestTuple()

	// Create template of same length with matching typeFields
	testFields := make([]interface{}, 5)
	testFields[0] = "Field 1"
	testFields[1] = 2
	testFields[2] = 3.14
	testFields[3] = false
	testFields[4] = 0
	testTemplateSameLengthNotMatchingFields := NewTemplate(testFields...)

	if testTuple.Match(testTemplateSameLengthNotMatchingFields) {
		t.Errorf("%q.match(%q) was expected to be true as the tuple looks like, %q, and the template looks like, %q.", testTuple, testTemplateSameLengthNotMatchingFields, testTuple, testTemplateSameLengthNotMatchingFields)
	}
}

func TestWriteToVariables(t *testing.T) {
	tuple := NewTuple([]interface{}{2, "hello"}...)
	var i int
	var s string
	variables := []interface{}{&i, &s}
	(&tuple).WriteToVariables(variables...)
	if i != 2 || s != "hello" {
		t.Errorf("Write tuple to variable did not work as expected")
	}
}

// Method used for setup.
func createTestTuple() Tuple {
	testFields := make([]interface{}, 5)
	testFields[0] = "Field 1"
	testFields[1] = 2
	testFields[2] = 3.14
	testFields[3] = false
	testFields[4] = 991

	return NewTuple(testFields...)
}
