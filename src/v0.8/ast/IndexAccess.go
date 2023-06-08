package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type IndexAccess struct {
	baseExpression   ASTNode
	ID               int `json:"id"`
	indexExpression  ASTNode
	IsConstant       bool   `json:"isConstant"`
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

func (ia *IndexAccess) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	// baseExpression
	{
		if ia.baseExpression != nil {
			switch baseExpression := ia.baseExpression.(type) {
			case *Identifier:
				code = code + baseExpression.SourceCode(false, false, indent, logger)
			case *IndexAccess:
				code = code + baseExpression.SourceCode(false, false, indent, logger)
			case *MemberAccess:
				code = code + baseExpression.SourceCode(false, false, indent, logger)
			case *ElementaryTypeNameExpression:
				code = code + baseExpression.SourceCode(false, false, indent, logger)
			default:
				if baseExpression != nil {
					logger.Warnf("Unknown baseExpression nodeType [%s] for IndexAccess [src:%s].", baseExpression.Type(), ia.Src)
				} else {
					logger.Warnf("Unknown baseExpression nodeType for IndexAccess [src:%s].", ia.Src)
				}
			}
		}
	}

	code = code + "["

	// indexExpression
	{
		if ia.indexExpression != nil {
			switch indexExpression := ia.indexExpression.(type) {
			case *Identifier:
				code = code + indexExpression.SourceCode(false, false, indent, logger)
			case *FunctionCall:
				code = code + indexExpression.SourceCode(false, false, indent, logger)
			case *Literal:
				code = code + indexExpression.SourceCode(false, false, indent, logger)
			case *IndexAccess:
				code = code + indexExpression.SourceCode(false, false, indent, logger)
			case *MemberAccess:
				code = code + indexExpression.SourceCode(false, false, indent, logger)
			case *BinaryOperation:
				code = code + indexExpression.SourceCode(false, false, indent, logger)
			case *UnaryOperation:
				code = code + indexExpression.SourceCode(false, false, indent, logger)
			default:
				if indexExpression != nil {
					logger.Warnf("Unknown indexExpression nodeType [%s] for IndexAccess [src:%s].", indexExpression.Type(), ia.Src)
				} else {
					logger.Warnf("Unknown indexExpression nodeType for IndexAccess [src:%s].", ia.Src)
				}
			}
		}
	}

	code = code + "]"

	if isSc {
		code = code + ";"
	}

	return code
}

func (ia *IndexAccess) Type() string {
	return ia.NodeType
}

func (ia *IndexAccess) Nodes() []ASTNode {
	return nil
}

func (ia *IndexAccess) NodeID() int {
	return ia.ID
}

func GetIndexAccess(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*IndexAccess, error) {
	ia := new(IndexAccess)
	if err := json.Unmarshal([]byte(raw.ToString()), ia); err != nil {
		logger.Errorf("Failed to unmarshal IndexAccess: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal IndexAccess: [%v]", err)
	}

	// baseExpression
	{
		baseExpression := raw.Get("baseExpression")
		if baseExpression.Size() > 0 {
			baseExpressionNodeType := baseExpression.Get("nodeType").ToString()
			var iaBaseExpression ASTNode
			var err error

			switch baseExpressionNodeType {
			case "Identifier":
				iaBaseExpression, err = GetIdentifier(gn, baseExpression, logger)
			case "IndexAccess":
				iaBaseExpression, err = GetIndexAccess(gn, baseExpression, logger)
			case "MemberAccess":
				iaBaseExpression, err = GetMemberAccess(gn, baseExpression, logger)
			case "ElementaryTypeNameExpression":
				iaBaseExpression, err = GetElementaryTypeNameExpression(gn, baseExpression, logger)
			default:
				logger.Warnf("Unknown baseExpression [%s] for IndexAccess [src:%s].", baseExpressionNodeType, ia.Src)
			}

			if err != nil {
				return nil, err
			}

			if iaBaseExpression != nil {
				ia.baseExpression = iaBaseExpression
			}
		}
	}

	// indexExpression
	{
		indexExpression := raw.Get("indexExpression")
		if indexExpression.Size() > 0 {
			indexExpressionNodeType := indexExpression.Get("nodeType").ToString()
			var iaIndexExpression ASTNode
			var err error

			switch indexExpressionNodeType {
			case "Identifier":
				iaIndexExpression, err = GetIdentifier(gn, indexExpression, logger)
			case "FunctionCall":
				iaIndexExpression, err = GetFunctionCall(gn, indexExpression, logger)
			case "Literal":
				iaIndexExpression, err = GetLiteral(gn, indexExpression, logger)
			case "IndexAccess":
				iaIndexExpression, err = GetIndexAccess(gn, indexExpression, logger)
			case "MemberAccess":
				iaIndexExpression, err = GetMemberAccess(gn, indexExpression, logger)
			case "BinaryOperation":
				iaIndexExpression, err = GetBinaryOperation(gn, indexExpression, logger)
			case "UnaryOperation":
				iaIndexExpression, err = GetUnaryOperation(gn, indexExpression, logger)
			default:
				logger.Warnf("Unknown indexExpression [%s] for IndexAccess [src:%s].", indexExpressionNodeType, ia.Src)
			}

			if err != nil {
				return nil, err
			}

			if iaIndexExpression != nil {
				ia.indexExpression = iaIndexExpression
			}
		}
	}

	gn.AddASTNode(ia)

	return ia, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (ia *IndexAccess) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	// baseExpression
	{
		if ia.baseExpression != nil {
			switch baseExpression := ia.baseExpression.(type) {
			case *IndexAccess:
				baseExpression.TraverseFunctionCall(ncp, gn, opt, logger)
			}
		}
	}

	// indexExpression
	{
		if ia.indexExpression != nil {
			switch indexExpression := ia.indexExpression.(type) {
			case *FunctionCall:
				indexExpression.TraverseFunctionCall(ncp, gn, opt, logger)
			case *IndexAccess:
				indexExpression.TraverseFunctionCall(ncp, gn, opt, logger)
			case *MemberAccess:
				indexExpression.TraverseFunctionCall(ncp, gn, opt, logger)
			}
		}
	}
}
