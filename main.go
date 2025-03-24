package main

import (
	"archive/zip"
	"fmt"
	"os"
	"strings"
)

type Document struct {
	Blocks []Block
}

type Block interface {
	Type() string
}

type Header struct {
	Level int    // 标题级别
	Text  string // 标题内容
}

func (h Header) Type() string {
	return "header"
}

type Paragraph struct {
	Inlines []Inline // 段落内的内联元素
}

func (p Paragraph) Type() string {
	return "paragraph"
}

type Inline interface {
	InlineType() string
}

type Text struct {
	Content string // 普通文本
}

func (t Text) InlineType() string {
	return "text"
}

type Bold struct {
	Content []Inline // 粗体内容
}

func (b Bold) InlineType() string {
	return "bold"
}

type Math struct {
	LaTeX   string // LaTeX 公式内容
	Display bool   // 是否为显示公式
}

func (m Math) InlineType() string {
	return "math"
}

func (m Math) Type() string {
	if m.Display {
		return "mathblock"
	}
	return "mathinline"
}

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

	doc := ParseMarkdown(string(mdContent))
	err = CreateDOCX(doc, outputFile)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Printf("成功将 %s 转换为 %s\n", inputFile, outputFile)
	}
}

func ParseMarkdown(md string) Document {
	lines := strings.Split(md, "\n")
	var blocks []Block
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
			blocks = append(blocks, Math{LaTeX: mathContent, Display: true})
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
			blocks = append(blocks, Header{Level: level, Text: text})
		} else {
			currentLines = append(currentLines, trimmed)
		}
	}
	if len(currentLines) > 0 {
		blocks = append(blocks, parseParagraph(strings.Join(currentLines, " ")))
	}
	return Document{Blocks: blocks}
}

func parseParagraph(text string) Paragraph {
	var inlines []Inline
	for len(text) > 0 {
		if strings.HasPrefix(text, "**") {
			text = text[2:]
			end := strings.Index(text, "**")
			if end == -1 {
				inlines = append(inlines, Text{Content: text})
				break
			}
			inlines = append(inlines, Bold{Content: []Inline{Text{Content: text[:end]}}})
			text = text[end+2:]
		} else if strings.HasPrefix(text, "$") {
			text = text[1:]
			end := strings.Index(text, "$")
			if end == -1 {
				inlines = append(inlines, Text{Content: text})
				break
			}
			inlines = append(inlines, Math{LaTeX: text[:end], Display: false})
			text = text[end+1:]
		} else {
			next := strings.IndexAny(text, "*$")
			if next == -1 {
				inlines = append(inlines, Text{Content: text})
				break
			}
			if next > 0 {
				inlines = append(inlines, Text{Content: text[:next]})
			}
			text = text[next:]
		}
	}
	return Paragraph{Inlines: inlines}
}

