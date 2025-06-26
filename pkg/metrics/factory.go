// Package metrics 提供代码质量分析指标
// 创建者：Done-0
package metrics

import (
	"github.com/Done-0/fuck-u-code/pkg/common"
	"github.com/Done-0/fuck-u-code/pkg/i18n"
	"github.com/Done-0/fuck-u-code/pkg/parser"
)

// MetricFactory 指标工厂结构体
type MetricFactory struct {
	translator i18n.Translator
}

// NewMetricFactory 创建指标工厂
func NewMetricFactory(translator i18n.Translator) *MetricFactory {
	return &MetricFactory{
		translator: translator,
	}
}

// SetTranslator 设置翻译器
func (f *MetricFactory) SetTranslator(translator i18n.Translator) {
	f.translator = translator
}

// CreateAllMetrics 创建所有指标
func (f *MetricFactory) CreateAllMetrics() []Metric {
	return []Metric{
		f.CreateCyclomaticComplexity(),
		f.CreateFunctionLength(),
		f.CreateCommentRatio(),
		f.CreateErrorHandling(),
		f.CreateNamingConvention(),
		f.CreateCodeDuplication(),
		f.CreateStructureAnalysis(),
	}
}

// CreateCyclomaticComplexity 创建循环复杂度指标
func (f *MetricFactory) CreateCyclomaticComplexity() Metric {
	metric := NewCyclomaticComplexityMetric()
	if f.translator != nil {
		metric.SetTranslator(f.translator)
	}
	return metric
}

// CreateFunctionLength 创建函数长度指标
func (f *MetricFactory) CreateFunctionLength() Metric {
	metric := NewFunctionLengthMetric()
	if f.translator != nil {
		metric.SetTranslator(f.translator)
	}
	return metric
}

// CreateCommentRatio 创建注释覆盖率指标
func (f *MetricFactory) CreateCommentRatio() Metric {
	metric := NewCommentRatioMetric()
	if f.translator != nil {
		metric.SetTranslator(f.translator)
	}
	return metric
}

// CreateErrorHandling 创建错误处理指标
func (f *MetricFactory) CreateErrorHandling() Metric {
	return f.createSimpleMetric("error_handling", 0.15)
}

// CreateNamingConvention 创建命名规范指标
func (f *MetricFactory) CreateNamingConvention() Metric {
	return f.createSimpleMetric("naming_convention", 0.1)
}

// CreateCodeDuplication 创建代码重复度指标
func (f *MetricFactory) CreateCodeDuplication() Metric {
	return f.createSimpleMetric("code_duplication", 0.15)
}

// CreateStructureAnalysis 创建代码结构分析指标
func (f *MetricFactory) CreateStructureAnalysis() Metric {
	return f.createSimpleMetric("structure_analysis", 0.2)
}

// createSimpleMetric 创建简单指标通用方法
func (f *MetricFactory) createSimpleMetric(metricKey string, weight float64) Metric {
	metric := &SimpleMetric{
		baseMetric: NewBaseMetric(
			f.translator.Translate(i18n.FormatKey("metric", metricKey)),
			f.translator.Translate("metric."+metricKey+".description"),
			weight,
			nil,
		),
		translator: f.translator,
	}

	return metric
}

// SimpleMetric 简单指标实现
type SimpleMetric struct {
	baseMetric *BaseMetric
	translator i18n.Translator
}

// Name 返回指标名称
func (m *SimpleMetric) Name() string {
	return m.baseMetric.Name()
}

// Description 返回指标描述
func (m *SimpleMetric) Description() string {
	return m.baseMetric.Description()
}

// Weight 返回指标权重
func (m *SimpleMetric) Weight() float64 {
	return m.baseMetric.Weight()
}

// SupportedLanguages 返回支持的语言类型
func (m *SimpleMetric) SupportedLanguages() []common.LanguageType {
	return m.baseMetric.SupportedLanguages()
}

// SetTranslator 设置翻译器
func (m *SimpleMetric) SetTranslator(translator i18n.Translator) {
	m.translator = translator
	m.baseMetric.name = translator.Translate(i18n.FormatKey("metric", m.Name()))
}

// Analyze 分析代码并返回结果
func (m *SimpleMetric) Analyze(parseResult parser.ParseResult) MetricResult {
	// 返回基本分析结果
	return MetricResult{
		Score:       0.5,
		Issues:      []string{},
		Description: m.Description(),
		Weight:      m.Weight(),
	}
}
