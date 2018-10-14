package ymdCppHeaderParser

import "bytes"

func find_first_not_of(s, sub string) (idx int, ok bool) {
	for idx, one := range []byte(s) {
		if !bytes.Contains([]byte(sub), []byte{one}) {
			return idx, true
		}
	}
	return -1, false
}

func sAppend(s string, c byte) string {
	return string(append([]byte(s), c))
}
