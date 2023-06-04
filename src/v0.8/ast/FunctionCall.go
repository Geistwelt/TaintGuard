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

	return code
}

func (fc *FunctionCall) Type() string {
	return fc.NodeType
}

func (fc *FunctionCall) Nodes() []ASTNode {
	return fc.arguments
}

func GetFunctionCall(raw jsoniter.Any, logger logging.Logger) (*FunctionCall, error) {
	fc := new(FunctionCall)
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
				fcExpression, err = GetIdentifier(expression, logger)
			case "MemberAccess":
				fcExpression, err = GetMemberAccess(expression, logger)
			case "ElementaryTypeNameExpression":
				fcExpression, err = GetElementaryTypeNameExpression(expression, logger)
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
					fcArgument, err = GetFunctionCall(argument, logger)
				case "Identifier":
					fcArgument, err = GetIdentifier(argument, logger)
				case "Literal":
					fcArgument, err = GetLiteral(argument, logger)
				case "BinaryOperation":
					fcArgument, err = GetBinaryOperation(argument, logger)
				case "MemberAccess":
					fcArgument, err = GetMemberAccess(argument, logger)
				case "UnaryOperation":
					fcArgument, err = GetUnaryOperation(argument, logger)
				case "Conditional":
					fcArgument, err = GetConditional(argument, logger)
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

	return fc, nil
}