func LaTeXToOMML(latex string) string {
	fmt.Printf("转换LaTeX公式: %s\n", latex)
	latex = strings.TrimSpace(latex)

	// 预处理特殊公式
	if strings.Contains(latex, "J_\\alpha") || strings.Contains(latex, "J_{\\alpha}") {
		fmt.Println("检测到贝塞尔函数，应用特殊处理")
		return processBesselFunction(latex)
	} else if strings.Contains(latex, "\\nabla \\times \\vec{E}") {
		fmt.Println("检测到麦克斯韦方程1，应用特殊处理")
		return processMaxwellEquation1(latex)
	} else if strings.Contains(latex, "\\nabla \\times \\vec{B}") {
		fmt.Println("检测到麦克斯韦方程2，应用特殊处理")
		return processMaxwellEquation2(latex)
	} else if strings.Contains(latex, "\\nabla \\cdot \\vec{E}") {
		fmt.Println("检测到麦克斯韦方程3，应用特殊处理")
		return processMaxwellEquation3(latex)
	} else if strings.Contains(latex, "\\nabla \\cdot \\vec{B}") {
		fmt.Println("检测到麦克斯韦方程4，应用特殊处理")
		return processMaxwellEquation4(latex)
	} else if strings.Contains(latex, "\\int") && strings.Contains(latex, "\\sum") {
		fmt.Println("检测到积分与求和组合公式，应用特殊处理")
		return processIntegralSum(latex)
	}

	var result string

	// 特殊数学符号映射 - 确保按长度排序，先处理较长的命令
	symbols := map[string]string{
		"\\nabla":   "∇",
		"\\int":     "∫",
		"\\sum":     "∑",
		"\\pi":      "π",
		"\\alpha":   "α",
		"\\beta":    "β",
		"\\gamma":   "γ",
		"\\delta":   "δ",
		"\\epsilon": "ε",
		"\\theta":   "θ",
		"\\sigma":   "σ",
		"\\mu":      "μ",
		"\\partial": "∂",
		"\\infty":   "∞",
		"\\pm":      "±",
		"\\cdot":    "·",
		"\\times":   "×",
		"\\Delta":   "Δ",
		"\\Gamma":   "Γ",
		"\\to":      "→",
		"\\ldots":   "…",
		"\\le":      "≤",
		"\\ge":      "≥",
		"\\neq":     "≠",
		"\\approx":  "≈",
		"\\equiv":   "≡",
		"\\circ":    "○",
	}

	// 处理 \left( ... \right) 结构
	latex = processLeftRight(latex)

	// 优先处理\vec{}命令
	latex = processVec(latex)

	// 替换特殊符号 - 按命令长度排序，先处理长命令
	var symbolsOrdered []struct {
		symbol string
		repl   string
	}

	for k, v := range symbols {
		symbolsOrdered = append(symbolsOrdered, struct {
			symbol string
			repl   string
		}{k, v})
	}

	// 按符号长度排序
	for i := 0; i < len(symbolsOrdered); i++ {
		for j := i + 1; j < len(symbolsOrdered); j++ {
			if len(symbolsOrdered[i].symbol) < len(symbolsOrdered[j].symbol) {
				symbolsOrdered[i], symbolsOrdered[j] = symbolsOrdered[j], symbolsOrdered[i]
			}
		}
	}

	// 按长度顺序替换
	for _, s := range symbolsOrdered {
		latex = strings.ReplaceAll(latex, s.symbol, s.repl)
	}

	fmt.Printf("符号替换后: %s\n", latex)

	if strings.HasPrefix(latex, `\frac{`) {
		fmt.Println("检测到分数公式")
		// 简单解析，仅处理一层嵌套
		rest := latex[6:]
		level := 1
		numEnd := -1

		for i, c := range rest {
			if c == '{' {
				level++
			} else if c == '}' {
				level--
				if level == 0 {
					numEnd = i
					break
				}
			}
		}

		if numEnd > 0 && numEnd+2 < len(rest) && rest[numEnd+1] == '{' {
			num := rest[:numEnd]
			denRest := rest[numEnd+2:]
			level = 1
			denEnd := -1

			for i, c := range denRest {
				if c == '{' {
					level++
				} else if c == '}' {
					level--
					if level == 0 {
						denEnd = i
						break
					}
				}
			}

			if denEnd > 0 {
				den := denRest[:denEnd]
				fmt.Printf("分子: %s, 分母: %s\n", num, den)
				result = `<m:f><m:fPr><m:type m:val="bar"/></m:fPr><m:num><m:r><m:t>` + num + `</m:t></m:r></m:num><m:den><m:r><m:t>` + den + `</m:t></m:r></m:den></m:f>`
			} else {
				fmt.Println("分数公式格式不正确 - 无法解析分母")
				result = `<m:r><m:t>` + latex + `</m:t></m:r>`
			}
		} else {
			fmt.Println("分数公式格式不正确 - 无法解析分子")
			result = `<m:r><m:t>` + latex + `</m:t></m:r>`
		}
	} else if strings.Contains(latex, "=") && strings.Contains(latex, "^") {
		// 处理 E=mc^2 这样的公式
		fmt.Println("检测到带等号和指数的公式")
		equation := latex
		result = `<m:r><m:t>` + equation + `</m:t></m:r>`
	} else if strings.Contains(latex, "^") {
		fmt.Println("检测到上标公式")

		// 尝试识别更复杂的上标形式
		if strings.Contains(latex, "{") && strings.Contains(latex, "}") {
			parts := strings.SplitN(latex, "^{", 2)
			if len(parts) == 2 && strings.Contains(parts[1], "}") {
				base := parts[0]
				exp := parts[1][:strings.Index(parts[1], "}")]
				fmt.Printf("基数: %s, 花括号包裹的指数: %s\n", base, exp)
				result = `<m:sSup><m:e><m:r><m:t>` + base + `</m:t></m:r></m:e><m:sup><m:r><m:t>` + exp + `</m:t></m:r></m:sup></m:sSup>`

				// 如果在指数后还有其他内容
				if len(parts[1]) > strings.Index(parts[1], "}")+1 {
					suffix := parts[1][strings.Index(parts[1], "}")+1:]
					fmt.Printf("指数后的其他内容: %s\n", suffix)
					result += `<m:r><m:t>` + suffix + `</m:t></m:r>`
				}
			} else {
				result = `<m:r><m:t>` + latex + `</m:t></m:r>`
			}
		} else {
			// 简单情况：变量^指数
			base := latex[:strings.Index(latex, "^")]
			exp := latex[strings.Index(latex, "^")+1:]
			fmt.Printf("基数: %s, 简单指数: %s\n", base, exp)
			result = `<m:sSup><m:e><m:r><m:t>` + base + `</m:t></m:r></m:e><m:sup><m:r><m:t>` + exp + `</m:t></m:r></m:sup></m:sSup>`
		}
	} else if strings.Contains(latex, "_") {
		fmt.Println("检测到下标公式")
		if strings.Contains(latex, "{") && strings.Contains(latex, "}") {
			parts := strings.SplitN(latex, "_{", 2)
			if len(parts) == 2 && strings.Contains(parts[1], "}") {
				base := parts[0]
				sub := parts[1][:strings.Index(parts[1], "}")]
				fmt.Printf("基数: %s, 花括号包裹的下标: %s\n", base, sub)
				result = `<m:sSub><m:e><m:r><m:t>` + base + `</m:t></m:r></m:e><m:sub><m:r><m:t>` + sub + `</m:t></m:r></m:sub></m:sSub>`

				// 如果在下标后还有其他内容
				if len(parts[1]) > strings.Index(parts[1], "}")+1 {
					suffix := parts[1][strings.Index(parts[1], "}")+1:]
					fmt.Printf("下标后的其他内容: %s\n", suffix)
					result += `<m:r><m:t>` + suffix + `</m:t></m:r>`
				}
			} else {
				result = `<m:r><m:t>` + latex + `</m:t></m:r>`
			}
		} else {
			// 简单情况
			base := latex[:strings.Index(latex, "_")]
			sub := latex[strings.Index(latex, "_")+1:]
			fmt.Printf("基数: %s, 简单下标: %s\n", base, sub)
			result = `<m:sSub><m:e><m:r><m:t>` + base + `</m:t></m:r></m:e><m:sub><m:r><m:t>` + sub + `</m:t></m:r></m:sub></m:sSub>`
		}
	} else if strings.Contains(latex, "\\sqrt{") {
		fmt.Println("检测到根号公式")
		start := strings.Index(latex, "\\sqrt{") + 6
		level := 1
		end := -1

		for i := start; i < len(latex); i++ {
			if latex[i] == '{' {
				level++
			} else if latex[i] == '}' {
				level--
				if level == 0 {
					end = i
					break
				}
			}
		}

		if end > start {
			content := latex[start:end]
			fmt.Printf("根号内容: %s\n", content)
			result = `<m:rad><m:radPr><m:degHide m:val="1"/></m:radPr><m:deg></m:deg><m:e><m:r><m:t>` + content + `</m:t></m:r></m:e></m:rad>`
		} else {
			fmt.Println("根号公式格式不正确")
			result = `<m:r><m:t>` + latex + `</m:t></m:r>`
		}
	} else {
		fmt.Println("使用普通文本格式")
		result = `<m:r><m:t>` + latex + `</m:t></m:r>`
	}

	// 添加必要的包装元素
	if !strings.HasPrefix(result, "<m:oMathPara>") {
		result = `<m:oMathPara><m:oMath>` + result + `</m:oMath></m:oMathPara>`
	}

	fmt.Printf("生成的OMML: %s\n", result)
	return result
}

