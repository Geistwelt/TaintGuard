package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type YulAssignment struct {
	NodeType      string `json:"nodeType"`
	Src           string `json:"src"`
	value         ASTNode
	variableNames []ASTNode
}

func (ya *YulAssignment) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	// variableNames
	{
		if len(ya.variableNames) > 0 {
			for index, variableName := range ya.variableNames {
				switch variable := variableName.(type) {
				case *YulIdentifier:
					code = code + variable.SourceCode(false, false, indent, logger)
				default:
					if variable != nil {
						logger.Warnf("Unknown variable nodeType [%s] for YulAssignment [src:%s].", variable.Type(), ya.Src)
					} else {
						logger.Warnf("Unknown variable nodeType for YulAssignment [src:%s].", ya.Src)
					}
				}
				if index < len(ya.variableNames)-1 {
					code = code + ", "
				}
			}
		} else {
			logger.Warnf("YulAssignment [src:%s] should have more than 0 variables.", ya.Src)
		}
	}

	code = code + " := "

	// value
	{
		if ya.value != nil {
			switch value := ya.value.(type) {
			case *YulFunctionCall:
				code = code + value.SourceCode(false, false, indent, logger)
			case *YulIdentifier:
				code = code + value.SourceCode(false, false, indent, logger)
			default:
				if value != nil {
					logger.Warnf("Unknown value nodeType [%s] for YulAssignment [src:%s].", value.Type(), ya.Src)
				} else {
					logger.Warnf("Unknown value nodeType for YulAssignment [src:%s].", ya.Src)
				}
			}
		}
	}

	return code
}

func (ya *YulAssignment) Type() string {
	return ya.NodeType
}

func (ya *YulAssignment) Nodes() []ASTNode {
	return nil
}

func (ya *YulAssignment) NodeID() int {
	return -1
}

func GetYulAssignment(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*YulAssignment, error) {
	ya := new(YulAssignment)
	if err := json.Unmarshal([]byte(raw.ToString()), ya); err != nil {
		logger.Errorf("Failed to unmarshal YulAssignment: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal YulAssignment: [%v]", err)
	}

	// variableNames
	{
		variableNames := raw.Get("variableNames")
		if variableNames.Size() > 0 {
			for i := 0; i < variableNames.Size(); i++ {
				variableName := variableNames.Get(i)
				if variableName.Size() > 0 {
					variableNameNodeType := variableName.Get("nodeType").ToString()
					var yaVariableName ASTNode
					var err error

					switch variableNameNodeType {
					case "YulIdentifier":
						yaVariableName, err = GetYulIdentifier(gn, variableName, logger)
					default:
						logger.Warnf("Unknown variableName [%s] for YulIdentifier [src:%s].", variableNameNodeType, ya.Src)
					}

					if err != nil {
						return nil, err
					}

					if yaVariableName != nil {
						ya.variableNames = append(ya.variableNames, yaVariableName)
					}
				} else {
					logger.Warnf("Variable name in YulAssignment [src:%s] should not be empty.", ya.Src)
				}
			}
		}
	}

	// value
	{
		value := raw.Get("value")
		if value.Size() > 0 {
			valueNodeType := value.Get("nodeType").ToString()
			var yaValue ASTNode
			var err error

			switch valueNodeType {
			case "YulFunctionCall":
				yaValue, err = GetYulFunctionCall(gn, value, logger)
			case "YulIdentifier":
				yaValue, err = GetYulIdentifier(gn, value, logger)
			default:
				logger.Warnf("Unknown value nodeType [%s] for YulAssignment [src:%s].", valueNodeType, ya.Src)
			}

			if err != nil {
				return nil, err
			}

			if yaValue != nil {
				ya.value = yaValue
			}
		}
	}

	return ya, nil
}
