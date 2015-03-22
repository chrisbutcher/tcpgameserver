package main

import (
  "bufio"
  "encoding/json"
  "fmt"
  "io"
  "log"
  "net"
  "os"
)

import (
  "github.com/chrisbutcher/tcpgameserver/types"
  "github.com/chrisbutcher/tcpgameserver/util"
)

const (
  Protocol      string = "tcp"
  ServerAddress string = "localhost:8080"
)

func handleConnection(conn net.Conn) {
  decoder := json.NewDecoder(conn)
  newGameState := &types.GameState{}

  for {
    if err := decoder.Decode(&newGameState); err == io.EOF {
      break
    } else if err != nil {
      log.Fatal(err)
    } else {
      log.Printf("New gamestate received: %+v\nCommand> ", newGameState)
    }
  }
}

func main() {
  conn, err := net.Dial(Protocol, ServerAddress)

  if err != nil {
    log.Fatal("Connection error", err)
  }

  log.Println("Connected")

  go handleConnection(conn)

  encoder := json.NewEncoder(conn)
  keyboard := bufio.NewReader(os.Stdin)

  fmt.Print("Command> ")

  for {
    command, err := keyboard.ReadString('\n')

    if err != nil {
      log.Fatal(err)
      continue
    }

    fmt.Print("Message> ")
    message, err := keyboard.ReadString('\n')

    if err != nil {
      log.Fatal(err)
      continue
    }

    command = util.CleanInput(command)
    message = util.CleanInput(message)

    err = encoder.Encode(&types.Request{Command: command, Message: message, Timestamp: util.Timestamp()})

    if err != nil {
      log.Println(err)
      os.Exit(1)
    }
  }

  conn.Close()
}
