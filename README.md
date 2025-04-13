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