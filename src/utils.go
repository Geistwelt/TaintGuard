package src

import (
	"os"
	"strings"
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
	s = strings.Trim(s, "\n")
	s = strings.Trim(s, "\t")
	s = strings.Trim(s, "\"")
	return s
}
