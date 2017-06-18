package tuplespace

import (
	"reflect"
	"sync"
	"testing"
)

func TestCreateTupleSpace(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9000)

	// Manually create tuple space
	actualMuTuple := new(sync.RWMutex)
	actualMuWaitingClients := new(sync.Mutex)
	actualTuples := []Tuple{}
	actualPort := ":9000"
	actualWaitingClients := []WaitingClient{}
	actualTupleSpace := &TupleSpace{muTuples: actualMuTuple, muWaitingClients: actualMuWaitingClients, tuples: actualTuples, port: actualPort, waitingClients: actualWaitingClients}

	// Test that the two templates are equal.
	tupleSpacesEqual := true
	if reflect.TypeOf(testTupleSpace.muTuples) != reflect.TypeOf(actualTupleSpace.muTuples) {
		tupleSpacesEqual = false
	} else if reflect.TypeOf(testTupleSpace.muWaitingClients) != reflect.TypeOf(actualTupleSpace.muWaitingClients) {
		tupleSpacesEqual = false
	} else if reflect.TypeOf(testTupleSpace.tuples) != reflect.TypeOf(actualTupleSpace.tuples) && len(testTupleSpace.tuples) == 0 && len(actualTupleSpace.tuples) == 0 {
		tupleSpacesEqual = false
	} else if !reflect.DeepEqual(testTupleSpace.port, actualTupleSpace.port) {
		tupleSpacesEqual = false
	} else if reflect.TypeOf(testTupleSpace.waitingClients) != reflect.TypeOf(actualTupleSpace.waitingClients) && len(testTupleSpace.waitingClients) == 0 && len(actualTupleSpace.waitingClients) == 0 {
		tupleSpacesEqual = false
	}

	if !tupleSpacesEqual {
		t.Errorf("CreateTupleSpace() gave %+v, should be %+v", testTupleSpace, actualTupleSpace)
	}
}

func TestSize(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9001)

	actualSize := 0

	testSize := testTupleSpace.Size()

	sizesEqual := testSize == actualSize

	if !sizesEqual {
		t.Errorf("Size() on tuple space: %+v == %v, should be %v", testTupleSpace, testSize, actualSize)
	}
}

func TestAddNewClient(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9002)

	// Make a client.
	temp := CreateTemplate([]interface{}{"Field 1"})
	response := make(chan *Tuple)
	remove := false // QueryRequest
	actualWaitingClient := CreateWaitingClient(temp, response, remove)

	// Add client to the data structure with method.
	testTupleSpace.addNewClient(actualWaitingClient)

	// Check that the size of the waitingClients is 1.
	if len(testTupleSpace.waitingClients) != 1 {
		t.Errorf("The size of %+v was %d but was expected to have the size 1 after the client %+v was added to an empty data structure", testTupleSpace.waitingClients, len(testTupleSpace.waitingClients), actualWaitingClient)
	}
}

func TestRemoveClientAt(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9003)
	actualWaitingClient := CreateWaitingClient(CreateTemplate([]interface{}{"Field 1"}), make(chan *Tuple), false)
	testTupleSpace.addNewClient(actualWaitingClient)

	// Remove client with method
	testTupleSpace.removeClientAt(0)

	isWaitingClientsEmpty := len(testTupleSpace.waitingClients) == 0

	if !isWaitingClientsEmpty {
		t.Errorf("The size of %+v was %d but was expected to have the size 0 after the client %+v was removed.", testTupleSpace.waitingClients, len(testTupleSpace.waitingClients), actualWaitingClient)
	}
}

// TestPutPOneMatchingQueryOneMatchingGet will make sure that both the
// QueryRequest and the GetRequest get the tuple.
func TestPutPOneMatchingQueryOneMatchingGet(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9004)
	// Add one matching QueryRequest and one matching GetRequest
	queryChan := make(chan *Tuple)
	getChan := make(chan *Tuple)
	waitingClientQuery := CreateWaitingClient(CreateTemplate([]interface{}{"Matching field"}), queryChan, false)
	waitingClientGet := CreateWaitingClient(CreateTemplate([]interface{}{"Matching field"}), getChan, true)
	testTupleSpace.addNewClient(waitingClientQuery)
	testTupleSpace.addNewClient(waitingClientGet)
	testTuple := CreateTuple([]interface{}{"Matching field"})

	go testTupleSpace.putP(&testTuple)

	queryResponse := <-queryChan

	if !reflect.DeepEqual(*queryResponse, testTuple) {
		t.Errorf("The tuple received by %+v is %+v but was expected to be %+v.", waitingClientQuery, *queryResponse, testTuple)
	}

	getResponse := <-getChan

	if !reflect.DeepEqual(*getResponse, testTuple) {
		t.Errorf("The tuple received by %+v is %+v but was expected to be %+v.", waitingClientGet, *getResponse, testTuple)
	}

	isTupleSpaceEmpty := testTupleSpace.Size() == 0

	if !isTupleSpaceEmpty {
		t.Errorf("The size of %+v was %d but was expected to have the size 0.", testTupleSpace.tuples, testTupleSpace.Size())
	}
}

