package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type ForStatement struct {
	body                     ASTNode
	condition                ASTNode
	ID                       int `json:"id"`
	initializationExpression ASTNode
	loopExpression           ASTNode
	NodeType                 string `json:"nodeType"`
	Src                      string `json:"src"`
}

func (fs *ForStatement) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	code = code + "for"

	if fs.initializationExpression != nil || fs.condition != nil || fs.loopExpression != nil {
		code = code + " " + "("
	}

	if fs.initializationExpression != nil {
		switch initializationExpression := fs.initializationExpression.(type) {
		case *VariableDeclarationStatement:
			code = code + initializationExpression.SourceCode(false, false, indent, logger)
		default:
			if initializationExpression != nil {
				logger.Warnf("Unknown initializationExpression nodeType [%s] for ForStatement [src:%s].", initializationExpression.Type(), fs.Src)
			} else {
				logger.Warnf("Unknown initializationExpression nodeType for ForStatement [src:%s].", fs.Src)
			}
		}
	}

	code = code + ";"

	if fs.condition != nil {
		switch condition := fs.condition.(type) {
		case *BinaryOperation:
			code = code + " " + condition.SourceCode(false, false, indent, logger)
		default:
			if condition != nil {
				logger.Warnf("Unknown condition nodeType [%s] for ForStatement [src:%s].", condition.Type(), fs.Src)
			} else {
				logger.Warnf("Unknown condition nodeType for ForStatement [src:%s].", fs.Src)
			}
		}
	}

	code = code + ";"

	if fs.loopExpression != nil {
		switch loopExpression := fs.loopExpression.(type) {
		case *ExpressionStatement:
			code = code + " " + loopExpression.SourceCode(false, false, indent, logger)
		default:
			if loopExpression != nil {
				logger.Warnf("Unknown loopExpression nodeType [%s] for ForStatement [src:%s].", loopExpression.Type(), fs.Src)
			} else {
				logger.Warnf("Unknown loopExpression nodeType for ForStatement [src:%s].", fs.Src)
			}
		}
	}

	code = code + ") {\n"

	if fs.body != nil {
		switch body := fs.body.(type) {
		case *Block:
			code = code + body.SourceCode(false, false, indent, logger)
		default:
			if body != nil {
				logger.Warnf("Unknown body nodeType [%s] for ForStatement [src:%s].", body.Type(), fs.Src)
			} else {
				logger.Warnf("Unknown body nodeType for ForStatement [src:%s].", fs.Src)
			}
		}
	}

	code = code + "\n"

	if isIndent {
		code = code + indent
	}

	code = code + "}"

	return code
}

func (fs *ForStatement) Type() string {
	return fs.NodeType
}

func (fs *ForStatement) Nodes() []ASTNode {
	return nil
}

func (fs *ForStatement) NodeID() int {
	return fs.ID
}

func GetForStatement(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*ForStatement, error) {
	fs := new(ForStatement)
	if err := json.Unmarshal([]byte(raw.ToString()), fs); err != nil {
		logger.Errorf("Failed to unmarshal ForStatement: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal ForStatement: [%v]", err)
	}

	// initializationExpression
	{
		initializationExpression := raw.Get("initializationExpression")
		if initializationExpression.Size() > 0 {
			initializationExpressionNodeType := initializationExpression.Get("nodeType").ToString()
			var fsInitializationExpression ASTNode
			var err error

			switch initializationExpressionNodeType {
			case "VariableDeclarationStatement":
				fsInitializationExpression, err = GetVariableDeclarationStatement(gn, initializationExpression, logger)
			default:
				logger.Warnf("Unknown initializationExpression nodeType [%s] for ForStatement [src:%s].", initializationExpressionNodeType, fs.Src)
			}

			if err != nil {
				return nil, err
			}

			if fsInitializationExpression != nil {
				fs.initializationExpression = fsInitializationExpression
			}
		}
	}

	// condition
	{
		condition := raw.Get("condition")
		if condition.Size() > 0 {
			conditionNodeType := condition.Get("nodeType").ToString()
			var fsCondition ASTNode
			var err error

			switch conditionNodeType {
			case "BinaryOperation":
				fsCondition, err = GetBinaryOperation(gn, condition, logger)
			default:
				logger.Warnf("Unknown condition nodeType [%s] for ForStatement [src:%s].", conditionNodeType, fs.Src)
			}

			if err != nil {
				return nil, err
			}

			if fsCondition != nil {
				fs.condition = fsCondition
			}
		}
	}

	// loopExpression
	{
		loopExpression := raw.Get("loopExpression")
		if loopExpression.Size() > 0 {
			loopExpressionNodeType := loopExpression.Get("nodeType").ToString()
			var fsLoopExpression ASTNode
			var err error

			switch loopExpressionNodeType {
			case "ExpressionStatement":
				fsLoopExpression, err = GetExpressionStatement(gn, loopExpression, logger)
			default:
				logger.Warnf("Unknown loopExpression nodeType [%s] for ForStatement [src:%s].", loopExpressionNodeType, fs.Src)
			}

			if err != nil {
				return nil, err
			}

			if fsLoopExpression != nil {
				fs.loopExpression = fsLoopExpression
			}
		}
	}

	// body
	{
		body := raw.Get("body")
		if body.Size() > 0 {
			bodyNodeType := body.Get("nodeType").ToString()
			var fsBody ASTNode
			var err error

			switch bodyNodeType {
			case "Block":
				fsBody, err = GetBlock(gn, body, logger)
			default:
				logger.Warnf("Unknown body nodeType [%s] for ForStatement [src:%s].", bodyNodeType, fs.Src)
			}

			if err != nil {
				return nil, err
			}

			if fsBody != nil {
				fs.body = fsBody
			}
		}
	}

	gn.AddASTNode(fs)

	return fs, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (fs *ForStatement) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes) {
	if fs.initializationExpression != nil {
		switch initializationExpression := fs.initializationExpression.(type) {
		case *VariableDeclarationStatement:
			initializationExpression.TraverseFunctionCall(ncp, gn)
		}
	}

	if fs.condition != nil {
		switch condition := fs.condition.(type) {
		case *BinaryOperation:
			condition.TraverseFunctionCall(ncp, gn)
		}
	}

	if fs.loopExpression != nil {
		switch loopExpression := fs.loopExpression.(type) {
		case *ExpressionStatement:
			loopExpression.TraverseFunctionCall(ncp, gn)
		}
	}

	if fs.body != nil {
		switch body := fs.body.(type) {
		case *Block:
			body.TraverseFunctionCall(ncp, gn)
		}
	}
}
