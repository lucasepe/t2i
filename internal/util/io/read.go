package io

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

func ReadInput(args []string) ([]byte, error) {
	// Caso 1 — c’è un pipe?
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return io.ReadAll(os.Stdin)
	}

	// Caso 2 — c’è un filename come primo argomento?
	if len(args) > 0 {
		return os.ReadFile(args[0])
	}

	// Caso 3 — input interattivo fino a CTRL+D
	fmt.Println("» Enter text (press CTRL+D to finish):")

	reader := bufio.NewReader(os.Stdin)
	var buf []byte

	for {
		line, err := reader.ReadBytes('\n')
		buf = append(buf, line...)

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}
