package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type ErrorDefinition struct {
	ErrorSelector string `json:"errorSelector"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	NameLocation  string `json:"nameLocation"`
	NodeType      string `json:"nodeType"`
	parameters    ASTNode
	Src           string `json:"src"`
}

func (ed *ErrorDefinition) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "error" + " " + ed.Name

	if ed.parameters != nil {
		switch parameters := ed.parameters.(type) {
		case *ParameterList:
			code = code + "(" + parameters.SourceCode(false, false, indent, logger)
		default:
			if parameters != nil {
				logger.Warnf("Unknown parameters nodeType [%s] for ErrorDefinition [src:%s].", parameters.Type(), ed.Src)
			} else {
				logger.Warnf("Unknown parameters nodeType for ErrorDefinition [src:%s].", ed.Src)
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (ed *ErrorDefinition) Type() string {
	return ed.NodeType
}

func (ed *ErrorDefinition) Nodes() []ASTNode {
	return nil
}

func (ed *ErrorDefinition) NodeID() int {
	return ed.ID
}

func GetErrorDefinition(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*ErrorDefinition, error) {
	ed := new(ErrorDefinition)
	if err := json.Unmarshal([]byte(raw.ToString()), ed); err != nil {
		logger.Errorf("Failed to unmarshal ErrorDefinition: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal ErrorDefinition: [%v]", err)
	}

	// parameters
	{
		parameters := raw.Get("parameters")
		if parameters.Size() > 0 {
			parametersNodeType := parameters.Get("nodeType").ToString()
			var edParameters ASTNode
			var err error

			switch parametersNodeType {
			case "ParameterList":
				edParameters, err = GetParameterList(gn, parameters, logger)
			default:
				logger.Warnf("Unknown parameters nodeType [%s] for ErrorDefinition [src:%s].", parametersNodeType, ed.Src)
			}

			if err != nil {
				return nil, err
			}

			if edParameters != nil {
				ed.parameters = edParameters
			}
		}
	}

	return ed, nil
}
