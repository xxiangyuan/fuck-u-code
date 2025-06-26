// Package parser 提供多语言代码解析功能
package parser

import (
	"fmt"
	"strings"

	"github.com/Done-0/fuck-u-code/pkg/common"
	"github.com/robertkrimen/otto"
)

// JavaScriptParser JavaScript语言解析器
type JavaScriptParser struct{}

// NewJavaScriptParser 创建新的JavaScript语言解析器
func NewJavaScriptParser() Parser {
	return &JavaScriptParser{}
}

// Parse 解析JavaScript代码
func (p *JavaScriptParser) Parse(filePath string, content []byte) (ParseResult, error) {
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	result := &BaseParseResult{
		Functions:    make([]Function, 0),
		CommentLines: 0,
		TotalLines:   len(lines),
		Language:     common.JavaScript,
	}

	// 计算注释行数
	result.CommentLines = p.countCommentLines(contentStr)

	// 使用otto执行解析JavaScript的辅助脚本
	vm := otto.New()
	_, err := vm.Run(jsASTParserScript)
	if err != nil {
		return nil, fmt.Errorf("加载JS解析脚本失败: %w", err)
	}

	// 调用JS函数解析代码
	jsValue, err := vm.Call("parseJavaScript", nil, contentStr)
	if err != nil {
		return nil, fmt.Errorf("解析JavaScript代码失败: %w", err)
	}

	// 获取解析结果
	parsedObj, err := jsValue.Export()
	if err != nil {
		return nil, fmt.Errorf("导出解析结果失败: %w", err)
	}

	// 保存AST根节点
	result.ASTRoot = parsedObj

	// 转换函数信息
	if parsedData, ok := parsedObj.(map[string]interface{}); ok {
		if functions, ok := parsedData["functions"].([]interface{}); ok {
			for _, funcData := range functions {
				if funcMap, ok := funcData.(map[string]interface{}); ok {
					function := Function{
						Name:       fmt.Sprintf("%v", funcMap["name"]),
						StartLine:  int(funcMap["startLine"].(float64)),
						EndLine:    int(funcMap["endLine"].(float64)),
						Complexity: int(funcMap["complexity"].(float64)),
						Parameters: int(funcMap["parameters"].(float64)),
						Node:       funcMap["node"],
					}
					result.Functions = append(result.Functions, function)
				}
			}
		}
	}

	return result, nil
}

// SupportedLanguages 返回支持的语言类型
func (p *JavaScriptParser) SupportedLanguages() []common.LanguageType {
	return []common.LanguageType{common.JavaScript}
}

// countCommentLines 计算JavaScript代码中的注释行数
func (p *JavaScriptParser) countCommentLines(content string) int {
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

// JavaScript AST解析脚本
// 利用JavaScript引擎本身解析JavaScript代码
const jsASTParserScript = `
function parseJavaScript(code) {
    // 简单的JavaScript解析器
    var functions = [];
    var result = { functions: functions };
    
    try {
        // 使用Function构造器来模拟解析，但不执行
        // 这是一个简化的方法，实际项目中应该使用专门的解析器库
        new Function(code);
        
        // 提取函数定义
        var functionRegex = /function\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\(([^)]*)\)|([a-zA-Z_$][a-zA-Z0-9_$]*)\s*=\s*function\s*\(([^)]*)\)|([a-zA-Z_$][a-zA-Z0-9_$]*)\s*:\s*function\s*\(([^)]*)\)|(const|let|var)\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*=\s*\(([^)]*)\)\s*=>/g;
        
        var lines = code.split('\n');
        var match;
        
        while ((match = functionRegex.exec(code)) !== null) {
            var funcCode = match[0];
            var startPos = match.index;
            
            // 计算行号
            var startLine = 1;
            for (var i = 0; i < startPos; i++) {
                if (code[i] === '\n') startLine++;
            }
            
            // 估算函数结束位置
            var endLine = startLine;
            var bracketCount = 0;
            var inString = false;
            var inComment = false;
            var bodyStart = code.indexOf('{', startPos);
            
            if (bodyStart !== -1) {
                for (var i = bodyStart; i < code.length; i++) {
                    var char = code[i];
                    
                    if (char === '\n') {
                        endLine++;
                    }
                    
                    if (char === '"' && code[i-1] !== '\\') {
                        inString = !inString;
                        continue;
                    }
                    
                    if (!inString) {
                        if (char === '/' && code[i+1] === '/') {
                            inComment = true;
                            continue;
                        }
                        
                        if (inComment && char === '\n') {
                            inComment = false;
                            continue;
                        }
                        
                        if (!inComment) {
                            if (char === '{') bracketCount++;
                            if (char === '}') {
                                bracketCount--;
                                if (bracketCount === 0) break;
                            }
                        }
                    }
                }
            }
            
            // 提取函数名
            var funcName = "";
            for (var i = 1; i < match.length; i++) {
                if (match[i] && typeof match[i] === 'string' && 
                    !['function', 'const', 'let', 'var', '=', ':'].includes(match[i])) {
                    funcName = match[i];
                    break;
                }
            }
            
            // 提取参数
            var params = 0;
            for (var i = 2; i < match.length; i++) {
                if (match[i] && typeof match[i] === 'string' && match[i].indexOf('(') === -1) {
                    var paramStr = match[i].trim();
                    if (paramStr) {
                        params = paramStr.split(',').length;
                    }
                    break;
                }
            }
            
            // 计算复杂度
            var complexity = 1;
            var funcBody = code.substring(startPos, i);
            var complexityPatterns = [/if\s*\(/g, /else\s+if/g, /else/g, /for\s*\(/g, /while\s*\(/g, /case\s+/g, /catch\s*\(/g, /\&\&/g, /\|\|/g, /\?/g];
            
            for (var j = 0; j < complexityPatterns.length; j++) {
                var matches = funcBody.match(complexityPatterns[j]);
                if (matches) complexity += matches.length;
            }
            
            functions.push({
                name: funcName || "anonymous",
                startLine: startLine,
                endLine: endLine,
                complexity: complexity,
                parameters: params,
                node: { type: "FunctionDeclaration", body: funcBody }
            });
        }
        
    } catch (e) {
        result.error = e.message;
    }
    
    return result;
}
`
