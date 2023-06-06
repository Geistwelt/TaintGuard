package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type YulExpressionStatement struct {
	expression ASTNode
	NodeType   string `json:"nodeType"`
	Src        string `json:"src"`
}

func (yes *YulExpressionStatement) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	if yes.expression != nil {
		switch expression := yes.expression.(type) {
		case *YulFunctionCall:
			code = code + expression.SourceCode(false, false, indent, logger)
		default:
			if expression != nil {
				logger.Warnf("Unknown expression nodeType [%s] for YulExpressionStatement [src:%s].", expression.Type(), yes.Src)
			} else {
				logger.Warnf("Unknown expression nodeType for YulExpressionStatement [src:%s].", yes.Src)
			}
		}
	}

	return code
}

func (yes *YulExpressionStatement) Type() string {
	return yes.NodeType
}

func (yes *YulExpressionStatement) Nodes() []ASTNode {
	return nil
}

func (yes *YulExpressionStatement) NodeID() int {
	return -1
}

func GetYulExpressionStatement(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*YulExpressionStatement, error) {
	yes := new(YulExpressionStatement)
	if err := json.Unmarshal([]byte(raw.ToString()), yes); err != nil {
		logger.Errorf("Failed to unmarshal YulExpressionStatement: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal YulExpressionStatement: [%v]", err)
	}

	expression := raw.Get("expression")
	if expression.Size() > 0 {
		expressionNodeType := expression.Get("nodeType").ToString()
		var yesExpression ASTNode
		var err error

		switch expressionNodeType {
		case "YulFunctionCall":
			yesExpression, err = GetYulFunctionCall(gn, expression, logger)
		default:
			logger.Warnf("Unknown expression nodeType [%s] for YulExpressionStatement [src:%s].", expressionNodeType, yes.Src)
		}

		if err != nil {
			return nil, err
		}

		if yesExpression != nil {
			yes.expression = yesExpression
		}
	}

	return yes, nil
}
