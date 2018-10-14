package ymdCppHeaderParser

import (
	"bytes"
	"strconv"
	"fmt"
)

type Tokenizer struct {
	input          []byte
	cursorPos      int
	cursorLine     int
	prevCursorPos  int
	prevCursorLine int
	comment        Comment
	lastComment    Comment
}

const (
	EndOfFileChar byte = 0xFF
)

func NewTokenizer(input []byte, startingLine int) *Tokenizer {
	this := &Tokenizer{}
	this.input = input
	this.cursorPos = 0
	this.cursorLine = startingLine
	return this
}

func (this *Tokenizer) GetChar() byte {
	c := this.input[this.cursorPos]
	this.prevCursorPos = this.cursorPos
	this.prevCursorLine = this.cursorLine

	//New line moves the cursor to the new line
	if c == '\n' {
		this.cursorLine++
	}
	this.cursorPos++
	return c
}

func (this *Tokenizer) UngetChar() {
	this.cursorLine = this.prevCursorLine
	this.cursorPos = this.prevCursorPos
}

func (this *Tokenizer) peek() byte {
	if this.is_eof() {
		return EndOfFileChar
	}
	return this.input[this.cursorPos]
}

func (this *Tokenizer) GetLeadingChar() byte {
	if this.comment.text == `` {
		this.lastComment = this.comment
	}

	this.comment.text = ``
	this.comment.startLine = this.cursorLine
	this.comment.endLine = this.cursorLine

	var c = EndOfFileChar
	for !this.is_eof() {
		c = this.GetChar()
		if this.is_eof() {
			// If this is a whitespace character skip it
			if isSpace(c) || isControl(c) || c == '\n' {
				c = EndOfFileChar
			}
			break
		}
		if c == '\n' {
			if this.comment.text != `` {
				this.comment.text += "\n"
			}
			continue
		}
		if isSpace(c) || isControl(c) {
			continue
		}
		// If this is a single line comment
		next := this.peek()
		if c == '/' && next == '/' {
			lines := []string{}
			indentationLastLine := 0

			for !this.is_eof() && c == '/' && next == '/' {
				// Search for the end of the line
				line := ``
				for c = this.GetChar(); c != EndOfFileChar && c != '\n'; c = this.GetChar() {
					line = sAppend(line, c)
				}
				// Store the line
				if lastSlashIndex, ok := find_first_not_of(line, "/"); !ok {
					line = ""
				} else {
					line = line[:lastSlashIndex]
				}

				firstCharIndex, ok := find_first_not_of(line, " \t")
				if ok {
					line = ""
				} else {
					line = line[:firstCharIndex]
				}

				if firstCharIndex > indentationLastLine && len(lines) != 0 {
					lines[len(lines)-1 ] += " " + line
				} else {
					lines = append(lines, line)
					indentationLastLine = firstCharIndex
				}

				// Check the next line
				for c = this.GetChar(); !this.is_eof() && isSpace(c); c = this.GetChar() {
					// nothing
				}

				if !this.is_eof() {
					next = this.peek()
				}
			}

			if !this.is_eof() {
				this.UngetChar()
			}
			// Build comment string
			ss := bytes.NewBuffer(nil)
			for i := 0; i < len(lines); i++ {
				if i > 0 {
					ss.WriteByte('\n')
				}
				ss.WriteString(lines[i])
			}
			this.comment.text = ss.String()
			this.comment.endLine = this.cursorLine
			// Go to the next
			continue
		}
		// If this is a block comment
		if c == '/' && next == '*' {
			// Search for the end of the block comment
			lines := []string{}
			line := ``

			c = this.GetChar()
			next = this.peek()
			for c != EndOfFileChar && (c != '*' || next != '/') {
				if c == '\n' {
					if len(lines) != 0 || line == `` {
						lines = append(lines, line)
					}
					line = ``
				} else {
					if line != `` || !isSpace(c) || c == '*' {
						line = sAppend(line, c)
					}
				}
				c = this.GetChar()
				next = this.peek()
			}

			// Skip past the slash
			if c != EndOfFileChar {
				this.GetChar()
			}

			// Skip past new lines and spaces
			c = this.GetChar()
			for !this.is_eof() && isSpace(c) {
				c = this.GetChar()
			}
			if !this.is_eof() {
				this.UngetChar()
			}

			// Remove empty lines from the back
			for len(lines) != 0 && lines[len(lines)-1] == `` {
				lines = lines[len(lines)-1:]
			}
			// Build comment string
			ss := bytes.NewBuffer(nil)
			for i := 0; i < len(lines); i++ {
				if i > 0 {
					ss.WriteString("\n")
				}
				ss.WriteString(lines[i])
			}

			this.comment.text = ss.String()
			this.comment.endLine = this.cursorLine

			// Move to the next character
			continue
		}
		break
	}
	return c
}

