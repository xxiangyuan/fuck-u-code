// Package report 提供代码质量分析报告生成功能
// 创建者：Done-0
// 创建时间：2023-10-01
package report

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/fatih/color"

	"github.com/Done-0/fuck-u-code/pkg/analyzer"
	"github.com/Done-0/fuck-u-code/pkg/i18n"
)

// 颜色风格定义
var (
	titleStyle   = color.New(color.FgHiYellow, color.Bold)
	scoreStyle   = color.New(color.FgHiCyan, color.Bold)
	goodStyle    = color.New(color.FgHiGreen)
	warningStyle = color.New(color.FgHiYellow)
	dangerStyle  = color.New(color.FgHiRed)
	headerStyle  = color.New(color.FgMagenta, color.Bold)
	sectionStyle = color.New(color.FgHiMagenta, color.Bold)
	infoStyle    = color.New(color.FgBlue)
	successStyle = color.New(color.FgGreen, color.Bold)
	detailStyle  = color.New(color.FgCyan)
	metricStyle  = color.New(color.FgCyan)
	fileStyle    = color.New(color.FgMagenta)
	levelStyle   = color.New(color.FgCyan)
	numberStyle  = color.New(color.FgHiWhite)
)

// QualityLevels 定义代码质量等级（每10分一个段位）
var QualityLevels = []struct {
	MinScore    float64
	NameKey     string
	Description string
	Emoji       string
}{
	{0, "level.clean", "level.clean.description", "🌱"},
	{5, "level.mild", "level.mild.description", "🌸"},
	{15, "level.moderate", "level.moderate.description", "😐"},
	{25, "level.bad", "level.bad.description", "😷"},
	{40, "level.terrible", "level.terrible.description", "💩"},
	{55, "level.disaster", "level.disaster.description", "🤕"},
	{65, "level.disaster.severe", "level.disaster.severe.description", "☣️"},
	{75, "level.disaster.very_bad", "level.disaster.very_bad.description", "🧟"},
	{85, "level.disaster.extreme", "level.disaster.extreme.description", "☢️"},
	{95, "level.disaster.worst", "level.disaster.worst.description", "🪦"},
	{100, "level.disaster.ultimate", "level.disaster.ultimate.description", "👑💩"},
}

// Report 表示代码分析报告对象
type Report struct {
	result     *analyzer.AnalysisResult
	translator i18n.Translator
}

// NewReport 创建一个新的报告实例
func NewReport(result *analyzer.AnalysisResult) *Report {
	return &Report{
		result:     result,
		translator: i18n.NewTranslator(i18n.ZhCN), // 默认使用中文
	}
}

// SetTranslator 设置翻译器
func (r *Report) SetTranslator(translator i18n.Translator) {
	r.translator = translator
}

// ReportOptions 定义报告生成的选项
type ReportOptions struct {
	Verbose     bool // 是否显示详细报告
	TopFiles    int  // 显示最差文件的数量
	MaxIssues   int  // 每个文件显示的问题数量
	SummaryOnly bool // 是否只显示摘要
}

// DefaultReportOptions 默认报告选项
var DefaultReportOptions = &ReportOptions{
	Verbose:     false,
	TopFiles:    3,
	MaxIssues:   3,
	SummaryOnly: false,
}

// GenerateConsoleReport 生成控制台报告
func (r *Report) GenerateConsoleReport(options *ReportOptions) {
	if options == nil {
		options = DefaultReportOptions
	}

	score := r.result.CodeQualityScore
	level := r.getQualityLevel(score)

	// 打印标题和总体评分
	printDivider()
	titleStyle.Printf("\n  %s %s %s\n", level.Emoji, r.translator.Translate("report.title"), level.Emoji)
	printDivider()

	fmt.Printf("\n")
	scoreStyle.Printf("  %s", r.translator.Translate("report.overall_score", math.Round(score*10000)/100))
	fmt.Printf(" - ")
	r.printScoreComment(score)
	fmt.Printf("\n")

	// 打印质量等级
	detailStyle.Printf("  %s", r.translator.Translate("report.level", r.translator.Translate(level.NameKey)))
	detailStyle.Printf(" - %s\n\n", r.translator.Translate(level.Description))

	if !options.SummaryOnly {
		r.printMetricItems()

		if options.Verbose {
			r.printAllFiles(options)
		} else {
			r.printTopIssues(options)
		}
	}

	r.printSummary(level)

	if options.Verbose {
		r.printVerboseInfo()
	}

	printDivider()
	fmt.Println()
}

