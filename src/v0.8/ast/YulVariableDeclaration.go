package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type YulVariableDeclaration struct {
	NodeType  string `json:"nodeType"`
	Src       string `json:"src"`
	value     ASTNode
	variables []ASTNode
}

func (yvd *YulVariableDeclaration) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if len(yvd.variables) > 0 {
		for index, variable := range yvd.variables {
			switch v := variable.(type) {
			case *YulTypedName:
				code = code + v.SourceCode(false, false, indent, logger)
			default:
				if v != nil {
					logger.Warnf("Unknown variable nodeType [%s] for YulVariableDeclaration [src:%s].", v.Type(), yvd.Src)
				} else {
					logger.Warnf("Unknown variable nodeType for YulVariableDeclaration [src:%s].", yvd.Src)
				}
			}

			if index < len(yvd.variables)-1 {
				code = code + ", "
			}
		}
	}

	if len(yvd.variables) > 1 {
		code = "(" + code + ")"
	}

	code = "let" + " " + code

	if isIndent {
		code = indent + code
	}

	if yvd.value != nil {
		switch value := yvd.value.(type) {
		case *YulFunctionCall:
			code = code + " := " + value.SourceCode(false, false, indent, logger)
		case *YulIdentifier:
			code = code + " := " + value.SourceCode(false, false, indent, logger)
		default:
			if value != nil {
				logger.Warnf("Unknown value nodeType [%s] for YulVariableDeclaration [src:%s].", value.Type(), yvd.Src)
			} else {
				logger.Warnf("Unknown value nodeType for YulVariableDeclaration [src:%s].", yvd.Src)
			}
		}
	}

	return code
}

func (yvd *YulVariableDeclaration) Type() string {
	return yvd.NodeType
}

func (yvd *YulVariableDeclaration) Nodes() []ASTNode {
	return yvd.variables
}

func (yvd *YulVariableDeclaration) NodeID() int {
	return -1
}

func GetYulVariableDeclaration(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*YulVariableDeclaration, error) {
	yvd := new(YulVariableDeclaration)
	if err := json.Unmarshal([]byte(raw.ToString()), yvd); err != nil {
		logger.Errorf("Failed to unmarshal YulVariableDeclaration: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal YulVariableDeclaration: [%v]", err)
	}

	// value
	{
		value := raw.Get("value")
		if value.Size() > 0 {
			valueNodeType := value.Get("nodeType").ToString()
			var yvdValue ASTNode
			var err error

			switch valueNodeType {
			case "YulFunctionCall":
				yvdValue, err = GetYulFunctionCall(gn, value, logger)
			case "YulIdentifier":
				yvdValue, err = GetYulIdentifier(gn, value, logger)
			default:
				logger.Warnf("Unknown value nodeType [%s] for YulVariableDeclaration [src:%s].", valueNodeType, yvd.Src)
			}

			if err != nil {
				return nil, err
			}

			if yvdValue != nil {
				yvd.value = yvdValue
			}
		}
	}

	// variables
	{
		variables := raw.Get("variables")
		if variables.Size() > 0 {
			yvd.variables = make([]ASTNode, 0)

			for i := 0; i < variables.Size(); i++ {
				variable := variables.Get(i)
				if variable.Size() > 0 {
					variableNodeType := variable.Get("nodeType").ToString()
					var yvdVariable ASTNode
					var err error

					switch variableNodeType {
					case "YulTypedName":
						yvdVariable, err = GetYulTypedName(gn, variable, logger)
					default:
						logger.Warnf("Unknown variable nodeType [%s] for YulVariableDeclaration [src:%s].", variableNodeType, yvd.Src)
					}

					if err != nil {
						return nil, err
					}

					if yvdVariable != nil {
						yvd.variables = append(yvd.variables, yvdVariable)
					}
				}	
			}
		}
	}

	return yvd, nil
}
