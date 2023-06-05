package ast

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type PragmaDirective struct {
	ID       int      `json:"id"`
	Literals []string `json:"literals"`
	NodeType string   `json:"nodeType"`
	Src      string   `json:"src"`
}

func (pd *PragmaDirective) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	reNum := regexp.MustCompile(`0*\.\d*`)
	var code string = "pragma"
	for _, literal := range pd.Literals {
		literal = strings.Trim(literal, "\"")
		if !reNum.Match([]byte(literal)) {
			code = code + " " + literal
		} else {
			code = code + literal
		}
	}
	code = code + ";"
	return code
}

func (pd *PragmaDirective) Type() string {
	return pd.NodeType
}

func (pd *PragmaDirective) Nodes() []ASTNode {
	return nil
}

func (pd *PragmaDirective) NodeID() int {
	return pd.ID
}

func GetPragmaDirective(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*PragmaDirective, error) {
	pd := new(PragmaDirective)
	if err := json.Unmarshal([]byte(raw.ToString()), pd); err != nil {
		logger.Errorf("Failed to unmarshal PragmaDirective: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal PragmaDirective: [%v]", err)
	}

	gn.AddASTNode(pd)
	
	return pd, nil
}
