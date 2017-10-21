package protocol

import (
	"reflect"
	"testing"
)

// Test to see if Message is creating correct.
func TestCreateMessage(t *testing.T) {
	// Setup
	testMessage := createTestMessage()

	// Create Message manually.
	actualOperation := GetRequest
	actualT := []interface{}{"3", true, 4}
	actualMessage := Message{actualOperation, actualT}

	// Test that the two templates are equal.
	messagesEqual := reflect.DeepEqual(testMessage, actualMessage)

	if !messagesEqual {
		t.Errorf("CreateMessage() gave %+v, should be %+v", testMessage, actualMessage)
	}
}

// Test to see if GetOperation returns the correct operation.
func TestMessageGetOperation(t *testing.T) {
	// Setup
	testMessage := createTestMessage()

	actualOperation := GetRequest

	// Get the operation of the message with method.
	testOperation := testMessage.GetOperation()

	operationsEqual := reflect.DeepEqual(testOperation, actualOperation)

	if !operationsEqual {
		t.Errorf("GetOperation() on message: %+v == %v, should be %v", testMessage, testOperation, actualOperation)
	}
}

// Test to see if GetBody returns the correct body.
func TestGetBody(t *testing.T) {
	// Setup
	testMessage := createTestMessage()

	actualBody := []interface{}{"3", true, 4}

	// Get the operation of the message with method.
	testBody := testMessage.GetBody()

	bodiesEqual := reflect.DeepEqual(testBody, actualBody)

	if !bodiesEqual {
		t.Errorf("GetBody() on message: %+v == %v, should be %v", testMessage, testBody, actualBody)
	}
}

func createTestMessage() Message {
	testOperation := GetRequest
	testT := []interface{}{"3", true, 4}

	return CreateMessage(testOperation, testT)
}
