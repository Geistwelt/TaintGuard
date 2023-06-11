package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type StructDefinition struct {
	CanonicalName string `json:"canonicalName"`
	ID            int    `json:"id"`
	members       []ASTNode
	Name          string `json:"name"`
	NameLocation  string `json:"nameLocation"`
	NodeType      string `json:"nodeType"`
	Scope         int    `json:"scope"`
	Src           string `json:"src"`
	Visibility    string `json:"visibility"`
}

func (sd *StructDefinition) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "struct"

	code = code + " " + sd.Name + "{\n"

	if len(sd.members) > 0 {
		for _, member := range sd.members {
			switch m := member.(type) {
			case *VariableDeclaration:
				code = code + m.SourceCode(true, true, indent+"    ", logger)
			default:
				if m != nil {
					logger.Warnf("Unknown member nodeType [%s] for StructDefinition [src:%s].", m.Type(), sd.Src)
				} else {
					logger.Warnf("Unknown member nodeType for StructDefinition [src:%s].", sd.Src)
				}
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

func (sd *StructDefinition) Type() string {
	return sd.NodeType
}

func (sd *StructDefinition) Nodes() []ASTNode {
	return sd.members
}

func (sd *StructDefinition) NodeID() int {
	return sd.ID
}

func GetStructDefinition(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*StructDefinition, error) {
	sd := new(StructDefinition)
	if err := json.Unmarshal([]byte(raw.ToString()), sd); err != nil {
		logger.Errorf("Failed to unmarshal StructDefinition: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal StructDefinition: [%v]", err)
	}

	// members
	{
		members := raw.Get("members")
		if members.Size() > 0 {
			sd.members = make([]ASTNode, 0)

			for i := 0; i< members.Size(); i++ {
				member := members.Get(i)
				if member.Size() > 0 {
					memberNodeType := member.Get("nodeType").ToString()
					var sdMember ASTNode
					var err error

					switch memberNodeType {
					case "VariableDeclaration":
						sdMember, err = GetVariableDeclaration(gn, member, logger)
					default:
						logger.Warnf("Unknown member nodeType [%s] for StructDefinition [src:%s].", memberNodeType, sd.Src)
					}

					if err != nil {
						return nil, err
					}

					if sdMember != nil {
						sd.members = append(sd.members, sdMember)
					}
				}
			}
		}
	}

	gn.AddASTNode(sd)

	return sd, nil
}
