package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type TupleExpression struct {
	components       []ASTNode
	ID               int    `json:"id"`
	IsConstant       bool   `json:"isConstant"`
	IsInlineArray    bool   `json:"isInlineArray"`
	IsLValue         bool   `json:"isLValue"`
	IsPure           bool   `json:"isPure"`
	LValueRequested  bool   `json:"lValueRequested"`
	NodeType         string `json:"nodeType"`
	Src              string `json:"src"`
	TypeDescriptions struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
}

func (te *TupleExpression) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "("

	if len(te.components) > 0 {
		for index, component := range te.components {
			switch c := component.(type) {
			case *BinaryOperation:
				code = code + c.SourceCode(false, false, indent, logger)
			default:
				if c != nil {
					logger.Warnf("Unknown component nodeType [%s] for TupleExpression [src:%s].", c.Type(), te.Src)
				} else {
					logger.Warnf("Unknown component nodeType for TupleExpression [src:%s].", te.Src)
				}
			}
			if index < len(te.components)-1 {
				code = code + ", "
			}
		}
	} else {
		logger.Warnf("TupleExpression [src:%s] should have more than 0 components.", te.Src)
	}

	code = code + ")"

	return code
}

func (te *TupleExpression) Type() string {
	return te.NodeType
}

func (te *TupleExpression) Nodes() []ASTNode {
	return te.components
}

func GetTupleExpression(raw jsoniter.Any, logger logging.Logger) (*TupleExpression, error) {
	te := new(TupleExpression)
	if err := json.Unmarshal([]byte(raw.ToString()), te); err != nil {
		logger.Errorf("Failed to unmarshal TupleExpression: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal TupleExpression: [%v]", err)
	}

	// components
	{
		components := raw.Get("components")
		if components.Size() > 0 {
			te.components = make([]ASTNode, 0)
			for i := 0; i < components.Size(); i++ {
				component := components.Get(i)
				componentNodeType := component.Get("nodeType").ToString()
				var teComponent ASTNode
				var err error

				switch componentNodeType {
				case "BinaryOperation":
					teComponent, err = GetBinaryOperation(component, logger)
				default:
					logger.Warnf("Unknown component nodeType [%s] for TupleExpression [src:%s].", componentNodeType, te.Src)
				}

				if err != nil {
					return nil, err
				}

				if teComponent != nil {
					te.components = append(te.components, teComponent)
				}
			}
		}
	}

	return te, nil
}