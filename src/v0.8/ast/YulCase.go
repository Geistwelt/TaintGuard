package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type YulCase struct {
	body     ASTNode
	NodeType string `json:"nodeType"`
	Src      string `json:"src"`
	value    ASTNode
}

func (yc *YulCase) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	if yc.value != nil {
		switch value := yc.value.(type) {
		case *YulLiteral:
			if value.SourceCode(false, false, indent, logger) != "default" {
				code = code + "case" + " " + value.SourceCode(false, false, indent, logger) + " {\n"
			} else {
				code = code + "default" + " " + "{\n"
			}
		default:
			if value != nil {
				logger.Warnf("Unknown value nodeType [%s] for YulCase [src:%s].", value.Type(), yc.Src)
			} else {
				logger.Warnf("Unknown value nodeType for YulCase [src:%s].", yc.Src)
			}
		}
	}

	if yc.body != nil {
		switch body := yc.body.(type) {
		case *YulBlock:
			code = code + body.SourceCode(false, false, indent, logger)
		default:
			if body != nil {
				logger.Warnf("Unknown body nodeType [%s] for YulCase [src:%s].", body.Type(), yc.Src)
			} else {
				logger.Warnf("Unknown body nodeType for YulCase [src:%s].", yc.Src)
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

func (yc *YulCase) Type() string {
	return yc.NodeType
}

func (yc *YulCase) Nodes() []ASTNode {
	return nil
}

func (yc *YulCase) NodeID() int {
	return -1
}

func GetYulCase(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*YulCase, error) {
	yc := new(YulCase)
	if err := json.Unmarshal([]byte(raw.ToString()), yc); err != nil {
		logger.Errorf("Failed to unmarshal YulCase: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal YulCase: [%v]", err)
	}

	// body
	{
		body := raw.Get("body")
		if body.Size() > 0 {
			bodyNodeType := body.Get("nodeType").ToString()
			var ycBody ASTNode
			var err error

			switch bodyNodeType {
			case "YulBlock":
				ycBody, err = GetYulBlock(gn, body, logger)
			default:
				logger.Warnf("Unknown body nodeType [%s] for YulCase [src:%s].", bodyNodeType, yc.Src)
			}

			if err != nil {
				return nil, err
			}

			if ycBody != nil {
				yc.body = ycBody
			}
		}
	}

	// value
	{
		value := raw.Get("value")
		var ycValue ASTNode
		var err error

		if value.Size() > 0 {
			valueNodeType := value.Get("nodeType").ToString()

			switch valueNodeType {
			case "YulLiteral":
				ycValue, err = GetYulLiteral(gn, value, logger)
			default:
				logger.Warnf("Unknown value nodeType [%s] for YulCase [src:%s].", valueNodeType, yc.Src)
			}
		} else {
			if yc.body != nil {
				ycValue = &YulLiteral{Kind: "number", NodeType: "YulLiteral", Src: "xxx", Value: "default"}
			}
		}

		if err != nil {
			return nil, err
		}

		if ycValue != nil {
			yc.value = ycValue
		}
	}

	return yc, nil
}
