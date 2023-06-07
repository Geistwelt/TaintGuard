package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type FunctionCall struct {
	arguments        []ASTNode
	expression       ASTNode
	ID               int      `json:"id"`
	IsConstant       bool     `json:"isConstant"`
	IsLValue         bool     `json:"isLValue"`
	IsPure           bool     `json:"isPure"`
	Kind             string   `json:"kind"`
	LValueRequested  bool     `json:"lValueRequested"`
	NameLocations    []int    `json:"nameLocations"`
	Names            []string `json:"names"`
	NodeType         string   `json:"nodeType"`
	Src              string   `json:"src"`
	TryCall          bool     `json:"tryCall"`
	TypeDescriptions struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`

	referencedFunctionDefinition int
}

func (fc *FunctionCall) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	// expression
	{
		if fc.expression != nil {
			switch expression := fc.expression.(type) {
			case *Identifier:
				if len(expression.ArgumentTypes) != len(fc.arguments) {
					logger.Warnf("Number of arguments mismatch [%d:%d] in FunctionCall: [src:%s].", len(expression.ArgumentTypes), len(fc.arguments), fc.Src)
				}
				code = code + expression.SourceCode(false, false, indent, logger)
			case *MemberAccess:
				if len(expression.ArgumentTypes) != len(fc.arguments) {
					logger.Warnf("Number of arguments mismatch [%d:%d] in FunctionCall: [src:%s].", len(expression.ArgumentTypes), len(fc.arguments), fc.Src)
				}
				code = code + expression.SourceCode(false, false, indent, logger)
			case *ElementaryTypeNameExpression:
				if len(expression.ArgumentTypes) != len(fc.arguments) {
					logger.Warnf("Number of arguments mismatch [%d:%d] in FunctionCall: [src:%s].", len(expression.ArgumentTypes), len(fc.arguments), fc.Src)
				}
				code = code + expression.SourceCode(false, false, indent, logger)
			case *NewExpression:
				if len(expression.ArgumentTypes) != len(fc.arguments) {
					logger.Warnf("Number of arguments mismatch [%d:%d] in FunctionCall: [src:%s].", len(expression.ArgumentTypes), len(fc.arguments), fc.Src)
				}
				code = code + expression.SourceCode(false, false, indent, logger)
			case *FunctionCallOptions:
				if len(expression.ArgumentTypes) != len(fc.arguments) {
					logger.Warnf("Number of arguments mismatch [%d:%d] in FunctionCall: [src:%s].", len(expression.ArgumentTypes), len(fc.arguments), fc.Src)
				}
				code = code + expression.SourceCode(false, false, indent, logger)
			default:
				if expression != nil {
					logger.Warnf("Unknown expression nodeType [%s] for FunctionCall [src:%s].", expression.Type(), fc.Src)
				} else {
					logger.Warnf("Unknown expression nodeType for FunctionCall [src:%s].", fc.Src)
				}
			}
		}
	}

	code = code + "("

	// arguments
	{
		if len(fc.arguments) > 0 {
			for index, argument := range fc.arguments {
				switch arg := argument.(type) {
				case *FunctionCall:
					code = code + arg.SourceCode(false, false, indent, logger)
				case *Identifier:
					code = code + arg.SourceCode(false, false, indent, logger)
				case *Literal:
					code = code + arg.SourceCode(false, false, indent, logger)
				case *BinaryOperation:
					code = code + arg.SourceCode(false, false, indent, logger)
				case *MemberAccess:
					code = code + arg.SourceCode(false, false, indent, logger)
				case *UnaryOperation:
					code = code + arg.SourceCode(false, false, indent, logger)
				case *Conditional:
					code = code + arg.SourceCode(false, false, indent, logger)
				case *ElementaryTypeNameExpression:
					code = code + arg.SourceCode(false, false, indent, logger)
				case *TupleExpression:
					code = code + arg.SourceCode(false, false, indent, logger)
				case *IndexAccess:
					code = code + arg.SourceCode(false, false, indent, logger)
				default:
					if arg != nil {
						logger.Warnf("Unknown argument nodeType [%s] for FunctionCall [src:%s].", arg.Type(), fc.Src)
					} else {
						logger.Warnf("Unknown argument nodeType for FunctionCall [src:%s].", fc.Src)
					}
				}
				if index < len(fc.arguments)-1 {
					code = code + ", "
				}
			}
		}
	}

	code = code + ")"

	if isSc {
		code = code + ";"
	}

	if fc.Kind == "typeConversion" && len(code) >= 15 && code[0:15] == "address payable" {
		code = code[8:]
	}

	return code
}

func (fc *FunctionCall) Type() string {
	return fc.NodeType
}

func (fc *FunctionCall) Nodes() []ASTNode {
	return fc.arguments
}

func (fc *FunctionCall) NodeID() int {
	return fc.ID
}

func GetFunctionCall(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*FunctionCall, error) {
	fc := new(FunctionCall)
	fc.referencedFunctionDefinition = -1
	if err := json.Unmarshal([]byte(raw.ToString()), fc); err != nil {
		logger.Errorf("Failed to unmarshal FunctionCall: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal FunctionCall: [%v]", err)
	}

	// expression
	{
		expression := raw.Get("expression")
		if expression.Size() > 0 {
			expressionNodeType := expression.Get("nodeType").ToString()
			var fcExpression ASTNode
			var err error

			switch expressionNodeType {
			case "Identifier":
				fcExpression, err = GetIdentifier(gn, expression, logger)
				identifier, _ := fcExpression.(*Identifier)
				fc.referencedFunctionDefinition = identifier.ReferencedDeclaration
			case "MemberAccess":
				fcExpression, err = GetMemberAccess(gn, expression, logger)
				memberAccess, _ := fcExpression.(*MemberAccess)
				fc.referencedFunctionDefinition = memberAccess.ReferencedDeclaration
			case "ElementaryTypeNameExpression":
				fcExpression, err = GetElementaryTypeNameExpression(gn, expression, logger)
			case "NewExpression":
				fcExpression, err = GetNewExpression(gn, expression, logger)
			case "FunctionCallOptions":
				fcExpression, err = GetFunctionCallOptions(gn, expression, logger)
			default:
				logger.Warnf("Unknown expression nodeType [%s] for FunctionCall [src:%s].", expressionNodeType, fc.Src)
			}

			if err != nil {
				return nil, err
			}

			if fcExpression != nil {
				fc.expression = fcExpression
			}
		}
	}

	// arguments
	{
		arguments := raw.Get("arguments")
		if arguments.Size() > 0 {
			fc.arguments = make([]ASTNode, 0)

			for i := 0; i < arguments.Size(); i++ {
				argument := arguments.Get(i)
				argumentNodeType := argument.Get("nodeType").ToString()
				var fcArgument ASTNode
				var err error

				switch argumentNodeType {
				case "FunctionCall":
					fcArgument, err = GetFunctionCall(gn, argument, logger)
				case "Identifier":
					fcArgument, err = GetIdentifier(gn, argument, logger)
				case "Literal":
					fcArgument, err = GetLiteral(gn, argument, logger)
				case "BinaryOperation":
					fcArgument, err = GetBinaryOperation(gn, argument, logger)
				case "MemberAccess":
					fcArgument, err = GetMemberAccess(gn, argument, logger)
				case "UnaryOperation":
					fcArgument, err = GetUnaryOperation(gn, argument, logger)
				case "Conditional":
					fcArgument, err = GetConditional(gn, argument, logger)
				case "ElementaryTypeNameExpression":
					fcArgument, err = GetElementaryTypeNameExpression(gn, argument, logger)
				case "TupleExpression":
					fcArgument, err = GetTupleExpression(gn, argument, logger)
				case "IndexAccess":
					fcArgument, err = GetIndexAccess(gn, argument, logger)
				default:
					logger.Warnf("Unknown argument nodeType [%s] for FunctionCall [src:%s].", argumentNodeType, fc.Src)
				}

				if err != nil {
					return nil, err
				}

				if fcArgument != nil {
					fc.arguments = append(fc.arguments, fcArgument)
				}
			}
		}
	}

	gn.AddASTNode(fc)

	return fc, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (fc *FunctionCall) ReferencedFunctionDefinition() int {
	return fc.referencedFunctionDefinition
}

func (fc *FunctionCall) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes) {
	if fc.ReferencedFunctionDefinition() != -1 {
		fd := gn.Functions()[fc.ReferencedFunctionDefinition()]
		function, ok := fd.(*FunctionDefinition)
		if ok {
			ncp.SetCallee(function.Signature(), fc.ReferencedFunctionDefinition())
		}
	}
}