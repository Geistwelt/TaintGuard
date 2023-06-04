package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type Literal struct {
	HexValue         string `json:"hexValue"`
	ID               int    `json:"id"`
	IsConstant       bool   `json:"isConstant"`
	IsValue          bool   `json:"isValue"`
	IsPure           bool   `json:"isPure"`
	Kind             string `json:"kind"`
	LValueRequested  bool   `json:"lValueRequested"`
	NodeType         string `json:"nodeType"`
	Src              string `json:"src"`
	TypeDescriptions struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
	Value string `json:"value"`
}

func (l *Literal) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	switch l.Kind {
	case "number", "bool":
		code = l.Value
	case "unicodeString":
		code = code + "unicode" + fmt.Sprintf("\"%s\"", l.Value)
	default:
		logger.Warnf("Unknown kind [%s] for Literal [src:%s].", l.Kind, l.Src)
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (l *Literal) Type() string {
	return l.NodeType
}

func (l *Literal) Nodes() []ASTNode {
	return nil
}

func GetLiteral(raw jsoniter.Any, logger logging.Logger) (*Literal, error) {
	l := new(Literal)
	if err := json.Unmarshal([]byte(raw.ToString()), l); err != nil {
		logger.Errorf("Failed to unmarshal for Literal: [%s].", err)
		return nil, fmt.Errorf("failed to unmarshal for Literal: [%s]", err)
	}

	return l, nil
}