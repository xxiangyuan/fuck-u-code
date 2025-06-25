package metrics

import (
	"fmt"
	"go/ast"
	"go/token"
)

// ErrorHandlingMetric 评估代码错误处理情况
type ErrorHandlingMetric struct{}

// NewErrorHandlingMetric 创建错误处理指标
func NewErrorHandlingMetric() *ErrorHandlingMetric {
	return &ErrorHandlingMetric{}
}

// Name 返回指标名称
func (m *ErrorHandlingMetric) Name() string {
	return "错误处理"
}

// Description 返回指标描述
func (m *ErrorHandlingMetric) Description() string {
	return "评估代码中错误处理的完整性和合理性"
}

// Weight 返回指标权重
func (m *ErrorHandlingMetric) Weight() float64 {
	return 0.2
}

// Analyze 分析错误处理情况
func (m *ErrorHandlingMetric) Analyze(file *ast.File, fileSet *token.FileSet, content []byte) (float64, []string) {
	var issues []string

	// 记录各种错误处理情况
	var (
		totalErrorProducingCalls int // 产生错误的函数调用总数
		handledErrors            int // 被正确处理的错误数
		ignoredErrors            int // 被忽略的错误数
		blankIdentErrors         int // 使用_忽略的错误数
	)

	// 查找赋值语句中的错误处理
	ast.Inspect(file, func(n ast.Node) bool {
		assignStmt, ok := n.(*ast.AssignStmt)
		if !ok || len(assignStmt.Lhs) == 0 || len(assignStmt.Rhs) == 0 {
			return true
		}

		// 检查多返回值函数调用，最后一个返回值是error
		_, isCall := assignStmt.Rhs[0].(*ast.CallExpr)
		if isCall && len(assignStmt.Lhs) >= 2 {
			lastExpr := assignStmt.Lhs[len(assignStmt.Lhs)-1]

			if ident, ok := lastExpr.(*ast.Ident); ok {
				switch ident.Name {
				case "_":
					totalErrorProducingCalls++
					blankIdentErrors++
					pos := fileSet.Position(assignStmt.Pos())
					issues = append(issues, fmt.Sprintf("行 %d: 使用_忽略了错误返回值", pos.Line))
				case "err", "error":
					totalErrorProducingCalls++
					handledErrors++ // 暂时认为被赋值给变量就是处理了
				}
			}
		}
		return true
	})

	// 查找单独的函数调用表达式中可能忽略的错误
	ast.Inspect(file, func(n ast.Node) bool {
		exprStmt, ok := n.(*ast.ExprStmt)
		if !ok {
			return true
		}

		call, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			return true
		}

		// 检查是否可能是返回error却没有检查的函数
		if m.isPotentialErrorReturningFunc(call) {
			totalErrorProducingCalls++
			ignoredErrors++
			pos := fileSet.Position(exprStmt.Pos())
			issues = append(issues, fmt.Sprintf("行 %d: 忽略了可能的错误返回值", pos.Line))
		}
		return true
	})

	if totalErrorProducingCalls == 0 {
		return 0.0, issues
	}

	// 计算错误处理率和得分
	errorHandlingRate := float64(handledErrors) / float64(totalErrorProducingCalls)
	return m.calculateScore(errorHandlingRate), issues
}

// isPotentialErrorReturningFunc 检查函数调用是否可能返回error
func (m *ErrorHandlingMetric) isPotentialErrorReturningFunc(call *ast.CallExpr) bool {
	fun, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// 检查包名
	if x, ok := fun.X.(*ast.Ident); ok {
		// 一些常见的可能返回error的包
		if x.Name == "os" || x.Name == "ioutil" || x.Name == "io" ||
			x.Name == "net" || x.Name == "http" {
			return true
		}
	}

	// 检查常见的返回error的方法名
	methodName := fun.Sel.Name
	return methodName == "Open" || methodName == "Close" ||
		methodName == "Read" || methodName == "Write" ||
		methodName == "ReadFile" || methodName == "WriteFile" ||
		methodName == "Stat" || methodName == "Connect" ||
		methodName == "Dial"
}

// calculateScore 根据错误处理率计算得分
func (m *ErrorHandlingMetric) calculateScore(errorHandlingRate float64) float64 {
	switch {
	case errorHandlingRate >= 0.9:
		return 0.0 // 几乎所有错误都被处理，最佳情况
	case errorHandlingRate >= 0.7:
		return (1.0 - errorHandlingRate) * 0.5 // 0.7-0.9 -> 0.15-0.0
	case errorHandlingRate >= 0.5:
		return 0.15 + (0.7-errorHandlingRate)/0.2*0.35 // 0.5-0.7 -> 0.5-0.15
	case errorHandlingRate >= 0.3:
		return 0.5 + (0.5-errorHandlingRate)/0.2*0.3 // 0.3-0.5 -> 0.8-0.5
	default:
		return 0.8 + (0.3-errorHandlingRate)/0.3*0.2 // 0.0-0.3 -> 1.0-0.8
	}
}
