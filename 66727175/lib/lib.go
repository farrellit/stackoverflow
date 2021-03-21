package lib

import(
  "fmt"
  "io"
)

func Hello(w io.Writer) {
  fmt.Fprintf(w, "Hello, world!\n")
}
