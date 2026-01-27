package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: billiemuk <build|serve|new> [args]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		fmt.Println("build: not yet implemented")
	case "serve":
		fmt.Println("serve: not yet implemented")
	case "new":
		fmt.Println("new: not yet implemented")
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
