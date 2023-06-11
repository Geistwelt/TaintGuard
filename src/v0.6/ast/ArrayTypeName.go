package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type ArrayTypeName struct {
	baseType         ASTNode
	ID               int `json:"id"`
	length           ASTNode
	NodeType         string `json:"nodeType"`
	Src              string `json:"src"`
	TypeDescriptions struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
}

func (atn *ArrayTypeName) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	if atn.baseType != nil {
		switch baseType := atn.baseType.(type) {
		case *ElementaryTypeName:
			code = code + baseType.SourceCode(false, false, indent, logger)
		case *UserDefinedTypeName:
			code = code + baseType.SourceCode(false, false, indent, logger)
		default:
			if baseType != nil {
				logger.Warnf("Unknown baseType nodeType [%s] for ArrayTypeName [src:%s].", baseType.Type(), atn.Src)
			} else {
				logger.Warnf("Unknown baseType nodeType for ArrayTypeName [src:%s].", atn.Src)
			}
		}
	}

	code = code + "["

	if atn.length != nil {
		switch length := atn.length.(type) {
		case *Literal:
			code = code + length.SourceCode(false, false, indent, logger)
		default:
			if length != nil {
				logger.Warnf("Unknown length nodeType [%s] for ArrayTypeName [src:%s].", length.Type(), atn.Src)
			} else {
				logger.Warnf("Unknown length nodeType for ArrayTypeName [src:%s].", atn.Src)
			}
		}
	}

	code = code + "]"

	return code
}

func (atn *ArrayTypeName) Type() string {
	return atn.NodeType
}

func (atn *ArrayTypeName) Nodes() []ASTNode {
	return nil
}

func (atn *ArrayTypeName) NodeID() int {
	return atn.ID
}

func GetArrayTypeName(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*ArrayTypeName, error) {
	atn := new(ArrayTypeName)
	if err := json.Unmarshal([]byte(raw.ToString()), atn); err != nil {
		logger.Errorf("Failed to unmarshal ArrayTypeName: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal ArrayTypeName: [%v]", err)
	}

	// baseType
	{
		baseType := raw.Get("baseType")
		if baseType.Size() > 0 {
			baseTypeNodeType := baseType.Get("nodeType").ToString()
			var atnBaseType ASTNode
			var err error

			switch baseTypeNodeType {
			case "ElementaryTypeName":
				atnBaseType, err = GetElementaryTypeName(gn, baseType, logger)
			case "UserDefinedTypeName":
				atnBaseType, err = GetUserDefinedTypeName(gn, baseType, logger)
			default:
				logger.Warnf("Unknown baseType nodeType [%s] for ArrayTypeName [src:%s].", baseTypeNodeType, atn.Src)
			}

			if err != nil {
				return nil, err
			}

			if atnBaseType != nil {
				atn.baseType = atnBaseType
			}
		}
	}

	// length
	{
		length := raw.Get("length")
		if length.Size() > 0 {
			lengthTypeNodeType := length.Get("nodeType").ToString()
			var atnLength ASTNode
			var err error

			switch lengthTypeNodeType {
			case "Literal":
				atnLength, err = GetLiteral(gn, length, logger)
			default:
				logger.Warnf("Unknown length nodeType [%s] for ArrayTypeName [src:%s].", lengthTypeNodeType, atn.Src)
			}

			if err != nil {
				return nil, err
			}

			if atnLength != nil {
				atn.length = atnLength
			}
		}
	}

	gn.AddASTNode(atn)

	return atn, nil
}
