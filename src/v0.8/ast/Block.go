package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type Block struct {
	ID         int    `json:"id"`
	NodeType   string `json:"nodeType"`
	Src        string `json:"src"`
	statements []ASTNode
}

func (b *Block) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if len(b.statements) > 0 {
		for index, statement := range b.statements {
			switch stat := statement.(type) {
			case *ExpressionStatement:
				code = code + stat.SourceCode(true, true, indent+"    ", logger)
			case *PlaceholderStatement:
				code = code + stat.SourceCode(true, true, indent+"    ", logger)
			case *Return:
				code = code + stat.SourceCode(true, true, indent+"    ", logger)
			case *EmitStatement:
				code = code + stat.SourceCode(true, true, indent+"    ", logger)
			case *IfStatement:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			case *VariableDeclarationStatement:
				code = code + stat.SourceCode(true, true, indent+"    ", logger)
			case *InlineAssembly:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			case *ForStatement:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			case *RevertStatement:
				code = code + stat.SourceCode(true, true, indent+"    ", logger)
			case *Block:
				code = code + indent + "    " + "{" + "\n" + stat.SourceCode(false, true, indent+"    ", logger) + "\n" + indent + "    " + "}"
			case *UncheckedBlock:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			case *WhileStatement:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			case *TryStatement:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			case *DoWhileStatement:
				code = code + stat.SourceCode(true, true, indent+"    ", logger)
			default:
				if stat != nil {
					logger.Warnf("Unknown statement nodeType [%s] for Block [src:%s].", stat.Type(), b.Src)
				} else {
					logger.Warnf("Unknown statement nodeType for Block [src:%s].", b.Src)
				}
			}
			if index < len(b.statements)-1 {
				code = code + "\n"
			}
		}
	}

	return code
}

func (b *Block) Type() string {
	return b.NodeType
}

func (b *Block) Nodes() []ASTNode {
	return b.statements
}

func (b *Block) NodeID() int {
	return b.ID
}

func GetBlock(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*Block, error) {
	b := new(Block)
	if err := json.Unmarshal([]byte(raw.ToString()), b); err != nil {
		logger.Errorf("Failed to unmarshal Block: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal Block: [%v]", err)
	}

	// statements
	{
		statements := raw.Get("statements")
		if statements.Size() > 0 {
			b.statements = make([]ASTNode, 0)

			for i := 0; i < statements.Size(); i++ {
				statement := statements.Get(i)
				statementNodeType := statement.Get("nodeType").ToString()
				var bStatement ASTNode
				var err error

				switch statementNodeType {
				case "PlaceholderStatement":
					bStatement, err = GetPlaceholderStatement(gn, statement, logger)
				case "ExpressionStatement":
					bStatement, err = GetExpressionStatement(gn, statement, logger)
				case "Return":
					bStatement, err = GetReturn(gn, statement, logger)
				case "EmitStatement":
					bStatement, err = GetEmitStatement(gn, statement, logger)
				case "IfStatement":
					bStatement, err = GetIfStatement(gn, statement, logger)
				case "VariableDeclarationStatement":
					bStatement, err = GetVariableDeclarationStatement(gn, statement, logger)
				case "InlineAssembly":
					bStatement, err = GetInlineAssembly(gn, statement, logger)
				case "ForStatement":
					bStatement, err = GetForStatement(gn, statement, logger)
				case "RevertStatement":
					bStatement, err = GetRevertStatement(gn, statement, logger)
				case "Block":
					bStatement, err = GetBlock(gn, statement, logger)
				case "UncheckedBlock":
					bStatement, err = GetUncheckedBlock(gn, statement, logger)
				case "WhileStatement":
					bStatement, err = GetWhileStatement(gn, statement, logger)
				case "TryStatement":
					bStatement, err = GetTryStatement(gn, statement, logger)
				case "DoWhileStatement":
					bStatement, err = GetDoWhileStatement(gn, statement, logger)
				default:
					logger.Warnf("Unknown statement nodeType [%s] for Block [src:%s].", statementNodeType, b.Src)
				}

				if err != nil {
					return nil, err
				}
				if bStatement != nil {
					b.statements = append(b.statements, bStatement)
				}
			}
		}
	}

	gn.AddASTNode(b)

	return b, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (b *Block) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	if len(b.statements) > 0 {
		for _, statement := range b.statements {
			switch stat := statement.(type) {
			case *ExpressionStatement:
				stat.TraverseFunctionCall(ncp, gn, opt, logger)
			case *Return:
				stat.TraverseFunctionCall(ncp, gn, opt, logger)
			case *EmitStatement:
				stat.TraverseFunctionCall(ncp, gn, opt, logger)
			case *IfStatement:
				stat.TraverseFunctionCall(ncp, gn, opt, logger)
			case *VariableDeclarationStatement:
				stat.TraverseFunctionCall(ncp, gn, opt, logger)
			case *ForStatement:
				stat.TraverseFunctionCall(ncp, gn, opt, logger)
			case *RevertStatement:
				stat.TraverseFunctionCall(ncp, gn, opt, logger)
			case *Block:
				stat.TraverseFunctionCall(ncp, gn, opt, logger)
			case *UncheckedBlock:
				stat.TraverseFunctionCall(ncp, gn, opt, logger)
			case *WhileStatement:
				stat.TraverseFunctionCall(ncp, gn, opt, logger)
			case *TryStatement:
				stat.TraverseFunctionCall(ncp, gn, opt, logger)
			case *DoWhileStatement:
				stat.TraverseFunctionCall(ncp, gn, opt, logger)
			}
		}
	}
}

func (b *Block) TraverseTaintOwner(opt *Option, logger logging.Logger) {
	if len(b.statements) > 0 {
		for _, statement := range b.statements {
			switch stat := statement.(type) {
			case *ExpressionStatement:
				stat.TraverseTaintOwner(opt, logger)
			case *IfStatement:
				stat.TraverseTaintOwner(opt, logger)
			case *ForStatement:
				stat.TraverseTaintOwner(opt, logger)
			case *Block:
				stat.TraverseTaintOwner(opt, logger)
			case *UncheckedBlock:
				stat.TraverseTaintOwner(opt, logger)
			case *WhileStatement:
				stat.TraverseTaintOwner(opt, logger)
			case *TryStatement:
				stat.TraverseTaintOwner(opt, logger)
			case *DoWhileStatement:
				stat.TraverseTaintOwner(opt, logger)
			}
		}
	}
}
