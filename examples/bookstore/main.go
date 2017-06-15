package main

import (
	"fmt"
	"svn/bachelorProject/tupleSpaceFramework/topology"
	"svn/bachelorProject/tupleSpaceFramework/tuplespace"
	"time"
)

func main() {
	ts := tuplespace.CreateTupleSpace(8080)
	fmt.Println(ts)
	ptp := topology.CreatePointToPoint("Bookstore", "localhost", 8080)
	addBooks(ptp)
	go cashier(ptp)
	costumer(ptp)
	fmt.Println(1234)
	time.Sleep(2 * time.Second)
	fmt.Println(ts)

}

func addBooks(ptp topology.PointToPoint) {
	//s := "Of Mice and Men"
	tuplespace.Put(ptp, "Of Mice and Men", 200)
}

func cashier(ptp topology.PointToPoint) {
	for {
		var i int
		var book string
		tuplespace.Get(ptp, "Payment", &book, &i)
		var price int
		tuplespace.Query(ptp, book, &price)
		fmt.Println(book, price, i)
		if price == i {
			fmt.Printf("Recieved payment of %d for %s\n", i, book)
			tuplespace.Get(ptp, book, i)
		}
	}
}

func costumer(ptp topology.PointToPoint) {
	var i int
	book := "Of Mice and Men"
	//search for book and save price in i
	tuplespace.Query(ptp, book, &i)
	//place payment
	tuplespace.Put(ptp, "Payment", book, i)
	fmt.Printf("Puchased %s, for %d\n", book, i)
}
