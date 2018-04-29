package uri

import (
	"testing"
)

func TestSpaceURI(t *testing.T) {
	rawuri := createTestRawURI()

	uri, err := NewSpaceURI(rawuri)

	valid := true
	if err != nil {
		valid = false
		t.Errorf("NewSpaceURI() experienced an error when parsing")
	} else {
		valid = valid && uri.Hostname() == "host"
		valid = valid && uri.Port() == "0"
		valid = valid && uri.Path() == "/space_name"
		valid = valid && uri.Mode() == "CONN"
		valid = valid && uri.Scheme() == "scheme"
		valid = valid && uri.Space() == "space_name"
	}

	if !valid {
		t.Errorf("NewSpaceURI returned %v, expected %v", uri, rawuri)
	}
}

func createTestRawURI() string {
	return "scheme://host:0/space_name?CONN"
}
