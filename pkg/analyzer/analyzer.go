// Package analyzer 提供代码分析功能
// 创建者：Done-0
// 创建时间：2023-10-01
package analyzer

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Done-0/fuck-u-code/pkg/common"
	"github.com/Done-0/fuck-u-code/pkg/i18n"
	"github.com/Done-0/fuck-u-code/pkg/metrics"
	"github.com/Done-0/fuck-u-code/pkg/parser"
)

// Analyzer 代码分析器接口
type Analyzer interface {
	// Analyze 分析指定路径的代码
	Analyze(path string) (*AnalysisResult, error)

	// AnalyzeFile 分析单个文件
	AnalyzeFile(filePath string) (*AnalysisResult, error)

	// AnalyzeWithExcludes 使用指定的包含/排除模式分析目录
	AnalyzeWithExcludes(path string, includePatterns []string, excludePatterns []string) (*AnalysisResult, error)

	// SetLanguage 设置分析器使用的语言
	SetLanguage(lang i18n.Language)

	// SetSilent 设置静默模式
	SetSilent(silent bool)
}

// AnalysisResult 分析结果
type AnalysisResult struct {
	CodeQualityScore float64                 // 代码质量得分
	Metrics          map[string]MetricResult // 各项指标结果
	FilesAnalyzed    []FileAnalysisResult    // 分析的文件结果
	TotalFiles       int                     // 总文件数
	TotalLines       int                     // 总代码行数
}

// MetricResult 指标结果
type MetricResult struct {
	Name        string  // 指标名称
	Score       float64 // 得分(0-1，越高越差)
	Description string  // 描述
	Weight      float64 // 权重
}

// FileAnalysisResult 文件分析结果
type FileAnalysisResult struct {
	FilePath  string   // 文件路径
	FileScore float64  // 文件得分
	Issues    []string // 问题列表
}

// DefaultAnalyzer 默认分析器实现
type DefaultAnalyzer struct {
	codeAnalyzer *CodeAnalyzer
	translator   i18n.Translator
	silent       bool // 静默模式，不输出进度信息
}

// NewAnalyzer 创建新的代码分析器
func NewAnalyzer() Analyzer {
	translator := i18n.NewTranslator(i18n.ZhCN)
	return &DefaultAnalyzer{
		codeAnalyzer: NewCodeAnalyzer(translator),
		translator:   translator,
	}
}

// SetLanguage 设置分析器使用的语言
func (a *DefaultAnalyzer) SetLanguage(lang i18n.Language) {
	a.translator = i18n.NewTranslator(lang)
	a.codeAnalyzer.SetTranslator(a.translator)
}

// SetSilent 设置静默模式
func (a *DefaultAnalyzer) SetSilent(silent bool) {
	a.silent = silent
}

// Analyze 分析指定路径的代码
func (a *DefaultAnalyzer) Analyze(path string) (*AnalysisResult, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf(a.translator.Translate("error.path_not_accessible"), err)
	}

	if info.IsDir() {
		return a.AnalyzeWithExcludes(path, nil, []string{"*/vendor/*", "*/node_modules/*", "*/.git/*"})
	}

	return a.AnalyzeFile(path)
}

// AnalyzeFile 分析单个文件
func (a *DefaultAnalyzer) AnalyzeFile(filePath string) (*AnalysisResult, error) {
	// 使用内部的CodeAnalyzer分析文件
	fileResult, err := a.codeAnalyzer.AnalyzeFile(filePath)
	if err != nil {
		return nil, err
	}

	// 转换为AnalysisResult
	result := &AnalysisResult{
		CodeQualityScore: fileResult.GetOverallScore(),
		Metrics:          make(map[string]MetricResult),
		FilesAnalyzed:    make([]FileAnalysisResult, 0, 1),
		TotalFiles:       1,
		TotalLines:       fileResult.TotalLines,
	}

	// 添加指标结果
	for name, metricResult := range fileResult.MetricResults {
		result.Metrics[name] = MetricResult{
			Name:        name,
			Score:       metricResult.Score,
			Description: metricResult.Description,
			Weight:      metricResult.Weight,
		}
	}

	// 添加文件分析结果
	result.FilesAnalyzed = append(result.FilesAnalyzed, FileAnalysisResult{
		FilePath:  filePath,
		FileScore: fileResult.GetOverallScore(),
		Issues:    fileResult.GetIssues(),
	})

	return result, nil
}

