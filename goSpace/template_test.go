package tuplespace

import (
	"reflect"
	"testing"
)

var testTemplateOneField Template
var testTemplateFourFields Template

func init() {
	// Create simple template. Default value of matchNumberOfFields is true
	testTemplateOneField = CreateTemplate("Test Field")
	// Create simple template with four fields.
	testTemplateFourFields = CreateTemplate("Field 1", 2, 3.14, false)
}

// Test to see if template is creating correct.
func TestCreateTemplate(t *testing.T) {
	fields := make([]interface{}, 4)
	fields[0] = "Field 1"
	fields[1] = 2
	fields[2] = 3.14
	fields[3] = false

	// First argument is true, as it is the default value when creating a
	// template.
	var testTemplateManual = Template{true, fields}

	// Test that the two templates are equal.
	templatesEqual := reflect.DeepEqual(testTemplateManual, testTemplateFourFields)

	if !templatesEqual {
		t.Errorf("CreateTemplate() gave %+v, should be %+v", testTemplateFourFields, testTemplateManual)
	}
}

// Test that will check if the correct boolean value of matchNumberOfFields is
// returned.
func TestTemplateGetMatchNumberOfFields(t *testing.T) {
	defaultvalueOfMatchNumberOfField := true
	valueOfMatchNumberOfFields := testTemplateOneField.getMatchNumberOfFields()

	if valueOfMatchNumberOfFields != defaultvalueOfMatchNumberOfField {
		t.Errorf("testTemplate's matchNumberOfFields was not returned correctly, should have been true.")
	}
}

// Test that will change the boolean value of matchNumberOfFields back and forth.
func TestChangeMatchNumberOfFields(t *testing.T) {
	// Make a copy of the constant defined template with default value of
	// matchNumberOfFields to true.
	testTemplate := testTemplateOneField

	if testTemplate.getMatchNumberOfFields() != true {
		t.Errorf("testTemplate's matchNumberOfFields was not set to the correct default value being true.")
	}

	// Set boolean value to false.
	testTemplate.DontMatchNumberOfFields()

	if testTemplate.getMatchNumberOfFields() != false {
		t.Errorf("testTemplate's matchNumberOfFields was not set to the correct value being false by method DontMatchNumberOfFields().")
	}

	// Set boolean value back to true.
	testTemplate.DoMatchNumberOfFields()

	if testTemplate.getMatchNumberOfFields() != true {
		t.Errorf("testTemplate's matchNumberOfFields was not set to the correct value being true by method MatchNumberOfFields().")
	}
}

// Test to see if template has correct length.
func TestTemplateLength(t *testing.T) {
	// Create template with same length as defined by variable.
	var templateLength = 4
	testTemplateLength := testTemplateFourFields.length()
	if templateLength != testTemplateLength {
		t.Errorf("Length(%+v) == %d, should be %d", testTemplateFourFields, testTemplateLength, templateLength)
	}
}

// Test to see if the correct template field is returned.
func TestTemplateGetFieldAt(t *testing.T) {
	var wantedField = false
	templateField := testTemplateFourFields.getFieldAt(3)
	if wantedField != templateField {
		t.Errorf("GetFieldAt(%d) on template: %+v == %v, should be %v", 4, testTemplateFourFields, templateField, wantedField)
	}
}
