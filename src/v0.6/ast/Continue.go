package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type Continue struct {
	ID int `json:"id"`
	NodeType string `json:"nodeType"`
	Src string `json:"src"`
}

func (b *Continue) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "continue"

	if isSc {
		code = code + ";"
	}

	return code
}

func (b *Continue) Type() string {
	return b.NodeType
}

func (b *Continue) Nodes() []ASTNode {
	return nil
}

func (b *Continue) NodeID() int {
	return b.ID
}

func GetContinue(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*Continue, error) {
	b := new(Continue)
	if err := json.Unmarshal([]byte(raw.ToString()), b); err != nil {
		logger.Errorf("Failed to unmarshal Continue: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal Continue: [%v]", err)
	}

	return b, nil
}