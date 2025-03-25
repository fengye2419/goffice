package parser

import (
	"fmt"
	"strings"

	"goffice/internal/models"
)

// ParseMarkdown 将Markdown文本解析为文档模型
func ParseMarkdown(md string) models.Document {
	lines := strings.Split(md, "\n")
	var blocks []models.Block
	var currentLines []string
	var inMathBlock bool = false
	var mathContent string = ""

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// 检测数学代码块开始
		if strings.HasPrefix(trimmed, "```math") {
			if len(currentLines) > 0 {
				blocks = append(blocks, parseParagraph(strings.Join(currentLines, " ")))
				currentLines = nil
			}
			inMathBlock = true
			mathContent = ""
			continue
		}

		// 检测数学代码块结束
		if inMathBlock && strings.HasPrefix(trimmed, "```") {
			fmt.Printf("检测到块级数学公式: %s\n", mathContent)
			blocks = append(blocks, models.Math{LaTeX: mathContent, Display: true})
			inMathBlock = false
			continue
		}

		// 收集数学代码块内容
		if inMathBlock {
			if mathContent != "" {
				mathContent += " "
			}
			mathContent += trimmed
			continue
		}

		// 正常Markdown解析
		if trimmed == "" {
			if len(currentLines) > 0 {
				blocks = append(blocks, parseParagraph(strings.Join(currentLines, " ")))
				currentLines = nil
			}
		} else if strings.HasPrefix(trimmed, "#") {
			if len(currentLines) > 0 {
				blocks = append(blocks, parseParagraph(strings.Join(currentLines, " ")))
				currentLines = nil
			}
			level := 0
			for _, c := range trimmed {
				if c == '#' {
					level++
				} else {
					break
				}
			}
			text := strings.TrimSpace(trimmed[level:])
			blocks = append(blocks, models.Header{Level: level, Text: text})
		} else {
			currentLines = append(currentLines, trimmed)
		}
	}
	if len(currentLines) > 0 {
		blocks = append(blocks, parseParagraph(strings.Join(currentLines, " ")))
	}
	return models.Document{Blocks: blocks}
}

// parseParagraph 解析段落文本，识别内联元素
func parseParagraph(text string) models.Paragraph {
	var inlines []models.Inline
	for len(text) > 0 {
		if strings.HasPrefix(text, "**") {
			text = text[2:]
			end := strings.Index(text, "**")
			if end == -1 {
				inlines = append(inlines, models.Text{Content: text})
				break
			}
			inlines = append(inlines, models.Bold{Content: []models.Inline{models.Text{Content: text[:end]}}})
			text = text[end+2:]
		} else if strings.HasPrefix(text, "$") {
			text = text[1:]
			end := strings.Index(text, "$")
			if end == -1 {
				inlines = append(inlines, models.Text{Content: text})
				break
			}
			inlines = append(inlines, models.Math{LaTeX: text[:end], Display: false})
			text = text[end+1:]
		} else {
			next := strings.IndexAny(text, "*$")
			if next == -1 {
				inlines = append(inlines, models.Text{Content: text})
				break
			}
			if next > 0 {
				inlines = append(inlines, models.Text{Content: text[:next]})
			}
			text = text[next:]
		}
	}
	return models.Paragraph{Inlines: inlines}
}