func TestPutPNoWaitingClients(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9005)
	testTuple := CreateTuple([]interface{}{"Matching field"})

	testTupleSpace.putP(&testTuple)

	if testTupleSpace.Size() != 1 {
		t.Errorf("The size of %+v was %d but was expected to have size 1 after %+v was put.", testTupleSpace.tuples, testTupleSpace.Size(), testTuple)
	}
}

func TestPut(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9006)
	testTuple := CreateTuple([]interface{}{"Matching field"})

	testChan := make(chan bool)

	go testTupleSpace.put(&testTuple, testChan)

	testResponse := <-testChan

	if testResponse != true {
		t.Errorf("The response of put was %t but was expected to be true.", testResponse)
	}
}

func TestClearTupleSpace(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9007)
	testTuple := CreateTuple([]interface{}{"Matching field"})
	testTupleSpace.putP(&testTuple)

	testTupleSpace.clearTupleSpace()

	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %+v was %d but was expected to have the size 0.", testTupleSpace.tuples, testTupleSpace.Size())
	}
}

func TestTupleSpaceRemoveTupleAt(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9008)
	testTuple1 := CreateTuple([]interface{}{"Tuple", 1})
	testTuple2 := CreateTuple([]interface{}{"Tuple", 2})
	testTuple3 := CreateTuple([]interface{}{"Tuple", 3})
	testTupleSpace.putP(&testTuple1)
	testTupleSpace.putP(&testTuple2)
	testTupleSpace.putP(&testTuple3)

	// Remove the middle tuple, index 1 and see the last tuple is moved to
	// that index.
	testTupleSpace.removeTupleAt(1)

	if !reflect.DeepEqual(testTupleSpace.tuples[0], testTuple1) {
		t.Errorf("The tuple at 0 is %+v but was expected to be %+v.", testTupleSpace.tuples[0], testTuple1)
	}

	if !reflect.DeepEqual(testTupleSpace.tuples[1], testTuple3) {
		t.Errorf("The tuple at 1 is %+v but was expected to be %+v.", testTupleSpace.tuples[1], testTuple3)
	}
}

func TestFindTupleMatchingTupleRemove(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9009)
	testTuple := CreateTuple([]interface{}{"Matching field"})
	testTupleSpace.putP(&testTuple)

	testTemplate := CreateTemplate([]interface{}{"Matching field"})
	findTupleResult := testTupleSpace.findTuple(testTemplate, true)

	// Check the correct tuple was found.
	if !reflect.DeepEqual(*findTupleResult, testTuple) {
		t.Errorf("The tuple found %+v was expected to look like %+v.", findTupleResult, testTuple)
	}

	// Check that the tuple was removed from the tuple space.
	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %+v was %d but was expected to have the size 0.", testTupleSpace.tuples, testTupleSpace.Size())
	}
}

func TestFindTupleMatchingTupleNoRemove(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9010)
	testTuple := CreateTuple([]interface{}{"Matching field"})
	testTupleSpace.putP(&testTuple)

	testTemplate := CreateTemplate([]interface{}{"Matching field"})
	findTupleResult := testTupleSpace.findTuple(testTemplate, false)

	// Check the correct tuple was found.
	if !reflect.DeepEqual(*findTupleResult, testTuple) {
		t.Errorf("The tuple found %+v was expected to look like %+v.", findTupleResult, testTuple)
	}

	// Check that the tuple was removed from the tuple space.
	if testTupleSpace.Size() != 1 {
		t.Errorf("The size of %+v was %d but was expected to have the size 1.", testTupleSpace.tuples, testTupleSpace.Size())
	}
}

