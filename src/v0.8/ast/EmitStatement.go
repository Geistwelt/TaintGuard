package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type EmitStatement struct {
	eventCall ASTNode
	ID        int    `json:"id"`
	NodeType  string `json:"nodeType"`
	Src       string `json:"src"`
}

func (es *EmitStatement) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "emit"

	if es.eventCall != nil {
		switch eventCall := es.eventCall.(type) {
		case *FunctionCall:
			code = code + " " + eventCall.SourceCode(false, false, indent, logger)
		default:
			if eventCall != nil {
				logger.Warnf("Unknown eventCall nodeType [%s] for EmitStatement [src:%s].", eventCall.Type(), es.Src)
			} else {
				logger.Warnf("Unknown eventCall nodeType for EmitStatement [src:%s].", es.Src)
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (es *EmitStatement) Type() string {
	return es.NodeType
}

func (es *EmitStatement) Nodes() []ASTNode {
	return nil
}

func GetEmitStatement(raw jsoniter.Any, logger logging.Logger) (*EmitStatement, error) {
	es := new(EmitStatement)
	if err := json.Unmarshal([]byte(raw.ToString()), es); err != nil {
		logger.Errorf("Failed to unmarshal EmitStatement: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal EmitStatement: [%v]", err)
	}

	// eventCall
	{
		eventCall := raw.Get("eventCall")
		if eventCall.Size() > 0 {
			eventCallNodeType := eventCall.Get("nodeType").ToString()
			var esEventCall ASTNode
			var err error

			switch eventCallNodeType {
			case "FunctionCall":
				esEventCall, err = GetFunctionCall(eventCall, logger)
			default:
				logger.Warnf("Unknown eventCall nodeType [%s] for EmitStatement [src:%s].", eventCallNodeType, es.Src)
			}

			if err != nil {
				return nil, err
			}

			if esEventCall != nil {
				es.eventCall = esEventCall
			}
		}
	}

	return es, nil
}
