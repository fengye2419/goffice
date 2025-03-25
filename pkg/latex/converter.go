package latex

import (
	"fmt"
	"strings"
)

// ToOMML 将LaTeX公式转换为Office Math Markup Language (OMML)
func ToOMML(latex string) string {
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
	} else {
		// 其他简单公式
		result = `<m:r><m:t>` + latex + `</m:t></m:r>`
	}

	return result
}
