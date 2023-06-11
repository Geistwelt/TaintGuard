package ast

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type InlineAssembly struct {
	ExternalReferences []map[string]struct {
		Declaration int    `json:"declaration"`
		IsOffset    bool   `json:"isOffset"`
		IsSlot      bool   `json:"isSlot"`
		Src         string `json:"src"`
		ValueSize   int    `json:"valueSize"`
	} `json:"externalReferences"`
	ID         int    `json:"id"`
	NodeType   string `json:"nodeType"`
	Operations string `json:"operations"`
	Src        string `json:"src"`
}

func (ia *InlineAssembly) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "assembly {\n"

	operations := strings.Split(ia.Operations, "\n")
	for i := 1; i < len(operations); i++ {
		if isIndent {
			code = code + indent + operations[i]
		}
		if i < len(operations)-1 {
			code = code +"\n"
		}
	}

	return code
}

func (ia *InlineAssembly) Type() string {
	return ia.NodeType
}

func (ia *InlineAssembly) Nodes() []ASTNode {
	return nil
}

func (ia *InlineAssembly) NodeID() int {
	return ia.ID
}

func GetInlineAssembly(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*InlineAssembly, error) {
	ia := new(InlineAssembly)
	if err := json.Unmarshal([]byte(raw.ToString()), ia); err != nil {
		logger.Errorf("Failed to unmarshal InlineAssembly: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal InlineAssembly: [%v]", err)
	}

	gn.AddASTNode(ia)

	return ia, nil
}

func (ia *InlineAssembly) TraverseFunctionCall() {

}
