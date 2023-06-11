package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type FunctionDefinition struct {
	body             ASTNode
	ID               int    `json:"id"`
	Implemented      bool   `json:"implemented"`
	Kind             string `json:"kind"`
	modifiers        []ASTNode
	Name             string `json:"name"`
	NodeType         string `json:"nodeType"`
	overrides        ASTNode
	parameters       ASTNode
	returnParameters ASTNode
	Scope            int    `json:"scope"`
	Src              string `json:"src"`
	StateMutability  string `json:"stateMutability"`
	Virtual          bool   `json:"virtual"`
	Visibility       string `json:"visibility"`

	signature string
}

func (fd *FunctionDefinition) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	if fd.Kind == "function" {
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

		// modifiers
		if len(fd.modifiers) > 0 {
			for index, modifier := range fd.modifiers {
				switch m := modifier.(type) {
				case *ModifierInvocation:
					code = code + " " + m.SourceCode(false, false, indent, logger)
				default:
					if m != nil {
						logger.Warnf("Unknown modifier nodeType [%s] for FunctionDefinition [src:%s].", m.Type(), fd.Src)
					} else {
						logger.Warnf("Unknown modifier nodeType for FunctionDefinition [src:%s].", fd.Src)
					}
				}

				if index < len(fd.modifiers)-1 {
					code = code + ","
				}
			}
		}

		// stateMutability
		if fd.StateMutability != "" && fd.StateMutability != "nonpayable" {
			code = code + " " + fd.StateMutability
		}

		// overrides
		if fd.overrides != nil {
			switch overrides := fd.overrides.(type) {
			// case *OverrideSpecifier:
			// 	code = code + " " + overrides.SourceCode(false, false, indent, logger)
			default:
				if overrides != nil {
					logger.Warnf("Unknown overrides nodeType [%s] for FunctionDefinition [src:%s].", overrides.Type(), fd.Src)
				} else {
					logger.Warnf("Unknown overrides nodeType for FunctionDefinition [src:%s].", fd.Src)
				}
			}
		}

		if fd.Virtual {
			code = code + " " + "virtual"
		}

		// returnParameters
		if fd.returnParameters != nil {
			switch returnParameters := fd.returnParameters.(type) {
			case *ParameterList:
				rpl := returnParameters.SourceCode(false, false, indent, logger)
				if rpl != "" {
					code = code + " " + "returns" + " " + "(" + rpl + ")"
				}
			default:
				if returnParameters != nil {
					logger.Warnf("Unknown returnParameters nodeType [%s] for FunctionDefinition [src:%s].", returnParameters.Type(), fd.Src)
				} else {
					logger.Warnf("Unknown returnParameters nodeType for FunctionDefinition [src:%s].", fd.Src)
				}
			}
		}
	} else if fd.Kind == "constructor" {
		code = code + "constructor("

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

		// modifiers
		if len(fd.modifiers) > 0 {
			for index, modifier := range fd.modifiers {
				switch m := modifier.(type) {
				case *ModifierInvocation:
					code = code + " " + m.SourceCode(false, false, indent, logger)
				default:
					if m != nil {
						logger.Warnf("Unknown modifier nodeType [%s] for FunctionDefinition [src:%s].", m.Type(), fd.Src)
					} else {
						logger.Warnf("Unknown modifier nodeType for FunctionDefinition [src:%s].", fd.Src)
					}
				}

				if index < len(fd.modifiers)-1 {
					code = code + ","
				}
			}
		}

		// visibility
		if fd.Visibility != "" {
			code = code + " " + fd.Visibility
		}
		
	} else if fd.Kind == "receive" {
		code = code + "receive("

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
		if fd.StateMutability != "" && fd.StateMutability != "nonpayable" {
			code = code + " " + fd.StateMutability
		}

		code = code + " {}"

		return code
	} else if fd.Kind == "fallback" {
		code = code + "fallback()"

		// visibility
		if fd.Visibility != "" {
			code = code + " " + fd.Visibility
		}

		// stateMutability
		if fd.StateMutability != "" && fd.StateMutability != "nonpayable" {
			code = code + " " + fd.StateMutability
		}
	}

	// body
	{
		if fd.body != nil {
			code = code + " " + "{\n"

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

			code = code + "\n"

			if isIndent {
				code = code + indent
			}

			code = code + "}"
		}
	}

	if !fd.Implemented {
		code = code + ";"
	}

	// logger.Debug(fd.Signature())

	return code
}

