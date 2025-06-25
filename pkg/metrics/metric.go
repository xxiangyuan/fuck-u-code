// Package metrics 提供代码质量评估的各种指标实现
package metrics

import (
	"go/ast"
	"go/token"
)

// Metric 定义代码质量指标接口
type Metric interface {
	// Name 返回指标名称
	Name() string

	// Description 返回指标描述
	Description() string

	// Weight 返回指标在总评分中的权重
	Weight() float64

	// Analyze 分析代码并返回得分和问题列表
	// 得分范围: 0.0-1.0，0表示最佳，1表示最差
	Analyze(file *ast.File, fileSet *token.FileSet, content []byte) (float64, []string)
}
