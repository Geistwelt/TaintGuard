package v08

import (
	"fmt"

	"github.com/geistwelt/logging"
	"github.com/geistwelt/taintguard/src/v0.8/ast"
)

func IsInheritFromOwnableContract(contract *ast.ContractDefinition, gn *ast.GlobalNodes, variables []string) (bool, *ast.ContractDefinition) {
	for _, node := range contract.Nodes() {
		switch n := node.(type) {
		case *ast.VariableDeclaration:
			// Here the parameter names can be defined via the command line.
			for _, variable := range variables {
				if n.Name == variable {
					return false, nil
				}
			}
		}
	}

	for _, contractNodeID := range contract.LinearizedBaseContracts {
		if contractNodeID == contract.NodeID() {
			continue
		}
		baseContract := gn.ContractsByID()[contractNodeID].(*ast.ContractDefinition)
		if baseContract.Name == "Ownable" {
			return true, baseContract
		} else {
			for _, node := range baseContract.Nodes() {
				switch n := node.(type) {
				case *ast.VariableDeclaration:
					for _, variable := range variables {
						if n.Name == variable {
							return true, baseContract
						}
					}
				}
			}
		}
	}

	return false, nil
}

func InstrumentCodeForOwner(contract *ast.ContractDefinition, variables []string) string {
	var ownerVariableName string
	for _, node := range contract.Nodes() {
		if node.Type() == "VariableDeclaration" {
			// check
			vdNode, _ := node.(*ast.VariableDeclaration)
			// Here the parameter names can be defined via the command line.
			var ok bool
			for _, variable := range variables {
				if vdNode.Name == variable {
					ok = true
				}
			}
			if ok {
				ownerVariableName = vdNode.Name

				instReturnOwnerFunction := &ast.FunctionDefinition{
					Implemented:     true,
					Kind:            "function",
					Name:            fmt.Sprintf("xxx_track_func_%s", ownerVariableName),
					NodeType:        "FunctionDefinition",
					Src:             "xxx",
					StateMutability: "view",
					Visibility:      "internal",
				}
				body := &ast.Block{
					NodeType: "Block",
					Src:      "xxx",
				}
				statement := &ast.Return{
					NodeType: "Return",
					Src:      "xxx",
				}
				expression := &ast.Identifier{
					Name:     ownerVariableName,
					NodeType: "Identifier",
					Src:      "xxx",
				}
				statement.SetExpression(expression)
				body.AppendStatement(statement)
				instReturnOwnerFunction.SetBody(body)
				returnParameters := &ast.ParameterList{
					NodeType: "ParameterList",
					Src:      "xxx",
				}
				parameter := &ast.VariableDeclaration{
					NodeType:        "VariableDeclaration",
					Src:             "xxx",
					StorageLocation: "default",
					Visibility:      "internal",
				}
				typeName := &ast.ElementaryTypeName{
					Name:            "address",
					NodeType:        "ElementaryTypeName",
					Src:             "xxx",
					StateMutability: "nonpayable",
				}
				parameter.SetTypeName(typeName)
				returnParameters.AppendParameter(parameter)
				instReturnOwnerFunction.SetReturnParameters(returnParameters)
				contract.InsertReturnOwnerFunction(instReturnOwnerFunction)

				protect1 := &ast.VariableDeclaration{
					Mutability:      "mutable",
					Name:            fmt.Sprintf("xxx_track_%s", vdNode.Name),
					NodeType:        "VariableDeclaration",
					Scope:           contract.NodeID(),
					Src:             "xxx",
					StateVariable:   true,
					StorageLocation: "default",
					Visibility:      "internal",
				}
				protect1Etn := &ast.ElementaryTypeName{
					Name:     "bytes",
					NodeType: "ElementaryTypeName",
					Src:      "xxx",
				}
				protect1.SetTypeName(protect1Etn)

				protect2 := &ast.VariableDeclaration{
					Mutability:      "mutable",
					Name:            fmt.Sprintf("xxx_track_mapping_%s", vdNode.Name),
					NodeType:        "VariableDeclaration",
					Src:             "xxx",
					StateVariable:   true,
					StorageLocation: "default",
					Visibility:      "internal",
				}

				protect2m := &ast.Mapping{
					NodeType: "Mapping",
					Src:      "xxx",
				}

				protect2Kt := &ast.ElementaryTypeName{
					Name:     "bytes",
					NodeType: "ElementaryTypeName",
					Src:      "xxx",
				}

				protect2Vt := &ast.ElementaryTypeName{
					Name:            "address",
					NodeType:        "ElementaryTypeName",
					Src:             "xxx",
					StateMutability: "nonpayable",
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
						TrackOwnerVariableName:   fmt.Sprintf("xxx_track_%s", vdNode.Name),
						TrackOwnerMappingName:    fmt.Sprintf("xxx_track_mapping_%s", vdNode.Name),
						SimilarOwnerVariableName: vdNode.Name,
					}, logging.MustNewLogger())
				}
			}
		}
	}
	return ownerVariableName
}

