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
The `goSpace/examples` folder contains a few examples.

Let's start by looking at the `bookstore` example. To run this go to `*/goSpace/examples/bookstore` and run the following Go command to install the example, as was described in the video

```terminal
go install
```

This will create an executable that is named as the folder where the `main.go` is located. For the `bookstore` example the `go install` command will create an executable called `bookstore`. The example can now be run by the following command regardless of the location on the system.

```terminal
bookstore
```



## How to perform tuple space operations.
