package ast

import "github.com/geistwelt/logging"

type ASTNode interface {
	Type() string
	SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string
	Nodes() []ASTNode
}