func InstrumentCodeForAssert(ownerVariableName string, contract *ast.ContractDefinition) {
	expressionStatement := &ast.ExpressionStatement{
		NodeType: "ExpressionStatement",
		Src:      "xxx",
	}
	functionCall := &ast.FunctionCall{
		Kind:     "functionCall",
		NodeType: "FunctionCall",
		Src:      "xxx",
	}
	functionCallExpression := &ast.Identifier{
		ArgumentTypes: []struct {
			TypeIdentifier string `json:"typeIdentifier"`
			TypeString     string `json:"typeString"`
		}{{TypeIdentifier: "t_bool", TypeString: "bool"}},
		Name:     "assert",
		NodeType: "Identifier",
		Src:      "xxx",
	}
	functionCallArgument := &ast.BinaryOperation{
		NodeType: "BinaryOperation",
		Operator: "==",
		Src:      "xxx",
	}
	binaryOperationRightExpression := &ast.FunctionCall{
		Kind:     "functionCall",
		NodeType: "FunctionCall",
		Src:      "xxx",
	}
	binaryOperationRightExpressionFunctionCallExpression := &ast.Identifier{
		Name:     fmt.Sprintf("xxx_track_func_%s", ownerVariableName),
		NodeType: "Identifier",
		Src:      "xxx",
	}
	binaryOperationLeftExpression := &ast.IndexAccess{
		NodeType: "IndexAccess",
		Src:      "xxx",
	}
	binaryOperationLeftExpressionBaseExpression := &ast.Identifier{
		Name:     fmt.Sprintf("xxx_track_mapping_%s", ownerVariableName),
		NodeType: "Identifier",
		Src:      "xxx",
	}
	binaryOperationLeftExpressionIndexExpression := &ast.Identifier{
		Name:     fmt.Sprintf("xxx_track_%s", ownerVariableName),
		NodeType: "Identifier",
		Src:      "xxx",
	}
	expressionStatement.SetExpression(functionCall)
	functionCall.SetExpression(functionCallExpression)
	functionCallArgument.SetLeftExpression(binaryOperationLeftExpression)
	functionCallArgument.SetRightExpression(binaryOperationRightExpression)
	binaryOperationRightExpression.SetExpression(binaryOperationRightExpressionFunctionCallExpression)
	binaryOperationLeftExpression.SetBaseExpression(binaryOperationLeftExpressionBaseExpression)
	binaryOperationLeftExpression.SetIndexExpression(binaryOperationLeftExpressionIndexExpression)
	functionCall.AppendArgument(functionCallArgument)

	contract.TraverseDelegatecall(&ast.Option{ExpressionStatement: expressionStatement}, logging.MustNewLogger())
}

func VerifyVariableDeclarationOrder(callerContract, calleeContract *ast.ContractDefinition, gn *ast.GlobalNodes, variables []string) bool {
	var callerContractVariables []*variable = make([]*variable, 0) // variableName => variableType
	var calleeContractVariables []*variable = make([]*variable, 0) // variableName => variableType
	var same bool = false

	// caller contract
	{
		lookupBaseContract(callerContract, &callerContractVariables, gn)
	}

	// callee contract
	{
		lookupBaseContract(calleeContract, &calleeContractVariables, gn)
	}

	for i := 0; i < len(callerContractVariables) && i < len(calleeContractVariables); i++ {
		for _, variable := range variables {
			if calleeContractVariables[i].variableName == variable && 
			calleeContractVariables[i].variableType == callerContractVariables[i].variableType {
				same = true
			}
		}
	}

	return same
}

func lookupBaseContract(contract *ast.ContractDefinition, variables *[]*variable, gn *ast.GlobalNodes) {
	for _, contractID := range contract.LinearizedBaseContracts {
		if contractID == contract.NodeID() {
			continue
		}
		lookupBaseContract(gn.ContractsByID()[contractID].(*ast.ContractDefinition), variables, gn)
	}
	for _, node := range contract.Nodes() {
		switch n := node.(type) {
		case *ast.VariableDeclaration:
			if n.StateVariable {
				*variables = append(*variables, &variable{n.Name, n.TypeDescriptions.TypeString})
			}
		}
	}
}

type variable struct {
	variableName string
	variableType string
}