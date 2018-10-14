package ymdCppHeaderParser

import (
	"testing"
)

func TestTokenizer_MatchIdentifier(t *testing.T) {
	tn := NewTokenizer([]byte(`class `), 1)
	ok := tn.MatchIdentifier(`class`)
	assert(ok)
}

func TestTokenizer_GetChar(t *testing.T) {
	tn := NewTokenizer([]byte(` int cc = 1; `), 1)
	assert(' ' == tn.GetChar())
}

func TestTokenizer_GetIdentifier(t *testing.T) {
	tn := NewTokenizer([]byte(` token `), 1)
	var token Token
	ok := tn.GetIdentifier(&token)
	assert(ok)
}

func TestTokenizer_GetLeadingChar(t *testing.T) {
	tn := NewTokenizer([]byte(`/* comment data 
*/ user `), 1)
	c := tn.GetLeadingChar()
	assert(c == 'u', c)
}

func TestTokenizer_GetToken(t *testing.T) {
	tn := NewTokenizer([]byte(`int j = 10;`), 1)
	var token Token
	ok := tn.GetToken(&token, false, false)
	assert(ok)
	assert(token.Mtoken == `int`)
	assert(token.MtokenType == kIdentifier)
}

func TestTokenizer_MatchSymbol(t *testing.T) {
	tn := NewTokenizer([]byte(`;`), 1)
	ok := tn.MatchSymbol(`;`)
	assert(ok)
}

func TestTokenizer_RequireSymbol(t *testing.T) {
	tn := NewTokenizer([]byte(`;`), 1)
	tn.RequireSymbol(`;`)
}

func TestTokenizer_UngetChar(t *testing.T) {
	tn := NewTokenizer([]byte(`if (i > 10 ) {}`), 1)
	c1 := tn.GetChar()
	tn.UngetChar()
	c2 := tn.GetChar()
	assert(c1 == c2)
}

func TestTokenizer_UngetToken(t *testing.T) {
	tn := NewTokenizer([]byte(`int input = 0; `), 1)
	var token Token
	ok := tn.GetToken(&token, false, false)
	assert(ok)
	tn.UngetToken(&token)
	assert(tn.cursorPos == 0)
	assert(tn.cursorLine == 1)
}

func TestNewTokenizer(t *testing.T) {
	tn := NewTokenizer([]byte(`sx`), 1)
	assert(tn.cursorLine == 1)
	assert(tn.cursorPos == 0)
}

func TestTokenizer_GetToken2(t *testing.T) {
	tn := NewTokenizer([]byte(` `), 1)
	var token Token
	ok := tn.GetToken(&token, false, false)
	assert(!ok)
}