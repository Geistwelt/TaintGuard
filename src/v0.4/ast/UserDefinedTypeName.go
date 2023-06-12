package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type UserDefinedTypeName struct {
	ContractScope    int `json:"contractScope"`
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	NodeType              string `json:"nodeType"`
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

	code = code + udtn.Name

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

func (udtn *UserDefinedTypeName) NodeID() int {
	return udtn.ID
}

func GetUserDefinedTypeName(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*UserDefinedTypeName, error) {
	udtn := new(UserDefinedTypeName)
	if err := json.Unmarshal([]byte(raw.ToString()), udtn); err != nil {
		logger.Errorf("Failed to unmarshal UserDefinedTypeName: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal UserDefinedTypeName: [%v]", err)
	}

	gn.AddASTNode(udtn)

	return udtn, nil
}
