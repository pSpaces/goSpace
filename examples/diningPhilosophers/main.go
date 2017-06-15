package main

import (
	"fmt"
	"svn/bachelorProject/tupleSpaceFramework/topology"
	"svn/bachelorProject/tupleSpaceFramework/tuplespace"
	"time"
)

func main() {
	port := 8081
	ts := tuplespace.CreateTupleSpace(port)
	fmt.Println(ts)
	ptp := topology.CreatePointToPoint("me", "localhost", port)
	n := 2
	placeForks(ptp, n)
	go philosopher(ptp, 0, 0, n-1)
	for i := 1; i < n; i++ {
		go philosopher(ptp, i, i-1, i)
	}
	time.Sleep(7200 * time.Second)
}

func placeForks(ptp topology.PointToPoint, n int) {
	for i := 0; i < n; i++ {
		tuplespace.Put(ptp, "fork", i)
	}
}

func philosopher(ptp topology.PointToPoint, n int, fork1 int, fork2 int) {
	i := 0
	for {
		fmt.Printf("Philosopher %d is thinking\n", n)
		//why does this not work?
		//time.Sleep(rand.Intn(2000) * time.Millisecond)
		time.Sleep(1000 * time.Millisecond)
		fmt.Printf("Philosopher %d is hungry\n", n)
		tuplespace.Get(ptp, "fork", min(fork1, fork2))
		tuplespace.Get(ptp, "fork", max(fork1, fork2))
		//fmt.Printf("Philosopher %d got fork %d\n", n, min(fork1, fork2))
		//time.Sleep(500 * time.Millisecond)
		//fmt.Printf("Philosopher %d got fork %d\n", n, max(fork1, fork2))
		//time.Sleep(500 * time.Millisecond)
		fmt.Printf("Philosopher %d is eating for the %d. time\n", n, i)
		i++
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
