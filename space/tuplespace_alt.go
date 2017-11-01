package space

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/pspaces/gospace/container"
	"github.com/pspaces/gospace/function"
	"github.com/pspaces/gospace/policy"
	"github.com/pspaces/gospace/protocol"
	"github.com/pspaces/gospace/space/uri"
)

// Constants for logging errors occuring in this file.
var (
	tsAltBuf    bytes.Buffer
	tsAltLogger = log.New(&tsAltBuf, "gospace: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC|log.Lshortfile)
)

// tsAltLog logs errors occuring in this file.
func tsAltLog(fun interface{}, e *error) {
	if *e != nil {
		tsAltLogger.Printf("%s: %s\n", function.Name(fun), *e)
	}
}

// TODO: This is a part of hack, don't touch this.
// TODO: This can be removed once rearchitecting begins.
var localChanMap = new(sync.Map)

// NewSpaceAlt creates a representation of a new tuple space.
func NewSpaceAlt(url string, cp ...*policy.Composable) (ptp *protocol.PointToPoint, ts *TupleSpace) {
	registerTypes()

	u, err := uri.NewSpaceURI(url)

	if err == nil {
		// TODO: Exchange capabilities instead and
		// TODO: make a mechanism capable of doing that.
		if function.GlobalRegistry == nil {
			fr := function.NewRegistry()
			function.GlobalRegistry = &fr
		}
		funcReg := *function.GlobalRegistry

		// TODO: Create a better condition if a hosts name resolves to a local address.
		ips, err := net.LookupIP(u.Hostname())

		// Test if we are connecting locally to avoid TCP port issue.
		localhost := false
		for _, a := range ips {
			localhost = localhost || a.IsLoopback()
		}

		/*l, lerr := net.Listen("tcp4", ":"+u.Port())
		if lerr == nil {
			l.Close()
		}*/

		// TODO: Embrace the following hack, and remove it with architecting differently.
		connc := make(chan *net.Conn)
		val, exists := localChanMap.LoadOrStore(u, connc)

		//if exists {
		//close(connc)
		//connc = (val.(chan *net.Conn))
		//}

		if !exists && err == nil {
			muTuples := new(sync.RWMutex)
			muWaitingClients := new(sync.Mutex)
			tuples := []container.Tuple{}

			addr := strings.Join([]string{u.Hostname(), u.Port()}, ":")

			ts = &TupleSpace{
				muTuples:         muTuples,
				muWaitingClients: muWaitingClients,
				tuples:           tuples,
				pol:              nil,
				funReg:           &funcReg,
				port:             addr,
				connc:            connc,
			}

			if len(cp) == 1 {
				(*ts).pol = cp[0]
			}

			go ts.Listen()

			ptp = protocol.CreatePointToPoint(u.Space(), u.Hostname(), "0", connc, &funcReg)
		} else {
			for localhost && !exists {
				val, exists = localChanMap.Load(u)

				if exists {
					connc = (val.(chan *net.Conn))
					break
				}
			}

			var port string

			if !localhost {
				port = u.Port()
			} else {
				port = "0"
			}

			ptp = protocol.CreatePointToPoint(u.Space(), u.Hostname(), port, connc, &funcReg)
		}
	} else {
		ts = nil
		ptp = nil
	}

	return ptp, ts
}

// NewRemoteSpaceAlt creates a remote representation of a tuple space.
func NewRemoteSpaceAlt(url string, cp ...*policy.Composable) (ptp *protocol.PointToPoint, ts *TupleSpace) {
	registerTypes()

	u, err := uri.NewSpaceURI(url)

	if err == nil {
		// TODO: Exchange capabilities instead and
		// TODO: make a mechanism capable of doing that.
		if function.GlobalRegistry == nil {
			fr := function.NewRegistry()
			function.GlobalRegistry = &fr
		}
		funcReg := *function.GlobalRegistry

		// TODO: Create a better condition if a hosts name resolves to a local address.
		ips, _ := net.LookupIP(u.Hostname())

		// Test if we are connecting locally to avoid TCP port issue.
		localhost := false
		for _, a := range ips {
			localhost = localhost || a.IsLoopback()
		}

		var connc chan *net.Conn

		for localhost {
			val, exists := localChanMap.Load(u)

			if exists {
				connc = (val.(chan *net.Conn))
				break
			}
		}

		var port string

		if !localhost {
			port = u.Port()
		} else {
			port = "0"
		}

		// NOTE: It is not possible to connect to localhost as a remote space
		// NOTE: or if the port is taken.

		ptp = protocol.CreatePointToPoint(u.Space(), u.Hostname(), port, connc, &funcReg)
	} else {
		ts = nil
		ptp = nil
	}

	return ptp, ts
}

// registerTypes registers all the types necessary for the implementation.
func registerTypes() {
	gob.Register(container.Label{})
	gob.Register(container.Labels{})
	gob.Register(container.Template{})
	gob.Register(container.Tuple{})
	gob.Register(container.TypeField{})
	gob.Register([]interface{}{})
}

// Size will open a TCP connection to the PointToPoint and request the size of the tuple space.
func Size(ptp protocol.PointToPoint) (sz int, b bool) {
	var conn *net.Conn
	var err error

	defer tsAltLog(Size, &err)

	sz = -1
	b = false

	conn, err = establishConnection(ptp)

	if err != nil {
		return sz, b
	}

	defer (*conn).Close()

	err = sendMessage(conn, protocol.SizeRequest, "")

	if err != nil {
		return sz, b
	}

	sz, err = receiveMessageInt(conn)

	if err != nil {
		sz = -1
	} else {
		b = true
	}

	return sz, b
}

// Put will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation and tuple specified by the user.
// The method returns a boolean to inform if the operation was carried out with
// success or not.
func Put(ptp protocol.PointToPoint, tupleFields ...interface{}) (t container.Tuple, b bool) {
	var conn *net.Conn
	var err error

	defer tsAltLog(Put, &err)

	b = false

	t = container.NewTuple(tupleFields...)

	// Never time out and block until connection will be established.
	conn, err = establishConnection(ptp)

	// TODO: Yes this is a bad idea, and we are doing it for now until semantics
	// TODO: for what it means to block is established.
	for {
		if err == nil {
			break
		}

		if *conn != nil {
			(*conn).Close()
		}

		conn, err = establishConnection(ptp)
	}

	defer (*conn).Close()

	funcEncode(ptp.GetRegistry(), &t)
	defer funcDecode(ptp.GetRegistry(), &t)

	err = sendMessage(conn, protocol.PutRequest, t)

	if err != nil {
		return container.NewTuple(nil), b
	}

	b, err = receiveMessageBool(conn)

	if err != nil {
		b = false
		return container.NewTuple(nil), b
	}

	return t, b
}

// PutP will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation and tuple specified by the user.
// As the method is nonblocking it wont wait for a response whether or not the
// operation was successful.
// The method returns a boolean to inform if the operation was carried out with
// any errors with communication.
func PutP(ptp protocol.PointToPoint, tupleFields ...interface{}) (t container.Tuple, b bool) {
	var conn *net.Conn
	var err error

	defer tsAltLog(PutP, &err)

	b = false

	t = container.NewTuple(tupleFields...)

	conn, err = establishConnection(ptp)

	if err != nil {
		return container.NewTuple(nil), b
	}

	defer (*conn).Close()

	funcEncode(ptp.GetRegistry(), &t)
	defer funcDecode(ptp.GetRegistry(), &t)

	err = sendMessage(conn, protocol.PutPRequest, t)

	if err != nil {
		return container.NewTuple(nil), b
	}

	b = true

	return t, b
}

// Get will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation and template specified by the user.
// The method returns a boolean to inform if the operation was carried out with
// any errors with communication.
func Get(ptp protocol.PointToPoint, tempFields ...interface{}) (t container.Tuple, b bool) {
	t, b = getAndQuery(ptp, protocol.GetRequest, tempFields...)
	return t, b
}

// Query will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation and template specified by the user.
// The method returns a boolean to inform if the operation was carried out with
// any errors with communication.
func Query(ptp protocol.PointToPoint, tempFields ...interface{}) (t container.Tuple, b bool) {
	t, b = getAndQuery(ptp, protocol.QueryRequest, tempFields...)
	return t, b
}

func getAndQuery(ptp protocol.PointToPoint, operation string, tempFields ...interface{}) (t container.Tuple, b bool) {
	var conn *net.Conn
	var err error

	defer tsAltLog(getAndQuery, &err)

	b = false

	tp := container.NewTemplate(tempFields...)

	// Never time out and block until connection will be established.
	conn, err = establishConnection(ptp)

	// Busy loop until a successful connection is established.
	// TODO: Yes this is a bad idea, and we are doing it for now until semantics
	// TODO: for what it means to block is established.
	for {
		if err == nil {
			break
		}

		if *conn != nil {
			(*conn).Close()
		}

		conn, err = establishConnection(ptp)
	}

	defer (*conn).Close()

	funcEncode(ptp.GetRegistry(), &tp)
	defer funcDecode(ptp.GetRegistry(), &t)

	err = sendMessage(conn, operation, tp)

	if err != nil {
		return container.NewTuple(nil), b
	}

	t, err = receiveMessageTuple(conn)

	if err != nil {
		return container.NewTuple(nil), b
	}

	b = true

	funcDecode(ptp.GetRegistry(), &t)

	return t, b
}

// GetP will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation and template specified by the user.
// The function will return two bool values. The first denotes if a tuple was
// found, the second if there were any erors with communication.
func GetP(ptp protocol.PointToPoint, tempFields ...interface{}) (container.Tuple, bool, bool) {
	return getPAndQueryP(ptp, protocol.GetPRequest, tempFields...)
}

// QueryP will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation and template specified by the user.
// The function will return two bool values. The first denotes if a tuple was
// found, the second if there were any erors with communication.
func QueryP(ptp protocol.PointToPoint, tempFields ...interface{}) (container.Tuple, bool, bool) {
	return getPAndQueryP(ptp, protocol.QueryPRequest, tempFields...)
}

func getPAndQueryP(ptp protocol.PointToPoint, operation string, tempFields ...interface{}) (t container.Tuple, tb bool, sb bool) {
	var conn *net.Conn
	var err error

	defer tsAltLog(getPAndQueryP, &err)

	tb = false
	sb = false

	tp := container.NewTemplate(tempFields...)

	conn, err = establishConnection(ptp)

	if err != nil {
		return container.NewTuple(nil), tb, sb
	}

	defer (*conn).Close()

	funcEncode(ptp.GetRegistry(), &tp)
	defer funcDecode(ptp.GetRegistry(), &tp)

	err = sendMessage(conn, operation, tp)

	if err != nil {
		return container.NewTuple(nil), tb, sb
	}

	tb, t, err = receiveMessageBoolAndTuple(conn)

	if err != nil {
		return container.NewTuple(nil), tb, sb
	}

	sb = true

	return t, tb, sb
}

// GetAll will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation specified by the user.
// The method is nonblocking and will return all tuples found in the tuple
// space as well as a bool to denote if there were any errors with the
// communication.
func GetAll(ptp protocol.PointToPoint, tempFields ...interface{}) (ts []container.Tuple, b bool) {
	ts, b = getAllAndQueryAll(ptp, protocol.GetAllRequest, tempFields...)
	return ts, b
}

// QueryAll will open a TCP connection to the PointToPoint and send the message,
// which includes the type of operation specified by the user.
// The method is nonblocking and will return all tuples found in the tuple
// space as well as a bool to denote if there were any errors with the
// communication.
func QueryAll(ptp protocol.PointToPoint, tempFields ...interface{}) (ts []container.Tuple, b bool) {
	ts, b = getAllAndQueryAll(ptp, protocol.QueryAllRequest, tempFields...)
	return ts, b
}

func getAllAndQueryAll(ptp protocol.PointToPoint, operation string, tempFields ...interface{}) (ts []container.Tuple, b bool) {
	var conn *net.Conn
	var err error

	defer tsAltLog(getAllAndQueryAll, &err)

	ts = []container.Tuple{}
	b = false

	tp := container.NewTemplate(tempFields...)

	conn, err = establishConnection(ptp)

	if err != nil {
		return ts, b
	}

	defer (*conn).Close()

	funcEncode(ptp.GetRegistry(), &tp)

	err = sendMessage(conn, operation, tp)

	if err != nil {
		return ts, b
	}

	ts, err = receiveMessageTupleList(conn)

	if err != nil {
		return ts, b
	}

	b = true

	for _, t := range ts {
		funcDecode(ptp.GetRegistry(), &t)
	}

	return ts, b
}

// PutAgg will connect to a space and aggregate on matched tuples from the space according to a template.
// This method is nonblocking and will return a tuple and boolean state to denote if there were any errors
// with the communication. The tuple returned is the aggregation of tuples in the space.
// If no tuples are found it will create and put a new tuple from the template itself.
func PutAgg(ptp protocol.PointToPoint, fun interface{}, tempFields ...interface{}) (t container.Tuple, b bool) {
	t, b = aggOperation(ptp, protocol.PutAggRequest, fun, tempFields...)
	return t, b
}

// GetAgg will connect to a space and aggregate on matched tuples from the space.
// This method is nonblocking and will return a tuple perfomed by the aggregation
// as well as a boolean state to denote if there were any errors with the communication.
// The resulting tuple is empty if no matching occurs or the aggregation function can not aggregate the matched tuples.
func GetAgg(ptp protocol.PointToPoint, fun interface{}, tempFields ...interface{}) (t container.Tuple, b bool) {
	t, b = aggOperation(ptp, protocol.GetAggRequest, fun, tempFields...)
	return t, b
}

// QueryAgg will connect to a space and aggregated on matched tuples from the space.
// which includes the type of operation specified by the user.
// The method is nonblocking and will return a tuple found by aggregating the matched typles.
// The resulting tuple is empty if no matching occurs or the aggregation function can not aggregate the matched tuples.
func QueryAgg(ptp protocol.PointToPoint, fun interface{}, tempFields ...interface{}) (t container.Tuple, b bool) {
	t, b = aggOperation(ptp, protocol.QueryAggRequest, fun, tempFields...)
	return t, b
}

func aggOperation(ptp protocol.PointToPoint, operation string, fun interface{}, tempFields ...interface{}) (t container.Tuple, b bool) {
	var conn *net.Conn
	var err error

	defer tsAltLog(aggOperation, &err)

	t = container.NewTuple()
	b = false

	fields := make([]interface{}, len(tempFields)+1)
	fields[0] = fun
	copy(fields[1:], tempFields)
	tp := container.NewTemplate(fields...)

	conn, err = establishConnection(ptp)

	if err != nil {
		return t, b
	}

	defer (*conn).Close()

	funcEncode(ptp.GetRegistry(), &tp)

	err = sendMessage(conn, operation, tp)

	if err != nil {
		return t, b
	}

	t, err = receiveMessageTuple(conn)

	if err != nil {
		return t, b
	}

	b = true

	funcDecode(ptp.GetRegistry(), &t)

	return t, b
}

// establishConnection will establish a connection to the PointToPoint ptp and
// return the Conn and error.
func establishConnection(ptp protocol.PointToPoint, timeout ...time.Duration) (*net.Conn, error) {
	var conn net.Conn
	var err error

	addr := ptp.GetAddress()

	host, _, err := net.SplitHostPort(addr)

	proto := "tcp4"

	ips, err := net.LookupIP(host)

	if err == nil {
		// Test if we are connecting locally to avoid TCP port issue.
		localhost := false
		for _, a := range ips {
			localhost = localhost || a.IsLoopback()
		}

		connc := ptp.GetConnectionChannel()

		if localhost && connc != nil {
			r, w := net.Pipe()
			conn = r
			(*connc) <- &w
		} else {
			if len(timeout) == 0 {
				conn, err = net.Dial(proto, addr)
			} else {
				conn, err = net.DialTimeout(proto, addr, timeout[0])
			}
		}
	}

	return &conn, err
}

func sendMessage(conn *net.Conn, operation string, t interface{}) (err error) {
	gob.Register(t)
	gob.Register(container.TypeField{})

	enc := gob.NewEncoder(*conn)

	message := protocol.CreateMessage(operation, t)

	err = enc.Encode(message)

	return err
}

func receiveMessageBool(conn *net.Conn) (b bool, err error) {
	dec := gob.NewDecoder(*conn)

	err = dec.Decode(&b)

	return b, err
}

func receiveMessageInt(conn *net.Conn) (i int, err error) {
	dec := gob.NewDecoder(*conn)

	err = dec.Decode(&i)

	return i, err
}

func receiveMessageTuple(conn *net.Conn) (t container.Tuple, err error) {
	dec := gob.NewDecoder(*conn)

	err = dec.Decode(&t)

	return t, err
}

func receiveMessageBoolAndTuple(conn *net.Conn) (b bool, t container.Tuple, err error) {
	dec := gob.NewDecoder(*conn)

	var result []interface{}
	err = dec.Decode(&result)

	b = result[0].(bool)
	t = result[1].(container.Tuple)

	return b, t, err
}

func receiveMessageTupleList(conn *net.Conn) (ts []container.Tuple, err error) {
	dec := gob.NewDecoder(*conn)

	err = dec.Decode(&ts)

	return ts, err
}
