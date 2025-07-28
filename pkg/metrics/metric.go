// Package metrics 提供代码质量分析指标
// 创建者：Done-0
package metrics

import (
	"go/ast"
	"go/token"

	"github.com/Done-0/fuck-u-code/pkg/common"
	"github.com/Done-0/fuck-u-code/pkg/i18n"
	"github.com/Done-0/fuck-u-code/pkg/parser"
)

// MetricResult 表示指标分析结果
type MetricResult struct {
	Score       float64  // 得分 (0-1，越低越好)
	Issues      []string // 问题列表
	Description string   // 结果描述
	Weight      float64  // 权重 (0-1)
}

// Metric 代码质量指标接口
type Metric interface {
	// Name 返回指标名称
	Name() string

	// Description 返回指标描述
	Description() string

	// Weight 返回指标权重
	Weight() float64

	// Analyze 分析代码并返回结果
	Analyze(parseResult parser.ParseResult) MetricResult

	// SupportedLanguages 返回支持的语言类型
	SupportedLanguages() []common.LanguageType

	// SetTranslator 设置翻译器
	SetTranslator(translator i18n.Translator)
}

// AnalysisResult 表示分析结果
type AnalysisResult struct {
	FilePath      string                  // 文件路径
	TotalLines    int                     // 总行数
	CommentLines  int                     // 注释行数
	Imports       []string                // 导入的包
	MetricResults map[string]MetricResult // 各指标结果
	Functions     []parser.Function       // 函数列表
	Language      common.LanguageType     // 语言类型
	ParseResult   parser.ParseResult      // 解析结果
}

// GetOverallScore 获取总体评分
func (r *AnalysisResult) GetOverallScore() float64 {
	if len(r.MetricResults) == 0 {
		return 0.0
	}

	totalScore := 0.0
	totalWeight := 0.0

	for _, result := range r.MetricResults {
		metricScore := result.Score
		totalScore += metricScore * result.Weight
		totalWeight += result.Weight
	}

	if totalWeight == 0 {
		return 0.0
	}

	finalScore := totalScore / totalWeight

	// 确保分数在0-1范围内
	if finalScore > 1.0 {
		finalScore = 1.0
	} else if finalScore < 0.0 {
		finalScore = 0.0
	}

	return finalScore
}

// AddMetricResult 添加指标结果
func (r *AnalysisResult) AddMetricResult(name string, result MetricResult) {
	r.MetricResults[name] = result
}

// GetIssues 获取所有问题
func (r *AnalysisResult) GetIssues() []string {
	issues := make([]string, 0, len(r.MetricResults)*2)
	for _, result := range r.MetricResults {
		issues = append(issues, result.Issues...)
	}
	return issues
}

// BaseMetric 提供指标的基础实现
type BaseMetric struct {
	name               string
	description        string
	weight             float64
	supportedLanguages []common.LanguageType
	translator         i18n.Translator
}

// NewBaseMetric 创建基础指标
func NewBaseMetric(
	name string,
	description string,
	weight float64,
	supportedLanguages []common.LanguageType,
) *BaseMetric {
	return &BaseMetric{
		name:               name,
		description:        description,
		weight:             weight,
		supportedLanguages: supportedLanguages,
	}
}

// Name 返回指标名称
func (m *BaseMetric) Name() string {
	return m.name
}

// Description 返回指标描述
func (m *BaseMetric) Description() string {
	return m.description
}

// Weight 返回指标权重
func (m *BaseMetric) Weight() float64 {
	return m.weight
}

// SupportedLanguages 返回支持的语言类型
func (m *BaseMetric) SupportedLanguages() []common.LanguageType {
	return m.supportedLanguages
}

// SetTranslator 设置翻译器
func (m *BaseMetric) SetTranslator(translator i18n.Translator) {
	m.translator = translator
}

// NewAnalysisResult 创建新的分析结果
func NewAnalysisResult(filePath string, parseResult parser.ParseResult) *AnalysisResult {
	// 获取导入列表
	var imports []string
	if importer, ok := parseResult.(interface{ GetImportPaths() []string }); ok {
		imports = importer.GetImportPaths()
	}

	return &AnalysisResult{
		FilePath:      filePath,
		TotalLines:    parseResult.GetTotalLines(),
		CommentLines:  parseResult.GetCommentLines(),
		Imports:       imports,
		MetricResults: make(map[string]MetricResult),
		Functions:     parseResult.GetFunctions(),
		Language:      parseResult.GetLanguage(),
		ParseResult:   parseResult,
	}
}

// ExtractGoAST 从解析结果中提取Go语言的AST信息
func ExtractGoAST(parseResult parser.ParseResult) (*ast.File, *token.FileSet, []byte) {
	// 检查是否为Go语言
	if parseResult.GetLanguage() != common.Go {
		return nil, nil, nil
	}

	// 获取AST根节点
	astRoot := parseResult.GetASTRoot()
	if astRoot == nil {
		return nil, nil, nil
	}

	// 类型转换
	file, ok := astRoot.(*ast.File)
	if !ok {
		return nil, nil, nil
	}

	// 创建文件集
	fileSet := token.NewFileSet()

	// 获取源代码内容
	var content []byte
	if contentProvider, ok := parseResult.(interface{ GetContent() []byte }); ok {
		content = contentProvider.GetContent()
	}

	return file, fileSet, content
}
