package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type TryStatement struct {
	clauses      []ASTNode
	externalCall ASTNode
	ID           int    `json:"id"`
	NodeType     string `json:"nodeType"`
	Src          string `json:"src"`
}

func (ts *TryStatement) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "try"

	if ts.externalCall != nil {
		switch externalCall := ts.externalCall.(type) {
		case *FunctionCall:
			code = code + " " + externalCall.SourceCode(false, false, indent, logger)
		default:
			if externalCall != nil {
				logger.Warnf("Unknown externalCall nodeType [%s] for TryStatement [src:%s].", externalCall.Type(), ts.Src)
			} else {
				logger.Warnf("Unknown externalCall nodeType for TryStatement [src:%s].", ts.Src)
			}
		}
	}

	if len(ts.clauses) > 0 {
		for index, clause := range ts.clauses {
			switch c := clause.(type) {
			case *TryCatchClause:
				if index == 0 {
					code = code + " " + "returns" + " " + c.SourceCode(false, true, indent, logger)
				} else {
					code = code + " " + "catch" + " " + c.SourceCode(false, true, indent, logger)
				}
			default:
				if c != nil {
					logger.Warnf("Unknown clause nodeType [%s] for TryStatement [src:%s].", c.Type(), ts.Src)
				} else {
					logger.Warnf("Unknown clause nodeType for TryStatement [src:%s].", ts.Src)
				}
			}
		}
	}

	return code
}

func (ts *TryStatement) Type() string {
	return ts.NodeType
}

func (ts *TryStatement) Nodes() []ASTNode {
	return ts.clauses
}

func (ts *TryStatement) NodeID() int {
	return ts.ID
}

func GetTryStatement(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*TryStatement, error) {
	ts := new(TryStatement)
	if err := json.Unmarshal([]byte(raw.ToString()), ts); err != nil {
		logger.Errorf("Failed to unmarshal TryStatement: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal TryStatement: [%v]", err)
	}

	// clauses
	{
		clauses := raw.Get("clauses")
		if clauses.Size() > 0 {
			ts.clauses = make([]ASTNode, 0)

			for i := 0; i < clauses.Size(); i++ {
				clause := clauses.Get(i)
				if clause.Size() > 0 {
					clauseNodeType := clause.Get("nodeType").ToString()
					var tsClause ASTNode
					var err error

					switch clauseNodeType {
					case "TryCatchClause":
						tsClause, err = GetTryCatchClause(gn, clause, logger)
					default:
						logger.Warnf("Unknown clause nodeType [%s] for TryStatement [src:%s].", clauseNodeType, ts.Src)
					}

					if err != nil {
						return nil, err
					}

					if tsClause != nil {
						ts.clauses = append(ts.clauses, tsClause)
					}
				} else {
					logger.Warnf("Clause in TryStatement [src:%s] should not be empty.", ts.Src)
				}
			}
		}
	}

	// externalCall
	{
		externalCall := raw.Get("externalCall")
		if externalCall.Size() > 0 {
			externalCallNodeType := externalCall.Get("nodeType").ToString()
			var tsExternalCall ASTNode
			var err error

			switch externalCallNodeType {
			case "FunctionCall":
				tsExternalCall, err = GetFunctionCall(gn, externalCall, logger)
			default:
				logger.Warnf("Unknown externalCall nodeType [%s] for TryStatement [src:%s].", externalCallNodeType, ts.Src)
			}

			if err != nil {
				return nil, err
			}

			if tsExternalCall != nil {
				ts.externalCall = tsExternalCall
			}
		}
	}

	gn.AddASTNode(ts)

	return ts, nil
}

func (ts *TryStatement) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	if ts.externalCall != nil {
		switch externalCall := ts.externalCall.(type) {
		case *FunctionCall:
			externalCall.TraverseFunctionCall(ncp, gn, opt, logger)
		}
	}

	if len(ts.clauses) > 0 {
		for _, clause := range ts.clauses {
			switch c := clause.(type) {
			case *TryCatchClause:
				c.TraverseFunctionCall(ncp, gn, opt, logger)
			}
		}
	}
}

func (ts *TryStatement) TraverseTaintOwner(opt *Option, logger logging.Logger) {
	if len(ts.clauses) > 0 {
		for _, clause := range ts.clauses {
			switch c := clause.(type) {
			case *TryCatchClause:
				c.TraverseTaintOwner(opt, logger)
			}
		}
	}
}
