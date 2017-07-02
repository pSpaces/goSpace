package tuplespace

import (
	"encoding/gob"
	"fmt"
	"goSpace/goSpace/constants"
	"goSpace/goSpace/topology"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

// TupleSpace contains a set of tuples and it has a mutex lock associated with
// it to secure mutual exclusion.
// Furthermore a port number to locate it.
type TupleSpace struct {
	muTuples         *sync.RWMutex   // Lock for the tuples[].
	muWaitingClients *sync.Mutex     // Lock for the waitingClients[].
	tuples           []Tuple         // Tuples in the tuple space.
	port             string          // Port number of the tuple space.
	waitingClients   []WaitingClient // Structure for clients that couldn't initially find a matching tuple.
}

// CreateTupleSpace will create the tuple space with the specified port and
// initialise the lock.
// It will run the listen function in a go routine to listen for incoming
// messages.
// The tuple space ts is returned. The list of tuples is initially empty.
func CreateTupleSpace(port int) *TupleSpace {
	// Register default structures for communication.
	gob.Register(Template{})
	gob.Register(Tuple{})
	gob.Register(TypeField{})

	// The specified port number by the user needs to be converted to a string
	// with the following format ":<portNumber>".
	portNumber := strings.Join([]string{"", strconv.Itoa(port)}, ":")
	ts := TupleSpace{muTuples: new(sync.RWMutex), muWaitingClients: new(sync.Mutex), port: portNumber}
	go ts.listen()
	return &ts
}

// Size return the number of tuples in the tuple space.
func (ts *TupleSpace) Size() int {
	return len(ts.tuples)
}

// put will call the nonblocking put method and places a success response on
// the response channel.
func (ts *TupleSpace) put(t *Tuple, response chan<- bool) {
	ts.putP(t)
	response <- true
}

// putP will put a lock on the tuple space add the tuple to the list of
// tuples and unlock the list.
func (ts *TupleSpace) putP(t *Tuple) {

	ts.muWaitingClients.Lock()
	// Check if someone is waiting for the tuple that is about to be placed.
	for i := 0; i < len(ts.waitingClients); i++ {
		waitingClient := ts.waitingClients[i]
		// Extract the template from the waiting client and check if it
		// matches the tuple.
		temp := waitingClient.GetTemplate()
		if t.match(temp) {
			// If this is reached, the tuple matched the template and the
			// tuple is send to the response channel of the waiting client.
			clientResponse := waitingClient.GetResponseChan()
			clientResponse <- t
			// Check if the client who was waiting for the tuple performed a get
			// or query operation.
			ts.removeClientAt(i)
			i--
			clientOperation := waitingClient.GetOperation()
			if clientOperation == constants.GetRequest {
				// Unlock before exiting the method.
				ts.muWaitingClients.Unlock()
				return
			}
		}
	}
	// No waiting client performing Get matched the tuple. So unlock.
	ts.muWaitingClients.Unlock()

	// Place lock on tuples[] before adding the new tuple.
	ts.muTuples.Lock()
	defer ts.muTuples.Unlock()
	ts.tuples = append(ts.tuples, *t)
}

func (ts *TupleSpace) removeClientAt(i int) {
	ts.waitingClients = append(ts.waitingClients[:i], ts.waitingClients[i+1:]...)
}

// get will find the first tuple that matches the template temp and remove the
// tuple from the tuple space.
func (ts *TupleSpace) get(temp Template, response chan<- *Tuple) {
	ts.findTupleBlocking(temp, response, true)
}

// query will find the first tuple that matches the template temp.
func (ts *TupleSpace) query(temp Template, response chan<- *Tuple) {
	ts.findTupleBlocking(temp, response, false)
}

// findTupleBlocking will continuously search for a tuple that matches the
// template temp, making the method blocking.
// The boolean remove will denote if the found tuple should be removed or not
// from the tuple space.
// The found tuple is written to the channel response.
func (ts *TupleSpace) findTupleBlocking(temp Template, response chan<- *Tuple, remove bool) {
	tuple := new(Tuple)
	// Seach for the a tuple in the tuple space.
	tuple = ts.findTuple(temp, remove)

	// Check if there was a tuple matching the template in the tuple space.
	if tuple != nil {
		// There was a tuple that matched the template. Write it to the
		// channel and return.
		response <- tuple
		return
	}
	// There was no tuple matching the template. Enter sleep.
	newWaitingClient := CreateWaitingClient(temp, response, remove)
	ts.addNewClient(newWaitingClient)
	return
}

// addNewClient will add the client to the list of waiting clients.
func (ts *TupleSpace) addNewClient(client WaitingClient) {
	ts.muWaitingClients.Lock()
	defer ts.muWaitingClients.Unlock()
	ts.waitingClients = append(ts.waitingClients, client)
}

// getP will find the first tuple that matches the template temp and remove the
// tuple from the tuple space.
func (ts *TupleSpace) getP(temp Template, response chan<- *Tuple) {
	ts.findTupleNonblocking(temp, response, true)
}

// queryP will find the first tuple that matches the template temp.
func (ts *TupleSpace) queryP(temp Template, response chan<- *Tuple) {
	ts.findTupleNonblocking(temp, response, false)
}

// findTupleNonblocking will search for a tuple that matches the template temp.
// The boolean remove will denote if the tuple should be removed or not from
// the tuple space.
// The found tuple is written to the channel response.
func (ts *TupleSpace) findTupleNonblocking(temp Template, response chan<- *Tuple, remove bool) {
	tuplePtr := ts.findTuple(temp, remove)
	response <- tuplePtr
}

// findTuple will run through the tuple space to see if it contains a tuple that
// matches the template temp.
// A lock is placed around the shared data.
// The boolean remove will denote if the tuple should be removed or not from
// the tuple space.
// If a match is found a pointer to the tuple is returned, otherwise nil is.
func (ts *TupleSpace) findTuple(temp Template, remove bool) *Tuple {
	if remove {
		ts.muTuples.Lock()
		defer ts.muTuples.Unlock()
	} else {
		ts.muTuples.RLock()
		defer ts.muTuples.RUnlock()
	}

	for i, t := range ts.tuples {
		if t.match(temp) {
			if remove {
				ts.removeTupleAt(i)
			}
			return &t
		}
	}
	return nil
}

// getAll will return and remove every tuple from the tuple space.
func (ts *TupleSpace) getAll(temp Template, response chan<- []Tuple) {
	ts.findAllTuples(temp, response, true)
}

// queryAll will return every tuple from the tuple space.
func (ts *TupleSpace) queryAll(temp Template, response chan<- []Tuple) {
	ts.findAllTuples(temp, response, false)
}

// findAllTuples will make a copy a the tuples in the tuple space to a list.
// The boolean remove will denote if the tuple should be removed or not from
// the tuple space.
// NOTE: an empty list of tuples is a legal return value.
func (ts *TupleSpace) findAllTuples(temp Template, response chan<- []Tuple, remove bool) {
	if remove {
		ts.muTuples.Lock()
		defer ts.muTuples.Unlock()
	} else {
		ts.muTuples.RLock()
		defer ts.muTuples.RUnlock()
	}
	var tuples []Tuple
	var removeIndex []int
	// Go through tuple space and collects matching tuples
	for i, t := range ts.tuples {
		if t.match(temp) {
			if remove {
				removeIndex = append(removeIndex, i)
			}
			tuples = append(tuples, t)
		}
	}
	// Remove tuples from tuple space if it is a get operations
	for i := len(removeIndex) - 1; i >= 0; i-- {
		ts.removeTupleAt(removeIndex[i])
	}

	response <- tuples
}

// clearTupleSpace will reinitialise the list of tuples in the tuple space.
func (ts *TupleSpace) clearTupleSpace() {
	ts.tuples = []Tuple{}
}

// removeTupleAt will removeTupleAt the tuple in the tuples space at index i.
func (ts *TupleSpace) removeTupleAt(i int) {
	//moves last tuple to place i, then removes last element from slice
	ts.tuples[i] = ts.tuples[ts.Size()-1]
	ts.tuples = ts.tuples[:ts.Size()-1]
}

// listen will listen and accept all incoming connections. Once a connection has
// been established, the connection is passed on to the handler.
func (ts *TupleSpace) listen() {

	// Create the listener to listen on the tuple space's port and using TCP.
	listener, errListen := net.Listen("tcp", ts.port)

	// Error check for creating listener.
	if errListen != nil {
		log.Fatal("Following error occured when creating the listener:", errListen)
	}

	defer listener.Close()

	// Infinite loop to accept all incoming connections.
	for {

		// Accept incoming connection and create a connection.
		conn, errAccept := listener.Accept()

		// Error check for accepting connection.
		if errAccept != nil {
			continue
		}

		// Pass on the connection to the handler.
		go ts.handle(conn)
	}
}

// handle will read and decode the message from the connection.
// The decoded message will be passed on to the respective method.
func (ts *TupleSpace) handle(conn net.Conn) {
	// Make sure the connection closes when method returns.
	defer conn.Close()

	// Create decoder to the connection to receive the message.
	dec := gob.NewDecoder(conn)

	// Read the message from the connection through the decoder.
	var message topology.Message
	errDec := dec.Decode(&message)

	// Error check for receiving message.
	if errDec != nil {
		return
	}

	operation := message.GetOperation()

	switch operation {
	case constants.PutRequest:
		// Body of message must be tuple.
		tuple := message.GetBody().(Tuple)
		ts.handlePut(conn, tuple)
	case constants.PutPRequest:
		// Body of message must be tuple.
		tuple := message.GetBody().(Tuple)
		ts.handlePutP(tuple)
	case constants.GetRequest:
		// Body of message must be template.
		template := message.GetBody().(Template)
		ts.handleGet(conn, template)
	case constants.GetPRequest:
		// Body of message must be template.
		template := message.GetBody().(Template)
		ts.handleGetP(conn, template)
	case constants.GetAllRequest:
		// Body of message must be empty.
		template := message.GetBody().(Template)
		ts.handleGetAll(conn, template)
	case constants.QueryRequest:
		// Body of message must be template.
		template := message.GetBody().(Template)
		ts.handleQuery(conn, template)
	case constants.QueryPRequest:
		// Body of message must be template.
		template := message.GetBody().(Template)
		ts.handleQueryP(conn, template)
	case constants.QueryAllRequest:
		// Body of message must be empty.
		template := message.GetBody().(Template)
		ts.handleQueryAll(conn, template)
	default:
		fmt.Println("Can't handle operation. Contact client at ", conn.RemoteAddr())
		return
	}
}

// handlePut is a blocking method.
// The method will place the tuple t in the tuple space ts.
// The method will send a boolean value to the connection conn to tell whether
// or not the placement succeeded
func (ts *TupleSpace) handlePut(conn net.Conn, t Tuple) {
	defer handleRecover()

	readChannel := make(chan bool)
	go ts.put(&t, readChannel)
	result := <-readChannel

	enc := gob.NewEncoder(conn)

	errEnc := enc.Encode(result)

	// NOTE: What happens here, if the tuple has been placed in the tuple
	// space, but something goes wrong in encoding the reponse?
	if errEnc != nil {
		panic("handlePut")
	}
}

// handlePutP is a nonblocking method.
// The method will try and place the tuple t in the tuple space ts.
func (ts *TupleSpace) handlePutP(t Tuple) {
	go ts.putP(&t)
}

// handleGet is a blocking method.
// It will find a tuple matching the template temp and return it.
func (ts *TupleSpace) handleGet(conn net.Conn, temp Template) {
	defer handleRecover()

	readChannel := make(chan *Tuple)
	go ts.get(temp, readChannel)
	resultTuplePtr := <-readChannel

	enc := gob.NewEncoder(conn)

	errEnc := enc.Encode(*resultTuplePtr)

	if errEnc != nil {
		panic("handleGet")
	}
}

// handleGetP is a nonblocking method.
// It will try to find a tuple matching the template temp and remove the tuple
// from the tuple space.
// As it may not find it, the method will send a boolean as well as the tuple
// to the connection conn.
func (ts *TupleSpace) handleGetP(conn net.Conn, temp Template) {
	defer handleRecover()

	readChannel := make(chan *Tuple)
	go ts.getP(temp, readChannel)
	resultTuplePtr := <-readChannel

	enc := gob.NewEncoder(conn)

	if resultTuplePtr == nil {
		result := []interface{}{false, Tuple{}}

		errEnc := enc.Encode(result)

		if errEnc != nil {
			panic("handleGetP")
		}
	} else {
		result := []interface{}{true, *resultTuplePtr}

		errEnc := enc.Encode(result)

		if errEnc != nil {
			panic("handleGetP")
		}
	}
}

// handleGetAll is a nonblocking method that will remove all tuples from the tuple
// space and send them in a list through the connection conn.
func (ts *TupleSpace) handleGetAll(conn net.Conn, temp Template) {
	defer handleRecover()

	readChannel := make(chan []Tuple)
	go ts.getAll(temp, readChannel)
	tupleList := <-readChannel

	enc := gob.NewEncoder(conn)

	errEnc := enc.Encode(tupleList)

	if errEnc != nil {
		panic("handleGetAll")
	}
}

// handleQuery is a blocking method.
// It will find a tuple matching the template temp.
// The found tuple will be send to the connection conn.
func (ts *TupleSpace) handleQuery(conn net.Conn, temp Template) {
	defer handleRecover()

	readChannel := make(chan *Tuple)
	go ts.query(temp, readChannel)
	resultTuplePtr := <-readChannel

	enc := gob.NewEncoder(conn)

	errEnc := enc.Encode(*resultTuplePtr)

	if errEnc != nil {
		panic("handleQuery")
	}
}

// handleQueryP is a nonblocking method.
// It will try to find a tuple matching the template temp.
// As it may not find it, the method returns a boolean as well as the tuple.
func (ts *TupleSpace) handleQueryP(conn net.Conn, temp Template) {
	defer handleRecover()

	readChannel := make(chan *Tuple)
	go ts.queryP(temp, readChannel)
	resultTuplePtr := <-readChannel

	enc := gob.NewEncoder(conn)

	if resultTuplePtr == nil {
		result := []interface{}{false, Tuple{}}

		errEnc := enc.Encode(result)

		if errEnc != nil {
			panic("handleQueryP")
		}
	} else {
		result := []interface{}{true, *resultTuplePtr}

		errEnc := enc.Encode(result)

		if errEnc != nil {
			panic("handleQueryP")
		}
	}
}

// handleQueryAll is a blocking method that will return all tuples from the tuple
// space in a list.
func (ts *TupleSpace) handleQueryAll(conn net.Conn, temp Template) {
	defer handleRecover()

	readChannel := make(chan []Tuple)
	go ts.queryAll(temp, readChannel)
	tupleList := <-readChannel

	enc := gob.NewEncoder(conn)

	errEnc := enc.Encode(tupleList)

	if errEnc != nil {
		panic("handleQueryAll")
	}
}

func handleRecover() {
	if error := recover(); error != nil {
		fmt.Println("Recovered from", error)
	}
}
