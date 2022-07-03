package main

import (
	"bufio"
	"fmt"
	"golox/scanner"
	"golox/tokens"
	"io/ioutil"
	"os"
)

func parse() {}

func run(input string) {
	scan := scanner.NewScanner(input)
	for {
		t, _ := scan.Read()
		if t.Type == tokens.Eof {
			break
		}

		switch t.Type {
		case tokens.String:
			fmt.Printf("token: \"%s\" \n", t.Lexeme)
		default:
            fmt.Printf("%d:%d: %s \n", t.Position.Row, t.Position.Col, t.Lexeme)
		}
	}
	if len(scan.Errors()) != 0 {
		for _, err := range scan.Errors() {
			// fmt.Printf("Error: %s\n", err.Error())
			fmt.Println(err.Error())

		}
	}

}

func reportError(t tokens.Token, err error) {
	fmt.Printf("Error on line %d, column %d: %s",
		t.Position.Row, t.Position.Col, err.Error())
}

func runPrompt() {
	s := bufio.NewScanner(os.Stdin)
	var line string = "\n"
	for {
		fmt.Print("> ")
		s.Scan()
		line = s.Text()
		run(line)
	}
}

func runFile(fileName string) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print(err)
	}

	run(string(b))
}

func main() {
	if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else if len(os.Args) == 1 {
		runPrompt()
	} else {
		fmt.Println("run without arguments to enter a repl, or with a filename to run a file")
	}
}
