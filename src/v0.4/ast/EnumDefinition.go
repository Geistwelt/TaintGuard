package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type EnumDefinition struct {
	CanonicalName string `json:"canonicalName"`
	ID            int    `json:"id"`
	members       []ASTNode
	Name          string `json:"name"`
	NodeType      string `json:"nodeType"`
	Src           string `json:"src"`
}

func (ed *EnumDefinition) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "enum " + ed.Name + " {\n"

	if len(ed.members) > 0 {
		for index, member := range ed.members {
			switch m := member.(type) {
			// case *EnumValue:
			// 	code = code + m.SourceCode(false, true, indent+"    ", logger)
			default:
				if m != nil {
					logger.Warnf("Unknown member nodeType [%s] for EnumValue [src:%s].", m.Type(), ed.Src)
				} else {
					logger.Warnf("Unknown member nodeType for EnumValue [src:%s].", ed.Src)
				}
			}

			if index < len(ed.members)-1 {
				code = code + ","
			}
			code = code + "\n"
		}
	}

	if isIndent {
		code = code + indent
	}
	code = code + "}"

	return code
}

func (ed *EnumDefinition) Type() string {
	return ed.NodeType
}

func (ed *EnumDefinition) Nodes() []ASTNode {
	return nil
}

func (ed *EnumDefinition) NodeID() int {
	return ed.ID
}

func GetEnumDefinition(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*EnumDefinition, error) {
	ed := new(EnumDefinition)
	if err := json.Unmarshal([]byte(raw.ToString()), ed); err != nil {
		logger.Errorf("Failed to unmarshal EnumDefinition: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal EnumDefinition: [%v]", err)
	}

	// members
	{
		members := raw.Get("members")
		if members.Size() > 0 {
			ed.members = make([]ASTNode, 0)

			for i := 0; i < members.Size(); i++ {
				member := members.Get(i)
				if member.Size() > 0 {
					memberNodeType := member.Get("nodeType").ToString()
					var edMember ASTNode
					var err error

					switch memberNodeType {
					// case "EnumValue":
					// 	edMember, err = GetEnumValue(gn, member, logger)
					default:
						logger.Warnf("Unknown member nodeType [%s] for EnumDefinition [src:%s].", memberNodeType, ed.Src)
					}

					if err != nil {
						return nil, err
					}

					if edMember != nil {
						ed.members = append(ed.members, edMember)
					}

				} else {
					logger.Warnf("Member in EnumDefinition [%s] should not be empty.", ed.Src)
				}
			}
		}
	}

	return ed, nil
}
