// Package analyzer 提供代码质量分析功能
package analyzer

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Done-0/fuck-u-code/pkg/metrics"
)

// AnalysisResult 表示分析结果
type AnalysisResult struct {
	TotalFiles       int
	TotalLines       int
	CodeQualityScore float64
	Metrics          map[string]MetricResult
	FilesAnalyzed    []FileAnalysisResult
}

// FileAnalysisResult 表示单文件分析结果
type FileAnalysisResult struct {
	FilePath     string
	FileScore    float64
	MetricScores map[string]float64
	Issues       []string
	Lines        int
}

// MetricResult 表示指标结果
type MetricResult struct {
	Name        string
	Score       float64
	Description string
	Weight      float64
}

// Analyzer 代码分析器
type Analyzer struct {
	metrics []metrics.Metric
}

// NewAnalyzer 创建新的分析器
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		metrics: []metrics.Metric{
			metrics.NewCyclomaticComplexityMetric(),
			metrics.NewFunctionLengthMetric(),
			metrics.NewCommentRatioMetric(),
			metrics.NewErrorHandlingMetric(),
			metrics.NewNamingConventionMetric(),
			metrics.NewCodeDuplicationMetric(),
		},
	}
}

// Analyze 分析代码目录
func (a *Analyzer) Analyze(path string) (*AnalysisResult, error) {
	result := &AnalysisResult{
		Metrics:       make(map[string]MetricResult),
		FilesAnalyzed: []FileAnalysisResult{},
	}

	err := filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 只分析Go文件
		if !d.IsDir() && strings.HasSuffix(filePath, ".go") {
			fileResult, err := a.analyzeFile(filePath)
			if err != nil {
				return err
			}

			result.FilesAnalyzed = append(result.FilesAnalyzed, *fileResult)
			result.TotalFiles++
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历项目失败: %w", err)
	}

	a.calculateOverallScore(result)
	return result, nil
}

// AnalyzeFile 分析单个文件
func (a *Analyzer) AnalyzeFile(filePath string) (*AnalysisResult, error) {
	result := &AnalysisResult{
		Metrics:       make(map[string]MetricResult),
		FilesAnalyzed: []FileAnalysisResult{},
	}

	fileResult, err := a.analyzeFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("分析文件失败: %w", err)
	}

	result.FilesAnalyzed = append(result.FilesAnalyzed, *fileResult)
	result.TotalFiles = 1

	a.calculateOverallScore(result)
	return result, nil
}

// analyzeFile 分析单个文件的内部实现
func (a *Analyzer) analyzeFile(filePath string) (*FileAnalysisResult, error) {
	fileSet := token.NewFileSet()
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	astFile, err := parser.ParseFile(fileSet, filePath, content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("解析文件失败: %w", err)
	}

	result := &FileAnalysisResult{
		FilePath:     filePath,
		MetricScores: make(map[string]float64),
		Issues:       []string{},
	}

	// 计算文件行数
	lines := strings.Split(string(content), "\n")
	result.Lines = len(lines)

	// 应用每个指标
	for _, metric := range a.metrics {
		score, issues := metric.Analyze(astFile, fileSet, content)
		result.MetricScores[metric.Name()] = score
		result.Issues = append(result.Issues, issues...)
	}

	result.FileScore = a.calculateFileScore(result.MetricScores)
	return result, nil
}

// calculateFileScore 计算文件总得分
func (a *Analyzer) calculateFileScore(metricScores map[string]float64) float64 {
	var totalScore float64
	var totalWeight float64

	for _, metric := range a.metrics {
		totalScore += metricScores[metric.Name()] * metric.Weight()
		totalWeight += metric.Weight()
	}

	if totalWeight == 0 {
		return 0
	}
	return totalScore / totalWeight
}

// calculateOverallScore 计算整体得分
func (a *Analyzer) calculateOverallScore(result *AnalysisResult) {
	metricScores := make(map[string]float64)
	metricCounts := make(map[string]int)

	// 汇总每个指标的总分和文件行数
	for _, file := range result.FilesAnalyzed {
		for name, score := range file.MetricScores {
			metricScores[name] += score
			metricCounts[name]++
		}
		result.TotalLines += file.Lines
	}

	// 计算每个指标的平均分
	totalScore := 0.0
	totalWeight := 0.0

	for _, metric := range a.metrics {
		name := metric.Name()
		count := metricCounts[name]

		if count > 0 {
			avgScore := metricScores[name] / float64(count)

			result.Metrics[name] = MetricResult{
				Name:        name,
				Score:       avgScore,
				Description: metric.Description(),
				Weight:      metric.Weight(),
			}

			totalScore += avgScore * metric.Weight()
			totalWeight += metric.Weight()
		}
	}

	// 最终分数 (0-100)，越高表示代码越"屎山"
	if totalWeight > 0 {
		result.CodeQualityScore = (totalScore / totalWeight) * 100
	}
}
