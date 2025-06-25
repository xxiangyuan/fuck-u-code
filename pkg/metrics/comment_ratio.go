package metrics

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// CommentRatioMetric 评估代码注释比例
type CommentRatioMetric struct{}

// NewCommentRatioMetric 创建注释比例指标
func NewCommentRatioMetric() *CommentRatioMetric {
	return &CommentRatioMetric{}
}

// Name 返回指标名称
func (m *CommentRatioMetric) Name() string {
	return "注释覆盖率"
}

// Description 返回指标描述
func (m *CommentRatioMetric) Description() string {
	return "评估代码注释的比例，注释太少不利于维护，太多则可能表明代码逻辑不清晰"
}

// Weight 返回指标权重
func (m *CommentRatioMetric) Weight() float64 {
	return 0.1
}

// Analyze 分析注释比例
func (m *CommentRatioMetric) Analyze(file *ast.File, fileSet *token.FileSet, content []byte) (float64, []string) {
	var issues []string

	lines := strings.Split(string(content), "\n")
	totalLines := len(lines)
	if totalLines == 0 {
		return 0.0, issues
	}

	commentMap := m.buildCommentMap(file, fileSet)
	commentLines := len(commentMap)
	commentRatio := float64(commentLines) / float64(totalLines)

	score := m.evaluateCommentRatio(commentRatio, &issues)
	return score, issues
}

// buildCommentMap 构建注释行映射
func (m *CommentRatioMetric) buildCommentMap(file *ast.File, fileSet *token.FileSet) map[int]bool {
	commentMap := make(map[int]bool)
	for _, cg := range file.Comments {
		for _, comment := range cg.List {
			pos := fileSet.Position(comment.Pos())
			commentLines := strings.Count(comment.Text, "\n") + 1
			for i := 0; i < commentLines; i++ {
				commentMap[pos.Line+i] = true
			}
		}
	}
	return commentMap
}

// evaluateCommentRatio 评估注释比例并生成问题列表
func (m *CommentRatioMetric) evaluateCommentRatio(commentRatio float64, issues *[]string) float64 {
	if commentRatio < 0.05 {
		*issues = append(*issues, fmt.Sprintf("代码注释率过低 (%.1f%%)，难以维护", commentRatio*100))
		return 0.9 // 严重问题
	} else if commentRatio < 0.1 {
		*issues = append(*issues, fmt.Sprintf("代码注释率较低 (%.1f%%)，建议增加注释", commentRatio*100))
		return 0.7 // 中度问题
	} else if commentRatio > 0.5 {
		*issues = append(*issues, fmt.Sprintf("注释率过高 (%.1f%%)，可能表明代码本身不够清晰", commentRatio*100))
		return 0.6 // 轻度问题
	} else if commentRatio >= 0.1 && commentRatio <= 0.3 {
		return 0.0 // 良好比例
	} else {
		// 0.3-0.5 的注释率，得分在 0.0-0.3 之间
		return (commentRatio - 0.3) / 0.2 * 0.3
	}
}
