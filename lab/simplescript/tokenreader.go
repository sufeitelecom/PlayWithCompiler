package simplescript

import "compiler/lab/simplelexer"

type TokenReader struct {
	tokens []simplelexer.SimpleToken
	curPos int
}

func NewToKenReader(tokens []simplelexer.SimpleToken) *TokenReader {
	return &TokenReader{
		tokens: tokens,
		curPos: 0,
	}
}

// 从tokens流中返回一个token
func (t *TokenReader) Read() *simplelexer.SimpleToken {
	if t.curPos < len(t.tokens) {
		token := t.tokens[t.curPos]
		t.curPos++
		return &token
	}
	return nil
}

//从tokens流中查看一个token
func (t *TokenReader) Peak() *simplelexer.SimpleToken {
	if t.curPos < len(t.tokens) {
		token := t.tokens[t.curPos]
		return &token
	}
	return nil
}

//放回一个token
func (t *TokenReader) Unread() {
	if t.curPos > 0 {
		t.curPos--
	}
}

//得到位置pos
func (t *TokenReader) GetPosition() int {
	return t.curPos
}

//设置位置pos
func (t *TokenReader) SetPosition(pos int) {
	if pos >= 0 && pos < len(t.tokens) {
		t.curPos = pos
	}
}
