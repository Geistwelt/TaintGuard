package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type ContractDefinition struct {
	Abstract                bool `json:"abstract"`
	baseContracts           []ASTNode
	CanonicalName           string   `json:"canonicalName"`
	ContractDependencies    []string `json:"contractDependencies"`
	ContractKind            string   `json:"contractKind"`
	FullyImplemented        bool     `json:"fullyImplemented"`
	ID                      int      `json:"id"`
	LinearizedBaseContracts []int    `json:"linearizedBaseContracts"`
	Name                    string   `json:"name"`
	NameLocation            string   `json:"nameLocation"`
	NodeType                string   `json:"nodeType"`
	nodes                   []ASTNode
	Scope                   int    `json:"scope"`
	Src                     string `json:"src"`
	UsedErrors              []int  `json:"usedErrors"`
	UsedEvents              []int  `json:"usedEvents"`
}

func (cd *ContractDefinition) SourceCode() string {
	var code string

	code = code + cd.ContractKind + " "
	code = code + cd.Name + " "

	if len(cd.baseContracts) > 0 {
		code = code + "is" + " "
		for index, baseContract := range cd.baseContracts {
			code = code + baseContract.SourceCode()
			if index != len(cd.baseContracts)-1 {
				code = code + "," + " "
			} else {
				code = code + " "
			}
		}
	}

	code = code + "{\n"

	return code
}

func (cd *ContractDefinition) Type() string {
	return cd.NodeType
}

func (cd *ContractDefinition) Nodes() []ASTNode {
	return cd.nodes
}

func GetContractDefinition(raw jsoniter.Any, logger logging.Logger) (*ContractDefinition, error) {
	cd := new(ContractDefinition)
	if err := json.Unmarshal([]byte(raw.ToString()), cd); err != nil {
		logger.Errorf("Failed to unmarshal ContractDefinition: [%v]", err)
		return nil, fmt.Errorf("failed to unmarshal ContractDefinition: [%v]", err)
	}
	baseContracts := raw.Get("baseContracts")
	if baseContracts.Size() > 0 {
		cd.baseContracts = make([]ASTNode, 0)
		for i := 0; i < baseContracts.Size(); i++ {
			baseContract := baseContracts.Get(i)
			baseContractNodeType := baseContract.Get("nodeType").ToString()
			switch baseContractNodeType {
			case "InheritanceSpecifier":
				bc, err := GetInheritanceSpecifier(baseContract, logger)
				if err != nil {
					return nil, err
				}
				cd.baseContracts = append(cd.baseContracts, bc)
			default:
				logger.Errorf("Unknown baseContract nodeType: [%v-%s]", baseContractNodeType, baseContract.Get("src").ToString())
			}
		}
	}

	return cd, nil
}
