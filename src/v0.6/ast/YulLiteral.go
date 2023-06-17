package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type YulLiteral struct {
	Kind     string `json:"kind"`
	NodeType string `json:"nodeType"`
	Src      string `json:"src"`
	Type_    string `json:"type"`
	Value    string `json:"value"`
}

func (yl *YulLiteral) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	switch yl.Kind {
	case "number":
		code = code + yl.Value
	default:
		code = code + "HHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH"
	}

	return code
}

func (yl *YulLiteral) Type() string {
	return yl.NodeType
}

func (yl *YulLiteral) Nodes() []ASTNode {
	return nil
}

func (yl *YulLiteral) NodeID() int {
	return -1
}

func GetYulLiteral(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*YulLiteral, error) {
	yl := new(YulLiteral)
	if err := json.Unmarshal([]byte(raw.ToString()), yl); err != nil {
		logger.Errorf("Failed to unmarshal YulLiteral: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal YulLiteral: [%v]", err)
	}

	return yl, nil
}
