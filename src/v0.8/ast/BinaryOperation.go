package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type BinaryOperation struct {
	CommonType struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"commonType"`
	ID               int  `json:"id"`
	IsConstant       bool `json:"isConstant"`
	IsValue          bool `json:"isValue"`
	IsPure           bool `json:"isPure"`
	LValueRequested  bool `json:"lValueRequested"`
	leftExpression   ASTNode
	NodeType         string `json:"nodeType"`
	Operator         string `json:"operator"`
	rightExpression  ASTNode
	Src              string `json:"src"`
	TypeDescriptions struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
}

func (bo *BinaryOperation) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	if bo.leftExpression != nil {
		switch leftExpression := bo.leftExpression.(type) {
		case *Literal:
			code = code + leftExpression.SourceCode(false, false, indent, logger)
		case *Identifier:
			code = code + leftExpression.SourceCode(false, false, indent, logger)
		case *BinaryOperation:
			code = code + leftExpression.SourceCode(false, false, indent, logger)
		case *UnaryOperation:
			code = code + leftExpression.SourceCode(false, false, indent, logger)
		case *FunctionCall:
			code = code + leftExpression.SourceCode(false, false, indent, logger)
		default:
			if leftExpression != nil {
				logger.Warnf("Unknown leftExpression nodeType [%s] for BinaryOperation [src:%s].", leftExpression.Type(), bo.Src)
			} else {
				logger.Warnf("Unknown leftExpression nodeType for BinaryOperation [src:%s].", bo.Src)
			}
		}
	}

	if bo.Operator != "" {
		code = code + " " + bo.Operator
	} else {
		logger.Warnf("Missing operator in BinaryOperation.")
	}

	if bo.rightExpression != nil {
		switch rightExpression := bo.rightExpression.(type) {
		case *BinaryOperation:
			code = code + " " + rightExpression.SourceCode(false, false, indent, logger)
		case *Identifier:
			code = code + " " + rightExpression.SourceCode(false, false, indent, logger)
		case *FunctionCall:
			code = code + " " + rightExpression.SourceCode(false, false, indent, logger)
		case *Literal:
			code = code + " " + rightExpression.SourceCode(false, false, indent, logger)
		case *UnaryOperation:
			code = code + " " + rightExpression.SourceCode(false, false, indent, logger)
		case *MemberAccess:
			code = code + " " + rightExpression.SourceCode(false, false, indent, logger)
		default:
			if rightExpression != nil {
				logger.Warnf("Unknown rightExpression nodeType [%s] for BinaryOperation [src:%s].", rightExpression.Type(), bo.Src)
			} else {
				logger.Warnf("Unknown rightExpression nodeType for BinaryOperation [src:%s].", bo.Src)
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (bo *BinaryOperation) Type() string {
	return bo.NodeType
}

func (bo *BinaryOperation) Nodes() []ASTNode {
	return nil
}

func GetBinaryOperation(raw jsoniter.Any, logger logging.Logger) (*BinaryOperation, error) {
	bo := new(BinaryOperation)
	if err := json.Unmarshal([]byte(raw.ToString()), bo); err != nil {
		logger.Errorf("Failed to unmarshal BinaryOperation: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal BinaryOperation: [%v]", err)
	}
	// leftExpression
	{
		leftExpression := raw.Get("leftExpression")
		if leftExpression.Size() > 0 {
			leftExpressionNodeType := leftExpression.Get("nodeType").ToString()
			var boLeftExpression ASTNode
			var err error

			switch leftExpressionNodeType {
			case "Literal":
				boLeftExpression, err = GetLiteral(leftExpression, logger)
			case "Identifier":
				boLeftExpression, err = GetIdentifier(leftExpression, logger)
			case "BinaryOperation":
				boLeftExpression, err = GetBinaryOperation(leftExpression, logger)
			case "UnaryOperation":
				boLeftExpression, err = GetUnaryOperation(leftExpression, logger)
			case "FunctionCall":
				boLeftExpression, err = GetFunctionCall(leftExpression, logger)
			default:
				logger.Warnf("Unknown leftExpression nodeType [%s] for BinaryOperation [src:%s].", leftExpressionNodeType, bo.Src)
			}

			if err != nil {
				return nil, err
			}
			if boLeftExpression != nil {
				bo.leftExpression = boLeftExpression
			}
		}
	}

	// rightExpression
	{
		rightExpression := raw.Get("rightExpression")
		if rightExpression.Size() > 0 {
			rightExpressionNodeType := rightExpression.Get("nodeType").ToString()
			var boRightExpression ASTNode
			var err error

			switch rightExpressionNodeType {
			case "Identifier":
				boRightExpression, err = GetIdentifier(rightExpression, logger)
			case "BinaryOperation":
				boRightExpression, err = GetBinaryOperation(rightExpression, logger)
			case "FunctionCall":
				boRightExpression, err = GetFunctionCall(rightExpression, logger)
			case "Literal":
				boRightExpression, err = GetLiteral(rightExpression, logger)
			case "UnaryOperation":
				boRightExpression, err = GetUnaryOperation(rightExpression, logger)
			case "MemberAccess":
				boRightExpression, err = GetMemberAccess(rightExpression, logger)
			default:
				logger.Warnf("Unknown rightExpression nodeType [%s] for BinaryOperation [src:%s].", rightExpressionNodeType, bo.Src)
			}

			if err != nil {
				return nil, err
			}
			if boRightExpression != nil {
				bo.rightExpression = boRightExpression
			}
		}
	}

	return bo, nil
}
