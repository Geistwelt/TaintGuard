package src

import (
	"os"
)

// MustReadFile 读取指定文件的内容，一旦出错，则直接 panic。
func MustReadFile(filePath string) []byte {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return content
}

func Trim(s string) string {
	var pre, post int
	for i, b := range s {
		if b != '\n' && b != '\t' && b != ' ' && b != '"' {
			pre = i
			break
		}
	}
	s = s[pre:]
	for i := len(s) - 1; i > 0; i-- {
		if s[i] != '\n' && s[i] != '\t' && s[i] != ' ' && s[i] != '"' {
			post = i + 1
			break
		}
	}
	return s[:post]
}
