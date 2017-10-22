package shared

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Mode defines the connection mode to a resource.
type Mode int

const (
	ConnKeep Mode = iota
	ConnOnce
	ConnPush
	ConnPull
)

// modeName is mapping from connection modes to their string values.
var modeName = map[Mode]string{
	ConnKeep: "KEEP",
	ConnOnce: "CONN",
	ConnPush: "PUSH",
	ConnPull: "PULL",
}

// SpaceURI is a structure for containing information about a resource location.
type SpaceURI struct {
	scheme string
	host   string
	port   string
	path   string
	mode   string
}

// NewSpaceURI creates a new URI location to a resource.
// NewSpaceURI follows a restricted URI scheme defined by pSpace specification.
func NewSpaceURI(rawurl string) (su *SpaceURI, e error) {
	uri, err := url.Parse(rawurl)

	su = &SpaceURI{}

	if err != nil || (uri.Hostname() == "" && uri.Path == "" && uri.Scheme == "") {
		return nil, err
	}

	if uri.Hostname() == "" && uri.Path != "" && uri.Scheme == "" {
		(*su).host = "localhost"
		(*su).scheme = "tcp"
	} else {
		(*su).host = uri.Hostname()
		(*su).scheme = uri.Scheme
	}

	if uri.Port() != "" {
		(*su).port = uri.Port()
	} else {
		(*su).port = "31415"
	}

	if uri.Path != "" {
		(*su).path = strings.Join([]string{"/", strings.Trim(uri.EscapedPath(), "/")}, "")
	} else {
		(*su).path = uri.EscapedPath()
	}

	modes := strings.Join([]string{"(?i)((",
		modeName[ConnKeep], ")|(",
		modeName[ConnOnce], ")|(",
		modeName[ConnPush], ")|(",
		modeName[ConnPull], "))"}, "")
	modere, _ := regexp.Compile(modes)
	mode := modere.FindString(uri.RawQuery)

	if mode != "" {
		(*su).mode = mode
	} else {
		(*su).mode = modeName[ConnKeep]
	}

	return su, err
}

// Hostname returns the hostname contained in the URI.
func (su *SpaceURI) Hostname() (hostname string) {
	return (*su).host
}

// Mode returns the connection mode contained in the URI.
func (su *SpaceURI) Mode() (mode string) {
	return (*su).mode
}

// Path returns the path contained in the URI.
func (su *SpaceURI) Path() (path string) {
	return (*su).path
}

// Port returns the port contained in the URI.
func (su *SpaceURI) Port() (port string) {
	return (*su).port
}

// Scheme returns the scheme contained in the URI.
func (su *SpaceURI) Scheme() (hostname string) {
	return (*su).scheme
}

// Space returns the space name contained in the URI.
func (su *SpaceURI) Space() (hostname string) {
	return strings.TrimLeft(su.Path(), "/")
}

// String returns a print friendly representation of the URI.
func (su SpaceURI) String() (str string) {
	return fmt.Sprintf("%s://%s:%s%s?%s", (&su).Scheme(), (&su).Hostname(), (&su).Port(), (&su).Path(), (&su).Mode())
}
