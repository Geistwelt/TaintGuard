package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type UncheckedBlock struct {
	ID         int    `json:"id"`
	NodeType   string `json:"nodeType"`
	Src        string `json:"src"`
	statements []ASTNode
}

func (ub *UncheckedBlock) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "unchecked {\n"

	if len(ub.statements) > 0 {
		for _, statement := range ub.statements {
			switch stat := statement.(type) {
			case *ExpressionStatement:
				code = code + stat.SourceCode(true, true, indent+"    ", logger)
			case *VariableDeclarationStatement:
				code = code + stat.SourceCode(true, true, indent+"    ", logger)
			case *InlineAssembly:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			case *IfStatement:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			case *Return:
				code = code + stat.SourceCode(true, true, indent+"    ", logger)
			case *WhileStatement:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			case *EmitStatement:
				code = code + stat.SourceCode(true, true, indent+"    ", logger)
			default:
				if stat != nil {
					logger.Warnf("Unknown statement nodeType [%s] for UncheckedBlock [src:%s].", stat.Type(), ub.Src)
				} else {
					logger.Warnf("Unknown statement nodeType for UncheckedBlock [src:%s].", ub.Src)
				}
			}
			code = code + "\n"
		}
	}

	if isIndent {
		code = code + indent
	}
	code = code + "}"

	return code
}

func (ub *UncheckedBlock) Type() string {
	return ub.NodeType
}

func (ub *UncheckedBlock) Nodes() []ASTNode {
	return ub.statements
}

func (ub *UncheckedBlock) NodeID() int {
	return ub.ID
}

func GetUncheckedBlock(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*UncheckedBlock, error) {
	ub := new(UncheckedBlock)
	if err := json.Unmarshal([]byte(raw.ToString()), ub); err != nil {
		logger.Errorf("Failed to unmarshal UncheckedBlock: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal UncheckedBlock: [%v]", err)
	}

	// statements
	{
		statements := raw.Get("statements")
		if statements.Size() > 0 {
			ub.statements = make([]ASTNode, 0)

			for i := 0; i < statements.Size(); i++ {
				statement := statements.Get(i)
				if statement.Size() > 0 {
					statementNodeType := statement.Get("nodeType").ToString()
					var ubStatement ASTNode
					var err error

					switch statementNodeType {
					case "ExpressionStatement":
						ubStatement, err = GetExpressionStatement(gn, statement, logger)
					case "VariableDeclarationStatement":
						ubStatement, err = GetVariableDeclarationStatement(gn, statement, logger)
					case "InlineAssembly":
						ubStatement, err = GetInlineAssembly(gn, statement, logger)
					case "IfStatement":
						ubStatement, err = GetIfStatement(gn, statement, logger)
					case "Return":
						ubStatement, err = GetReturn(gn, statement, logger)
					case "WhileStatement":
						ubStatement, err = GetWhileStatement(gn, statement, logger)
					case "EmitStatement":
						ubStatement, err = GetEmitStatement(gn, statement, logger)
					default:
						logger.Warnf("Unknown statement nodeType [%s] for UncheckedBlock [src:%s].", statementNodeType, ub.Src)
					}

					if err != nil {
						return nil, err
					}

					if ubStatement != nil {
						ub.statements = append(ub.statements, ubStatement)
					}
				}
			}
		}
	}

	gn.AddASTNode(ub)

	return ub, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (ub *UncheckedBlock) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	if len(ub.statements) > 0 {
		for _, statement := range ub.statements {
			switch stat := statement.(type) {
			case *ExpressionStatement:
				stat.TraverseFunctionCall(ncp, gn, opt, logger)
			case *VariableDeclarationStatement:
				stat.TraverseFunctionCall(ncp, gn, opt, logger)
			}
		}
	}
}

func (ub *UncheckedBlock) TraverseTaintOwner(opt *Option, logger logging.Logger) {
	if len(ub.statements) > 0 {
		for _, statement := range ub.statements {
			switch stat := statement.(type) {
			case *ExpressionStatement:
				stat.TraverseTaintOwner(opt, logger)
			}
		}
	}
}
