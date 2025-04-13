package parser

import (
	"fmt"
	"melhorzin-lang/internal/lexer"
	"regexp"
	"strconv"
	"strings"
)

// Node representa um nó da AST.
type Node interface {
	Evaluate(vars map[string]interface{}) interface{}
}

// PrintNode para instruções de impressão.
type PrintNode struct {
	Text         string
	Interpolated bool
	Variables    map[string]struct{} // Variáveis que serão interpoladas
}

func (n *PrintNode) Evaluate(vars map[string]interface{}) interface{} {
	if !n.Interpolated {
		// Garantir que o texto seja impresso
		fmt.Println(n.Text)
		return n.Text
	}

	// Processa a interpolação
	result := n.Text
	for varName := range n.Variables {
		if val, exists := vars[varName]; exists {
			placeholder := fmt.Sprintf("💱{%s}", varName)
			result = strings.Replace(result, placeholder, fmt.Sprintf("%v", val), -1)
		}
	}

	// Garantir que o resultado seja impresso
	fmt.Println(result)
	return result
}

// AssignNode para atribuições.
type AssignNode struct {
	Name  string
	Value interface{}
}

func (n *AssignNode) Evaluate(vars map[string]interface{}) interface{} {
	// Se o valor for um nó, avaliá-lo primeiro
	if node, ok := n.Value.(Node); ok {
		vars[n.Name] = node.Evaluate(vars)
	} else {
		vars[n.Name] = n.Value
	}
	return nil
}

// EqualNode para comparações.
type EqualNode struct {
	Name  string
	Value interface{}
}

func (n *EqualNode) Evaluate(vars map[string]interface{}) interface{} {
	if val, exists := vars[n.Name]; exists {
		return val == n.Value
	}
	return false
}

// BinaryOpNode para operações binárias
type BinaryOpNode struct {
	Left  Node
	Op    lexer.TokenType
	Right Node
}

func (n *BinaryOpNode) Evaluate(vars map[string]interface{}) interface{} {
	leftVal := n.Left.Evaluate(vars)
	rightVal := n.Right.Evaluate(vars)

	switch n.Op {
	case lexer.TokenPlus:
		// + agora é só para soma numérica
		if leftInt, ok := leftVal.(int); ok {
			if rightInt, ok := rightVal.(int); ok {
				return leftInt + rightInt
			}
		}
		return nil
	case lexer.TokenConcat:
		// . é para concatenação de strings
		return fmt.Sprintf("%v%v", leftVal, rightVal)
	case lexer.TokenNumPlus:
		// ➕ é para soma numérica (manter por compatibilidade)
		if leftInt, ok := leftVal.(int); ok {
			if rightInt, ok := rightVal.(int); ok {
				return leftInt + rightInt
			}
		}
		return nil
	case lexer.TokenMult:
		// ✖️ é para multiplicação numérica
		if leftInt, ok := leftVal.(int); ok {
			if rightInt, ok := rightVal.(int); ok {
				return leftInt * rightInt
			}
		}
	}
	return nil
}

// VariableNode para acessar variáveis
type VariableNode struct {
	Name string
}

func (n *VariableNode) Evaluate(vars map[string]interface{}) interface{} {
	if val, exists := vars[n.Name]; exists {
		return val
	}
	return nil
}

// FunctionNode para definição de funções
type FunctionNode struct {
	Name       string
	Parameters []string
	Body       []Node
}

func (n *FunctionNode) Evaluate(vars map[string]interface{}) interface{} {
	// Armazena a função no mapa de variáveis
	vars[n.Name] = n
	return nil
}

// ReturnNode para retorno de valores
type ReturnNode struct {
	Value Node
}

func (n *ReturnNode) Evaluate(vars map[string]interface{}) interface{} {
	return n.Value.Evaluate(vars)
}

// FunctionCallNode para chamadas de função
type FunctionCallNode struct {
	Name      string
	Arguments []Node
}

func (n *FunctionCallNode) Evaluate(vars map[string]interface{}) interface{} {
	if fnValue, exists := vars[n.Name]; exists {
		if fn, ok := fnValue.(*FunctionNode); ok {
			// Criar ambiente local para a função
			localVars := make(map[string]interface{})
			for k, v := range vars {
				localVars[k] = v
			}

			// Avaliar argumentos e associá-los aos parâmetros
			for i, argNode := range n.Arguments {
				if i < len(fn.Parameters) {
					localVars[fn.Parameters[i]] = argNode.Evaluate(vars)
				}
			}

			// Executar o corpo da função
			var result interface{}
			for _, node := range fn.Body {
				if returnNode, ok := node.(*ReturnNode); ok {
					// Se encontrar um return, avalia e retorna
					return returnNode.Value.Evaluate(localVars)
				}
				// Avaliar cada nó no corpo da função
				nodeResult := node.Evaluate(localVars)

				// Se o nó for um PrintNode, ele já imprimiu seu conteúdo
				// Só precisamos atualizar o resultado para o próximo nó
				result = nodeResult
			}
			return result
		}
	}
	return nil
}

