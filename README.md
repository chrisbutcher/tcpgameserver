# tcpgameserver

Run the server with
```go run server.go```

Run the client with
```go run client.go```

You can then issue game client commands ('move' or 'attack') with messages associated with the commands, and all clients will receive an updated gamestate (one that is approved by the server). Everything is asynchronous, and uses JSON-formatted []byte slices as message payloads.
