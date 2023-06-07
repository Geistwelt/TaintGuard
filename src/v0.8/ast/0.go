package ast

import (
	"sync"

	"github.com/geistwelt/logging"
)

type ASTNode interface {
	Type() string
	SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string
	Nodes() []ASTNode
	NodeID() int
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type traverseFunctionCall interface {
	TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes)
}

var _ traverseFunctionCall = (*Assignment)(nil)
var _ traverseFunctionCall = (*BinaryOperation)(nil)
var _ traverseFunctionCall = (*Block)(nil)
var _ traverseFunctionCall = (*Conditional)(nil)
var _ traverseFunctionCall = (*EmitStatement)(nil)
var _ traverseFunctionCall = (*EventDefinition)(nil)
var _ traverseFunctionCall = (*ExpressionStatement)(nil)
var _ traverseFunctionCall = (*ForStatement)(nil)
var _ traverseFunctionCall = (*FunctionCall)(nil)
var _ traverseFunctionCall = (*FunctionDefinition)(nil)
var _ traverseFunctionCall = (*IfStatement)(nil)
var _ traverseFunctionCall = (*IndexAccess)(nil)
var _ traverseFunctionCall = (*Mapping)(nil)
var _ traverseFunctionCall = (*MemberAccess)(nil)
var _ traverseFunctionCall = (*ParameterList)(nil)
var _ traverseFunctionCall = (*Return)(nil)
var _ traverseFunctionCall = (*TupleExpression)(nil)
var _ traverseFunctionCall = (*UnaryOperation)(nil)
var _ traverseFunctionCall = (*VariableDeclaration)(nil)
var _ traverseFunctionCall = (*VariableDeclarationStatement)(nil)
var _ traverseFunctionCall = (*UncheckedBlock)(nil)
var _ traverseFunctionCall = (*WhileStatement)(nil)
var _ traverseFunctionCall = (*FunctionCallOptions)(nil)
var _ traverseFunctionCall = (*TryCatchClause)(nil)
var _ traverseFunctionCall = (*TryStatement)(nil)
var _ traverseFunctionCall = (*DoWhileStatement)(nil)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Analysis

type GlobalNodes struct {
	nodes     map[int]ASTNode // id => all ASTNode
	contracts map[int]ASTNode // id => all ContractDefinition
	functions map[int]ASTNode // id => all FunctionDefinition
	mu        sync.RWMutex
}

func NewGlobalNodes() *GlobalNodes {
	gn := new(GlobalNodes)
	gn.nodes = make(map[int]ASTNode)
	gn.contracts = make(map[int]ASTNode)
	gn.functions = make(map[int]ASTNode)
	return gn
}

func (gn *GlobalNodes) AddASTNode(node ASTNode) {
	gn.mu.Lock()
	gn.nodes[node.NodeID()] = node
	if node.Type() == "FunctionDefinition" {
		gn.functions[node.NodeID()] = node
	}
	if node.Type() == "ContractDefinition" {
		gn.contracts[node.NodeID()] = node
	}
	gn.mu.Unlock()
}

func (gn *GlobalNodes) Nodes() map[int]ASTNode {
	return gn.nodes
}

func (gn *GlobalNodes) Functions() map[int]ASTNode {
	return gn.functions
}

func (gn *GlobalNodes) Contracts() map[int]ASTNode {
	return gn.contracts
}

type NormalCallPath struct {
	caller  *NormalCallPath   // caller function
	name    string            // my function name
	id      int               // my function id
	callees []*NormalCallPath // callee functions
}

func NewNormalCallPath() *NormalCallPath {
	ncp := new(NormalCallPath)
	return ncp
}

func (ncp *NormalCallPath) SetCaller(caller string, callerID int) {
	ncp.name = caller
	ncp.id = callerID
}

func (ncp *NormalCallPath) SetCallee(callee string, calleeID int) {
	if ncp.callees == nil {
		ncp.callees = make([]*NormalCallPath, 0)
	}
	ncp.callees = append(ncp.callees, &NormalCallPath{caller: ncp, name: callee, id: calleeID})
}

func (ncp *NormalCallPath) Callees() []*NormalCallPath {
	return ncp.callees
}

func (ncp *NormalCallPath) Caller() *NormalCallPath {
	return ncp.caller
}

func (ncp *NormalCallPath) ID() int {
	return ncp.id
}

func (ncp *NormalCallPath) Name() string {
	return ncp.name
}

func Lookup(ncp *NormalCallPath, id int) *NormalCallPath {
	if ncp.id == id {
		return ncp
	}

	var ret *NormalCallPath

	for _, callee := range ncp.callees {
		if Lookup(callee, id) != nil {
			ret = Lookup(callee, id)
		}
	}

	return ret
}

func Path(ncp *NormalCallPath, path *string) {
	*path = *path + ncp.name + ", "
	if len(ncp.callees) > 0 {
		for _, callee := range ncp.callees {
			Path(callee, path)
		}
	}
}
