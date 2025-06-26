// Package metrics 提供代码质量分析指标
// 创建者：Done-0
package metrics

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/Done-0/fuck-u-code/pkg/i18n"
	"github.com/Done-0/fuck-u-code/pkg/parser"
)

// FunctionLengthMetric 检测函数长度及状态变量管理
type FunctionLengthMetric struct {
	*BaseMetric
	translator i18n.Translator
}

// NewFunctionLengthMetric 创建函数长度指标
func NewFunctionLengthMetric() *FunctionLengthMetric {
	translator := i18n.NewTranslator(i18n.ZhCN)
	return &FunctionLengthMetric{
		BaseMetric: NewBaseMetric(
			i18n.FormatKey("metric", "function_length"),
			"检测函数长度及状态变量管理，合理的函数长度和状态管理能提高代码可维护性",
			0.15,
			nil,
		),
		translator: translator,
	}
}

// SetTranslator 设置翻译器
func (m *FunctionLengthMetric) SetTranslator(translator i18n.Translator) {
	m.translator = translator
	m.name = translator.Translate(i18n.FormatKey("metric", "function_length"))
}

// Analyze 实现指标接口分析方法
func (m *FunctionLengthMetric) Analyze(parseResult parser.ParseResult) MetricResult {
	file, fileSet, content := ExtractGoAST(parseResult)

	// 如果content为空，使用解析结果获取
	if len(content) == 0 {
		content = []byte(strings.Repeat("\n", parseResult.GetTotalLines()))
	}

	score, issues := m.analyzeFunctions(file, fileSet, content, parseResult)

	return MetricResult{
		Score:       score,
		Issues:      issues,
		Description: m.Description(),
		Weight:      m.Weight(),
	}
}

// analyzeFunctions 分析函数长度及状态变量管理
func (m *FunctionLengthMetric) analyzeFunctions(file *ast.File, fileSet *token.FileSet, content []byte, parseResult parser.ParseResult) (float64, []string) {
	var issues []string

	// 使用解析器获取的函数信息
	functions := parseResult.GetFunctions()
	if len(functions) == 0 {
		return 0.0, issues
	}

	totalComplexity := 0
	longFunctions := 0
	totalFunctions := len(functions)

	// 分析每个函数
	for _, fn := range functions {
		lineCount := fn.EndLine - fn.StartLine + 1

		// 检查函数长度
		if lineCount > 50 {
			issues = append(issues, fmt.Sprintf("函数 '%s' 过长 (%d 行)，建议拆分", fn.Name, lineCount))
			longFunctions++
		} else if lineCount > 30 {
			issues = append(issues, fmt.Sprintf("函数 '%s' 较长 (%d 行)，可考虑重构", fn.Name, lineCount))
			longFunctions++
		}

		// 检查函数复杂度
		totalComplexity += fn.Complexity
		if fn.Complexity > 10 {
			issues = append(issues, fmt.Sprintf("函数 '%s' 复杂度过高 (%d)，建议简化", fn.Name, fn.Complexity))
		}

		// 检查参数数量
		if fn.Parameters > 5 {
			issues = append(issues, fmt.Sprintf("函数 '%s' 参数过多 (%d 个)，建议使用结构体封装", fn.Name, fn.Parameters))
		}
	}

	// 如果存在 Go AST，进行更深入的状态分析
	if file != nil {
		stateIssues, stateScore := m.analyzeStateManagement(file)
		issues = append(issues, stateIssues...)

		// 综合得分，函数长度占 60%，状态管理占 40%
		longFunctionRatio := float64(longFunctions) / float64(totalFunctions)
		avgComplexity := float64(totalComplexity) / float64(totalFunctions) / 10.0
		if avgComplexity > 1.0 {
			avgComplexity = 1.0
		}

		return longFunctionRatio*0.4 + avgComplexity*0.2 + stateScore*0.4, issues
	}

	// 对于非 Go 语言或无法进行 AST 分析的情况，使用简化评分
	longFunctionRatio := float64(longFunctions) / float64(totalFunctions)
	avgComplexity := float64(totalComplexity) / float64(totalFunctions) / 10.0
	if avgComplexity > 1.0 {
		avgComplexity = 1.0
	}

	return longFunctionRatio*0.6 + avgComplexity*0.4, issues
}

// analyzeStateManagement 分析状态变量管理
func (m *FunctionLengthMetric) analyzeStateManagement(file *ast.File) ([]string, float64) {
	var issues []string
	stateVars := make(map[string]stateVarInfo)
	globalVars := 0
	mutableVars := 0
	totalVars := 0

	// 检测全局变量
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						if name.Name != "_" {
							stateVars[name.Name] = stateVarInfo{
								isGlobal:  true,
								isMutable: true,
							}
							globalVars++
							mutableVars++
							totalVars++
							issues = append(issues, fmt.Sprintf("全局变量 '%s' 可能导致状态难以追踪", name.Name))
						}
					}
				}
			}
		}
	}

	// 检测函数内部的状态变量
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			// 检测函数参数中的指针类型
			if node.Type.Params != nil {
				for _, field := range node.Type.Params.List {
					if _, ok := field.Type.(*ast.StarExpr); ok {
						for _, name := range field.Names {
							issues = append(issues, fmt.Sprintf("函数 '%s' 的指针参数 '%s' 可能导致状态被外部修改", node.Name.Name, name.Name))
							mutableVars++
							totalVars++
						}
					}
				}
			}

			// 检测函数内部的状态变量修改
			if node.Body != nil {
				ast.Inspect(node.Body, func(n ast.Node) bool {
					if assign, ok := n.(*ast.AssignStmt); ok {
						for _, lhs := range assign.Lhs {
							if ident, ok := lhs.(*ast.Ident); ok {
								if info, exists := stateVars[ident.Name]; exists && !info.isMutable {
									issues = append(issues, fmt.Sprintf("在函数 '%s' 中修改了不应该变化的状态变量 '%s'", node.Name.Name, ident.Name))
								}
							}
						}
					}
					return true
				})
			}
		}
		return true
	})

	// 如果没有检测到变量
	if totalVars == 0 {
		return issues, 0.0
	}

	// 计算状态管理得分
	globalRatio := float64(globalVars) / float64(totalVars)
	mutableRatio := float64(mutableVars) / float64(totalVars)

	// 全局变量比例和可变状态比例越高，得分越差
	score := globalRatio*0.6 + mutableRatio*0.4
	if score > 1.0 {
		score = 1.0
	}

	return issues, score
}

// stateVarInfo 状态变量信息
type stateVarInfo struct {
	isGlobal  bool
	isMutable bool
}