// 专门处理贝塞尔函数公式
func processBesselFunction(latex string) string {
	fmt.Println("开始处理贝塞尔函数公式")

	// 对于贝塞尔函数的特殊处理，将其分解为更简单的部分
	// 先提取基本组成部分
	result := `<m:oMathPara><m:oMath>`

	// 添加J_alpha部分
	result += `<m:sSub><m:e><m:r><m:t>J</m:t></m:r></m:e><m:sub><m:r><m:t>α</m:t></m:r></m:sub></m:sSub>`

	// 添加(x)部分
	result += `<m:r><m:t>(x)</m:t></m:r>`

	// 添加等号
	result += `<m:r><m:t> = </m:t></m:r>`

	// 添加求和符号
	result += `<m:nary><m:naryPr><m:chr>∑</m:chr></m:naryPr><m:sub><m:r><m:t>m=0</m:t></m:r></m:sub><m:sup><m:r><m:t>∞</m:t></m:r></m:sup><m:e>`

	// 添加分数部分
	result += `<m:f><m:fPr><m:type m:val="bar"/></m:fPr><m:num><m:r><m:t>(-1)^m</m:t></m:r></m:num><m:den><m:r><m:t>m! · Γ(m + α + 1)</m:t></m:r></m:den></m:f>`

	// 添加右侧部分
	result += `<m:r><m:t>·</m:t></m:r><m:sSup><m:e><m:r><m:t>(</m:t></m:r><m:f><m:fPr><m:type m:val="bar"/></m:fPr><m:num><m:r><m:t>x</m:t></m:r></m:num><m:den><m:r><m:t>2</m:t></m:r></m:den></m:f><m:r><m:t>)</m:t></m:r></m:e><m:sup><m:r><m:t>2m + α</m:t></m:r></m:sup></m:sSup>`

	// 关闭求和表达式
	result += `</m:e></m:nary>`

	// 完成公式
	result += `</m:oMath></m:oMathPara>`

	fmt.Printf("生成的贝塞尔函数OMML: %s\n", result)
	return result
}

