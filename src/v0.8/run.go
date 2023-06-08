package v08

import (
	"github.com/geistwelt/logging"
	"github.com/geistwelt/taintguard/src/v0.8/ast"
	"github.com/geistwelt/taintguard/src/v0.8/cfg"

	// "github.com/geistwelt/taintguard/src/v0.8/cfg"
	jsoniter "github.com/json-iterator/go"
)

func Run(jsonBytes []byte, logger logging.Logger) (ast.ASTNode, error) {
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
		f.TraverseFunctionCall(ncp, gn)
		ncps = append(ncps, ncp)
	}

	for _, ncp := range ncps {
		TraverseFunctionCallAll(ncp.Callees(), gn)
	}

	for _, ncp := range ncps {
		cfg.MakeCG(ncp, logger)
	}

	return sourceUnit, nil
}

func TraverseFunctionCallAll(ncps []*ast.NormalCallPath, gn *ast.GlobalNodes) {
	for _, ncp := range ncps {
		fd := gn.Functions()[ncp.ID()]
		f, ok := fd.(*ast.FunctionDefinition)
		if ok {
			f.TraverseFunctionCall(ncp, gn)
		}
		if len(ncp.Callees()) > 0 {
			TraverseFunctionCallAll(ncp.Callees(), gn)
		}
	}
}
