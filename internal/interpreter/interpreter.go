package interpreter

import (
	"fmt"
	"melhorzin-lang/internal/parser"
)

// Interpreter executa a AST.
type Interpreter struct {
	variables map[string]interface{}
	types     map[string]parser.Type // Mapa para guardar os tipos das variáveis
	result    interface{}
}

// NewInterpreter cria um novo interpretador.
func NewInterpreter() *Interpreter {
	return &Interpreter{
		variables: make(map[string]interface{}),
		types:     make(map[string]parser.Type),
	}
}

// Interpret executa os nós da AST.
func (i *Interpreter) Interpret(nodes []parser.Node) interface{} {
	i.result = nil

	for _, node := range nodes {
		result := node.Evaluate(i.variables)
		i.result = result

		// Se for um nó de atribuição, armazenar o tipo da variável
		if assignNode, ok := node.(*parser.AssignNode); ok {
			i.types[assignNode.Name] = assignNode.GetType()
		}

		// Remover prints duplicados - o PrintNode já imprime diretamente
		// Apenas mostrar outros tipos de resultados
		if result != nil {
			if _, ok := node.(*parser.PrintNode); !ok {
				switch v := result.(type) {
				case int:
					fmt.Println(v)
				case bool:
					fmt.Println(v)
				case string:
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

// GetVariableType retorna o tipo de uma variável
func (i *Interpreter) GetVariableType(name string) parser.Type {
	if t, exists := i.types[name]; exists {
		return t
	}
	return parser.TypeAny
}
