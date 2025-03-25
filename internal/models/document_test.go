package models

import (
	"testing"
)

func TestDocumentStructure(t *testing.T) {
	// 创建一个包含各种元素的文档
	doc := Document{
		Blocks: []Block{
			Header{Level: 1, Text: "测试标题"},
			Paragraph{
				Inlines: []Inline{
					Text{Content: "这是普通文本"},
					Bold{Content: []Inline{Text{Content: "这是粗体文本"}}},
					Math{LaTeX: "E=mc^2", Display: false},
				},
			},
			Math{LaTeX: "\\int_0^\\infty e^{-x} dx = 1", Display: true},
		},
	}

	// 测试文档结构
	if len(doc.Blocks) != 3 {
		t.Errorf("期望文档包含3个块元素，实际为%d", len(doc.Blocks))
	}

	// 测试标题
	if header, ok := doc.Blocks[0].(Header); ok {
		if header.Level != 1 {
			t.Errorf("期望标题级别为1，实际为%d", header.Level)
		}
		if header.Text != "测试标题" {
			t.Errorf("期望标题文本为'测试标题'，实际为'%s'", header.Text)
		}
		if header.Type() != "header" {
			t.Errorf("期望标题类型为'header'，实际为'%s'", header.Type())
		}
	} else {
		t.Error("第一个块元素应该是Header类型")
	}

	// 测试段落
	if paragraph, ok := doc.Blocks[1].(Paragraph); ok {
		if len(paragraph.Inlines) != 3 {
			t.Errorf("期望段落包含3个内联元素，实际为%d", len(paragraph.Inlines))
		}
		if paragraph.Type() != "paragraph" {
			t.Errorf("期望段落类型为'paragraph'，实际为'%s'", paragraph.Type())
		}

		// 测试普通文本
		if text, ok := paragraph.Inlines[0].(Text); ok {
			if text.Content != "这是普通文本" {
				t.Errorf("期望文本内容为'这是普通文本'，实际为'%s'", text.Content)
			}
			if text.InlineType() != "text" {
				t.Errorf("期望内联元素类型为'text'，实际为'%s'", text.InlineType())
			}
		} else {
			t.Error("第一个内联元素应该是Text类型")
		}

		// 测试粗体文本
		if bold, ok := paragraph.Inlines[1].(Bold); ok {
			if len(bold.Content) != 1 {
				t.Errorf("期望粗体包含1个内联元素，实际为%d", len(bold.Content))
			}
			if bold.InlineType() != "bold" {
				t.Errorf("期望内联元素类型为'bold'，实际为'%s'", bold.InlineType())
			}
			if text, ok := bold.Content[0].(Text); ok {
				if text.Content != "这是粗体文本" {
					t.Errorf("期望粗体文本内容为'这是粗体文本'，实际为'%s'", text.Content)
				}
			} else {
				t.Error("粗体内容应该是Text类型")
			}
		} else {
			t.Error("第二个内联元素应该是Bold类型")
		}

		// 测试内联数学公式
		if math, ok := paragraph.Inlines[2].(Math); ok {
			if math.LaTeX != "E=mc^2" {
				t.Errorf("期望公式内容为'E=mc^2'，实际为'%s'", math.LaTeX)
			}
			if math.Display {
				t.Error("期望内联公式Display为false")
			}
			if math.InlineType() != "math" {
				t.Errorf("期望内联元素类型为'math'，实际为'%s'", math.InlineType())
			}
		} else {
			t.Error("第三个内联元素应该是Math类型")
		}
	} else {
		t.Error("第二个块元素应该是Paragraph类型")
	}

	// 测试块级数学公式
	if math, ok := doc.Blocks[2].(Math); ok {
		if math.LaTeX != "\\int_0^\\infty e^{-x} dx = 1" {
			t.Errorf("期望公式内容为'\\int_0^\\infty e^{-x} dx = 1'，实际为'%s'", math.LaTeX)
		}
		if !math.Display {
			t.Error("期望块级公式Display为true")
		}
		if math.Type() != "mathblock" {
			t.Errorf("期望块级元素类型为'mathblock'，实际为'%s'", math.Type())
		}
	} else {
		t.Error("第三个块元素应该是Math类型")
	}
}
