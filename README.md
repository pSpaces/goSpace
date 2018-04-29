# goSpace

## Synopsis
goSpace, a pSpace implementation in Go.

## Importing goSpace
To get goSpace, do:

```terminal
go get -u github.com/pspaces/gospace
```
To import goSpace into your project, add:

```go
import (
      "github.com/pspaces/gospace"
)
```

## Documentation
Documentation for goSpace can be accessed through the terminal with `go doc` or viewed online with `godoc`. 

## Library

### Usage
A space can be created by using `NewSpace` for creating a local space or `NewRemoteSpace` for connecting to a remote space.

To create a space on the local host, one can do:

```go
spc := gospace.NewSpace("space")
```

To connect to a remote space with name `space`, one can do:

```go
spc := gospace.NewRemoteSpace("tcp://example.com/space")
```

In order to use goSpace efficiently, there are certain rules one needs to be aware of:

   1. An operation acts on a `Space` structure.
   2. Tuples contain values. Templates contain values and types.
   3. Ordinary operations such as `Put`, `Get` and so forth are blocking operations.
   4. Ordinary operations postfixed by a `P` such as `PutP`, `GetP` and so forth are non-blocking operations.
   5. Aggregation operations postfixed by a `Agg` such as `PutAgg`, `GetAgg` and so forth are non-blocking operations.
   6. Aggregation operations require inclusion of a user defined function `f` which operate on a `Tuple` structure and has the signature `func(Tuple, Tuple) Tuple`.
   7. Aggregation operations require that the user defined function `f` exists on all peers that will use that function.
   8. For `Put` and `PutP` operations the parameters must be values. For the remaining operations the parameters must be values or binding variables for pattern matching.
   9. Return signature of an operations is always `func(...) (Tuple, error)`.

Due to aggregation operators, binding variables are not written to when an operation completes.

### Operators
goSpace strives to follow the pSpace specification. It contains the following operations:

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
Do note that `GetAll` and `QueryAll` are non-blocking operators.

goSpace has experimental operators for aggregating tuples in a space. It contains the following operations:

```go
PutAgg(f, x_1, x_2, ..., x_n)
GetAgg(f, x_1, x_2, ..., x_n)
QueryAgg(f, x_1, x_2, ..., x_n)
```

## Specification
The specification for the pSpace can be found [here](https://github.com/pspaces/Programming-with-Spaces/blob/master/guide.md).

## Limitations
There are currently some limitations to the implementation:
 - Strict 4 GiB size limit on the tuple space.
 - Only TCP over IPv4 is supported.
 - Gates and space repositories are not supported.
 - Multiplexing of multiple spaces over a single gate is not supported.
 - Performance is bouded by the amount of connections established due to networking code.

These limitations are currently being resolved.

## Examples
Examples and cases for goSpace can be found [here](https://github.com/pspaces/gospace-examples).

## Contributing
Want to contribute? See [CONTRIBUTING.md](https://github.com/pspaces/gospace/blob/master/CONTRIBUTING.md).

## Authors
See [AUTHORS.md](https://github.com/pspaces/gospace/blob/master/AUTHORS.md).

## License
MIT License. See [LICENSE.md](https://github.com/pspaces/gospace/blob/master/LICENSE.md).
