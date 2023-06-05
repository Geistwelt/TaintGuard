package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type IfStatement struct {
	condition ASTNode
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
			default:
				if trueBody != nil {
					logger.Warnf("Unknown trueBody nodeType [%s] for IfStatement [src:%s].", trueBody.Type(), is.Src)
				} else {
					logger.Warnf("Unknown trueBody nodeType for IfStatement [src:%s].", is.Src)
				}
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

func (is *IfStatement) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes) {
	// condition
	{
		if is.condition != nil {
			switch condition := is.condition.(type) {
			case *BinaryOperation:
				condition.TraverseFunctionCall(ncp, gn)
			}
		}
	}

	// trueBody
	{
		if is.trueBody != nil {
			switch trueBody := is.trueBody.(type) {
			case *Block:
				trueBody.TraverseFunctionCall(ncp, gn)
			}
		}
	}
}