// AnalyzeWithExcludes 使用指定的排除模式分析目录
func (a *DefaultAnalyzer) AnalyzeWithExcludes(path string, includePatterns []string, excludePatterns []string) (*AnalysisResult, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf(a.translator.Translate("error.path_not_accessible"), err)
	}

	if !info.IsDir() {
		return a.AnalyzeFile(path)
	}

	// 只在非静默模式下显示文件搜索进度
	var lastFoundCount int
	var progressCallback func(int)

	if !a.silent {
		fmt.Printf("🔍 %s\n", a.translator.Translate("analyzer.searching_files"))
		progressCallback = func(found int) {
			if found > lastFoundCount {
				fmt.Printf("\r📂 %s: %d", a.translator.Translate("analyzer.files_found"), found)
				lastFoundCount = found
			}
		}
	} else {
		// 静默模式下的空回调
		progressCallback = func(int) {}
	}

	// 分析目录中的所有文件
	fileResults, err := a.codeAnalyzer.AnalyzeDirectory(path, includePatterns, excludePatterns, progressCallback)
	if err != nil {
		return nil, err
	}

	// 只在非静默模式下清除进度显示
	if !a.silent && lastFoundCount > 0 {
		fmt.Printf("\r%s\r", strings.Repeat(" ", 50))
	}

	// 创建结果对象
	result := &AnalysisResult{
		Metrics:       make(map[string]MetricResult),
		FilesAnalyzed: make([]FileAnalysisResult, 0, len(fileResults)),
		TotalFiles:    len(fileResults),
	}

	// 收集所有指标结果
	allMetrics := make(map[string][]metrics.MetricResult)
	totalLines := 0

	// 处理每个文件的结果
	for _, fileResult := range fileResults {
		totalLines += fileResult.TotalLines

		// 添加文件分析结果
		result.FilesAnalyzed = append(result.FilesAnalyzed, FileAnalysisResult{
			FilePath:  fileResult.FilePath,
			FileScore: fileResult.GetOverallScore(),
			Issues:    fileResult.GetIssues(),
		})

		// 收集各指标结果
		for name, metricResult := range fileResult.MetricResults {
			if _, ok := allMetrics[name]; !ok {
				allMetrics[name] = make([]metrics.MetricResult, 0, len(fileResults))
			}
			allMetrics[name] = append(allMetrics[name], metricResult)
		}
	}

	// 计算各指标的平均分数
	for name, metricResults := range allMetrics {
		if len(metricResults) == 0 {
			continue
		}

		// 计算平均分
		totalScore := 0.0
		totalWeight := 0.0
		description := ""

		for _, m := range metricResults {
			totalScore += m.Score
			totalWeight = m.Weight
			description = m.Description
		}

		avgScore := totalScore / float64(len(metricResults))

		// 添加到结果中
		result.Metrics[name] = MetricResult{
			Name:        name,
			Score:       avgScore,
			Description: description,
			Weight:      totalWeight,
		}
	}

	// 设置总行数
	result.TotalLines = totalLines

	// 计算总体评分
	result.CodeQualityScore = a.codeAnalyzer.CalculateOverallScore(fileResults)

	return result, nil
}

// CodeAnalyzer 代码分析器
type CodeAnalyzer struct {
	metricFactory *metrics.MetricFactory
	translator    i18n.Translator
}