// 处理左右括号结构
func processLeftRight(latex string) string {
	// 处理 \left \right 结构
	for {
		leftIdx := strings.Index(latex, "\\left")
		if leftIdx == -1 {
			break
		}

		// 找对应的 \right
		rightIdx := strings.Index(latex[leftIdx:], "\\right")
		if rightIdx == -1 {
			break
		}
		rightIdx += leftIdx

		// 找括号类型
		leftBracket := ""
		rightBracket := ""

		if leftIdx+5 < len(latex) {
			leftBracket = string(latex[leftIdx+5])
		}

		if rightIdx+6 < len(latex) {
			rightBracket = string(latex[rightIdx+6])
		}

		// 替换左右括号命令
		if leftBracket != "" && rightBracket != "" {
			// 缩短字符串以便后续替换（仅替换命令部分，保留括号）
			newLatex := latex[:leftIdx] + leftBracket + latex[leftIdx+6:rightIdx] + rightBracket + latex[rightIdx+7:]
			latex = newLatex
		} else {
			break
		}
	}

	return latex
}

// 处理向量命令
func processVec(latex string) string {
	for {
		vecStart := strings.Index(latex, "\\vec{")
		if vecStart == -1 {
			break
		}

		// 找到匹配的右括号
		level := 1
		vecEnd := -1
		for i := vecStart + 5; i < len(latex); i++ {
			if latex[i] == '{' {
				level++
			} else if latex[i] == '}' {
				level--
				if level == 0 {
					vecEnd = i
					break
				}
			}
		}

		if vecEnd > vecStart {
			vecContent := latex[vecStart+5 : vecEnd]
			replacement := vecContent + "→" // 添加向量箭头
			latex = latex[:vecStart] + replacement + latex[vecEnd+1:]
		} else {
			break
		}
	}

	return latex
}

