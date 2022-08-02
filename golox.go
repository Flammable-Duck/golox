package main

import (
	"bufio"
	"fmt"
	"golox/scanner"
	"golox/tokens"
	"golox/parser"
	"golox/parser/astPrinter"
	"io/ioutil"
	"os"
)

func run(input string) {
	scan := scanner.NewScanner(input)
    var toks []tokens.Token
	for {
		t, _ := scan.Read()
        toks = append(toks, t)
		if t.Type == tokens.Eof {
			break
		}

		// switch t.Type {
		// case tokens.String:
		// 	fmt.Printf("token: \"%s\" \n", t.Lexeme)
		// default:
  //           fmt.Printf("%d:%d: %s \n", t.Position.Row, t.Position.Col, t.Lexeme)
		// }
	}
	if len(scan.Errors()) != 0 {
		for _, err := range scan.Errors() {
			fmt.Println(err.Error())

		}
        return
	}

    p := parser.New(toks)
    expr := p.Parse()
    if expr != nil {
        astPrinter.PrintAst(expr)
    } else {
        fmt.Println("Parse Error: expected expression")
        return
    }

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
    // run("\"apple")
}
