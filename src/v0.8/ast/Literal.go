package ast

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type Literal struct {
	HexValue         string `json:"hexValue"`
	ID               int    `json:"id"`
	IsConstant       bool   `json:"isConstant"`
	IsValue          bool   `json:"isValue"`
	IsPure           bool   `json:"isPure"`
	Kind             string `json:"kind"`
	LValueRequested  bool   `json:"lValueRequested"`
	NodeType         string `json:"nodeType"`
	Src              string `json:"src"`
	Subdenomination  string `json:"subdenomination"`
	TypeDescriptions struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
	Value string `json:"value"`
}

func (l *Literal) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	switch l.Kind {
	case "number", "bool":
		code = l.Value
	case "unicodeString":
		code = code + "unicode" + fmt.Sprintf("\"%s\"", l.Value)
	case "string":
		if strings.Contains(l.TypeDescriptions.TypeString, "literal_string hex") {
			var c string
			for i := 0; i <= len(l.HexValue)-2; i += 2 {
				c = c + `\x` + l.HexValue[i:i+2]
			}
			code = code + fmt.Sprintf("\"%s\"", c)
		} else {
			code = code + fmt.Sprintf("\"%s\"", l.Value)
		}
	case "hexString":
		code = code + "hex" + fmt.Sprintf("\"%s\"", l.Value)
	default:
		logger.Warnf("Unknown kind [%s] for Literal [src:%s].", l.Kind, l.Src)
	}

	if l.Subdenomination != "" {
		code = code + " " + l.Subdenomination
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (l *Literal) Type() string {
	return l.NodeType
}

func (l *Literal) Nodes() []ASTNode {
	return nil
}

func (l *Literal) NodeID() int {
	return l.ID
}

func GetLiteral(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*Literal, error) {
	l := new(Literal)
	if err := json.Unmarshal([]byte(raw.ToString()), l); err != nil {
		logger.Errorf("Failed to unmarshal for Literal: [%s].", err)
		return nil, fmt.Errorf("failed to unmarshal for Literal: [%s]", err)
	}

	gn.AddASTNode(l)

	return l, nil
}
