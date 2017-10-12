# goSpace

## Synopsis
goSpace, a pSpace implementation in Go.

## Importing goSpace
To get goSpace, do:

```terminal
go get -u github.com/luhac/gospace
```
To import goSpace into your project, add:

```go
import (
      . "github.com/luhac/gospace"
)
```

## goSpace API

The goSpace API follows the Space API Specification. It contains the following operations:
```go
Put(x_1, x_2, ..., x_n)
PutP(x_1, x_2, ..., x_n)
Get(x_1, x_2, ..., x_n)
GetP(x_1, x_2, ..., x_n)
GetAll(x_1, x_2, ..., x_n)
Query(x_1, x_2, ..., x_n)
QueryP(x_1, x_2, ..., x_n)
QueryAll(x_1, x_2, ..., x_n)
```
A space can be created by using `NewSpace` for creating a local space, or `NewRemoteSpace` for connecting to a remote space. For example:
```go
spc := NewSpace("8080")
```

All operations act on a `Space` structures and `x_1, x_2, ..., x_n` are terms in a tuple.

Operations such as `Put`, `Get`, and so forth are blocking operations.

Operations postfixed by a `P` such as `PutP`, `GetP`, and so forth are non-blocking operations.

For `Put` and `PutP` operations, the terms are values and for the remaining operations terms are either values or binding variables.

Pattern matching can be achieved by passing a binding variable, that is, passing a pointer to a variabe by adding an `&` infront of the variable.

Binding variables can only be passed to `Get*` and `Query*` operations.

## Space API Specification
The specification for the pSpace Space API can be found [here](https://github.com/pSpaces/Programming-with-Spaces/blob/master/guide.md).

## Examples
Examples and cases for goSpace can be found [here](https://github.com/luhac/gospace-examples).
