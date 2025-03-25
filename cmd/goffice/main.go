package main

import (
	"fmt"
	"os"

	"goffice/internal/docx"
	"goffice/internal/parser"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("用法: ./程序名 输入文件.md 输出文件.docx")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	mdContent, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("无法读取输入文件:", err)
		return
	}

	doc := parser.ParseMarkdown(string(mdContent))
	err = docx.CreateDOCX(doc, outputFile)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Printf("成功将 %s 转换为 %s\n", inputFile, outputFile)
	}
}
