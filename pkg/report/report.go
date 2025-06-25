// Package report 提供代码质量分析报告生成功能
package report

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"

	"github.com/Done-0/fuck-u-code/pkg/analyzer"
)

// 颜色定义
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
	issueStyle   = color.New(color.FgYellow)
	detailStyle  = color.New(color.FgCyan)
	metricStyle  = color.New(color.FgCyan)
	fileStyle    = color.New(color.FgMagenta)
	levelStyle   = color.New(color.FgCyan)
	numberStyle  = color.New(color.FgHiWhite)
)

// ShiShanLevels 定义屎山等级（每10分一个段位）
var ShiShanLevels = []struct {
	MinScore    float64
	Name        string
	Description string
	Emoji       string
}{
	{0, "清新可人", "代码洁净，令人赏心悦目", "🌱"},
	{10, "偶有异味", "基本没事，但是有伤风化", "🌸"},
	{20, "微臭青年", "略有异味，建议适量通风", "😐"},
	{30, "屎气扑鼻", "代码开始散发气味，谨慎维护", "😷"},
	{40, "中度屎山", "臭味明显，开窗也救不了", "💩"},
	{50, "隐性毒瘤", "写的时候爽，改的时候哭", "🤕"},
	{60, "重度屎山", "毒气弥漫，建议戴防毒面具", "☣️"},
	{70, "代码化尸场", "进去的程序员没有一个活着出来", "🧟"},
	{80, "核平级灾难", "反人类罪行，建议火化", "☢️"},
	{90, "祖传老屎", "历代工程师共创的遗产，无法维护", "🪦"},
	{100, "终极屎王", "写的时候热血澎湃，改的时候亲妈不认", "👑💩"},
}

// Report 表示屎山分析报告对象
type Report struct {
	result *analyzer.AnalysisResult
}

// NewReport 创建一个新的报告实例
func NewReport(result *analyzer.AnalysisResult) *Report {
	return &Report{result: result}
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
	level := r.getShiShanLevel(score)

	// 打印标题和总体评分
	printDivider()
	titleStyle.Printf("\n  %s 屎山代码分析报告 %s\n", level.Emoji, level.Emoji)
	printDivider()

	fmt.Printf("\n")
	scoreStyle.Printf("  总体评分: %.1f / 100", score)
	fmt.Printf(" - ")
	r.printScoreComment(score)
	fmt.Printf("\n")

	// 打印屎山等级
	detailStyle.Printf("  屎山等级: ")
	levelStyle.Printf("%s ", level.Name)
	detailStyle.Printf("- %s\n\n", level.Description)

	if !options.SummaryOnly {
		r.printMetricItems() // 打印各项评分指标

		// 详细模式下显示所有文件，否则只显示最差的几个
		if options.Verbose {
			r.printAllFiles(options)
		} else {
			r.printTopIssues(options) // 打印最屎的代码
		}
	}

	r.printSummary(level) // 打印总结建议

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
	sectionStyle.Printf("\n◆ 评分指标详情\n\n")

	metrics := r.getSortedMetrics()

	// 计算对齐所需的宽度
	maxNameLen := 0
	for _, m := range metrics {
		if len(m.Name) > maxNameLen {
			maxNameLen = len(m.Name)
		}
	}

	// 格式化模板
	nameFormat := fmt.Sprintf("  %%s %%-%ds", maxNameLen+2)
	scoreFormat := "%-8s"

	for _, m := range metrics {
		scorePercentage := m.Score * 100
		addedPoints := scorePercentage / 5.0

		// 确定状态图标和颜色
		var statusEmoji string
		var statusColor *color.Color

		if scorePercentage < 20 {
			statusEmoji = "✓"
			statusColor = goodStyle
		} else if scorePercentage < 50 {
			statusEmoji = "!"
			statusColor = warningStyle
		} else {
			statusEmoji = "✗"
			statusColor = dangerStyle
		}

		// 格式化分数
		scoreStr := fmt.Sprintf("%.1f分", addedPoints)

		// 打印一行，确保对齐
		statusColor.Printf(nameFormat, statusEmoji, m.Name)
		metricStyle.Printf(scoreFormat, scoreStr)
		detailStyle.Printf("%s\n", r.getMetricComment(m.Name, scorePercentage))
	}
	fmt.Println()
}

