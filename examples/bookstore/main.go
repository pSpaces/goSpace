package main

import (
	"fmt"
	"goSpace/goSpace/topology"
	"goSpace/goSpace/tuplespace"
	"reflect"
	"time"
)

func main() {
	ts := tuplespace.CreateTupleSpace(8080)
	ptp := topology.CreatePointToPoint("Bookstore", "localhost", 8080)
	addBooks(ptp)
	fmt.Println(ts)
	go cashier(ptp)
	costumer(ptp)
	time.Sleep(2 * time.Second)
}

// addBooks adds books to the store.
func addBooks(ptp topology.PointToPoint) {
	book := "Of Mice and Men"
	tuplespace.Put(ptp, book, 200)
}

func cashier(ptp topology.PointToPoint) {
	for {
		// Get the payment from the tuple space.
		var i int
		var book string
		tuplespace.Get(ptp, "Payment", &book, &i)
		var price int
		// Find the price of the book
		tuplespace.Query(ptp, book, &price)
		// Check if the priced paid is equal to what the book costs.
		if price == i {
			fmt.Printf("Recieved payment of %d for the book \"%s\".\n", i, book)
			// Remove the book from the store.
			tuplespace.Get(ptp, book, i)
		}
	}
}

func costumer(ptp topology.PointToPoint) {
	// Search for book and save price in i
	var i int
	book := "Of Mice and Men"
	tuplespace.Query(ptp, book, &i)
	fmt.Printf("Checked price for book \"%s\". The price is %d.\n", book, i)
	// Place payment for the book
	tuplespace.Put(ptp, "Payment", book, i)
	fmt.Printf("Placed payment for book \"%s\", at the price of %d.\n", book, i)
}

func printType(v interface{}) {
	fmt.Println(reflect.TypeOf(v))
}
