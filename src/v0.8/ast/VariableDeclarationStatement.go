package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type VariableDeclarationStatement struct {
	Assignments  []int `json:"assignments"`
	declarations []ASTNode
	ID           int `json:"id"`
	initialValue ASTNode
	NodeType     string `json:"nodeType"`
	Src          string `json:"src"`
}

func (vds *VariableDeclarationStatement) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	if len(vds.declarations) > 0 {
		for index, declaration := range vds.declarations {
			switch d := declaration.(type) {
			case *VariableDeclaration:
				code = code + d.SourceCode(false, false, indent, logger)
			default:
				if d != nil {
					logger.Warnf("Unknown declaration nodeType [%s] for VariableDeclarationStatement [src:%s]", d.Type(), vds.Src)
				} else {
					logger.Warnf("Unknown declaration nodeType for VariableDeclarationStatement [src:%s]", vds.Src)
				}
			}

			if index < len(vds.declarations)-1 {
				code = code + ","
			}
		}
	}

	if vds.initialValue != nil {
		switch initialValue := vds.initialValue.(type) {
		case *FunctionCall:
			code = code + " = " + initialValue.SourceCode(false, false, indent, logger)
		case *MemberAccess:
			code = code + " = " + initialValue.SourceCode(false, false, indent, logger)
		case *BinaryOperation:
			code = code + " = " + initialValue.SourceCode(false, false, indent, logger)
		case *Literal:
			code = code + " = " + initialValue.SourceCode(false, false, indent, logger)
		default:
			if initialValue != nil {
				logger.Warnf("Unknown initialValue nodeType [%s] for VariableDeclarationStatement [src:%s]", initialValue.Type(), vds.Src)
			} else {
				logger.Warnf("Unknown initialValue nodeType for VariableDeclarationStatement [src:%s]", vds.Src)
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (vds *VariableDeclarationStatement) Type() string {
	return vds.NodeType
}

func (vds *VariableDeclarationStatement) Nodes() []ASTNode {
	return vds.declarations
}

func GetVariableDeclarationStatement(raw jsoniter.Any, logger logging.Logger) (*VariableDeclarationStatement, error) {
	vds := new(VariableDeclarationStatement)
	if err := json.Unmarshal([]byte(raw.ToString()), vds); err != nil {
		logger.Errorf("Failed to unmarshal VariableDeclarationStatement: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal VariableDeclarationStatement: [%v]", err)
	}

	// declarations
	{
		declarations := raw.Get("declarations")
		if declarations.Size() > 0 {
			vds.declarations = make([]ASTNode, 0)

			for i := 0; i < declarations.Size(); i++ {
				declaration := declarations.Get(i)
				if declaration.Size() > 0 {
					var declarationNodeType = declaration.Get("nodeType").ToString()
					var err error
					var vdsDeclaration ASTNode
	
					switch declarationNodeType {
					case "VariableDeclaration":
						vdsDeclaration, err = GetVariableDeclaration(declaration, logger)
					default:
						logger.Warnf("Unknown declaration nodeType [%s] for VariableDeclarationStatement [src:%s].", declarationNodeType, vds.Src)
					}
	
					if err != nil {
						return nil, err
					}
	
					if vdsDeclaration != nil {
						vds.declarations = append(vds.declarations, vdsDeclaration)
					}
				} else {
					logger.Warnf("Declaration in VariableDeclarationStatement [src:%s] should not be nil.", vds.Src)
				}
			}
		}
	}

	// initialValue
	{
		initialValue := raw.Get("initialValue")
		if initialValue.Size() > 0 {
			initialValueNodeType := initialValue.Get("nodeType").ToString()
			var vdsInitialValue ASTNode
			var err error

			switch initialValueNodeType {
			case "FunctionCall":
				vdsInitialValue, err = GetFunctionCall(initialValue, logger)
			case "MemberAccess":
				vdsInitialValue, err = GetMemberAccess(initialValue, logger)
			case "BinaryOperation":
				vdsInitialValue, err = GetBinaryOperation(initialValue, logger)
			case "Literal":
				vdsInitialValue, err = GetLiteral(initialValue, logger)
			default:
				logger.Warnf("Unknown initialValue nodeType [%s] for VariableDeclarationStatement [src:%s].", initialValueNodeType, vds.Src)
			}

			if err != nil {
				return nil, err
			}

			if vdsInitialValue != nil {
				vds.initialValue = vdsInitialValue
			}
		}
	}

	return vds, nil
}
