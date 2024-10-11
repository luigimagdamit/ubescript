package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		repl()
	} else if len(args) == 2 {
		readFile(args[1])
	} else if len(args) == 3 {
		fmt.Println("Usage: ube [path]")
		os.Exit(64)
	}
}
func readFile(path string) {
	source, err := preprocessFile(path)
	if err != nil {
		fmt.Println("Error reading file from main")
		os.Exit(65)
	}
	var result InterpretResult = interpret(&source)
	if result == INTERPRET_COMPILE_ERROR {
		os.Exit(65)
	}
	if result == INTERPRET_RUNTIME_ERROR {
		os.Exit(70)
	}

}
func repl() {
	reader := bufio.NewReader(os.Stdin)
	var line string
	for {
		fmt.Printf(">")
		line, _ = reader.ReadString('\n')
		line += "\x01"
		interpret(&line)

	}
}
