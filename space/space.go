package space

import (
	"encoding/gob"
	. "github.com/luhac/gospace/protocol"
	. "github.com/luhac/gospace/shared"
	"strings"
	"sync"
)

// Space is a structure for interacting with tuple spaces.
type Space struct {
	id string
	p  PointToPoint
}

// NewSpace will create an empty tuple space with the specified URL.
// It uses a listener function in a go routine to listen for incoming messages.
func NewSpace(url string) Space {
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

	return Space{url, space}
}

// NewRemoteSpace connects to a remote space with the specified URL.
func NewRemoteSpace(url string) Space {
	// Register default structures for communication.
	gob.Register(Template{})
	gob.Register(Tuple{})
	gob.Register(TypeField{})

	space := CreatePointToPoint("whatever", "localhost", url)
	return Space{url, space}
}

func (s *Space) Put(tuple ...interface{}) bool {
	return Put((*s).p, tuple...)
}

func (s *Space) Get(template ...interface{}) bool {
	return Get((*s).p, template...)
}

func (s *Space) Query(template ...interface{}) bool {
	return Query((*s).p, template...)
}

func (s *Space) PutP(tuple ...interface{}) bool {
	return PutP((*s).p, tuple...)
}

func (s *Space) GetP(template ...interface{}) (bool, bool) {
	return GetP((*s).p, template...)
}

func (s *Space) QueryP(template ...interface{}) (bool, bool) {
	return QueryP((*s).p, template...)
}

func (s *Space) GetAll(template ...interface{}) ([]Tuple, bool) {
	return GetAll((*s).p, template...)
}

func (s *Space) QueryAll(template ...interface{}) ([]Tuple, bool) {
	return GetAll((*s).p, template...)
}

// Interspace defines the internal space interface.
// The content of this interface may change without notice.
type Interspace interface {
	Put(tuple ...interface{}) bool
	Get(template ...interface{}) bool
	Query(template ...interface{}) bool
	PutP(tuple ...interface{}) bool
	GetP(template ...interface{}) (bool, bool)
	QueryP(template ...interface{}) (bool, bool)
	GetAll(template ...interface{}) ([]Tuple, bool)
	QueryAll(template ...interface{}) ([]Tuple, bool)
}
