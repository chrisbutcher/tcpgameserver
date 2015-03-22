package types

type ReceivedRequest struct {
  Timestamp  string
  Request    Request
  ClientUuid string
}

type Request struct {
  Timestamp string
  Command   string
  Message   string
}

type GameState struct {
  Timestamp string
  State     string
}