func TestFindTupleNoMatchingTuple(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9011)

	testTemplate := CreateTemplate([]interface{}{"Matching field"})
	findTupleResult := testTupleSpace.findTuple(testTemplate, true)

	// Check the correct tuple was found.
	if findTupleResult != nil {
		t.Errorf("The tuple found %+v was expected to be nil.", findTupleResult)
	}

	// Check that the tuple was removed from the tuple space.
	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %+v was %d but was expected to have the size 0.", testTupleSpace.tuples, testTupleSpace.Size())
	}
}

func TestFindTupleBlockingMatchingTuple(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9012)
	testTuple := CreateTuple([]interface{}{"Matching field"})
	testTupleSpace.putP(&testTuple)

	testTemplate := CreateTemplate([]interface{}{"Matching field"})
	testChan := make(chan *Tuple)
	go testTupleSpace.findTupleBlocking(testTemplate, testChan, true)

	testResponse := <-testChan

	// Check the correct tuple was found.
	if !reflect.DeepEqual(*testResponse, testTuple) {
		t.Errorf("The tuple found %+v was expected to look like %+v.", *testResponse, testTuple)
	}

	// Check that the tuple was removed from the tuple space.
	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %+v was %d but was expected to have the size 0.", testTupleSpace.tuples, testTupleSpace.Size())
	}
}

func TestFindTupleBlockingNoMatchingTuple(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9013)

	testTemplate := CreateTemplate([]interface{}{"Matching field"})
	testChan := make(chan *Tuple)
	testTupleSpace.findTupleBlocking(testTemplate, testChan, true)

	// Check that the client was added to waitingClients.
	if len(testTupleSpace.waitingClients) != 1 {
		t.Errorf("The size of %+v was %d but was expected to have the size 1.", testTupleSpace.waitingClients, len(testTupleSpace.waitingClients))
	}
}

func TestGet(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9014)
	testTuple := CreateTuple([]interface{}{"Matching field"})
	testTupleSpace.putP(&testTuple)

	testTemplate := CreateTemplate([]interface{}{"Matching field"})
	testChan := make(chan *Tuple)
	go testTupleSpace.get(testTemplate, testChan)

	testResponse := <-testChan

	if !reflect.DeepEqual(*testResponse, testTuple) {
		t.Errorf("get() gave %+v but was expected to return %+v.", *testResponse, testTuple)
	}

	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %+v was %d but was expected to have the size 0.", testTupleSpace.tuples, testTupleSpace.Size())
	}
}

func TestQuery(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9015)
	testTuple := CreateTuple([]interface{}{"Matching field"})
	testTupleSpace.putP(&testTuple)

	testTemplate := CreateTemplate([]interface{}{"Matching field"})
	testChan := make(chan *Tuple)
	go testTupleSpace.query(testTemplate, testChan)

	testResponse := <-testChan

	if !reflect.DeepEqual(*testResponse, testTuple) {
		t.Errorf("query() gave %+v but was expected to return %+v.", *testResponse, testTuple)
	}

	if testTupleSpace.Size() != 1 {
		t.Errorf("The size of %+v was %d but was expected to have the size 1.", testTupleSpace.tuples, testTupleSpace.Size())
	}
}

func TestFindTupleNonblocking(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9016)
	testTuple := CreateTuple([]interface{}{"Matching field"})
	testTupleSpace.putP(&testTuple)

	testTemplate := CreateTemplate([]interface{}{"Matching field"})
	testChan := make(chan *Tuple)
	go testTupleSpace.findTupleNonblocking(testTemplate, testChan, true)

	testResponse := <-testChan

	if !reflect.DeepEqual(*testResponse, testTuple) {
		t.Errorf("FindTupleNonblocking gave %+v but was expected to return %+v.", *testResponse, testTuple)
	}
}

func TestGetP(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9017)
	testTuple := CreateTuple([]interface{}{"Matching field"})
	testTupleSpace.putP(&testTuple)

	testTemplate := CreateTemplate([]interface{}{"Matching field"})
	testChan := make(chan *Tuple)
	go testTupleSpace.getP(testTemplate, testChan)

	testResponse := <-testChan

	if !reflect.DeepEqual(*testResponse, testTuple) {
		t.Errorf("getP() gave %+v but was expected to return %+v.", *testResponse, testTuple)
	}

	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %+v was %d but was expected to have the size 0.", testTupleSpace.tuples, testTupleSpace.Size())
	}
}

