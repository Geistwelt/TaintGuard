package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type ModifierInvocation struct {
	arguments    []ASTNode
	ID           int    `json:"id"`
	Kind         string `json:"kind"`
	modifierName ASTNode
	NodeType     string `json:"nodeType"`
	Src          string `json:"src"`
}

func (mi *ModifierInvocation) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	// modifierName
	if mi.modifierName != nil {
		switch modifierName := mi.modifierName.(type) {
		case *IdentifierPath:
			code = code + modifierName.SourceCode(false, false, indent, logger)
		default:
			if modifierName != nil {
				logger.Warnf("Unknown modifierName nodeType [%s] for ModifierInvocation [src:%s].", modifierName.Type(), mi.Src)
			} else {
				logger.Warnf("Unknown modifierName nodeType for ModifierInvocation [src:%s].", mi.Src)
			}
		}
	}

	// arguments
	if len(mi.arguments) > 0 {
		code = code + "("
		for index, argument := range mi.arguments {
			switch arg := argument.(type) {
			case *Identifier:
				code = code + arg.SourceCode(false, false, indent, logger)
			default:
				if arg != nil {
					logger.Warnf("Unknown argument nodeType [%s] for ModifierInvocation [src:%s].", arg.Type(), mi.Src)
				} else {
					logger.Warnf("Unknown argument nodeType for ModifierInvocation [src:%s].", mi.Src)
				}
			}

			if index < len(mi.arguments)-1 {
				code = code + ", "
			}
		}
		code = code + ")"
	}

	return code
}

func (mi *ModifierInvocation) Type() string {
	return mi.NodeType
}

func (mi *ModifierInvocation) Nodes() []ASTNode {
	return nil
}

func (mi *ModifierInvocation) NodeID() int {
	return mi.ID
}

func GetModifierInvocation(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*ModifierInvocation, error) {
	mi := new(ModifierInvocation)
	if err := json.Unmarshal([]byte(raw.ToString()), mi); err != nil {
		logger.Errorf("Failed to unmarshal ModifierInvocation: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal ModifierInvocation: [%v]", err)
	}

	// arguments
	{
		arguments := raw.Get("arguments")
		if arguments.Size() > 0 {
			mi.arguments = make([]ASTNode, 0)

			for i := 0; i < arguments.Size(); i++ {
				argument := arguments.Get(i)
				if argument.Size() > 0 {
					argumentNodeType := argument.Get("nodeType").ToString()
					var miArgument ASTNode
					var err error

					switch argumentNodeType {
					case "Identifier":
						miArgument, err = GetIdentifier(gn, argument, logger)
					default:
						logger.Warnf("Unknown argument nodeType [%s] for ModifierInvocation [src:%s].", argumentNodeType, mi.Src)
					}

					if err != nil {
						return nil, err
					}

					if miArgument != nil {
						mi.arguments = append(mi.arguments, miArgument)
					}
				} else {
					logger.Warnf("Argument in ModifierInvocation [src:%s] should not be empty.", mi.Src)
				}
			}
		}
	}

	// modifierName
	{
		modifierName := raw.Get("modifierName")
		if modifierName.Size() > 0 {
			modifierNameNodeType := modifierName.Get("nodeType").ToString()
			var miModifierName ASTNode
			var err error

			switch modifierNameNodeType {
			case "IdentifierPath":
				miModifierName, err = GetIdentifierPath(gn, modifierName, logger)
			default:
				logger.Warnf("Unknown modifierName nodeType [%s] for ModifierInvocation [src:%s].", modifierNameNodeType, mi.Src)
			}

			if err != nil {
				return nil, err
			}

			if miModifierName != nil {
				mi.modifierName = miModifierName
			}
		}
	}

	return mi, nil
}
