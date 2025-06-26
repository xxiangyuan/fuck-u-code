// Package parser 提供多语言代码解析功能
package parser

import (
	"context"
	"strings"

	"github.com/Done-0/fuck-u-code/pkg/common"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/c"
	"github.com/smacker/go-tree-sitter/cpp"
)

// CParser C/C++语言解析器
type CParser struct{}

// NewCParser 创建新的C/C++语言解析器
func NewCParser() Parser {
	return &CParser{}
}

// Parse 解析C/C++代码
func (p *CParser) Parse(filePath string, content []byte) (ParseResult, error) {
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	// 检测语言类型
	detector := common.NewLanguageDetector()
	language := detector.DetectLanguage(filePath)

	result := &BaseParseResult{
		Functions:    make([]Function, 0),
		CommentLines: 0,
		TotalLines:   len(lines),
		Language:     language,
	}

	// 计算注释行数
	result.CommentLines = p.countCommentLines(contentStr)

	// 创建适合的解析器
	var parser *sitter.Parser
	var treeSitterLang *sitter.Language

	if language == common.CPlusPlus {
		treeSitterLang = cpp.GetLanguage()
	} else {
		treeSitterLang = c.GetLanguage()
	}

	parser = sitter.NewParser()
	parser.SetLanguage(treeSitterLang)

	// 解析代码
	tree, err := parser.ParseCtx(context.TODO(), nil, []byte(contentStr))
	if err != nil {
		result.Functions = p.detectFunctions(contentStr, lines)
		return result, nil
	}

	// 保存AST根节点
	result.ASTRoot = tree.RootNode()

	// 提取函数信息
	p.extractFunctions(tree.RootNode(), contentStr, &result.Functions)

	return result, nil
}

// SupportedLanguages 返回支持的语言类型
func (p *CParser) SupportedLanguages() []common.LanguageType {
	return []common.LanguageType{common.C, common.CPlusPlus}
}

// countCommentLines 计算C/C++代码中的注释行数
func (p *CParser) countCommentLines(content string) int {
	commentCount := 0
	lines := strings.Split(content, "\n")

	// 处理 // 和 /* */ 注释
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

// detectFunctions 基于文本分析的C/C++函数检测
func (p *CParser) detectFunctions(_ string, lines []string) []Function {
	functions := make([]Function, 0)

	// 简单的函数定义模式检测
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasSuffix(trimmed, "{") && !strings.HasPrefix(trimmed, "{") {
			// 可能是函数定义
			name := p.extractFunctionName(trimmed)
			if name != "" {
				params := p.countParameters(trimmed)

				// 寻找函数结束位置
				endLine := p.findFunctionEnd(lines, i)

				// 计算复杂度
				complexity := p.estimateComplexity(lines, i, endLine-i)

				function := Function{
					Name:       name,
					StartLine:  i + 1,
					EndLine:    endLine,
					Complexity: complexity,
					Parameters: params,
				}

				functions = append(functions, function)
			}
		}
	}

	return functions
}

// extractFunctionName 从函数定义行提取函数名
func (p *CParser) extractFunctionName(line string) string {
	// 去掉结尾的 {
	line = strings.TrimSuffix(line, "{")
	line = strings.TrimSpace(line)

	// 去掉参数部分
	if idx := strings.LastIndex(line, "("); idx != -1 {
		line = line[:idx]
		line = strings.TrimSpace(line)

		// 提取最后一个单词，即函数名
		parts := strings.Fields(line)
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}

	return ""
}

// countParameters 计算函数参数数量
func (p *CParser) countParameters(line string) int {
	// 提取括号内的内容
	start := strings.LastIndex(line, "(")
	end := strings.LastIndex(line, ")")

	if start != -1 && end != -1 && end > start {
		params := line[start+1 : end]
		params = strings.TrimSpace(params)

		if params == "" || params == "void" {
			return 0
		}

		return strings.Count(params, ",") + 1
	}

	return 0
}

// findFunctionEnd 查找函数结束位置
func (p *CParser) findFunctionEnd(lines []string, startLine int) int {
	bracketCount := 1

	for i := startLine + 1; i < len(lines); i++ {
		line := lines[i]

		for _, char := range line {
			switch char {
			case '{':
				bracketCount++
			case '}':
				bracketCount--
				if bracketCount == 0 {
					return i + 1
				}
			}
		}
	}

	return len(lines)
}

