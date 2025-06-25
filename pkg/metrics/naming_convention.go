package metrics

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"
)

// NamingConventionMetric 评估代码命名规范
type NamingConventionMetric struct{}

// NewNamingConventionMetric 创建命名规范指标
func NewNamingConventionMetric() *NamingConventionMetric {
	return &NamingConventionMetric{}
}

// Name 返回指标名称
func (m *NamingConventionMetric) Name() string {
	return "命名规范"
}

// Description 返回指标描述
func (m *NamingConventionMetric) Description() string {
	return "评估代码是否遵循标准命名规范，如驼峰命名法、合适的包名等"
}

// Weight 返回指标权重
func (m *NamingConventionMetric) Weight() float64 {
	return 0.15
}

// 命名检查结果
type namingCheckResult struct {
	totalNames int      // 检查的标识符总数
	badNames   int      // 不符合规范的标识符数
	issues     []string // 发现的问题
}

// Analyze 分析命名规范
func (m *NamingConventionMetric) Analyze(file *ast.File, fileSet *token.FileSet, content []byte) (float64, []string) {
	result := &namingCheckResult{issues: []string{}}

	m.checkPackageName(file, result)
	m.checkFunctionNames(file, fileSet, result)
	m.checkVariableNames(file, fileSet, result)

	score := m.calculateScore(result)
	return score, result.issues
}

// checkPackageName 检查包名
func (m *NamingConventionMetric) checkPackageName(file *ast.File, result *namingCheckResult) {
	if !m.isValidPackageName(file.Name.Name) {
		result.issues = append(result.issues, fmt.Sprintf("包名 '%s' 不符合规范（应使用小写字母，避免下划线）", file.Name.Name))
	}
	result.totalNames++
}

// checkFunctionNames 检查函数、方法、参数的命名
func (m *NamingConventionMetric) checkFunctionNames(file *ast.File, fileSet *token.FileSet, result *namingCheckResult) {
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		m.checkFuncName(funcDecl, fileSet, result)
		m.checkReceiverName(funcDecl, fileSet, result)
		m.checkParamNames(funcDecl, fileSet, result)
	}
}

// checkFuncName 检查函数名
func (m *NamingConventionMetric) checkFuncName(funcDecl *ast.FuncDecl, fileSet *token.FileSet, result *namingCheckResult) {
	funcName := funcDecl.Name.Name
	if len(funcName) > 0 {
		result.totalNames++

		if !m.isValidFunctionName(funcName) {
			result.badNames++
			pos := fileSet.Position(funcDecl.Pos())
			result.issues = append(result.issues, fmt.Sprintf("行 %d: 函数名 '%s' 不符合驼峰命名规范", pos.Line, funcName))
		}
	}
}

// checkReceiverName 检查方法接收器名
func (m *NamingConventionMetric) checkReceiverName(funcDecl *ast.FuncDecl, fileSet *token.FileSet, result *namingCheckResult) {
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		for _, field := range funcDecl.Recv.List {
			if len(field.Names) > 0 {
				recvName := field.Names[0].Name
				result.totalNames++

				if !m.isValidReceiverName(recvName) {
					result.badNames++
					pos := fileSet.Position(field.Pos())
					result.issues = append(result.issues, fmt.Sprintf("行 %d: 方法接收器名 '%s' 不符合规范（应为1-2个小写字母）", pos.Line, recvName))
				}
			}
		}
	}
}

// checkParamNames 检查函数参数名
func (m *NamingConventionMetric) checkParamNames(funcDecl *ast.FuncDecl, fileSet *token.FileSet, result *namingCheckResult) {
	if funcDecl.Type.Params != nil {
		for _, field := range funcDecl.Type.Params.List {
			for _, name := range field.Names {
				paramName := name.Name
				if len(paramName) > 0 {
					result.totalNames++

					if !m.isValidVariableName(paramName) {
						result.badNames++
						pos := fileSet.Position(name.Pos())
						result.issues = append(result.issues, fmt.Sprintf("行 %d: 参数名 '%s' 不符合规范", pos.Line, paramName))
					}
				}
			}
		}
	}
}

// checkVariableNames 检查变量和常量名
func (m *NamingConventionMetric) checkVariableNames(file *ast.File, fileSet *token.FileSet, result *namingCheckResult) {
	ast.Inspect(file, func(n ast.Node) bool {
		valueSpec, ok := n.(*ast.ValueSpec)
		if !ok {
			return true
		}

		// 检查变量和常量声明
		for _, name := range valueSpec.Names {
			if name.Name != "_" { // 忽略匿名变量
				result.totalNames++

				if !m.isValidVariableName(name.Name) {
					result.badNames++
					pos := fileSet.Position(name.Pos())
					result.issues = append(result.issues, fmt.Sprintf("行 %d: 变量名 '%s' 不符合规范", pos.Line, name.Name))
				}
			}
		}
		return true
	})
}

// calculateScore 计算最终得分
func (m *NamingConventionMetric) calculateScore(result *namingCheckResult) float64 {
	if result.totalNames == 0 {
		return 0.0
	}

	badNameRatio := float64(result.badNames) / float64(result.totalNames)

	switch {
	case badNameRatio == 0:
		return 0.0 // 完全符合规范
	case badNameRatio <= 0.1:
		return badNameRatio // 0-0.1
	case badNameRatio <= 0.25:
		return 0.1 + (badNameRatio-0.1)/0.15*0.4 // 0.1-0.5
	case badNameRatio <= 0.5:
		return 0.5 + (badNameRatio-0.25)/0.25*0.3 // 0.5-0.8
	default:
		return 0.8 + (badNameRatio-0.5)/0.5*0.2 // 0.8-1.0
	}
}

// isValidPackageName 检查包名是否符合规范
func (m *NamingConventionMetric) isValidPackageName(name string) bool {
	if len(name) == 0 {
		return false
	}

	// 包名应该全部小写，没有下划线
	return !strings.Contains(name, "_") && strings.ToLower(name) == name
}

// isValidFunctionName 检查函数名是否符合规范
func (m *NamingConventionMetric) isValidFunctionName(name string) bool {
	if len(name) == 0 {
		return false
	}

	// 测试函数例外 (如 TestXxx)
	if strings.HasPrefix(name, "Test") || strings.HasPrefix(name, "Benchmark") {
		return true
	}

	// 普通函数应该使用驼峰命名法
	return m.isCamelCase(name) && !strings.Contains(name, "_")
}

// isValidVariableName 检查变量名是否符合规范
func (m *NamingConventionMetric) isValidVariableName(name string) bool {
	if len(name) == 0 {
		return false
	}

	// 变量名应该使用驼峰命名法，无下划线
	return m.isCamelCase(name) && !strings.Contains(name, "_")
}

// isValidReceiverName 检查方法接收器名是否符合规范
func (m *NamingConventionMetric) isValidReceiverName(name string) bool {
	// 接收器名通常为1-2个小写字母
	return len(name) <= 2 && strings.ToLower(name) == name
}

// isCamelCase 检查名称是否符合驼峰命名法
func (m *NamingConventionMetric) isCamelCase(name string) bool {
	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}
