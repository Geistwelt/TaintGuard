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
			case *IndexAccess:
				code = code + leftHandSide.SourceCode(false, false, indent, logger)
			case *TupleExpression:
				code = code + leftHandSide.SourceCode(false, false, indent, logger)
			case *MemberAccess:
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

	code = code + " " + a.Operator

	//rightHandSide
	{
		if a.rightHandSide != nil {
			switch rightHandSide := a.rightHandSide.(type) {
			case *Literal:
				code = code + " " + rightHandSide.SourceCode(false, false, indent, logger)
			case *Identifier:
				code = code + " " + rightHandSide.SourceCode(false, false, indent, logger)
			case *FunctionCall:
				code = code + " " + rightHandSide.SourceCode(false, false, indent, logger)
			case *MemberAccess:
				code = code + " " + rightHandSide.SourceCode(false, false, indent, logger)
			case *BinaryOperation:
				code = code + " " + rightHandSide.SourceCode(false, false, indent, logger)
			case *IndexAccess:
				code = code + " " + rightHandSide.SourceCode(false, false, indent, logger)
			case *UnaryOperation:
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

func (a *Assignment) NodeID() int {
	return a.ID
}

func GetAssignment(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*Assignment, error) {
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
				aLeftHandSide, err = GetIdentifier(gn, leftHandSide, logger)
			case "IndexAccess":
				aLeftHandSide, err = GetIndexAccess(gn, leftHandSide, logger)
			case "TupleExpression":
				aLeftHandSide, err = GetTupleExpression(gn, leftHandSide, logger)
			case "MemberAccess":
				aLeftHandSide, err = GetMemberAccess(gn, leftHandSide, logger)
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
				aRightHandSide, err = GetLiteral(gn, rightHandSide, logger)
			case "Identifier":
				aRightHandSide, err = GetIdentifier(gn, rightHandSide, logger)
			case "FunctionCall":
				aRightHandSide, err = GetFunctionCall(gn, rightHandSide, logger)
			case "MemberAccess":
				aRightHandSide, err = GetMemberAccess(gn, rightHandSide, logger)
			case "BinaryOperation":
				aRightHandSide, err = GetBinaryOperation(gn, rightHandSide, logger)
			case "IndexAccess":
				aRightHandSide, err = GetIndexAccess(gn, rightHandSide, logger)
			case "UnaryOperation":
				aRightHandSide, err = GetUnaryOperation(gn, rightHandSide, logger)
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

	gn.AddASTNode(a)

	return a, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (a *Assignment) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	// leftHandSide
	{
		if a.leftHandSide != nil {
			switch leftHandSide := a.leftHandSide.(type) {
			case *IndexAccess:
				leftHandSide.TraverseFunctionCall(ncp, gn, opt, logger)
			case *MemberAccess:
				leftHandSide.TraverseFunctionCall(ncp, gn, opt, logger)
			}
		}
	}

	//rightHandSide
	{
		if a.rightHandSide != nil {
			switch rightHandSide := a.rightHandSide.(type) {
			case *FunctionCall:
				rightHandSide.TraverseFunctionCall(ncp, gn, opt, logger)
			case *MemberAccess:
				rightHandSide.TraverseFunctionCall(ncp, gn, opt, logger)
			case *BinaryOperation:
				rightHandSide.TraverseFunctionCall(ncp, gn, opt, logger)
			case *IndexAccess:
				rightHandSide.TraverseFunctionCall(ncp, gn, opt, logger)
			}
		}
	}
}

func (a *Assignment) TraverseTaintOwner(opt *Option, logger logging.Logger) {
	// leftHandSide
	{
		if a.leftHandSide != nil {
			switch leftHandSide := a.leftHandSide.(type) {
			case *Identifier:
				if leftHandSide.Name == opt.SimilarOwnerVariableName {
					opt.IsTainted = true
					trackAssignment := &Assignment{
						leftHandSide: &IndexAccess{
							baseExpression: &Identifier{
								Name:     opt.TrackOwnerMappingName,
								NodeType: "Identifier",
								Src:      "xxx",
							},
							ID: 0,
							indexExpression: &Literal{
								Kind:     "string",
								NodeType: "Literal",
								Src:      "xxx",
								Value:    opt.TrackFunctionDefinitionName,
							},
							NodeType: "IndexAccess",
							Src:      "xxx",
						},
						NodeType:      "Assignment",
						Operator:      "=",
						rightHandSide: a.rightHandSide,
						Src:           "xxx",
					}
					opt.TrackAssignment = trackAssignment
				}
			}
		}
	}

	// //rightHandSide
	// {
	// 	if a.rightHandSide != nil {
	// 		switch rightHandSide := a.rightHandSide.(type) {
	// 		case *FunctionCall:
	// 			rightHandSide.TraverseFunctionCall(ncp, gn, opt, logger)
	// 		case *MemberAccess:
	// 			rightHandSide.TraverseFunctionCall(ncp, gn, opt, logger)
	// 		case *BinaryOperation:
	// 			rightHandSide.TraverseFunctionCall(ncp, gn, opt, logger)
	// 		case *IndexAccess:
	// 			rightHandSide.TraverseFunctionCall(ncp, gn, opt, logger)
	// 		}
	// 	}
	// }
}

func (a *Assignment) TraverseDelegatecall(opt *Option, logger logging.Logger) {
	// leftHandSide
	{
		if a.leftHandSide != nil {
			switch leftHandSide := a.leftHandSide.(type) {
			case *MemberAccess:
				leftHandSide.TraverseDelegatecall(opt, logger)
			}
		}
	}

	//rightHandSide
	{
		if a.rightHandSide != nil {
			switch rightHandSide := a.rightHandSide.(type) {
			case *FunctionCall:
				rightHandSide.TraverseDelegatecall(opt, logger)
			case *MemberAccess:
				rightHandSide.TraverseDelegatecall(opt, logger)
			case *BinaryOperation:
				rightHandSide.TraverseDelegatecall(opt, logger)
			}
		}
	}
}

func (a *Assignment) SetLeft(left ASTNode) {
	a.leftHandSide = left
}

func (a *Assignment) SetRight(right ASTNode) {
	a.rightHandSide = right
}
