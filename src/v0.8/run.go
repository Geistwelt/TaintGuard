package v08

import (
	"github.com/geistwelt/logging"
	"github.com/geistwelt/taintguard/src/v0.8/ast"
	jsoniter "github.com/json-iterator/go"
)

func Run(jsonBytes []byte, logger logging.Logger) (ast.ASTNode, error) {
	fullFile := jsoniter.Get(jsonBytes)
	sourceUnit, err := ast.GetSourceUnit(fullFile, logger)
	if err != nil {
		return nil, err
	}
	return sourceUnit, nil
}
