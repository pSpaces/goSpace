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
	port := 8081
	ts := tuplespace.CreateTupleSpace(port)
	fmt.Println(ts)
	ptp := topology.CreatePointToPoint("me", "localhost", port)

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

	placeForks(ptp, n)
	go philosopher(ptp, 0, 0, n-1)
	for i := 1; i < n; i++ {
		go philosopher(ptp, i, i-1, i)
	}
	time.Sleep(time.Duration(t) * time.Second)
}

// placeForks places forks on the table.
func placeForks(ptp topology.PointToPoint, n int) {
	for i := 0; i < n; i++ {
		tuplespace.Put(ptp, "fork", i)
	}
}

func philosopher(ptp topology.PointToPoint, n int, fork1 int, fork2 int) {
	i := 0
	for {
		fmt.Printf("Philosopher %d is thinking\n", n)
		time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
		fmt.Printf("Philosopher %d is hungry\n", n)
		// Get the fork with the lowest number.
		tuplespace.Get(ptp, "fork", min(fork1, fork2))
		// Get the fork with the highest number.
		tuplespace.Get(ptp, "fork", max(fork1, fork2))
		fmt.Printf("Philosopher %d is eating for the %d. time\n", n, i)
		i++
		// Put forks back when done eating.
		tuplespace.Put(ptp, "fork", fork1)
		tuplespace.Put(ptp, "fork", fork2)

	}
}

//min function for intergers is not in standard library because reasons
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
