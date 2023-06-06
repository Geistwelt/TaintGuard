package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type MemberAccess struct {
	ArgumentTypes []struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"argumentTypes"`
	expression            ASTNode
	ID                    int    `json:"id"`
	IsConstant            bool   `json:"isConstant"`
	IsLValue              bool   `json:"isLValue"`
	IsPure                bool   `json:"isPure"`
	LValueRequested       bool   `json:"lValueRequested"`
	MemberLocation        string `json:"memberLocation"`
	MemberName            string `json:"memberName"`
	NodeType              string `json:"nodeType"`
	ReferencedDeclaration int    `json:"referencedDeclaration"`
	Src                   string `json:"src"`
	TypeDescriptions      struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
}

func (ma *MemberAccess) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	// expression
	{
		if ma.expression != nil {
			switch expression := ma.expression.(type) {
			case *IndexAccess:
				code = code + expression.SourceCode(false, false, indent, logger)
			case *Identifier:
				code = code + expression.SourceCode(false, false, indent, logger)
			case *FunctionCall:
				code = code + expression.SourceCode(false, false, indent, logger)
			case *MemberAccess:
				code = code + expression.SourceCode(false, false, indent, logger)
			default:
				if expression != nil {
					logger.Warnf("Unknown expression nodeType [%s] for MemberAccess [src:%s].", expression.Type(), ma.Src)
				} else {
					logger.Warnf("Unknown expression nodeType for MemberAccess [src:%s].", ma.Src)
				}
			}
		}
	}

	code = code + "."

	if ma.MemberName != "" {
		code = code + ma.MemberName
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (ma *MemberAccess) Type() string {
	return ma.NodeType
}

func (ma *MemberAccess) Nodes() []ASTNode {
	return nil
}

func (ma *MemberAccess) NodeID() int {
	return ma.ID
}

func GetMemberAccess(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*MemberAccess, error) {
	ma := new(MemberAccess)
	if err := json.Unmarshal([]byte(raw.ToString()), ma); err != nil {
		logger.Errorf("Failed to unmarshal MemberAccess: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal MemberAccess: [%v]", err)
	}

	// expression
	{
		expression := raw.Get("expression")
		if expression.Size() > 0 {
			expressionNodeType := expression.Get("nodeType").ToString()
			var maExpression ASTNode
			var err error

			switch expressionNodeType {
			case "IndexAccess":
				maExpression, err = GetIndexAccess(gn, expression, logger)
			case "Identifier":
				maExpression, err = GetIdentifier(gn, expression, logger)
			case "FunctionCall":
				maExpression, err = GetFunctionCall(gn, expression, logger)
			case "MemberAccess":
				maExpression, err = GetMemberAccess(gn, expression, logger)
			default:
				logger.Warnf("Unknown expression nodeType [%s] for MemberAccess [src:%s].", expressionNodeType, ma.Src)
			}

			if err != nil {
				return nil, err
			}

			if maExpression != nil {
				ma.expression = maExpression
			}
		}
	}

	gn.AddASTNode(ma)

	return ma, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (ma *MemberAccess) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes) {
	// expression
	{
		if ma.expression != nil {
			switch expression := ma.expression.(type) {
			case *IndexAccess:
				expression.TraverseFunctionCall(ncp, gn)
			case *FunctionCall:
				expression.TraverseFunctionCall(ncp, gn)
			case *MemberAccess:
				expression.TraverseFunctionCall(ncp, gn)
			}
		}
	}
}