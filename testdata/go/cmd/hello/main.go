// Command hello is a fixture for the shared Go and image workflows.
package main

import (
	"fmt"

	"example.com/fixture/internal/greet"
)

func main() {
	fmt.Println(greet.Hello("Frazer"))
}
