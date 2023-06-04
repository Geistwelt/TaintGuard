package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type EventDefinition struct {
	Anonymous     bool   `json:"anonymous"`
	EventSelector string `json:"eventSelector"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	NameLocation  string `json:"nameLocation"`
	NodeType      string `json:"nodeType"`
	parameters    ASTNode
	Src           string `json:"src"`
}

func (ed *EventDefinition) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "event" + " " + ed.Name

	if ed.parameters != nil {
		switch parameters := ed.parameters.(type) {
		case *ParameterList:
			code = code + "(" + parameters.SourceCode(false, false, indent, logger) + ")"
		default:
			if parameters != nil {
				logger.Warnf("Unknown parameters nodeType [%s] for EventDefinition [src:%s].", parameters.Type(), ed.Src)
			} else {
				logger.Warnf("Unknown parameters nodeType for EventDefinition [src:%s].", ed.Src)
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (ed *EventDefinition) Type() string {
	return ed.NodeType
}

func (ed *EventDefinition) Nodes() []ASTNode {
	return nil
}

func GetEventDefinition(raw jsoniter.Any, logger logging.Logger) (*EventDefinition, error) {
	ed := new(EventDefinition)
	if err := json.Unmarshal([]byte(raw.ToString()), ed); err != nil {
		logger.Errorf("Failed to unmarshal EventDefinition: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal EventDefinition: [%v]", err)
	}

	parameters := raw.Get("parameters")
	if parameters.Size() > 0 {
		var edParameters ASTNode
		var err error
		parametersNodeType := parameters.Get("nodeType").ToString()

		switch parametersNodeType {
		case "ParameterList":
			edParameters, err = GetParameterList(parameters, logger)
		default:
			logger.Warnf("Unknown parameters nodeType [%s] for EventDefinition [src:%s].", parametersNodeType, ed.Src)
		}

		if err != nil {
			return nil, err
		}
		if edParameters != nil {
			ed.parameters = edParameters
		}
	}

	return ed, nil
}