// printDivider 打印分隔线
func printDivider() {
	fmt.Printf("\n%s\n", strings.Repeat("─", 80))
}

// printMetricItems 打印各项评分指标及简评
func (r *Report) printMetricItems() {
	sectionStyle.Printf("\n◆ %s\n\n", r.translator.Translate("report.metrics_details"))

	metrics := r.getSortedMetrics()

	maxNameLen := 0
	for _, m := range metrics {
		if len(m.Name) > maxNameLen {
			maxNameLen = len(m.Name)
		}
	}

	nameFormat := fmt.Sprintf("  %%s %%-%ds", maxNameLen+2)
	scoreFormat := "%-8s"

	var totalWeight float64
	var weightedScore float64

	for _, m := range metrics {
		totalWeight += m.Weight
		weightedScore += m.Score * m.Weight
	}

	for _, m := range metrics {
		scorePercentage := math.Round(m.Score*10000) / 100

		// 确定状态图标和颜色
		var statusEmoji string
		var statusColor *color.Color

		switch {
		case scorePercentage < 20:
			statusEmoji = "✓✓" // 双勾 - 优秀
			statusColor = color.New(color.FgHiGreen, color.Bold)
		case scorePercentage < 35:
			statusEmoji = "✓" // 单勾 - 良好
			statusColor = goodStyle
		case scorePercentage < 50:
			statusEmoji = "○" // 圆圈 - 一般偏好
			statusColor = color.New(color.FgHiCyan)
		case scorePercentage < 60:
			statusEmoji = "•" // 点 - 一般
			statusColor = color.New(color.FgBlue)
		case scorePercentage < 70:
			statusEmoji = "⚠" // 警告 - 一般偏差
			statusColor = color.New(color.FgHiYellow)
		case scorePercentage < 80:
			statusEmoji = "!" // 感叹号 - 较差
			statusColor = warningStyle
		case scorePercentage < 90:
			statusEmoji = "!!" // 双感叹号 - 差
			statusColor = color.New(color.FgHiRed)
		default:
			statusEmoji = "✗" // 叉号 - 极差
			statusColor = dangerStyle
		}

		// 格式化分数
		scoreStr := fmt.Sprintf("%.2f%s", scorePercentage, r.translator.Translate("metric.score.suffix"))

		statusColor.Printf(nameFormat, statusEmoji, m.Name)
		metricStyle.Printf(scoreFormat, scoreStr)
		detailStyle.Printf("  %s\n", r.getMetricComment(m.Name, scorePercentage))
	}
	fmt.Println()

	if totalWeight > 0 {
		infoStyle.Printf("  %s", r.translator.Translate("report.score_calc"))

		// 计算公式的第一部分
		first := true
		infoStyle.Printf("(")
		for _, m := range metrics {
			if !first {
				infoStyle.Printf(" + ")
			}

			scorePercentage := math.Round(m.Score*10000) / 100
			infoStyle.Printf("%.2f×%.2f", scorePercentage, m.Weight)

			first = false
		}

		// 计算公式的第二部分
		overallScore := math.Round(weightedScore/totalWeight*10000) / 100
		infoStyle.Printf(") ÷ %.2f = %.2f\n\n", totalWeight, overallScore)
	}
}

