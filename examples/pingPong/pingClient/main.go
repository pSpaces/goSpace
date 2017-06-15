package main

import (
	"fmt"
	"goSpace/goSpace/topology"
	"goSpace/goSpace/tuplespace"
	"os"
	"strconv"
	"time"
)

func main() {
	// Get the port number that the user which to run the application on.
	pingPort, err := strconv.Atoi(os.Args[1])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	// Get the IP address and port number of the client running the pong
	// application.
	pongIP := os.Args[2]
	pongPort, err := strconv.Atoi(os.Args[3])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}

	tsPtr := tuplespace.CreateTupleSpace(pingPort)
	fmt.Println(tsPtr)

	ownPtP, theirPtP := createPtP(pingPort, pongIP, pongPort)

	ping(ownPtP, theirPtP)
}

// Ping will initially act as client
func ping(ownPtP topology.PointToPoint, theirPtP topology.PointToPoint) {
	// Create template to find the "Pong" tuple in own tuple space.
	//template := tuplespace.CreateTemplate("Pong")
	// Create "Ping" tuple to send to their tuple space.
	//tuple := tuplespace.CreateTuple("Ping")

	// Initialise the ping-pong by sending a "Ping" tuple to the pong
	// application's tuple space.
	tuplespace.Put(theirPtP, "Ping")
	for {
		// Find a "Pong" tuple in own tuple space
		tuplespace.Get(ownPtP, "Pong")
		// A "Pong" tuple was found.
		fmt.Println("Pong recieved")
		// Send back a "Ping" tuple.
		tuplespace.Put(theirPtP, "Ping")
		fmt.Println("Ping Send")
		time.Sleep(500 * time.Millisecond)
	}
}

func createPtP(ownPort int, theirIP string, theirPort int) (topology.PointToPoint, topology.PointToPoint) {
	ownPtP := topology.CreatePointToPoint("Ping client", "localhost", ownPort)
	theirPtP := topology.CreatePointToPoint("Pong client", theirIP, theirPort)

	return ownPtP, theirPtP
}
