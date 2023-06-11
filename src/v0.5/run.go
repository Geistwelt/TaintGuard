package v05

import (
	"github.com/geistwelt/logging"
	"github.com/geistwelt/taintguard/src/v0.5/ast"
	jsoniter "github.com/json-iterator/go"
)

func Run(jsonBytes []byte, isCfg bool, logger logging.Logger) (ast.ASTNode, error) {
	gn := ast.NewGlobalNodes()
	fullFile := jsoniter.Get(jsonBytes)
	sourceUnit, err := ast.GetSourceUnit(gn, fullFile, logger)
	if err != nil {
		return nil, err
	}

	return sourceUnit, nil
}