// getMetricComment 返回指标评论
func (r *Report) getMetricComment(metricName string, score float64) string {
	var commentKey string

	var level string
	switch {
	case score < 20:
		level = "good"
	case score < 40:
		level = "medium"
	case score < 60:
		level = "medium"
	case score < 80:
		level = "bad"
	default:
		level = "bad"
	}

	nameKey := strings.ToLower(metricName)

	var metricType string
	switch {
	case strings.Contains(nameKey, "complexity") || strings.Contains(nameKey, "复杂度"):
		metricType = "complexity"
	case strings.Contains(nameKey, "function") || strings.Contains(nameKey, "length") || strings.Contains(nameKey, "状态") || strings.Contains(nameKey, "长度"):
		metricType = "length"
	case strings.Contains(nameKey, "comment") || strings.Contains(nameKey, "注释"):
		metricType = "comment"
	case strings.Contains(nameKey, "error") || strings.Contains(nameKey, "错误"):
		metricType = "error"
	case strings.Contains(nameKey, "naming") || strings.Contains(nameKey, "命名"):
		metricType = "naming"
	case strings.Contains(nameKey, "duplication") || strings.Contains(nameKey, "重复"):
		metricType = "duplication"
	case strings.Contains(nameKey, "structure") || strings.Contains(nameKey, "结构"):
		metricType = "structure"
	}

	if metricType != "" {
		commentKey = "metric." + metricType + "." + level
		return r.translator.Translate(commentKey)
	}

	switch r.translator.GetLanguage() {
	case i18n.EnUS:
		if score < 20 {
			return "Like a spring breeze, code kissed by angels"
		} else if score < 40 {
			return "Clean and polished, a craftsman's work"
		} else if score < 60 {
			return "Not bad, not great, perfectly balanced"
		} else if score < 80 {
			return "Shows signs of neglect, needs attention"
		} else {
			return "Needs serious improvement, like yesterday"
		}
	default:
		if score < 20 {
			return "如沐春风，代码仿佛被天使亲吻过"
		} else if score < 40 {
			return "干净整洁，工匠级的作品"
		} else if score < 60 {
			return "不咸不淡，刚刚好，就像人生的平凡日子"
		} else if score < 80 {
			return "显露疏于打理的痕迹，需要关注"
		} else {
			return "惨不忍睹，建议重写，或者假装没看见"
		}
	}
}

// printScoreComment 根据得分打印带颜色的总评
func (r *Report) printScoreComment(score float64) {
	comment := r.getScoreComment(score)

	switch {
	case score < 0.2:
		successStyle.Printf("%s", comment) // 深绿色 - 优秀
	case score < 0.4:
		goodStyle.Printf("%s", comment) // 浅绿色 - 良好
	case score < 0.55:
		color.New(color.FgHiBlue).Printf("%s", comment) // 蓝色 - 一般
	case score < 0.7:
		warningStyle.Printf("%s", comment) // 黄色 - 较差
	case score < 0.85:
		color.New(color.FgHiRed).Printf("%s", comment) // 亮红色 - 差
	default:
		dangerStyle.Printf("%s", comment) // 深红色 - 极差
	}
}

