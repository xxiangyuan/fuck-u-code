// Package parser 提供多语言代码解析功能
package parser

import (
    "regexp"
    "strings"

    "github.com/Done-0/fuck-u-code/pkg/common"
)

// CSharpParser C#语言解析器
type CSharpParser struct{}

// NewCSharpParser 创建新的C#语言解析器
func NewCSharpParser() Parser {
    return &CSharpParser{}
}

// Parse 解析C#代码
func (p *CSharpParser) Parse(filePath string, content []byte) (ParseResult, error) {
    contentStr := string(content)
    lines := strings.Split(contentStr, "\n")

    // 检测是否为Razor文件
    isRazor := strings.HasSuffix(strings.ToLower(filePath), ".razor") || 
               strings.HasSuffix(strings.ToLower(filePath), ".cshtml")

    result := &BaseParseResult{
        Functions:    make([]Function, 0),
        CommentLines: 0,
        TotalLines:   len(lines),
        Language:     common.CSharp,
    }

    if isRazor {
        // Razor文件处理
        result.CommentLines = p.countRazorCommentLines(contentStr)
        result.Functions = p.detectRazorMethods(contentStr, lines)
    } else {
        // 普通C#文件处理
        result.CommentLines = p.countCommentLines(contentStr)
        result.Functions = p.detectCSharpMethods(contentStr, lines)
    }

    return result, nil
}

// SupportedLanguages 返回支持的语言类型
func (p *CSharpParser) SupportedLanguages() []common.LanguageType {
    return []common.LanguageType{common.CSharp}
}

