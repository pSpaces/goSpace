package space

import (
	"encoding/gob"
	. "github.com/luhac/gospace/protocol"
	. "github.com/luhac/gospace/shared"
	"strings"
	"sync"
)

// NewSpace will create the tuple space with the specified url and
// initialise the lock.
// It will run the listen function in a go routine to listen for incoming
// messages.
// The list of tuples is initially empty.
// It returns a pointToPoint
func NewSpace(url string) PointToPoint {
	// Register default structures for communication.
	gob.Register(Template{})
	gob.Register(Tuple{})
	gob.Register(TypeField{})

	// For now, we only accept port numbers

	// The specified port number by the user needs to be converted to a string
	// with the following format ":<port>".
	ts := TupleSpace{muTuples: new(sync.RWMutex), muWaitingClients: new(sync.Mutex), port: strings.Join([]string{"", url}, ":")}
	go ts.Listen()
	space := CreatePointToPoint("whatever", "localhost", url)
	return space
}

// RemoteSpace is like CreateSpace but instead of creating a space
// It connects to a remote one
func RemoteSpace(url string) PointToPoint {
	// Register default structures for communication.
	gob.Register(Template{})
	gob.Register(Tuple{})
	gob.Register(TypeField{})

	space := CreatePointToPoint("whatever", "localhost", url)
	return space
}

type SpaceInterface interface {
	Put(ptp PointToPoint, tupleFields ...interface{}) bool
	Get(ptp PointToPoint, tempFields ...interface{}) bool
	Query(ptp PointToPoint, tempFields ...interface{}) bool
	PutP(ptp PointToPoint, tupleFields ...interface{}) bool
	GetP(ptp PointToPoint, tempFields ...interface{}) (bool, bool)
	QueryP(ptp PointToPoint, tempFields ...interface{}) (bool, bool)
	GetAll(ptp PointToPoint, tempFields ...interface{}) ([]Tuple, bool)
	QueryAll(ptp PointToPoint, tempFields ...interface{}) ([]Tuple, bool)
}
