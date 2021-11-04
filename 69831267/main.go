package main
import(
	"fmt"
	"io"
)
func Print(out io.Writer) {
		fmt.Fprint(out, "Hello, world")
	}
