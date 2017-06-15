# goSpace

## Setting up Go
To work with your Go project in the most seamlessly way, it should be set up correctly. This is perfectly described in the following video where you're taken through each and every step of the set up.

[Writing, building, installing, and testing Go code](https://www.youtube.com/watch?v=XCsL89YtqCs)

## How to import goSpace
To import the goSpace into your project, clone the repository and place it in the `src` folder, which is located at the root of your GOPATH.

```terminal
GOPATH/src/
```
Now that the repository has been placed correctly, we can start using the actual framework. This is done by importing the `tuplespace` and `topology` packages into your project. This will like

```go
import (
      "goSpace/goSpace/tuplespace"
      "goSpace/goSpace/topology"
)
```

## Run test examples

### Install the examples
The `goSpace/examples` folder contains a few examples.

Let's start by looking at the `bookstore` example. To run this go to `*/goSpace/examples/bookstore` and run the following Go command to install the example, as was described in the video
```terminal
go install
```

### Run the examples
This will create an executable that is named as the folder where the `main.go` is located. For the `bookstore` example the `go install` command will create an executable called `bookstore`. The example can now be run by the following command regardless of the location on the system.

```terminal
bookstore
```
The remaining examples are installed and run in the same way.

### Running the bookstore example
If Go and the framework has been set up and imported correctly you should see the following being printed
```terminal
&{0xc420010d60 0xc42000cb10 [{[Of Mice and Men 200]}] :8080 []}
Checked price for book "Of Mice and Men". The price is 200.
Placed payment for book "Of Mice and Men", at the price of 200.
Recieved payment of 200 for the book "Of Mice and Men".
```
What you see is
* The first line is the tuple space that consists of one book being "Of Mice and Men" at the price of 200.
* The second line is right the after the customer performed a `Query` to check the price of the book he wants.
* The third line is after the customer has placed a payment for the book with a `Put`.
* The forth line is after the cashier has received the payment with a `Get`. He then checked if the price that the customer paid, matches the price of the book by a `Query`. As the price mathed, he removed the book from the store with a `Get` and handed it to the customer.

## How to perform tuple space operations.
The framework contains the following tuple space operation and the syntax of them looks as follows
```go
Put(ptp, x1, x2, ..., xn)
PutP(ptp, x1, x2, ..., xn)
Get(ptp, x1, x2, ..., xn)
GetP(ptp, x1, x2, ..., xn)
GetAll(ptp)
Query(ptp, x1, x2, ..., xn)
QueryP(ptp, x1, x2, ..., xn)
QueryAll(ptp)
```
Where `ptp` is a pointToPoint structure and `x1, x2, ..., xn` are values.

Pattern matching can be achieved by passing a pointer to the variable, this by can be done by adding a `&` infront of the variable. This is used in the `bookstore` example and looks a follows
```go
var i int
book := "Of Mice and Men"
tuplespace.Query(ptp, book, &i)
```
This is used to look up the price of the book "Of Mice and Men". After the `Query` operation the value of `i` is 200.
