// Package main 提供屎山代码检测器的命令行入口
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Done-0/fuck-u-code/pkg/analyzer"
	"github.com/Done-0/fuck-u-code/pkg/report"
)

var (
	verbose     bool // 是否输出屎味细节
	topFiles    int  // 最臭的文件数量
	maxIssues   int  // 每个文件最多列出的问题数
	summaryOnly bool // 是否只显示结论，不看过程
)

func main() {
	// 创建根命令
	rootCmd := &cobra.Command{
		Use:   "fuck-u-code [path]",
		Short: "💩 史上最毒的屎山分析器",
		Long: `🔥 fuck-u-code - 给你的代码打一针"自知之明"疫苗

它可以评估代码质量、输出屎山等级评分，带你直面自己写过的"屎味之作"。适用于：
- 遗留项目自我审判
- 团队重构前的"灾难认定"
- 离职前的"放毒"回顾

💀 友情提示：结果可能伤害心灵，谨慎使用。`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// 如果没有参数，显示帮助信息
			if len(args) == 0 && !cmd.Flags().Changed("help") {
				return cmd.Help()
			}
			// 有参数时正常执行分析
			runAnalysis(cmd, args)
			return nil
		},
	}

	// 添加命令行参数
	addFlags(rootCmd)

	// 创建分析子命令
	analyzeCmd := &cobra.Command{
		Use:   "analyze [path]",
		Short: "分析代码质量并输出屎山指数",
		Long:  "深入分析代码库，检测各种代码坏味道，输出屎山报告。不指定路径时分析当前目录。",
		Args:  cobra.MaximumNArgs(1),
		Run:   runAnalysis,
	}

	// 为子命令添加相同的参数
	addFlags(analyzeCmd)

	// 添加到根命令
	rootCmd.AddCommand(analyzeCmd)

	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// addFlags 添加命令行参数
func addFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "显示详细屎山成分分析报告")
	cmd.Flags().IntVarP(&topFiles, "top", "t", 5, "显示最臭的文件数量（默认5个）")
	cmd.Flags().IntVarP(&maxIssues, "issues", "i", 5, "每个文件显示多少条臭味来源（默认5个）")
	cmd.Flags().BoolVarP(&summaryOnly, "summary", "s", false, "只看屎山结论，过程略过")
}

// runAnalysis 执行代码分析
func runAnalysis(cmd *cobra.Command, args []string) {
	// 读取路径（默认为当前目录）
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// 解析成绝对路径
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("🚫 无法解析路径：%v\n", err)
		os.Exit(1)
	}

	// 确认路径存在
	info, err := os.Stat(absPath)
	if err != nil {
		fmt.Printf("🚫 路径不可访问 '%s': %v\n", absPath, err)
		os.Exit(1)
	}

	fmt.Printf("🔍 开始嗅探：%s\n", absPath)

	// 创建分析器
	codeAnalyzer := analyzer.NewAnalyzer()
	var result *analyzer.AnalysisResult

	// 判断是文件还是目录
	if info.IsDir() {
		result, err = codeAnalyzer.Analyze(absPath)
	} else {
		result, err = codeAnalyzer.AnalyzeFile(absPath)
	}

	if err != nil {
		fmt.Printf("💥 分析失败：%v\n", err)
		os.Exit(1)
	}

	// 生成报告
	reportGen := report.NewReport(result)
	reportOptions := &report.ReportOptions{
		Verbose:     verbose,
		TopFiles:    topFiles,
		MaxIssues:   maxIssues,
		SummaryOnly: summaryOnly,
	}

	// 输出报告
	reportGen.GenerateConsoleReport(reportOptions)
}
