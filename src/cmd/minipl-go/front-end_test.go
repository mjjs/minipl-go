package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

var testCases = []struct {
	name           string
	sourceCode     string
	expectedOutput *bytes.Buffer
	userInput      io.Reader
}{
	{
		name: "Factorial program",
		sourceCode: `
		print "Give a number";
		var n : int;
		read n;
		print "\n";
		var v : int := 1;
		var i : int;
		for i in 1..n do
			v := v * i;
		end for;
		print "The result is: ";
		print v;
		`,
		userInput:      bytes.NewBufferString("5\n"),
		expectedOutput: bytes.NewBufferString("Give a number\nThe result is: 24"),
	},
	{
		name: "Boolean logic",
		sourceCode: `
			var true : bool := (1=1);
			var false : bool := !true;
			print true;
			print "\n";
			print false;
			print "\n";
			print (true & true);
			print "\n";
			print (false & false);
			print "\n";
			print (true & false);
		`,
		expectedOutput: bytes.NewBufferString("true\nfalse\ntrue\nfalse\nfalse"),
	},
}

func TestEndToEndInterpreter(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := writeTempFile(t, tc.name, tc.sourceCode)
			defer removeTempFile(t, f)

			w := &bytes.Buffer{}

			fe := &frontEnd{out: w, in: tc.userInput}

			fe.Execute(f.Name())

			if w.String() != tc.expectedOutput.String() {
				t.Errorf("Expected: %s\ngot: %s", tc.expectedOutput.String(), w.String())
			}
		})
	}
}

func writeTempFile(t *testing.T, name string, source string) *os.File {
	t.Helper()

	f, err := os.CreateTemp("", fmt.Sprintf("%s.*.minipl", name))
	if err != nil {
		panic(err)
	}

	if _, err := f.Write([]byte(source)); err != nil {
		panic(err)
	}

	return f
}

func removeTempFile(t *testing.T, f *os.File) {
	t.Helper()

	f.Close()
	os.Remove(f.Name())
}
