package main

import (
	"fmt"
)

func main() {
	fmt.Println("colors={{ .Scaffold.colors | join `, ` }} description={{ .Scaffold.description }}")
}
