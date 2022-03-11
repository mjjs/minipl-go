package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <file_path>\n", os.Args[0])
		return
	}

	filePath := os.Args[1]

	fe := &frontEnd{}
	fe.Execute(filePath)
}
