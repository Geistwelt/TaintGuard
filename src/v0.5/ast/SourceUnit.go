package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type SourceUnit struct {
	AbsolutePath    string           `json:"absolutePath"`
	ExportedSymbols map[string][]int `json:"exportedSymbols"`
	ID              int              `json:"id"`
	NodeType        string           `json:"nodeType"`
	Src             string           `json:"src"`
	nodes           []ASTNode
}

func (su *SourceUnit) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	
	for _, node := range su.nodes {
		switch node.Type() {
		case "PragmaDirective":
			code = code + node.SourceCode(false, false, indent, logger) + "\n"
		case "ContractDefinition":
			code = code + node.SourceCode(false, false, indent, logger) + "\n"
		case "ErrorDefinition":
			code = code + node.SourceCode(true, false, indent, logger) + "\n"
		case "VariableDeclaration":
			code = code + node.SourceCode(true, false, indent, logger) + "\n"
		default:
			logger.Warnf("Unknown node nodeType [%s] for SourceUnit [src:%s].", node.Type(), su.Src)
		}
	}

	return code
}

func (su *SourceUnit) Type() string {
	return su.NodeType
}

func (su *SourceUnit) Nodes() []ASTNode {
	return su.nodes
}

func (su *SourceUnit) NodeID() int {
	return su.ID
}

func (su *SourceUnit) AppendNode(node ASTNode) {
	if su.nodes == nil {
		su.nodes = make([]ASTNode, 0)
	}
	su.nodes = append(su.nodes, node)
}

func GetSourceUnit(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*SourceUnit, error) {
	su := new(SourceUnit)
	if err := json.Unmarshal([]byte(raw.ToString()), su); err != nil {
		logger.Errorf("Failed to unmarshal SourceUnit: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal SourceUnit: [%v]", err)
	}

	sourceUnitNodes := raw.Get("nodes")
	for i := 0; i < sourceUnitNodes.Size(); i++ {
		sourceUnitChild := sourceUnitNodes.Get(i)
		sourceUnitChildType := sourceUnitChild.Get("nodeType").ToString()
		switch sourceUnitChildType {
		case "PragmaDirective":
			pragmaDirective, err := GetPragmaDirective(gn, sourceUnitChild, logger)
			if err != nil {
				return nil, err
			}
			su.AppendNode(pragmaDirective)
		case "ContractDefinition":
			contractDefinition, err := GetContractDefinition(gn, sourceUnitChild, logger)
			if err != nil {
				return nil, err
			}
			su.AppendNode(contractDefinition)
		// case "ErrorDefinition":
		// 	errorDefinition, err := GetErrorDefinition(gn, sourceUnitChild, logger)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	su.AppendNode(errorDefinition)
		case "VariableDeclaration":
			variableDeclaration, err := GetVariableDeclaration(gn, sourceUnitChild, logger)
			if err != nil {
				return nil, err
			}
			su.AppendNode(variableDeclaration)
		default:
			logger.Warnf("Unknown node nodeType [%s] for SourceUnit [src:%s].", sourceUnitChildType, su.Src)
		}

	}

	gn.AddASTNode(su)
	return su, nil
}
