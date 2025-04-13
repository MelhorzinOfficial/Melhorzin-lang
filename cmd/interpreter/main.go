package main

import (
	"fmt"
	"io/ioutil"
	"melhorzin-lang/internal/interpreter"
	"melhorzin-lang/internal/lexer"
	"melhorzin-lang/internal/parser"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: emojilang <arquivo.mlz>")
		os.Exit(1)
	}

	filename := os.Args[1]
	code, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Erro ao ler arquivo: %v\n", err)
		os.Exit(1)
	}

	lex := lexer.NewLexer(string(code))
	tokens := lex.Lex()
	pars := parser.NewParser(tokens)
	nodes := pars.Parse()
	interp := interpreter.NewInterpreter()

	result := interp.Interpret(nodes)
	if result != nil {
		fmt.Printf("Resultado final: %v\n", result)
	}
}
