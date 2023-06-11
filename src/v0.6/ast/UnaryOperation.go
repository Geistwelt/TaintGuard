package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type UnaryOperation struct {
	ID               int    `json:"id"`
	IsConstant       bool   `json:"isConstant"`
	IsLValue         bool   `json:"isLValue"`
	IsPure           bool   `json:"isPure"`
	LValueRequested  bool   `json:"lValueRequested"`
	NodeType         string `json:"nodeType"`
	Operator         string `json:"operator"`
	Prefix           bool   `json:"prefix"`
	Src              string `json:"src"`
	subExpression    ASTNode
	TypeDescriptions struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
}

func (uo *UnaryOperation) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	if uo.subExpression != nil {
		var expression string
		switch subExpression := uo.subExpression.(type) {
		case *Identifier:
			expression = subExpression.SourceCode(false, false, indent, logger)
		case *FunctionCall:
			expression = subExpression.SourceCode(false, false, indent, logger)
		case *IndexAccess:
			expression = subExpression.SourceCode(false, false, indent, logger)
		default:
			if subExpression != nil {
				logger.Warnf("Unknown subExpression nodeType [%s] for UnaryOperation [src:%s].", subExpression.Type(), uo.Src)
			} else {
				logger.Warnf("Unknown subExpression nodeType for UnaryOperation [src:%s].", uo.Src)
			}
		}
		if uo.Prefix {
			if uo.Operator == "delete" {
				code = code + uo.Operator + " " + expression
			} else {
				code = code + uo.Operator + expression
			}
		} else {
			code = code + expression + uo.Operator
		}
	} else {
		logger.Warnf("UnaryOperation [src:%s] should have subExpression.", uo.Src)
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (uo *UnaryOperation) Type() string {
	return uo.NodeType
}

func (uo *UnaryOperation) Nodes() []ASTNode {
	return nil
}

func (uo *UnaryOperation) NodeID() int {
	return uo.ID
}

func GetUnaryOperation(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*UnaryOperation, error) {
	uo := new(UnaryOperation)
	if err := json.Unmarshal([]byte(raw.ToString()), uo); err != nil {
		logger.Errorf("Failed to unmarshal UnaryOperation: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal UnaryOperation: [%v]", err)
	}

	// subExpression
	{
		subExpression := raw.Get("subExpression")
		if subExpression.Size() > 0 {
			subExpressionNodeType := subExpression.Get("nodeType").ToString()
			var uoSubExpression ASTNode
			var err error

			switch subExpressionNodeType {
			case "Identifier":
				uoSubExpression, err = GetIdentifier(gn, subExpression, logger)
			case "FunctionCall":
				uoSubExpression, err = GetFunctionCall(gn, subExpression, logger)
			case "IndexAccess":
				uoSubExpression, err = GetIndexAccess(gn, subExpression, logger)
			default:
				logger.Warnf("Unknown subExpression nodeType [%s] for UnaryOperation [src:%s].", subExpressionNodeType, uo.Src)
			}

			if err != nil {
				return nil, err
			}

			if uoSubExpression != nil {
				uo.subExpression = uoSubExpression
			}
		}
	}

	gn.AddASTNode(uo)

	return uo, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (uo *UnaryOperation) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	if uo.subExpression != nil {
		switch subExpression := uo.subExpression.(type) {
		case *FunctionCall:
			subExpression.TraverseFunctionCall(ncp, gn, opt, logger)
		case *IndexAccess:
			subExpression.TraverseFunctionCall(ncp, gn, opt, logger)
		}
	}
}
