package ast

type ASTNode interface {
	Type() string
	SourceCode() string
	Nodes() []ASTNode
}