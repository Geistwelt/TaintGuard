package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type Mapping struct {
	ID               int    `json:"id"`
	KeyName          string `json:"keyName"`
	keyType          ASTNode
	NodeType         string `json:"nodeType"`
	Src              string `json:"src"`
	TypeDescriptions struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
	ValueName         string `json:"valueName"`
	ValueNameLocation string `json:"valueNameLocation"`
	valueType         ASTNode
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

func GetMapping(raw jsoniter.Any, logger logging.Logger) (*Mapping, error) {
	m := new(Mapping)
	if err := json.Unmarshal([]byte(raw.ToString()), m); err != nil {
		logger.Errorf("Failed to unmarshal Mapping: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal Mapping: [%v]", err)
	}

	// keyType
	{
		keyType := raw.Get("keyType")
		keyTypeNodeType := keyType.Get("nodeType").ToString()
		switch keyTypeNodeType {
		case "ElementaryTypeName":
			etn, err := GetElementaryTypeName(keyType, logger)
			if err != nil {
				return nil, err
			}
			m.keyType = etn
		default:
			logger.Warnf("Unknown keyType nodeType [%s] for Mapping [src:%s].", keyTypeNodeType, m.Src)
		}
	}

	// valueType
	{
		valueType := raw.Get("valueType")
		valueTypeNodeType := valueType.Get("nodeType").ToString()
		switch valueTypeNodeType {
		case "ElementaryTypeName":
			etn, err := GetElementaryTypeName(valueType, logger)
			if err != nil {
				return nil, err
			}
			m.valueType = etn
		default:
			logger.Warnf("Unknown valueType nodeType [%s] for Mapping [src:%s].", valueTypeNodeType, m.Src)
		}
	}

	return m, nil
}
