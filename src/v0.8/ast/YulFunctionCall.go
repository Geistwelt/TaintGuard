package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type YulFunctionCall struct {
	arguments    []ASTNode
	functionName ASTNode
	NodeType     string `json:"nodeType"`
	Src          string `json:"src"`
}

func (yfc *YulFunctionCall) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	// functionName
	{
		if yfc.functionName != nil {
			switch functionName := yfc.functionName.(type) {
			case *YulIdentifier:
				code = code + functionName.SourceCode(false, false, indent, logger)
			default:
				if functionName != nil {
					logger.Warnf("Unknown functionName nodeType [%s] for YulFunctionCall [src:%s].", functionName.Type(), yfc.Src)
				} else {
					logger.Warnf("Unknown functionName nodeType for YulFunctionCall [src:%s].", yfc.Src)
				}
			}
		}
	}

	code = code + "("

	// arguments
	{
		if len(yfc.arguments) > 0 {
			for index, argument := range yfc.arguments {
				switch arg := argument.(type) {
				case *YulIdentifier:
					code = code + arg.SourceCode(false, false, indent, logger)
				default:
					if arg != nil {
						logger.Warnf("Unknown argument nodeType [%s] for YulFunctionCall [src:%s].", arg.Type(), yfc.Src)
					} else {
						logger.Warnf("Unknown argument nodeType for YulFunctionCall [src:%s].", yfc.Src)
					}
				}

				if index < len(yfc.arguments)-1 {
					code = code + ", "
				}
			}
		}
	}

	code = code + ")"

	return code
}

func (yfc *YulFunctionCall) Type() string {
	return yfc.NodeType
}

func (yfc *YulFunctionCall) Nodes() []ASTNode {
	return yfc.arguments
}

func GetYulFunctionCall(raw jsoniter.Any, logger logging.Logger) (*YulFunctionCall, error) {
	yfc := new(YulFunctionCall)
	if err := json.Unmarshal([]byte(raw.ToString()), yfc); err != nil {
		logger.Errorf("Failed to unmarshal YulFunctionCall: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal YulFunctionCall: [%v]", err)
	}

	// arguments
	{
		arguments := raw.Get("arguments")
		if arguments.Size() > 0 {
			yfc.arguments = make([]ASTNode, 0)

			for i := 0; i < arguments.Size(); i++ {
				argument := arguments.Get(i)
				argumentNodeType := argument.Get("nodeType").ToString()
				var yfcArgument ASTNode
				var err error

				switch argumentNodeType {
				case "YulIdentifier":
					yfcArgument, err = GetYulIdentifier(argument, logger)
				default:
					logger.Warnf("Unknown argument nodeType [%s] for YulFunctionCall [src:%s].", argumentNodeType, yfc.Src)
				}

				if err != nil {
					return nil, err
				}

				if yfcArgument != nil {
					yfc.arguments = append(yfc.arguments, yfcArgument)
				}
			}
		}
	}

	// functionName
	{
		functionName := raw.Get("functionName")
		if functionName.Size() > 0 {
			functionNameNodeType := functionName.Get("nodeType").ToString()
			var fn ASTNode
			var err error

			switch functionNameNodeType {
			case "YulIdentifier":
				fn, err = GetYulIdentifier(functionName, logger)
			default:
				logger.Warnf("Unknown functionName nodeType [%s] for YulIdentifier [src:%s].", functionNameNodeType, yfc.Src)
			}

			if err != nil {
				return nil, err
			}

			if fn != nil {
				yfc.functionName = fn
			}
		} else {
			logger.Warnf("Function name in YulFunctionCall [src:%s] should not be empty.", yfc.Src)
		}
	}

	return yfc, nil
}
