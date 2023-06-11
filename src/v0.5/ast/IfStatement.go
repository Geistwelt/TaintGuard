package ast

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type IfStatement struct {
	condition ASTNode
	falseBody ASTNode
	ID        int    `json:"id"`
	NodeType  string `json:"nodeType"`
	Src       string `json:"src"`
	trueBody  ASTNode
}

func (is *IfStatement) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	code = code + "if"

	// condition
	{
		if is.condition != nil {
			switch condition := is.condition.(type) {
			case *BinaryOperation:
				code = code + "(" + condition.SourceCode(false, false, indent, logger) + ")"
			case *UnaryOperation:
				code = code + "(" + condition.SourceCode(false, false, indent, logger) + ")"
			case *Literal:
				code = code + "(" + condition.SourceCode(false, false, indent, logger) + ")"
			case *IndexAccess:
				code = code + "(" + condition.SourceCode(false, false, indent, logger) + ")"
			case *Identifier:
				code = code + "(" + condition.SourceCode(false, false, indent, logger) + ")"
			default:
				if condition != nil {
					logger.Warnf("Unknown condition nodeType [%s] for IfStatement [src:%s].", condition.Type(), is.Src)
				} else {
					logger.Warnf("Unknown condition nodeType for IfStatement [src:%s].", is.Src)
				}
			}
		}
	}

	code = code + " {\n"

	// trueBody
	{
		if is.trueBody != nil {
			switch trueBody := is.trueBody.(type) {
			case *Block:
				code = code + trueBody.SourceCode(false, true, indent, logger)
			case *ExpressionStatement:
				code = code + trueBody.SourceCode(true, true, indent+"    ", logger)
			// case *RevertStatement:
			// 	code = code + trueBody.SourceCode(true, true, indent, logger)
			case *Return:
				code = code + trueBody.SourceCode(true, true, indent, logger)
			case *IfStatement:
				code = code + trueBody.SourceCode(false, true, indent, logger)
			// case *Break:
			// 	code = code + trueBody.SourceCode(true, true, indent+"    ", logger)
			default:
				if trueBody != nil {
					logger.Warnf("Unknown trueBody nodeType [%s] for IfStatement [src:%s].", trueBody.Type(), is.Src)
				} else {
					logger.Warnf("Unknown trueBody nodeType for IfStatement [src:%s].", is.Src)
				}
			}
		}

		code = code + "\n"

		if isIndent {
			code = code + indent
		}

		code = code + "}"
	}

	// falseBody
	{
		if is.falseBody != nil {
			switch falseBody := is.falseBody.(type) {
			case *Block:
				code = code + " " + "else {\n"
				code = code + falseBody.SourceCode(false, true, indent, logger)
				code = code + "\n"
				if isIndent {
					code = code + indent
				}
				code = code + "}"
			case *IfStatement:
				code = code + " " + "else "
				code = code + strings.TrimLeft(falseBody.SourceCode(false, true, indent, logger), " ")
			case *ExpressionStatement:
				code = code + " " + "else {\n"
				code = code + falseBody.SourceCode(true, true, indent, logger)
				code = code + "\n"
				if isIndent {
					code = code + indent
				}
				code = code + "}"
			default:
				if falseBody != nil {
					logger.Warnf("Unknown falseBody nodeType [%s] for IfStatement [src:%s].", falseBody.Type(), is.Src)
				} else {
					logger.Warnf("Unknown falseBody nodeType for IfStatement [src:%s].", is.Src)
				}
			}
		}
	}

	return code
}

func (is *IfStatement) Type() string {
	return is.NodeType
}

func (is *IfStatement) Nodes() []ASTNode {
	return nil
}

func (is *IfStatement) NodeID() int {
	return is.ID
}

