package ast

import "github.com/geistwelt/logging"

type ASTNode interface {
	Type() string
	SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string
	Nodes() []ASTNode
	NodeID() int
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Option struct {
	// search object that delegatecall unknown contract
	delegatecallUnknownContractCh chan struct{}

	// search object that delegatecall known contract
	delegatecallKnownContractCh chan string

	// instrument track code
	TrackFunctionDefinitionName string
	TrackOwnerVariableName      string
	TrackOwnerMappingName       string
	SimilarOwnerVariableName    string
	IsTainted                   bool
	TrackAssignment             ASTNode

	// search sentence that use delegatecall
	ExpressionStatement ASTNode
}

func (opt *Option) MakeDelegatecallUnknownContractCh(size int) {
	opt.delegatecallUnknownContractCh = make(chan struct{}, size)
}

func (opt *Option) MakeDelegatecallKnownContractCh(size int) {
	opt.delegatecallKnownContractCh = make(chan string, size)
}

func (opt *Option) DelegatecallUnknownContractCh() <-chan struct{} {
	return opt.delegatecallUnknownContractCh
}

func (opt *Option) DelegatecallKnownContractCh() <-chan string {
	return opt.delegatecallKnownContractCh
}

type traverseFunctionCall interface {
	TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger)
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

type traverseTaintOwner interface {
	TraverseTaintOwner(opt *Option, logger logging.Logger)
}

var _ traverseTaintOwner = (*ContractDefinition)(nil)
var _ traverseTaintOwner = (*FunctionDefinition)(nil)
var _ traverseTaintOwner = (*Block)(nil)
var _ traverseTaintOwner = (*ExpressionStatement)(nil)
var _ traverseTaintOwner = (*IfStatement)(nil)
var _ traverseTaintOwner = (*ForStatement)(nil)
var _ traverseTaintOwner = (*UncheckedBlock)(nil)
var _ traverseTaintOwner = (*WhileStatement)(nil)
var _ traverseTaintOwner = (*TryStatement)(nil)
var _ traverseTaintOwner = (*DoWhileStatement)(nil)
var _ traverseTaintOwner = (*Assignment)(nil)
var _ traverseTaintOwner = (*TryCatchClause)(nil)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type traverseDelegatecall interface {
	TraverseDelegatecall(opt *Option, logger logging.Logger)
}

var _ traverseDelegatecall = (*ContractDefinition)(nil)
var _ traverseDelegatecall = (*FunctionDefinition)(nil)
var _ traverseDelegatecall = (*Block)(nil)
var _ traverseDelegatecall = (*ExpressionStatement)(nil)
var _ traverseDelegatecall = (*IfStatement)(nil)
var _ traverseDelegatecall = (*ForStatement)(nil)
var _ traverseDelegatecall = (*UncheckedBlock)(nil)
var _ traverseDelegatecall = (*WhileStatement)(nil)
var _ traverseDelegatecall = (*TryStatement)(nil)
var _ traverseDelegatecall = (*DoWhileStatement)(nil)
var _ traverseDelegatecall = (*VariableDeclarationStatement)(nil)
var _ traverseDelegatecall = (*FunctionCall)(nil)
var _ traverseDelegatecall = (*FunctionCallOptions)(nil)
var _ traverseDelegatecall = (*MemberAccess)(nil)
var _ traverseDelegatecall = (*BinaryOperation)(nil)
var _ traverseDelegatecall = (*TryCatchClause)(nil)
var _ traverseDelegatecall = (*Assignment)(nil)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Analysis

type GlobalNodes struct {
	nodes           map[int]ASTNode    // id => all ASTNode
	contractsByID       map[int]ASTNode    // id => all ContractDefinition
	contractsByName map[string]ASTNode // name => all ContractDefinition
	functions       map[int]ASTNode    // id => all FunctionDefinition
	mu              sync.RWMutex
}

func NewGlobalNodes() *GlobalNodes {
	gn := new(GlobalNodes)
	gn.nodes = make(map[int]ASTNode)
	gn.contractsByID = make(map[int]ASTNode)
	gn.contractsByName = make(map[string]ASTNode)
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
		gn.contractsByID[node.NodeID()] = node
		cd := node.(*ContractDefinition)
		gn.contractsByName[cd.Name] = node
	}
	gn.mu.Unlock()
}

func (gn *GlobalNodes) Nodes() map[int]ASTNode {
	return gn.nodes
}

func (gn *GlobalNodes) Functions() map[int]ASTNode {
	return gn.functions
}

func (gn *GlobalNodes) ContractsByID() map[int]ASTNode {
	return gn.contractsByID
}

func (gn *GlobalNodes) ContractsByName() map[string]ASTNode {
	return gn.contractsByName
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

