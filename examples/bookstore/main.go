package main

import (
	"fmt"
	"goSpace/goSpace"
	"reflect"
	"time"
)

func main() {
	store := goSpace.NewSpace("8080")
	addBooks(store)
	go cashier(store)
	customer(store)
	time.Sleep(2 * time.Second)
}

// addBooks adds books to the store.
func addBooks(store goSpace.PointToPoint) {
	book := "Of Mice and Men"
	goSpace.Put(store, book, 200)
}

func cashier(store goSpace.PointToPoint) {
	for {
		// Get the payment from the tuple space.
		var i int
		var book string
		goSpace.Get(store, "Payment", &book, &i)
		var price int
		// Find the price of the book
		goSpace.Query(store, book, &price)
		// Check if the priced paid is equal to what the book costs.
		if price == i {
			fmt.Printf("Received payment of %d for the book \"%s\".\n", i, book)
			// Remove the book from the store.
			goSpace.Get(store, book, i)
		}
	}
}

func customer(store goSpace.PointToPoint) {
	// Search for book and save price in i
	var i int
	book := "Of Mice and Men"
	goSpace.Query(store, book, &i)
	fmt.Printf("Checked price for book \"%s\". The price is %d.\n", book, i)
	// Place payment for the book
	goSpace.Put(store, "Payment", book, i)
	fmt.Printf("Placed payment for book \"%s\", at the price of %d.\n", book, i)
}

func printType(v interface{}) {
	fmt.Println(reflect.TypeOf(v))
}
