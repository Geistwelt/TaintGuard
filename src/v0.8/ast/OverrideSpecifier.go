package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type OverrideSpecifier struct {
	ID        int      `json:"id"`
	NodeType  string   `json:"nodeType"`
	Overrides []string `json:"overrides"`
	Src       string   `json:"src"`
}

func (os *OverrideSpecifier) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	code = code + "override"

	if isSc {
		code = code + ";"
	}

	return code
}

func (os *OverrideSpecifier) Type() string {
	return os.NodeType
}

func (os *OverrideSpecifier) Nodes() []ASTNode {
	return nil
}

func GetOverrideSpecifier(raw jsoniter.Any, logger logging.Logger) (*OverrideSpecifier, error) {
	os := new(OverrideSpecifier)
	if err := json.Unmarshal([]byte(raw.ToString()), os); err != nil {
		logger.Errorf("Failed to unmarshal OverrideSpecifier: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal OverrideSpecifier: [%v]", err)
	}

	return os, nil
}
