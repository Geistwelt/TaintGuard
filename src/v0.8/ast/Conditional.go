package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type Conditional struct {
	condition        ASTNode
	falseExpression  ASTNode
	ID               int    `json:"id"`
	IsConstant       bool   `json:"isConstant"`
	IsLValue         bool   `json:"isLValue"`
	IsPure           bool   `json:"isPure"`
	LValueRequested  bool   `json:"lValueRequested"`
	NodeType         string `json:"nodeType"`
	Src              string `json:"src"`
	trueExpression   ASTNode
	TypeDescriptions struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
}

func (c *Conditional) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	if c.condition != nil {
		switch condition := c.condition.(type) {
		case *TupleExpression:
			code = code + condition.SourceCode(false, false, indent, logger)
		case *BinaryOperation:
			code = code + condition.SourceCode(false, false, indent, logger)
		case *Identifier:
			code = code + condition.SourceCode(false, false, indent, logger)
		default:
			if condition != nil {
				logger.Warnf("Unknown condition nodeType [%s] for Conditional [src:%s].", condition.Type(), c.Src)
			} else {
				logger.Warnf("Unknown condition nodeType for Conditional [src:%s].", c.Src)
			}
		}
	} else {
		logger.Warnf("Condition in Conditional [src:%s] should not be nil.", c.Src)
	}

	code = code + "?"

	if c.trueExpression != nil {
		switch trueExpression := c.trueExpression.(type) {
		case *Identifier:
			code = code + trueExpression.SourceCode(false, false, indent, logger)
		case *MemberAccess:
			code = code + trueExpression.SourceCode(false, false, indent, logger)
		case *Literal:
			code = code + trueExpression.SourceCode(false, false, indent, logger)
		case *FunctionCall:
			code = code + trueExpression.SourceCode(false, false, indent, logger)
		case *BinaryOperation:
			code = code + trueExpression.SourceCode(false, false, indent, logger)
		default:
			if trueExpression != nil {
				logger.Warnf("Unknown trueExpression nodeType [%s] for Conditional [src:%s].", trueExpression.Type(), c.Src)
			} else {
				logger.Warnf("Unknown trueExpression nodeType for Conditional [src:%s].", c.Src)
			}
		}
	} else {
		logger.Warnf("TrueExpression in Conditional [src:%s] should not be nil.", c.Src)
	}

	code = code + ":"

	if c.falseExpression != nil {
		switch falseExpression := c.falseExpression.(type) {
		case *Identifier:
			code = code + falseExpression.SourceCode(false, false, indent, logger)
		case *FunctionCall:
			code = code + falseExpression.SourceCode(false, false, indent, logger)
		case *Literal:
			code = code + falseExpression.SourceCode(false, false, indent, logger)
		case *BinaryOperation:
			code = code + falseExpression.SourceCode(false, false, indent, logger)
		default:
			if falseExpression != nil {
				logger.Warnf("Unknown falseExpression nodeType [%s] for Conditional [src:%s].", falseExpression.Type(), c.Src)
			} else {
				logger.Warnf("Unknown falseExpression nodeType for Conditional [src:%s].", c.Src)
			}
		}
	} else {
		logger.Warnf("FalseExpression in Conditional [src:%s] should not be nil.", c.Src)
	}

	return code
}

func (c *Conditional) Type() string {
	return c.NodeType
}

func (c *Conditional) Nodes() []ASTNode {
	return nil
}

func (c *Conditional) NodeID() int {
	return c.ID
}

func GetConditional(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*Conditional, error) {
	c := new(Conditional)
	if err := json.Unmarshal([]byte(raw.ToString()), c); err != nil {
		logger.Errorf("Failed to unmarshal Conditional: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal Conditional: [%v]", err)
	}

	// condition
	{
		condition := raw.Get("condition")
		if condition.Size() > 0 {
			conditionNodeType := condition.Get("nodeType").ToString()
			var cCondition ASTNode
			var err error

			switch conditionNodeType {
			case "TupleExpression":
				cCondition, err = GetTupleExpression(gn, condition, logger)
			case "BinaryOperation":
				cCondition, err = GetBinaryOperation(gn, condition, logger)
			case "Identifier":
				cCondition, err = GetIdentifier(gn, condition, logger)
			default:
				logger.Warnf("Unknown condition nodeType [%s] for Conditional [src:%s].", conditionNodeType, c.Src)
			}

			if err != nil {
				return nil, err
			}

			if cCondition != nil {
				c.condition = cCondition
			}
		}
	}

	// falseExpression
	{
		falseExpression := raw.Get("falseExpression")
		if falseExpression.Size() > 0 {
			falseExpressionNodeType := falseExpression.Get("nodeType").ToString()
			var cFalseExpression ASTNode
			var err error

			switch falseExpressionNodeType {
			case "Identifier":
				cFalseExpression, err = GetIdentifier(gn, falseExpression, logger)
			case "FunctionCall":
				cFalseExpression, err = GetFunctionCall(gn, falseExpression, logger)
			case "Literal":
				cFalseExpression, err = GetLiteral(gn, falseExpression, logger)
			case "BinaryOperation":
				cFalseExpression, err = GetBinaryOperation(gn, falseExpression, logger)
			default:
				logger.Warnf("Unknown faleExpression nodeType [%s] for Conditional [src:%s].", falseExpressionNodeType, c.Src)
			}

			if err != nil {
				return nil, err
			}

			if cFalseExpression != nil {
				c.falseExpression = cFalseExpression
			}
		}
	}

	// trueExpression
	{
		trueExpression := raw.Get("trueExpression")
		if trueExpression.Size() > 0 {
			trueExpressionNodeType := trueExpression.Get("nodeType").ToString()
			var cTrueExpression ASTNode
			var err error

			switch trueExpressionNodeType {
			case "Identifier":
				cTrueExpression, err = GetIdentifier(gn, trueExpression, logger)
			case "MemberAccess":
				cTrueExpression, err = GetMemberAccess(gn, trueExpression, logger)
			case "Literal":
				cTrueExpression, err = GetLiteral(gn, trueExpression, logger)
			case "FunctionCall":
				cTrueExpression, err = GetFunctionCall(gn, trueExpression, logger)
			case "BinaryOperation":
				cTrueExpression, err = GetBinaryOperation(gn, trueExpression, logger)
			default:
				logger.Warnf("Unknown trueExpression nodeType [%s] for Conditional [src:%s].", trueExpressionNodeType, c.Src)
			}

			if err != nil {
				return nil, err
			}

			if cTrueExpression != nil {
				c.trueExpression = cTrueExpression
			}
		}
	}

	gn.AddASTNode(c)

	return c, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (c *Conditional) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	if c.condition != nil {
		switch condition := c.condition.(type) {
		case *TupleExpression:
			condition.TraverseFunctionCall(ncp, gn, opt, logger)
		case *BinaryOperation:
			condition.TraverseFunctionCall(ncp, gn, opt, logger)
		}
	}

	if c.falseExpression != nil {
		switch falseExpression := c.falseExpression.(type) {
		case *FunctionCall:
			falseExpression.TraverseFunctionCall(ncp, gn, opt, logger)
		}
	}
}
