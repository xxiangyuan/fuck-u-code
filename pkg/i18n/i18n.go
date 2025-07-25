// Package i18n 提供国际化和本地化支持
// 创建者：Done-0
package i18n

import (
	"fmt"
	"strings"
)

// Language 表示支持的语言
type Language string

const (
	// ZhCN 简体中文
	ZhCN Language = "zh-CN"

	// EnUS 英文（美国）
	EnUS Language = "en-US"
)

// Translator 翻译器接口
type Translator interface {
	// Translate 翻译指定的键
	Translate(key string, args ...interface{}) string

	// GetLanguage 获取当前语言
	GetLanguage() Language
}

// DefaultTranslator 默认翻译器实现
type DefaultTranslator struct {
	language Language
	messages map[string]string
}

// NewTranslator 创建新的翻译器
func NewTranslator(language Language) Translator {
	translator := &DefaultTranslator{
		language: language,
		messages: make(map[string]string),
	}

	// 加载语言包
	translator.loadMessages()

	return translator
}

// Translate 翻译指定的键
func (t *DefaultTranslator) Translate(key string, args ...interface{}) string {
	if msg, ok := t.messages[key]; ok {
		if len(args) > 0 {
			return fmt.Sprintf(msg, args...)
		}
		return msg
	}

	// 如果找不到翻译，返回键本身
	return key
}

// GetLanguage 获取当前语言
func (t *DefaultTranslator) GetLanguage() Language {
	return t.language
}

// loadMessages 加载语言包
func (t *DefaultTranslator) loadMessages() {
	switch t.language {
	case ZhCN:
		t.messages = zhCNMessages
	case EnUS:
		t.messages = enUSMessages
	default:
		t.messages = enUSMessages // 默认使用英文
	}
}

// FormatKey 格式化翻译键，将多个部分组合成一个键
func FormatKey(parts ...string) string {
	return strings.Join(parts, ".")
}

