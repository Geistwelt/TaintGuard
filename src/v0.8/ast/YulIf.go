package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type YulIf struct {
	body      ASTNode
	condition ASTNode
	NodeType  string `json:"nodeType"`
	Src       string `json:"src"`
}

func (yi *YulIf) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "if"

	// condition
	{
		if yi.condition != nil {
			switch condition := yi.condition.(type) {
			case *YulFunctionCall:
				code = code + " " + condition.SourceCode(false, false, indent, logger)
			case *YulLiteral:
				code = code + " " + condition.SourceCode(false, false, indent, logger)
			case *YulIdentifier:
				code = code + " " + condition.SourceCode(false, false, indent, logger)
			default:
				if condition != nil {
					logger.Warnf("Unknown condition nodeType [%s] for YulIf [src:%s].", condition.Type(), yi.Src)
				} else {
					logger.Warnf("Unknown condition nodeType for YulIf [src:%s].", yi.Src)
				}
			}
		}
	}

	code = code + "{\n"

	// body
	{
		if yi.body != nil {
			switch body := yi.body.(type) {
			case *YulBlock:
				code = code + body.SourceCode(false, false, indent, logger)
			default:
				if body != nil {
					logger.Warnf("Unknown body nodeType [%s] for YulIf [src:%s].", body.Type(), yi.Src)
				} else {
					logger.Warnf("Unknown body nodeType for YulIf [src:%s].", yi.Src)
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

func (yi *YulIf) Type() string {
	return yi.NodeType
}

func (yi *YulIf) Nodes() []ASTNode {
	return nil
}

func (yi *YulIf) NodeID() int {
	return -1
}

func GetYulIf(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*YulIf, error) {
	yi := new(YulIf)
	if err := json.Unmarshal([]byte(raw.ToString()), yi); err != nil {
		logger.Errorf("Failed to unmarshal YulIf: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal YulIf: [%v]", err)
	}

	// condition
	{
		condition := raw.Get("condition")
		if condition.Size() > 0 {
			conditionNodeType := condition.Get("nodeType").ToString()
			var yiCondition ASTNode
			var err error

			switch conditionNodeType {
			case "YulFunctionCall":
				yiCondition, err = GetYulFunctionCall(gn, condition, logger)
			case "YulLiteral":
				yiCondition, err = GetYulLiteral(gn, condition, logger)
			case "YulIdentifier":
				yiCondition, err = GetYulIdentifier(gn, condition, logger)
			default:
				logger.Warnf("Unknown condition nodeType [%s] for YulIf [src:%s].", conditionNodeType, yi.Src)
			}

			if err != nil {
				return nil, err
			}

			if yiCondition != nil {
				yi.condition = yiCondition
			}
		}
	}

	// body
	{
		body := raw.Get("body")
		if body.Size() > 0 {
			bodyNodeType := body.Get("nodeType").ToString()
			var yiBody ASTNode
			var err error

			switch bodyNodeType {
			case "YulBlock":
				yiBody, err = GetYulBlock(gn, body, logger)
			default:
				logger.Warnf("Unknown body nodeType [%s] for YulIf [src:%s].", bodyNodeType, yi.Src)
			}

			if err != nil {
				return nil, err
			}

			if yiBody != nil {
				yi.body = yiBody
			}
		}
	}

	return yi, nil
}
