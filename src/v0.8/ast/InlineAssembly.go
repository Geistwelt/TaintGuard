package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type InlineAssembly struct {
	ast                ASTNode
	EvmVersion         string `json:"evmVersion"`
	ExternalReferences []struct {
		Declaration int    `json:"declaration"`
		IsOffset    bool   `json:"isOffset"`
		IsSlot      bool   `json:"isSlot"`
		Src         string `json:"src"`
		ValueSize   int    `json:"valueSize"`
	} `json:"externalReferences"`
	ID       int    `json:"id"`
	NodeType string `json:"nodeType"`
	Src      string `json:"src"`
}

func (ia *InlineAssembly) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "assembly {\n"

	if ia.ast != nil {
		switch ast := ia.ast.(type) {
		case *YulBlock:
			code = code + ast.SourceCode(false, true, indent, logger)
		default:
			if ast != nil {
				logger.Warnf("Unknown ast nodeType [%s] for InlineAssembly [src:%s]", ast.Type(), ia.Src)
			} else {
				logger.Warnf("Unknown ast nodeType for InlineAssembly [src:%s]", ia.Src)
			}
		}
	}

	code = code + "\n"

	if isIndent {
		code = code + indent
	}

	code = code + "}"

	return code
}

func (ia *InlineAssembly) Type() string {
	return ia.NodeType
}

func (ia *InlineAssembly) Nodes() []ASTNode {
	return nil
}

func GetInlineAssembly(raw jsoniter.Any, logger logging.Logger) (*InlineAssembly, error) {
	ia := new(InlineAssembly)
	if err := json.Unmarshal([]byte(raw.ToString()), ia); err != nil {
		logger.Errorf("Failed to unmarshal InlineAssembly: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal InlineAssembly: [%v]", err)
	}

	// ast
	{
		ast := raw.Get("AST")
		if ast.Size() > 0 {
			astNodeType := ast.Get("nodeType").ToString()
			var iaAST ASTNode
			var err error

			switch astNodeType {
			case "YulBlock":
				iaAST, err = GetYulBlock(ast, logger)
			default:
				logger.Warnf("Unknown ast nodeType [%s] for InlineAssembly [src:%s].", astNodeType, ia.Src)
			}

			if err != nil {
				return nil, err
			}

			if iaAST != nil {
				ia.ast = iaAST
			}
		}
	}

	return ia, nil
}
