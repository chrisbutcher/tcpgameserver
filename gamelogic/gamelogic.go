package gamelogic

func AllowedCommands() map[string]bool {
  return map[string]bool{
    "move":   true,
    "attack": true,
  }
}
