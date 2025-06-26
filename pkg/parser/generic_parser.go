// Package parser 提供多语言代码解析功能
package parser

import (
	"regexp"
	"strings"

	"github.com/Done-0/fuck-u-code/pkg/common"
)

// GenericParser 通用文本解析器，适用于任何语言
type GenericParser struct{}

// NewGenericParser 创建新的通用文本解析器
func NewGenericParser() Parser {
	return &GenericParser{}
}

// Parse 通过文本分析解析代码
func (p *GenericParser) Parse(filePath string, content []byte) (ParseResult, error) {
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	result := &BaseParseResult{
		Functions:    make([]Function, 0),
		CommentLines: 0,
		TotalLines:   len(lines),
		Language:     common.Unsupported,
	}

	// 检测语言类型
	detector := common.NewLanguageDetector()
	result.Language = detector.DetectLanguage(filePath)

	// 计算注释行数
	result.CommentLines = p.countCommentLines(contentStr, result.Language)

	// 尝试识别函数
	result.Functions = p.detectFunctions(contentStr, lines, result.Language)

	return result, nil
}

// SupportedLanguages 返回支持的语言类型
func (p *GenericParser) SupportedLanguages() []common.LanguageType {
	return []common.LanguageType{
		common.JavaScript,
		common.TypeScript,
		common.Python,
		common.Java,
		common.CPlusPlus,
		common.C,
		common.Unsupported,
	}
}

// countCommentLines 根据语言类型计算注释行数
func (p *GenericParser) countCommentLines(content string, language common.LanguageType) int {
	commentCount := 0
	lines := strings.Split(content, "\n")

	switch language {
	case common.JavaScript, common.TypeScript, common.Java, common.CPlusPlus, common.C:
		// C风格注释处理
		commentCount = p.countCStyleComments(lines)
	case common.Python:
		// Python风格注释处理
		commentCount = p.countPythonComments(lines)
	default:
		// 通用注释处理
		commentCount = p.countGenericComments(lines)
	}

	return commentCount
}

// countCStyleComments 计算C风格注释行数
func (p *GenericParser) countCStyleComments(lines []string) int {
	commentCount := 0
	inBlockComment := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if inBlockComment {
			commentCount++
			if strings.Contains(trimmedLine, "*/") {
				inBlockComment = false
			}
			continue
		}

		if strings.HasPrefix(trimmedLine, "//") {
			commentCount++
			continue
		}

		if strings.HasPrefix(trimmedLine, "/*") {
			commentCount++
			inBlockComment = true
			if strings.Contains(trimmedLine, "*/") {
				inBlockComment = false
			}
			continue
		}
	}

	return commentCount
}

// countPythonComments 计算Python风格注释行数
func (p *GenericParser) countPythonComments(lines []string) int {
	commentCount := 0
	inDocString := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if inDocString {
			commentCount++
			if strings.Contains(trimmedLine, "\"\"\"") || strings.Contains(trimmedLine, "'''") {
				inDocString = false
			}
			continue
		}

		if strings.HasPrefix(trimmedLine, "#") {
			commentCount++
			continue
		}

		if strings.HasPrefix(trimmedLine, "\"\"\"") || strings.HasPrefix(trimmedLine, "'''") {
			commentCount++
			inDocString = true
			if (strings.Count(trimmedLine, "\"\"\"") > 1) || (strings.Count(trimmedLine, "'''") > 1) {
				inDocString = false
			}
			continue
		}
	}

	return commentCount
}

// countGenericComments 计算通用注释行数
func (p *GenericParser) countGenericComments(lines []string) int {
	commentCount := 0
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "//") ||
			strings.HasPrefix(trimmedLine, "#") ||
			strings.HasPrefix(trimmedLine, "/*") ||
			strings.HasPrefix(trimmedLine, "*") {
			commentCount++
		}
	}
	return commentCount
}

// detectFunctions 根据语言类型检测函数
func (p *GenericParser) detectFunctions(content string, lines []string, language common.LanguageType) []Function {
	switch language {
	case common.JavaScript, common.TypeScript:
		return p.detectJSFunctions(content, lines)
	case common.Python:
		return p.detectPythonFunctions(content, lines)
	case common.Java:
		return p.detectJavaFunctions(content, lines)
	case common.CPlusPlus, common.C:
		return p.detectCFunctions(content, lines)
	default:
		return p.detectGenericFunctions(content, lines)
	}
}

// 语言特定的函数模式
var (
	jsPattern      = regexp.MustCompile(`(?m)(function\s+([a-zA-Z_$][a-zA-Z0-9_$]*)|([a-zA-Z_$][a-zA-Z0-9_$]*)\s*=\s*function|([a-zA-Z_$][a-zA-Z0-9_$]*)\s*:\s*function|(?:const|let|var)\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*=\s*\([^)]*\)\s*=>)`)
	pythonPattern  = regexp.MustCompile(`(?m)^\s*def\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(([^)]*)\)`)
	javaPattern    = regexp.MustCompile(`(?m)(public|private|protected|static|\s)+[\w\<\>\[\]]+\s+([\w]+)\s*\(([^\)]*)\)\s*(\{|throws)`)
	cPattern       = regexp.MustCompile(`(?m)([\w\*]+\s+)+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(([^;]*)\)\s*\{`)
	genericPattern = regexp.MustCompile(`(?m)(function|def|void|int|bool|string|double|float)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`)
)

