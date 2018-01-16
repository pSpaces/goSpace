package space

import (
	"crypto/tls"
	"encoding/gob"
	"fmt"
	"net"
	"strings"
	"sync"

	. "github.com/pspaces/gospace/protocol"
	. "github.com/pspaces/gospace/shared"
)

// NewSpaceAlt creates a representation of a new tuple space.
func NewSpaceAlt(url string, config *tls.Config) (ptp *PointToPoint, ts *TupleSpace) {
	registerTypes()

	uri, err := NewSpaceURI(url)

	if err == nil {
		// TODO: This is not the best way of doing it since
		// TODO: a host can resolve to multiple addresses.
		// TODO: For now, accept this limitation, and fix it soon.
		ts = &TupleSpace{muTuples: new(sync.RWMutex), muWaitingClients: new(sync.Mutex), port: strings.Join([]string{"", uri.Port()}, ":")}

		go ts.Listen(config)

		ptp = CreatePointToPoint(uri.Space(), "localhost", uri.Port())
	} else {
		ts = nil
		ptp = nil
	}

	return ptp, ts
}

// NewRemoteSpaceAlt creates a representaiton of a remote tuple space.
func NewRemoteSpaceAlt(url string) (ptp *PointToPoint, ts *TupleSpace) {
	registerTypes()

	uri, err := NewSpaceURI(url)

	if err == nil {
		// TODO: This is not the best way of doing it since
		// TODO: a host can resolve to multiple addresses.
		// TODO: For now, accept this limitation, and fix it soon.
		ptp = CreatePointToPoint(uri.Space(), uri.Hostname(), uri.Port())
	} else {
		ts = nil
		ptp = nil
	}

	return ptp, ts
}

// registerTypes registers all the types necessary for the implementation.
func registerTypes() {
	// Register default structures for communication.
	gob.Register(Template{})
	gob.Register(Tuple{})
	gob.Register(TypeField{})
	gob.Register([]interface{}{})
}

// Put will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation and tuple specified by the user.
// The method returns a boolean to inform if the operation was carried out with
// success or not.
func Put(ptp PointToPoint, tupleFields ...interface{}) bool {
	t := CreateTuple(tupleFields...)
	conn, errDial := establishConnection(ptp)

	// Error check for establishing connection.
	if errDial != nil {
		fmt.Println("ErrDial:", errDial)
		return false
	}

	// Make sure the connection closes when method returns.
	defer conn.Close()

	errSendMessage := sendMessage(conn, PutRequest, t)

	// Error check for sending message.
	if errSendMessage != nil {
		fmt.Println("ErrSendMessage:", errSendMessage)
		return false
	}

	b, errReceiveMessage := receiveMessageBool(conn)

	// Error check for receiving response.
	if errReceiveMessage != nil {
		fmt.Println("ErrReceiveMessage:", errReceiveMessage)
		return false
	}

	// Return result.
	return b
}

// PutP will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation and tuple specified by the user.
// As the method is nonblocking it wont wait for a response whether or not the
// operation was successful.
// The method returns a boolean to inform if the operation was carried out with
// any errors with communication.
func PutP(ptp PointToPoint, tupleFields ...interface{}) bool {
	t := CreateTuple(tupleFields...)
	conn, errDial := establishConnection(ptp)

	// Error check for establishing connection.
	if errDial != nil {
		fmt.Println("ErrDial:", errDial)
		return false
	}

	// Make sure the connection closes when method returns.
	defer conn.Close()

	errSendMessage := sendMessage(conn, PutPRequest, t)

	// Error check for sending message.
	if errSendMessage != nil {
		fmt.Println("ErrSendMessage:", errSendMessage)
		return false
	}

	return true
}

// Get will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation and template specified by the user.
// The method returns a boolean to inform if the operation was carried out with
// any errors with communication.
func Get(ptp PointToPoint, tempFields ...interface{}) bool {
	return getAndQuery(ptp, GetRequest, tempFields...)
}

// Query will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation and template specified by the user.
// The method returns a boolean to inform if the operation was carried out with
// any errors with communication.
func Query(ptp PointToPoint, tempFields ...interface{}) bool {
	return getAndQuery(ptp, QueryRequest, tempFields...)
}

func getAndQuery(ptp PointToPoint, operation string, tempFields ...interface{}) bool {
	t := CreateTemplate(tempFields...)
	conn, errDial := establishConnection(ptp)

	// Error check for establishing connection.
	if errDial != nil {
		fmt.Println("ErrDial:", errDial)
		return false
	}

	// Make sure the connection closes when method returns.
	defer conn.Close()

	errSendMessage := sendMessage(conn, operation, t)

	// Error check for sending message.
	if errSendMessage != nil {
		fmt.Println("ErrSendMessage:", errSendMessage)
		return false
	}

	tuple, errReceiveMessage := receiveMessageTuple(conn)

	// Error check for receiving response.
	if errReceiveMessage != nil {
		fmt.Println("ErrReceiveMessage:", errReceiveMessage)
		return false
	}

	tuple.WriteToVariables(tempFields...)

	// Return result.
	return true
}

// GetP will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation and template specified by the user.
// The function will return two bool values. The first denotes if a tuple was
// found, the second if there were any erors with communication.
func GetP(ptp PointToPoint, tempFields ...interface{}) (bool, bool) {
	return getPAndQueryP(ptp, GetPRequest, tempFields...)
}