// MainNode para a função main.
type MainNode struct {
	Body []Node
}

func (n *MainNode) Evaluate(vars map[string]interface{}) interface{} {
	for _, node := range n.Body {
		node.Evaluate(vars)
	}
	return nil
}

// TryCatchNode para try-catch.
type TryCatchNode struct {
	TryBody   []Node
	CatchBody []Node
}

func (n *TryCatchNode) Evaluate(vars map[string]interface{}) interface{} {
	for _, node := range n.TryBody {
		if _, ok := node.(*TryCatchNode); ok {
			for _, catchNode := range n.CatchBody {
				return catchNode.Evaluate(vars)
			}
		}
	}
	return nil
}

// Parser contém o estado do parser.
type Parser struct {
	tokens []lexer.Token
	pos    int
}

// NewParser cria um novo parser.
func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

// Parse analisa os tokens e retorna a AST.
func (p *Parser) Parse() []Node {
	var nodes []Node
	for p.currentToken().Type != lexer.TokenEOF {
		node := p.parseStatement()
		if node != nil {
			nodes = append(nodes, node)
		} else {
			p.pos++ // Avança para evitar loop infinito
		}
	}
	return nodes
}

func (p *Parser) currentToken() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF, Value: ""}
	}
	return p.tokens[p.pos]
}

func (p *Parser) consume(typ lexer.TokenType) lexer.Token {
	if p.currentToken().Type != typ {
		panic(fmt.Sprintf("Esperado %s, encontrado %s (valor: %s)", typ, p.currentToken().Type, p.currentToken().Value))
	}
	token := p.currentToken()
	p.pos++
	return token
}

func (p *Parser) parseStatement() Node {
	switch p.currentToken().Type {
	case lexer.TokenPrint:
		return p.parsePrint()
	case lexer.TokenAssign:
		return p.parseAssign()
	case lexer.TokenMain:
		return p.parseMain()
	case lexer.TokenTryStart:
		return p.parseTryCatch()
	case lexer.TokenFunction:
		return p.parseFunction()
	case lexer.TokenReturn:
		return p.parseReturn()
	case lexer.TokenIdentifier:
		if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == lexer.TokenEqual {
			return p.parseEqual()
		}
		if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == lexer.TokenLParen {
			return p.parseFunctionCall()
		}
		// Se for apenas um identificador, é uma referência a variável
		name := p.consume(lexer.TokenIdentifier).Value
		return &VariableNode{Name: name}
	default:
		return nil // Ignora tokens desconhecidos
	}
}

func (p *Parser) parsePrint() Node {
	p.consume(lexer.TokenPrint)

	if p.currentToken().Type != lexer.TokenString {
		panic(fmt.Sprintf("Esperado STRING após PRINT, encontrado %s (valor: %s)", p.currentToken().Type, p.currentToken().Value))
	}

	// Obter o texto da string
	text := p.consume(lexer.TokenString).Value

	// Verificar se tem interpolação
	interpolated := false
	variables := make(map[string]struct{})

	// Usar regex para encontrar todas as ocorrências de 💱{varname}
	interpolationRegex := regexp.MustCompile(`💱\{([a-zA-Z0-9_]+)\}`)
	matches := interpolationRegex.FindAllStringSubmatch(text, -1)

	if len(matches) > 0 {
		interpolated = true
		for _, match := range matches {
			variables[match[1]] = struct{}{}
		}
	}

	return &PrintNode{
		Text:         text,
		Interpolated: interpolated,
		Variables:    variables,
	}
}

func (p *Parser) parseAssign() Node {
	p.consume(lexer.TokenAssign)
	name := p.consume(lexer.TokenIdentifier).Value
	p.consume(lexer.TokenEqualSign)

	// Verificar se é uma função sendo chamada
	if p.currentToken().Type == lexer.TokenIdentifier && p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == lexer.TokenLParen {
		return &AssignNode{Name: name, Value: p.parseFunctionCall()}
	}

	if p.currentToken().Type == lexer.TokenString {
		value := p.consume(lexer.TokenString).Value
		return &AssignNode{Name: name, Value: value}
	}

	value, err := strconv.Atoi(p.consume(lexer.TokenNumber).Value)
	if err != nil {
		panic(fmt.Sprintf("Número inválido: %s", p.currentToken().Value))
	}
	return &AssignNode{Name: name, Value: value}
}

func (p *Parser) parseEqual() Node {
	name := p.consume(lexer.TokenIdentifier).Value
	p.consume(lexer.TokenEqual)
	if p.currentToken().Type == lexer.TokenString {
		value := p.consume(lexer.TokenString).Value
		return &EqualNode{Name: name, Value: value}
	}
	value, err := strconv.Atoi(p.consume(lexer.TokenNumber).Value)
	if err != nil {
		panic(fmt.Sprintf("Número inválido: %s", p.currentToken().Value))
	}
	return &EqualNode{Name: name, Value: value}
}

