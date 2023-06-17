package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type YulTypedName struct {
	Name     string `json:"name"`
	NodeType string `json:"nodeType"`
	Src      string `json:"src"`
	Type_    string `json:"type"`
}

func (ytn *YulTypedName) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + ytn.Name

	return code
}

func (ytn *YulTypedName) Type() string {
	return ytn.NodeType
}

func (ytn *YulTypedName) Nodes() []ASTNode {
	return nil
}

func (ytn *YulTypedName) NodeID() int {
	return -1
}

func GetYulTypedName(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*YulTypedName, error) {
	ytn := new(YulTypedName)
	if err := json.Unmarshal([]byte(raw.ToString()), ytn); err != nil {
		logger.Errorf("Failed to unmarshal YulTypedName: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal YulTypedName: [%v]", err)
	}

	return ytn, nil
}