// detectJSFunctions 检测JavaScript/TypeScript函数
func (p *GenericParser) detectJSFunctions(content string, lines []string) []Function {
	return p.detectFunctionsWithPattern(content, lines, jsPattern, common.JavaScript)
}

// detectPythonFunctions 检测Python函数
func (p *GenericParser) detectPythonFunctions(content string, lines []string) []Function {
	functions := make([]Function, 0)
	matches := pythonPattern.FindAllStringSubmatchIndex(content, -1)

	for _, match := range matches {
		startPos := match[0]
		startLine := p.getLineNumber(content, startPos)

		// 提取函数名和参数
		submatch := pythonPattern.FindStringSubmatch(content[startPos : startPos+match[1]-match[0]])
		funcName := submatch[1]
		paramStr := submatch[2]

		// 查找函数结束行（Python特有的缩进检测）
		endLine := p.findPythonBlockEnd(lines, startLine-1)

		// 计算复杂度
		complexity := p.estimateComplexity(content, startPos, endLine-startLine, common.Python)

		// 计算参数数量
		params := 0
		if len(strings.TrimSpace(paramStr)) > 0 {
			params = strings.Count(paramStr, ",") + 1
		}

		functions = append(functions, Function{
			Name:       funcName,
			StartLine:  startLine,
			EndLine:    endLine,
			Complexity: complexity,
			Parameters: params,
		})
	}

	return functions
}

// detectJavaFunctions 检测Java方法
func (p *GenericParser) detectJavaFunctions(content string, lines []string) []Function {
	return p.detectFunctionsWithPattern(content, lines, javaPattern, common.Java)
}

// detectCFunctions 检测C/C++函数
func (p *GenericParser) detectCFunctions(content string, lines []string) []Function {
	return p.detectFunctionsWithPattern(content, lines, cPattern, common.C)
}

// detectGenericFunctions 通用函数检测（用于未知语言）
func (p *GenericParser) detectGenericFunctions(content string, lines []string) []Function {
	return p.detectFunctionsWithPattern(content, lines, genericPattern, common.Unsupported)
}

// detectFunctionsWithPattern 使用正则表达式检测函数（通用方法）
func (p *GenericParser) detectFunctionsWithPattern(content string, lines []string, pattern *regexp.Regexp, language common.LanguageType) []Function {
	functions := make([]Function, 0)
	matches := pattern.FindAllStringSubmatchIndex(content, -1)

	for _, match := range matches {
		startPos := match[0]
		startLine := p.getLineNumber(content, startPos)

		// 提取函数名（根据不同语言，提取位置可能不同）
		submatch := pattern.FindStringSubmatch(content[startPos : startPos+match[1]-match[0]])
		funcName := p.extractFunctionName(submatch, language)

		// 查找函数结束位置
		endLine := p.findBlockEnd(content, startPos, lines, startLine)

		// 计算参数数量
		params := p.countParameters(content[startPos:startPos+500], "(", ")")

		// 计算复杂度
		complexity := p.estimateComplexity(content, startPos, endLine-startLine, language)

		functions = append(functions, Function{
			Name:       funcName,
			StartLine:  startLine,
			EndLine:    endLine,
			Complexity: complexity,
			Parameters: params,
		})
	}

	return functions
}

// extractFunctionName 从正则匹配结果中提取函数名
func (p *GenericParser) extractFunctionName(submatch []string, language common.LanguageType) string {
	switch language {
	case common.JavaScript, common.TypeScript:
		// JS/TS函数名在不同位置
		for i := 2; i < len(submatch); i++ {
			if submatch[i] != "" {
				return submatch[i]
			}
		}
	case common.Java, common.C, common.CPlusPlus:
		// Java和C/C++函数名在第2位
		if len(submatch) > 2 {
			return submatch[2]
		}
	case common.Python, common.Unsupported:
		// Python和通用的函数名在第2位
		if len(submatch) > 2 {
			return submatch[2]
		}
	}
	return "anonymous"
}

// getLineNumber 计算字符串位置对应的行号
func (p *GenericParser) getLineNumber(content string, pos int) int {
	line := 1
	for i := 0; i < pos && i < len(content); i++ {
		if content[i] == '\n' {
			line++
		}
	}
	return line
}

