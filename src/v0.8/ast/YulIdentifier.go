package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type YulIdentifier struct {
	Name     string `json:"name"`
	NodeType string `json:"nodeType"`
	Src      string `json:"src"`
}

func (yi *YulIdentifier) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	code = code + yi.Name

	return code
}

func (yi *YulIdentifier) Type() string {
	return yi.NodeType
}

func (yi *YulIdentifier) Nodes() []ASTNode {
	return nil
}

func GetYulIdentifier(raw jsoniter.Any, logger logging.Logger) (*YulIdentifier, error) {
	yi := new(YulIdentifier)
	if err := json.Unmarshal([]byte(raw.ToString()), yi); err != nil {
		logger.Errorf("Failed to unmarshal YulIdentifier: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal YulIdentifier: [%v]", err)
	}

	return yi, nil
}
