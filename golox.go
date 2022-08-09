package main

import (
	"bufio"
	"fmt"
	"golox/interpreter"
	"golox/parser"
	"golox/parser/astPrinter"
	"golox/scanner"
	"golox/tokens"
	"io/ioutil"
	"os"
)

func run(input string) error {
	scan := scanner.NewScanner(input)
	var toks []tokens.Token
	for {
		t, _ := scan.Read()
		toks = append(toks, t)
		if t.Type == tokens.Eof {
			break
		}
	}

	for _, t := range toks {
		fmt.Printf("token: %s:%s\n", t.Type.String(), t.Lexeme)
	}

	if len(scan.Errors()) != 0 {
		var errs string
		for _, err := range scan.Errors() {
			errs += err.Error() + "\n"
		}
		return fmt.Errorf(errs)
	}

	expr, err := parser.Parse(toks)
	if err != nil {
		return err
	}
	astPrinter.PrintAst(expr)

	return interpreter.Interpret(expr)
}

func runPrompt() {
	s := bufio.NewScanner(os.Stdin)
	var line string = "\n"
	for {
		fmt.Print("> ")
		s.Scan()
		if len(s.Bytes()) == 0 {
			os.Exit(0)
		}
		line = s.Text()
		err := run(line)
		if err != nil {
			fmt.Println(err.Error())
		}
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
	// if err := run("1+1"); err != nil {
	// 	fmt.Println(err)
	// }
}
