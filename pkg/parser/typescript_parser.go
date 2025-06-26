// Package parser 提供多语言代码解析功能
package parser

import (
	"github.com/Done-0/fuck-u-code/pkg/common"
)

// TypeScriptParser TypeScript语言解析器
type TypeScriptParser struct {
	JavaScriptParser
}

// NewTypeScriptParser 创建新的TypeScript语言解析器
func NewTypeScriptParser() Parser {
	return &TypeScriptParser{}
}

// SupportedLanguages 返回支持的语言类型
func (p *TypeScriptParser) SupportedLanguages() []common.LanguageType {
	return []common.LanguageType{common.TypeScript}
}

// Parse 解析TypeScript代码
func (p *TypeScriptParser) Parse(filePath string, content []byte) (ParseResult, error) {
	result, err := p.JavaScriptParser.Parse(filePath, content)
	if err != nil {
		return nil, err
	}

	if baseResult, ok := result.(*BaseParseResult); ok {
		baseResult.Language = common.TypeScript
	}

	return result, nil
}
