package main

import (
  "encoding/json"
  "fmt"
  "io"
  "log"
  "net"
  "os"
)

import (
  "github.com/chrisbutcher/tcpgameserver/gamelogic"
  "github.com/chrisbutcher/tcpgameserver/types"
  "github.com/chrisbutcher/tcpgameserver/util"
  "github.com/satori/go.uuid"
)

const (
  Protocol   string = "tcp"
  ListenPort string = ":8080"
)

func listenForConnections(listener net.Listener, connections chan<- net.Conn) {
  for {
    conn, err := listener.Accept()

    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    connections <- conn
  }
}

func handleConnection(conn net.Conn, clientUuid string, requests chan<- types.ReceivedRequest, expiredConnections chan<- net.Conn) {
  decoder := json.NewDecoder(conn)
  newRequest := &types.Request{}

  for {
    if err := decoder.Decode(&newRequest); err == io.EOF {
      break
    } else if err != nil {
      log.Fatal(err)
      expiredConnections <- conn
    } else {
      log.Printf("Received request from client: %s", clientUuid)
      requests <- types.ReceivedRequest{util.Timestamp(), *newRequest, clientUuid}
    }
  }
}

func writeToClient(conn net.Conn, expiredConnections chan<- net.Conn, gameState types.GameState, clientUuid string) {
  encoder := json.NewEncoder(conn)
  err := encoder.Encode(gameState)

  log.Printf("Wrote gamestate to client %s => %+v", clientUuid, gameState)

  if err != nil {
    expiredConnections <- conn
  }
}

func updateGameState(receivedRequest types.ReceivedRequest, lastGameState types.GameState, gameStates chan<- types.GameState) {
  command := receivedRequest.Request.Command

  if gamelogic.AllowedCommands()[command] {
    lastGameState.Timestamp = util.Timestamp()
    lastGameState.State += receivedRequest.Request.Command + "," + receivedRequest.Request.Message + ":"
    gameStates <- lastGameState
  } else {
    // Request invalid. Do not update gamestate
  }
}

func main() {
  log.Printf("Server started listening, on %s", ListenPort)

  clientCount := 0
  clients := make(map[net.Conn]string)
  connections := make(chan net.Conn)
  expiredConnections := make(chan net.Conn)
  requests := make(chan types.ReceivedRequest)
  gameStates := make(chan types.GameState)
  lastGameState := types.GameState{}

  listener, err := net.Listen(Protocol, ListenPort)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  go listenForConnections(listener, connections)

  for {
    select {
    case conn := <-connections:
      clientCount += 1
      log.Printf("Client connected. Client count: %d", clientCount)

      clientUuid := fmt.Sprintf("%s", uuid.NewV4())
      clients[conn] = clientUuid
      go handleConnection(conn, clients[conn], requests, expiredConnections)

    case request := <-requests:
      go updateGameState(request, lastGameState, gameStates)

    case gameState := <-gameStates:
      for conn, _ := range clients {
        clientUuid := clients[conn]
        go writeToClient(conn, expiredConnections, gameState, clientUuid)
      }
      lastGameState = gameState

    case conn := <-expiredConnections:
      log.Printf("Client %s disconnected", clients[conn])
      delete(clients, conn)
    }
  }
}
