package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type Return struct {
	expression               ASTNode
	FunctionReturnParameters int    `json:"functionReturnParameters"`
	ID                       int    `json:"id"`
	NodeType                 string `json:"nodeType"`
	Src                      string `json:"src"`
}

func (r *Return) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	code = code + "return"

	// expression
	{
		if r.expression != nil {
			switch expression := r.expression.(type) {
			case *IndexAccess:
				code = code + " " + expression.SourceCode(false, false, indent, logger)
			case *Literal:
				code = code + " " + expression.SourceCode(false, false, indent, logger)
			default:
				if expression != nil {
					logger.Warnf("Unknown expression nodeType [%s] for Return [src:%s].", expression.Type(), r.Src)
				} else {
					logger.Warnf("Unknown expression nodeType for Return [src:%s].", r.Src)
				}
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (r *Return) Type() string {
	return r.NodeType
}

func (r *Return) Nodes() []ASTNode {
	return nil
}

func GetReturn(raw jsoniter.Any, logger logging.Logger) (*Return, error) {
	r := new(Return)
	if err := json.Unmarshal([]byte(raw.ToString()), r); err != nil {
		logger.Errorf("Failed to unmarshal Return: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal Return: [%v]", err)
	}

	// expression
	{
		expression := raw.Get("expression")
		if expression.Size() > 0 {
			var rExpression ASTNode
			var err error
			var expressionNodeType = expression.Get("nodeType").ToString()

			switch expressionNodeType {
			case "IndexAccess":
				rExpression, err = GetIndexAccess(expression, logger)
			case "Literal":
				rExpression, err = GetLiteral(expression, logger)
			default:
				logger.Warnf("Unknown expression nodeType [%s] for Return [src:%s].", expressionNodeType, r.Src)
			}

			if err != nil {
				return nil, err
			}

			if rExpression != nil {
				r.expression = rExpression
			}
		}
	}

	return r, nil
}
