package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type ElementaryTypeName struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	NodeType         string `json:"nodeType"`
	Src              string `json:"src"`
	StateMutability  string `json:"stateMutability"`
	TypeDescriptions struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
}

func (etn *ElementaryTypeName) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + etn.Name
	if etn.StateMutability != "" && etn.StateMutability != "nonpayable" {
		code = code + " " + etn.StateMutability
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (etn *ElementaryTypeName) Type() string {
	return etn.NodeType
}

func (etn *ElementaryTypeName) Nodes() []ASTNode {
	return nil
}

func (etn *ElementaryTypeName) NodeID() int {
	return etn.ID
}

func GetElementaryTypeName(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*ElementaryTypeName, error) {
	etn := new(ElementaryTypeName)
	if err := json.Unmarshal([]byte(raw.ToString()), etn); err != nil {
		logger.Errorf("Failed to unmarshal ElementaryTypeName: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal ElementaryTypeName: [%v]", err)
	}

	gn.AddASTNode(etn)

	return etn, nil
}
