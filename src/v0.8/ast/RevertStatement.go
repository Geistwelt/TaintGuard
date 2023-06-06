package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type RevertStatement struct {
	errorCall ASTNode
	ID        int    `json:"id"`
	NodeType  string `json:"nodeType"`
	Src       string `json:"src"`
}

func (rs *RevertStatement) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "revert"

	if rs.errorCall != nil {
		switch errorCall := rs.errorCall.(type) {
		case *FunctionCall:
			code = code + " " + errorCall.SourceCode(false, false, indent, logger)
		default:
			if errorCall != nil {
				logger.Warnf("Unknown errorCall nodeType [%s] for RevertStatement [src:%s].", errorCall.Type(), rs.Src)
			} else {
				logger.Warnf("Unknown errorCall nodeType for RevertStatement [src:%s].", rs.Src)
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (rs *RevertStatement) Type() string {
	return rs.NodeType
}

func (rs *RevertStatement) Nodes() []ASTNode {
	return nil
}

func (rs *RevertStatement) NodeID() int {
	return rs.ID
}

func GetRevertStatement(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*RevertStatement, error) {
	rs := new(RevertStatement)
	if err := json.Unmarshal([]byte(raw.ToString()), rs); err != nil {
		logger.Errorf("Failed to unmarshal RevertStatement: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal RevertStatement: [%v]", err)
	}

	// errorCall
	{
		errorCall := raw.Get("errorCall")
		if errorCall.Size() > 0 {
			errorCallNodeType := errorCall.Get("nodeType").ToString()
			var rsErrorCall ASTNode
			var err error

			switch errorCallNodeType {
			case "FunctionCall":
				rsErrorCall, err = GetFunctionCall(gn, errorCall, logger)
			default:
				logger.Warnf("Unknown errorCall nodeType [%s] for RevertStatement [src:%s].", errorCallNodeType, rs.Src)
			}

			if err != nil {
				return nil, err
			}

			if rsErrorCall != nil {
				rs.errorCall = rsErrorCall
			}
		}
	}

	gn.AddASTNode(rs)

	return rs, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (rs *RevertStatement) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes) {
	if rs.errorCall != nil {
		switch errorCall := rs.errorCall.(type) {
		case *FunctionCall:
			errorCall.TraverseFunctionCall(ncp, gn)
		}
	}
}