// getMetricComment 返回指标评论
func (r *Report) getMetricComment(metricName string, score float64) string {
	switch metricName {
	case "循环复杂度":
		if score > 50 {
			return "函数像迷宫，维护像打副本"
		} else if score > 20 {
			return "绕来绕去，跟你脑子一样乱"
		} else {
			return "结构清晰，不绕弯子，赞👍"
		}
	case "函数长度":
		if score > 50 {
			return "这是函数？这怕不是《代码圣经》"
		} else if score > 20 {
			return "代码小作文，读完脑子疼"
		} else {
			return "函数短小精悍，一眼一个小目标"
		}
	case "注释覆盖率":
		if score > 50 {
			return "没有注释，靠缘分理解"
		} else if score > 20 {
			return "注释稀薄，读者全靠脑补"
		} else {
			return "注释不错，能靠它活下来"
		}
	case "错误处理":
		if score > 50 {
			return "err 见了就跳过？宛如人生"
		} else if score > 20 {
			return "有处理，但处理得跟没处理一样"
		} else {
			return "错误都照顾到了，代码有大爱"
		}
	case "命名规范":
		if score > 50 {
			return "变量名像键盘砸出来的：x, y, z, tmp, xxx"
		} else if score > 20 {
			return "命名还行，有些得猜"
		} else {
			return "命名清晰，程序员的文明之光"
		}
	case "代码重复度":
		if score > 50 {
			return "一眼复制痕迹，Ctrl+C/V 荣誉勋章"
		} else if score > 20 {
			return "有点重复，抽象一下不难吧"
		} else {
			return "该抽象的都抽象了，强迫症舒服了"
		}
	default:
		return "有点问题，但还没臭到家"
	}
}

// printScoreComment 根据得分打印带颜色的总评
func (r *Report) printScoreComment(score float64) {
	comment := r.getScoreComment(score)

	if score < 30 {
		goodStyle.Printf("%s", comment)
	} else if score < 70 {
		warningStyle.Printf("%s", comment)
	} else {
		dangerStyle.Printf("%s", comment)
	}
}

// printTopIssues 打印最屎的几个代码文件及其问题
func (r *Report) printTopIssues(options *ReportOptions) {
	sectionStyle.Printf("\n◆ 最屎代码排行榜\n\n")

	worstFiles := r.getSortedFiles()
	maxFiles := min(options.TopFiles, len(worstFiles))

	if maxFiles == 0 {
		successStyle.Println("  🎉 恭喜！没有特别屎的文件！")
		return
	}

	// 计算文件路径最大长度，用于对齐
	maxPathLen := 0
	for i := 0; i < maxFiles; i++ {
		pathLen := len(shortenPath(worstFiles[i].FilePath))
		if pathLen > maxPathLen {
			maxPathLen = pathLen
		}
	}

	// 限制最大宽度
	maxPathLen = min(maxPathLen, 60)

	for i := 0; i < maxFiles; i++ {
		f := worstFiles[i]

		// 根据得分选择颜色
		fileScoreColor := getScoreColor(f.FileScore)

		// 打印文件名和得分
		fmt.Printf("  ")
		numberStyle.Printf("%d. ", i+1)
		fileStyle.Printf("%-*s", maxPathLen+2, shortenPath(f.FilePath))
		fileScoreColor.Printf("(屎气指数: %.1f)\n", f.FileScore*100)

		// 打印问题列表
		maxIssues := len(f.Issues)
		if !options.Verbose {
			maxIssues = min(options.MaxIssues, maxIssues)
		}

		indent := "     "
		for j := 0; j < maxIssues; j++ {
			issueStyle.Printf("%s⚠️  %s\n", indent, f.Issues[j])
		}

		// 只在非详细模式下显示"还有更多问题"的提示
		if !options.Verbose && len(f.Issues) > maxIssues {
			warningStyle.Printf("%s🔥 ...还有 %d 个问题实在太屎，列不完了\n",
				indent, len(f.Issues)-maxIssues)
		}

		if i < maxFiles-1 {
			fmt.Println()
		}
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
	Name        string
	Description string
	Emoji       string
}) {
	score := r.result.CodeQualityScore

	sectionStyle.Printf("\n◆ 诊断结论\n\n")

	// 使用levelStyle打印等级名称
	fmt.Printf("  %s ", level.Emoji)
	levelStyle.Printf("%s", level.Name)
	detailStyle.Printf(" - %s\n\n", level.Description)

	switch {
	case score < 30:
		successStyle.Println("  👍 继续保持，你是编码界的一股清流")
	case score < 60:
		warningStyle.Println("  🔧 建议：适量重构，减少技术债，不然以后哭的是你")
	default:
		dangerStyle.Println("  🧨 建议：这代码建议删库跑路，或者封印它别再动了")
	}
	fmt.Println()
}

