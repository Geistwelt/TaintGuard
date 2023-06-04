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

func GetIndexAccess(raw jsoniter.Any, logger logging.Logger) (*IndexAccess, error) {
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
				iaBaseExpression, err = GetIdentifier(baseExpression, logger)
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
				iaIndexExpression, err = GetIdentifier(indexExpression, logger)
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

	return ia, nil
}
