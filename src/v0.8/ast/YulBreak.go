package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type YulBreak struct {
	NodeType string `json:"nodeType"`
	Src string `json:"src"`
}

func (yb *YulBreak) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "break"

	return code
}

func (yb *YulBreak) Type() string {
	return yb.NodeType
}

func (yb *YulBreak) Nodes() []ASTNode {
	return nil
}

func (yb *YulBreak) NodeID() int {
	return -1
}

func GetYulBreak(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*YulBreak, error) {
	yb := new(YulBreak)
	if err := json.Unmarshal([]byte(raw.ToString()), yb); err != nil {
		logger.Errorf("Failed to unmarshal YulBreak: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal YulBreak: [%v]", err)
	}

	return yb, nil
}