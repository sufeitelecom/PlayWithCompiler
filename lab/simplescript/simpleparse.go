package simplescript

import (
	"compiler/lab/simplelexer"
	"fmt"
	"strconv"
)

/*
	该解析器，只对语法树进行解析，至于各节点的含义由语义解析负责
    它支持的语法规则为：
 *
 * programm -> intDeclare | expressionStatement | assignmentStatement
 * intDeclare -> 'int' Id ( = additive) ';'
 * expressionStatement -> addtive ';'
 * addtive -> multiplicative ( (+ | -) multiplicative)*
 * multiplicative -> primary ( (* | /) primary)*
 * primary -> IntLiteral | Id | (additive)
*/
type SimpleParser struct {
	AstTree *ASTNode
}

func NewSimpleParser() *SimpleParser {
	return &SimpleParser{
		AstTree: nil,
	}
}

func (s *SimpleParser) primary(tokens *TokenReader) *ASTNode {
	var node *ASTNode = nil
	token := tokens.Peak()
	if token != nil {
		if token.Token == simplelexer.IntLiteral {
			token = tokens.Read()
			node = NewAstNode(IntLiteral, token.Val)
		} else if token.Token == simplelexer.Identifier {
			token = tokens.Read()
			node = NewAstNode(Identifier, token.Val)
		} else if token.Token == simplelexer.LeftParen {
			tokens.Read()
			node = s.additive(tokens)
			if node != nil {
				token = tokens.Peak()
				if token != nil && token.Token == simplelexer.RightParen {
					tokens.Read()
				} else {
					fmt.Printf("expecting right parenthesis")
					return nil
				}
			} else {
				fmt.Printf("expecting an additive expression inside parenthesis")
				return nil
			}
		}
	}
	return node
}
func (s *SimpleParser) multiplicative(tokens *TokenReader) *ASTNode {
	var node *ASTNode = nil
	child1 := s.primary(tokens)
	node = child1
	for {
		token := tokens.Peak()
		if token != nil && (token.Token == simplelexer.Star || token.Token == simplelexer.Slash) {
			token = tokens.Read()
			child2 := s.multiplicative(tokens)
			if child2 != nil {
				node = NewAstNode(Multiplicative, token.Val)
				node.children = append(node.children, child1)
				node.children = append(node.children, child2)
				child1 = node
			} else {
				fmt.Printf("invalid multiplicative expression, expecting the right part.")
				return nil
			}
		} else {
			break
		}
	}
	return node
}
func (s *SimpleParser) additive(tokens *TokenReader) *ASTNode {
	var node *ASTNode = nil
	child1 := s.multiplicative(tokens)
	node = child1
	if child1 != nil {
		for {
			token := tokens.Peak()
			if token != nil && (token.Token == simplelexer.Plus || token.Token == simplelexer.Minus) {
				token = tokens.Read()
				child2 := s.multiplicative(tokens)
				if child2 != nil {
					node = NewAstNode(Additive, token.Val)
					node.children = append(node.children, child1)
					node.children = append(node.children, child2)
					child1 = node
				} else {
					fmt.Printf("invalid additive expression, expecting the right part.")
					return nil
				}
			} else {
				break
			}
		}
	}
	return node
}
func (s *SimpleParser) IntDeclare(tokens *TokenReader) *ASTNode {
	var node *ASTNode = nil
	token := tokens.Peak()
	if token != nil && token.Token == simplelexer.Int {
		token = tokens.Read() // 取出int关键字
		token = tokens.Peak() // 预读下一个
		if token != nil && token.Token == simplelexer.Identifier {
			token = tokens.Read()
			node = NewAstNode(IntDeclaration, token.Val)
			token = tokens.Peak() //预读等号
			if token != nil && token.Token == simplelexer.Assignment {
				token = tokens.Read()
				child := s.additive(tokens) // 构建add表达式
				if child == nil {
					fmt.Printf("invalide variable initialization, expecting an expression")
					return nil
				} else {
					node.children = append(node.children, child)
				}
			}
		} else {
			fmt.Printf("variable name expected")
			return nil
		}
		if node != nil {
			token = tokens.Peak()
			if token != nil && token.Token == simplelexer.SemiColon {
				tokens.Read()
			} else {
				fmt.Printf("invalid statement, expecting semicolon")
				return nil
			}
		}
	}
	return node
}
func (s *SimpleParser) AssignmentStatement(tokens *TokenReader) *ASTNode {
	var node *ASTNode = nil
	token := tokens.Peak()
	if token != nil && token.Token == simplelexer.Identifier {
		token = tokens.Read() //读入标识符
		node = NewAstNode(AssignmentStmt, token.Val)
		token = tokens.Peak() //预读，看看下面是不是等号
		if token != nil && token.Token == simplelexer.Assignment {
			tokens.Read() //取出等号
			child := s.additive(tokens)
			if child == nil { //出错，等号右面没有一个合法的表达式
				fmt.Printf("invalide assignment statement, expecting an expression")
				return nil
			} else {
				node.children = append(node.children, child) //添加子节点
				token = tokens.Peak()                        //预读，看看后面是不是分号
				if token != nil && token.Token == simplelexer.SemiColon {
					tokens.Read() //消耗掉这个分号
				} else { //报错，缺少分号
					fmt.Printf("invalid statement, expecting semicolon")
					return nil
				}
			}
		} else {
			tokens.Unread() //回溯，吐出之前消化掉的标识符
			node = nil
		}
	}
	return node
}
func (s *SimpleParser) ExpressionStatement(tokens *TokenReader) *ASTNode {
	var node *ASTNode = nil
	pos := tokens.GetPosition()
	node = s.additive(tokens)
	if node != nil {
		token := tokens.Peak()
		if token != nil && token.Token == simplelexer.SemiColon {
			tokens.Read()
		} else {
			node = nil
			tokens.SetPosition(pos) // 回溯
		}
	}
	return node
}
func (s *SimpleParser) Prog(tokens *TokenReader) *ASTNode {
	s.AstTree = nil
	node := NewAstNode(Programm, "pwc")
	for tokens.Peak() != nil {
		child := s.IntDeclare(tokens)

		if child == nil {
			child = s.ExpressionStatement(tokens)
		}
		if child == nil {
			child = s.AssignmentStatement(tokens)
		}

		if child != nil {
			node.children = append(node.children, child)
		} else {
			fmt.Printf("unknown statement")
			return nil
		}
	}
	return node
}
func (s *SimpleParser) Parse(str string) *ASTNode {
	lexer := simplelexer.NewSimpleLexer()
	lexer.Tokenize(str)
	tokenreader := NewToKenReader(lexer.GetTokenList())
	node := s.Prog(tokenreader)
	return node
}

func (s *SimpleParser) Evaluate(node *ASTNode, str string) int {
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
