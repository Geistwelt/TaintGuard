package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type ElementaryTypeNameExpression struct {
	ArgumentTypes []struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"argumentTypes"`
	ID               int    `json:"id"`
	IsConstant       bool   `json:"isConstant"`
	IsValue          bool   `json:"isValue"`
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

func (etne *ElementaryTypeNameExpression) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string 
	if isIndent {
		code = code + indent
	}

	code = code + etne.typeName.SourceCode(false, false, indent, logger)

	if isSc {
		code = code + ";"
	}

	return code
}

func (etne *ElementaryTypeNameExpression) Type() string {
	return etne.NodeType
}

func (etne *ElementaryTypeNameExpression) Nodes() []ASTNode {
	return nil
}

func (etne *ElementaryTypeNameExpression) NodeID() int {
	return etne.ID
}

func GetElementaryTypeNameExpression(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*ElementaryTypeNameExpression, error) {
	etne := new(ElementaryTypeNameExpression)
	if err := json.Unmarshal([]byte(raw.ToString()), etne); err != nil {
		logger.Errorf("Failed to unmarshal ElementaryTypeNameExpression: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal ElementaryTypeNameExpression: [%v]", err)
	}

	// typeName
	{
		typeName := raw.Get("typeName")
		if typeName.Size() > 0 {
			typeNameNodeType := typeName.Get("nodeType").ToString()
			var etneTypeName ASTNode
			var err error

			switch typeNameNodeType {
			case "ElementaryTypeName":
				etneTypeName, err = GetElementaryTypeName(gn, typeName, logger)
			default:
				logger.Warnf("Unknown typeName nodeType [%s] for ElementaryTypeNameExpression [src:%s].", typeNameNodeType, etne.Src)
			}

			if err != nil {
				return nil, err
			}

			if etneTypeName != nil {
				etne.typeName = etneTypeName
			}
		}
	}

	gn.AddASTNode(etne)

	return etne, nil
}
 