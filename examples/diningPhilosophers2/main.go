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
	waiterts := tuplespace.CreateTupleSpace(portWaiter)
	fmt.Println(waiterts)
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

func placeForks(ptp topology.PointToPoint, n int) {
	for i := 0; i < n; i++ {
		tuplespace.Put(ptp, "fork", i)
	}
}

func philosopher(pttable topology.PointToPoint, ptwaiter topology.PointToPoint, n int, fork1 int, fork2 int) {
	i := 0
	for {
		fmt.Printf("Philosopher %d is thinking\n", n)
		//why does this not work?
		time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
		//time.Sleep(1000 * time.Millisecond)
		fmt.Printf("Philosopher %d is hungry\n", n)
		tuplespace.Put(ptwaiter, "request", n)
		tuplespace.Query(ptwaiter, "permission", n)
		tuplespace.Get(pttable, "fork", fork1)
		tuplespace.Get(pttable, "fork", fork2)
		fmt.Printf("Philosopher %d is eating for the %d. time\n", n, i)
		i++
		tuplespace.Put(pttable, "fork", fork1)
		tuplespace.Put(pttable, "fork", fork2)
		tuplespace.Get(ptwaiter, "permission", n)
	}
}

func waiter(ptwaiter topology.PointToPoint, n int) {
	for {
		var philosopher int
		tuplespace.Get(ptwaiter, "request", &philosopher)
		var neighbor1 int
		var neighbor2 int
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
		for tuplespace.QueryP(ptwaiter, "permission", neighbor1) || tuplespace.QueryP(ptwaiter, "permission", neighbor2) {

		}
		fmt.Printf("Waiter gave permission to philosopher %d\n", philosopher)
		tuplespace.Put(ptwaiter, "permission", philosopher)
	}
}
