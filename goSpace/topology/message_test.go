package topology

import (
	"goSpace/goSpace/constants"
	"reflect"
	"testing"
)

var testMessage Message

func init() {
	testMessage = CreateMessage(constants.PutRequest, []interface{}{"3", true, 4})
}

// TestMessageCreateMessage will see if the method for creating the message does
// it as intended.
func TestMessageCreateMessage(t *testing.T) {
	messageManual := Message{Operation: constants.PutRequest, T: []interface{}{"3", true, 4}}

	// Check equality between messages.
	if !reflect.DeepEqual(testMessage, messageManual) {
		t.Errorf("The method generated message %q, didn't look as expected: %q", testMessage, messageManual)
	}
}

// TestMessageGetOperation will make sure the right part of the message is
// returned, the operation of the message.
func TestMessageGetOperation(t *testing.T) {
	// Create operation string for comparison.
	messageManualOperation := constants.PutRequest

	// Get the operation of the message with method.
	messageOperation := testMessage.GetOperation()

	if !reflect.DeepEqual(messageManualOperation, messageOperation) {
		t.Errorf("The operation from the message %q, didn't look as expected: %q", messageOperation, messageManualOperation)
	}
}

// TestMessageGetBody will make sure the right part of the message is
// returned, the body of the message.
func TestMessageBodyOperation(t *testing.T) {
	// Create body for comparison.
	messageManualBody := []interface{}{"3", true, 4}

	// Get the body of the message with method.
	messageBody := testMessage.GetBody()

	if !reflect.DeepEqual(messageManualBody, messageBody) {
		t.Errorf("The body from the message %q, didn't look as expected: %q", messageBody, messageManualBody)
	}
}
