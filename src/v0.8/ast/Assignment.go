package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type Assignment struct {
	ID               int  `json:"id"`
	IsConstant       bool `json:"isConstant"`
	IsLValue         bool `json:"isLValue"`
	IsPure           bool `json:"isPure"`
	LValueRequested  bool `json:"lValueRequested"`
	leftHandSide     ASTNode
	NodeType         string `json:"nodeType"`
	Operator         string `json:"operator"`
	rightHandSide    ASTNode
	Src              string `json:"src"`
	TypeDescriptions struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
}

func (a *Assignment) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	// leftHandSide
	{
		if a.leftHandSide != nil {
			switch leftHandSide := a.leftHandSide.(type) {
			case *Identifier:
				code = code + leftHandSide.SourceCode(false, false, indent, logger)
			default:
				if leftHandSide != nil {
					logger.Warnf("Unknown leftHandSide nodeType [%s] for Assignment [src:%s].", leftHandSide.Type(), a.Src)
				} else {
					logger.Warnf("Unknown leftHandSide nodeType for Assignment [src:%s].", a.Src)
				}
			}
		}
	}

	code = code + " " + "="

	//rightHandSide
	{
		if a.rightHandSide != nil {
			switch rightHandSide := a.rightHandSide.(type) {
			case *Literal:
				code = code + " " + rightHandSide.SourceCode(false, false, indent, logger)
			default:
				if rightHandSide != nil {
					logger.Warnf("Unknown rightHandSide nodeType [%s] for Assignment [src:%s].", rightHandSide.Type(), a.Src)
				} else {
					logger.Warnf("Unknown rightHandSide nodeType for Assignment [src:%s].", a.Src)
				}
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (a *Assignment) Type() string {
	return a.NodeType
}

func (a *Assignment) Nodes() []ASTNode {
	return nil
}

func GetAssignment(raw jsoniter.Any, logger logging.Logger) (*Assignment, error) {
	a := new(Assignment)
	if err := json.Unmarshal([]byte(raw.ToString()), a); err != nil {
		logger.Errorf("Failed to unmarshal Assignment: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal Assignment: [%v]", err)
	}

	// leftHandSide
	{
		leftHandSide := raw.Get("leftHandSide")
		if leftHandSide.Size() > 0 {
			leftHandSideNodeType := leftHandSide.Get("nodeType").ToString()
			var aLeftHandSide ASTNode
			var err error

			switch leftHandSideNodeType {
			case "Identifier":
				aLeftHandSide, err = GetIdentifier(leftHandSide, logger)
			default:
				logger.Warnf("Unknown leftHandSide nodeType [%s] for Assignment [src:%s].", leftHandSideNodeType, a.Src)
			}

			if err != nil {
				return nil, err
			}

			if aLeftHandSide != nil {
				a.leftHandSide = aLeftHandSide
			}
		}
	}

	// rightHandSide
	{
		rightHandSide := raw.Get("rightHandSide")
		if rightHandSide.Size() > 0 {
			rightHandSideNodeType := rightHandSide.Get("nodeType").ToString()
			var aRightHandSide ASTNode
			var err error

			switch rightHandSideNodeType {
			case "Literal":
				aRightHandSide, err = GetLiteral(rightHandSide, logger)
			default:
				logger.Warnf("Unknown rightHandSide nodeType [%s] for Assignment [src:%s].", rightHandSideNodeType, a.Src)
			}

			if err != nil {
				return nil, err
			}

			if aRightHandSide != nil {
				a.rightHandSide = aRightHandSide
			}
		}
	}

	return a, nil
}
