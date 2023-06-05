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
				code = code + stat.SourceCode(true, true, indent+"\t", logger)
			case *PlaceholderStatement:
				code = code + stat.SourceCode(true, true, indent+"\t", logger)
			case *Return:
				code = code + stat.SourceCode(true, true, indent+"\t", logger)
			case *EmitStatement:
				code = code + stat.SourceCode(true, true, indent+"\t", logger)
			case *IfStatement:
				code = code + stat.SourceCode(false, true, indent+"\t", logger)
			case *VariableDeclarationStatement:
				code = code + stat.SourceCode(true, true, indent+"\t", logger)
			case *InlineAssembly:
				code = code + stat.SourceCode(false, true, indent+"\t", logger)
			case *ForStatement:
				code = code + stat.SourceCode(false, true, indent+"\t", logger)
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

func GetBlock(raw jsoniter.Any, logger logging.Logger) (*Block, error) {
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
					bStatement, err = GetPlaceholderStatement(statement, logger)
				case "ExpressionStatement":
					bStatement, err = GetExpressionStatement(statement, logger)
				case "Return":
					bStatement, err = GetReturn(statement, logger)
				case "EmitStatement":
					bStatement, err = GetEmitStatement(statement, logger)
				case "IfStatement":
					bStatement, err = GetIfStatement(statement, logger)
				case "VariableDeclarationStatement":
					bStatement, err = GetVariableDeclarationStatement(statement, logger)
				case "InlineAssembly":
					bStatement, err = GetInlineAssembly(statement, logger)
				case "ForStatement":
					bStatement, err = GetForStatement(statement, logger)
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

	return b, nil
}
