package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type NewExpression struct {
	ArgumentTypes []struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"argumentTypes"`
	ID               int    `json:"id"`
	IsConstant       bool   `json:"isConstant"`
	IsLValue         bool   `json:"isLValue"`
	IsPure           bool   `json:"isPure"`
	LValueRequested  bool   `json:"lValueRequested"`
	NodeType         string `json:"nodeType"`
	Src              string `json:"src"`
	TypeDescriptions struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
	typeName ASTNode
}

func (ne *NewExpression) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	code = code + "new"

	if ne.typeName != nil {
		switch typeName := ne.typeName.(type) {
		case *ArrayTypeName:
			code = code + " " + typeName.SourceCode(false, false, indent, logger)
		case *ElementaryTypeName:
			code = code + " " + typeName.SourceCode(false, false, indent, logger)
		default:
			if typeName != nil {
				logger.Warnf("Unknown typeName nodeType [%s] for NewExpression [src:%s].", typeName.Type(), ne.Src)
			} else {
				logger.Warnf("Unknown typeName nodeType for NewExpression [src:%s].", ne.Src)
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (ne *NewExpression) Type() string {
	return ne.NodeType
}

func (ne *NewExpression) Nodes() []ASTNode {
	return nil
}

func (ne *NewExpression) NodeID() int {
	return ne.ID
}

func GetNewExpression(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*NewExpression, error) {
	ne := new(NewExpression)
	if err := json.Unmarshal([]byte(raw.ToString()), ne); err != nil {
		logger.Errorf("Failed to unmarshal NewExpression: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal NewExpression: [%v]", err)
	}

	// typeName
	{
		typeName := raw.Get("typeName")
		if typeName.Size() > 0 {
			typeNameNodeType := typeName.Get("nodeType").ToString()
			var neTypeName ASTNode
			var err error

			switch typeNameNodeType {
			case "ArrayTypeName":
				neTypeName, err = GetArrayTypeName(gn, typeName, logger)
			case "ElementaryTypeName":
				neTypeName, err = GetElementaryTypeName(gn, typeName, logger)
			default:
				logger.Warnf("Unknown typeName nodeType [%s] for NewExpression [src:%s].", typeNameNodeType, ne.Src)
			}

			if err != nil {
				return nil, err
			}

			if neTypeName != nil {
				ne.typeName = neTypeName
			}
		}
	}

	gn.AddASTNode(ne)

	return ne, nil
}
