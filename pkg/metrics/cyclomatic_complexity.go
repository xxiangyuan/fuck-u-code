package metrics

import (
	"fmt"
	"go/ast"
	"go/token"
)

// CyclomaticComplexityMetric 计算代码的循环复杂度
type CyclomaticComplexityMetric struct{}

// NewCyclomaticComplexityMetric 创建循环复杂度指标
func NewCyclomaticComplexityMetric() *CyclomaticComplexityMetric {
	return &CyclomaticComplexityMetric{}
}

// Name 返回指标名称
func (m *CyclomaticComplexityMetric) Name() string {
	return "循环复杂度"
}

// Description 返回指标描述
func (m *CyclomaticComplexityMetric) Description() string {
	return "测量函数的控制流复杂度，复杂度越高，代码越难理解和测试"
}

// Weight 返回指标权重
func (m *CyclomaticComplexityMetric) Weight() float64 {
	return 0.25
}

// Analyze 分析代码的循环复杂度
func (m *CyclomaticComplexityMetric) Analyze(file *ast.File, fileSet *token.FileSet, content []byte) (float64, []string) {
	var issues []string
	funcCount := 0
	totalComplexity := 0

	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		funcName := funcDecl.Name.Name
		complexity := m.calculateComplexity(funcDecl)
		funcCount++
		totalComplexity += complexity

		if complexity > 15 {
			issues = append(issues, fmt.Sprintf("函数 %s 的循环复杂度过高 (%d)，考虑重构", funcName, complexity))
		} else if complexity > 10 {
			issues = append(issues, fmt.Sprintf("函数 %s 的循环复杂度较高 (%d)，建议简化", funcName, complexity))
		}
		return true
	})

	if funcCount == 0 {
		return 0.0, issues
	}

	avgComplexity := float64(totalComplexity) / float64(funcCount)
	return m.calculateScore(avgComplexity), issues
}

// calculateComplexity 计算函数的复杂度
func (m *CyclomaticComplexityMetric) calculateComplexity(funcDecl *ast.FuncDecl) int {
	complexity := 1

	ast.Inspect(funcDecl, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.IfStmt:
			complexity++
		case *ast.ForStmt, *ast.RangeStmt:
			complexity++
		case *ast.CaseClause:
			complexity++
		case *ast.BinaryExpr:
			if node.Op == token.LAND || node.Op == token.LOR {
				complexity++
			}
		}
		return true
	})

	return complexity
}

// calculateScore 根据平均复杂度计算得分
func (m *CyclomaticComplexityMetric) calculateScore(avgComplexity float64) float64 {
	switch {
	case avgComplexity <= 5:
		return 0.0
	case avgComplexity <= 10:
		return (avgComplexity - 5) / 5 * 0.5 // 5-10 -> 0.0-0.5
	case avgComplexity <= 20:
		return 0.5 + (avgComplexity-10)/10*0.3 // 10-20 -> 0.5-0.8
	default:
		score := 0.8 + (avgComplexity-20)/10*0.2 // >20 -> 0.8-1.0
		if score > 1.0 {
			score = 1.0
		}
		return score
	}
}
