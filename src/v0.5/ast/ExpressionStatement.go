package ast

var instrumentedTrack map[string]bool = make(map[string]bool)

type ExpressionStatement struct {
	expression ASTNode
	ID         int    `json:"id"`
	NodeType   string `json:"nodeType"`
	Src        string `json:"src"`

	trackMapping  ASTNode
	trackVariable ASTNode
}

func (es *ExpressionStatement) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	// expression
	if es.expression != nil {
		switch expression := es.expression.(type) {
		case *Assignment:
			code = code + expression.SourceCode(false, false, indent, logger)
		case *FunctionCall:
			code = code + expression.SourceCode(false, false, indent, logger)
		case *UnaryOperation:
			code = code + expression.SourceCode(false, false, indent, logger)
		case *Identifier:
			code = code + expression.SourceCode(false, false, indent, logger)
		default:
			if expression != nil {
				logger.Warnf("Unknown expression nodeType [%s] for ExpressionStatement [src:%s].", expression.Type(), es.Src)
			} else {
				logger.Warnf("Unknown expression nodeType for ExpressionStatement [src:%s].", es.Src)
			}
		}

		if es.trackMapping != nil && !instrumentedTrack[es.trackMapping.SourceCode(false, true, indent, logger)] {
			code = code + ";"
		}
	}

	if es.trackVariable != nil {
		switch other := es.trackVariable.(type) {
		case *ExpressionStatement:
			if !instrumentedTrack[other.SourceCode(true, true, indent, logger)] {
				code = code + "\n" + other.SourceCode(true, true, indent, logger)
				instrumentedTrack[other.SourceCode(true, true, indent, logger)] = true
			}
		}
	}

	if es.trackMapping != nil {
		switch other := es.trackMapping.(type) {
		case *ExpressionStatement:
			if !instrumentedTrack[other.SourceCode(false, true, indent, logger)] {
				code = code + "\n" + other.SourceCode(false, true, indent, logger)
				instrumentedTrack[other.SourceCode(false, true, indent, logger)] = true
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (es *ExpressionStatement) Type() string {
	return es.NodeType
}

func (es *ExpressionStatement) Nodes() []ASTNode {
	return nil
}

func (es *ExpressionStatement) NodeID() int {
	return es.ID
}

func GetExpressionStatement(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*ExpressionStatement, error) {
	es := new(ExpressionStatement)
	if err := json.Unmarshal([]byte(raw.ToString()), es); err != nil {
		logger.Errorf("Failed to unmarshal ExpressionStatement: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal ExpressionStatement: [%v]", err)
	}

	// expression
	{
		expression := raw.Get("expression")
		if expression.Size() > 0 {
			var expressionNodeType = expression.Get("nodeType").ToString()
			var esExpression ASTNode
			var err error

			switch expressionNodeType {
			case "Assignment":
				esExpression, err = GetAssignment(gn, expression, logger)
			case "FunctionCall":
				esExpression, err = GetFunctionCall(gn, expression, logger)
			case "UnaryOperation":
				esExpression, err = GetUnaryOperation(gn, expression, logger)
			case "Identifier":
				esExpression, err = GetIdentifier(gn, expression, logger)
			default:
				logger.Warnf("Unknown expression nodeType [%s] for ExpressionStatement [src:%s].", expressionNodeType, es.Src)
			}

			if err != nil {
				return nil, err
			}

			if esExpression != nil {
				es.expression = esExpression
			}
		}
	}

	gn.AddASTNode(es)

	return es, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (es *ExpressionStatement) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	if es.expression != nil {
		switch expression := es.expression.(type) {
		case *Assignment:
			expression.TraverseFunctionCall(ncp, gn, opt, logger)
		case *FunctionCall:
			expression.TraverseFunctionCall(ncp, gn, opt, logger)
		case *UnaryOperation:
			expression.TraverseFunctionCall(ncp, gn, opt, logger)
		}
	}
}

func (es *ExpressionStatement) TraverseTaintOwner(opt *Option, logger logging.Logger) {
	if es.expression != nil {
		switch expression := es.expression.(type) {
		case *Assignment:
			expression.TraverseTaintOwner(opt, logger)
			if opt.IsTainted {
				trackExpressionStatement := &ExpressionStatement{
					expression: opt.TrackAssignment,
					NodeType:   "ExpressionStatement",
					Src:        "xxx",
				}
				es.trackMapping = trackExpressionStatement
				es.trackVariable = &ExpressionStatement{
					expression: &Assignment{
						LValueRequested: false,
						leftHandSide: &Identifier{
							Name:     opt.TrackOwnerVariableName,
							NodeType: "Identifier",
							Src:      "xxx",
						},
						NodeType: "Assignment",
						Operator: "=",
						rightHandSide: &Literal{
							Kind:     "string",
							NodeType: "Literal",
							Src:      "xxx",
							Value:    opt.TrackFunctionDefinitionName,
						},
						Src: "xxx",
						TypeDescriptions: struct {
							TypeIdentifier string "json:\"typeIdentifier\""
							TypeString     string "json:\"typeString\""
						}{},
					},
					NodeType: "ExpressionStatement",
					Src:      "xxx",
				}
			}
		}
	}
}

func (es *ExpressionStatement) TraverseDelegatecall(opt *Option, logger logging.Logger) {
	if es.expression != nil {
		switch expression := es.expression.(type) {
		case *Assignment:
			expression.TraverseDelegatecall(opt, logger)
		}
	}
}

func (es *ExpressionStatement) SetExpression(expression ASTNode) {
	es.expression = expression
}
