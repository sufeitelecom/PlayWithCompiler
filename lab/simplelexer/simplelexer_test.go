package simplelexer

import (
	"fmt"
	"testing"
)

func Test_Simplelexer(t *testing.T) {
	lexer := NewSimpleLexer()
	str := "int age = 45"
	fmt.Printf("\n============lexer string :%s============\n", str)
	lexer.Tokenize(str)
	lexer.Dump()

	str = "inta name = 45"
	fmt.Printf("\n============lexer string :%s============\n", str)
	lexer.Tokenize(str)
	lexer.Dump()

	str = "age >= 45"
	fmt.Printf("\n============lexer string :%s============\n", str)
	lexer.Tokenize(str)
	lexer.Dump()

	str = "age > 45"
	fmt.Printf("\n============lexer string :%s============\n", str)
	lexer.Tokenize(str)
	lexer.Dump()
}