// printTopIssues 打印问题最多的几个代码文件及其问题
func (r *Report) printTopIssues(options *ReportOptions) {
	sectionStyle.Printf("\n◆ %s\n\n", r.translator.Translate("report.worst_files"))

	// 获取所有文件，按问题数量排序
	allFiles := r.getSortedFiles()

	// 如果没有文件，显示提示信息
	if len(allFiles) == 0 {
		successStyle.Println("  🎉 " + r.translator.Translate("report.no_issues"))
		return
	}

	// 计算文件路径最大长度，用于对齐
	maxPathLen := 0
	for _, file := range allFiles {
		pathLen := len(shortenPath(file.FilePath))
		if pathLen > maxPathLen {
			maxPathLen = pathLen
		}
	}

	maxPathLen = min(maxPathLen, 60)

	maxFiles := min(options.TopFiles, len(allFiles))

	for i := 0; i < maxFiles; i++ {
		f := allFiles[i]

		fileScoreColor := getScoreColor(f.FileScore)

		fmt.Printf("  ")
		numberStyle.Printf("%d. ", i+1)
		fileStyle.Printf("%-*s", maxPathLen+2, shortenPath(f.FilePath))
		fileScoreColor.Printf("(%s)\n", r.translator.Translate("report.file_score", math.Round(f.FileScore*10000)/100))

		issuesByCategory := r.categorizeIssues(f.Issues)

		if len(issuesByCategory) > 0 {
			categoryInfo := map[string]struct {
				Color *color.Color
				Icon  string
			}{
				"complexity":  {color.New(color.FgMagenta), "🔄 "},
				"comment":     {color.New(color.FgBlue), "📝 "},
				"naming":      {color.New(color.FgCyan), "🏷️  "},
				"structure":   {color.New(color.FgYellow), "🏗️  "},
				"duplication": {color.New(color.FgRed), "📋 "},
				"error":       {color.New(color.FgHiRed), "❌ "},
				"other":       {color.New(color.FgHiYellow), "⚠️  "},
			}

			categoryOrder := []string{"complexity", "comment", "naming", "structure", "duplication", "error", "other"}

			var categories []string
			for _, category := range categoryOrder {
				if count, exists := issuesByCategory[category]; exists {
					var categoryStr strings.Builder

					info := categoryInfo[category]
					categoryStr.WriteString(info.Icon)
					categoryStr.WriteString(r.translator.Translate("issue.category." + category))
					categoryStr.WriteString(": ")

					categories = append(categories, fmt.Sprintf("%s%d", categoryStr.String(), count))
				}
			}

			categoriesPerLine := 3
			if len(categories) <= 2 {
				categoriesPerLine = len(categories)
			}

			indent := "     "
			for i, category := range categories {
				if i > 0 && i%categoriesPerLine == 0 {
					fmt.Printf("\n%s", indent)
				} else if i > 0 {
					fmt.Printf("   ")
				} else {
					fmt.Printf("%s", indent)
				}

				parts := strings.SplitN(category, ":", 2)
				if len(parts) == 2 {
					for catName, info := range categoryInfo {
						catKey := "issue.category." + catName
						catTrans := r.translator.Translate(catKey)

						if strings.Contains(parts[0], catTrans) {
							info.Color.Printf("%s:", parts[0])
							numberStyle.Printf("%s", parts[1])
							break
						}
					}
				} else {
					fmt.Printf("%s", category)
				}
			}
			fmt.Println()
		}

		fmt.Println()
		indent := "     "

		if len(f.Issues) == 0 {
			successStyle.Printf("%s✓ %s\n", indent, r.translator.Translate("verbose.file_good_quality"))
		} else {
			maxIssues := min(options.MaxIssues, len(f.Issues))

			// 打印问题
			for j := 0; j < maxIssues; j++ {
				issueIcon, issueColor := r.getIssueIconAndColor(f.Issues[j])
				fmt.Printf("%s", indent)
				issueColor.Printf("%s%s\n", issueIcon, f.Issues[j])
			}

			if !options.Verbose && len(f.Issues) > maxIssues {
				warningStyle.Printf("%s🔍 %s\n",
					indent, r.translator.Translate("report.more_issues", len(f.Issues)-maxIssues))
			}
		}

		if i < maxFiles-1 {
			fmt.Println()
		}
	}
}

// categorizeIssues 将问题按类别分类统计
func (r *Report) categorizeIssues(issues []string) map[string]int {
	categories := map[string]int{
		"complexity":  0, // 复杂度问题
		"comment":     0, // 注释问题
		"naming":      0, // 命名问题
		"structure":   0, // 结构问题
		"duplication": 0, // 重复问题
		"error":       0, // 错误处理问题
		"other":       0, // 其他问题
	}

	for _, issue := range issues {
		lowerIssue := strings.ToLower(issue)

		switch {
		case strings.Contains(lowerIssue, "复杂度") || strings.Contains(lowerIssue, "complexity"):
			categories["complexity"]++
		case strings.Contains(lowerIssue, "注释") || strings.Contains(lowerIssue, "comment"):
			categories["comment"]++
		case strings.Contains(lowerIssue, "命名") || strings.Contains(lowerIssue, "name") || strings.Contains(lowerIssue, "naming"):
			categories["naming"]++
		case strings.Contains(lowerIssue, "结构") || strings.Contains(lowerIssue, "嵌套") || strings.Contains(lowerIssue, "structure") || strings.Contains(lowerIssue, "nest"):
			categories["structure"]++
		case strings.Contains(lowerIssue, "重复") || strings.Contains(lowerIssue, "duplication"):
			categories["duplication"]++
		case strings.Contains(lowerIssue, "错误") || strings.Contains(lowerIssue, "error"):
			categories["error"]++
		default:
			categories["other"]++
		}
	}

	// 删除计数为0的类别
	for category, count := range categories {
		if count == 0 {
			delete(categories, category)
		}
	}

	return categories
}

