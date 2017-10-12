package protocol

import (
	. "github.com/luhac/gospace/shared"
)

// WaitingClient is used as a structure for clients who performed an
// unsuccessful Get or Query operation, in the sense that it didn't initially
// found a tuple in the tuple space. This structure stores the necessary
// information about the client, so that a tuple can be send to it once a
// matching tuple arrives.
type WaitingClient struct {
	template     Template      // Template that the waiting client is using to search for a tuple.
	responseChan chan<- *Tuple // Channel where the response can be send through.
	operation    string        // String that will denote the type of operation the client is trying to carry out.
}

// CreateWaitingClient will create the waiting client with the template that
// should be used for tuple matching, response channel for the matched tuple to
// be send to. The remove value will be used to determine if the client
// performed a Get or Query operation.
func CreateWaitingClient(temp Template, tupleChan chan<- *Tuple, remove bool) WaitingClient {
	var o string
	if remove {
		o = GetRequest
	} else {
		o = QueryRequest
	}

	waitingClient := WaitingClient{template: temp, responseChan: tupleChan, operation: o}
	return waitingClient
}

// GetTemplate will return the template of the waiting client.
func (waitingClient *WaitingClient) GetTemplate() Template {
	return waitingClient.template
}

// GetResponseChan will return the response channel of the waiting client.
func (waitingClient *WaitingClient) GetResponseChan() chan<- *Tuple {
	return waitingClient.responseChan
}

// GetOperation will return the operation of the waiting client.
func (waitingClient *WaitingClient) GetOperation() string {
	return waitingClient.operation
}
