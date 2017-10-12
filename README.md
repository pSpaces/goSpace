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
The goSpace API follows the Space API Specification. It contains the following tuple space operations:
```go
Put(ptp, x_1, x_2, ..., x_n)
PutP(ptp, x_1, x_2, ..., x_n)
Get(ptp, x_1, x_2, ..., x_n)
GetP(ptp, x_1, x_2, ..., x_n)
GetAll(ptp, x_1, x_2, ..., x_n)
Query(ptp, x_1, x_2, ..., x_n)
QueryP(ptp, x_1, x_2, ..., x_n)
QueryAll(ptp, x_1, x_2, ..., x_n)
```

Operations such as `Put`, `Get`, and so forth are blocking operations.

Operations postfixed by a `P` such as `PutP`, `GetP`, and so forth are non-blocking operations.

For all operations `ptp` is a PointToPoint structure and `x_1, x_2, ..., x_n` are terms.

For `Put` and `PutP` operations, the terms are values and for the remaining operations terms are either values or binding variables.

Pattern matching can be achieved by passing a binding variable, that is, passing a pointer to a variabe by adding a `&` infront of the variable.

Binding variables can only be passed to `Get*` and `Query*` operations.

## Space API Specification
The specification for the pSpace Space API can be found [here](https://github.com/pSpaces/Programming-with-Spaces/blob/master/guide.md).

## Examples
Examples and cases for goSpace can be found [here](https://github.com/luhac/gospace-examples).