// countCommentLines 计算C#代码中的注释行数
func (p *CSharpParser) countCommentLines(content string) int {
    commentCount := 0
    lines := strings.Split(content, "\n")

    inBlockComment := false
    inXmlDoc := false

    for _, line := range lines {
        trimmedLine := strings.TrimSpace(line)

        if inBlockComment {
            commentCount++
            if strings.Contains(trimmedLine, "*/") {
                inBlockComment = false
            }
            continue
        }

        if inXmlDoc {
            commentCount++
            if !strings.HasPrefix(trimmedLine, "///") && 
               !strings.HasPrefix(trimmedLine, "*") && 
               len(trimmedLine) > 0 {
                inXmlDoc = false
            }
            continue
        }

        // 单行注释
        if strings.HasPrefix(trimmedLine, "//") {
            commentCount++
            continue
        }

        // XML文档注释
        if strings.HasPrefix(trimmedLine, "///") {
            commentCount++
            inXmlDoc = true
            continue
        }

        // 块注释开始
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

// countRazorCommentLines 计算Razor文件中的注释行数
func (p *CSharpParser) countRazorCommentLines(content string) int {
    commentCount := 0
    lines := strings.Split(content, "\n")

    inBlockComment := false
    inHtmlComment := false
    inRazorComment := false
    inCodeBlock := false

    for _, line := range lines {
        trimmedLine := strings.TrimSpace(line)

        // 检测是否在@code或@functions块中
        if strings.Contains(trimmedLine, "@code") || strings.Contains(trimmedLine, "@functions") {
            inCodeBlock = true
        }
        if inCodeBlock && strings.Contains(trimmedLine, "}") && 
           !strings.Contains(trimmedLine, "{") {
            inCodeBlock = false
        }

        // Razor注释 @* ... *@
        if inRazorComment {
            commentCount++
            if strings.Contains(trimmedLine, "*@") {
                inRazorComment = false
            }
            continue
        }

        // HTML注释 <!-- ... -->
        if inHtmlComment {
            commentCount++
            if strings.Contains(trimmedLine, "-->") {
                inHtmlComment = false
            }
            continue
        }

        // C#块注释 /* ... */
        if inBlockComment {
            commentCount++
            if strings.Contains(trimmedLine, "*/") {
                inBlockComment = false
            }
            continue
        }

        // 检测注释开始
        if strings.HasPrefix(trimmedLine, "@*") {
            commentCount++
            inRazorComment = true
            if strings.Contains(trimmedLine, "*@") {
                inRazorComment = false
            }
            continue
        }

        if strings.HasPrefix(trimmedLine, "<!--") {
            commentCount++
            inHtmlComment = true
            if strings.Contains(trimmedLine, "-->") {
                inHtmlComment = false
            }
            continue
        }

        // 在代码块中的C#注释
        if inCodeBlock {
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
            if strings.HasPrefix(trimmedLine, "///") {
                commentCount++
                continue
            }
        }
    }

    return commentCount
}

// detectRazorMethods 检测Razor文件中的方法
func (p *CSharpParser) detectRazorMethods(content string, lines []string) []Function {
    functions := make([]Function, 0)

    // 提取Razor文件中的C#代码块
    codeBlocks := p.extractRazorCodeBlocks(content)
    
    for _, block := range codeBlocks {
        // 在每个代码块中检测方法
        blockLines := strings.Split(block.Content, "\n")
        blockFunctions := p.detectCSharpMethods(block.Content, blockLines)
        
        // 调整行号偏移
        for i := range blockFunctions {
            blockFunctions[i].StartLine += block.StartLine
            blockFunctions[i].EndLine += block.StartLine
        }
        
        functions = append(functions, blockFunctions...)
    }

    // 检测页面级别
    pageDirectives := p.detectRazorPageDirectives(content, lines)
    functions = append(functions, pageDirectives...)

    return functions
}

// RazorCodeBlock 表示Razor文件中的C#代码块
type RazorCodeBlock struct {
    Content   string
    StartLine int
    EndLine   int
    BlockType string // "code", "functions", "inline"
}

// extractRazorCodeBlocks 提取Razor文件中的C#代码块
func (p *CSharpParser) extractRazorCodeBlocks(content string) []RazorCodeBlock {
    var blocks []RazorCodeBlock
    lines := strings.Split(content, "\n")

    // 匹配@code和@functions块
    codeBlockPattern := regexp.MustCompile(`@(code|functions)\s*\{`)

    // 处理@code和@functions块
    for i, line := range lines {
        if matches := codeBlockPattern.FindStringSubmatch(line); matches != nil {
            blockType := matches[1]
            startLine := i
            bracketCount := 0
            
            // 查找第一个大括号
            if openBraceIndex := strings.Index(line, "{"); openBraceIndex != -1 {
                bracketCount = 1
            }
            
            var blockContent strings.Builder
            blockContent.WriteString(line)
            blockContent.WriteString("\n")

            // 查找块结束
            for j := i + 1; j < len(lines) && bracketCount > 0; j++ {
                currentLine := lines[j]
                blockContent.WriteString(currentLine)
                blockContent.WriteString("\n")
                
                // 大括号计数，不考虑字符串内的大括号
                for _, char := range currentLine {
                    switch char {
                    case '{':
                        bracketCount++
                    case '}':
                        bracketCount--
                    }
                }
                
                if bracketCount == 0 {
                    blocks = append(blocks, RazorCodeBlock{
                        Content:   blockContent.String(),
                        StartLine: startLine,
                        EndLine:   j,
                        BlockType: blockType,
                    })
                    break
                }
            }
        }
    }

    // 处理内联代码块 @{ ... } 
    inlinePattern := regexp.MustCompile(`@\{`)
    for i, line := range lines {
        if inlinePattern.MatchString(line) {
            startLine := i
            bracketCount := 0
            
            // 查找开始位置
            if openBraceIndex := strings.Index(line, "@{"); openBraceIndex != -1 {
                bracketCount = 1
                
                var blockContent strings.Builder
                // 从@{开始提取
                blockContent.WriteString(line[openBraceIndex:])
                blockContent.WriteString("\n")
                
                // 检查同行是否结束
                for _, char := range line[openBraceIndex+2:] {
                    if char == '{' {
                        bracketCount++
                    } else if char == '}' {
                        bracketCount--
                    }
                }
                
                if bracketCount == 0 {
                    // 单行内联代码
                    blocks = append(blocks, RazorCodeBlock{
                        Content:   blockContent.String(),
                        StartLine: startLine,
                        EndLine:   startLine,
                        BlockType: "inline",
                    })
                } else {
                    // 多行内联代码
                    for j := i + 1; j < len(lines) && bracketCount > 0; j++ {
                        currentLine := lines[j]
                        blockContent.WriteString(currentLine)
                        blockContent.WriteString("\n")
                        
                        for _, char := range currentLine {
                            if char == '{' {
                                bracketCount++
                            } else if char == '}' {
                                bracketCount--
                            }
                        }
                        
                        if bracketCount == 0 {
                            blocks = append(blocks, RazorCodeBlock{
                                Content:   blockContent.String(),
                                StartLine: startLine,
                                EndLine:   j,
                                BlockType: "inline",
                            })
                            break
                        }
                    }
                }
            }
        }
    }

    return blocks
}

// detectRazorPageDirectives 检测Razor页面指令和生命周期方法
func (p *CSharpParser) detectRazorPageDirectives(content string, lines []string) []Function {
    var functions []Function

    // 生命周期方法检测
    lifecycleMethods := []string{
        "OnInitialized", "OnInitializedAsync",
        "OnParametersSet", "OnParametersSetAsync", 
        "OnAfterRender", "OnAfterRenderAsync",
        "Dispose", "DisposeAsync",
        "SetParametersAsync",
    }

    for _, methodName := range lifecycleMethods {
        // 正则表达式
        pattern := regexp.MustCompile(`(?m)(?:protected\s+)?(?:override\s+)?(?:async\s+)?(?:void|Task)\s+` + methodName + `\s*\(`)
        matches := pattern.FindAllStringIndex(content, -1)
        
        for _, match := range matches {
            startLine := p.getLineNumber(content, match[0])
            
            // 查找方法结束位置
            endLine := p.findMethodEnd(content, match[0], lines, startLine)
            
            functions = append(functions, Function{
                Name:       methodName,
                StartLine:  startLine,
                EndLine:    endLine,
                Complexity: p.estimateComplexity(content, match[0], endLine-startLine),
                Parameters: 0, // 生命周期方法通常无参数
            })
        }
    }

    // 检测事件处理器
    eventHandlerPattern := regexp.MustCompile(`@on\w+\s*=\s*"?([a-zA-Z_][a-zA-Z0-9_]*)"?`)
    matches := eventHandlerPattern.FindAllStringSubmatch(content, -1)
    
    handlerNames := make(map[string]bool) // 避免重复
    for _, match := range matches {
        if len(match) > 1 {
            handlerName := match[1]
            if handlerName != "" && !handlerNames[handlerName] {
                handlerNames[handlerName] = true
                
                functions = append(functions, Function{
                    Name:       handlerName,
                    StartLine:  1, // 事件处理器可能在任何地方定义
                    EndLine:    1,
                    Complexity: 1,
                    Parameters: 1, // 通常有一个事件参数
                })
            }
        }
    }

    return functions
}

// C#方法检测的正则表达式模式
var (
    // 匹配C#方法声明
    csharpMethodPattern = regexp.MustCompile(`(?m)^\s*(?:(?:public|private|protected|internal|static|virtual|override|abstract|sealed|async)\s+)*([a-zA-Z_][a-zA-Z0-9_<>\[\]]*(?:\?)?)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(([^)]*)\)\s*(?:\{|=>)`)
    
    // 匹配属性
    csharpPropertyPattern = regexp.MustCompile(`(?m)^\s*(?:(?:public|private|protected|internal|static|virtual|override|abstract)\s+)*([a-zA-Z_][a-zA-Z0-9_<>\[\]]*(?:\?)?)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\{\s*(?:get|set)`)
    
    // 匹配构造函数
    csharpConstructorPattern = regexp.MustCompile(`(?m)^\s*(?:(?:public|private|protected|internal)\s+)*([A-Z][a-zA-Z0-9_]*)\s*\(([^)]*)\)\s*(?:\{|:)`)
    
    // Razor组件参数
    razorParameterPattern = regexp.MustCompile(`(?m)^\s*\[Parameter\]\s*(?:public\s+)?([a-zA-Z_][a-zA-Z0-9_<>\[\]]*\s+)([a-zA-Z_][a-zA-Z0-9_]*)\s*\{\s*get;\s*set;\s*\}`)
)

// detectCSharpMethods 检测C#方法
func (p *CSharpParser) detectCSharpMethods(content string, lines []string) []Function {
    functions := make([]Function, 0)

    // 检测普通方法
    methods := p.detectMethodsWithPattern(content, lines, csharpMethodPattern, "method")
    functions = append(functions, methods...)

    // 检测属性
    properties := p.detectMethodsWithPattern(content, lines, csharpPropertyPattern, "property")
    functions = append(functions, properties...)

    // 检测构造函数
    constructors := p.detectMethodsWithPattern(content, lines, csharpConstructorPattern, "constructor")
    functions = append(functions, constructors...)

    return functions
}

// detectMethodsWithPattern 使用指定模式检测方法
func (p *CSharpParser) detectMethodsWithPattern(content string, lines []string, pattern *regexp.Regexp, methodType string) []Function {
    functions := make([]Function, 0)
    matches := pattern.FindAllStringSubmatchIndex(content, -1)

    for _, match := range matches {
        startPos := match[0]
        startLine := p.getLineNumber(content, startPos)

        // 提取方法信息
        submatch := pattern.FindStringSubmatch(content[startPos:match[1]])
        
        var funcName string
        var paramStr string

        switch methodType {
        case "method":
            if len(submatch) >= 3 {
                funcName = submatch[2] // 方法名
                paramStr = submatch[3] // 参数
            }
        case "property":
            if len(submatch) >= 3 {
                funcName = submatch[2] // 属性名
                paramStr = "" // 属性没有参数
            }
        case "constructor":
            if len(submatch) >= 3 {
                funcName = submatch[1] // 构造函数名
                paramStr = submatch[2] // 参数
            }
        }

        if funcName == "" {
            continue
        }

        // 查找方法结束位置
        endLine := p.findMethodEnd(content, startPos, lines, startLine)

        // 计算参数数量
        params := p.countParameters(paramStr)

        // 计算复杂度
        complexity := p.estimateComplexity(content, startPos, endLine-startLine)

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

// getLineNumber 计算字符串位置对应的行号
func (p *CSharpParser) getLineNumber(content string, pos int) int {
    line := 1
    for i := 0; i < pos && i < len(content); i++ {
        if content[i] == '\n' {
            line++
        }
    }
    return line
}

// findMethodEnd 查找方法结束位置
func (p *CSharpParser) findMethodEnd(content string, startPos int, lines []string, startLine int) int {
    // 检查是否是表达式方法体 (=>)
    lineContent := lines[startLine-1]
    if strings.Contains(lineContent, "=>") {
        // 表达式方法体，查找分号或下一行
        for i := startLine; i < len(lines); i++ {
            if strings.Contains(lines[i], ";") {
                return i + 1
            }
        }
        return startLine + 1
    }

    // 普通方法体，查找匹配的大括号
    bracketCount := 0
    inString := false
    inChar := false
    inLineComment := false
    inBlockComment := false

    for i := startPos; i < len(content); i++ {
        c := content[i]

        // 处理字符串和字符字面量
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

        // 处理注释
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

        // 统计大括号
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

        if c == '\n' {
            inLineComment = false
        }
    }

    return len(lines)
}

// countParameters 计算参数数量
func (p *CSharpParser) countParameters(paramStr string) int {
    paramStr = strings.TrimSpace(paramStr)
    if paramStr == "" {
        return 0
    }

    // 处理C#特有的参数修饰符
    paramStr = strings.ReplaceAll(paramStr, "ref ", "")
    paramStr = strings.ReplaceAll(paramStr, "out ", "")
    paramStr = strings.ReplaceAll(paramStr, "in ", "")
    paramStr = strings.ReplaceAll(paramStr, "params ", "")

    // 简单计算逗号数量
    count := strings.Count(paramStr, ",") + 1
    
    // 验证是否真的有参数
    if strings.TrimSpace(strings.ReplaceAll(paramStr, ",", "")) == "" {
        return 0
    }

    return count
}

// estimateComplexity 估算方法复杂度
func (p *CSharpParser) estimateComplexity(content string, startPos, lineCount int) int {
    complexity := 1

    // C#特有的复杂度关键字
    keywords := []string{
        "if", "else", "for", "foreach", "while", "do", "switch", "case",
        "try", "catch", "finally", "throw", "?", "??", "&&", "||",
        "when", "is", "as", "using", "lock", "yield",
    }

    // 提取方法内容
    endPos := p.findContentEndPosition(content, startPos, lineCount)
    methodContent := content[startPos:endPos]

    // 计算关键字出现次数
    for _, keyword := range keywords {
        // 使用单词边界匹配，避免误匹配
        re := regexp.MustCompile(`\b` + regexp.QuoteMeta(keyword) + `\b`)
        matches := re.FindAllStringIndex(methodContent, -1)
        if matches != nil {
            complexity += len(matches)
        }
    }

    // 特殊处理LINQ表达式
    linqKeywords := []string{"where", "select", "from", "join", "group", "order"}
    for _, keyword := range linqKeywords {
        re := regexp.MustCompile(`\b` + keyword + `\b`)
        matches := re.FindAllStringIndex(methodContent, -1)
        if matches != nil {
            complexity += len(matches)
        }
    }

    return complexity
}

// findContentEndPosition 找到内容的结束位置
func (p *CSharpParser) findContentEndPosition(content string, startPos, lineCount int) int {
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