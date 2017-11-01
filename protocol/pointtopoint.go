package protocol

import (
	"net"
	"strings"

	"github.com/pspaces/gospace/function"
)

// PointToPoint contains information about the receiver, being a user specified
// name, the IP address and the port number.
type PointToPoint struct {
	name    string             // Name of receiver.
	address string             // IP address and port number of receiver separated by ":".
	connc   chan *net.Conn     // Active connection channel.
	funReg  *function.Registry // Function registry.
}

// CreatePointToPoint will concatenate the ip and the port to a string to create
// an address of the receiver. The created PointToPoint is then returned.
// One can optionally pass an active connection channel and a function registry.
func CreatePointToPoint(name string, ip string, port string, connc chan *net.Conn, fr *function.Registry) (ptp *PointToPoint) {
	address := strings.Join([]string{ip, port}, ":")
	ptp = &PointToPoint{name: name, address: address, connc: connc, funReg: fr}
	return ptp
}

// GetConnectionChannel gets the connection of the PointToPoint structure.
func (ptp *PointToPoint) GetConnectionChannel() (connc *chan *net.Conn) {
	b := ptp != nil

	if b {
		connc = &(*ptp).connc
	}

	return connc
}

// SetConnectionChannel sets the connection of the PointToPoint structure.
func (ptp *PointToPoint) SetConnectionChannel(connc *chan *net.Conn) (b bool) {
	b = ptp != nil && connc != nil

	if b {
		(*ptp).connc = *connc
	}

	return b
}

// ToString will combine the name and address of the PointToPoint in a readable
// string and return it.
func (ptp *PointToPoint) ToString() string {
	sName := strings.Join([]string{"Name", ptp.name}, ": ")
	sAddress := strings.Join([]string{"address", ptp.address}, ": ")

	s := strings.Join([]string{sName, sAddress}, ", ")

	return s
}

// GetAddress will return the address of the PointToPoint.
func (ptp *PointToPoint) GetAddress() string {
	return ptp.address
}

// GetName will return the name of the PointToPoint.
func (ptp *PointToPoint) GetName() string {
	return ptp.name
}

// GetRegistry will return the function registry associated to ptp.
func (ptp *PointToPoint) GetRegistry() (fr *function.Registry) {
	return ptp.funReg
}
