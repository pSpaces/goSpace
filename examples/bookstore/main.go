package main

import (
	"fmt"
	"goSpace/goSpace/topology"
	"goSpace/goSpace/tuplespace"
)

func main() {
	ts := tuplespace.CreateTupleSpace(8080)
	ptp := topology.CreatePointToPoint("Bookstore", "localhost", 8080)
	tuplespace.Put(ptp, 2, 2)
	tuplespace.Put(ptp, 2, 2)
	tuplespace.Put(ptp, 2, 3)
	tuplespace.Put(ptp, 2, 3)
	tuplespace.Put(ptp, 2, false)
	fmt.Println(ts)
	var i int
	fmt.Println(tuplespace.QueryAll(ptp, 2, 2))
	fmt.Println(tuplespace.GetAll(ptp, 2, &i))
	fmt.Println(tuplespace.QueryAll(ptp, 2, 2))
	fmt.Println(ts)
	/*
		addBooks(ptp)
		fmt.Println(ts)
		go cashier(ptp)
		costumer(ptp)
		time.Sleep(2 * time.Second)
	*/
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
		if price == i {
			fmt.Printf("Recieved payment of %d for the book \"%s\".\n", i, book)
			tuplespace.Get(ptp, book, i)
		}
	}
}

func costumer(ptp topology.PointToPoint) {
	var i int
	book := "Of Mice and Men"
	//search for book and save price in i
	tuplespace.Query(ptp, book, &i)
	fmt.Printf("Checked price for book \"%s\". The price is %d.\n", book, i)
	//place payment
	tuplespace.Put(ptp, "Payment", book, i)
	fmt.Printf("Placed payment for book \"%s\", at the price of %d.\n", book, i)
}
