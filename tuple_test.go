package tuplespace

import (
	"reflect"
	"testing"
)

var testTupleTwoFields Tuple
var testTupleFourFields Tuple

func init() {
	// Create simple tuple with four fields.
	testTupleTwoFields = CreateTuple("Field 1", true)
	testTupleFourFields = CreateTuple("Field1", 2, 3.3, false)
}

// Test to see if tuple is creating correct.
func TestCreateTuple(t *testing.T) {
	// Create tuple with different kind of types, manually.
	fields := make([]interface{}, 4)
	fields[0] = "Field1"
	fields[1] = 2
	fields[2] = 3.3
	fields[3] = false
	var testTuple = Tuple{fields}

	// Test the manually created tuple against the tuple created by the method.
	if !reflect.DeepEqual(testTuple, testTupleFourFields) {
		t.Errorf("CreateTuple() gave %q, should be %q", testTupleFourFields, testTuple)
	}
}

// Test to see if tuple has correct length.
func TestTupleLength(t *testing.T) {
	// Create tuple with same length as defined by variable.
	var tupleLength = 4
	var testTuple = CreateTuple("Field1", 2, 3.3, false)
	testTupleLength := testTuple.Length()
	if tupleLength != testTupleLength {
		t.Errorf("Length(%q) == %q, should be %q", testTuple, testTupleLength, tupleLength)
	}
}

// Test to see if the correct tuple field is returned.
func TestTupleGetFieldAt(t *testing.T) {
	var wantedField = 2
	var testTuple = CreateTuple("Field1", 2, 3.3, false)
	tupleField := testTuple.GetFieldAt(1)
	if wantedField != tupleField {
		t.Errorf("GetFieldAt(%q) == %q, should be %q", testTuple, tupleField, wantedField)
	}
}

// Test to see if the correct tuple field is set.
func TestSetFieldAt(t *testing.T) {
	newField := "val"

	var testTuple = CreateTuple("Field1", 2, 3.3, false)
	testTuple.SetFieldAt(1, newField)

	newTupleField := testTuple.GetFieldAt(1)

	if newField != newTupleField {
		t.Errorf("SetFieldAt(%d, %q) == %q, should be %q", 1, newField, newTupleField, newField)
	}
}

func TestMatchSameLengthAndMatchingNumberOfFields(t *testing.T) {
	var testTemplate = CreateTemplate("Field1", 2, 3.3, false)
	var testTuple = CreateTuple("Field1", 2, 3.3, false)

	if !testTuple.match(testTemplate) {
		t.Errorf("TestMatchSameLengthAndMatchingNumberOfFields should be true")
	}
}

func TestMatchSameLengthNotMatchingNumberOfFields(t *testing.T) {
	var testTemplate = CreateTemplate("Field1", 2)
	var testTuple = CreateTuple("Field1", 2, 3.3, false)

	if testTuple.match(testTemplate) {
		t.Errorf("TestMatchSameLengthNotMatchingNumberOfFields should be false")
	}
}

func TestMatchNotSameLengthAndNotMatchingNumberOfFields(t *testing.T) {
	var testTemplate = CreateTemplate("Field1", 2, 3.3, false)
	testTemplate.DontMatchNumberOfFields()
	var testTuple = CreateTuple("Field1", 2, 3.3, false)

	if !testTuple.match(testTemplate) {
		t.Errorf("TestMatchNotSameLengthAndNotMatchingNumberOfFields should be true")
	}
}

func TestMatchNotSameLengthTemplateLongerNotMatchingNumberOfFields(t *testing.T) {
	var testTemplate = CreateTemplate("Field1", 2, 3.3, false, "extra field")
	testTemplate.DontMatchNumberOfFields()
	var testTuple = CreateTuple("Field1", 2, 3.3, false)

	if testTuple.match(testTemplate) {
		t.Errorf("TestMatchNotSameLengthTemplateLongerNotMatchingNumberOfFields should be false")
	}
}

func TestMatchFieldsOfSameType(t *testing.T) {
	var testTemplate = CreateTemplate(reflect.TypeOf("string"), reflect.TypeOf(32), reflect.TypeOf(3.14), reflect.TypeOf(true))
	testTemplate.DontMatchNumberOfFields()
	var testTuple = CreateTuple("Field1", 2, 3.3, false)

	if !testTuple.matchFieldsOf(testTemplate) {
		t.Errorf("TestMatchFieldsOfSameType")
	}
}

func TestMatchFieldsOfDifferentType(t *testing.T) {
	var testTemplate = CreateTemplate(reflect.TypeOf(false))
	testTemplate.DontMatchNumberOfFields()
	var testTuple = CreateTuple("Field1", 2, 3.3, false)

	if testTuple.matchFieldsOf(testTemplate) {
		t.Errorf("TestMatchFieldsOfSameType")
	}
}
