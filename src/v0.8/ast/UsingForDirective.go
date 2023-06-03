package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type UsingForDirective struct {
	Global      bool `json:"global"`
	ID          int  `json:"id"`
	libraryName ASTNode
	NodeType    string `json:"nodeType"`
	Src         string `json:"src"`
	typeName    ASTNode
}

func (ufd *UsingForDirective) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}
	code = code + "using"

	if ufd.libraryName != nil {
		switch libraryNameType := ufd.libraryName.(type) {
		case *IdentifierPath:
			code = code + " " + libraryNameType.SourceCode(false, false, indent, logger)
		default:
			logger.Errorf("Unknown libraryName nodeType [%s] for UsingForDirective [src:%s].", libraryNameType.Type(), ufd.Src)
		}
	}

	code = code + " " + "for"

	switch typeNameType := ufd.typeName.(type) {
	case *ElementaryTypeName:
		code = code + " " + typeNameType.SourceCode(false, false, indent, logger)
	default:
		logger.Errorf("Unknown typeName nodeType [%s] for UsingForDirective [src:%s].", typeNameType.Type(), ufd.Src)
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (ufd *UsingForDirective) Type() string {
	return ufd.NodeType
}

func (ufd *UsingForDirective) Nodes() []ASTNode {
	return nil
}

func GetUsingForDirective(raw jsoniter.Any, logger logging.Logger) (*UsingForDirective, error) {
	ufd := new(UsingForDirective)

	if err := json.Unmarshal([]byte(raw.ToString()), ufd); err != nil {
		logger.Errorf("Failed to unmarshal UsingForDirective: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal UsingForDirective: [%v]", err)
	}

	// libraryName
	{
		libraryName := raw.Get("libraryName")
		switch libraryName.Get("nodeType").ToString() {
		case "IdentifierPath":
			ip, err := GetIdentifierPath(libraryName, logger)
			if err != nil {
				return nil, err
			}
			ufd.libraryName = ip
		default:
			logger.Errorf("Unknown libraryName nodeType [%s] for UsingForDirective [src:%s].", libraryName.Get("nodeType").ToString(), ufd.Src)
		}
	}

	// typeName
	{
		typeName := raw.Get("typeName")
		switch typeName.Get("nodeType").ToString() {
		case "ElementaryTypeName":
			etn, err := GetElementaryTypeName(typeName, logger)
			if err != nil {
				return nil, err
			}
			ufd.typeName = etn
		default:
			logger.Errorf("Unknown typeName nodeType [%s] for UsingForDirective [src:%s].", typeName.Get("nodeType").ToString(), ufd.Src)
		}
	}


	return ufd, nil
}
