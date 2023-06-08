package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type ModifierDefinition struct {
	body         ASTNode
	ID           int    `json:"id"`
	Name         string `json:"name"`
	NameLocation string `json:"nameLocation"`
	NodeType     string `json:"nodeType"`
	parameters   ASTNode
	Src          string `json:"src"`
	Virtual      bool   `json:"virtual"`
	Visibility   string `json:"visibility"`
}

func (md *ModifierDefinition) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "modifier"

	code = code + " " + md.Name

	// parameters
	{
		if md.parameters != nil {
			switch parameters:= md.parameters.(type) {
			case *ParameterList:
				code = code + "(" + parameters.SourceCode(false, false, indent, logger) + ")"
			default:
				if parameters != nil {
					logger.Warnf("Unknown parameters nodeType [%s] for ModifierDefinition [src:%s].", parameters.Type(), md.Src)
				} else {
					logger.Warnf("Unknown parameters nodeType for ModifierDefinition [src:%s].", md.Src)
				}
			}
		}
	}

	code = code + " " + "{\n"

	// body
	{
		if md.body != nil {
			switch body := md.body.(type) {
			case *Block:
				code = code + body.SourceCode(false, false, indent, logger)
			default:
				if body != nil {
					logger.Warnf("Unknown body nodeType [%s] for ModifierDefinition [src:%s].", body.Type(), md.Src)
				} else {
					logger.Errorf("Unknown body nodeType for ModifierDefinition [src:%s].", md.Src)
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

func (md *ModifierDefinition) Type() string {
	return md.NodeType
}

func (md *ModifierDefinition) Nodes() []ASTNode {
	return nil
}

func (md *ModifierDefinition) NodeID() int {
	return md.ID
}

func GetModifierDefinition(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*ModifierDefinition, error) {
	md := new(ModifierDefinition)
	if err := json.Unmarshal([]byte(raw.ToString()), md); err != nil {
		logger.Errorf("Failed to unmarshal ModifierDefinition: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal ModifierDefinition: [%v]", err)
	}

	// parameters
	{
		parameters := raw.Get("parameters")
		if parameters.Size() > 0 {
			parametersNodeType := parameters.Get("nodeType").ToString()
			var mdParameters ASTNode
			var err error

			switch parametersNodeType {
			case "ParameterList":
				mdParameters, err = GetParameterList(gn, parameters, logger)
			default:
				logger.Warnf("Unknown parameters nodeType [%s] for ModifierDefinition [src:%s].", parametersNodeType, md.Src)
			}

			if err != nil {
				return nil, err
			}

			if mdParameters != nil {
				md.parameters = mdParameters
			}
		}
	}

	// body
	{
		body := raw.Get("body")
		if body.Size() > 0 {
			bodyNodeType := body.Get("nodeType").ToString()
			var mdBody ASTNode
			var err error

			switch bodyNodeType {
			case "Block":
				mdBody, err = GetBlock(gn, body, logger)
			default:
				logger.Warnf("Unknown body nodeType [%s] for ModifierDefinition [src:%s].", bodyNodeType, md.Src)
			}

			if err != nil {
				return nil, err
			}

			if mdBody != nil {
				md.body = mdBody
			}
		}
	}

	return md, nil
}
