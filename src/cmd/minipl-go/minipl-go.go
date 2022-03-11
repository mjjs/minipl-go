package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <file_path>\n", os.Args[0])
		return
	}

	filePath := os.Args[1]

	sourceCode, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	fe := &frontEnd{}
	fe.Execute(string(sourceCode))
}
