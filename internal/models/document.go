package models

// Document 表示一个文档对象，包含多个块元素
type Document struct {
	Blocks []Block
}

// Block 是文档中的块级元素接口
type Block interface {
	Type() string
}

// Header 表示标题元素
type Header struct {
	Level int    // 标题级别
	Text  string // 标题内容
}

// Type 返回块类型
func (h Header) Type() string {
	return "header"
}

// Paragraph 表示段落元素
type Paragraph struct {
	Inlines []Inline // 段落内的内联元素
}

// Type 返回块类型
func (p Paragraph) Type() string {
	return "paragraph"
}

// Inline 是文档中的内联元素接口
type Inline interface {
	InlineType() string
}

// Text 表示普通文本
type Text struct {
	Content string // 普通文本
}

// InlineType 返回内联元素类型
func (t Text) InlineType() string {
	return "text"
}

// Bold 表示粗体文本
type Bold struct {
	Content []Inline // 粗体内容
}

// InlineType 返回内联元素类型
func (b Bold) InlineType() string {
	return "bold"
}

// Math 表示数学公式
type Math struct {
	LaTeX   string // LaTeX 公式内容
	Display bool   // 是否为显示公式
}

// InlineType 返回内联元素类型
func (m Math) InlineType() string {
	return "math"
}

// Type 返回块类型（如果是块级数学公式）
func (m Math) Type() string {
	if m.Display {
		return "mathblock"
	}
	return "mathinline"
}
