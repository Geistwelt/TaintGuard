package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type ExpressionStatement struct {
	expression ASTNode
	ID         int    `json:"id"`
	NodeType   string `json:"nodeType"`
	Src        string `json:"src"`
}

func (es *ExpressionStatement) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	// expression
	if es.expression != nil {
		switch expression := es.expression.(type) {
		case *Assignment:
			code = code + expression.SourceCode(false, false, indent, logger)
		case *FunctionCall:
			code = code + expression.SourceCode(false, false, indent, logger)
		case *UnaryOperation:
			code = code + expression.SourceCode(false, false, indent, logger)
		case *Identifier:
			code = code + expression.SourceCode(false, false, indent, logger)
		default:
			if expression != nil {
				logger.Warnf("Unknown expression nodeType [%s] for ExpressionStatement [src:%s].", expression.Type(), es.Src)
			} else {
				logger.Warnf("Unknown expression nodeType for ExpressionStatement [src:%s].", es.Src)
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (es *ExpressionStatement) Type() string {
	return es.NodeType
}

func (es *ExpressionStatement) Nodes() []ASTNode {
	return nil
}

func (es *ExpressionStatement) NodeID() int {
	return es.ID
}

func GetExpressionStatement(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*ExpressionStatement, error) {
	es := new(ExpressionStatement)
	if err := json.Unmarshal([]byte(raw.ToString()), es); err != nil {
		logger.Errorf("Failed to unmarshal ExpressionStatement: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal ExpressionStatement: [%v]", err)
	}

	// expression
	{
		expression := raw.Get("expression")
		if expression.Size() > 0 {
			var expressionNodeType = expression.Get("nodeType").ToString()
			var esExpression ASTNode
			var err error

			switch expressionNodeType {
			case "Assignment":
				esExpression, err = GetAssignment(gn, expression, logger)
			case "FunctionCall":
				esExpression, err = GetFunctionCall(gn, expression, logger)
			case "UnaryOperation":
				esExpression, err = GetUnaryOperation(gn, expression, logger)
			case "Identifier":
				esExpression, err = GetIdentifier(gn, expression, logger)
			default:
				logger.Warnf("Unknown expression nodeType [%s] for ExpressionStatement [src:%s].", expressionNodeType, es.Src)
			}

			if err != nil {
				return nil, err
			}

			if esExpression != nil {
				es.expression = esExpression
			}
		}
	}

	gn.AddASTNode(es)

	return es, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (es *ExpressionStatement) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes) {
	if es.expression != nil {
		switch expression := es.expression.(type) {
		case *Assignment:
			expression.TraverseFunctionCall(ncp, gn)
		case *FunctionCall:
			expression.TraverseFunctionCall(ncp, gn)
		case *UnaryOperation:
			expression.TraverseFunctionCall(ncp, gn)
		}
	}
}