// NewCodeAnalyzer 创建新的代码分析器
func NewCodeAnalyzer(translator i18n.Translator) *CodeAnalyzer {
	metricFactory := metrics.NewMetricFactory(translator)
	return &CodeAnalyzer{
		metricFactory: metricFactory,
		translator:    translator,
	}
}

// SetTranslator 设置翻译器
func (a *CodeAnalyzer) SetTranslator(translator i18n.Translator) {
	a.translator = translator
	a.metricFactory.SetTranslator(translator)
}

// GetMetrics 获取所有指标
func (a *CodeAnalyzer) GetMetrics() []metrics.Metric {
	return a.metricFactory.CreateAllMetrics()
}

// AnalyzeFile 分析单个文件
func (a *CodeAnalyzer) AnalyzeFile(filePath string) (*metrics.AnalysisResult, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf(a.translator.Translate("error.file_read_failed"), filePath, err)
	}

	// 创建适合该文件的解析器
	codeParser := parser.CreateParserForFile(filePath)

	// 解析代码
	parseResult, err := codeParser.Parse(filePath, content)
	if err != nil {
		return nil, fmt.Errorf(a.translator.Translate("error.code_parse_failed"), filePath, err)
	}

	// 创建分析结果对象
	result := metrics.NewAnalysisResult(filePath, parseResult)

	// 应用每个指标进行分析
	for _, metric := range a.GetMetrics() {
		if !a.isLanguageSupported(metric, parseResult.GetLanguage()) {
			continue
		}

		metricResult := metric.Analyze(parseResult)
		result.AddMetricResult(metric.Name(), metricResult)
	}

	return result, nil
}

// isLanguageSupported 检查指标是否支持指定语言
func (a *CodeAnalyzer) isLanguageSupported(metric metrics.Metric, language common.LanguageType) bool {
	supportedLanguages := metric.SupportedLanguages()
	if len(supportedLanguages) == 0 {
		return true // 支持所有语言
	}

	for _, lang := range supportedLanguages {
		if lang == language {
			return true
		}
	}

	return false
}

// AnalyzeDirectory 分析目录
func (a *CodeAnalyzer) AnalyzeDirectory(dirPath string, includePatterns []string, excludePatterns []string, progressCallback func(found int)) ([]*metrics.AnalysisResult, error) {
	// 查找所有符合条件的文件
	files, err := common.FindSourceFiles(dirPath, includePatterns, excludePatterns, progressCallback)
	if err != nil {
		return nil, fmt.Errorf(a.translator.Translate("error.source_files_not_found"), err)
	}

	// 使用并发加速分析
	results := make([]*metrics.AnalysisResult, 0, len(files))
	resultsChan := make(chan *metrics.AnalysisResult, len(files))
	errChan := make(chan error, len(files))

	// 并发控制
	concurrencyLimit := min(8, len(files)) // 最大并发数
	semaphore := make(chan struct{}, concurrencyLimit)
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result, err := a.AnalyzeFile(filePath)
			if err != nil {
				errChan <- fmt.Errorf(a.translator.Translate("error.file_analysis_failed"), filePath, err)
				return
			}
			resultsChan <- result
		}(file)
	}

	// 等待所有分析完成
	wg.Wait()
	close(resultsChan)
	close(errChan)

	// 收集结果
	for result := range resultsChan {
		results = append(results, result)
	}

	// 检查错误
	for err := range errChan {
		fmt.Fprintf(os.Stderr, a.translator.Translate("warning.format"), err)
	}

	return results, nil
}

// CalculateOverallScore 计算总体评分
func (a *CodeAnalyzer) CalculateOverallScore(results []*metrics.AnalysisResult) float64 {
	if len(results) == 0 {
		return 0.0
	}

	totalScore := 0.0
	for _, result := range results {
		totalScore += result.GetOverallScore()
	}

	return totalScore / float64(len(results))
}

// min 返回两个整数中较小的一个
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
