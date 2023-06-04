package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type FunctionDefinition struct {
	BaseFunctions    []int `json:"baseFunctions"`
	body             ASTNode
	FunctionSelector string `json:"functionSelector"`
	ID               int    `json:"id"`
	Implemented      bool   `json:"implemented"`
	Kind             string `json:"kind"`
	modifiers        []ASTNode
	Name             string `json:"name"`
	NameLocation     string `json:"nameLocation"`
	NodeType         string `json:"nodeType"`
	overrides        ASTNode
	parameters       ASTNode
	returnParameters ASTNode
	Scope            int    `json:"scope"`
	Src              string `json:"src"`
	StateMutability  string `json:"stateMutability"`
	Virtual          bool   `json:"virtual"`
	Visibility       string `json:"visibility"`
}

func (fd *FunctionDefinition) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	code = code + "function" + " " + fd.Name + "("

	// parameters
	{
		if fd.parameters != nil {
			switch parameters := fd.parameters.(type) {
			case *ParameterList:
				code = code + parameters.SourceCode(false, false, indent, logger)
			default:
				if parameters != nil {
					logger.Warnf("Unknown parameters nodeType [%s] for FunctionDefinition [src:%s].", parameters.Type(), fd.Src)
				} else {
					logger.Warnf("Unknown parameters nodeType for FunctionDefinition [src:%s].", fd.Src)
				}
			}
		}
	}

	code = code + ")"

	// visibility
	if fd.Visibility != "" {
		code = code + " " + fd.Visibility
	}

	// stateMutability
	if fd.StateMutability != "" {
		code = code + " " + fd.StateMutability
	}

	// overrides
	if fd.overrides != nil {
		switch overrides := fd.overrides.(type) {
		case *OverrideSpecifier:
			code = code + " " + overrides.SourceCode(false, false, indent, logger)
		default:
			if overrides != nil {
				logger.Warnf("Unknown overrides nodeType [%s] for FunctionDefinition [src:%s].", overrides.Type(), fd.Src)
			} else {
				logger.Warnf("Unknown overrides nodeType for FunctionDefinition [src:%s].", fd.Src)
			}
		}
	}

	// returnParameters
	if fd.returnParameters != nil {
		switch returnParameters := fd.returnParameters.(type) {
		case *ParameterList:
			rpl := returnParameters.SourceCode(false, false, indent, logger)
			if rpl != "" {
				code = code + " " + "retruns" + " " + "(" + rpl + ")"
			}
		default:
			if returnParameters != nil {
				logger.Warnf("Unknown returnParameters nodeType [%s] for FunctionDefinition [src:%s].", returnParameters.Type(), fd.Src)
			} else {
				logger.Warnf("Unknown returnParameters nodeType for FunctionDefinition [src:%s].", fd.Src)
			}
		}
	}

	code = code + " " + "{\n"

	// body
	{
		if fd.body != nil {
			switch body := fd.body.(type) {
			case *Block:
				code = code + body.SourceCode(false, false, indent, logger)
			default:
				if body != nil {
					logger.Warnf("Unknown body nodeType [%s] for FunctionDefinition [src:%s].", body.Type(), fd.Src)
				} else {
					logger.Warnf("Unknown body nodeType for FunctionDefinition [src:%s].", fd.Src)
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

func (fd *FunctionDefinition) Type() string {
	return fd.NodeType
}

func (fd *FunctionDefinition) Nodes() []ASTNode {
	return fd.modifiers
}

func GetFunctionDefinition(raw jsoniter.Any, logger logging.Logger) (*FunctionDefinition, error) {
	fd := new(FunctionDefinition)
	if err := json.Unmarshal([]byte(raw.ToString()), fd); err != nil {
		logger.Errorf("Failed to unmarshal FunctionDefinition: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal FunctionDefinition: [%v]", err)
	}

	// body
	{
		body := raw.Get("body")
		if body.Size() > 0 {
			bodyNodeType := body.Get("nodeType").ToString()
			var fdBody ASTNode
			var err error

			switch bodyNodeType {
			case "Block":
				fdBody, err = GetBlock(body, logger)
			default:
				logger.Warnf("Unknown body nodeType [%s] for FunctionDefinition [src:%s].", bodyNodeType, fd.Src)
			}

			if err != nil {
				return nil, err
			}

			if fdBody != nil {
				fd.body = fdBody
			}
		}
	}

	// overrides
	{
		overrides := raw.Get("overrides")
		if overrides.Size() > 0 {
			overridesNodeType := overrides.Get("nodeType").ToString()
			var fdOverrides ASTNode
			var err error

			switch overridesNodeType {
			case "OverrideSpecifier":
				fdOverrides, err = GetOverrideSpecifier(overrides, logger)
			default:
				logger.Warnf("Unknown overrides nodeType [%s] for FunctionDefinition [src:%s].", overridesNodeType, fd.Src)
			}

			if err != nil {
				return nil, err
			}

			if fdOverrides != nil {
				fd.overrides = fdOverrides
			}
		}
	}

	// parameters
	{
		parameters := raw.Get("parameters")
		if parameters.Size() > 0 {
			parametersNodeType := parameters.Get("nodeType").ToString()
			var fdParameters ASTNode
			var err error

			switch parametersNodeType {
			case "ParameterList":
				fdParameters, err = GetParameterList(parameters, logger)
			default:
				logger.Warnf("Unknown parameters nodeType [%s] for FunctionDefinition [src:%s].", parametersNodeType, fd.Src)
			}

			if err != nil {
				return nil, err
			}

			if fdParameters != nil {
				fd.parameters = fdParameters
			}
		}
	}

	// returnParameters
	{
		returnParameters := raw.Get("returnParameters")
		if returnParameters.Size() > 0 {
			returnParametersNodeType := returnParameters.Get("nodeType").ToString()
			var fdReturnParameters ASTNode
			var err error

			switch returnParametersNodeType {
			case "ParameterList":
				fdReturnParameters, err = GetParameterList(returnParameters, logger)
			default:
				logger.Warnf("Unknown returnParameters nodeType [%s] for FunctionDefinition [src:%s].", returnParametersNodeType, fd.Src)
			}

			if err != nil {
				return nil, err
			}

			if fdReturnParameters != nil {
				fd.returnParameters = fdReturnParameters
			}
		}
	}

	return fd, nil
}