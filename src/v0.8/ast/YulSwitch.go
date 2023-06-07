package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type YulSwitch struct {
	cases      []ASTNode
	expression ASTNode
	NodeType   string `json:"nodeType"`
	Src        string `json:"src"`
}

func (ys *YulSwitch) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "switch"

	if ys.expression != nil {
		switch expression := ys.expression.(type) {
		case *YulFunctionCall:
			code = code + " " + expression.SourceCode(false, false, indent, logger) + "\n"
		case *YulIdentifier:
			code = code + " " + expression.SourceCode(false, false, indent, logger) + "\n"
		default:
			if expression != nil {
				logger.Warnf("Unknown expression nodeType [%s] for YulSwitch [src:%s].", expression.Type(), ys.Src)
			} else {
				logger.Warnf("Unknown expression nodeType for YulSwitch [src:%s].", ys.Src)
			}
		}
	}

	if len(ys.cases) > 0 {
		for _, c := range ys.cases {
			switch ce := c.(type) {
			case *YulCase:
				code = code + ce.SourceCode(false, false, indent, logger) + "\n"
			default:
				if ce != nil {
					logger.Warnf("Unknown case nodeType [%s] for YulSwitch [src:%s].", ce.Type(), ys.Src)
				} else {
					logger.Warnf("Unknown case nodeType for YulSwitch [src:%s].", ys.Src)
				}
			}
		}
	}

	return code
}

func (ys *YulSwitch) Type() string {
	return ys.NodeType
}

func (ys *YulSwitch) Nodes() []ASTNode {
	return ys.cases
}

func (ys *YulSwitch) NodeID() int {
	return ys.NodeID()
}

func GetYulSwitch(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*YulSwitch, error) {
	ys := new(YulSwitch)
	if err := json.Unmarshal([]byte(raw.ToString()), ys); err != nil {
		logger.Errorf("Failed to unmarshal YulSwitch: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal YulSwitch: [%v]", err)
	}

	// cases
	{
		cases := raw.Get("cases")
		if cases.Size() > 0 {
			ys.cases = make([]ASTNode, 0)
			for i := 0; i < cases.Size(); i++ {
				c := cases.Get(i)
				if c.Size() > 0 {
					caseNodeType := c.Get("nodeType").ToString()
					var ysCase ASTNode
					var err error

					switch caseNodeType {
					case "YulCase":
						ysCase, err = GetYulCase(gn, c, logger)
					default:
						logger.Warnf("Unknown case nodeType [%s] for YulSwitch [src:%s].", caseNodeType, ys.Src)
					}

					if err != nil {
						return nil, err
					}

					if ysCase != nil {
						ys.cases = append(ys.cases, ysCase)
					}
				} else {
					logger.Warnf("Case in YulSwitch [src:%s] should not be empty.", ys.Src)
				}
			}
		}
	}

	// expression
	{
		expression := raw.Get("expression")
		if expression.Size() > 0 {
			expressionNodeType := expression.Get("nodeType").ToString()
			var ysExpression ASTNode
			var err error

			switch expressionNodeType {
			case "YulFunctionCall":
				ysExpression, err = GetYulFunctionCall(gn, expression, logger)
			case "YulIdentifier":
				ysExpression, err = GetYulIdentifier(gn, expression, logger)
			default:
				logger.Warnf("Unknown expression nodeType [%s] for YulSwitch [src:%s].", expressionNodeType, ys.Src)
			}

			if err != nil {
				return nil, err
			}

			if ysExpression != nil {
				ys.expression = ysExpression
			}
		}
	}

	return ys, nil
}
