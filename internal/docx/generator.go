package docx

import (
	"archive/zip"
	"fmt"
	"os"

	"goffice/internal/models"
	"goffice/pkg/latex"
)

// GenerateDocumentXML 将文档模型转换为XML
func GenerateDocumentXML(doc models.Document) string {
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
		case models.Header:
			fmt.Printf("标题级别: %d, 内容: %s\n", b.Level, b.Text)
			xml += fmt.Sprintf(`<w:p><w:pPr><w:pStyle w:val="Heading%d"/></w:pPr><w:r><w:t>%s</w:t></w:r></w:p>`, b.Level, b.Text)
		case models.Paragraph:
			fmt.Printf("段落包含 %d 个内联元素\n", len(b.Inlines))
			xml += `<w:p><w:pPr><w:rPr></w:rPr></w:pPr>`
			for j, inline := range b.Inlines {
				fmt.Printf("  处理段落中第 %d 个内联元素，类型: %s\n", j+1, inline.InlineType())
				switch i := inline.(type) {
				case models.Text:
					fmt.Printf("  文本内容: %s\n", i.Content)
					xml += fmt.Sprintf(`<w:r><w:t>%s</w:t></w:r>`, i.Content)
				case models.Bold:
					if len(i.Content) > 0 {
						fmt.Printf("  粗体内容: %v\n", i.Content[0])
						if textContent, ok := i.Content[0].(models.Text); ok {
							xml += `<w:r><w:rPr><w:b/></w:rPr><w:t>` + textContent.Content + `</w:t></w:r>`
						}
					}
				case models.Math:
					fmt.Printf("  数学公式(LaTeX): %s\n", i.LaTeX)
					mathXml := latex.ToOMML(i.LaTeX)
					fmt.Printf("  生成的数学XML: %s\n", mathXml)
					xml += `<m:oMathPara><m:oMath>` + mathXml + `</m:oMath></m:oMathPara>`
				}
			}
			xml += `</w:p>`
		case models.Math:
			// 处理块级数学公式
			fmt.Printf("块级数学公式(LaTeX): %s\n", b.LaTeX)
			mathXml := latex.ToOMML(b.LaTeX)
			fmt.Printf("生成的块级数学XML: %s\n", mathXml)
			xml += `<w:p><w:pPr><w:jc w:val="center"/></w:pPr><m:oMathPara><m:oMath>` + mathXml + `</m:oMath></m:oMathPara></w:p>`
		}
	}
	xml += `</w:body></w:document>`
	return xml
}

// CreateDOCX 创建DOCX文件
func CreateDOCX(doc models.Document, filename string) error {
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

// addFileToZip 向zip文件添加内容
func addFileToZip(w *zip.Writer, name, content string) error {
	fw, err := w.Create(name)
	if err != nil {
		return err
	}
	_, err = fw.Write([]byte(content))
	return err
}
