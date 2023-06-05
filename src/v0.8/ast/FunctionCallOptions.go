package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type FunctionCallOptions struct {
	ArgumentTypes []struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"argumentTypes"`
	expression       ASTNode
	ID               int      `json:"id"`
	IsConstant       bool     `json:"isConstant"`
	IsLValue         bool     `json:"isLValue"`
	IsPure           bool     `json:"isPure"`
	LValueRequested  bool     `json:"lValueRequested"`
	Names            []string `json:"names"`
	NodeType         string   `json:"nodeType"`
	options          []ASTNode
	Src              string `json:"src"`
	TypeDescriptions struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
}

func (fco *FunctionCallOptions) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	if fco.expression != nil {
		switch expression := fco.expression.(type) {
		case *MemberAccess:
			code = code + expression.SourceCode(false, false, indent, logger)
		default:
			if expression != nil {
				logger.Warnf("Unknown expression nodeType [%s] for FunctionCallOptions [src:%s].", expression.Type(), fco.Src)
			} else {
				logger.Warnf("Unknown expression nodeType for FunctionCallOptions [src:%s].", fco.Src)
			}
		}
	}

	code = code + "{"

	if len(fco.Names) != len(fco.options) {
		logger.Warnf("The number of names and options mismatch [%d:%d] in FunctionCallOptions [src:%s]", len(fco.Names), len(fco.options), fco.Src)
	}

	for index, option := range fco.options {
		name := fco.Names[index]
		switch opt := option.(type) {
		case *MemberAccess:
			code = code + name + ": " + opt.SourceCode(false, false, indent, logger)
		default:
			if opt != nil {
				logger.Warnf("Unknown option nodeType [%s] for FunctionCallOptions [src:%s].", opt.Type(), fco.Src)
			} else {
				logger.Warnf("Unknown option nodeType for FunctionCallOptions [src:%s].", fco.Src)
			}
		}
		if index < len(fco.options)-1 {
			code = code + ", "
		}
	}

	code = code + "}"

	return code
}

func (fco *FunctionCallOptions) Type() string {
	return fco.NodeType
}

func (fco *FunctionCallOptions) Nodes() []ASTNode {
	return fco.options
}

func (fco *FunctionCallOptions) NodeID() int {
	return fco.ID
}

func GetFunctionCallOptions(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*FunctionCallOptions, error) {
	fco := new(FunctionCallOptions)
	if err := json.Unmarshal([]byte(raw.ToString()), fco); err != nil {
		logger.Errorf("Failed to unmarshal FunctionCallOptions: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal FunctionCallOptions: [%v]", err)
	}

	// expression
	{
		expression := raw.Get("expression")
		if expression.Size() > 0 {
			expressionNodeType := expression.Get("nodeType").ToString()
			var fcoExpression ASTNode
			var err error

			switch expressionNodeType {
			case "MemberAccess":
				fcoExpression, err = GetMemberAccess(gn, expression, logger)
			default:
				logger.Warnf("Unknown expression nodeType [%s] for FunctionCallOptions [src:%s].", expressionNodeType, fco.Src)
			}

			if err != nil {
				return nil, err
			}

			if fcoExpression != nil {
				fco.expression = fcoExpression
			}
		} else {
			logger.Warnf("Expression in FunctionCallOptions [src:%s] should not be nil.", fco.Src)
		}
	}

	// options
	{
		options := raw.Get("options")
		if options.Size() > 0 {
			fco.options = make([]ASTNode, 0)

			for i := 0; i < options.Size(); i++ {
				option := options.Get(i)
				if option.Size() > 0 {
					optionNodeType := option.Get("nodeType").ToString()
					var fcoOption ASTNode
					var err error

					switch optionNodeType {
					case "MemberAccess":
						fcoOption, err = GetMemberAccess(gn, option, logger)
					default:
						logger.Warnf("Unknown option nodeType [%s] for FunctionCallOptions [src:%s].", optionNodeType, fco.Src)
					}

					if err != nil {
						return nil, err
					}

					if fcoOption != nil {
						fco.options = append(fco.options, fcoOption)
					}
				}
			}
		}
	}

	gn.AddASTNode(fco)

	return fco, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (fco *FunctionCallOptions) TraverseFunctionCall() {
	
}