// getIssueIconAndColor 根据问题内容返回合适的图标和颜色
func (r *Report) getIssueIconAndColor(issue string) (string, *color.Color) {
	lowerIssue := strings.ToLower(issue)

	switch {
	case strings.Contains(lowerIssue, "复杂度") || strings.Contains(lowerIssue, "complexity"):
		return "🔄 ", color.New(color.FgMagenta) // 窄图标，只需一个空格
	case strings.Contains(lowerIssue, "注释") || strings.Contains(lowerIssue, "comment"):
		return "📝 ", color.New(color.FgBlue) // 窄图标，只需一个空格
	case strings.Contains(lowerIssue, "命名") || strings.Contains(lowerIssue, "name") || strings.Contains(lowerIssue, "naming"):
		return "🏷️  ", color.New(color.FgCyan) // 宽图标，需要两个空格
	case strings.Contains(lowerIssue, "结构") || strings.Contains(lowerIssue, "嵌套") || strings.Contains(lowerIssue, "structure") || strings.Contains(lowerIssue, "nest"):
		return "🏗️  ", color.New(color.FgYellow) // 宽图标，需要两个空格
	case strings.Contains(lowerIssue, "重复") || strings.Contains(lowerIssue, "duplication"):
		return "📋 ", color.New(color.FgRed) // 窄图标，只需一个空格
	case strings.Contains(lowerIssue, "错误") || strings.Contains(lowerIssue, "error"):
		return "❌ ", color.New(color.FgHiRed) // 窄图标，只需一个空格
	default:
		return "⚠️  ", color.New(color.FgHiYellow) // 宽图标，需要两个空格
	}
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// shortenPath 缩短文件路径，只显示最后几个部分
func shortenPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) <= 4 {
		return path
	}

	return "./" + strings.Join(parts[len(parts)-3:], "/")
}

// printSummary 打印最终诊断结论和建议
func (r *Report) printSummary(level struct {
	MinScore    float64
	NameKey     string
	Description string
	Emoji       string
}) {
	sectionStyle.Printf("\n◆ %s\n\n", r.translator.Translate("report.conclusion"))

	// 使用levelStyle打印等级名称和表情符号
	fmt.Printf("  %s ", level.Emoji)
	levelStyle.Printf("%s", r.translator.Translate(level.NameKey))
	detailStyle.Printf(" - %s\n\n", r.translator.Translate(level.Description))

	// 根据不同等级提供相应的建议
	switch {
	case level.MinScore < 30:
		successStyle.Println("  " + r.translator.Translate("advice.good"))
	case level.MinScore < 60:
		warningStyle.Println("  " + r.translator.Translate("advice.moderate"))
	default:
		dangerStyle.Println("  " + r.translator.Translate("advice.bad"))
	}
	fmt.Println()
}

// getScoreComment 根据得分生成总评
func (r *Report) getScoreComment(score float64) string {
	score = score * 100

	scoreRange := int(score) / 10 * 10
	if scoreRange > 90 {
		scoreRange = 90
	}

	commentKey := fmt.Sprintf("score.comment.%d", scoreRange)
	return r.translator.Translate(commentKey)
}

// getSortedMetrics 按照分数升序排列各项指标
func (r *Report) getSortedMetrics() []analyzer.MetricResult {
	var metrics []analyzer.MetricResult
	for _, m := range r.result.Metrics {
		metrics = append(metrics, m)
	}
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].Score < metrics[j].Score
	})
	return metrics
}

// getSortedFiles 获取按分数排序的问题文件列表（降序）
func (r *Report) getSortedFiles() []analyzer.FileAnalysisResult {
	worstFiles := append([]analyzer.FileAnalysisResult{}, r.result.FilesAnalyzed...)
	sort.Slice(worstFiles, func(i, j int) bool {
		return worstFiles[i].FileScore > worstFiles[j].FileScore
	})
	return worstFiles
}

// getQualityLevel 根据总分获取对应的质量等级
func (r *Report) getQualityLevel(score float64) struct {
	MinScore    float64
	NameKey     string
	Description string
	Emoji       string
} {
	adjustedScore := score

	if score > 0.2 && score < 0.8 {
		// 使用S形曲线公式: 1/(1+e^(-k*(x-0.5)))
		shift := (score - 0.5) * 10
		sigmoidValue := 1.0 / (1.0 + math.Exp(-shift))

		adjustedScore = sigmoidValue
	}

	// 转换为百分比
	adjustedScore = adjustedScore * 100

	level := QualityLevels[0]
	for i := len(QualityLevels) - 1; i >= 0; i-- {
		if adjustedScore >= QualityLevels[i].MinScore {
			level = QualityLevels[i]
			break
		}
	}
	return level
}

