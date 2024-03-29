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
			case *ElementaryTypeNameExpression:
				code = code + c.SourceCode(false, false, indent, logger)
			case *Identifier:
				code = code + c.SourceCode(false, false, indent, logger)
			case *Literal:
				code = code + c.SourceCode(false, false, indent, logger)
			case *Conditional:
				code = code + c.SourceCode(false, false, indent, logger)
			case *IndexAccess:
				code = code + c.SourceCode(false, false, indent, logger)
			case *FunctionCall:
				code = code + c.SourceCode(false, false, indent, logger)
			case *MemberAccess:
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

func (te *TupleExpression) NodeID() int {
	return te.ID
}

func GetTupleExpression(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*TupleExpression, error) {
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
				var teComponent ASTNode
				var err error

				if component.Size() > 0 {
					componentNodeType := component.Get("nodeType").ToString()

					switch componentNodeType {
					case "BinaryOperation":
						teComponent, err = GetBinaryOperation(gn, component, logger)
					case "ElementaryTypeNameExpression":
						teComponent, err = GetElementaryTypeNameExpression(gn, component, logger)
					case "Identifier":
						teComponent, err = GetIdentifier(gn, component, logger)
					case "Literal":
						teComponent, err = GetLiteral(gn, component, logger)
					case "Conditional":
						teComponent, err = GetConditional(gn, component, logger)
					case "IndexAccess":
						teComponent, err = GetIndexAccess(gn, component, logger)
					case "FunctionCall":
						teComponent, err = GetFunctionCall(gn, component, logger)
					case "MemberAccess":
						teComponent, err = GetMemberAccess(gn, component, logger)
					default:
						logger.Warnf("Unknown component nodeType [%s] for TupleExpression [src:%s].", componentNodeType, te.Src)
					}
				} else {
					teComponent = &Literal{
						ID:       0,
						Kind:     "number",
						NodeType: "Literal",
						Src:      "xxx",
						Value:    "",
					}
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

	gn.AddASTNode(te)

	return te, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (te *TupleExpression) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	if len(te.components) > 0 {
		for _, component := range te.components {
			switch c := component.(type) {
			case *BinaryOperation:
				c.TraverseFunctionCall(ncp, gn, opt, logger)
			}
		}
	}
}
