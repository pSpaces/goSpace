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
	portWaiter, err := strconv.Atoi(os.Args[1])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	tsWaiter := tuplespace.CreateTupleSpace(portWaiter)
	fmt.Println(tsWaiter)
	ptWaiter := topology.CreatePointToPoint("waiter", "localhost", portWaiter)

	n, err := strconv.Atoi(os.Args[2])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	t, err := strconv.Atoi(os.Args[3])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}

	go waiter(ptWaiter, n)
	time.Sleep(time.Duration(t) * time.Second)
}

func waiter(ptWaiter topology.PointToPoint, n int) {
	for {
		var philosopher int
		//get requests
		if tuplespace.Get(ptWaiter, "request", &philosopher) {

			var neighbor1 int
			var neighbor2 int
			//Find neighbors
			if philosopher == 0 {
				neighbor1 = n - 1
				neighbor2 = 1
			} else if philosopher == n-1 {
				neighbor1 = n - 2
				neighbor2 = 0
			} else {
				neighbor1 = philosopher - 1
				neighbor2 = philosopher + 1
			}
			found1, conn1 := tuplespace.QueryP(ptWaiter, "permission", neighbor1)
			found2, conn2 := tuplespace.QueryP(ptWaiter, "permission", neighbor2)
			//see if neighbors have permission to eat
			for found1 || found2 || !conn1 || !conn2 {
				found1, conn1 = tuplespace.QueryP(ptWaiter, "permission", neighbor1)
				found2, conn2 = tuplespace.QueryP(ptWaiter, "permission", neighbor2)
			}
			fmt.Printf("Waiter gave permission to philosopher %d\n", philosopher)
			//Give permission to philosopher
			tuplespace.Put(ptWaiter, "permission", philosopher)
		}
	}
}
