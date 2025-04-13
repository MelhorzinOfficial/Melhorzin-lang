package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType define os tipos de tokens.
type TokenType string

const (
	TokenPrint       TokenType = "PRINT"       // 🖨️
	TokenAssign      TokenType = "ASSIGN"      // ✍️
	TokenEqual       TokenType = "EQUAL"       // 🟰
	TokenMain        TokenType = "MAIN"        // main
	TokenTry         TokenType = "TRY"         // 👨🏿‍💻
	TokenCatch       TokenType = "CATCH"       // 🤦🏿‍♂️
	TokenTryStart    TokenType = "TRY_START"   // 🚀
	TokenFunction    TokenType = "FUNCTION"    // ▶️
	TokenReturn      TokenType = "RETURN"      // ↩️
	TokenInterpolate TokenType = "INTERPOLATE" // 💱
	TokenMult        TokenType = "MULT"        // ✖️
	TokenNumPlus     TokenType = "NUMPLUS"     // ➕
	TokenConcat      TokenType = "CONCAT"      // .
	TokenIdentifier  TokenType = "IDENTIFIER"  // nomedavariavel
	TokenString      TokenType = "STRING"      // "texto"
	TokenNumber      TokenType = "NUMBER"      // 10
	TokenBoolean     TokenType = "BOOLEAN"     // true/false
	TokenLBrace      TokenType = "LBRACE"      // {
	TokenRBrace      TokenType = "RBRACE"      // }
	TokenLParen      TokenType = "LPAREN"      // (
	TokenRParen      TokenType = "RPAREN"      // )
	TokenComma       TokenType = "COMMA"       // ,
	TokenEqualSign   TokenType = "EQUALSIGN"   // =
	TokenPlus        TokenType = "PLUS"        // +
	TokenEOF         TokenType = "EOF"

	// Tokens for type system
	TokenTypeNumber TokenType = "TYPE_NUMBER" // 🔢
	TokenTypeString TokenType = "TYPE_STRING" // 📝
	TokenTypeBool   TokenType = "TYPE_BOOL"   // ⚖️
	TokenTypeAny    TokenType = "TYPE_ANY"    // 🗑️
	TokenTypeColon  TokenType = "TYPE_COLON"  // :
)

// Token representa um token com tipo e valor.
type Token struct {
	Type  TokenType
	Value string
}

// Lexer contém o estado do lexer.
type Lexer struct {
	input  string
	pos    int
	tokens []Token
}

// NewLexer cria um novo lexer.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0, tokens: []Token{}}
}

