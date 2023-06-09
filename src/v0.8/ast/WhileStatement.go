package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type WhileStatement struct {
	body      ASTNode
	condition ASTNode
	ID        int    `json:"id"`
	NodeType  string `json:"nodeType"`
	Src       string `json:"src"`
}

func (ws *WhileStatement) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "while"

	if ws.condition != nil {
		switch condition := ws.condition.(type) {
		case *FunctionCall:
			code = code + " (" + condition.SourceCode(false, false, indent, logger) + ") "
		case *BinaryOperation:
			code = code + " (" + condition.SourceCode(false, false, indent, logger) + ") "
		case *Literal:
			code = code + " (" + condition.SourceCode(false, false, indent, logger) + ") "
		default:
			if condition != nil {
				logger.Warnf("Unknown condition nodeType [%s] for WhileStatement [src:%s].", condition.Type(), ws.Src)
			} else {
				logger.Warnf("Unknown condition nodeType for WhileStatement [src:%s].", ws.Src)
			}
		}
	}

	code = code + "{\n"

	if ws.body != nil {
		switch body := ws.body.(type) {
		case *Block:
			code = code + body.SourceCode(false, false, indent, logger)
		default:
			if body != nil {
				logger.Warnf("Unknown body nodeType [%s] for WhileStatement [src:%s].", body.Type(), ws.Src)
			} else {
				logger.Warnf("Unknown body nodeType for WhileStatement [src:%s].", ws.Src)
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

func (ws *WhileStatement) Type() string {
	return ws.NodeType
}

func (ws *WhileStatement) Nodes() []ASTNode {
	return nil
}

func (ws *WhileStatement) NodeID() int {
	return ws.ID
}

func GetWhileStatement(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*WhileStatement, error) {
	ws := new(WhileStatement)
	if err := json.Unmarshal([]byte(raw.ToString()), ws); err != nil {
		logger.Errorf("Failed to unmarshal WhileStatement: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal WhileStatement: [%v]", err)
	}

	// body
	{
		body := raw.Get("body")
		if body.Size() > 0 {
			bodyNodeType := body.Get("nodeType").ToString()
			var wsBody ASTNode
			var err error

			switch bodyNodeType {
			case "Block":
				wsBody, err = GetBlock(gn, body, logger)
			default:
				logger.Warnf("Unknown body nodeType [%s] for WhileStatement [src:%s].", bodyNodeType, ws.Src)
			}

			if err != nil {
				return nil, err
			}

			if wsBody != nil {
				ws.body = wsBody
			}
		}
	}

	// condition
	{
		condition := raw.Get("condition")
		if condition.Size() > 0 {
			conditionNodeType := condition.Get("nodeType").ToString()
			var wsCondition ASTNode
			var err error

			switch conditionNodeType {
			case "FunctionCall":
				wsCondition, err = GetFunctionCall(gn, condition, logger)
			case "BinaryOperation":
				wsCondition, err = GetBinaryOperation(gn, condition, logger)
			case "Literal":
				wsCondition, err = GetLiteral(gn, condition, logger)
			default:
				logger.Warnf("Unknown condition nodeType [%s] for WhileStatement [src:%s].", conditionNodeType, ws.Src)
			}

			if err != nil {
				return nil, err
			}

			if wsCondition != nil {
				ws.condition = wsCondition
			}
		}
	}

	gn.AddASTNode(ws)

	return ws, nil
}

func (ws *WhileStatement) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	if ws.body != nil {
		switch body := ws.body.(type) {
		case *Block:
			body.TraverseFunctionCall(ncp, gn, opt, logger)
		}
	}

	if ws.condition != nil {
		switch condition := ws.condition.(type) {
		case *FunctionCall:
			condition.TraverseFunctionCall(ncp, gn, opt, logger)
		}
	}
}

func (ws *WhileStatement) TraverseTaintOwner(opt *Option, logger logging.Logger) {
	if ws.body != nil {
		switch body := ws.body.(type) {
		case *Block:
			body.TraverseTaintOwner(opt, logger)
		}
	}
}

func (ws *WhileStatement) TraverseDelegatecall(opt *Option, logger logging.Logger) {
	if ws.body != nil {
		switch body := ws.body.(type) {
		case *Block:
			body.TraverseDelegatecall(opt, logger)
		}
	}

	if ws.condition != nil {
		switch condition := ws.condition.(type) {
		case *FunctionCall:
			condition.TraverseDelegatecall(opt, logger)
		}
	}
}