// printVerboseInfo 打印详细的代码分析信息
func (r *Report) printVerboseInfo() {
	sectionStyle.Printf("\n◆ %s\n\n", r.translator.Translate("verbose.basic_statistics"))

	// 打印基本统计数据
	headerStyle.Println("  📊 " + r.translator.Translate("verbose.basic_statistics"))
	detailStyle.Printf("    %-15s %d\n", r.translator.Translate("verbose.total_files"), r.result.TotalFiles)
	detailStyle.Printf("    %-15s %d\n", r.translator.Translate("verbose.total_lines"), r.result.TotalLines)
	detailStyle.Printf("    %-15s %d\n", r.translator.Translate("verbose.total_issues"), r.getTotalIssues())

	// 打印各指标详细信息
	headerStyle.Println("\n  🔍 " + r.translator.Translate("verbose.metric_details"))

	metrics := r.getSortedMetrics()
	maxNameLen := 0
	for _, metric := range metrics {
		if len(metric.Name) > maxNameLen {
			maxNameLen = len(metric.Name)
		}
	}

	nameFormat := fmt.Sprintf("\n    【%%-%ds】", maxNameLen)

	for _, metric := range metrics {
		scoreColor := getScoreColor(metric.Score)
		metricStyle.Printf(nameFormat, metric.Name)
		infoStyle.Printf("(%s %.2f)\n", r.translator.Translate("verbose.weight"), metric.Weight)
		detailStyle.Printf("      %s %s\n", r.translator.Translate("verbose.description"), metric.Description)
		fmt.Printf("      %s ", r.translator.Translate("verbose.score"))
		// 精确到小数点后2位
		scoreColor.Printf("%.2f/100\n", math.Round(metric.Score*10000)/100)
	}
}

// getTotalIssues 获取所有文件的问题总数
func (r *Report) getTotalIssues() int {
	total := 0
	for _, file := range r.result.FilesAnalyzed {
		total += len(file.Issues)
	}
	return total
}

// getScoreColor 根据得分返回对应的颜色
func getScoreColor(score float64) *color.Color {
	switch {
	case score < 0.2:
		return color.New(color.FgHiGreen, color.Bold) // 加粗亮绿色 - 优秀
	case score < 0.35:
		return goodStyle // 绿色 - 良好
	case score < 0.5:
		return color.New(color.FgHiCyan) // 青色 - 一般偏好
	case score < 0.6:
		return color.New(color.FgBlue) // 蓝色 - 一般
	case score < 0.7:
		return color.New(color.FgHiYellow) // 亮黄色 - 一般偏差
	case score < 0.8:
		return warningStyle // 黄色 - 较差
	case score < 0.9:
		return color.New(color.FgHiRed) // 亮红色 - 差
	default:
		return dangerStyle // 红色 - 极差
	}
}

