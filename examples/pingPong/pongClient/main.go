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
	pongPort, err := strconv.Atoi(os.Args[1])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	// Get the IP address and port number of the client running the pong
	// application.
	pingIP := os.Args[2]
	pingPort, err := strconv.Atoi(os.Args[3])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}

	tsPtr := tuplespace.CreateTupleSpace(pongPort)
	fmt.Println(tsPtr)

	ownPtP, theirPtP := createPtP(pongPort, pingIP, pingPort)

	pong(ownPtP, theirPtP)
}

// Pong will initially act as server
func pong(ownPtP topology.PointToPoint, theirPtP topology.PointToPoint) {

	for {
		// Find a "Ping" tuple in own tuple space
		tuplespace.Get(ownPtP, "Ping")
		// A "Ping" tuple was found.
		fmt.Println("Ping Received.")
		// Send back a "Pong" tuple.
		tuplespace.Put(theirPtP, "Pong")
		fmt.Println("Pong Send")

		time.Sleep(500 * time.Millisecond)
	}
}

func createPtP(ownPort int, theirIP string, theirPort int) (topology.PointToPoint, topology.PointToPoint) {
	ownPtP := topology.CreatePointToPoint("Pong client", "localhost", ownPort)
	theirPtP := topology.CreatePointToPoint("Ping client", theirIP, theirPort)

	return ownPtP, theirPtP
}
