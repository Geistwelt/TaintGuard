package src

import "os"

// MustReadFile 读取指定文件的内容，一旦出错，则直接 panic。
func MustReadFile(filePath string) []byte {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return content
}
