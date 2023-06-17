package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type OverrideSpecifier struct {
	ID        int    `json:"id"`
	NodeType  string `json:"nodeType"`
	overrides []ASTNode
	Src       string `json:"src"`
}

func (os *OverrideSpecifier) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	code = code + "override"

	if len(os.overrides) > 0 {
		code = code + "("
		for index, o := range os.overrides {
			code = code + o.SourceCode(false, false, indent, logger)
			if index < len(os.overrides)-1 {
				code = code + ", "
			}
		}
		code = code + ")"
	}

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

func (os *OverrideSpecifier) NodeID() int {
	return os.ID
}

func GetOverrideSpecifier(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*OverrideSpecifier, error) {
	os := new(OverrideSpecifier)
	if err := json.Unmarshal([]byte(raw.ToString()), os); err != nil {
		logger.Errorf("Failed to unmarshal OverrideSpecifier: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal OverrideSpecifier: [%v]", err)
	}

	// overrides
	{
		overrides := raw.Get("overrides")
		if overrides.Size() > 0 {
			os.overrides = make([]ASTNode, 0)

			for i := 0; i < overrides.Size(); i++ {
				override := overrides.Get(i)
				if override.Size() > 0 {
					overrideNodeType := override.Get("nodeType").ToString()
					var osOverride ASTNode
					var err error

					switch overrideNodeType {
					case "IdentifierPath":
						osOverride, err = GetIdentifierPath(gn, override, logger)
					default:
						logger.Warnf("Unknown override nodeType [%s] for OverrideSpecifier [src:%s].", overrideNodeType, os.Src)
					}

					if err != nil {
						return nil, err
					}

					if osOverride != nil {
						os.overrides = append(os.overrides, osOverride)
					}
				} else {
					logger.Errorf("Override in OverrideSpecifier [src:%s] should not be empty.", os.Src)
				}
			}
		}
	}

	gn.AddASTNode(os)

	return os, nil
}
