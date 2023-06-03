package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type Identifier struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	NodeType               string `json:"nodeType"`
	OverloadedDeclarations []int  `json:"overloadedDeclarations"`
	ReferencedDeclaration  int    `json:"referencedDeclaration"` // 被引用的变量的声明语句ID
	Src                    string `json:"src"`
	TypeDescriptions       struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
}

func (i *Identifier) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}
	code = code + i.Name
	if isSc {
		code = code + ";"
	}
	return code
}

func (i *Identifier) Type() string {
	return i.NodeType
}

func (i *Identifier) Nodes() []ASTNode {
	return nil
}

func GetIdentifier(raw jsoniter.Any, logger logging.Logger) (*Identifier, error) {
	i := new(Identifier)
	if err := json.Unmarshal([]byte(raw.ToString()), i); err != nil {
		logger.Errorf("Failed to unmarshal for Identifier: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal for Identifier: [%v]", err)
	}
	return i, nil
}
