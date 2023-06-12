package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type Mapping struct {
	ID               int `json:"id"`
	keyType          ASTNode
	NodeType         string `json:"nodeType"`
	Src              string `json:"src"`
	TypeDescriptions struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
	valueType ASTNode
}

func (m *Mapping) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string = "mapping"
	if isIndent {
		code = code + indent
	}

	// keyType
	{
		if m.keyType != nil {
			switch keyType := m.keyType.(type) {
			case *ElementaryTypeName:
				code = code + " " + "(" + keyType.SourceCode(false, false, indent, logger)
			default:
				if keyType != nil {
					logger.Warnf("Unknown keyType nodeType [%s] for Mapping [src:%s].", keyType.Type(), m.Src)
				} else {
					logger.Warnf("Unknown keyType nodeType for Mapping [src:%s].", m.Src)
				}
			}
		}
	}

	// valueType
	{
		if m.valueType != nil {
			switch valueType := m.valueType.(type) {
			case *ElementaryTypeName:
				code = code + " " + "=>" + " " + valueType.SourceCode(false, false, indent, logger) + ")"
			case *Mapping:
				code = code + " " + "=>" + " " + valueType.SourceCode(false, false, indent, logger) + ")"
			case *UserDefinedTypeName:
				code = code + " " + "=>" + " " + valueType.SourceCode(false, false, indent, logger) + ")"
			case *ArrayTypeName:
				code = code + " " + "=>" + " " + valueType.SourceCode(false, false, indent, logger) + ")"
			default:
				if valueType != nil {
					logger.Warnf("Unknown valueType nodeType [%s] for Mapping [src:%s].", valueType.Type(), m.Src)
				} else {
					logger.Warnf("Unknown valueType nodeType for Mapping [src:%s].", m.Src)
				}
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (m *Mapping) Type() string {
	return m.NodeType
}

func (m *Mapping) Nodes() []ASTNode {
	return nil
}

func (m *Mapping) NodeID() int {
	return m.ID
}

func GetMapping(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*Mapping, error) {
	m := new(Mapping)
	if err := json.Unmarshal([]byte(raw.ToString()), m); err != nil {
		logger.Errorf("Failed to unmarshal Mapping: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal Mapping: [%v]", err)
	}

	// keyType
	{
		keyType := raw.Get("keyType")
		keyTypeNodeType := keyType.Get("nodeType").ToString()
		var mKeyType ASTNode
		var err error

		switch keyTypeNodeType {
		case "ElementaryTypeName":
			mKeyType, err = GetElementaryTypeName(gn, keyType, logger)
		default:
			logger.Warnf("Unknown keyType nodeType [%s] for Mapping [src:%s].", keyTypeNodeType, m.Src)
		}

		if err != nil {
			return nil, err
		}

		if mKeyType != nil {
			m.keyType = mKeyType
		}
	}

	// valueType
	{
		valueType := raw.Get("valueType")
		valueTypeNodeType := valueType.Get("nodeType").ToString()
		var mValueType ASTNode
		var err error

		switch valueTypeNodeType {
		case "ElementaryTypeName":
			mValueType, err = GetElementaryTypeName(gn, valueType, logger)
		case "Mapping":
			mValueType, err = GetMapping(gn, valueType, logger)
		case "UserDefinedTypeName":
			mValueType, err = GetUserDefinedTypeName(gn, valueType, logger)
		case "ArrayTypeName":
			mValueType, err = GetArrayTypeName(gn, valueType, logger)
		default:
			logger.Warnf("Unknown valueType nodeType [%s] for Mapping [src:%s].", valueTypeNodeType, m.Src)
		}

		if err != nil {
			return nil, err
		}

		if mValueType != nil {
			m.valueType = mValueType
		}
	}

	gn.AddASTNode(m)

	return m, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (m *Mapping) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	// keyType
	{

	}

	// valueType
	{
		if m.valueType != nil {
			switch valueType := m.valueType.(type) {
			case *Mapping:
				valueType.TraverseFunctionCall(ncp, gn, opt, logger)
			}
		}
	}
}

func (m *Mapping) SetKeyType(kt ASTNode) {
	m.keyType = kt
}

func (m *Mapping) SetValueType(vt ASTNode) {
	m.valueType = vt
}
