package parser

import (
	"context"
	"github.com/Done-0/fuck-u-code/pkg/common"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/rust"
)

// RustParser RustParser
type RustParser struct {
	rootNode *sitter.Node
}

// NewRustParser create a new rust parser
func NewRustParser() Parser {
	return &RustParser{}
}

// Parse parse rust codes
func (r *RustParser) Parse(filePath string, content []byte) (ParseResult, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(rust.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, err
	}

	rootNode := tree.RootNode()
	r.rootNode = rootNode

	functions, err := r.parseFunction(content)
	if err != nil {
		return nil, err
	}

	totalLines := int(rootNode.EndPoint().Row) + 1
	commentLines, err := r.calculateCommentLines()
	if err != nil {
		return nil, err
	}

	return &BaseParseResult{
		Functions:    functions,
		CommentLines: commentLines,
		TotalLines:   totalLines,
		Language:     common.Rust,
		ASTRoot:      nil,
	}, nil
}

// SupportedLanguages return supported languages
func (r *RustParser) SupportedLanguages() []common.LanguageType {
	return []common.LanguageType{common.Rust}
}

// parseFunction parse functions
func (r *RustParser) parseFunction(allContent []byte) ([]Function, error) {
	queryStr := `
((function_item
  name: (identifier) @function.name
  parameters: (parameters) @function.parameters
  body: (block) @function.body) @function.definition)
`

	query, err := sitter.NewQuery([]byte(queryStr), rust.GetLanguage())
	if err != nil {
		return nil, err
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(query, r.rootNode)

	var functions []Function

	for {
		match, ok := qc.NextMatch()
		if !ok {
			break
		}

		var funcNameNode, paramsNode, bodyNode, functionNode *sitter.Node

		for _, capture := range match.Captures {
			captureName := query.CaptureNameForId(capture.Index)

			switch captureName {
			case "function.name":
				funcNameNode = capture.Node
			case "function.parameters":
				paramsNode = capture.Node
			case "function.body":
				bodyNode = capture.Node
			case "function.definition":
				functionNode = capture.Node
			}

			if functionNode != nil && funcNameNode != nil && paramsNode != nil && bodyNode != nil {
				function, err := parseSingleFunction(funcNameNode, paramsNode, bodyNode, functionNode, allContent)
				if err != nil {
					return nil, err
				}

				functions = append(functions, function)
			}
		}
	}

	return functions, nil
}

// parseSingleFunction parse a single function
func parseSingleFunction(
	funcNameNode *sitter.Node,
	paramsNode *sitter.Node,
	bodyNode *sitter.Node,
	functionNode *sitter.Node,
	src []byte,
) (Function, error) {
	name := funcNameNode.Content(src)

	// row line number starts with 0
	startLine := int(functionNode.StartPoint().Row) + 1
	endLine := int(functionNode.EndPoint().Row) + 1

	// parameters may contain `self`, or may not have type annotations. So we should find parameter nodes
	paramQuery, err := sitter.NewQuery([]byte("((parameter) @param)"), rust.GetLanguage())
	if err != nil {
		return Function{}, err
	}

	paramCursor := sitter.NewQueryCursor()
	paramCursor.Exec(paramQuery, paramsNode)
	paramCount := 0

	for {
		_, ok := paramCursor.NextMatch()
		if !ok {
			break
		}
		paramCount++
	}

	complexity := 1
	complexityQuery, err := sitter.NewQuery(
		[]byte(`[ (if_expression) (for_expression) (while_expression) (match_arm) (binary_expression) ] @complexity`),
		rust.GetLanguage(),
	)
	if err != nil {
		return Function{}, err
	}

	complexityCursor := sitter.NewQueryCursor()
	complexityCursor.Exec(complexityQuery, bodyNode)
	for {
		_, ok := complexityCursor.NextMatch()
		if !ok {
			break
		}
		complexity++
	}

	return Function{
		Name:       name,
		StartLine:  startLine,
		EndLine:    endLine,
		Parameters: paramCount,
		Complexity: complexity,
		Node:       nil,
	}, nil
}

// calculateCommentLines calculate comment lines
func (r *RustParser) calculateCommentLines() (int, error) {
	commentQueryStr := `
[
  (line_comment)
  (block_comment)
] @comment
`
	commentQuery, err := sitter.NewQuery([]byte(commentQueryStr), rust.GetLanguage())
	if err != nil {
		return 0, err
	}

	commentCursor := sitter.NewQueryCursor()
	commentCursor.Exec(commentQuery, r.rootNode)
	commentedLines := make(map[int]bool)

	for {
		match, ok := commentCursor.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			commentNode := capture.Node

			startLine := int(commentNode.StartPoint().Row)
			endLine := int(commentNode.EndPoint().Row)

			for line := startLine; line <= endLine; line++ {
				commentedLines[line] = true
			}
		}
	}

	return len(commentedLines), nil
}
