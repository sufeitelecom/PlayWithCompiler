package simplescript

import (
	"fmt"
	"strconv"
)

type SimpleSript struct {
	version  bool
	valuemap map[string]int
}

func NewSimpleScript(ver bool) *SimpleSript {
	return &SimpleSript{
		version:  ver,
		valuemap: make(map[string]int),
	}
}

// 该函数就是一个计算器语法树的语义解析器
func (s *SimpleSript) Calculator(node *ASTNode, str string) int {
	result := 0
	if s.version {
		fmt.Printf("%s Calculator: %s\n", str, node.text)
	}
	switch node.Type {
	case Programm:
		for _, child := range node.children {
			result = s.Calculator(child, str)
		}
		break
	case Additive:
		child1 := node.children[0]
		result1 := s.Calculator(child1, str+"\t")
		child2 := node.children[1]
		result2 := s.Calculator(child2, str+"\t")
		if node.text == "+" {
			result = result1 + result2
		} else {
			result = result1 - result2
		}
		break
	case Multiplicative:
		child1 := node.children[0]
		result1 := s.Calculator(child1, str+"\t")
		child2 := node.children[1]
		result2 := s.Calculator(child2, str+"\t")
		if node.text == "*" {
			result = result1 * result2
		} else {
			result = result1 / result2
		}
		break
	case IntLiteral:
		result, _ = strconv.Atoi(node.text)
		break
	case Identifier:
		name := node.text
		if val, ok := s.valuemap[name]; ok {
			return val
		} else {
			fmt.Printf("variable  %s has not been set any value", name)
		}
		break
	case AssignmentStmt:
		name := node.text
		if _, ok := s.valuemap[name]; ok {
			res := s.Calculator(node.children[0], str+"\t")
			s.valuemap[name] = res
		} else {
			fmt.Printf("unknown variable  %s", name)
		}
		break
	case IntDeclaration:
		name := node.text
		if len(node.children) > 0 {
			res := s.Calculator(node.children[0], str+"\t")
			s.valuemap[name] = res
		}
		break
	default:
		break
	}
	if s.version {
		fmt.Printf("%s Result: %d\n", str, result)
	}
	return result
}
