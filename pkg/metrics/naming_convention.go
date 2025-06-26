package metrics

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"github.com/Done-0/fuck-u-code/pkg/common"
)

// NamingConventionMetric 检测命名规范
type NamingConventionMetric struct {
	*BaseMetric
}

// NewNamingConventionMetric 创建命名规范指标
func NewNamingConventionMetric() *NamingConventionMetric {
	return &NamingConventionMetric{
		BaseMetric: NewBaseMetric(
			"命名规范",
			"检测代码中的命名规范，良好的命名能提高代码可读性",
			0.1,
			[]common.LanguageType{common.Go},
		),
	}
}

// Analyze 分析命名规范
func (m *NamingConventionMetric) Analyze(file *ast.File, fileSet *token.FileSet, content []byte) (float64, []string) {
	var issues []string

	// 统计各种命名问题
	badNames := 0
	totalNames := 0

	// 分析包名
	if !m.isValidPackageName(file.Name.Name) {
		issues = append(issues, fmt.Sprintf("包名 '%s' 不符合规范，应使用小写字母且不包含下划线", file.Name.Name))
		badNames++
	}
	totalNames++

	// 分析变量、常量、函数和类型名称
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			totalNames++
			if !m.isValidFuncName(node.Name.Name) {
				issues = append(issues, fmt.Sprintf("函数名 '%s' 不符合规范", node.Name.Name))
				badNames++
			}

		case *ast.TypeSpec:
			totalNames++
			if !m.isValidTypeName(node.Name.Name) {
				issues = append(issues, fmt.Sprintf("类型名 '%s' 不符合规范", node.Name.Name))
				badNames++
			}

		case *ast.ValueSpec:
			for _, name := range node.Names {
				totalNames++
				if name.Name == "_" {
					continue
				}

				// 检查是否是常量
				if m.isConstDecl(n) {
					if !m.isValidConstName(name.Name) {
						issues = append(issues, fmt.Sprintf("常量名 '%s' 不符合规范", name.Name))
						badNames++
					}
				} else {
					// 变量名检查
					if !m.isValidVarName(name.Name) {
						issues = append(issues, fmt.Sprintf("变量名 '%s' 不符合规范", name.Name))
						badNames++
					}
				}
			}

		case *ast.AssignStmt:
			// 检查短变量声明
			if node.Tok == token.DEFINE {
				for _, expr := range node.Lhs {
					if ident, ok := expr.(*ast.Ident); ok {
						if ident.Name != "_" {
							totalNames++
							if !m.isValidVarName(ident.Name) {
								issues = append(issues, fmt.Sprintf("变量名 '%s' 不符合规范", ident.Name))
								badNames++
							}
						}
					}
				}
			}
		}
		return true
	})

	// 如果没有名称，返回0分
	if totalNames == 0 {
		return 0.0, issues
	}

	// 计算命名规范得分
	badRatio := float64(badNames) / float64(totalNames)
	return m.calculateScore(badRatio), issues
}

// isConstDecl 检查是否是常量声明
func (m *NamingConventionMetric) isConstDecl(node ast.Node) bool {
	var isConst bool
	ast.Inspect(node, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok {
			isConst = genDecl.Tok == token.CONST
			return false
		}
		return true
	})
	return isConst
}

// isValidPackageName 检查包名是否符合规范
func (m *NamingConventionMetric) isValidPackageName(name string) bool {
	// 包名应该是小写字母，不含下划线
	return m.isLowerCase(name) && !strings.Contains(name, "_")
}

// isValidVarName 检查变量名是否符合规范
func (m *NamingConventionMetric) isValidVarName(name string) bool {
	// 检查是否是有效的骆驼命名法
	return m.isCamelCase(name) || m.isLowerCase(name)
}

// isValidFuncName 检查函数名是否符合规范
func (m *NamingConventionMetric) isValidFuncName(name string) bool {
	// 函数名应该是驼峰命名法
	return m.isCamelCase(name)
}

// isValidTypeName 检查类型名是否符合规范
func (m *NamingConventionMetric) isValidTypeName(name string) bool {
	// 类型名应该是大写开头的驼峰命名法
	return m.isPascalCase(name)
}

// isValidConstName 检查常量名是否符合规范
func (m *NamingConventionMetric) isValidConstName(name string) bool {
	// 常量名应该是全大写加下划线，或者驼峰命名法
	return m.isUpperSnakeCase(name) || m.isPascalCase(name)
}

// isLowerCase 检查是否是小写字母
func (m *NamingConventionMetric) isLowerCase(name string) bool {
	for _, r := range name {
		if !unicode.IsLower(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// isCamelCase 检查是否是驼峰命名法（小写开头）
func (m *NamingConventionMetric) isCamelCase(name string) bool {
	if name == "" || strings.Contains(name, "_") || unicode.IsUpper([]rune(name)[0]) {
		return false
	}

	// 检查是否只包含字母和数字
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// isPascalCase 检查是否是帕斯卡命名法（大写开头）
func (m *NamingConventionMetric) isPascalCase(name string) bool {
	if name == "" || strings.Contains(name, "_") || !unicode.IsUpper([]rune(name)[0]) {
		return false
	}

	// 检查是否只包含字母和数字
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// isUpperSnakeCase 检查是否是大写蛇形命名法
func (m *NamingConventionMetric) isUpperSnakeCase(name string) bool {
	for _, r := range name {
		if !unicode.IsUpper(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

// calculateScore 根据不良命名比例计算得分
func (m *NamingConventionMetric) calculateScore(badRatio float64) float64 {
	switch {
	case badRatio == 0:
		return 0.0 // 完全符合命名规范
	case badRatio <= 0.1:
		return badRatio * 3 // 0.0-0.3
	case badRatio <= 0.3:
		return 0.3 + (badRatio-0.1)/0.2*0.4 // 0.3-0.7
	case badRatio <= 0.5:
		return 0.7 + (badRatio-0.3)/0.2*0.2 // 0.7-0.9
	default:
		return 0.9 + (badRatio-0.5)/0.5*0.1 // 0.9-1.0
	}
}