func (this *Tokenizer) GetToken(token *Token, angleBracketsForStrings bool, seperateBraces bool) bool {
	// Get the next character
	c := this.GetLeadingChar()
	p := this.peek()

	if c == EndOfFileChar {
		return false
	}

	// Record the start of the token position
	token.MstartPos = this.prevCursorPos
	token.MstartLine = this.prevCursorLine
	token.Mtoken = ``
	token.MtokenType = kNone

	// Alphanumeric token
	if isAlpha(c) || c == '_' {
		// Read the rest of the alphanumeric characters
		for {
			token.Mtoken = sAppend(token.Mtoken, c)
			c = this.GetChar()
			if isAlnum(c) || c == '_' { // 字母数字或 _
				continue
			} else {
				// Put back the last read character since it's not part of the identifier
				this.UngetChar()
				break
			}
		}

		// Set the type of the token
		token.MtokenType = kIdentifier

		if token.Mtoken == `true` {
			token.MtokenType = kConst
			token.MconstType = kBoolean
			token.MboolConst = true
		} else if token.Mtoken == `false` {
			token.MtokenType = kConst
			token.MconstType = kBoolean
			token.MboolConst = false
		}

		return true
	} else if isDigit(c) || ((c == '-' || c == '+') && isDigit(p)) { // Constant
		isFloat := false
		isHex := false
		for {
			if c == '.' {
				isFloat = true
			}

			if c == 'x' || c == 'X' {
				isHex = true
			}
			token.Mtoken = sAppend(token.Mtoken, c)
			c = this.GetChar()

			if isDigit(c) {
				continue
			}
			if (!isFloat && c == '.') {
				continue
			}
			if (!isHex && (c == 'X' || c == 'x')) {
				continue
			}
			if (isHex && isXDigit(c)) {
				continue
			}
			break
		}

		if !isFloat || (c != 'f' && c != 'F') {
			this.UngetChar()
		}

		token.MtokenType = kConst
		if !isFloat {
			i, err := strconv.ParseInt(token.Mtoken, 0, 64)
			assert(err == nil, token.Mtoken)
			token.Mint64Const = i
			token.MconstType = kInt64
		} else {
			f, err := strconv.ParseFloat(token.Mtoken, 64)
			assert(err == nil, token.Mtoken)
			token.Mfloat64Const = f
			token.MconstType = kFloat64
		}
		return true
	} else if c == '"' || (angleBracketsForStrings && c == '<') {
		var closingElement byte = '"'
		if c != '"' {
			closingElement = '>'
		}

		c = this.GetChar()
		for c != closingElement && c != EndOfFileChar {
			if c == '\\' {
				c = this.GetChar()
				if c == EndOfFileChar {
					break
				} else if c == 'n' {
					c = '\n'
				} else if c == 't' {
					c = '\t'
				} else if c == 'r' {
					c = '\r'
				} else if c == '"' {
					c = '"'
				}
			}

			token.Mtoken = sAppend(token.Mtoken, c)
			c = this.GetChar()
		}
		if c != closingElement {
			this.UngetChar()
		}

		token.MtokenType = kConst
		token.MconstType = kString
		token.MstringConst = token.Mtoken

		return true
	} else { // Symbol
		// Push back the symbol
		token.MtokenType = kSymbol
		token.Mtoken = sAppend(token.Mtoken, c)
		d := this.peek()
		switch string([]byte{c, d}) {
		case
			`<>`, `->`, `!=`, `<=`, `>=`, `++`, `--`,
			`+=`, `-=`, `*=`, `/=`, `^=`, `|=`, `&=`,
			`~=`, `%=`, `||`, `==`, `::`:
			token.Mtoken = sAppend(token.Mtoken, d)
			this.GetChar()
		case `>>`:
			if !seperateBraces {
				token.Mtoken = sAppend(token.Mtoken, d)
				this.GetChar()
			}
		}
		return true
	}
	return false
}

func (this *Tokenizer) is_eof() bool {
	return this.cursorPos >= len(this.input)
}

func (this *Tokenizer) GetIdentifier(token *Token) bool {
	if !this.GetToken(token, false, false) {
		return false
	}

	if token.MtokenType == kIdentifier {
		return true
	}

	this.UngetToken(token)
	return false
}

func (this *Tokenizer) UngetToken(token *Token) {
	this.cursorLine = token.MstartLine
	this.cursorPos = token.MstartPos
}

func (this *Tokenizer) MatchIdentifier(identifier string) bool {
	var token Token
	if this.GetToken(&token, false, false) {
		if token.MtokenType == kIdentifier && token.Mtoken == identifier {
			return true
		}
		this.UngetToken(&token)
	}
	return false
}

func (this *Tokenizer) MatchSymbol(symbol string) bool {
	var token Token
	if this.GetToken(&token, false, len([]byte(symbol)) == 1 && symbol[0] == '>') {
		if token.MtokenType == kSymbol && token.Mtoken == symbol {
			return true
		}

		this.UngetToken(&token)
	}

	return false
}

func (this *Tokenizer) RequireSymbol(symbol string) {
	if !this.MatchSymbol(symbol) {
		this.panicf(`8jn0qzkn`, `Missing symbol %v`, symbol)
	}
}

func (this *Tokenizer) panicf(funcId string, msg string, a ... interface{}) {
	pre := this.cursorPos
	if this.cursorPos > 5 {
		pre = 5
	}
	as := a
	as = append(as, string(this.input[this.cursorPos-pre:this.cursorPos]))
	panic(fmt.Errorf(funcId+msg+" after \"%v\"\n", as...))
}
