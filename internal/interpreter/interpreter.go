package interpreter

import (
	"fmt"
	"melhorzin-lang/internal/parser"
)

// Interpreter executa a AST.
type Interpreter struct {
	variables map[string]interface{}
	result    interface{}
}

// NewInterpreter cria um novo interpretador.
func NewInterpreter() *Interpreter {
	return &Interpreter{variables: make(map[string]interface{})}
}

// Interpret executa os nós da AST.
func (i *Interpreter) Interpret(nodes []parser.Node) interface{} {
	i.result = nil

	for _, node := range nodes {
		result := node.Evaluate(i.variables)
		i.result = result

		// Remover prints duplicados - o PrintNode já imprime diretamente
		// Apenas mostrar outros tipos de resultados
		if result != nil {
			if _, ok := node.(*parser.PrintNode); !ok {
				switch v := result.(type) {
				case int:
					fmt.Println(v)
				case bool:
					fmt.Println(v)
				case *parser.FunctionNode:
					// Não exibe nada quando define uma função
				default:
					// Se o resultado for um valor de retorno de chamada de função, exibe
					if _, ok := node.(*parser.FunctionCallNode); ok {
						// Não imprimir novamente, o valor já foi impresso na função
					}
				}
			}
		}
	}

	return i.result
}

// GetResult retorna o último valor calculado
func (i *Interpreter) GetResult() interface{} {
	return i.result
}
