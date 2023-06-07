package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type EnumValue struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	NameLocation string `json:"nameLocation"`
	NodeType     string `json:"nodeType"`
	Src          string `json:"src"`
}

func (ev *EnumValue) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + ev.Name

	return code
}

func (ev *EnumValue) Type() string {
	return ev.NodeType
}

func (ev *EnumValue) Nodes() []ASTNode {
	return nil
}

func (ev *EnumValue) NodeID() int {
	return ev.ID
}

func GetEnumValue(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*EnumValue, error) {
	ev := new(EnumValue)
	if err := json.Unmarshal([]byte(raw.ToString()), ev); err != nil {
		logger.Errorf("Failed to unmarshal EnumValue: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal EnumValue: [%v]", err)
	}

	return ev, nil
}
