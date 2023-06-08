package v08

import (
	"fmt"

	"github.com/geistwelt/logging"
	"github.com/geistwelt/taintguard/src/v0.8/ast"
)

func IsInheritFromOwnableContract(contract *ast.ContractDefinition, gn *ast.GlobalNodes) (bool, *ast.ContractDefinition) {
	for _, contractNodeID := range contract.LinearizedBaseContracts {
		if contractNodeID == contract.NodeID() {
			continue
		}
		baseContract := gn.Contracts()[contractNodeID].(*ast.ContractDefinition)
		if baseContract.Name == "Ownable" {
			return true, baseContract
		}
	}
	return false, nil
}

func InstrumentCode(contract *ast.ContractDefinition) {
	for _, node := range contract.Nodes() {
		if node.Type() == "VariableDeclaration" {
			// check
			vdNode, _ := node.(*ast.VariableDeclaration)
			// Here the parameter names can be defined via the command line.
			if vdNode.Name == "owner" || vdNode.Name == "_owner" {
				protect1 := &ast.VariableDeclaration{
					Constant:        false,
					ID:              -100,
					Indexed:         false,
					Mutability:      "mutable",
					Name:            fmt.Sprintf("xxx_track_%s", vdNode.Name),
					NameLocation:    "",
					NodeType:        "VariableDeclaration",
					Scope:           contract.NodeID(),
					Src:             "xxx",
					StateVariable:   true,
					StorageLocation: "default",
					TypeDescriptions: struct {
						TypeIdentifier string "json:\"typeIdentifier\""
						TypeString     string "json:\"typeString\""
					}{TypeIdentifier: "t_bytes_storage", TypeString: "bytes"},
					Visibility: "private",
				}
				protect1Etn := &ast.ElementaryTypeName{
					ID:              -101,
					Name:            "bytes",
					NodeType:        "ElementaryTypeName",
					Src:             "xxx",
					StateMutability: "",
					TypeDescriptions: struct {
						TypeIdentifier string "json:\"typeIdentifier\""
						TypeString     string "json:\"typeString\""
					}{TypeIdentifier: "t_bytes_storage_ptr", TypeString: "bytes"},
				}
				protect1.SetTypeName(protect1Etn)

				protect2 := &ast.VariableDeclaration{
					Constant:        false,
					ID:              -102,
					Indexed:         false,
					Mutability:      "mutable",
					Name:            fmt.Sprintf("xxx_track_mapping_%s", vdNode.Name),
					NameLocation:    "",
					NodeType:        "VariableDeclaration",
					Scope:           0,
					Src:             "xxx",
					StateVariable:   true,
					StorageLocation: "default",
					TypeDescriptions: struct {
						TypeIdentifier string "json:\"typeIdentifier\""
						TypeString     string "json:\"typeString\""
					}{TypeIdentifier: "t_mapping$_t_bytes_memory_ptr_$_t_address_$", TypeString: "mapping(bytes => address)"},
					Visibility: "private",
				}

				protect2m := &ast.Mapping{
					ID:       -103,
					KeyName:  "",
					NodeType: "Mapping",
					Src:      "xxx",
					TypeDescriptions: struct {
						TypeIdentifier string "json:\"typeIdentifier\""
						TypeString     string "json:\"typeString\""
					}{TypeIdentifier: "t_mapping$_t_bytes_memory_ptr_$_t_address_$", TypeString: "mapping(bytes => address)"},
					ValueName:         "",
					ValueNameLocation: "",
				}

				protect2Kt := &ast.ElementaryTypeName{
					ID:              -104,
					Name:            "bytes",
					NodeType:        "ElementaryTypeName",
					Src:             "xxx",
					StateMutability: "",
					TypeDescriptions: struct {
						TypeIdentifier string "json:\"typeIdentifier\""
						TypeString     string "json:\"typeString\""
					}{TypeIdentifier: "t_bytes_storage_ptr", TypeString: "bytes"},
				}

				protect2Vt := &ast.ElementaryTypeName{
					ID:              -105,
					Name:            "address",
					NodeType:        "ElementaryTypeName",
					Src:             "xxx",
					StateMutability: "nonpayable",
					TypeDescriptions: struct {
						TypeIdentifier string "json:\"typeIdentifier\""
						TypeString     string "json:\"typeString\""
					}{TypeIdentifier: "t_address", TypeString: "address"},
				}

				protect2m.SetKeyType(protect2Kt)
				protect2m.SetValueType(protect2Vt)
				protect2.SetTypeName(protect2m)

				var isExist bool = false
				for _, node_ := range contract.Nodes() {
					if node_.Type() == "VariableDeclaration" {
						if node_.SourceCode(false, false, "", nil) == protect1.SourceCode(false, false, "", nil) {
							isExist = true
						}
					}
				}
				if !isExist {
					contract.AppendNode(protect1)
					contract.AppendNode(protect2)
					contract.TraverseTaintOwner(&ast.Option{
						TrackFunctionDefinitionName: "",
						TrackOwnerVariableName:      fmt.Sprintf("xxx_track_%s", vdNode.Name),
						TrackOwnerMappingName:       fmt.Sprintf("xxx_track_mapping_%s", vdNode.Name),
						SimilarOwnerVariableName:    vdNode.Name,
					}, logging.MustNewLogger())
				}
			}
		}
	}
}