// printAllFiles 打印所有文件及其问题
func (r *Report) printAllFiles(options *ReportOptions) {
	sectionStyle.Printf("\n◆ %s\n\n", r.translator.Translate("verbose.all_files"))

	files := r.getSortedFiles()
	if len(files) == 0 {
		successStyle.Println("  " + r.translator.Translate("verbose.no_files_found"))
		return
	}

	// 计算文件路径最大长度，用于对齐
	maxPathLen := 0
	for _, file := range files {
		pathLen := len(shortenPath(file.FilePath))
		if pathLen > maxPathLen {
			maxPathLen = pathLen
		}
	}

	// 限制最大宽度
	maxPathLen = min(maxPathLen, 60)

	// 根据options.TopFiles决定显示多少文件
	maxFilesToShow := len(files)
	if !options.Verbose && options.TopFiles > 0 && options.TopFiles < maxFilesToShow {
		maxFilesToShow = options.TopFiles
	}

	for i, f := range files[:maxFilesToShow] {
		// 根据得分选择颜色
		fileScoreColor := getScoreColor(f.FileScore)

		// 打印文件名和得分
		fmt.Printf("  ")
		numberStyle.Printf("%d. ", i+1)
		fileStyle.Printf("%-*s", maxPathLen+2, shortenPath(f.FilePath))
		fileScoreColor.Printf("(%s)\n", r.translator.Translate("report.file_score", math.Round(f.FileScore*10000)/100))

		// 分类统计问题
		issuesByCategory := r.categorizeIssues(f.Issues)

		// 打印问题分类统计 - 使用更紧凑美观的布局
		if len(issuesByCategory) > 0 {
			// 定义优雅的颜色组合和图标
			categoryInfo := map[string]struct {
				Color *color.Color
				Icon  string
			}{
				"complexity":  {color.New(color.FgMagenta), "🔄 "},
				"comment":     {color.New(color.FgBlue), "📝 "},
				"naming":      {color.New(color.FgCyan), "🏷️  "},
				"structure":   {color.New(color.FgYellow), "🏗️  "},
				"duplication": {color.New(color.FgRed), "📋 "},
				"error":       {color.New(color.FgHiRed), "❌ "},
				"other":       {color.New(color.FgHiYellow), "⚠️  "},
			}

			// 定义问题类别的显示顺序
			categoryOrder := []string{"complexity", "comment", "naming", "structure", "duplication", "error", "other"}

			// 创建一个紧凑的类别统计字符串
			var categories []string
			for _, category := range categoryOrder {
				if count, exists := issuesByCategory[category]; exists {
					// 使用字符串构建器创建每个类别的显示
					var categoryStr strings.Builder

					// 使用颜色写入图标和类别名称
					info := categoryInfo[category]
					categoryStr.WriteString(info.Icon)
					categoryStr.WriteString(r.translator.Translate("issue.category." + category))
					categoryStr.WriteString(": ")

					// 添加到类别列表
					categories = append(categories, fmt.Sprintf("%s%d", categoryStr.String(), count))
				}
			}

			// 计算每行显示的类别数量
			categoriesPerLine := 3
			if len(categories) <= 2 {
				categoriesPerLine = len(categories)
			}

			// 打印类别统计
			indent := "     "
			for i, category := range categories {
				if i > 0 && i%categoriesPerLine == 0 {
					fmt.Printf("\n%s", indent)
				} else if i > 0 {
					fmt.Printf("   ")
				} else {
					fmt.Printf("%s", indent)
				}

				// 解析类别字符串并使用适当的颜色打印
				parts := strings.SplitN(category, ":", 2)
				if len(parts) == 2 {
					// 找出对应的类别以获取颜色
					for catName, info := range categoryInfo {
						catKey := "issue.category." + catName
						catTrans := r.translator.Translate(catKey)

						if strings.Contains(parts[0], catTrans) {
							// 使用颜色打印类别名称和图标
							info.Color.Printf("%s:", parts[0])
							// 使用数字样式打印计数
							numberStyle.Printf("%s", parts[1])
							break
						}
					}
				} else {
					// 回退方案
					fmt.Printf("%s", category)
				}
			}
			fmt.Println()
		}

		// 打印问题列表
		fmt.Println()
		indent := "     "
		if len(f.Issues) == 0 {
			// 如果没有问题，显示"无问题"提示，手动添加✓图标
			successStyle.Printf("%s✓ %s\n", indent, r.translator.Translate("verbose.file_good_quality"))
		} else {
			// 确定显示多少问题
			maxIssues := len(f.Issues)
			if !options.Verbose {
				maxIssues = min(options.MaxIssues, maxIssues)
			}

			for j := 0; j < maxIssues; j++ {
				issueIcon, issueColor := r.getIssueIconAndColor(f.Issues[j])
				fmt.Printf("%s", indent)
				issueColor.Printf("%s%s\n", issueIcon, f.Issues[j])
			}

			// 只在非详细模式下显示"还有更多问题"的提示
			if !options.Verbose && len(f.Issues) > maxIssues {
				warningStyle.Printf("%s🔍 %s\n",
					indent, r.translator.Translate("report.more_issues", len(f.Issues)-maxIssues))
			}
		}

		if i < maxFilesToShow-1 {
			fmt.Println()
		}
	}
}
