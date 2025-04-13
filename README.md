# Melhorzin-lang
Uma linguagem de programação baseada em emojis, escrita em Go.

## Instalação
1. Clone o repositório: `git clone <url>`
2. Compile: `go build -o emojilang cmd/interpreter/main.go`
3. Execute: `./emojilang examples/hello.mlz`

## Recursos da Linguagem

### Impressão
```emoji
🖨️ "Hello World"
```

### Interpolação de Strings
```emoji
✍️ nome = "Melhorzin"
🖨️ "Olá 💱{nome}!"
```

### Definição de Variáveis
```emoji
✍️ nome = "Valor"
✍️ numero = 42
```

### Sistema de Tipos
```emoji
// Tipos inferidos automaticamente
✍️ nome = "Melhorzin"   // Tipo String inferido automaticamente
✍️ idade = 25          // Tipo Number inferido automaticamente

// Tipos explícitos
✍️ pontos:🔢 = 100              // Número (🔢)
✍️ mensagem:📝 = "Olá mundo!"   // String (📝)
✍️ ativo:⚖️ = true              // Boolean (⚖️)
✍️ qualquer:🗑️ = "qualquer coisa"  // Any/Qualquer (🗑️)
```

### Funções com Tipos
```emoji
// Função com parâmetros e retorno tipados
▶️ soma(a:🔢, b:🔢):🔢 {
    ↩️ a + b
}

// Função com tipo de retorno string
▶️ saudacao(nome:📝):📝 {
    ↩️ "Olá, " . nome . "!"
}
```

### Operações
```emoji
+   # Soma numérica (anteriormente ➕)
*   # Multiplicação (anteriormente ✖️)
.   # Concatenação de strings (substitui o +)
```

### Funções
```emoji
▶️ soma(a, b) {
    🖨️ "Somando valores..."
    ↩️ a + b
}

▶️ multiplica(a, b) {
    ↩️ a * b
}

▶️ concatena(a, b) {
    ↩️ a . b
}

✍️ resultado = soma(5, 10)
```

## Exemplos
Veja pasta `examples/` para exemplos completos.