func GetIfStatement(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*IfStatement, error) {
	is := new(IfStatement)
	if err := json.Unmarshal([]byte(raw.ToString()), is); err != nil {
		logger.Errorf("Failed to unmarshal IfStatement: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal IfSatement: [%v]", err)
	}

	// condition
	{
		condition := raw.Get("condition")
		if condition.Size() > 0 {
			conditionNodeType := condition.Get("nodeType").ToString()
			var isCondition ASTNode
			var err error

			switch conditionNodeType {
			case "BinaryOperation":
				isCondition, err = GetBinaryOperation(gn, condition, logger)
			case "UnaryOperation":
				isCondition, err = GetUnaryOperation(gn, condition, logger)
			case "Literal":
				isCondition, err = GetLiteral(gn, condition, logger)
			case "IndexAccess":
				isCondition, err = GetIndexAccess(gn, condition, logger)
			case "Identifier":
				isCondition, err = GetIdentifier(gn, condition, logger)
			default:
				logger.Warnf("Unknown condition nodeType [%s] for IfStatement [src:%s].", conditionNodeType, is.Src)
			}

			if err != nil {
				return nil, err
			}

			if isCondition != nil {
				is.condition = isCondition
			}
		}
	}

	// falseBody
	{
		falseBody := raw.Get("falseBody")
		if falseBody.Size() > 0 {
			falseBodyNodeType := falseBody.Get("nodeType").ToString()
			var isFalseBody ASTNode
			var err error

			switch falseBodyNodeType {
			case "Block":
				isFalseBody, err = GetBlock(gn, falseBody, logger)
			case "IfStatement":
				isFalseBody, err = GetIfStatement(gn, falseBody, logger)
			case "ExpressionStatement":
				isFalseBody, err = GetExpressionStatement(gn, falseBody, logger)
			default:
				logger.Warnf("Unknown falseBody nodeType [%s] for IfStatement [src:%s].", falseBodyNodeType, is.Src)
			}

			if err != nil {
				return nil, err
			}

			if isFalseBody != nil {
				is.falseBody = isFalseBody
			}
		}
	}

	// trueBody
	{
		trueBody := raw.Get("trueBody")
		if trueBody.Size() > 0 {
			trueBodyNodeType := trueBody.Get("nodeType").ToString()
			var isTrueBody ASTNode
			var err error

			switch trueBodyNodeType {
			case "Block":
				isTrueBody, err = GetBlock(gn, trueBody, logger)
			case "ExpressionStatement":
				isTrueBody, err = GetExpressionStatement(gn, trueBody, logger)
			// case "RevertStatement":
			// 	isTrueBody, err = GetRevertStatement(gn, trueBody, logger)
			case "Return":
				isTrueBody, err = GetReturn(gn, trueBody, logger)
			case "IfStatement":
				isTrueBody, err = GetIfStatement(gn, trueBody, logger)
			// case "Break":
			// 	isTrueBody, err = GetBreak(gn, trueBody, logger)
			default:
				logger.Warnf("Unknown trueBody [%s] for IfStatement [src:%s].", trueBodyNodeType, is.Src)
			}

			if err != nil {
				return nil, err
			}

			if isTrueBody != nil {
				is.trueBody = isTrueBody
			}
		}
	}

	gn.AddASTNode(is)

	return is, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (is *IfStatement) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	// condition
	{
		if is.condition != nil {
			switch condition := is.condition.(type) {
			case *BinaryOperation:
				condition.TraverseFunctionCall(ncp, gn, opt, logger)
			case *UnaryOperation:
				condition.TraverseFunctionCall(ncp, gn, opt, logger)
			case *IndexAccess:
				condition.TraverseFunctionCall(ncp, gn, opt, logger)
			}
		}
	}

	//falseBody
	{
		if is.falseBody != nil {
			switch falseBody := is.falseBody.(type) {
			case *Block:
				falseBody.TraverseFunctionCall(ncp, gn, opt, logger)
			case *ExpressionStatement:
				falseBody.TraverseFunctionCall(ncp, gn, opt, logger)
			case *IfStatement:
				falseBody.TraverseFunctionCall(ncp, gn, opt, logger)
			}
		}
	}

	// trueBody
	{
		if is.trueBody != nil {
			switch trueBody := is.trueBody.(type) {
			case *Block:
				trueBody.TraverseFunctionCall(ncp, gn, opt, logger)
			case *ExpressionStatement:
				trueBody.TraverseFunctionCall(ncp, gn, opt, logger)
			}
		}
	}
}

func (is *IfStatement) TraverseTaintOwner(opt *Option, logger logging.Logger) {
	//falseBody
	{
		if is.falseBody != nil {
			switch falseBody := is.falseBody.(type) {
			case *Block:
				falseBody.TraverseTaintOwner(opt, logger)
			case *ExpressionStatement:
				falseBody.TraverseTaintOwner(opt, logger)
			case *IfStatement:
				falseBody.TraverseTaintOwner(opt, logger)
			}
		}
	}

	// trueBody
	{
		if is.trueBody != nil {
			switch trueBody := is.trueBody.(type) {
			case *Block:
				trueBody.TraverseTaintOwner(opt, logger)
			case *ExpressionStatement:
				trueBody.TraverseTaintOwner(opt, logger)
			}
		}
	}
}

func (is *IfStatement) TraverseDelegatecall(opt *Option, logger logging.Logger) {
	//falseBody
	{
		if is.falseBody != nil {
			switch falseBody := is.falseBody.(type) {
			case *Block:
				falseBody.TraverseDelegatecall(opt, logger)
			case *ExpressionStatement:
				falseBody.TraverseDelegatecall(opt, logger)
			case *IfStatement:
				falseBody.TraverseDelegatecall(opt, logger)
			}
		}
	}

	// trueBody
	{
		if is.trueBody != nil {
			switch trueBody := is.trueBody.(type) {
			case *Block:
				trueBody.TraverseDelegatecall(opt, logger)
			case *ExpressionStatement:
				trueBody.TraverseDelegatecall(opt, logger)
			}
		}
	}
}