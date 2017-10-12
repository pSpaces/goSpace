package protocol

import (
	. "github.com/pspaces/gospace/shared"
	"reflect"
	"strings"
	"testing"
)

// Test to see if WaitingClient is creating correct.
func TestCreateWaitingClientGet(t *testing.T) {
	// Setup
	testWaitingClient := createTestWaitingClient()

	// Create waitingClient manually.
	actualFields := make([]interface{}, 1)
	actualFields[0] = "Field 1"
	actualChan := make(chan<- *Tuple)
	actualOperation := GetRequest
	actualWaitingClient := WaitingClient{CreateTemplate(actualFields), actualChan, actualOperation}

	// Test that the two templates are equal.
	waitingClientsEqual := true
	if !reflect.DeepEqual(testWaitingClient.template, actualWaitingClient.template) {
		waitingClientsEqual = false
	} else if reflect.TypeOf(testWaitingClient.responseChan) != reflect.TypeOf(actualWaitingClient.responseChan) {
		waitingClientsEqual = false
	} else if strings.Compare(testWaitingClient.operation, actualWaitingClient.operation) != 0 {
		waitingClientsEqual = false
	}

	if !waitingClientsEqual {
		t.Errorf("CreateWaitingClient() gave %+v, should be %+v", testWaitingClient, actualWaitingClient)
	}
}

// Test to see if WaitingClient is creating correct.
func TestCreateWaitingClientQuery(t *testing.T) {
	// Create waitingClient manually.
	testFields := make([]interface{}, 1)
	testFields[0] = "Field 1"
	testChan := make(chan<- *Tuple)
	testBool := false
	testWaitingClient := CreateWaitingClient(CreateTemplate(testFields), testChan, testBool)

	// Test that the two templates are equal.
	operationsEqual := true
	if strings.Compare(testWaitingClient.operation, QueryRequest) != 0 {
		operationsEqual = false
	}

	if !operationsEqual {
		t.Errorf("CreateWaitingClient() passed %+v as %q, should be %q", testBool, testWaitingClient.operation, QueryRequest)
	}
}

// Test to see if template is returned correctly.
func TestGetTemplate(t *testing.T) {
	// Setup
	testWaitingClient := createTestWaitingClient()

	actualFields := make([]interface{}, 1)
	actualFields[0] = "Field 1"
	actualTemplate := CreateTemplate(actualFields)

	testTemplate := testWaitingClient.GetTemplate()

	templatesEqual := reflect.DeepEqual(testTemplate, actualTemplate)

	if !templatesEqual {
		t.Errorf("GetTemplate() on waitingClient: %+v == %v, should be %v", testWaitingClient, testTemplate, actualTemplate)
	}
}

// Test to see if channel is returned correctly.
func TestGetResponseChan(t *testing.T) {
	// Setup
	testWaitingClient := createTestWaitingClient()

	actualChan := make(chan<- *Tuple)

	testChan := testWaitingClient.GetResponseChan()

	channelsEqual := true
	if reflect.TypeOf(testChan) != reflect.TypeOf(actualChan) {
		channelsEqual = false
	}

	if !channelsEqual {
		t.Errorf("GetReponseChan() on waitingClient: %+v == %v, should be %v", testWaitingClient, testChan, actualChan)
	}
}

// Test to see if operation is returned correctly.
func TestGetOperation(t *testing.T) {
	// Setup
	testWaitingClient := createTestWaitingClient()

	actualOperation := GetRequest

	testOperation := testWaitingClient.GetOperation()

	operationsEqual := reflect.DeepEqual(testOperation, actualOperation)

	if !operationsEqual {
		t.Errorf("GetOperation() on waitingClient: %+v == %v, should be %v", testWaitingClient, testOperation, actualOperation)
	}
}

func createTestWaitingClient() WaitingClient {
	testFields := make([]interface{}, 1)
	testFields[0] = "Field 1"
	testTemplate := CreateTemplate(testFields)

	testChan := make(chan<- *Tuple)

	testBool := true

	return CreateWaitingClient(testTemplate, testChan, testBool)
}
