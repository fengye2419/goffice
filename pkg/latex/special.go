package latex

import (
	"strings"
)

// processLeftRight 处理LaTeX中的\left(...\right)结构
func processLeftRight(latex string) string {
	result := latex
	for {
		leftIdx := strings.Index(result, "\\left")
		if leftIdx == -1 {
			break
		}

		// 找到匹配的\right
		rightIdx := strings.Index(result[leftIdx:], "\\right")
		if rightIdx == -1 {
			break
		}
		rightIdx += leftIdx

		// 替换为简单括号
		leftChar := string(result[leftIdx+5])
		rightChar := string(result[rightIdx+6])

		// 构建新字符串
		newStr := result[:leftIdx] + leftChar + result[leftIdx+6:rightIdx] + rightChar + result[rightIdx+7:]
		result = newStr
	}
	return result
}

// processVec 处理LaTeX中的向量表示\vec{}
func processVec(latex string) string {
	result := latex
	vecIdx := strings.Index(result, "\\vec{")
	for vecIdx != -1 {
		startIdx := vecIdx + 5 // 跳过\vec{

		// 找到匹配的右花括号
		level := 1
		endIdx := startIdx
		for endIdx < len(result) && level > 0 {
			if result[endIdx] == '{' {
				level++
			} else if result[endIdx] == '}' {
				level--
			}
			if level > 0 {
				endIdx++
			}
		}

		if endIdx < len(result) {
			vecContent := result[startIdx:endIdx]
			// 用箭头符号替换
			newStr := result[:vecIdx] + vecContent + "→" + result[endIdx+1:]
			result = newStr

			// 继续查找下一个\vec
			vecIdx = strings.Index(result, "\\vec{")
		} else {
			break
		}
	}
	return result
}

// processBesselFunction 处理贝塞尔函数特殊情况
func processBesselFunction(latex string) string {
	// 生成贝塞尔函数的OMML表示
	return `<m:r><m:t>J</m:t></m:r><m:sSub><m:e><m:r><m:t></m:t></m:r></m:e><m:sub><m:r><m:t>α</m:t></m:r></m:sub></m:sSub><m:r><m:t>(x)</m:t></m:r>`
}

// processMaxwellEquation1 处理麦克斯韦方程1
func processMaxwellEquation1(latex string) string {
	return `<m:r><m:t>∇</m:t></m:r><m:r><m:t>×</m:t></m:r><m:r><m:t>E</m:t></m:r><m:r><m:t>→</m:t></m:r><m:r><m:t> = -</m:t></m:r><m:f><m:fPr><m:type m:val="bar"/></m:fPr><m:num><m:r><m:t>∂B</m:t></m:r><m:r><m:t>→</m:t></m:r></m:num><m:den><m:r><m:t>∂t</m:t></m:r></m:den></m:f>`
}

// processMaxwellEquation2 处理麦克斯韦方程2
func processMaxwellEquation2(latex string) string {
	return `<m:r><m:t>∇</m:t></m:r><m:r><m:t>×</m:t></m:r><m:r><m:t>B</m:t></m:r><m:r><m:t>→</m:t></m:r><m:r><m:t> = μ</m:t></m:r><m:r><m:t>₀</m:t></m:r><m:r><m:t>J</m:t></m:r><m:r><m:t>→</m:t></m:r><m:r><m:t>+μ</m:t></m:r><m:r><m:t>₀</m:t></m:r><m:r><m:t>ε</m:t></m:r><m:r><m:t>₀</m:t></m:r><m:f><m:fPr><m:type m:val="bar"/></m:fPr><m:num><m:r><m:t>∂E</m:t></m:r><m:r><m:t>→</m:t></m:r></m:num><m:den><m:r><m:t>∂t</m:t></m:r></m:den></m:f>`
}

// processMaxwellEquation3 处理麦克斯韦方程3
func processMaxwellEquation3(latex string) string {
	return `<m:r><m:t>∇</m:t></m:r><m:r><m:t>·</m:t></m:r><m:r><m:t>E</m:t></m:r><m:r><m:t>→</m:t></m:r><m:r><m:t> = </m:t></m:r><m:f><m:fPr><m:type m:val="bar"/></m:fPr><m:num><m:r><m:t>ρ</m:t></m:r></m:num><m:den><m:r><m:t>ε</m:t></m:r><m:r><m:t>₀</m:t></m:r></m:den></m:f>`
}

// processMaxwellEquation4 处理麦克斯韦方程4
func processMaxwellEquation4(latex string) string {
	return `<m:r><m:t>∇</m:t></m:r><m:r><m:t>·</m:t></m:r><m:r><m:t>B</m:t></m:r><m:r><m:t>→</m:t></m:r><m:r><m:t> = 0</m:t></m:r>`
}

// processIntegralSum 处理积分与求和组合公式
func processIntegralSum(latex string) string {
	return `<m:r><m:t>∫</m:t></m:r><m:r><m:t> </m:t></m:r><m:r><m:t>∑</m:t></m:r><m:r><m:t> f(x) dx</m:t></m:r>`
}