// 中文语言包
var zhCNMessages = map[string]string{
	// 通用
	"app.name":        "屎山代码鉴定器",
	"app.description": "一个专为挖掘 Go 项目\"屎坑\"设计的工具，能无情揭露代码的丑陋真相，并用毫不留情的幽默语言告诉你：你的代码到底有多烂。",

	// 指标名称
	"metric.cyclomatic_complexity": "循环复杂度",
	"metric.function_length":       "状态管理",
	"metric.comment_ratio":         "注释覆盖率",
	"metric.error_handling":        "错误处理",
	"metric.naming_convention":     "命名规范",
	"metric.code_duplication":      "代码重复度",
	"metric.structure_analysis":    "代码结构",

	// 分析器进度
	"analyzer.searching_files":   "正在搜索源代码文件...",
	"analyzer.files_found":       "已找到文件数",
	"analyzer.analyzing_files":   "正在分析文件...",
	"analyzer.analysis_complete": "分析完成",

	// 问题分类
	"report.no_issues":           "恭喜！没有特别多问题的文件！",
	"issue.category.complexity":  "复杂度问题",
	"issue.category.comment":     "注释问题",
	"issue.category.naming":      "命名问题",
	"issue.category.structure":   "结构问题",
	"issue.category.duplication": "重复问题",
	"issue.category.error":       "错误处理问题",
	"issue.category.other":       "其他问题",

	// 质量等级
	"level.clean":             "清新可人",
	"level.mild":              "偶有异味",
	"level.moderate":          "微臭青年",
	"level.bad":               "屎气扑鼻",
	"level.terrible":          "中度屎山",
	"level.disaster":          "隐性毒瘤",
	"level.disaster.severe":   "重度屎山",
	"level.disaster.very_bad": "代码化尸场",
	"level.disaster.extreme":  "核平级灾难",
	"level.disaster.worst":    "祖传老屎",
	"level.disaster.ultimate": "终极屎王",

	// 命令行
	"cmd.short":                      "💻 fuck-u-code",
	"cmd.long":                       "🔍 屎山代码鉴定器 - 客观评估您的代码质量\n\n它可以分析代码质量、输出评分，帮助您发现代码中的💩。适用于：\n- 项目重构前的质量评估\n- 团队代码审查辅助工具\n- 学习编程最佳实践",
	"cmd.analyze":                    "分析代码质量并输出评分",
	"cmd.analyze.long":               "深入分析代码库，检测各种代码潜在问题，输出质量报告。不指定路径时分析当前目录。",
	"cmd.completion":                 "生成自动补全脚本",
	"cmd.completion.long":            "为指定的shell生成自动补全脚本，支持bash、zsh、fish和PowerShell。",
	"cmd.completion.long_prefix":     "为指定的shell生成fuck-u-code的自动补全脚本。",
	"cmd.completion.long_suffix":     "查看每个子命令的帮助，了解如何使用生成的脚本。",
	"cmd.completion.bash":            "为bash生成自动补全脚本",
	"cmd.completion.zsh":             "为zsh生成自动补全脚本",
	"cmd.completion.fish":            "为fish生成自动补全脚本",
	"cmd.completion.powershell":      "为powershell生成自动补全脚本",
	"cmd.completion.bash.long":       "为bash shell生成自动补全脚本",
	"cmd.completion.zsh.long":        "为zsh shell生成自动补全脚本",
	"cmd.completion.fish.long":       "为fish shell生成自动补全脚本",
	"cmd.completion.powershell.long": "为powershell生成自动补全脚本",
	"cmd.help":                       "获取帮助信息",
	"cmd.help.long":                  "获取关于任何命令的帮助信息。",
	"cmd.help_flag":                  "获取关于fuck-u-code的帮助",
	"cmd.no_descriptions":            "禁用补全描述",
	"cmd.path_not_found":             "路径不可访问 '%s': %v",
	"cmd.analysis_failed":            "分析失败：%v",
	"cmd.lang":                       "指定输出语言（支持：zh-CN, en-US，默认：zh-CN）",
	"cmd.verbose":                    "显示详细分析报告",
	"cmd.top":                        "显示问题最多的文件数量（默认5个）",
	"cmd.issues":                     "每个文件显示多少条问题（默认5个）",
	"cmd.summary":                    "只看结论，过程略过",
	"cmd.markdown":                   "输出Markdown格式的精简报告，便于AI工具处理",
	"cmd.exclude":                    "排除的文件/目录模式 (可多次使用，默认已排除常见依赖目录)",
	"cmd.start_analyzing":            "开始嗅探：%s",
	"cmd.exclude_patterns":           "排除以下文件/目录模式:",

	// Cobra框架内部文本
	"cobra.available_commands": "可用命令",
	"cobra.flags":              "选项",
	"cobra.global_flags":       "全局选项",
	"cobra.additional_help":    "附加帮助主题",
	"cobra.use_help_cmd":       "使用",
	"cobra.for_more_info":      "获取关于命令的更多信息",
	"cobra.usage":              "用法",

	// 报告
	"report.title":                   "屎山代码分析报告",
	"report.overall_score":           "总体评分: %.2f / 100",
	"report.level":                   "屎山等级: %s",
	"report.metrics_details":         "评分指标详情",
	"report.worst_files":             "最屎代码排行榜",
	"report.conclusion":              "诊断结论",
	"report.file_score":              "屎气指数: %.2f",
	"report.more_issues":             "...还有 %d 个问题实在太屎，列不完了",
	"report.score_calc":              "评分计算: ",
	"report.overall_assessment":      "总体评估",
	"report.quality_score":           "质量评分",
	"report.quality_level":           "质量等级",
	"report.analyzed_files":          "分析文件数",
	"report.total_lines":             "代码总行数",
	"report.quality_metrics":         "质量指标",
	"report.metric":                  "指标",
	"report.score":                   "得分",
	"report.weight":                  "权重",
	"report.status":                  "状态",
	"report.problem_files":           "问题文件",
	"report.issue_categories":        "问题分类",
	"report.main_issues":             "主要问题",
	"report.and":                     "还有",
	"report.more_issues_short":       "个问题",
	"report.improvement_suggestions": "改进建议",

	// 指标评分后缀
	"metric.score.suffix": "分",

	// 循环复杂度评价
	"metric.complexity.good":   "结构清晰，不绕弯子，赞",
	"metric.complexity.medium": "绕来绕去，跟你脑子一样乱",
	"metric.complexity.bad":    "函数像迷宫，维护像打副本",

	// 函数长度评价
	"metric.length.good":   "状态管理清晰，变量作用域合理，状态可预测",
	"metric.length.medium": "状态管理一般，存在部分全局状态或状态变化不明确的情况",
	"metric.length.bad":    "状态管理混乱，大量使用全局变量，状态变化难以追踪",

	// 注释覆盖率评价
	"metric.comment.good":   "注释不错，能靠它活下来",
	"metric.comment.medium": "注释稀薄，读者全靠脑补",
	"metric.comment.bad":    "没有注释，靠缘分理解",

	// 错误处理评价
	"metric.error.good":   "错误都照顾到了，代码有大爱",
	"metric.error.medium": "有处理，但处理得跟没处理一样",
	"metric.error.bad":    "err 见了就跳过？宛如人生",

	// 命名规范评价
	"metric.naming.good":   "命名清晰，程序员的文明之光",
	"metric.naming.medium": "命名还行，有些得猜",
	"metric.naming.bad":    "变量名像键盘砸出来的：x, y, z, tmp, xxx",

	// 代码重复度评价
	"metric.duplication.good":   "该抽象的都抽象了，强迫症舒服了",
	"metric.duplication.medium": "有点重复，抽象一下不难吧",
	"metric.duplication.bad":    "一眼复制痕迹，Ctrl+C/V 荣誉勋章",

	// 代码结构评价
	"metric.structure.good":   "结构优美，不容易看岔",
	"metric.structure.medium": "结构还行，但有点混乱",
	"metric.structure.bad":    "层层嵌套，套娃结构，看完眼花",

	// 质量建议
	"advice.good":     "👍 继续保持，你是编码界的一股清流，代码洁癖者的骄傲",
	"advice.moderate": "🔧 建议：这代码像个叛逆期的青少年，需要适当管教才能成才",
	"advice.bad":      "🧨 建议：删库跑路是唯一出路，或者封印它，下辈子再来维护",

	// 改进建议优先级
	"advice.priority.high":   "高优先级",
	"advice.priority.medium": "中优先级",
	"advice.priority.low":    "低优先级",

	// 良好代码的建议
	"advice.good.maintain": "继续保持当前的代码质量标准",
	"advice.good.optimize": "可以考虑进一步优化性能和可读性",
	"advice.good.document": "完善文档和注释，便于团队协作",

	// 中等代码的建议
	"advice.moderate.refactor":    "重构复杂度过高的函数和模块",
	"advice.moderate.complexity":  "降低循环复杂度，简化控制流",
	"advice.moderate.naming":      "改善变量和函数命名规范",
	"advice.moderate.comments":    "增加代码注释覆盖率",
	"advice.moderate.duplication": "消除重复代码，提取公共方法",
	"advice.moderate.structure":   "优化代码结构，减少嵌套层级",
	"advice.moderate.style":       "统一代码风格和格式",

	// 较差代码的建议
	"advice.bad.urgent_refactor": "紧急重构过长函数，遵循单一职责原则",
	"advice.bad.complexity":      "大幅降低循环复杂度，拆分复杂逻辑",
	"advice.bad.error_handling":  "添加完善的错误处理机制",
	"advice.bad.naming":          "全面改善命名规范，避免使用无意义变量名",
	"advice.bad.duplication":     "彻底消除重复代码，建立代码复用机制",
	"advice.bad.comments":        "大幅增加代码注释，提高可读性",
	"advice.bad.structure":       "重新设计代码架构，改善整体结构",
	"advice.bad.style":           "建立并执行严格的代码规范",

	// 指标描述
	"metric.function_length.description":       "检测代码中状态变量的管理，良好的状态管理能提高代码可维护性和可预测性",
	"metric.comment_ratio.description":         "检测代码的注释覆盖率，良好的注释能提高代码可读性和可维护性",
	"metric.error_handling.description":        "检测代码中的错误处理情况，良好的错误处理能提高代码的健壮性",
	"metric.naming_convention.description":     "检测代码中的命名规范，良好的命名能提高代码可读性",
	"metric.code_duplication.description":      "评估代码中重复逻辑的比例，重复代码越多，越需要抽象和重构",
	"metric.structure_analysis.description":    "检测代码的嵌套深度和引用复杂度，评估结构清晰度",
	"metric.cyclomatic_complexity.description": "测量函数的控制流复杂度，复杂度越高，代码越难理解和测试",

	// 质量等级描述
	"level.clean.description":             "代码洁净，令人赏心悦目",
	"level.mild.description":              "基本没事，但是有伤风化",
	"level.moderate.description":          "略有异味，建议适量通风",
	"level.bad.description":               "代码开始散发气味，谨慎维护",
	"level.terrible.description":          "臭味明显，开窗也救不了",
	"level.disaster.description":          "写的时候爽，改的时候哭",
	"level.disaster.severe.description":   "毒气弥漫，建议戴防毒面具",
	"level.disaster.very_bad.description": "进去的程序员没有一个活着出来",
	"level.disaster.extreme.description":  "反人类罪行，建议火化",
	"level.disaster.worst.description":    "历代工程师共创的遗产，无法维护",
	"level.disaster.ultimate.description": "写的时候热血澎湃，改的时候亲妈不认",

	// 总体评分评价
	"score.comment.0":  "如沐春风，代码仿佛被天使亲吻过",
	"score.comment.10": "清新宜人，闻起来像早晨的露珠",
	"score.comment.20": "略带清香，偶尔飘过一丝酸爽",
	"score.comment.30": "有点臭味，但还不至于熏死人",
	"score.comment.40": "屎气扑鼻，建议戴口罩阅读",
	"score.comment.50": "毒气缭绕，代码审查犹如酷刑",
	"score.comment.60": "熏天臭气，维护者已开始咳血",
	"score.comment.70": "生化危机，接手前请立遗嘱",
	"score.comment.80": "核废料现场，需穿防护服维护",
	"score.comment.90": "宇宙级毒瘤，看一眼寿命减十年",

	// 错误消息
	"error.path_not_accessible":    "无法访问路径: %v",
	"error.file_read_failed":       "读取文件 %s 失败: %v",
	"error.code_parse_failed":      "解析代码 %s 失败: %v",
	"error.source_files_not_found": "查找源文件失败: %v",
	"error.file_analysis_failed":   "分析文件 %s 失败: %v",

	// 警告和提示
	"warning.format": "警告: %v\n",

	// 函数复杂度问题
	"issue.high_complexity":        "函数 %s 的循环复杂度过高 (%d)，考虑重构",
	"issue.medium_complexity":      "函数 %s 的循环复杂度较高 (%d)，建议简化",
	"issue.file_high_complexity":   "文件循环复杂度过高 (%d)，建议拆分为多个文件",
	"issue.file_medium_complexity": "文件循环复杂度较高 (%d)，建议优化",

	// 函数长度问题
	"issue.function_very_long": "函数 %s 代码行数过多 (%d 行)，极度建议拆分",
	"issue.function_long":      "函数 %s 代码行数较多 (%d 行)，建议拆分为多个小函数",
	"issue.function_medium":    "函数 %s 长度为 %d 行，考虑是否可以简化",
	"issue.file_very_long":     "文件代码行数过多 (%d 行)，建议拆分为多个文件",
	"issue.file_long":          "文件代码行数较多 (%d 行)，考虑是否可以优化结构",

	// 注释覆盖率问题
	"issue.comment_very_low":         "代码注释率极低 (%.2f%%)，几乎没有注释",
	"issue.comment_low":              "代码注释率较低 (%.2f%%)，建议增加注释",
	"issue.exported_func_no_comment": "导出函数 %s 缺少文档注释",
	"issue.exported_type_no_comment": "导出类型 %s 缺少文档注释",

	// 详细报告
	"verbose.basic_statistics":  "📊 基本统计:",
	"verbose.total_files":       "总文件数:",
	"verbose.total_lines":       "总代码行:",
	"verbose.total_issues":      "总问题数:",
	"verbose.metric_details":    "🔍 指标详细信息:",
	"verbose.weight":            "权重:",
	"verbose.description":       "描述:",
	"verbose.score":             "得分:",
	"verbose.all_files":         "全部代码文件分析",
	"verbose.no_files_found":    "🎉 没有找到需要分析的文件！",
	"verbose.file_good_quality": "代码质量良好，没有明显问题",
}

