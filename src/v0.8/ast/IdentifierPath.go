package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type IdentifierPath struct {
	ID                    int      `json:"id"`
	Name                  string   `json:"name"`
	NameLocations         []string `json:"nameLocations"`
	NodeType              string   `json:"nodeType"`
	ReferencedDeclaration int      `json:"referencedDeclaration"`
}

func (ip *IdentifierPath) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + ip.Name

	if isSc {
		code = code + ";"
	}
	return code
}

func (ip *IdentifierPath) Type() string {
	return ip.NodeType
}

func (ip *IdentifierPath) Nodes() []ASTNode {
	return nil
}

func (ip *IdentifierPath) NodeID() int {
	return ip.ID
}


func GetIdentifierPath(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*IdentifierPath, error) {
	ip := new(IdentifierPath)
	if err := json.Unmarshal([]byte(raw.ToString()), ip); err != nil {
		logger.Errorf("Failed to unmarshal IdentifierPath: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal IdentifierPath: [%v]", err)
	}

	gn.AddASTNode(ip)
	
	return ip, nil
}