// findBlockEnd 查找代码块的结束位置
func (p *GenericParser) findBlockEnd(content string, startPos int, lines []string, startLine int) int {
	bracketCount := 0
	inString := false
	inChar := false
	inLineComment := false
	inBlockComment := false

	for i := startPos; i < len(content); i++ {
		c := content[i]

		// 跳过字符串和字符字面量
		if !inLineComment && !inBlockComment {
			if c == '"' && (i == 0 || content[i-1] != '\\') {
				inString = !inString
				continue
			}
			if c == '\'' && (i == 0 || content[i-1] != '\\') {
				inChar = !inChar
				continue
			}
		}

		// 跳过注释
		if !inString && !inChar {
			if !inLineComment && !inBlockComment && c == '/' && i+1 < len(content) && content[i+1] == '/' {
				inLineComment = true
				continue
			}
			if !inLineComment && !inBlockComment && c == '/' && i+1 < len(content) && content[i+1] == '*' {
				inBlockComment = true
				i++
				continue
			}
			if inLineComment && c == '\n' {
				inLineComment = false
				continue
			}
			if inBlockComment && c == '*' && i+1 < len(content) && content[i+1] == '/' {
				inBlockComment = false
				i++
				continue
			}
		}

		// 统计括号
		if !inString && !inChar && !inLineComment && !inBlockComment {
			switch c {
			case '{':
				bracketCount++
			case '}':
				bracketCount--
				if bracketCount == 0 {
					// 计算结束行号
					endLine := startLine
					for j := startPos; j <= i; j++ {
						if content[j] == '\n' {
							endLine++
						}
					}
					return endLine
				}
			}
		}

		// 换行重置行注释状态
		if c == '\n' {
			inLineComment = false
		}
	}

	// 如果无法确定结束位置，返回文件末尾
	return len(lines)
}

// findPythonBlockEnd 查找Python函数结束位置
func (p *GenericParser) findPythonBlockEnd(lines []string, startLine int) int {
	if startLine >= len(lines) {
		return len(lines)
	}

	// 获取函数定义行的缩进级别
	defLine := lines[startLine]
	indentLevel := p.getIndentLevel(defLine)

	// 查找下一个缩进级别小于等于当前函数的非空行
	for i := startLine + 1; i < len(lines); i++ {
		if i >= len(lines) {
			return len(lines)
		}

		line := lines[i]
		trimmedLine := strings.TrimSpace(line)
		if len(trimmedLine) == 0 {
			continue
		}

		// 跳过注释行
		if strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		lineIndent := p.getIndentLevel(line)
		if lineIndent <= indentLevel {
			return i
		}
	}

	return len(lines)
}

// getIndentLevel 计算行的缩进级别
func (p *GenericParser) getIndentLevel(line string) int {
	indentLevel := 0
	for _, c := range line {
		switch c {
		case ' ':
			indentLevel++
		case '\t':
			indentLevel += 4 // 假设一个tab等于4个空格
		default:
			return indentLevel
		}
	}
	return indentLevel
}

// countParameters 从字符串中计算参数数量
func (p *GenericParser) countParameters(content string, openChar, closeChar string) int {
	start := strings.Index(content, openChar)
	if start == -1 {
		return 0
	}

	end := strings.Index(content[start:], closeChar)
	if end == -1 {
		return 0
	}

	params := content[start+1 : start+end]
	params = strings.TrimSpace(params)

	if params == "" || params == "void" {
		return 0
	}

	return strings.Count(params, ",") + 1
}

// estimateComplexity 根据语言估算函数复杂度
func (p *GenericParser) estimateComplexity(content string, startPos, lineCount int, language common.LanguageType) int {
	// 基础复杂度为1
	complexity := 1

	// 根据语言选择关键字
	var keywords []string
	switch language {
	case common.JavaScript, common.TypeScript:
		keywords = []string{"if", "else", "for", "while", "case", "catch", "&&", "||", "?", "switch"}
	case common.Python:
		keywords = []string{"if", "elif", "else", "for", "while", "except", "finally", "with", "and", "or"}
	case common.Java:
		keywords = []string{"if", "else", "for", "while", "do", "case", "catch", "finally", "?", "&&", "||", "switch"}
	case common.C, common.CPlusPlus:
		keywords = []string{"if", "else", "for", "while", "do", "case", "switch", "catch", "?", "&&", "||", "goto"}
	default:
		keywords = []string{"if", "else", "for", "while", "switch", "case", "try", "catch", "&&", "||"}
	}

	// 提取函数内容
	endPos := p.findContentEndPosition(content, startPos, lineCount)
	funcContent := content[startPos:endPos]

	// 计算关键字出现次数
	for _, keyword := range keywords {
		re := regexp.MustCompile(`\b` + keyword + `\b`)
		matches := re.FindAllStringIndex(funcContent, -1)
		if matches != nil {
			complexity += len(matches)
		}
	}

	return complexity
}

// findContentEndPosition 找到内容的结束位置
func (p *GenericParser) findContentEndPosition(content string, startPos, lineCount int) int {
	endPos := startPos
	lineCounter := 0

	for i := startPos; i < len(content) && lineCounter < lineCount; i++ {
		if content[i] == '\n' {
			lineCounter++
		}
		endPos = i
	}

	return endPos
}
