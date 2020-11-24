package simplecalculator

import (
	"compiler/lab/simplelexer"
	"fmt"
	"strconv"
)

/*
	实现一个简单的计算器，使用如下语法规则：
	add   ->  multi | multi Plus add
	multi ->  primary | primary Star multi
    primary -> IntLiteral | Identifier
	这里的规则递归项在右边，主要是解决左递归问题，但是导致右结合性
*/

type SimpleCalculator struct {
	AstTree *ASTNode
}

func NewCalculator() *SimpleCalculator {
	return &SimpleCalculator{
		AstTree: nil,
	}
}

func (s *SimpleCalculator) primary(tokens []simplelexer.SimpleToken) (*ASTNode, []simplelexer.SimpleToken) {
	var node *ASTNode = nil
	if len(tokens) > 0 {
		if tokens[0].Token == simplelexer.IntLiteral {
			node = NewAstNode(IntLiteral, tokens[0].Val)
			tokens = append(tokens[1:]) // 消耗一个token
		} else if tokens[0].Token == simplelexer.Identifier {
			node = NewAstNode(Identifier, tokens[0].Val)
			tokens = append(tokens[1:])
		}
	}
	return node, tokens
}
func (s *SimpleCalculator) multiplicative(tokens []simplelexer.SimpleToken) (*ASTNode, []simplelexer.SimpleToken) {
	var node *ASTNode = nil
	child1, tokens := s.primary(tokens)
	node = child1
	if child1 != nil && len(tokens) > 0 {
		if tokens[0].Token == simplelexer.Star || tokens[0].Token == simplelexer.Slash {
			text := tokens[0]
			tokens = append(tokens[1:])
			child2, tokens := s.multiplicative(tokens)
			if child2 != nil {
				node = NewAstNode(Multiplicative, text.Val)
				node.children = append(node.children, child1)
				node.children = append(node.children, child2)
			}
			return node, tokens
		}
	}
	return node, tokens
}
func (s *SimpleCalculator) additive(tokens []simplelexer.SimpleToken) (*ASTNode, []simplelexer.SimpleToken) {
	var node *ASTNode = nil
	child1, tokens := s.multiplicative(tokens)
	node = child1
	if child1 != nil && len(tokens) > 0 {
		if tokens[0].Token == simplelexer.Plus || tokens[0].Token == simplelexer.Minus {
			text := tokens[0]
			tokens = append(tokens[1:])
			child2, tokens := s.additive(tokens)
			if child2 != nil {
				node = NewAstNode(Additive, text.Val)
				node.children = append(node.children, child1)
				node.children = append(node.children, child2)
			}
			return node, tokens
		}
	}

	return node, tokens
}
func (s *SimpleCalculator) Prog(tokens []simplelexer.SimpleToken) {
	s.AstTree = nil
	node := NewAstNode(Programm, "calculator")
	child, tokens := s.additive(tokens) // 因为所有的算数表达式最后都可以变为加法ast语法树
	if child != nil {
		node.children = append(node.children, child)
	}
	s.AstTree = node
}
func (s *SimpleCalculator) Evaluate(node *ASTNode, str string) int {
	result := 0
	if node.Type == Programm {
		fmt.Printf("%s calculating:\n", str)
	}
	switch node.Type {
	case Programm:
		for _, child := range node.children {
			result = s.Evaluate(child, str)
		}
		break
	case Additive:
		child1 := node.children[0]
		result1 := s.Evaluate(child1, str)
		child2 := node.children[1]
		result2 := s.Evaluate(child2, str)
		if node.text == "+" {
			result = result1 + result2
		} else {
			result = result1 - result2
		}
		break
	case Multiplicative:
		child1 := node.children[0]
		result1 := s.Evaluate(child1, str)
		child2 := node.children[1]
		result2 := s.Evaluate(child2, str)
		if node.text == "*" {
			result = result1 * result2
		} else {
			result = result1 / result2
		}
		break
	case IntLiteral:
		result, _ = strconv.Atoi(node.text)
		break
	}
	fmt.Printf("%s result: %d\n", node.text, result)
	return result
}

type AstNodeType int32

const (
	Programm AstNodeType = iota //程序入口，根节点

	IntDeclaration //整型变量声明
	ExpressionStmt //表达式语句，即表达式后面跟个分号
	AssignmentStmt //赋值语句

	Primary        //基础表达式
	Multiplicative //乘法表达式
	Additive       //加法表达式

	Identifier //标识符
	IntLiteral //整型字面量
)

type ASTNode struct {
	Type     AstNodeType
	Parent   *ASTNode
	children []*ASTNode
	text     string
}

func NewAstNode(nodetype AstNodeType, text string) *ASTNode {
	return &ASTNode{
		Type: nodetype,
		text: text,
	}
}
