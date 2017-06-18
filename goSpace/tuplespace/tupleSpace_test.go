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

func createTestTupleSpace(testPort int) *TupleSpace {
	return CreateTupleSpace(testPort)
}

/*

// TestTupleSpacePutP will make sure a tuple is added to the tuple space
// correctly.
func TestTupleSpacePutPNoWaitingClients(t *testing.T) {
	// Initially check the tuple space is empty by securing the size of it is
	// 0.
	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %q was %d but was expected to have the size 0.", testTupleSpace.tuples, testTupleSpace.Size())
	}

	// Initially check there are no waiting clients by securing the size of it
	// is 0.
	if len(testTupleSpace.waitingClients) != 0 {
		t.Errorf("The size of %+v was %d but was expected to have the size 0.", testTupleSpace.waitingClients, len(testTupleSpace.waitingClients))
	}

	// Add a tuple to the tuple space.
	testTupleSpace.putP(&testTupleFourFields)

	// Check that the size of the tuples is 1.
	if testTupleSpace.Size() != 1 {
		t.Errorf("The size of %q was %d but was expected to have the size 1 after to tuple %q was added to an empty tuple space", testTupleSpace.tuples, testTupleSpace.Size(), testTupleFourFields)
	}

	// Extract the tuple from the data structure for tuples in the tuple space
	// and check it was added without being altered.
	tuple := testTupleSpace.tuples[0]

	if !reflect.DeepEqual(testTupleFourFields, tuple) {
		t.Errorf("The tuple %q from the tuple spaces isn't equal to the tuple %q that was added to the tuple space.", tuple, testTupleFourFields)
	}
}

// TestTupleSpaceClearTupleSpace will make sure that the method will clear the
// data structure tuples[] for tuples.
func TestTupleSpaceClearTupleSpace(t *testing.T) {
	// Initially check that the tuple space contains some number of tuples.
	if testTupleSpace.Size() == 0 {
		t.Errorf("The tuple space is empty, %q, and should contain some amount of tuples to test clearTupleSpace().", testTupleSpace.tuples)
	}

	// Clear the tuple space for exsisting tuples.
	testTupleSpace.clearTupleSpace()

	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %q was %d but was expected to have the size 0.", testTupleSpace.tuples, testTupleSpace.Size())
	}
}

// TestTupleSpacePut will make sure a tuple is the right response is returned
// when a tuple has been added to the tuple space.
func TestTupleSpacePut(t *testing.T) {
	// Make sure the tuple space is empty.
	testTupleSpace.clearTupleSpace()

	// Create channel to receive the response.
	response := make(chan bool)

	// Add the tuple to the tuple space with the put()
	go func() {
		testTupleSpace.put(&testTupleFourFields, response)
		time.Sleep(time.Second * 1)
	}()

	// Read the response from the channel.
	b := <-response

	// Check the response is as expected.
	if !b {
		t.Errorf("TestTupleSpaceFindTuple(%t) should be true", b)
	}
}

// TestTupleSpaceRemoveTupleAt will make sure the right tuple is removed.
func TestTupleSpaceRemoveTupleAt(t *testing.T) {
	// Initially check that the tuple space contains some number of tuples.
	if testTupleSpace.Size() == 0 {
		t.Errorf("The tuple space is empty, %q, and should contain some amount of tuples to test clearTupleSpace().", testTupleSpace.tuples)
	}

	// Get the first tuple.
	tuple := testTupleSpace.tuples[0]

	// Remove the first tuple from the tuple space.
	testTupleSpace.removeTupleAt(0)

	// Check the tuple isn't in the tuple space anymore.
	for _, tupleFromTS := range testTupleSpace.tuples {
		if reflect.DeepEqual(tuple, tupleFromTS) {
			t.Errorf("The tuple %q is still in the tuples space %q", tuple, testTupleSpace.tuples)
		}
	}
}

// TestTupleSpaceFindTuple will make sure that the right tuple is found of the
// data structure of tuples.
func TestTupleSpaceFindTuple(t *testing.T) {
	// Make sure the tuple space is empty.
	testTupleSpace.clearTupleSpace()

	// Add a tuple to the tuple space.
	testTupleSpace.putP(&testTupleTwoFields)

	// Create a template that matches the tuple just added to the tuple space.
	temp := CreateTemplate("Field 1", true)

	// Find the tuple using the method, without removing the tuple.
	t1 := testTupleSpace.findTuple(temp, false)

	// Check that the tuple found is equal to the tuple that was added.
	if !reflect.DeepEqual(*t1, testTupleTwoFields) {
		t.Errorf("The tuple %q was found by the method by should be %q", *t1, testTupleTwoFields)
	}

	// Check that the tuple wasn't removed from the tuple space
	if testTupleSpace.Size() != 1 {
		t.Errorf("The size of %q was %d but was expected to have the size 1 as the tuple shouldn't have been removed from the tuple space", testTupleSpace.tuples, testTupleSpace.Size())
	}

	// Find the tuple using the method and removing the tuple.
	t2 := testTupleSpace.findTuple(temp, true)

	// Check that the tuple found is equal to the tuple that was added.
	if !reflect.DeepEqual(*t2, testTupleTwoFields) {
		t.Errorf("The tuple %q was found by the method by should be %q", *t1, testTupleTwoFields)
	}

	// Check that the tuple was removed from the tuple space
	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %q was %d but was expected to have the size 0 as the tuple should have been removed from the tuple space", testTupleSpace.tuples, testTupleSpace.Size())
	}

	// Find a tuple that doesn't exists in the tuple space.
	t3 := testTupleSpace.findTuple(temp, true)

	// Check that the tuple is nil, meaning there was no tuple matching the
	// template.
	if t3 != nil {
		t.Errorf("TestTupleSpaceFindTuple(%q) should be nil", t3)
	}
}

// TestTupleSpacePutPWithOneGetWaitingClient will make sure that is there's a
// client performing a get who matches the tuple, the tuple is sent to it and
// isn't added to the tuple space.
func TestTupleSpacePutPWithOneGetWaitingClient(t *testing.T) {
	// Initially check the tuple space is empty by securing the size of it is
	// 0.
	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %q was %d but was expected to have the size 0.", testTupleSpace.tuples, testTupleSpace.Size())
	}

	// Initially check there are no waiting clients by securing the size of it
	// is 0.
	if len(testTupleSpace.waitingClients) != 0 {
		t.Errorf("The size of %+v was %d but was expected to have the size 0.", testTupleSpace.waitingClients, len(testTupleSpace.waitingClients))
	}

	temp := CreateTemplate("test field")
	tupleChan := make(chan *Tuple)
	remove := true
	client := CreateWaitingClient(temp, tupleChan, remove)
	tupleResponse := new(Tuple)
	testTupleSpace.addNewClient(client)
	go func() {
		tupleResponse = <-tupleChan
	}()

	// Add a matching tuple to the tuple space.
	tuple := CreateTuple("test field")
	testTupleSpace.putP(&tuple)

	// Check that the size of the tuples is 0.
	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %q was %d but was expected to have the size 0 as the waiting client was performing a get operation.", testTupleSpace.tuples, testTupleSpace.Size(), tuple)
	}

	if len(testTupleSpace.waitingClients) != 0 {
		t.Errorf("The size of %+v was %d but was expected to have the size 0", testTupleSpace.waitingClients, len(testTupleSpace.waitingClients))
	}

	if !reflect.DeepEqual(tuple, *tupleResponse) {
		t.Errorf("The tuples were not equal, they tuple place by putP looks like %+v and the found tuple looks like %+v", tuple, tupleResponse)
	}
}

// TestTupleSpacePutPWithOneQueryWaitingClient will make sure that is there's a
// client performing a query who matches the tuple, the tuple is sent to it and
// isn't added to the tuple space.
func TestTupleSpacePutPWithOneQueryWaitingClient(t *testing.T) {
	// Initially check the tuple space is empty by securing the size of it is
	// 0.
	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %q was %d but was expected to have the size 0.", testTupleSpace.tuples, testTupleSpace.Size())
	}

	// Initially check there are no waiting clients by securing the size of it
	// is 0.
	if len(testTupleSpace.waitingClients) != 0 {
		t.Errorf("The size of %+v was %d but was expected to have the size 0.", testTupleSpace.waitingClients, len(testTupleSpace.waitingClients))
	}

	temp := CreateTemplate("test field")
	tupleChan := make(chan *Tuple)
	remove := false
	client := CreateWaitingClient(temp, tupleChan, remove)
	tupleResponse := new(Tuple)
	testTupleSpace.addNewClient(client)
	go func() {
		tupleResponse = <-tupleChan
	}()

	// Add a matching tuple to the tuple space.
	tuple := CreateTuple("test field")
	testTupleSpace.putP(&tuple)

	// Check that the size of the tuples is 1.
	if testTupleSpace.Size() != 1 {
		t.Errorf("The size of %q was %d but was expected to have the size 1 as the waiting client was performing a get operation.", testTupleSpace.tuples, testTupleSpace.Size())
	}

	if len(testTupleSpace.waitingClients) != 0 {
		t.Errorf("The size of %+v was %d but was expected to have the size 0", testTupleSpace.waitingClients, len(testTupleSpace.waitingClients))
	}

	if !reflect.DeepEqual(tuple, *tupleResponse) {
		t.Errorf("The tuples were not equal, they tuple place by putP looks like %+v and the found tuple looks like %+v", tuple, tupleResponse)
	}
}

// TestTupleSpacePutPWithOneQueryWaitingClient will make sure that is there's a
// client performing a query who matches the tuple followed by a client
// performing a get who matches the tuple. The tuple is sent to both client due
// to their order and it isn't added to the tuple space.
func TestTupleSpacePutPWithOneQueryOneGetWaitingClient(t *testing.T) {
	// Make sure tuple space is empty.
	testTupleSpace.clearTupleSpace()

	// Initially check there are no waiting clients by securing the size of it
	// is 0.
	if len(testTupleSpace.waitingClients) != 0 {
		t.Errorf("The size of %+v was %d but was expected to have the size 0.", testTupleSpace.waitingClients, len(testTupleSpace.waitingClients))
	}

	temp := CreateTemplate("test field")
	tupleChanQ := make(chan *Tuple)
	tupleChanG := make(chan *Tuple)
	bQuery := false
	bGet := true
	cQuery := CreateWaitingClient(temp, tupleChanQ, bQuery)
	testTupleSpace.addNewClient(cQuery)
	cGet := CreateWaitingClient(temp, tupleChanG, bGet)
	testTupleSpace.addNewClient(cGet)
	tupleResponseQ := new(Tuple)
	tupleResponseG := new(Tuple)

	go func() {
		tupleResponseQ = <-tupleChanQ
	}()
	go func() {
		tupleResponseG = <-tupleChanG
	}()

	// Add a matching tuple to the tuple space.
	tuple := CreateTuple("test field")
	testTupleSpace.putP(&tuple)

	time.Sleep(1 * time.Second)

	// Check that the size of the tuples is 1.
	if testTupleSpace.Size() != 0 {
		t.Errorf("The size of %q was %d but was expected to have the size 0 as the waiting clients were performing a query followed by a get operation.", testTupleSpace.tuples, testTupleSpace.Size())
	}

	if len(testTupleSpace.waitingClients) != 0 {
		t.Errorf("The size of %+v was %d but was expected to have the size 0", testTupleSpace.waitingClients, len(testTupleSpace.waitingClients))
	}

	// Check that the Query found the correct tuple.
	if !reflect.DeepEqual(tuple, *tupleResponseQ) {
		t.Errorf("The tuples were not equal, they tuple place by putP looks like %+v and the found tuple by Query looks like %+v", tuple, tupleResponseQ)
	}

	// Check that the Get found the correct tuple.
	if !reflect.DeepEqual(tuple, *tupleResponseG) {
		t.Errorf("The tuples were not equal, they tuple place by putP looks like %+v and the found tuple by Get looks like %+v", tuple, tupleResponseG)
	}
}
*/
// TestTupleSpaceFindTupleBlocking will make sure that if a tuple matching the
// template doesn't exsist in the tuple space it will add the client to
// waitingClients.
/*
func TestTupleSpaceFindTupleBlocking(t *testing.T) {
	// Make sure the tuple space is empty.
	testTupleSpace.muTuples.Lock()
	testTupleSpace.clearTupleSpace()
	testTupleSpace.muTuples.Unlock()

	if len(testTupleSpace.waitingClients) != 0 {
		t.Errorf("The size of %+v was %d but was expected to have the size 0.", testTupleSpace.waitingClients, len(testTupleSpace.waitingClients))
	}

	// Create a template that matches the tuple just added to the tuple space.
	temp := CreateTemplate("Field 1")

	// Create a channel to receive the response
	response := make(chan<- *Tuple)

	testTupleSpace.findTupleBlocking(temp, response, true)

	if len(testTupleSpace.waitingClients) != 1 {
		t.Errorf("The size of %+v was %d but was expected to have the size 1.", testTupleSpace.waitingClients, len(testTupleSpace.waitingClients))
	}

	tuple := CreateTuple("Field 1")
	go testTupleSpace.putP(&tuple)

	if len(testTupleSpace.waitingClients) != 0 {
		t.Errorf("The size of %+v was %d but was expected to have the size 0.", testTupleSpace.waitingClients, len(testTupleSpace.waitingClients))
	}
}*/