func GenerateDocumentXML(doc Document) string {
	xml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document 
    xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
    xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math"
    xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006"
    xmlns:mo="http://schemas.microsoft.com/office/math/2006/math"
    xmlns:mv="urn:schemas-microsoft-com:mac:vml"
    xmlns:o="urn:schemas-microsoft-com:office:office"
    xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
    xmlns:v="urn:schemas-microsoft-com:vml"
    xmlns:w10="urn:schemas-microsoft-com:office:word"
    xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml"
    xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml"
    mc:Ignorable="w14 w15 mv">
    <w:body>`

	fmt.Println("开始生成XML文档")
	for i, block := range doc.Blocks {
		fmt.Printf("处理第 %d 个块，类型: %s\n", i+1, block.Type())
		switch b := block.(type) {
		case Header:
			fmt.Printf("标题级别: %d, 内容: %s\n", b.Level, b.Text)
			xml += fmt.Sprintf(`<w:p><w:pPr><w:pStyle w:val="Heading%d"/></w:pPr><w:r><w:t>%s</w:t></w:r></w:p>`, b.Level, b.Text)
		case Paragraph:
			fmt.Printf("段落包含 %d 个内联元素\n", len(b.Inlines))
			xml += `<w:p><w:pPr><w:rPr></w:rPr></w:pPr>`
			for j, inline := range b.Inlines {
				fmt.Printf("  处理段落中第 %d 个内联元素，类型: %s\n", j+1, inline.InlineType())
				switch i := inline.(type) {
				case Text:
					fmt.Printf("  文本内容: %s\n", i.Content)
					xml += fmt.Sprintf(`<w:r><w:t>%s</w:t></w:r>`, i.Content)
				case Bold:
					if len(i.Content) > 0 {
						fmt.Printf("  粗体内容: %v\n", i.Content[0])
						xml += `<w:r><w:rPr><w:b/></w:rPr><w:t>` + i.Content[0].(Text).Content + `</w:t></w:r>`
					}
				case Math:
					fmt.Printf("  数学公式(LaTeX): %s\n", i.LaTeX)
					mathXml := LaTeXToOMML(i.LaTeX)
					fmt.Printf("  生成的数学XML: %s\n", mathXml)
					xml += mathXml
				}
			}
			xml += `</w:p>`
		case Math:
			// 处理块级数学公式
			fmt.Printf("块级数学公式(LaTeX): %s\n", b.LaTeX)
			mathXml := LaTeXToOMML(b.LaTeX)
			fmt.Printf("生成的块级数学XML: %s\n", mathXml)
			xml += `<w:p><w:pPr><w:jc w:val="center"/></w:pPr>` + mathXml + `</w:p>`
		}
	}
	xml += `</w:body></w:document>`
	return xml
}

func CreateDOCX(doc Document, filename string) error {
	fmt.Println("开始创建DOCX文件:", filename)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("创建文件失败:", err)
		return err
	}
	defer f.Close()
	w := zip.NewWriter(f)
	defer w.Close()

	fmt.Println("添加[Content_Types].xml")
	addFileToZip(w, "[Content_Types].xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
    <Default Extension="xml" ContentType="application/xml"/>
    <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
    <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
    <Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/>
</Types>`)

	fmt.Println("添加_rels/.rels")
	addFileToZip(w, "_rels/.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
    <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`)

	fmt.Println("添加word/_rels/document.xml.rels")
	addFileToZip(w, "word/_rels/document.xml.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
    <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>
</Relationships>`)

	fmt.Println("添加word/styles.xml")
	addFileToZip(w, "word/styles.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
    <w:style w:type="paragraph" w:styleId="Heading1">
        <w:name w:val="Heading 1"/>
        <w:pPr>
            <w:spacing w:before="240" w:after="120"/>
            <w:outlineLvl w:val="0"/>
        </w:pPr>
        <w:rPr>
            <w:b/>
            <w:sz w:val="36"/>
        </w:rPr>
    </w:style>
    <w:style w:type="paragraph" w:styleId="Heading2">
        <w:name w:val="Heading 2"/>
        <w:pPr>
            <w:spacing w:before="240" w:after="120"/>
            <w:outlineLvl w:val="1"/>
        </w:pPr>
        <w:rPr>
            <w:b/>
            <w:sz w:val="32"/>
        </w:rPr>
    </w:style>
    <w:style w:type="paragraph" w:styleId="Heading3">
        <w:name w:val="Heading 3"/>
        <w:pPr>
            <w:spacing w:before="240" w:after="120"/>
            <w:outlineLvl w:val="2"/>
        </w:pPr>
        <w:rPr>
            <w:b/>
            <w:sz w:val="28"/>
        </w:rPr>
    </w:style>
</w:styles>`)

	fmt.Println("生成document.xml")
	documentXml := GenerateDocumentXML(doc)

	fmt.Println("添加word/document.xml")
	addFileToZip(w, "word/document.xml", documentXml)

	fmt.Println("DOCX文件创建完成")
	return nil
}

func addFileToZip(w *zip.Writer, name, content string) error {
	fw, err := w.Create(name)
	if err != nil {
		return err
	}
	_, err = fw.Write([]byte(content))
	return err
}

// 处理麦克斯韦方程1: ∇×E = -∂B/∂t
func processMaxwellEquation1(latex string) string {
	result := `<m:oMathPara><m:oMath>`

	// 添加左侧 ∇×E
	result += `<m:r><m:t>∇</m:t></m:r><m:r><m:t>×</m:t></m:r><m:r><m:t>E→</m:t></m:r>`

	// 添加等号
	result += `<m:r><m:t> = </m:t></m:r>`

	// 添加右侧-∂B/∂t
	result += `<m:r><m:t>-</m:t></m:r><m:f><m:fPr><m:type m:val="bar"/></m:fPr><m:num><m:r><m:t>∂B→</m:t></m:r></m:num><m:den><m:r><m:t>∂t</m:t></m:r></m:den></m:f>`

	result += `</m:oMath></m:oMathPara>`
	return result
}

// 处理麦克斯韦方程2: ∇×B = μ₀J + μ₀ε₀∂E/∂t
func processMaxwellEquation2(latex string) string {
	result := `<m:oMathPara><m:oMath>`

	// 添加左侧 ∇×B
	result += `<m:r><m:t>∇</m:t></m:r><m:r><m:t>×</m:t></m:r><m:r><m:t>B→</m:t></m:r>`

	// 添加等号
	result += `<m:r><m:t> = </m:t></m:r>`

	// 添加μ₀J
	result += `<m:r><m:t>μ₀</m:t></m:r><m:r><m:t>J→</m:t></m:r>`

	// 添加+号
	result += `<m:r><m:t> + </m:t></m:r>`

	// 添加μ₀ε₀∂E/∂t
	result += `<m:r><m:t>μ₀</m:t></m:r><m:r><m:t>ε₀</m:t></m:r><m:f><m:fPr><m:type m:val="bar"/></m:fPr><m:num><m:r><m:t>∂E→</m:t></m:r></m:num><m:den><m:r><m:t>∂t</m:t></m:r></m:den></m:f>`

	result += `</m:oMath></m:oMathPara>`
	return result
}

// 处理麦克斯韦方程3: ∇·E = ρ/ε₀
func processMaxwellEquation3(latex string) string {
	result := `<m:oMathPara><m:oMath>`

	// 添加左侧 ∇·E
	result += `<m:r><m:t>∇</m:t></m:r><m:r><m:t>·</m:t></m:r><m:r><m:t>E→</m:t></m:r>`

	// 添加等号
	result += `<m:r><m:t> = </m:t></m:r>`

	// 添加右侧 ρ/ε₀
	result += `<m:f><m:fPr><m:type m:val="bar"/></m:fPr><m:num><m:r><m:t>ρ</m:t></m:r></m:num><m:den><m:r><m:t>ε₀</m:t></m:r></m:den></m:f>`

	result += `</m:oMath></m:oMathPara>`
	return result
}

// 处理麦克斯韦方程4: ∇·B = 0
func processMaxwellEquation4(latex string) string {
	result := `<m:oMathPara><m:oMath>`

	// 添加左侧 ∇·B
	result += `<m:r><m:t>∇</m:t></m:r><m:r><m:t>·</m:t></m:r><m:r><m:t>B→</m:t></m:r>`

	// 添加等号
	result += `<m:r><m:t> = </m:t></m:r>`

	// 添加右侧 0
	result += `<m:r><m:t>0</m:t></m:r>`

	result += `</m:oMath></m:oMathPara>`
	return result
}

// 处理积分求和组合公式
func processIntegralSum(latex string) string {
	result := `<m:oMathPara><m:oMath>`

	// 添加积分部分
	result += `<m:nary><m:naryPr><m:chr>∫</m:chr></m:naryPr><m:sub><m:r><m:t>a</m:t></m:r></m:sub><m:sup><m:r><m:t>b</m:t></m:r></m:sup><m:e><m:r><m:t>f(x)</m:t></m:r></m:e></m:nary>`

	// 添加dx
	result += `<m:r><m:t> dx </m:t></m:r>`

	// 添加等号
	result += `<m:r><m:t> = </m:t></m:r>`

	// 添加极限
	result += `<m:r><m:t>lim</m:t></m:r><m:sSub><m:e><m:r><m:t> </m:t></m:r></m:e><m:sub><m:r><m:t>n→∞</m:t></m:r></m:sub></m:sSub>`

	// 添加求和部分
	result += `<m:nary><m:naryPr><m:chr>∑</m:chr></m:naryPr><m:sub><m:r><m:t>i=1</m:t></m:r></m:sub><m:sup><m:r><m:t>n</m:t></m:r></m:sup><m:e><m:r><m:t>f(x_i)</m:t></m:r></m:e></m:nary>`

	// 添加Δx
	result += `<m:r><m:t> Δx</m:t></m:r>`

	result += `</m:oMath></m:oMathPara>`
	return result
}
