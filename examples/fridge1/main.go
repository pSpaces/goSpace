package main

import (
	"fmt"
	"goSpace/goSpace"
)

func main() {
	fridge := goSpace.NewSpace("8080")

	// add some stuff to the grocery list
	goSpace.Put(fridge, "milk", 2)
	goSpace.Put(fridge, "butter", 3)

	// retrieve one item to
	var item string
	var quantity int
	goSpace.Get(fridge, &item, &quantity)

	// print whatever we retrieved
	fmt.Println("We found this tuple: (", item, ",", quantity, ")")
}
