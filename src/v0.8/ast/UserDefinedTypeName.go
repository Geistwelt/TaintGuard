package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type UserDefinedTypeName struct {
	ID                    int    `json:"id"`
	NodeType              string `json:"nodeType"`
	pathNode              ASTNode
	ReferencedDeclaration int    `json:"referencedDeclaration"`
	Src                   string `json:"src"`
	TypeDescriptions      struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
}

func (udtn *UserDefinedTypeName) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	if udtn.pathNode != nil {
		switch pathNode := udtn.pathNode.(type) {
		case *IdentifierPath:
			code = code + pathNode.SourceCode(false, false, indent, logger)
		default:
			if pathNode != nil {
				logger.Errorf("Unknown pathNode nodeType [%s] for UserDefinedTypeName [src:%s].", pathNode.Type(), udtn.Src)
			} else {
				logger.Errorf("Unknown pathNode nodeType for UserDefinedTypeName [src:%s].", udtn.Src)
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (udtn *UserDefinedTypeName) Type() string {
	return udtn.NodeType
}

func (udtn *UserDefinedTypeName) Nodes() []ASTNode {
	return nil
}

func GetUserDefinedTypeName(raw jsoniter.Any, logger logging.Logger) (*UserDefinedTypeName, error) {
	udtn := new(UserDefinedTypeName)
	if err := json.Unmarshal([]byte(raw.ToString()), udtn); err != nil {
		logger.Errorf("Failed to unmarshal UserDefinedTypeName: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal UserDefinedTypeName: [%v]", err)
	}

	// pathNode
	{
		pathNode := raw.Get("pathNode")
		if pathNode.Size() > 0 {
			pathNodeNodeType := pathNode.Get("nodeType").ToString()
			var udtnPathNode ASTNode
			var err error

			switch pathNodeNodeType {
			case "IdentifierPath":
				udtnPathNode, err = GetIdentifierPath(pathNode, logger)
			default:
				logger.Errorf("Unknown pathNode nodeType [%s] for UserDefinedTypeName [src:%s].", pathNodeNodeType, udtn.Src)
			}

			if err != nil {
				return nil, err
			}
			if udtnPathNode != nil {
				udtn.pathNode = udtnPathNode
			}
		}
	}

	return udtn, nil
}
