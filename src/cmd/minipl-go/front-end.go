package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/mjjs/minipl-go/pkg/interpreter"
	"github.com/mjjs/minipl-go/pkg/lexer"
	"github.com/mjjs/minipl-go/pkg/parser"
	"github.com/mjjs/minipl-go/pkg/symboltable"
	"github.com/mjjs/minipl-go/pkg/typechecker"
)

type frontEnd struct {
	out io.Writer
	in  io.Reader
}

func (fe *frontEnd) Execute(filepath string) {
	if fe.out == nil {
		fe.out = os.Stdout
	}

	if fe.in == nil {
		fe.in = os.Stdin
	}

	sourceBytes := readFile(filepath)

	lexer := lexer.New(string(sourceBytes))
	parser := parser.New(lexer)

	astRoot, errors := parser.Parse()
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Fprintln(fe.out, err)
		}
		return
	}

	stc := &symboltable.SymbolTableCreator{}
	symbols, errors := stc.Create(astRoot)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Fprintln(fe.out, err)
		}
		return
	}

	tc := typechecker.New(symbols)
	errors = tc.CheckTypes(astRoot)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Fprintln(fe.out, err)
		}
		return
	}

	i := interpreter.New(fe.out, fe.in)
	i.Run(astRoot)
}

func readFile(path string) []byte {
	sourceCode, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return sourceCode
}
