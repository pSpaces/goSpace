package topology

import (
	"strconv"
	"strings"
)

// PointToPoint contains information about the receiver, being a user specified
// name, the IP address and the port number.
type PointToPoint struct {
	name    string // Name of receiver.
	address string // IP address and port number of receiver separated by ":".
}

// CreatePointToPoint will concatenate the ip and the port to a string to create
// an address of the receiver. The created PointToPoint is then returned.
func CreatePointToPoint(name string, ip string, port int) PointToPoint {
	address := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	pointToPoint := PointToPoint{name, address}
	return pointToPoint
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