// getScoreComment 根据得分生成总评
func (r *Report) getScoreComment(score float64) string {
	switch {
	case score < 10:
		return "如沐春风"
	case score < 20:
		return "略带芳香"
	case score < 30:
		return "有点臭味，尚可忍受"
	case score < 40:
		return "气味渐浓，小心中毒"
	case score < 50:
		return "屎气扑鼻，慎入"
	case score < 60:
		return "屎味稳定释放中"
	case score < 70:
		return "代码已具毒性"
	case score < 80:
		return "臭不可闻，建议远离"
	case score < 90:
		return "深度污染区"
	default:
		return "屎味浓郁！建议火化"
	}
}

// getSortedMetrics 按照分数降序排列各项指标
func (r *Report) getSortedMetrics() []analyzer.MetricResult {
	var metrics []analyzer.MetricResult
	for _, m := range r.result.Metrics {
		metrics = append(metrics, m)
	}
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].Score > metrics[j].Score
	})
	return metrics
}

// getSortedFiles 获取按分数排序的屎文件列表（降序）
func (r *Report) getSortedFiles() []analyzer.FileAnalysisResult {
	worstFiles := append([]analyzer.FileAnalysisResult{}, r.result.FilesAnalyzed...)
	sort.Slice(worstFiles, func(i, j int) bool {
		return worstFiles[i].FileScore > worstFiles[j].FileScore
	})
	return worstFiles
}

// getShiShanLevel 根据总分获取对应的屎山等级
func (r *Report) getShiShanLevel(score float64) struct {
	MinScore    float64
	Name        string
	Description string
	Emoji       string
} {
	level := ShiShanLevels[0]
	for i := len(ShiShanLevels) - 1; i >= 0; i-- {
		if score >= ShiShanLevels[i].MinScore {
			level = ShiShanLevels[i]
			break
		}
	}
	return level
}

// printVerboseInfo 打印详细的代码分析信息
func (r *Report) printVerboseInfo() {
	sectionStyle.Printf("\n◆ 详细分析信息\n\n")

	// 打印基本统计数据
	headerStyle.Println("  📊 基本统计:")
	detailStyle.Printf("    %-15s %d\n", "总文件数:", r.result.TotalFiles)
	detailStyle.Printf("    %-15s %d\n", "总代码行:", r.result.TotalLines)
	detailStyle.Printf("    %-15s %d\n", "总问题数:", r.getTotalIssues())

	// 打印各指标详细信息
	headerStyle.Println("\n  🔍 指标详细信息:")

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
		infoStyle.Printf("(权重: %.2f)\n", metric.Weight)
		detailStyle.Printf("      描述: %s\n", metric.Description)
		fmt.Printf("      得分: ")
		scoreColor.Printf("%.1f/100\n", metric.Score*100)
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
	if score > 0.7 {
		return dangerStyle
	} else if score > 0.3 {
		return warningStyle
	}
	return goodStyle
}

// printAllFiles 打印所有文件及其问题
func (r *Report) printAllFiles(options *ReportOptions) {
	sectionStyle.Printf("\n◆ 全部代码文件分析\n\n")

	files := r.getSortedFiles()
	if len(files) == 0 {
		successStyle.Println("  🎉 没有找到需要分析的文件！")
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

		// 打印文件名和得分，序号使用numberStyle
		fmt.Printf("  ")
		numberStyle.Printf("%d. ", i+1)
		fileStyle.Printf("%-*s", maxPathLen+2, shortenPath(f.FilePath))
		fileScoreColor.Printf("(屎气指数: %.1f)\n", f.FileScore*100)

		// 打印问题列表
		indent := "     "
		if len(f.Issues) == 0 {
			successStyle.Printf("%s✓ 代码干净，没有明显问题\n", indent)
		} else {
			for _, issue := range f.Issues {
				issueStyle.Printf("%s⚠️  %s\n", indent, issue)
			}
		}

		if i < maxFilesToShow-1 {
			fmt.Println()
		}
	}
}
