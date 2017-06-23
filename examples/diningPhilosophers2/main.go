package main

import (
	"fmt"
	"goSpace/goSpace/topology"
	"goSpace/goSpace/tuplespace"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	portTable := 8081
	portWaiter := 8082
	table := tuplespace.CreateTupleSpace(portTable)
	fmt.Println(table)
	pttable := topology.CreatePointToPoint("me", "localhost", portTable)
	waiterTS := tuplespace.CreateTupleSpace(portWaiter)
	fmt.Println(waiterTS)
	ptwaiter := topology.CreatePointToPoint("me", "localhost", portWaiter)

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	t, err := strconv.Atoi(os.Args[2])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}

	placeForks(pttable, n)
	go waiter(ptwaiter, n)
	go philosopher(pttable, ptwaiter, 0, 0, n-1)
	for i := 1; i < n; i++ {
		go philosopher(pttable, ptwaiter, i, i-1, i)
	}
	time.Sleep(time.Duration(t) * time.Second)
}

// placeForks places forks on the table.
func placeForks(ptp topology.PointToPoint, n int) {
	for i := 0; i < n; i++ {
		tuplespace.Put(ptp, "fork", i)
	}
}

func philosopher(pttable topology.PointToPoint, ptwaiter topology.PointToPoint, n int, fork1 int, fork2 int) {
	i := 0
	for {
		// Philosopher thinks
		fmt.Printf("Philosopher %d is thinking\n", n)
		time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
		fmt.Printf("Philosopher %d is hungry\n", n)
		// Sends request to eat
		tuplespace.Put(ptwaiter, "request", n)
		// LookS for permission
		tuplespace.Query(ptwaiter, "permission", n)
		// Grab forks
		tuplespace.Get(pttable, "fork", fork1)
		tuplespace.Get(pttable, "fork", fork2)
		// Eats
		fmt.Printf("Philosopher %d is eating for the %d. time\n", n, i)
		i++
		// Return forks
		tuplespace.Put(pttable, "fork", fork1)
		tuplespace.Put(pttable, "fork", fork2)
		// Remove permission
		tuplespace.Get(ptwaiter, "permission", n)
	}
}

func waiter(ptwaiter topology.PointToPoint, n int) {
	for {
		var philosopher int
		// Get requests
		if tuplespace.Get(ptwaiter, "request", &philosopher) {

			var neighbor1 int
			var neighbor2 int
			// Find neighbours
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
			found1, conn1 := tuplespace.QueryP(ptwaiter, "permission", neighbor1)
			found2, conn2 := tuplespace.QueryP(ptwaiter, "permission", neighbor2)
			// See if neighbours have permission to eat
			for found1 || found2 || !conn1 || !conn2 {
				found1, conn1 = tuplespace.QueryP(ptwaiter, "permission", neighbor1)
				found2, conn2 = tuplespace.QueryP(ptwaiter, "permission", neighbor2)
			}
			fmt.Printf("Waiter gave permission to philosopher %d\n", philosopher)
			// Give permission to philosopher
			tuplespace.Put(ptwaiter, "permission", philosopher)
		}
	}
}
