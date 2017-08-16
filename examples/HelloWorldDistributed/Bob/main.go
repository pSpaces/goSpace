package main

import "goSpace/goSpace"

func main() {

	// Create a space S
	inbox := goSpace.RemoteSpace("8080")

	// Put a message in the space
	goSpace.Put(inbox, "Hello World!")

}
