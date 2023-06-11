package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type Break struct {
	ID int `json:"id"`
	NodeType string `json:"nodeType"`
	Src string `json:"src"`
}

func (b *Break) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "break"

	if isSc {
		code = code + ";"
	}

	return code
}

func (b *Break) Type() string {
	return b.NodeType
}

func (b *Break) Nodes() []ASTNode {
	return nil
}

func (b *Break) NodeID() int {
	return b.ID
}

func GetBreak(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*Break, error) {
	b := new(Break)
	if err := json.Unmarshal([]byte(raw.ToString()), b); err != nil {
		logger.Errorf("Failed to unmarshal Break: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal Break: [%v]", err)
	}

	return b, nil
}