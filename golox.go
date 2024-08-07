package main

import (
	"bufio"
	"fmt"
	"golox/interpreter"
	"golox/parser"
	// "golox/parser/astPrinter"
	"golox/scanner"
	"golox/tokens"
	"io/ioutil"
	"os"
)

func run(input string, intrpr *interpreter.Interpreter) (interface{}, error) {
	scan := scanner.NewScanner(input)
	var toks []tokens.Token
	for {
		t, _ := scan.Read()
		toks = append(toks, t)
		if t.Type == tokens.Eof {
			break
		}
	}

	// for _, t := range toks {
	// 	fmt.Printf("token: %s: %s\n", t.Type.String(), t.Lexeme)
	// }

	if len(scan.Errors()) != 0 {
		var errs string
		for _, err := range scan.Errors() {
			errs += err.Error() + "\n"
		}
		return nil, fmt.Errorf(errs)
	}


	stmts := parser.Parse(toks)
    var res interface{}
    var err error
    for _, stmt := range stmts {
        // astPrinter.PrintStmt(stmt)
        res, err = intrpr.Interpret(stmt)
        if err != nil {
            return nil, err
        }
    }

	return res, nil
}

func runPrompt() {
	s := bufio.NewScanner(os.Stdin)
    intrpr := interpreter.New()
	var line string = "\n"
	for {
		fmt.Print("> ")
		s.Scan()
		if len(s.Bytes()) == 0 {
			os.Exit(0)
		}
		line = s.Text() + ";"
		res, err := run(line, &intrpr)
		if err != nil {
			fmt.Printf("\u001b[31m%s\u001b[39m\n", err.Error())
		}
        switch v := res.(type) {
            case string:
                fmt.Printf("\u001b[2m\"%s\"\u001b[22m\n", v)
            case nil:
            default:
                fmt.Printf("\u001b[2m%v\u001b[22m\n", v)
        }
	}
}

func runFile(fileName string) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print(err)
	}

    intrpr := interpreter.New()
	run(string(b), &intrpr)
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
