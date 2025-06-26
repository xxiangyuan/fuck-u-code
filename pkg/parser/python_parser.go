// Package parser 提供多语言代码解析功能
package parser

import (
	"strings"

	"github.com/Done-0/fuck-u-code/pkg/common"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// PythonParser Python语言解析器
type PythonParser struct {
	parser *participle.Parser[PythonAST]
}

// NewPythonParser 创建新的Python语言解析器
func NewPythonParser() Parser {
	pythonLexer := lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Comment", Pattern: `#[^\n]*|"""(?:.|\n)*?"""|'''(?:.|\n)*?'''`},
		{Name: "Whitespace", Pattern: `\s+`},
		{Name: "Keyword", Pattern: `def|class|if|elif|else|for|while|try|except|finally|with|return|pass|break|continue|import|from|as|global|nonlocal|assert`},
		{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
		{Name: "Number", Pattern: `[0-9]+`},
		{Name: "String", Pattern: `"[^"]*"|'[^']*'`},
		{Name: "Punct", Pattern: `[\(\)\[\]\{\}:;,\.\+\-\*/%<>=!&\|\^~]`},
		{Name: "EOL", Pattern: `[\n\r]+`},
	})

	parser, _ := participle.Build[PythonAST](
		participle.Lexer(pythonLexer),
		participle.Elide("Whitespace"),
	)

	return &PythonParser{
		parser: parser,
	}
}

// Parse 解析Python代码
func (p *PythonParser) Parse(filePath string, content []byte) (ParseResult, error) {
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	result := &BaseParseResult{
		Functions:    make([]Function, 0),
		CommentLines: 0,
		TotalLines:   len(lines),
		Language:     common.Python,
	}

	// 计算注释行数
	result.CommentLines = p.countCommentLines(contentStr)

	// 解析AST
	ast, err := p.parser.ParseString("", contentStr)
	if err != nil {
		// 如果解析失败，使用备用的基于文本的方法
		functions := p.detectFunctions(contentStr, lines)
		result.Functions = functions
		return result, nil
	}

	// 保存AST根节点
	result.ASTRoot = ast

	// 提取函数信息
	p.extractFunctions(ast, contentStr, &result.Functions)

	return result, nil
}

// SupportedLanguages 返回支持的语言类型
func (p *PythonParser) SupportedLanguages() []common.LanguageType {
	return []common.LanguageType{common.Python}
}

// countCommentLines 计算Python代码中的注释行数
func (p *PythonParser) countCommentLines(content string) int {
	commentCount := 0
	lines := strings.Split(content, "\n")

	// 处理 # 和 """ """ 或 ''' ''' 注释
	inDocString := false
	docStringDelimiter := ""

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if inDocString {
			commentCount++
			// 检查文档字符串是否结束
			if strings.Contains(trimmedLine, docStringDelimiter) {
				inDocString = false
			}
			continue
		}

		// 检查行注释
		if strings.HasPrefix(trimmedLine, "#") {
			commentCount++
			continue
		}

		// 检查文档字符串开始
		if strings.HasPrefix(trimmedLine, "\"\"\"") {
			commentCount++
			inDocString = true
			docStringDelimiter = "\"\"\""
			// 检查单行文档字符串
			if strings.Count(trimmedLine, "\"\"\"") > 1 {
				inDocString = false
			}
			continue
		}

		if strings.HasPrefix(trimmedLine, "'''") {
			commentCount++
			inDocString = true
			docStringDelimiter = "'''"
			// 检查单行文档字符串
			if strings.Count(trimmedLine, "'''") > 1 {
				inDocString = false
			}
			continue
		}
	}

	return commentCount
}

// detectFunctions 基于文本的Python函数检测（备用方法）
func (p *PythonParser) detectFunctions(content string, lines []string) []Function {
	functions := make([]Function, 0)

	// 查找函数定义
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "def ") {
			// 提取函数名
			parts := strings.Split(trimmedLine, "(")
			if len(parts) > 0 {
				funcName := strings.TrimPrefix(parts[0], "def ")
				funcName = strings.TrimSpace(funcName)

				// 计算参数数量
				params := 0
				if len(parts) > 1 {
					paramPart := parts[1]
					if strings.Contains(paramPart, ")") {
						paramPart = strings.Split(paramPart, ")")[0]
					}
					paramPart = strings.TrimSpace(paramPart)
					if len(paramPart) > 0 {
						params = strings.Count(paramPart, ",") + 1
					}
				}

				// 计算函数结束行
				endLine := p.findFunctionEnd(lines, i)

				// 计算复杂度
				complexity := p.estimateComplexity(content, i, endLine-i)

				function := Function{
					Name:       funcName,
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

// findFunctionEnd 查找Python函数结束位置
func (p *PythonParser) findFunctionEnd(lines []string, startLine int) int {
	if startLine >= len(lines) {
		return len(lines)
	}

	// 获取函数定义行的缩进级别
	defLine := lines[startLine]
	indentLevel := 0
	for _, c := range defLine {
		switch c {
		case ' ':
			indentLevel++
		case '\t':
			indentLevel += 4
		default:
			goto BreakIndent
		}
	}
BreakIndent:

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

		lineIndent := 0
		for _, c := range line {
			switch c {
			case ' ':
				lineIndent++
			case '\t':
				lineIndent += 4
			default:
				goto BreakLoop
			}
		}
	BreakLoop:

		if lineIndent <= indentLevel {
			return i
		}
	}

	return len(lines)
}

// estimateComplexity 估算Python函数复杂度
func (p *PythonParser) estimateComplexity(content string, startLine int, lineCount int) int {
	// 基础复杂度为1
	complexity := 1

	// 从startLine开始提取函数内容
	lines := strings.Split(content, "\n")
	funcLines := lines[startLine : startLine+lineCount]
	funcContent := strings.Join(funcLines, "\n")

	// 统计复杂度关键字
	complexityKeywords := []string{"if", "elif", "else", "for", "while", "except", "finally", "with", "and", "or"}
	for _, keyword := range complexityKeywords {
		complexity += strings.Count(funcContent, " "+keyword+" ")
	}

	return complexity
}

// extractFunctions 从AST中提取函数信息
func (p *PythonParser) extractFunctions(ast *PythonAST, content string, functions *[]Function) {
	if ast == nil || len(ast.Statements) == 0 {
		return
	}

	lines := strings.Split(content, "\n")

	for _, stmt := range ast.Statements {
		if stmt.FunctionDef != nil {
			funcDef := stmt.FunctionDef

			// 提取函数名和参数
			name := funcDef.Name
			params := 0
			if funcDef.Parameters != nil && len(funcDef.Parameters.Params) > 0 {
				params = len(funcDef.Parameters.Params)
			}

			// 计算行号范围
			startLine := funcDef.Pos.Line
			endLine := p.findFunctionEnd(lines, startLine-1) // 转换为0-based

			// 计算复杂度
			complexity := 1
			p.calculateComplexity(funcDef.Body, &complexity)

			function := Function{
				Name:       name,
				StartLine:  startLine,
				EndLine:    endLine,
				Complexity: complexity,
				Parameters: params,
				Node:       funcDef,
			}

			*functions = append(*functions, function)
		}
	}
}

// calculateComplexity 递归计算AST节点的复杂度
func (p *PythonParser) calculateComplexity(stmts []PythonStatement, complexity *int) {
	for _, stmt := range stmts {
		if stmt.IfStmt != nil {
			*complexity++
			p.calculateComplexity(stmt.IfStmt.Body, complexity)
			for _, elifClause := range stmt.IfStmt.ElifClauses {
				*complexity++
				p.calculateComplexity(elifClause.Body, complexity)
			}
			if stmt.IfStmt.ElseClause != nil {
				p.calculateComplexity(stmt.IfStmt.ElseClause.Body, complexity)
			}
		} else if stmt.ForStmt != nil {
			*complexity++
			p.calculateComplexity(stmt.ForStmt.Body, complexity)
		} else if stmt.WhileStmt != nil {
			*complexity++
			p.calculateComplexity(stmt.WhileStmt.Body, complexity)
		} else if stmt.TryStmt != nil {
			p.calculateComplexity(stmt.TryStmt.Body, complexity)
			for _, exceptClause := range stmt.TryStmt.ExceptClauses {
				*complexity++
				p.calculateComplexity(exceptClause.Body, complexity)
			}
			if stmt.TryStmt.FinallyClause != nil {
				p.calculateComplexity(stmt.TryStmt.FinallyClause.Body, complexity)
			}
		}
	}
}

// Python AST结构定义
type PythonAST struct {
	Statements []PythonStatement `parser:"@@*"`
}

// PythonStatement 表示Python语句
type PythonStatement struct {
	FunctionDef *FunctionDef `parser:"@@"`
	IfStmt      *IfStmt      `parser:"| @@"`
	ForStmt     *ForStmt     `parser:"| @@"`
	WhileStmt   *WhileStmt   `parser:"| @@"`
	TryStmt     *TryStmt     `parser:"| @@"`
	Expression  *Expression  `parser:"| @@"`
	Comment     *string      `parser:"| @Comment"`
	EOL         *string      `parser:"| @EOL"`
}

// FunctionDef 表示函数定义
type FunctionDef struct {
	Pos        lexer.Position
	Name       string            `parser:"'def' @Ident"`
	Parameters *ParameterList    `parser:"'(' @@? ')'"`
	Body       []PythonStatement `parser:"':' @@ EOL*"`
}

// ParameterList 表示参数列表
type ParameterList struct {
	Params []Parameter `parser:"@@ (',' @@)*"`
}

// Parameter 表示函数参数
type Parameter struct {
	Name string `parser:"@Ident"`
}

// IfStmt 表示if语句
type IfStmt struct {
	Condition   *Expression       `parser:"'if' @@"`
	Body        []PythonStatement `parser:"':' @@ EOL*"`
	ElifClauses []ElifClause      `parser:"@@*"`
	ElseClause  *ElseClause       `parser:"@@?"`
}

// ElifClause 表示elif子句
type ElifClause struct {
	Condition *Expression       `parser:"'elif' @@"`
	Body      []PythonStatement `parser:"':' @@ EOL*"`
}

// ElseClause 表示else子句
type ElseClause struct {
	Body []PythonStatement `parser:"'else' ':' @@ EOL*"`
}

// ForStmt 表示for循环
type ForStmt struct {
	Target   string            `parser:"'for' @Ident"`
	Iterable *Expression       `parser:"'in' @@"`
	Body     []PythonStatement `parser:"':' @@ EOL*"`
}

// WhileStmt 表示while循环
type WhileStmt struct {
	Condition *Expression       `parser:"'while' @@"`
	Body      []PythonStatement `parser:"':' @@ EOL*"`
}

// TryStmt 表示try语句
type TryStmt struct {
	Body          []PythonStatement `parser:"'try' ':' @@ EOL*"`
	ExceptClauses []ExceptClause    `parser:"@@*"`
	FinallyClause *FinallyClause    `parser:"@@?"`
}

// ExceptClause 表示except子句
type ExceptClause struct {
	Body []PythonStatement `parser:"'except' ':' @@ EOL*"`
}

// FinallyClause 表示finally子句
type FinallyClause struct {
	Body []PythonStatement `parser:"'finally' ':' @@ EOL*"`
}

// Expression 表示表达式
type Expression struct {
	Ident  *string `parser:"@Ident"`
	Number *string `parser:"| @Number"`
	String *string `parser:"| @String"`
}
