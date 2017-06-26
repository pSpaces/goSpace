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
	portTable, err := strconv.Atoi(os.Args[1])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	table := tuplespace.CreateTupleSpace(portTable)
	fmt.Println(table)
	ptTable := topology.CreatePointToPoint("table", "localhost", portTable)

	IPWaiter := os.Args[2]

	portWaiter, err := strconv.Atoi(os.Args[3])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	ptWaiter := topology.CreatePointToPoint("waiter", IPWaiter, portWaiter)

	n, err := strconv.Atoi(os.Args[4])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	t, err := strconv.Atoi(os.Args[5])
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}

	placeForks(ptTable, n)
	// Inform waiter of how many philosopher are at the table.
	for !tuplespace.Put(ptWaiter, "philosophers", n) {
		fmt.Println("Error in placing the number of philosophers at Waiter.")
	}

	go philosopher(ptTable, ptWaiter, 0, 0, n-1)
	for i := 1; i < n; i++ {
		go philosopher(ptTable, ptWaiter, i, i-1, i)
	}
	time.Sleep(time.Duration(t) * time.Second)
}

func placeForks(ptp topology.PointToPoint, n int) {
	for i := 0; i < n; i++ {
		tuplespace.Put(ptp, "fork", i)
	}
}

func philosopher(ptTable topology.PointToPoint, ptWaiter topology.PointToPoint, n int, fork1 int, fork2 int) {
	i := 0
	for {
		// philosopher thinks
		fmt.Printf("Philosopher %d is thinking\n", n)
		time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
		fmt.Printf("Philosopher %d is hungry\n", n)
		//sends request to eat
		if tuplespace.Put(ptWaiter, "request", n) {
			//look for permission
			for !tuplespace.Query(ptWaiter, "permission", n) {
				fmt.Println("Error in receiving permission.")
			}
			//grab forks
			tuplespace.Get(ptTable, "fork", fork1)
			tuplespace.Get(ptTable, "fork", fork2)
			//eat
			fmt.Printf("Philosopher %d is eating for the %d. time\n", n, i)
			i++
			//return forks
			tuplespace.Put(ptTable, "fork", fork1)
			tuplespace.Put(ptTable, "fork", fork2)
			//remove permission
			for !tuplespace.Get(ptWaiter, "permission", n) {
				fmt.Println("Error in getting the permission.")
			}
		}
	}
}
