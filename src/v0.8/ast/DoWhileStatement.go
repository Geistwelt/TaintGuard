package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type DoWhileStatement struct {
	body      ASTNode
	condition ASTNode
	ID        int    `json:"id"`
	NodeType  string `json:"nodeType"`
	Src       string `json:"src"`
}

func (dws *DoWhileStatement) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "do {\n"

	if dws.body != nil {
		switch body := dws.body.(type) {
		case *Block:
			code = code + body.SourceCode(false, false, indent, logger)
		default:
			if body != nil {
				logger.Warnf("Unknown body nodeType [%s] for DoWhileStatement [src:%s].", body.Type(), dws.Src)
			} else {
				logger.Warnf("Unknown body nodeType for DoWhileStatement [src:%s].", dws.Src)
			}
		}
	}

	code = code + "\n"
	if isIndent {
		code = code + indent
	}
	code = code + "} while"

	if dws.condition != nil {
		switch condition := dws.condition.(type) {
		case *FunctionCall:
			code = code + " (" + condition.SourceCode(false, false, indent, logger) + ")"
		case *BinaryOperation:
			code = code + " (" + condition.SourceCode(false, false, indent, logger) + ")"
		case *Literal:
			code = code + " (" + condition.SourceCode(false, false, indent, logger) + ")"
		default:
			if condition != nil {
				logger.Warnf("Unknown condition nodeType [%s] for DoWhileStatement [src:%s].", condition.Type(), dws.Src)
			} else {
				logger.Warnf("Unknown condition nodeType for DoWhileStatement [src:%s].", dws.Src)
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (dws *DoWhileStatement) Type() string {
	return dws.NodeType
}

func (dws *DoWhileStatement) Nodes() []ASTNode {
	return nil
}

func (dws *DoWhileStatement) NodeID() int {
	return dws.ID
}

func GetDoWhileStatement(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*DoWhileStatement, error) {
	dws := new(DoWhileStatement)
	if err := json.Unmarshal([]byte(raw.ToString()), dws); err != nil {
		logger.Errorf("Failed to unmarshal DoWhileStatement: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal DoWhileStatement: [%v]", err)
	}

	// body
	{
		body := raw.Get("body")
		if body.Size() > 0 {
			bodyNodeType := body.Get("nodeType").ToString()
			var dwsBody ASTNode
			var err error

			switch bodyNodeType {
			case "Block":
				dwsBody, err = GetBlock(gn, body, logger)
			default:
				logger.Warnf("Unknown body nodeType [%s] for DoWhileStatement [src:%s].", bodyNodeType, dws.Src)
			}

			if err != nil {
				return nil, err
			}

			if dwsBody != nil {
				dws.body = dwsBody
			}
		}
	}

	// condition
	{
		condition := raw.Get("condition")
		if condition.Size() > 0 {
			conditionNodeType := condition.Get("nodeType").ToString()
			var dwsCondition ASTNode
			var err error

			switch conditionNodeType {
			case "FunctionCall":
				dwsCondition, err = GetFunctionCall(gn, condition, logger)
			case "BinaryOperation":
				dwsCondition, err = GetBinaryOperation(gn, condition, logger)
			case "Literal":
				dwsCondition, err = GetLiteral(gn, condition, logger)
			default:
				logger.Warnf("Unknown condition nodeType [%s] for DoWhileStatement [src:%s].", conditionNodeType, dws.Src)
			}

			if err != nil {
				return nil, err
			}

			if dwsCondition != nil {
				dws.condition = dwsCondition
			}
		}
	}

	gn.AddASTNode(dws)

	return dws, nil
}

func (dws *DoWhileStatement) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	if dws.body != nil {
		switch body := dws.body.(type) {
		case *Block:
			body.TraverseFunctionCall(ncp, gn, opt, logger)
		}
	}

	if dws.condition != nil {
		switch condition := dws.condition.(type) {
		case *FunctionCall:
			condition.TraverseFunctionCall(ncp, gn, opt, logger)
		}
	}
}

func (dws *DoWhileStatement) TraverseTaintOwner(opt *Option, logger logging.Logger) {
	if dws.body != nil {
		switch body := dws.body.(type) {
		case *Block:
			body.TraverseTaintOwner(opt, logger)
		}
	}
}

func (dws *DoWhileStatement) TraverseDelegatecall(opt *Option, logger logging.Logger) {
	if dws.body != nil {
		switch body := dws.body.(type) {
		case *Block:
			body.TraverseDelegatecall(opt, logger)
		}
	}

	if dws.condition != nil {
		switch condition := dws.condition.(type) {
		case *FunctionCall:
			condition.TraverseDelegatecall(opt, logger)
		}
	}
}