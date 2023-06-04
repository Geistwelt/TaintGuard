package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type ParameterList struct {
	ID         int    `json:"id"`
	NodeType   string `json:"nodeType"`
	parameters []ASTNode
	Src        string `json:"src"`
}

func (pl *ParameterList) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	if len(pl.parameters) > 0 {
		for index, parameter := range pl.parameters {
			switch p := parameter.(type) {
			case *VariableDeclaration:
				code = code + p.SourceCode(false, false, indent, logger)
			default:
				if p != nil {
					logger.Warnf("Unknown parameter nodeType [%s] for ParameterList [src:%s].", p.Type(), pl.Src)
				} else {
					logger.Warnf("Unknown parameter nodeType for ParameterList [src:%s].", pl.Src)
				}
			}
			if index < len(pl.parameters)-1 {
				code = code + ", "
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (pl *ParameterList) Type() string {
	return pl.NodeType
}

func (pl *ParameterList) Nodes() []ASTNode {
	return pl.parameters
}

func GetParameterList(raw jsoniter.Any, logger logging.Logger) (*ParameterList, error) {
	pl := new(ParameterList)
	if err := json.Unmarshal([]byte(raw.ToString()), pl); err != nil {
		logger.Errorf("Failed to unmarshal ParameterList: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal ParameterList: [%v]", err)
	}

	parameters := raw.Get("parameters")
	if parameters.Size() > 0 {
		pl.parameters = make([]ASTNode, 0)
	}

	for i := 0; i < parameters.Size(); i++ {
		parameter := parameters.Get(i)
		parameterNodeType := parameter.Get("nodeType").ToString()
		var plParameter ASTNode
		var err error

		switch parameterNodeType {
		case "VariableDeclaration":
			plParameter, err = GetVariableDeclaration(parameter, logger)
		default:
			logger.Warnf("Unknown parameter nodeType [%s] for ParameterList [src:%s].", parameterNodeType, pl.Src)
		}

		if err != nil {
			return nil, err
		}
		if plParameter != nil {
			pl.parameters = append(pl.parameters, plParameter)
		}
	}

	return pl, nil
}