func (fd *FunctionDefinition) Type() string {
	return fd.NodeType
}

func (fd *FunctionDefinition) Nodes() []ASTNode {
	return fd.modifiers
}

func (fd *FunctionDefinition) NodeID() int {
	return fd.ID
}

func GetFunctionDefinition(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*FunctionDefinition, error) {
	fd := new(FunctionDefinition)
	if err := json.Unmarshal([]byte(raw.ToString()), fd); err != nil {
		logger.Errorf("Failed to unmarshal FunctionDefinition: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal FunctionDefinition: [%v]", err)
	}

	// modifiers
	{
		modifiers := raw.Get("modifiers")
		if modifiers.Size() > 0 {
			fd.modifiers = make([]ASTNode, 0)

			for i := 0; i < modifiers.Size(); i++ {
				modifier := modifiers.Get(i)
				if modifier.Size() > 0 {
					modifierNodeType := modifier.Get("nodeType").ToString()
					var fdModifier ASTNode
					var err error

					switch modifierNodeType {
					case "ModifierInvocation":
						fdModifier, err = GetModifierInvocation(gn, modifier, logger)
					default:
						logger.Warnf("Unknown modifier nodeType [%s] for FunctionDefinition [src:%s].", modifierNodeType, fd.Src)
					}

					if err != nil {
						return nil, err
					}

					if fdModifier != nil {
						fd.modifiers = append(fd.modifiers, fdModifier)
					}
				} else {
					logger.Warnf("Modifier in FunctionDefinition [src:%s] should not be empty.", fd.Src)
				}
			}
		}
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
				fdBody, err = GetBlock(gn, body, logger)
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
			// case "OverrideSpecifier":
			// 	fdOverrides, err = GetOverrideSpecifier(gn, overrides, logger)
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
				fdParameters, err = GetParameterList(gn, parameters, logger)
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
				fdReturnParameters, err = GetParameterList(gn, returnParameters, logger)
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

	gn.AddASTNode(fd)

	return fd, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (fd *FunctionDefinition) MakeSignature(contractName string, logger logging.Logger) {
	var signature string = contractName
	if fd.Kind == "function" {
		signature = signature + "." + fd.Name
	} else if fd.Kind == "constructor" {
		signature = signature + "." + "constructor"
	} else if fd.Kind == "receive" {
		signature = signature + "." + "receive"
	}

	signature = signature + "("

	if fd.parameters != nil {
		switch parameters := fd.parameters.(type) {
		case *ParameterList:
			signature = signature + parameters.SourceCode(false, false, "", logger)
		default:
			if parameters != nil {
				logger.Warnf("Unknown parameters nodeType [%s] for FunctionDefinition [src:%s].", parameters.Type(), fd.Src)
			} else {
				logger.Warnf("Unknown parameters nodeType for FunctionDefinition [src:%s].", fd.Src)
			}
		}
	}

	signature = signature + ")"

	fd.signature = signature
}

func (fd *FunctionDefinition) Signature() string {
	return fd.signature
}

func (fd *FunctionDefinition) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	ncp.SetCaller(fd.Signature(), fd.NodeID())

	// Function call statements are generally inside functions.
	if fd.body != nil {
		switch body := fd.body.(type) {
		case *Block:
			body.TraverseFunctionCall(ncp, gn, opt, logger)
		}
	}
}

func (fd *FunctionDefinition) TraverseTaintOwner(opt *Option, logger logging.Logger) {
	opt.TrackFunctionDefinitionName = fd.Signature()
	if fd.body != nil {
		switch body := fd.body.(type) {
		case *Block:
			body.TraverseTaintOwner(opt, logger)
		}
	}
}

func (fd *FunctionDefinition) SetBody(body ASTNode) {
	fd.body = body
}

func (fd *FunctionDefinition) SetReturnParameters(returnParameters ASTNode) {
	fd.returnParameters = returnParameters
}

func (fd *FunctionDefinition) TraverseDelegatecall(opt *Option, logger logging.Logger) {
	opt.TrackFunctionDefinitionName = fd.Signature()
	if fd.body != nil {
		switch body := fd.body.(type) {
		case *Block:
			body.TraverseDelegatecall(opt, logger)
		}
	}
}