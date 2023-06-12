package src

import (
	"fmt"
	"os"
)

// MustReadFile 读取指定文件的内容，一旦出错，则直接 panic。
func MustReadFile(filePath string) []byte {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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

func EnsureDir(dir string) error {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return fmt.Errorf("could not create directory %q: %w", dir, err)
	}
	return nil
}
