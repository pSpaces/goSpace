The diningPhilosophersRemote contains a waiter and a philosophers application.

The waiter application will initially act as a server and should be started first.
It takes one input when running the example
```terminal
waiter waiterPort
```

The philosophers application will initially act as a client and should be started af the waiter.
It takes five inputs when running the example
```terminal
philosophers philosophersPort waiterIP waiterPort numberOfPhilosophers timeToRunApplication
```
