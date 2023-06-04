package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type PlaceholderStatement struct {
	ID       int    `json:"id"`
	NodeType string `json:"nodeType"`
	Src      string `json:"src"`
}

func (ps *PlaceholderStatement) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	code = code + "_"

	if isSc {
		code = code + ";"
	}

	return code
}

func (ps *PlaceholderStatement) Type() string {
	return ps.NodeType
}

func (ps *PlaceholderStatement) Nodes() []ASTNode {
	return nil
}

func GetPlaceholderStatement(raw jsoniter.Any, logger logging.Logger) (*PlaceholderStatement, error) {
	ps := new(PlaceholderStatement)
	if err := json.Unmarshal([]byte(raw.ToString()), ps); err != nil {
		logger.Errorf("Failed to unmarshal PlaceholderStatement: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal PlaceholderStatement: [%v]", err)
	}

	return ps, nil
}
