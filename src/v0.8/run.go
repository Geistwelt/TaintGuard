package v08

import (
	"fmt"

	"github.com/geistwelt/logging"
	"github.com/geistwelt/taintguard/src/v0.8/ast"
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
	i := 0
	for _, function := range gn.Functions() {
		ncp := ast.NewNormalCallPath()
		f, _ := function.(*ast.FunctionDefinition)
		f.TraverseFunctionCall(ncp, gn)
		ncps = append(ncps, ncp)
		fmt.Println(i, f.Signature())
		i++
	}

	for _, ncp := range ncps {
		test(ncp.Callees(), gn)
	}

	var path = new(string)
	ast.Path(ncps[35], path)

	fmt.Println(*path)

	return sourceUnit, nil
}

func test(ncps []*ast.NormalCallPath, gn *ast.GlobalNodes) {
	for _, ncp := range ncps {
		fd := gn.Functions()[ncp.ID()]
		f, ok := fd.(*ast.FunctionDefinition)
		if ok {
			f.TraverseFunctionCall(ncp, gn)
		}
		if len(ncp.Callees()) > 0 {
			test(ncp.Callees(), gn)
		}
	}
}