func (p *Parser) parseMain() Node {
	p.consume(lexer.TokenMain)
	p.consume(lexer.TokenAssign) // ◀️ tratado como ASSIGN
	p.consume(lexer.TokenAssign)
	p.consume(lexer.TokenLBrace)
	var body []Node
	for p.currentToken().Type != lexer.TokenRBrace && p.currentToken().Type != lexer.TokenEOF {
		if node := p.parseStatement(); node != nil {
			body = append(body, node)
		}
	}
	p.consume(lexer.TokenRBrace)
	return &MainNode{Body: body}
}

func (p *Parser) parseTryCatch() Node {
	p.consume(lexer.TokenTryStart)
	p.consume(lexer.TokenIdentifier) // verifyUser
	p.consume(lexer.TokenComma)
	p.consume(lexer.TokenNumber) // 2
	p.consume(lexer.TokenTry)
	p.consume(lexer.TokenLBrace)
	var tryBody []Node
	for p.currentToken().Type != lexer.TokenRBrace && p.currentToken().Type != lexer.TokenEOF {
		if node := p.parseStatement(); node != nil {
			tryBody = append(tryBody, node)
		}
	}
	p.consume(lexer.TokenRBrace)
	p.consume(lexer.TokenCatch)
	p.consume(lexer.TokenLBrace)
	var catchBody []Node
	for p.currentToken().Type != lexer.TokenRBrace && p.currentToken().Type != lexer.TokenEOF {
		if node := p.parseStatement(); node != nil {
			catchBody = append(catchBody, node)
		}
	}
	p.consume(lexer.TokenRBrace)
	return &TryCatchNode{TryBody: tryBody, CatchBody: catchBody}
}

// parseFunction analisa uma definição de função
func (p *Parser) parseFunction() Node {
	p.consume(lexer.TokenFunction)
	name := p.consume(lexer.TokenIdentifier).Value
	p.consume(lexer.TokenLParen)

	// Analisar parâmetros
	var params []string
	if p.currentToken().Type != lexer.TokenRParen {
		params = append(params, p.consume(lexer.TokenIdentifier).Value)
		for p.currentToken().Type == lexer.TokenComma {
			p.consume(lexer.TokenComma)
			params = append(params, p.consume(lexer.TokenIdentifier).Value)
		}
	}
	p.consume(lexer.TokenRParen)
	p.consume(lexer.TokenLBrace)

	// Analisar corpo da função
	var body []Node
	for p.currentToken().Type != lexer.TokenRBrace && p.currentToken().Type != lexer.TokenEOF {
		if node := p.parseStatement(); node != nil {
			body = append(body, node)
		}
	}
	p.consume(lexer.TokenRBrace)

	return &FunctionNode{Name: name, Parameters: params, Body: body}
}

// parseReturn analisa uma expressão de retorno
func (p *Parser) parseReturn() Node {
	p.consume(lexer.TokenReturn)
	return &ReturnNode{Value: p.parseExpression()}
}

// parseFunctionCall analisa uma chamada de função
func (p *Parser) parseFunctionCall() Node {
	name := p.consume(lexer.TokenIdentifier).Value
	p.consume(lexer.TokenLParen)

	// Analisar argumentos
	var args []Node
	if p.currentToken().Type != lexer.TokenRParen {
		args = append(args, p.parseExpression())
		for p.currentToken().Type == lexer.TokenComma {
			p.consume(lexer.TokenComma)
			args = append(args, p.parseExpression())
		}
	}
	p.consume(lexer.TokenRParen)

	return &FunctionCallNode{Name: name, Arguments: args}
}

// parseExpression analisa uma expressão matemática ou variável
func (p *Parser) parseExpression() Node {
	left := p.parseTerm()

	for p.currentToken().Type == lexer.TokenPlus ||
		p.currentToken().Type == lexer.TokenMult ||
		p.currentToken().Type == lexer.TokenNumPlus ||
		p.currentToken().Type == lexer.TokenConcat {

		operator := p.currentToken()
		p.pos++
		right := p.parseTerm()
		left = &BinaryOpNode{Left: left, Op: operator.Type, Right: right}
	}

	return left
}

// parseTerm analisa um termo (variável, número ou string)
func (p *Parser) parseTerm() Node {
	if p.currentToken().Type == lexer.TokenIdentifier {
		if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == lexer.TokenLParen {
			return p.parseFunctionCall()
		}
		name := p.consume(lexer.TokenIdentifier).Value
		return &VariableNode{Name: name}
	}

	if p.currentToken().Type == lexer.TokenNumber {
		value, _ := strconv.Atoi(p.consume(lexer.TokenNumber).Value)
		return &AssignNode{Value: value}
	}

	if p.currentToken().Type == lexer.TokenString {
		value := p.consume(lexer.TokenString).Value
		return &StringLiteralNode{Value: value}
	}

	panic(fmt.Sprintf("Termo inesperado: %s", p.currentToken().Value))
}

// StringLiteralNode representa uma string literal
type StringLiteralNode struct {
	Value string
}

func (n *StringLiteralNode) Evaluate(vars map[string]interface{}) interface{} {
	return n.Value
}
