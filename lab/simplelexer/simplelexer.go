package simplelexer

import (
	"fmt"
	"unicode"
)

/*
	本实验主要是词法解析器编写，书写一个词法解析器主要包含：
	1、确定状态数（包含一个初始状态），写出token的正则表达式；
	2、画出相应的有限自动机的图形
	3、根据图形直观地写出相应的代码

	本实验主要编写一个词法解析器，实现
	age >= 45
	int age = 40
	2+3*5
	词法分析
*/
type TokenType int32

const (
	Plus TokenType = iota
	Minus
	Star
	Slash

	GE
	GT
	EQ
	LE
	LT

	SemiColon
	LeftParen  //(
	RightParen //)

	Assignment // =

	If
	Else
	Int

	Identifier    //标识符
	IntLiteral    //整型字面量
	StringLiteral //字符串字面量
)

func (t TokenType) String() string {
	switch t {
	case Plus:
		return "TokenType_plus"
	case Minus:
		return "TokenType_minus"
	case Star:
		return "TokenType_star"
	case Slash:
		return "TokenType_slash"
	case GE:
		return "TokenType_GE"
	case GT:
		return "TokenType_GT"
	case SemiColon:
		return "TokenType_semicolon"
	case LeftParen:
		return "TokenType_leftparen"
	case RightParen:
		return "TokenType_rightparen"
	case Assignment:
		return "TokenType_assign"
	case IntLiteral:
		return "TokenType_intliteral"
	case Int:
		return "TokenType_int"
	case Identifier:
		return "TokenType_Identifier"
	default:
		return "no known!"
	}
}

// 用于保存词法解析器返回的token解析结果
type SimpleToken struct {
	Val   string
	Token TokenType
}

// 有限自动机状态
type DFSstatus int32

const (
	Status_Initial DFSstatus = iota // 初始状态

	Status_If
	Status_Id_if1
	Status_Id_if2
	Status_Else
	Status_Id_else1
	Status_Id_else2
	Status_Id_else3
	Status_Id_else4
	Status_Int
	Status_Id_int1
	Status_Id_int2
	Status_Id_int3
	Status_Id
	Status_GT
	Status_GE

	Status_Assignment

	Status_Plus
	Status_Minus
	Status_Star
	Status_Slash

	Status_SemiColon

	Status_LeftParen
	Status_RightParen

	Status_IntLiteral
)

type Simplelexer struct {
	tokenList []SimpleToken // 保存解析出来的Token
	buff      []rune        // 临时保存输入token的字符串
	token     SimpleToken   // 正在解析的token
}

func NewSimpleLexer() *Simplelexer {
	return &Simplelexer{}
}

/**
 * 该函数是生成token，加入list，然后更加现在字符进行状态转换函数
 * 这个初始状态有时并不做停留，它马上进入其他状态。
 * 开始解析的时候，进入初始状态；某个Token解析完毕，也进入初始状态，在这里把Token记下来，然后建立一个新的Token。
 * 返回DFA(有限自动机下一个状态)
 */

