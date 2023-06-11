package v05

import (
	"github.com/geistwelt/logging"
	"github.com/geistwelt/taintguard/src/v0.5/ast"
	"github.com/geistwelt/taintguard/src/v0.5/cfg"

	// "github.com/geistwelt/taintguard/src/v0.8/cfg"
	jsoniter "github.com/json-iterator/go"
)

func Run(jsonBytes []byte, isCfg bool, logger logging.Logger) (ast.ASTNode, error) {
	gn := ast.NewGlobalNodes()
	fullFile := jsoniter.Get(jsonBytes)
	sourceUnit, err := ast.GetSourceUnit(gn, fullFile, logger)
	if err != nil {
		return nil, err
	}

	// Get the call path of each function.
	ncps := make([]*ast.NormalCallPath, 0)
	for _, function := range gn.Functions() {
		ncp := ast.NewNormalCallPath()
		f, _ := function.(*ast.FunctionDefinition)
		opt := &ast.Option{}
		opt.MakeDelegatecallUnknownContractCh(1)
		opt.MakeDelegatecallKnownContractCh(1)

		f.TraverseFunctionCall(ncp, gn, opt, logger)
		ncps = append(ncps, ncp)
		select {
		case <-opt.DelegatecallUnknownContractCh():
			contract := gn.ContractsByID()[f.Scope].(*ast.ContractDefinition)
			logger.Infof("Contract [%s] should be instrumented directly, because it delegatecall to unknown contract.", contract.Name)
			if ok, c := IsInheritFromOwnableContract(contract, gn); ok {
				ownerVariableName := InstrumentCodeForOwner(c)
				InstrumentCodeForAssert(ownerVariableName, contract)
			} else {
				ownerVariableName := InstrumentCodeForOwner(contract)
				InstrumentCodeForAssert(ownerVariableName, contract)
			}
		case calleeContractName := <-opt.DelegatecallKnownContractCh():
			callerContract := gn.ContractsByID()[f.Scope].(*ast.ContractDefinition)
			calleeContract := gn.ContractsByName()[calleeContractName].(*ast.ContractDefinition)
			if VerifyVariableDeclarationOrder(callerContract, calleeContract, gn) {
				if ok, c := IsInheritFromOwnableContract(callerContract, gn); ok {
					ownerVariableName := InstrumentCodeForOwner(c)
					InstrumentCodeForAssert(ownerVariableName, callerContract)
				} else {
					ownerVariableName := InstrumentCodeForOwner(callerContract)
				InstrumentCodeForAssert(ownerVariableName, callerContract)
				logger.Debug("在本地插桩")
				}
			} else {
				logger.Debug("No instrumentation protection required.")
			}

		default:
		}
	}

	if isCfg {
		for _, ncp := range ncps {
			TraverseFunctionCallAll(ncp.Callees(), gn, logger)
		}

		for _, ncp := range ncps {
			cfg.MakeCG(ncp, logger)
		}
	}

	return sourceUnit, nil
}

func TraverseFunctionCallAll(ncps []*ast.NormalCallPath, gn *ast.GlobalNodes, logger logging.Logger) {
	for _, ncp := range ncps {
		fd := gn.Functions()[ncp.ID()]
		f, ok := fd.(*ast.FunctionDefinition)
		if ok {
			f.TraverseFunctionCall(ncp, gn, nil, logger)
		}
		if len(ncp.Callees()) > 0 {
			TraverseFunctionCallAll(ncp.Callees(), gn, logger)
		}
	}
}
