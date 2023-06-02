package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type InheritanceSpecifier struct {
	baseName ASTNode
	ID       int    `json:"id"`
	NodeType string `json:"nodeType"`
	Src      string `json:"src"`
}

func (is *InheritanceSpecifier) SourceCode() string {
	return is.baseName.SourceCode()
}

func (is *InheritanceSpecifier) Type() string {
	return is.NodeType
}

func (is *InheritanceSpecifier) Nodes() []ASTNode {
	return nil
}

func GetInheritanceSpecifier(raw jsoniter.Any, logger logging.Logger) (*InheritanceSpecifier, error) {
	is := new(InheritanceSpecifier)
	if err := json.Unmarshal([]byte(raw.ToString()), is); err != nil {
		logger.Error("Failed to unmarshal InheritanceSpecifier: [%v]", err)
		return nil, fmt.Errorf("failed to unmarshal InheritanceSpecifier: [%v]", err)
	}
	baseName := raw.Get("baseName")
	if baseName.Size() > 0 {
		var err error
		baseNameNodeType := baseName.Get("nodeType").ToString()
		switch baseNameNodeType {
		case "IdentifierPath":
			is.baseName, err = GetIdentifierPath(baseName, logger)
			if err != nil {
				return nil, err
			}
		default:
			logger.Errorf("Unknown baseName nodeType: [%s-%s]", baseNameNodeType, baseName.Get("src").ToString())
		}
	}
	return is, nil
}
