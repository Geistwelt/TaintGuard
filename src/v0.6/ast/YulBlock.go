package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type YulBlock struct {
	NodeType   string `json:"nodeType"`
	Src        string `json:"src"`
	statements []ASTNode
}

func (yb *YulBlock) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	if len(yb.statements) > 0 {
		for index, statement := range yb.statements {
			switch stat := statement.(type) {
			case *YulAssignment:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			case *YulExpressionStatement:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			case *YulVariableDeclaration:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			case *YulForLoop:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			case *YulSwitch:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			case *YulIf:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			case *YulBreak:
				code = code + stat.SourceCode(false, true, indent+"    ", logger)
			default:
				if stat != nil {
					logger.Warnf("Unknown statement nodeType [%s] for YulBlock [src:%s].", stat.Type(), yb.Src)
				} else {
					logger.Warnf("Unknown statement nodeType for YulBlock [src:%s].", yb.Src)
				}
			}

			if index < len(yb.statements)-1 {
				code = code + "\n"
			}
		}
	}

	return code
}

func (yb *YulBlock) Type() string {
	return yb.NodeType
}

func (yb *YulBlock) Nodes() []ASTNode {
	return yb.statements
}

func (yb *YulBlock) NodeID() int {
	return -1
}

func GetYulBlock(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*YulBlock, error) {
	yb := new(YulBlock)
	if err := json.Unmarshal([]byte(raw.ToString()), yb); err != nil {
		logger.Errorf("Failed to unmarshal YulBlock: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal YulBlock: [%v]", err)
	}

	// statements
	{
		statements := raw.Get("statements")
		if statements.Size() > 0 {
			yb.statements = make([]ASTNode, 0)

			for i := 0; i < statements.Size(); i++ {
				statement := statements.Get(i)
				if statement.Size() > 0 {
					statementNodeType := statement.Get("nodeType").ToString()
					var ybStatement ASTNode
					var err error

					switch statementNodeType {
					case "YulAssignment":
						ybStatement, err = GetYulAssignment(gn, statement, logger)
					case "YulExpressionStatement":
						ybStatement, err = GetYulExpressionStatement(gn, statement, logger)
					case "YulVariableDeclaration":
						ybStatement, err = GetYulVariableDeclaration(gn, statement, logger)
					case "YulForLoop":
						ybStatement, err = GetYulForLoop(gn, statement, logger)
					case "YulSwitch":
						ybStatement, err = GetYulSwitch(gn, statement, logger)
					case "YulIf":
						ybStatement, err = GetYulIf(gn, statement, logger)
					case "YulBreak":
						ybStatement, err = GetYulBreak(gn, statement, logger)
					default:
						logger.Warnf("Unknown statement nodeType [%s] for YulBlock [src:%s].", statementNodeType, yb.Src)
					}

					if err != nil {
						return nil, err
					}

					if ybStatement != nil {
						yb.statements = append(yb.statements, ybStatement)
					}
				} else {
					logger.Warnf("Statement in YulBlock [src:%s] should not be empty.", yb.Src)
				}
			}
		}
	}

	return yb, nil
}
