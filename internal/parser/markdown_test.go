package parser

import (
	"reflect"
	"testing"

	"goffice/internal/models"
)

func TestParseMarkdown(t *testing.T) {
	// 测试案例1：解析标题
	t.Run("解析标题", func(t *testing.T) {
		md := "# 一级标题\n## 二级标题\n### 三级标题"
		doc := ParseMarkdown(md)

		if len(doc.Blocks) != 3 {
			t.Fatalf("期望解析出3个块元素，实际为%d", len(doc.Blocks))
		}

		// 验证第一个标题
		if header, ok := doc.Blocks[0].(models.Header); ok {
			if header.Level != 1 || header.Text != "一级标题" {
				t.Errorf("第一个标题解析错误，期望为一级标题'一级标题'，实际为%d级'%s'", header.Level, header.Text)
			}
		} else {
			t.Errorf("第一个元素应为Header类型，实际为%s", reflect.TypeOf(doc.Blocks[0]))
		}

		// 验证第二个标题
		if header, ok := doc.Blocks[1].(models.Header); ok {
			if header.Level != 2 || header.Text != "二级标题" {
				t.Errorf("第二个标题解析错误，期望为二级标题'二级标题'，实际为%d级'%s'", header.Level, header.Text)
			}
		} else {
			t.Errorf("第二个元素应为Header类型，实际为%s", reflect.TypeOf(doc.Blocks[1]))
		}

		// 验证第三个标题
		if header, ok := doc.Blocks[2].(models.Header); ok {
			if header.Level != 3 || header.Text != "三级标题" {
				t.Errorf("第三个标题解析错误，期望为三级标题'三级标题'，实际为%d级'%s'", header.Level, header.Text)
			}
		} else {
			t.Errorf("第三个元素应为Header类型，实际为%s", reflect.TypeOf(doc.Blocks[2]))
		}
	})

	// 测试案例2：解析段落
	t.Run("解析段落", func(t *testing.T) {
		md := "这是一个普通段落。\n\n这是第二个段落。"
		doc := ParseMarkdown(md)

		if len(doc.Blocks) != 2 {
			t.Fatalf("期望解析出2个块元素，实际为%d", len(doc.Blocks))
		}

		// 验证第一个段落
		if paragraph, ok := doc.Blocks[0].(models.Paragraph); ok {
			if len(paragraph.Inlines) != 1 {
				t.Fatalf("期望段落包含1个内联元素，实际为%d", len(paragraph.Inlines))
			}
			if text, ok := paragraph.Inlines[0].(models.Text); ok {
				if text.Content != "这是一个普通段落。" {
					t.Errorf("段落内容解析错误，期望为'这是一个普通段落。'，实际为'%s'", text.Content)
				}
			} else {
				t.Errorf("段落内容应为Text类型，实际为%s", reflect.TypeOf(paragraph.Inlines[0]))
			}
		} else {
			t.Errorf("第一个元素应为Paragraph类型，实际为%s", reflect.TypeOf(doc.Blocks[0]))
		}

		// 验证第二个段落
		if paragraph, ok := doc.Blocks[1].(models.Paragraph); ok {
			if len(paragraph.Inlines) != 1 {
				t.Fatalf("期望段落包含1个内联元素，实际为%d", len(paragraph.Inlines))
			}
			if text, ok := paragraph.Inlines[0].(models.Text); ok {
				if text.Content != "这是第二个段落。" {
					t.Errorf("段落内容解析错误，期望为'这是第二个段落。'，实际为'%s'", text.Content)
				}
			} else {
				t.Errorf("段落内容应为Text类型，实际为%s", reflect.TypeOf(paragraph.Inlines[0]))
			}
		} else {
			t.Errorf("第二个元素应为Paragraph类型，实际为%s", reflect.TypeOf(doc.Blocks[1]))
		}
	})

	// 测试案例3：解析粗体文本
	t.Run("解析粗体文本", func(t *testing.T) {
		md := "这是**粗体**文本。"
		doc := ParseMarkdown(md)

		if len(doc.Blocks) != 1 {
			t.Fatalf("期望解析出1个块元素，实际为%d", len(doc.Blocks))
		}

		if paragraph, ok := doc.Blocks[0].(models.Paragraph); ok {
			if len(paragraph.Inlines) != 3 {
				t.Fatalf("期望段落包含3个内联元素，实际为%d", len(paragraph.Inlines))
			}

			// 第一部分："这是"
			if text, ok := paragraph.Inlines[0].(models.Text); ok {
				if text.Content != "这是" {
					t.Errorf("第一部分文本解析错误，期望为'这是'，实际为'%s'", text.Content)
				}
			} else {
				t.Errorf("第一部分应为Text类型，实际为%s", reflect.TypeOf(paragraph.Inlines[0]))
			}

			// 第二部分："粗体"（粗体）
			if bold, ok := paragraph.Inlines[1].(models.Bold); ok {
				if len(bold.Content) != 1 {
					t.Fatalf("期望粗体包含1个内联元素，实际为%d", len(bold.Content))
				}
				if text, ok := bold.Content[0].(models.Text); ok {
					if text.Content != "粗体" {
						t.Errorf("粗体文本解析错误，期望为'粗体'，实际为'%s'", text.Content)
					}
				} else {
					t.Errorf("粗体内容应为Text类型，实际为%s", reflect.TypeOf(bold.Content[0]))
				}
			} else {
				t.Errorf("第二部分应为Bold类型，实际为%s", reflect.TypeOf(paragraph.Inlines[1]))
			}

			// 第三部分："文本。"
			if text, ok := paragraph.Inlines[2].(models.Text); ok {
				if text.Content != "文本。" {
					t.Errorf("第三部分文本解析错误，期望为'文本。'，实际为'%s'", text.Content)
				}
			} else {
				t.Errorf("第三部分应为Text类型，实际为%s", reflect.TypeOf(paragraph.Inlines[2]))
			}
		} else {
			t.Errorf("块元素应为Paragraph类型，实际为%s", reflect.TypeOf(doc.Blocks[0]))
		}
	})

	// 测试案例4：解析行内数学公式
	t.Run("解析行内数学公式", func(t *testing.T) {
		md := "爱因斯坦方程：$E=mc^2$"
		doc := ParseMarkdown(md)

		if len(doc.Blocks) != 1 {
			t.Fatalf("期望解析出1个块元素，实际为%d", len(doc.Blocks))
		}

		if paragraph, ok := doc.Blocks[0].(models.Paragraph); ok {
			if len(paragraph.Inlines) != 2 {
				t.Fatalf("期望段落包含2个内联元素，实际为%d", len(paragraph.Inlines))
			}

			// 第一部分："爱因斯坦方程："
			if text, ok := paragraph.Inlines[0].(models.Text); ok {
				if text.Content != "爱因斯坦方程：" {
					t.Errorf("第一部分文本解析错误，期望为'爱因斯坦方程：'，实际为'%s'", text.Content)
				}
			} else {
				t.Errorf("第一部分应为Text类型，实际为%s", reflect.TypeOf(paragraph.Inlines[0]))
			}

			// 第二部分：数学公式 "E=mc^2"
			if math, ok := paragraph.Inlines[1].(models.Math); ok {
				if math.LaTeX != "E=mc^2" {
					t.Errorf("数学公式解析错误，期望为'E=mc^2'，实际为'%s'", math.LaTeX)
				}
				if math.Display {
					t.Errorf("行内数学公式Display应为false，实际为%v", math.Display)
				}
			} else {
				t.Errorf("第二部分应为Math类型，实际为%s", reflect.TypeOf(paragraph.Inlines[1]))
			}
		} else {
			t.Errorf("块元素应为Paragraph类型，实际为%s", reflect.TypeOf(doc.Blocks[0]))
		}
	})

	// 测试案例5：解析块级数学公式
	t.Run("解析块级数学公式", func(t *testing.T) {
		md := "```math\n\\int_0^{\\infty} e^{-x} dx = 1\n```"
		doc := ParseMarkdown(md)

		if len(doc.Blocks) != 1 {
			t.Fatalf("期望解析出1个块元素，实际为%d", len(doc.Blocks))
		}

		if math, ok := doc.Blocks[0].(models.Math); ok {
			if math.LaTeX != "\\int_0^{\\infty} e^{-x} dx = 1" {
				t.Errorf("数学公式解析错误，期望为'\\int_0^{\\infty} e^{-x} dx = 1'，实际为'%s'", math.LaTeX)
			}
			if !math.Display {
				t.Errorf("块级数学公式Display应为true，实际为%v", math.Display)
			}
		} else {
			t.Errorf("块元素应为Math类型，实际为%s", reflect.TypeOf(doc.Blocks[0]))
		}
	})
}
