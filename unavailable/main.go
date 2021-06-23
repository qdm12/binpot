package main

import (
	"fmt"
	"os"
)

var (
	name     = "this program"
	platform = "this platform"
)

func main() {
	fmt.Println(name + " is unavailable on " + platform)
	os.Exit(1)
}