// estimateComplexity 估算函数复杂度
func (p *CParser) estimateComplexity(lines []string, startLine, lineCount int) int {
	complexity := 1

	// 定义复杂度关键词
	keywords := []string{"if", "else", "for", "while", "do", "switch", "case", "catch", "?", "&&", "||"}

	for i := startLine; i < startLine+lineCount && i < len(lines); i++ {
		line := lines[i]

		for _, keyword := range keywords {
			// 简单统计关键词出现次数
			count := strings.Count(line, keyword)

			// 确保是独立的关键词，而不是变量名的一部分
			for j := 0; j < count; j++ {
				pos := strings.Index(line, keyword)
				if pos != -1 {
					// 检查前后是否是标识符的一部分
					if (pos == 0 || !isAlphaNum(rune(line[pos-1]))) &&
						(pos+len(keyword) >= len(line) || !isAlphaNum(rune(line[pos+len(keyword)]))) {
						complexity++
					}
					line = line[pos+len(keyword):]
				}
			}
		}
	}

	return complexity
}

// isAlphaNum 检查字符是否是字母、数字或下划线
func isAlphaNum(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}

// extractFunctions 从AST中提取函数信息
func (p *CParser) extractFunctions(rootNode *sitter.Node, content string, functions *[]Function) {
	if rootNode == nil {
		return
	}

	// 遍历AST查找函数定义
	cursor := sitter.NewTreeCursor(rootNode)
	defer cursor.Close()

	if cursor.GoToFirstChild() {
		for {
			node := cursor.CurrentNode()
			if node.Type() == "function_definition" || node.Type() == "method_definition" {
				function := p.processFunctionNode(node, content)
				*functions = append(*functions, function)
			}

			if !cursor.GoToNextSibling() {
				break
			}
		}
		cursor.GoToParent()
	}
}

// processFunctionNode 处理函数节点
func (p *CParser) processFunctionNode(node *sitter.Node, content string) Function {
	// 提取函数名
	var funcName string
	declarator := node.ChildByFieldName("declarator")
	if declarator != nil {
		// 继续深入查找函数名
		funcName = p.findFunctionName(declarator, content)
	}

	// 获取行号范围
	startLine := int(node.StartPoint().Row) + 1
	endLine := int(node.EndPoint().Row) + 1

	// 计算参数数量
	params := p.countParametersFromAST(node)

	// 计算复杂度
	complexity := p.calculateComplexity(node)

	return Function{
		Name:       funcName,
		StartLine:  startLine,
		EndLine:    endLine,
		Complexity: complexity,
		Parameters: params,
		Node:       node,
	}
}

// findFunctionName 从AST中提取函数名
func (p *CParser) findFunctionName(node *sitter.Node, content string) string {
	if node == nil {
		return "unknown"
	}

	// 针对不同语言和节点类型的不同处理
	if node.Type() == "function_declarator" {
		declarator := node.ChildByFieldName("declarator")
		if declarator != nil {
			return p.findFunctionName(declarator, content)
		}
	} else if node.Type() == "identifier" {
		// 获取标识符文本
		return content[node.StartByte():node.EndByte()]
	}

	// 递归查找子节点
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "identifier" {
			return content[child.StartByte():child.EndByte()]
		}

		name := p.findFunctionName(child, content)
		if name != "unknown" {
			return name
		}
	}

	return "unknown"
}

// countParametersFromAST 从AST中计算函数参数数量
func (p *CParser) countParametersFromAST(node *sitter.Node) int {
	// 寻找参数列表节点
	var paramNode *sitter.Node

	declarator := node.ChildByFieldName("declarator")
	if declarator != nil && declarator.Type() == "function_declarator" {
		paramNode = declarator.ChildByFieldName("parameters")
	}

	if paramNode == nil {
		return 0
	}

	// 计算参数数量
	count := 0
	for i := 0; i < int(paramNode.ChildCount()); i++ {
		child := paramNode.Child(i)
		if child.Type() == "parameter_declaration" {
			count++
		}
	}

	return count
}

// calculateComplexity 通过AST计算函数复杂度
func (p *CParser) calculateComplexity(node *sitter.Node) int {
	complexity := 1

	// 使用队列进行广度优先搜索
	queue := []*sitter.Node{node}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// 检查节点类型，增加复杂度
		switch current.Type() {
		case "if_statement", "for_statement", "while_statement", "do_statement",
			"case_statement", "catch_statement", "conditional_expression":
			complexity++
		}

		// 将所有子节点加入队列
		for i := 0; i < int(current.ChildCount()); i++ {
			child := current.Child(i)
			queue = append(queue, child)
		}
	}

	return complexity
}
