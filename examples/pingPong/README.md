The pongClient initially acts as a server and should therefore be started before the pingClient, which initially will act as a client.

The pongClient takes three arguments when running the example

```terminal
pongClient pongClientPort pingClientIP pingClientPort
```

The pingClient takes three arguments when running the example

```terminal
pingClient pingClientPort pongClientIP pongClientPort
```
