package goSpace

// Message is the package that is send across a connection.
// It contains the type of message, denoted by operation and either a tuple or
// template, depending on the type of operation.
type Message struct {
	Operation string
	T         interface{}
}

// CreateMessage will create the message and return it with the opertaion type
// and tuple or template specified by the user.
func CreateMessage(operation string, t interface{}) Message {
	message := Message{Operation: operation, T: t}
	return message
}

// GetOperation will return the operation of the message.
func (message *Message) GetOperation() string {
	return message.Operation
}

// GetBody will return the body of the message, which can be a template or a
// tuple.
func (message *Message) GetBody() interface{} {
	return message.T
}