// Lex analisa o input e retorna os tokens.
func (l *Lexer) Lex() []Token {
	for l.pos < len(l.input) {
		// Ler a próxima sequência de bytes como string para comparar emojis
		remaining := l.input[l.pos:]

		// Verificar emojis primeiro
		if strings.HasPrefix(remaining, "🖨️") {
			l.tokens = append(l.tokens, Token{Type: TokenPrint, Value: "🖨️"})
			l.pos += len("🖨️")
			continue
		}
		if strings.HasPrefix(remaining, "✍️") {
			l.tokens = append(l.tokens, Token{Type: TokenAssign, Value: "✍️"})
			l.pos += len("✍️")
			continue
		}
		if strings.HasPrefix(remaining, "🟰") {
			l.tokens = append(l.tokens, Token{Type: TokenEqual, Value: "🟰"})
			l.pos += len("🟰")
			continue
		}
		if strings.HasPrefix(remaining, "👨🏿‍💻") {
			l.tokens = append(l.tokens, Token{Type: TokenTry, Value: "👨🏿‍💻"})
			l.pos += len("👨🏿‍💻")
			continue
		}
		if strings.HasPrefix(remaining, "🤦🏿‍♂️") {
			l.tokens = append(l.tokens, Token{Type: TokenCatch, Value: "🤦🏿‍♂️"})
			l.pos += len("🤦🏿‍♂️")
			continue
		}
		if strings.HasPrefix(remaining, string('🚀')) {
			l.tokens = append(l.tokens, Token{Type: TokenTryStart, Value: string('🚀')})
			l.pos += len(string('🚀'))
			continue
		}
		if strings.HasPrefix(remaining, "▶️") {
			l.tokens = append(l.tokens, Token{Type: TokenFunction, Value: "▶️"})
			l.pos += len("▶️")
			continue
		}
		if strings.HasPrefix(remaining, "↩️") {
			l.tokens = append(l.tokens, Token{Type: TokenReturn, Value: "↩️"})
			l.pos += len("↩️")
			continue
		}
		if strings.HasPrefix(remaining, "✖️") {
			l.tokens = append(l.tokens, Token{Type: TokenMult, Value: "✖️"})
			l.pos += len("✖️")
			continue
		}
		if strings.HasPrefix(remaining, "➕") {
			l.tokens = append(l.tokens, Token{Type: TokenNumPlus, Value: "➕"})
			l.pos += len("➕")
			continue
		}
		if strings.HasPrefix(remaining, "💱") {
			// Não precisamos tokenizar o emoji de interpolação, a interpolação será tratada no parser
			l.pos += len("💱")
			continue
		}
		// Type emojis
		if strings.HasPrefix(remaining, "🔢") {
			l.tokens = append(l.tokens, Token{Type: TokenTypeNumber, Value: "🔢"})
			l.pos += len("🔢")
			continue
		}
		if strings.HasPrefix(remaining, "📝") {
			l.tokens = append(l.tokens, Token{Type: TokenTypeString, Value: "📝"})
			l.pos += len("📝")
			continue
		}
		if strings.HasPrefix(remaining, "⚖️") {
			l.tokens = append(l.tokens, Token{Type: TokenTypeBool, Value: "⚖️"})
			l.pos += len("⚖️")
			continue
		}
		if strings.HasPrefix(remaining, "🗑️") {
			l.tokens = append(l.tokens, Token{Type: TokenTypeAny, Value: "🗑️"})
			l.pos += len("🗑️")
			continue
		}

		// Ler o próximo rune
		r := rune(l.input[l.pos])

		switch {
		case r == '{':
			l.tokens = append(l.tokens, Token{Type: TokenLBrace, Value: "{"})
			l.pos++
		case r == '}':
			l.tokens = append(l.tokens, Token{Type: TokenRBrace, Value: "}"})
			l.pos++
		case r == '(':
			l.tokens = append(l.tokens, Token{Type: TokenLParen, Value: "("})
			l.pos++
		case r == ')':
			l.tokens = append(l.tokens, Token{Type: TokenRParen, Value: ")"})
			l.pos++
		case r == ',':
			l.tokens = append(l.tokens, Token{Type: TokenComma, Value: ","})
			l.pos++
		case r == '=':
			l.tokens = append(l.tokens, Token{Type: TokenEqualSign, Value: "="})
			l.pos++
		case r == '+':
			l.tokens = append(l.tokens, Token{Type: TokenPlus, Value: "+"})
			l.pos++
		case r == '*':
			l.tokens = append(l.tokens, Token{Type: TokenMult, Value: "*"})
			l.pos++
		case r == '"':
			l.pos++
			start := l.pos
			stringContent := ""

			for l.pos < len(l.input) && l.input[l.pos] != '"' {
				// Verificar se temos uma interpolação dentro da string
				if l.pos+2 < len(l.input) &&
					strings.HasPrefix(l.input[l.pos:], "💱{") {
					// Adiciona o conteúdo até aqui como uma string
					if l.pos > start {
						stringContent += l.input[start:l.pos]
					}

					// Adiciona o token de interpolação
					l.tokens = append(l.tokens, Token{Type: TokenString, Value: stringContent})
					stringContent = "" // Reseta o conteúdo

					// Avança além do 💱{
					l.pos += len("💱{")
					l.tokens = append(l.tokens, Token{Type: TokenInterpolate, Value: "💱"})
					l.tokens = append(l.tokens, Token{Type: TokenLBrace, Value: "{"})

					// Captura o identificador dentro das chaves
					identStart := l.pos
					for l.pos < len(l.input) && l.input[l.pos] != '}' {
						l.pos++
					}

					if l.pos >= len(l.input) {
						fmt.Printf("Interpolação não terminada na posição %d\n", l.pos)
						return l.tokens
					}

					// Adiciona o identificador
					identifier := l.input[identStart:l.pos]
					l.tokens = append(l.tokens, Token{Type: TokenIdentifier, Value: identifier})
					l.tokens = append(l.tokens, Token{Type: TokenRBrace, Value: "}"})

					l.pos++ // Pula o }
					start = l.pos
					continue
				}

				l.pos++
			}

			// Adiciona qualquer texto restante como uma string
			if l.pos > start {
				stringContent += l.input[start:l.pos]
			}

			if l.pos >= len(l.input) {
				fmt.Printf("String não terminada na posição %d\n", l.pos)
				return l.tokens
			}

			l.tokens = append(l.tokens, Token{Type: TokenString, Value: stringContent})
			l.pos++ // Pula o "
		case unicode.IsLetter(r):
			start := l.pos
			for l.pos < len(l.input) && (unicode.IsLetter(rune(l.input[l.pos])) || unicode.IsDigit(rune(l.input[l.pos]))) {
				l.pos++
			}
			value := l.input[start:l.pos]

			// Verificar palavras-chave e valores booleanos
			if value == "true" || value == "false" {
				l.tokens = append(l.tokens, Token{Type: TokenBoolean, Value: value})
			} else if value == "main" {
				l.tokens = append(l.tokens, Token{Type: TokenMain, Value: value})
			} else {
				l.tokens = append(l.tokens, Token{Type: TokenIdentifier, Value: value})
			}
			continue
		case unicode.IsDigit(r):
			start := l.pos
			for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
				l.pos++
			}
			l.tokens = append(l.tokens, Token{Type: TokenNumber, Value: l.input[start:l.pos]})
			continue
		case unicode.IsSpace(r) || r == '\n' || r == '\r' || r == '\t':
			l.pos++ // Ignora espaços, tabs e quebras de linha
		case r == 0xFE0F || r == 0xFEFF: // Ignora Variation Selector e BOM
			l.pos++
		case r == '.':
			l.tokens = append(l.tokens, Token{Type: TokenConcat, Value: "."})
			l.pos++
		case r == ':':
			l.tokens = append(l.tokens, Token{Type: TokenTypeColon, Value: ":"})
			l.pos++
		default:
			fmt.Printf("Caractere inesperado na posição %d: '%s' (Unicode: U+%X)\n", l.pos, string(r), r)
			l.pos++
		}
	}
	l.tokens = append(l.tokens, Token{Type: TokenEOF, Value: ""})
	fmt.Println("Tokens gerados:", l.tokens) // Log para depuração
	return l.tokens
}
