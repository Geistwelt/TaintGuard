package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type UsingForDirective struct {
	ID          int `json:"id"`
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
		case *UserDefinedTypeName:
			code = code + " " + libraryNameType.SourceCode(false, false, indent, logger)
		default:
			logger.Warnf("Unknown libraryName nodeType [%s] for UsingForDirective [src:%s].", libraryNameType.Type(), ufd.Src)
		}
	}

	code = code + " " + "for"

	switch typeNameType := ufd.typeName.(type) {
	case *ElementaryTypeName:
		code = code + " " + typeNameType.SourceCode(false, false, indent, logger)
	case *UserDefinedTypeName:
		code = code + " " + typeNameType.SourceCode(false, false, indent, logger)
	default:
		if typeNameType != nil {
			logger.Warnf("Unknown typeName nodeType [%s] for UsingForDirective [src:%s].", typeNameType.Type(), ufd.Src)
		} else {
			logger.Warnf("Unknown typeName nodeType for UsingForDirective [src:%s].", ufd.Src)
		}
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

func (ufd *UsingForDirective) NodeID() int {
	return ufd.ID
}

func GetUsingForDirective(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*UsingForDirective, error) {
	ufd := new(UsingForDirective)

	if err := json.Unmarshal([]byte(raw.ToString()), ufd); err != nil {
		logger.Errorf("Failed to unmarshal UsingForDirective: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal UsingForDirective: [%v]", err)
	}

	// libraryName
	{
		libraryName := raw.Get("libraryName")
		switch libraryName.Get("nodeType").ToString() {
		case "UserDefinedTypeName":
			ip, err := GetUserDefinedTypeName(gn, libraryName, logger)
			if err != nil {
				return nil, err
			}
			ufd.libraryName = ip
		default:
			logger.Warnf("Unknown libraryName nodeType [%s] for UsingForDirective [src:%s].", libraryName.Get("nodeType").ToString(), ufd.Src)
		}
	}

	// typeName
	{
		typeName := raw.Get("typeName")
		var ufdTypeName ASTNode
		var err error

		switch typeName.Get("nodeType").ToString() {
		case "ElementaryTypeName":
			ufdTypeName, err = GetElementaryTypeName(gn, typeName, logger)
		case "UserDefinedTypeName":
			ufdTypeName, err = GetUserDefinedTypeName(gn, typeName, logger)
		default:
			logger.Warnf("Unknown typeName nodeType [%s] for UsingForDirective [src:%s].", typeName.Get("nodeType").ToString(), ufd.Src)
		}

		if err != nil {
			return nil, err
		}

		if ufdTypeName != nil {
			ufd.typeName = ufdTypeName
		}
	}

	gn.AddASTNode(ufd)

	return ufd, nil
}