func TestQueryP(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9018)
	testTuple := CreateTuple([]interface{}{"Matching field"})
	testTupleSpace.putP(&testTuple)

	testTemplate := CreateTemplate([]interface{}{"Matching field"})
	testChan := make(chan *Tuple)
	go testTupleSpace.queryP(testTemplate, testChan)

	testResponse := <-testChan

	if !reflect.DeepEqual(*testResponse, testTuple) {
		t.Errorf("queryP() gave %+v but was expected to return %+v.", *testResponse, testTuple)
	}

	if testTupleSpace.Size() != 1 {
		t.Errorf("The size of %+v was %d but was expected to have the size 1.", testTupleSpace.tuples, testTupleSpace.Size())
	}
}

func TestFindAllTuplesMatchingTuplesRemove(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9019)
	testTuple := CreateTuple([]interface{}{"Matching field"})
	testTupleSpace.putP(&testTuple)
	testTupleSpace.putP(&testTuple)

	testTemplate := CreateTemplate([]interface{}{"Matching field"})
	testChan := make(chan []Tuple)
	go testTupleSpace.findAllTuples(testTemplate, testChan, true)

	testReponse := <-testChan

	if len(testReponse) != 2 {
		t.Errorf("The size of %+v was %d but was expected to have size 2", testReponse, len(testReponse))
	}

	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %+v was %d but was expected to have size 0", testTupleSpace.tuples, testTupleSpace.tuples)
	}
}

func TestFindAllTuplesMatchingTuplesNoRemove(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9020)
	testTuple := CreateTuple([]interface{}{"Matching field"})
	testTupleSpace.putP(&testTuple)
	testTupleSpace.putP(&testTuple)

	testTemplate := CreateTemplate([]interface{}{"Matching field"})
	testChan := make(chan []Tuple)
	go testTupleSpace.findAllTuples(testTemplate, testChan, false)

	testReponse := <-testChan

	if len(testReponse) != 2 {
		t.Errorf("The size of %+v was %d but was expected to have size 2", testReponse, len(testReponse))
	}

	if testTupleSpace.Size() != 2 {
		t.Errorf("The size of %+v was %d but was expected to have size 2", testTupleSpace.tuples, testTupleSpace.tuples)
	}
}

func TestFindAllTuplesNoMatchingTuples(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9021)

	testTemplate := CreateTemplate([]interface{}{"Matching field"})
	testChan := make(chan []Tuple)
	go testTupleSpace.findAllTuples(testTemplate, testChan, false)

	testReponse := <-testChan

	if len(testReponse) != 0 {
		t.Errorf("The size of %+v was %d but was expected to have size 0", testReponse, len(testReponse))
	}

	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %+v was %d but was expected to have size 0", testTupleSpace.tuples, testTupleSpace.tuples)
	}
}

func TestGetAll(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9022)
	testTuple := CreateTuple([]interface{}{"Matching field"})
	testTupleSpace.putP(&testTuple)
	testTupleSpace.putP(&testTuple)

	testTemplate := CreateTemplate([]interface{}{"Matching field"})
	testChan := make(chan []Tuple)
	go testTupleSpace.getAll(testTemplate, testChan)

	testReponse := <-testChan

	if len(testReponse) != 2 {
		t.Errorf("The size of %+v was %d but was expected to have size 2", testReponse, len(testReponse))
	}

	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %+v was %d but was expected to have size 0", testTupleSpace.tuples, testTupleSpace.tuples)
	}
}

func TestQueryAll(t *testing.T) {
	// Setup
	testTupleSpace := createTestTupleSpace(9023)
	testTuple := CreateTuple([]interface{}{"Matching field"})
	testTupleSpace.putP(&testTuple)
	testTupleSpace.putP(&testTuple)

	testTemplate := CreateTemplate([]interface{}{"Matching field"})
	testChan := make(chan []Tuple)
	go testTupleSpace.queryAll(testTemplate, testChan)

	testReponse := <-testChan

	if len(testReponse) != 2 {
		t.Errorf("The size of %+v was %d but was expected to have size 2", testReponse, len(testReponse))
	}

	if testTupleSpace.Size() != 2 {
		t.Errorf("The size of %+v was %d but was expected to have size 2", testTupleSpace.tuples, testTupleSpace.tuples)
	}
}

func createTestTupleSpace(testPort int) *TupleSpace {
	return CreateTupleSpace(testPort)
}
