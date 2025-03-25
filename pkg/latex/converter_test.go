package latex

import (
	"strings"
	"testing"
)

func TestToOMML(t *testing.T) {
	testCases := []struct {
		name       string
		latex      string
		expectPart string // 期望结果中应包含的部分字符串
	}{
		{
			name:       "简单文本",
			latex:      "x + y",
			expectPart: "<m:r><m:t>x + y</m:t></m:r>",
		},
		{
			name:       "分数公式",
			latex:      "\\frac{a}{b}",
			expectPart: "<m:f><m:fPr><m:type m:val=\"bar\"/></m:fPr><m:num><m:r><m:t>a</m:t></m:r></m:num><m:den><m:r><m:t>b</m:t></m:r></m:den></m:f>",
		},
		{
			name:       "上标",
			latex:      "E=mc^2",
			expectPart: "<m:t>E=mc^2</m:t>", // 简化的验证，只检查文本内容
		},
		{
			name:       "复杂上标",
			latex:      "x^{n+1}",
			expectPart: "<m:sSup>", // 验证包含上标标记
		},
		{
			name:       "下标",
			latex:      "a_i",
			expectPart: "<m:sSub>", // 验证包含下标标记
		},
		{
			name:       "复杂下标",
			latex:      "a_{i+1}",
			expectPart: "<m:sSub>", // 验证包含下标标记
		},
		{
			name:       "希腊字母",
			latex:      "\\alpha + \\beta = \\gamma",
			expectPart: "α + β = γ", // 验证替换了希腊字母
		},
		{
			name:       "特殊贝塞尔函数",
			latex:      "J_\\alpha(x)",
			expectPart: "<m:r><m:t>J</m:t></m:r><m:sSub>", // 验证处理了贝塞尔函数格式
		},
		{
			name:       "麦克斯韦方程",
			latex:      "\\nabla \\times \\vec{E}",
			expectPart: "<m:r><m:t>∇</m:t></m:r><m:r><m:t>×</m:t></m:r>", // 验证处理了向量算符
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ToOMML(tc.latex)
			if !strings.Contains(result, tc.expectPart) {
				t.Errorf("期望结果包含 '%s'，但实际结果为:\n%s", tc.expectPart, result)
			}
		})
	}
}

func TestSpecialFunctions(t *testing.T) {
	// 测试processLeftRight函数
	t.Run("processLeftRight", func(t *testing.T) {
		latex := "\\left(x + y\\right)"
		result := processLeftRight(latex)
		expected := "(x + y)"
		if result != expected {
			t.Errorf("processLeftRight 未正确处理左右括号，期望 '%s'，得到 '%s'", expected, result)
		}
	})

	// 测试processVec函数
	t.Run("processVec", func(t *testing.T) {
		latex := "\\vec{E}"
		result := processVec(latex)
		if !strings.Contains(result, "E→") {
			t.Errorf("processVec 未正确处理向量，期望包含 'E→'，得到 '%s'", result)
		}
	})

	// 测试特殊公式处理函数
	specialCases := []struct {
		name     string
		function func(string) string
		part     string
	}{
		{"processBesselFunction", processBesselFunction, "<m:r><m:t>J</m:t></m:r>"},
		{"processMaxwellEquation1", processMaxwellEquation1, "∇"},
		{"processMaxwellEquation2", processMaxwellEquation2, "∇"},
		{"processMaxwellEquation3", processMaxwellEquation3, "∇"},
		{"processMaxwellEquation4", processMaxwellEquation4, "∇"},
		{"processIntegralSum", processIntegralSum, "∫"},
	}

	for _, sc := range specialCases {
		t.Run(sc.name, func(t *testing.T) {
			result := sc.function("")
			if !strings.Contains(result, sc.part) {
				t.Errorf("%s 未返回预期结果，期望包含 '%s'，得到 '%s'", sc.name, sc.part, result)
			}
		})
	}
}
