package simplecalculator

import (
	"compiler/lab/simplelexer"
	"testing"
)

func TestSimpleCalculator(t *testing.T) {
	cal := NewCalculator()
	lex := simplelexer.NewSimpleLexer()
	str := "1+2*3+23"
	lex.Tokenize(str)

	cal.Prog(lex.GetTokenList())
	cal.Evaluate(cal.AstTree, str)

}
