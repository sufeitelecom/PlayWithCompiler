package main

import (
	"bufio"
	"compiler/lab/simplescript"
	"fmt"
	"os"
)

func main() {
	fmt.Printf("\n>")
	parse := simplescript.NewSimpleParser()
	scripe := simplescript.NewSimpleScript(true)
	reader := bufio.NewReader(os.Stdin)
	for {
		strBytes, _, err := reader.ReadLine()
		if err != nil {
			fmt.Printf("input error %s", err.Error())
			return
		}
		str := string(strBytes)
		if str == "exit()" {
			return
		}
		tree := parse.Parse(str)
		if tree != nil {
			scripe.Calculator(tree, str)
		}
		fmt.Printf("\n>")
		strBytes = []byte{}
	}
}