// QueryP will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation and template specified by the user.
// The function will return two bool values. The first denotes if a tuple was
// found, the second if there were any erors with communication.
func QueryP(ptp PointToPoint, tempFields ...interface{}) (bool, bool) {
	return getPAndQueryP(ptp, QueryPRequest, tempFields...)
}

func getPAndQueryP(ptp PointToPoint, operation string, tempFields ...interface{}) (bool, bool) {
	t := CreateTemplate(tempFields...)
	conn, errDial := establishConnection(ptp)

	// Error check for establishing connection.
	if errDial != nil {
		fmt.Println("ErrDial:", errDial)
		return false, false
	}

	// Make sure the connection closes when method returns.
	defer conn.Close()

	errSendMessage := sendMessage(conn, operation, t)

	// Error check for sending message.
	if errSendMessage != nil {
		fmt.Println("ErrSendMessage:", errSendMessage)
		return false, false
	}

	b, tuple, errReceiveMessage := receiveMessageBoolAndTuple(conn)

	// Error check for receiving response.
	if errReceiveMessage != nil {
		fmt.Println("ErrReceiveMessage:", errReceiveMessage)
		return false, false
	}

	if b {
		tuple.WriteToVariables(tempFields...)
	}

	// Return result.
	return b, true
}

// GetAll will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation specified by the user.
// The method is nonblocking and will return all tuples found in the tuple
// space as well as a bool to denote if there were any errors with the
// communication.
// NOTE: tuples is allowed to be an empty list, implying the tuple space was
// empty.
func GetAll(ptp PointToPoint, tempFields ...interface{}) ([]Tuple, bool) {
	return getAllAndQueryAll(ptp, GetAllRequest, tempFields...)
}

// QueryAll will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation specified by the user.
// The method is nonblocking and will return all tuples found in the tuple
// space as well as a bool to denote if there were any errors with the
// communication.
// NOTE: tuples is allowed to be an empty list, implying the tuple space was
// empty.
func QueryAll(ptp PointToPoint, tempFields ...interface{}) ([]Tuple, bool) {
	return getAllAndQueryAll(ptp, QueryAllRequest, tempFields...)
}

func getAllAndQueryAll(ptp PointToPoint, operation string, tempFields ...interface{}) ([]Tuple, bool) {
	t := CreateTemplate(tempFields...)
	conn, errDial := establishConnection(ptp)

	// Error check for establishing connection.
	if errDial != nil {
		fmt.Println("ErrDial:", errDial)
		return []Tuple{}, false
	}

	// Make sure the connection closes when method returns.
	defer conn.Close()

	// Initiallise dummy tuple.
	// TODO: Get rid of the dummy tuple.
	errSendMessage := sendMessage(conn, operation, t)

	// Error check for sending message.
	if errSendMessage != nil {
		fmt.Println("ErrSendMessage:", errSendMessage)
		return []Tuple{}, false
	}

	tuples, errReceiveMessage := receiveMessageTupleList(conn)

	// Error check for receiving response.
	if errReceiveMessage != nil {
		fmt.Println("ErrReceiveMessage:", errReceiveMessage)
		return []Tuple{}, false
	}

	// Return result.
	return tuples, true
}

// establishConnection will establish a connection to the PointToPoint ptp and
// return the Conn and error.
func establishConnection(ptp PointToPoint) (net.Conn, error) {
	addr := ptp.GetAddress()

	// Establish a connection to the PointToPoint using TCP to ensure reliability.
	conn, errDial := net.Dial("tcp4", addr)

	return conn, errDial
}

func sendMessage(conn net.Conn, operation string, t interface{}) error {
	// Create encoder to the connection.
	enc := gob.NewEncoder(conn)

	// Register the type of t for Encode to handle it.
	gob.Register(t)
	//registrer typefield to match types
	gob.Register(TypeField{})

	// Generate the message.
	message := CreateMessage(operation, t)

	// Sends the message to the connection through the encoder.
	errEnc := enc.Encode(message)

	return errEnc
}

func receiveMessageBool(conn net.Conn) (bool, error) {
	// Create decoder to the connection to receive the response.
	dec := gob.NewDecoder(conn)

	// Read the response from the connection through the decoder.
	var b bool
	errDec := dec.Decode(&b)

	return b, errDec
}

func receiveMessageTuple(conn net.Conn) (Tuple, error) {
	// Create decoder to the connection to receive the response.
	dec := gob.NewDecoder(conn)

	// Read the response from the connection through the decoder.
	var tuple Tuple
	errDec := dec.Decode(&tuple)

	return tuple, errDec
}

func receiveMessageBoolAndTuple(conn net.Conn) (bool, Tuple, error) {
	// Create decoder to the connection to receive the response.
	dec := gob.NewDecoder(conn)

	// Read the response from the connection through the decoder.
	var result []interface{}
	errDec := dec.Decode(&result)

	// Extract the boolean and tuple from the result.
	b := result[0].(bool)
	tuple := result[1].(Tuple)

	return b, tuple, errDec
}

func receiveMessageTupleList(conn net.Conn) ([]Tuple, error) {
	// Create decoder to the connection to receive the response.
	dec := gob.NewDecoder(conn)

	// Read the response from the connection through the decoder.
	var tuples []Tuple
	errDec := dec.Decode(&tuples)

	return tuples, errDec
}
