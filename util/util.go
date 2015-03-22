package util

import "fmt"
import "time"
import "strings"

func Timestamp() string {
  return fmt.Sprintf("%d", time.Now().Unix())
}

func CleanInput(input string) string {
  return strings.TrimSpace(input)
}
