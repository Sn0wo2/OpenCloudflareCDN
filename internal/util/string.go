package util

import (
	"net/http"
	"strings"
	"unsafe"
)

func BytesToString(v []byte) string {
	return *(*string)(unsafe.Pointer(&v)) //nolint:gosec
}

func StringToBytes(v string) []byte {
	return unsafe.Slice(unsafe.StringData(v), len(v)) //nolint:gosec
}

func TitleCase(s string) string {
	if len(s) == 0 {
		return s
	}

	if s[0] >= 'a' && s[0] <= 'z' {
		return string(s[0]-32) + s[1:]
	}

	return s
}

func GetStatusText(status int, msg ...string) string {
	var m string
	if len(msg) == 0 {
		m = http.StatusText(status)
	}

	m = strings.ToLower(m)

	return m
}