// 英文语言包
var enUSMessages = map[string]string{
	// 通用
	"app.name":        "Code Quality Detector",
	"app.description": "A tool designed to evaluate code quality and reveal potential issues",

	// 指标名称
	"metric.cyclomatic_complexity": "Cyclomatic Complexity",
	"metric.function_length":       "State Management",
	"metric.comment_ratio":         "Comment Ratio",
	"metric.error_handling":        "Error Handling",
	"metric.naming_convention":     "Naming Convention",
	"metric.code_duplication":      "Code Duplication",
	"metric.structure_analysis":    "Code Structure",

	// 分析器进度
	"analyzer.searching_files":   "Searching for source code files...",
	"analyzer.files_found":       "Files found",
	"analyzer.analyzing_files":   "Analyzing files...",
	"analyzer.analysis_complete": "Analysis complete",

	// 问题分类
	"report.no_issues":           "Congratulations! No problematic files found!",
	"issue.category.complexity":  "Complexity Issues",
	"issue.category.comment":     "Comment Issues",
	"issue.category.naming":      "Naming Issues",
	"issue.category.structure":   "Structure Issues",
	"issue.category.duplication": "Duplication Issues",
	"issue.category.error":       "Error Handling Issues",
	"issue.category.other":       "Other Issues",

	// 质量等级
	"level.clean":             "Excellent",
	"level.mild":              "Slightly Smelly",
	"level.moderate":          "Mildly Stinky",
	"level.bad":               "Code Odor",
	"level.terrible":          "Medium Mess",
	"level.disaster":          "Hidden Tumor",
	"level.disaster.severe":   "Severe Mess",
	"level.disaster.very_bad": "Code Graveyard",
	"level.disaster.extreme":  "Nuclear Disaster",
	"level.disaster.worst":    "Ancient Legacy",
	"level.disaster.ultimate": "Ultimate Mess King",

	// 命令行
	"cmd.short":                      "💻 fuck-u-code",
	"cmd.long":                       "🔍 Code Quality Detector - Objectively assess your code quality\n\nIt can analyze code quality, output scores, and help you find 💩 in your code. Suitable for:\n- Quality assessment before project refactoring\n- Team code review assistance tool\n- Learning programming best practices",
	"cmd.analyze":                    "Analyze code quality and output score",
	"cmd.analyze.long":               "Deeply analyze the codebase, detect various potential code issues, and output a quality report. When no path is specified, the current directory is analyzed.",
	"cmd.completion":                 "Generate the autocompletion script for the specified shell",
	"cmd.completion.long":            "Generate the autocompletion script for the specified shell, supporting bash, zsh, fish and PowerShell.",
	"cmd.completion.long_prefix":     "Generate the autocompletion script for fuck-u-code for the specified shell.",
	"cmd.completion.long_suffix":     "See each sub-command's help for details on how to use the generated script.",
	"cmd.completion.bash":            "Generate the autocompletion script for bash",
	"cmd.completion.zsh":             "Generate the autocompletion script for zsh",
	"cmd.completion.fish":            "Generate the autocompletion script for fish",
	"cmd.completion.powershell":      "Generate the autocompletion script for powershell",
	"cmd.completion.bash.long":       "Generate the autocompletion script for the bash shell",
	"cmd.completion.zsh.long":        "Generate the autocompletion script for the zsh shell",
	"cmd.completion.fish.long":       "Generate the autocompletion script for the fish shell",
	"cmd.completion.powershell.long": "Generate the autocompletion script for powershell",
	"cmd.help":                       "Help about any command",
	"cmd.help.long":                  "Help provides help for any command in the application.",
	"cmd.help_flag":                  "help for fuck-u-code",
	"cmd.no_descriptions":            "disable completion descriptions",
	"cmd.path_not_found":             "Path not accessible '%s': %v",
	"cmd.analysis_failed":            "Analysis failed: %v",
	"cmd.lang":                       "Specify output language (supported: zh-CN, en-US, default: zh-CN)",
	"cmd.verbose":                    "Show detailed analysis report",
	"cmd.top":                        "Show the number of files with the most issues (default 5)",
	"cmd.issues":                     "How many issues to show for each file (default 5)",
	"cmd.summary":                    "Show only conclusion, skip the process",
	"cmd.markdown":                   "Output streamlined Markdown format report, suitable for AI tool processing",
	"cmd.exclude":                    "Exclude file/directory patterns (can be used multiple times, common dependency directories are excluded by default)",
	"cmd.start_analyzing":            "Start analyzing: %s",
	"cmd.exclude_patterns":           "Excluding the following file/directory patterns:",

	// Cobra框架内部文本
	"cobra.available_commands": "Available Commands",
	"cobra.flags":              "Flags",
	"cobra.global_flags":       "Global Flags",
	"cobra.additional_help":    "Additional help topics",
	"cobra.use_help_cmd":       "Use",
	"cobra.for_more_info":      "for more information about a command",
	"cobra.usage":              "Usage",

	// 报告
	"report.title":                   "Code Quality Analysis Report",
	"report.overall_score":           "Overall Score: %.2f / 100",
	"report.level":                   "Quality Level: %s",
	"report.metrics_details":         "Metrics Details",
	"report.worst_files":             "Problem Files Ranking",
	"report.conclusion":              "Conclusion",
	"report.file_score":              "Issue Score: %.2f",
	"report.more_issues":             "...and %d more issues",
	"report.score_calc":              "Score Calculation: ",
	"report.overall_assessment":      "Overall Assessment",
	"report.quality_score":           "Quality Score",
	"report.quality_level":           "Quality Level",
	"report.analyzed_files":          "Analyzed Files",
	"report.total_lines":             "Total Lines",
	"report.quality_metrics":         "Quality Metrics",
	"report.metric":                  "Metric",
	"report.score":                   "Score",
	"report.weight":                  "Weight",
	"report.status":                  "Status",
	"report.problem_files":           "Problem Files",
	"report.issue_categories":        "Issue Categories",
	"report.main_issues":             "Main Issues",
	"report.and":                     "and",
	"report.more_issues_short":       "more issues",
	"report.improvement_suggestions": "Improvement Suggestions",

	// 指标评分后缀
	"metric.score.suffix": " pts",

	// 循环复杂度评价
	"metric.complexity.good":   "Clear structure, no unnecessary complexity, great!",
	"metric.complexity.medium": "Winding logic, like a maze for your brain",
	"metric.complexity.bad":    "Functions like labyrinths, maintenance like a dungeon raid",

	// 函数长度评价
	"metric.length.good":   "Clear state management, reasonable variable scope, predictable state",
	"metric.length.medium": "Average state management, some global state or unclear state changes",
	"metric.length.bad":    "Chaotic state management, excessive use of global variables, difficult to track state changes",

	// 注释覆盖率评价
	"metric.comment.good":   "Good comments, they'll help you survive",
	"metric.comment.medium": "Sparse comments, readers need imagination",
	"metric.comment.bad":    "No comments, understanding depends on luck",

	// 错误处理评价
	"metric.error.good":   "Errors are handled with care, code shows compassion",
	"metric.error.medium": "Error handling exists, but barely helps",
	"metric.error.bad":    "Errors ignored? Just like life's problems",

	// 命名规范评价
	"metric.naming.good":   "Clear naming, the light of programmer civilization",
	"metric.naming.medium": "Naming is okay, some guesswork needed",
	"metric.naming.bad":    "Variable names look like keyboard smashes: x, y, z, tmp, xxx",

	// 代码重复度评价
	"metric.duplication.good":   "Proper abstraction, satisfying for the OCD programmer",
	"metric.duplication.medium": "Some repetition, abstraction wouldn't hurt",
	"metric.duplication.bad":    "Copy-paste evidence everywhere, Ctrl+C/V medal earned",

	// 代码结构评价
	"metric.structure.good":   "Beautiful structure, easy to follow",
	"metric.structure.medium": "Structure is okay, but somewhat confusing",
	"metric.structure.bad":    "Nested like Russian dolls, dizzying to read",

	// 质量建议
	"advice.good":     "👍 Keep it up, you're a breath of fresh air in the coding world, a clean code champion",
	"advice.moderate": "🔧 Suggestion: This code is like a rebellious teenager, needs some discipline to reach its potential",
	"advice.bad":      "🧨 Suggestion: Delete the repo and run, or seal it away for future generations to deal with",

	// 改进建议优先级
	"advice.priority.high":   "High Priority",
	"advice.priority.medium": "Medium Priority",
	"advice.priority.low":    "Low Priority",

	// 良好代码的建议
	"advice.good.maintain": "Continue maintaining current code quality standards",
	"advice.good.optimize": "Consider further optimizing performance and readability",
	"advice.good.document": "Improve documentation and comments for better team collaboration",

	// 中等代码的建议
	"advice.moderate.refactor":    "Refactor overly complex functions and modules",
	"advice.moderate.complexity":  "Reduce cyclomatic complexity and simplify control flow",
	"advice.moderate.naming":      "Improve variable and function naming conventions",
	"advice.moderate.comments":    "Increase code comment coverage",
	"advice.moderate.duplication": "Eliminate duplicate code and extract common methods",
	"advice.moderate.structure":   "Optimize code structure and reduce nesting levels",
	"advice.moderate.style":       "Unify code style and formatting",

	// 较差代码的建议
	"advice.bad.urgent_refactor": "Urgently refactor overly long functions, follow single responsibility principle",
	"advice.bad.complexity":      "Significantly reduce cyclomatic complexity, break down complex logic",
	"advice.bad.error_handling":  "Add comprehensive error handling mechanisms",
	"advice.bad.naming":          "Comprehensively improve naming conventions, avoid meaningless variable names",
	"advice.bad.duplication":     "Thoroughly eliminate duplicate code, establish code reuse mechanisms",
	"advice.bad.comments":        "Significantly increase code comments to improve readability",
	"advice.bad.structure":       "Redesign code architecture and improve overall structure",
	"advice.bad.style":           "Establish and enforce strict coding standards",

	// 指标描述
	"metric.function_length.description":       "Evaluates state variable management; good state management improves code maintainability and predictability",
	"metric.comment_ratio.description":         "Evaluates code comment coverage; good comments improve code readability and maintainability",
	"metric.error_handling.description":        "Examines error handling practices; proper error handling improves code robustness",
	"metric.naming_convention.description":     "Checks naming standards; good naming improves code readability",
	"metric.code_duplication.description":      "Assesses code duplication; more duplication indicates need for abstraction and refactoring",
	"metric.structure_analysis.description":    "Analyzes code nesting depth and reference complexity to evaluate structural clarity",
	"metric.cyclomatic_complexity.description": "Measures control flow complexity; higher complexity makes code harder to understand and test",

	// 质量等级描述
	"level.clean.description":             "Clean code, a pleasure to read",
	"level.mild.description":              "Mostly fine, but with some issues",
	"level.moderate.description":          "Slight smell, needs some ventilation",
	"level.bad.description":               "Code starts to smell, maintain with caution",
	"level.terrible.description":          "Obvious smell, can't be saved by opening windows",
	"level.disaster.description":          "Fun to write, painful to modify",
	"level.disaster.severe.description":   "Toxic gas, gas mask recommended",
	"level.disaster.very_bad.description": "No programmer enters and leaves alive",
	"level.disaster.extreme.description":  "Crime against humanity, incineration advised",
	"level.disaster.worst.description":    "Legacy of generations of engineers, unmaintainable",
	"level.disaster.ultimate.description": "Exciting to write, disowned by your own mother",

	// 总体评分评价
	"score.comment.0":  "Like a spring breeze, code kissed by angels",
	"score.comment.10": "Fresh and pleasant, smells like morning dew",
	"score.comment.20": "Slightly fragrant, with occasional tang",
	"score.comment.30": "Slightly smelly, but not lethal",
	"score.comment.40": "Stench in the air, mask recommended for reading",
	"score.comment.50": "Toxic fumes, code review feels like torture",
	"score.comment.60": "Overwhelming stench, maintainers coughing blood",
	"score.comment.70": "Biohazard, write your will before taking over",
	"score.comment.80": "Nuclear waste site, protective suit required",
	"score.comment.90": "Cosmic level tumor, each glance reduces lifespan by ten years",

	// 错误消息
	"error.path_not_accessible":    "Cannot access path: %v",
	"error.file_read_failed":       "Failed to read file %s: %v",
	"error.code_parse_failed":      "Failed to parse code %s: %v",
	"error.source_files_not_found": "Failed to find source files: %v",
	"error.file_analysis_failed":   "Failed to analyze file %s: %v",

	// 警告和提示
	"warning.format": "Warning: %v\n",

	// 函数复杂度问题
	"issue.high_complexity":        "Function %s has very high cyclomatic complexity (%d), consider refactoring",
	"issue.medium_complexity":      "Function %s has high cyclomatic complexity (%d), consider simplifying",
	"issue.file_high_complexity":   "File has very high complexity (%d), consider splitting into multiple files",
	"issue.file_medium_complexity": "File has high complexity (%d), consider optimizing",

	// 函数长度问题
	"issue.function_very_long": "Function %s has too many lines of code (%d), strongly recommend splitting",
	"issue.function_long":      "Function %s has many lines of code (%d), consider splitting into smaller functions",
	"issue.function_medium":    "Function %s has %d lines of code, consider if it can be simplified",
	"issue.file_very_long":     "File has too many lines of code (%d), recommend splitting into multiple files",
	"issue.file_long":          "File has many lines of code (%d), consider optimizing the structure",

	// 注释覆盖率问题
	"issue.comment_very_low":         "Code comment ratio is extremely low (%.2f%%), almost no comments",
	"issue.comment_low":              "Code comment ratio is low (%.2f%%), consider adding more comments",
	"issue.exported_func_no_comment": "Exported function %s lacks documentation comment",
	"issue.exported_type_no_comment": "Exported type %s lacks documentation comment",

	// 详细报告
	"verbose.basic_statistics":  "📊 Basic Statistics:",
	"verbose.total_files":       "Total Files:",
	"verbose.total_lines":       "Total Lines:",
	"verbose.total_issues":      "Total Issues:",
	"verbose.metric_details":    "🔍 Metrics Details:",
	"verbose.weight":            "Weight:",
	"verbose.description":       "Description:",
	"verbose.score":             "Score:",
	"verbose.all_files":         "All Code Files Analysis",
	"verbose.no_files_found":    "🎉 No files found for analysis!",
	"verbose.file_good_quality": "Good code quality, no obvious issues",
}