/*
// TestTupleSpaceQueryAndGet will make sure that query() and get() returns the
// correct tuple.
func TestTupleSpaceQueryAndGet(t *testing.T) {
	// Make sure the tuple space is empty.
	testTupleSpace.clearTupleSpace()

	// Add a tuple to the tuple space.
	testTupleSpace.putP(&testTupleTwoFields)

	// Create a template that matches the tuple just added to the tuple space.
	temp := CreateTemplate("Field 1", true)

	// Create a channel to receive the response
	response := make(chan *Tuple)

	// Find a tuple matching the template with the query method.
	go testTupleSpace.query(temp, response)

	// Read the tuple found.
	tQuery := <-response

	// Check that it's equal to the tuple added.
	if !reflect.DeepEqual(*tQuery, testTupleTwoFields) {
		t.Errorf("The query found %q but should have found %q", tQuery, testTupleTwoFields)
	}

	// Find a tupe matching the template with the get method.
	go testTupleSpace.get(temp, response)

	// Read the tuple found.
	tGet := <-response

	// Check that it's equal to the tuple added.
	if !reflect.DeepEqual(*tGet, testTupleTwoFields) {
		t.Errorf("The get found %q but should have found %q", tGet, testTupleTwoFields)
	}
}

// TestTupleSpaceFindTupleNonlocking will make sure the method isn't blocking
// if there isn't a tuple matching a template in the tuple space.
func TestTupleSpaceFindTupleNonlocking(t *testing.T) {
	// Make sure the tuple space is empty.
	testTupleSpace.clearTupleSpace()

	// Create a template that matches the tuple just added to the tuple space.
	temp := CreateTemplate("Field 1", true)

	// Create a channel to receive the response
	response := make(chan *Tuple)

	// Run the findTupleNonblocking, with no tuples in the tuple space.
	go testTupleSpace.findTupleNonblocking(temp, response, true)

	// Read the tuple found.
	tNonblockingNil := <-response

	// Check the tuple found is nil as no tuple were in the tuple space.
	if tNonblockingNil != nil {
		t.Errorf("tNonblockingNil should have been nil but was %q instead", tNonblockingNil)
	}

	// Add a tuple to the tuple space.
	testTupleSpace.putP(&testTupleTwoFields)

	// Run the findTupleNonblocking, with a macthing tuple in the tuple space.
	go testTupleSpace.findTupleNonblocking(temp, response, true)

	// Read the tuple found.
	tNonblocking := <-response

	// Check that the found tuple is equal to the tuple added to the tuple
	// space.
	if !reflect.DeepEqual(*tNonblocking, testTupleTwoFields) {
		t.Errorf("The findTupleNonblocking found %q but should have found %q", tNonblocking, testTupleTwoFields)
	}
}

// TestTupleSpaceQueryPAndGetP will make sure that queryP() and getP() returns
// the correct tuple.
func TestTupleSpaceQueryPAndGetP(t *testing.T) {
	// Make sure the tuple space is empty.
	testTupleSpace.clearTupleSpace()

	// Add a tuple to the tuple space.
	testTupleSpace.putP(&testTupleTwoFields)

	// Create a template that matches the tuple just added to the tuple space.
	temp := CreateTemplate("Field 1", true)

	// Create a channel to receive the response
	response := make(chan *Tuple)

	// Find a tuple matching the template with the query method.
	go testTupleSpace.queryP(temp, response)

	// Read the tuple found.
	tQueryP := <-response

	// Check that it's equal to the tuple added.
	if !reflect.DeepEqual(*tQueryP, testTupleTwoFields) {
		t.Errorf("The query found %q but should have found %q", tQueryP, testTupleTwoFields)
	}

	// Find a tupe matching the template with the get method.
	go testTupleSpace.getP(temp, response)

	// Read the tuple found.
	tGetP := <-response

	// Check that it's equal to the tuple added.
	if !reflect.DeepEqual(*tGetP, testTupleTwoFields) {
		t.Errorf("The get found %q but should have found %q", tGetP, testTupleTwoFields)
	}
}

// TestTupleSpaceFindAllTuples will make sure all tuples from the tuple space
// are returned.
func TestTupleSpaceFindAllTuples(t *testing.T) {
	// Make sure the tuple space is empty.
	testTupleSpace.clearTupleSpace()

	// Add a tuple twice to the tuple space.
	testTupleSpace.putP(&testTupleTwoFields)
	testTupleSpace.putP(&testTupleTwoFields)

	// Make local list of tuples for comparison.
	tuplesLocal := []Tuple{testTupleTwoFields, testTupleTwoFields}

	// Create a channel to receive the response
	response := make(chan []Tuple)

	// Run the findAllTuples, without removing them.
	go testTupleSpace.findAllTuples(response, false)

	// Read the tuples.
	tuples := <-response

	// Check the list of tuples found is equal to the local list of tuples.
	if !reflect.DeepEqual(tuplesLocal, tuples) {
		t.Errorf("findAllTuples returned %q but was expected to return %q.", tuples, tuplesLocal)
	}

	// Run the findAllTuples and remove them.
	go testTupleSpace.findAllTuples(response, true)

	// Read the tuples.
	tuples = <-response

	// Check the list of tuples found is equal to the local list of tuples.
	if !reflect.DeepEqual(tuplesLocal, tuples) {
		t.Errorf("findAllTuples returned %q but was expected to return %q.", tuples, tuplesLocal)
	}
}

// TestTupleSpaceQueryAllAndGetAll will make sure that queryAll() returns every
// tuple without removing them and that getAll() returns every tuple along with
// removing them from the tuple space.
func TestTupleSpaceQueryAllAndGetAll(t *testing.T) {
	// Make sure the tuple space is empty.
	testTupleSpace.clearTupleSpace()

	// Add a tuple twice to the tuple space.
	testTupleSpace.putP(&testTupleTwoFields)
	testTupleSpace.putP(&testTupleTwoFields)

	// Make local list of tuples for comparison.
	tuplesLocal := []Tuple{testTupleTwoFields, testTupleTwoFields}

	// Create a channel to receive the response
	response := make(chan []Tuple)

	// Find a tuple matching the template with the query method.
	go testTupleSpace.queryAll(response)

	// Read the tuples.
	tuplesQuery := <-response

	// Check the list of tuples found is equal to the local list of tuples.
	if !reflect.DeepEqual(tuplesLocal, tuplesQuery) {
		t.Errorf("findAllTuples returned %q but was expected to return %q.", tuplesQuery, tuplesLocal)
	}

	// Find a tupe matching the template with the get method.
	go testTupleSpace.getAll(response)

	// Read the tuple found.
	tuplesGet := <-response

	// Check the list of tuples found is equal to the local list of tuples.
	if !reflect.DeepEqual(tuplesLocal, tuplesGet) {
		t.Errorf("findAllTuples returned %q but was expected to return %q.", tuplesGet, tuplesLocal)
	}
}
*/