func (s *Simplelexer) InitToken(ch rune) DFSstatus {
	if len(s.buff) > 0 { // 进行了状态切换，如果此时buff中存在数据，则输出到tokenlist中保存，并新建token和buff
		s.token.Val = string(s.buff)
		s.tokenList = append(s.tokenList, s.token)

		s.buff = []rune{}
		s.token = SimpleToken{}
	}
	newstatus := Status_Initial
	if IsAlpha(ch) {
		if ch == rune('i') {
			newstatus = Status_Id_int1
		} else {
			newstatus = Status_Id
		}
		s.token.Token = Identifier
		s.buff = append(s.buff, ch)
	} else if IsDigit(ch) {
		newstatus = Status_IntLiteral
		s.token.Token = IntLiteral
		s.buff = append(s.buff, ch)
	} else if ch == rune('>') {
		newstatus = Status_GT
		s.token.Token = GT
		s.buff = append(s.buff, ch)
	} else if ch == rune('+') {
		newstatus = Status_Plus
		s.token.Token = Plus
		s.buff = append(s.buff, ch)
	} else if ch == rune('-') {
		newstatus = Status_Minus
		s.token.Token = Minus
		s.buff = append(s.buff, ch)
	} else if ch == rune('*') {
		newstatus = Status_Star
		s.token.Token = Star
		s.buff = append(s.buff, ch)
	} else if ch == rune('/') {
		newstatus = Status_Slash
		s.token.Token = Slash
		s.buff = append(s.buff, ch)
	} else if ch == rune(';') {
		newstatus = Status_SemiColon
		s.token.Token = SemiColon
		s.buff = append(s.buff, ch)
	} else if ch == rune('(') {
		newstatus = Status_LeftParen
		s.token.Token = LeftParen
		s.buff = append(s.buff, ch)
	} else if ch == rune(')') {
		newstatus = Status_RightParen
		s.token.Token = RightParen
		s.buff = append(s.buff, ch)
	} else if ch == rune('=') {
		newstatus = Status_Assignment
		s.token.Token = Assignment
		s.buff = append(s.buff, ch)
	} else {
		newstatus = Status_Initial
	}
	return newstatus
}
func (s *Simplelexer) Tokenize(str string) {
	buff_split := []rune(str)
	var ch rune
	status := Status_Initial
	for _, ch = range buff_split {
		// 先处理特殊的，然后处理简单的。实际关键字的处理可以先识别出标识符，然后判断是否是关键字即可。这里采用状态机自己处理（不太可取）
		switch status {
		case Status_Initial:
			status = s.InitToken(ch) // 在初始化状态下，通过获取的字符来判断下一个状态
			break
		case Status_Id:
			if IsAlpha(ch) || IsDigit(ch) {
				s.buff = append(s.buff, ch)
			} else {
				status = s.InitToken(ch)
			}
			break
		case Status_GT:
			if ch == rune('=') { // 如果是大于等于，则进入处理
				s.token.Token = GE
				status = Status_GE
				s.buff = append(s.buff, ch)
			} else {
				status = s.InitToken(ch)
			}
		case Status_IntLiteral: // 字面值
			if IsDigit(ch) {
				s.buff = append(s.buff, ch) //如果是数字则一直保存字面值状态
			} else {
				status = s.InitToken(ch)
			}
			break
		case Status_Id_int1:
			if ch == rune('n') {
				status = Status_Id_int2
				s.buff = append(s.buff, ch)
			} else if IsAlpha(ch) || IsDigit(ch) {
				status = Status_Id
				s.token.Token = Identifier
				s.buff = append(s.buff, ch)
			} else {
				status = s.InitToken(ch)
			}
			break
		case Status_Id_int2:
			if ch == rune('t') {
				status = Status_Id_int3
				s.buff = append(s.buff, ch)
			} else if IsAlpha(ch) || IsDigit(ch) {
				status = Status_Id // 切换回id状态
				s.token.Token = Identifier
				s.buff = append(s.buff, ch)
			} else {
				status = s.InitToken(ch)
			}
			break
		case Status_Id_int3:
			if IsBlank(ch) {
				s.token.Token = Int
				status = s.InitToken(ch)
			} else {
				status = Status_Id
				status = s.InitToken(ch)
			}
			break
		case Status_GE, Status_Assignment,
			Status_Plus, Status_Minus,
			Status_Star, Status_Slash,
			Status_SemiColon, Status_LeftParen,
			Status_RightParen:
			status = s.InitToken(ch)
			break
		default:
			break
		}
	}
	// 将最后一个token送入tokenlist
	if len(s.buff) > 0 {
		s.InitToken(ch)
	}
}
func (s *Simplelexer) Dump() {
	if len(s.tokenList) > 0 {
		for _, v := range s.tokenList {
			fmt.Printf("%s\t\t\t%s\n", v.Token, v.Val)
		}
	}
	// 清空结果
	s.tokenList = []SimpleToken{}
	s.buff = []rune{}
}

func IsAlpha(ch rune) bool {
	return unicode.IsLetter(ch)
}
func IsDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func IsBlank(ch rune) bool {
	return unicode.IsSpace(ch)
}

func main() {

	return
}
