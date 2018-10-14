package ymdCppHeaderParser

import "unicode"

func isSpace(c byte) bool {
	return unicode.IsSpace(rune(c))
}

func isControl(c byte) bool {
	return unicode.IsControl(rune(c))
}

func isAlpha(c byte) bool {
	return unicode.IsLower(rune(c)) || unicode.IsUpper(rune(c))
}

func isAlnum(c byte) bool {
	return isAlpha(c) || unicode.IsNumber(rune(c))
}

func isDigit(c byte) bool {
	return unicode.IsNumber(rune(c))
}

func isXDigit(c byte) bool {
	return ('0' <= c && c <= '9') ||
		('a' <= c && c <= 'f') ||
		('A' <= c && c <= 'F')
}
