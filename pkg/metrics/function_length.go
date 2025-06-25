package metrics

import (
	"fmt"
	"go/ast"
	"go/token"
)

// FunctionLengthMetric 检测函数长度
type FunctionLengthMetric struct{}

// NewFunctionLengthMetric 创建函数长度指标
func NewFunctionLengthMetric() *FunctionLengthMetric {
	return &FunctionLengthMetric{}
}

// Name 返回指标名称
func (m *FunctionLengthMetric) Name() string {
	return "函数长度"
}

// Description 返回指标描述
func (m *FunctionLengthMetric) Description() string {
	return "检测函数长度，函数过长通常意味着职责不单一，难以理解和测试"
}

// Weight 返回指标权重
func (m *FunctionLengthMetric) Weight() float64 {
	return 0.15
}

// Analyze 分析函数长度
func (m *FunctionLengthMetric) Analyze(file *ast.File, fileSet *token.FileSet, content []byte) (float64, []string) {
	var issues []string
	functionLengths := make(map[string]int)

	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			return true
		}

		startPos := fileSet.Position(funcDecl.Body.Pos())
		endPos := fileSet.Position(funcDecl.Body.End())
		length := endPos.Line - startPos.Line + 1
		functionName := funcDecl.Name.Name

		functionLengths[functionName] = length

		if length > 100 {
			issues = append(issues, fmt.Sprintf("函数 %s 过长 (%d 行)，极度建议拆分", functionName, length))
		} else if length > 50 {
			issues = append(issues, fmt.Sprintf("函数 %s 较长 (%d 行)，建议拆分为多个小函数", functionName, length))
		} else if length > 30 {
			issues = append(issues, fmt.Sprintf("函数 %s 长度为 %d 行，考虑是否可以简化", functionName, length))
		}
		return true
	})

	if len(functionLengths) == 0 {
		return 0.0, issues
	}

	avgLength := m.calculateAverageLength(functionLengths)
	return m.calculateScore(avgLength), issues
}

// calculateAverageLength 计算平均函数长度
func (m *FunctionLengthMetric) calculateAverageLength(functionLengths map[string]int) float64 {
	var totalLength int
	for _, length := range functionLengths {
		totalLength += length
	}
	return float64(totalLength) / float64(len(functionLengths))
}

// calculateScore 根据平均长度计算得分
func (m *FunctionLengthMetric) calculateScore(avgLength float64) float64 {
	switch {
	case avgLength <= 15:
		return 0.0
	case avgLength <= 30:
		return (avgLength - 15) / 15 * 0.4 // 15-30 -> 0.0-0.4
	case avgLength <= 50:
		return 0.4 + (avgLength-30)/20*0.3 // 30-50 -> 0.4-0.7
	case avgLength <= 100:
		return 0.7 + (avgLength-50)/50*0.2 // 50-100 -> 0.7-0.9
	default:
		score := 0.9 + (avgLength-100)/100*0.1 // >100 -> 0.9-1.0
		if score > 1.0 {
			score = 1.0
		}
		return score
	}
}
