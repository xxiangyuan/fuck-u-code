package metrics

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// CodeDuplicationMetric 评估代码重复度
type CodeDuplicationMetric struct{}

// NewCodeDuplicationMetric 创建代码重复度指标
func NewCodeDuplicationMetric() *CodeDuplicationMetric {
	return &CodeDuplicationMetric{}
}

// Name 返回指标名称
func (m *CodeDuplicationMetric) Name() string {
	return "代码重复度"
}

// Description 返回指标描述
func (m *CodeDuplicationMetric) Description() string {
	return "评估代码中重复逻辑的比例，重复代码越多，越需要抽象和重构"
}

// Weight 返回指标权重
func (m *CodeDuplicationMetric) Weight() float64 {
	return 0.15
}

// Analyze 分析代码重复度
func (m *CodeDuplicationMetric) Analyze(file *ast.File, fileSet *token.FileSet, content []byte) (float64, []string) {
	var issues []string

	// 寻找类似的函数实现
	similarFunctions := make(map[string][]string) // 代码特征 -> 函数列表
	totalFunctions := 0

	// 遍历所有函数声明
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}

		totalFunctions++
		// 提取函数特征
		signature := m.extractFunctionSignature(funcDecl)
		similarFunctions[signature] = append(similarFunctions[signature], funcDecl.Name.Name)
	}

	// 识别可能重复的函数
	duplicatedCount := m.findDuplicatedFunctions(similarFunctions, &issues)

	// 如果没有足够的函数，返回低分
	if totalFunctions < 3 {
		return 0.1, issues
	}

	// 计算重复率和得分
	duplicationRate := float64(duplicatedCount) / float64(totalFunctions)
	return m.calculateScore(duplicationRate), issues
}

// findDuplicatedFunctions 查找重复函数并生成问题列表
func (m *CodeDuplicationMetric) findDuplicatedFunctions(similarFunctions map[string][]string, issues *[]string) int {
	duplicatedCount := 0

	for signature, functions := range similarFunctions {
		if len(functions) > 1 && len(signature) > 10 { // 忽略太短的函数
			duplicatedCount += len(functions) - 1
			*issues = append(*issues, fmt.Sprintf("可能存在重复实现: %s", strings.Join(functions, ", ")))
		}
	}

	return duplicatedCount
}

// calculateScore 根据重复率计算得分
func (m *CodeDuplicationMetric) calculateScore(duplicationRate float64) float64 {
	switch {
	case duplicationRate <= 0.05:
		return 0.0 // 几乎没有重复
	case duplicationRate <= 0.1:
		return duplicationRate * 5 // 0.05-0.1 -> 0.25-0.5
	case duplicationRate <= 0.2:
		return 0.5 + (duplicationRate-0.1)/0.1*0.3 // 0.1-0.2 -> 0.5-0.8
	default:
		score := 0.8 + (duplicationRate-0.2)/0.8*0.2 // >0.2 -> 0.8-1.0
		if score > 1.0 {
			score = 1.0
		}
		return score
	}
}

// extractFunctionSignature 提取函数的特征签名
func (m *CodeDuplicationMetric) extractFunctionSignature(funcDecl *ast.FuncDecl) string {
	var signature strings.Builder

	// 记录函数中关键节点类型，忽略具体标识符
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt:
			signature.WriteString("if")
		case *ast.ForStmt:
			signature.WriteString("for")
		case *ast.RangeStmt:
			signature.WriteString("range")
		case *ast.SwitchStmt:
			signature.WriteString("switch")
		case *ast.CaseClause:
			signature.WriteString("case")
		case *ast.AssignStmt:
			signature.WriteString("=")
		case *ast.ReturnStmt:
			signature.WriteString("return")
		}
		return true
	})

	return signature.String()
}
