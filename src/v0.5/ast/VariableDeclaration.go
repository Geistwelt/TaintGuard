package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type VariableDeclaration struct {
	Constant         bool   `json:"constant"`
	ID               int    `json:"id"`
	Indexed          bool   `json:"indexed"`
	Name             string `json:"name"`
	NodeType         string `json:"nodeType"`
	Scope            int    `json:"scope"`
	Src              string `json:"src"`
	StateVariable    bool   `json:"stateVariable"`
	StorageLocation  string `json:"storageLocation"`
	TypeDescriptions struct {
		TypeIdentifier string `json:"typeIdentifier"`
		TypeString     string `json:"typeString"`
	} `json:"typeDescriptions"`
	typeName   ASTNode
	value      ASTNode
	Visibility string `json:"visibility"`
}

func (vd *VariableDeclaration) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if isIndent {
		code = code + indent
	}

	if vd.typeName != nil {
		switch typeName := vd.typeName.(type) {
		case *Mapping:
			code = code + typeName.SourceCode(false, false, indent, logger)
		case *ElementaryTypeName:
			code = code + typeName.SourceCode(false, false, indent, logger)
		case *UserDefinedTypeName:
			code = code + typeName.SourceCode(false, false, indent, logger)
		case *ArrayTypeName:
			code = code + typeName.SourceCode(false, false, indent, logger)
		default:
			if typeName != nil {
				logger.Warnf("Unknown typeName nodeType [%s] for VariableDeclaration [src:%s].", typeName.Type(), vd.Src)
			} else {
				logger.Warnf("Unknown typeName nodeType for VariableDeclaration [src:%s].", vd.Src)
			}
		}
	}

	if vd.Constant {
		code = code + " " + "constant"
	}

	if vd.StorageLocation != "default" {
		code = code + " " + vd.StorageLocation
	}

	if vd.Visibility != "" && vd.Visibility != "internal" {
		code = code + " " + vd.Visibility
	}

	if vd.Indexed {
		code = code + " " + "indexed"
	}

	if vd.Name != "" {
		code = code + " " + vd.Name
	}

	if vd.value != nil {
		switch value := vd.value.(type) {
		case *Literal:
			code = code + " = " + value.SourceCode(false, false, indent, logger)
		case *BinaryOperation:
			code = code + " = " + value.SourceCode(false, false, indent, logger)
		case *FunctionCall:
			code = code + " = " + value.SourceCode(false, false, indent, logger)
		case *Identifier:
			code = code + " = " + value.SourceCode(false, false, indent, logger)
		case *MemberAccess:
			code = code + " = " + value.SourceCode(false, false, indent, logger)
		case *UnaryOperation:
			code = code + " = " + value.SourceCode(false, false, indent, logger)
		case *TupleExpression:
			code = code + " = " + value.SourceCode(false, false, indent, logger)
		case *ArrayTypeName:
			code = code + " = " + value.SourceCode(false, false, indent, logger)
		default:
			if value != nil {
				logger.Warnf("Unknown value nodeType [%s] for VariableDeclaration [src:%s].", value.Type(), vd.Src)
			} else {
				logger.Warnf("Unknown value nodeType for VariableDeclaration [src:%s].", vd.Src)
			}
		}
	}

	if isSc {
		code = code + ";"
	}

	return code
}

func (vd *VariableDeclaration) Type() string {
	return vd.NodeType
}

func (vd *VariableDeclaration) Nodes() []ASTNode {
	return nil
}

func (vd *VariableDeclaration) NodeID() int {
	return vd.ID
}

func GetVariableDeclaration(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*VariableDeclaration, error) {
	vd := new(VariableDeclaration)
	if err := json.Unmarshal([]byte(raw.ToString()), vd); err != nil {
		logger.Errorf("Failed to unmarshal VariableDeclaration: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal VariableDeclaration: [%v]", err)
	}

	// typeName
	{
		typeName := raw.Get("typeName")
		typeNameNodeType := typeName.Get("nodeType").ToString()

		var vdTypeName ASTNode
		var err error

		switch typeNameNodeType {
		case "Mapping":
			vdTypeName, err = GetMapping(gn, typeName, logger)
		case "ElementaryTypeName":
			vdTypeName, err = GetElementaryTypeName(gn, typeName, logger)
		case "UserDefinedTypeName":
			vdTypeName, err = GetUserDefinedTypeName(gn, typeName, logger)
		case "ArrayTypeName":
			vdTypeName, err = GetArrayTypeName(gn, typeName, logger)
		default:
			logger.Warnf("Unknown typeName nodeType [%s] for VariableDeclaration [src:%s].", typeNameNodeType, vd.Src)
		}

		if err != nil {
			return nil, err
		}
		vd.typeName = vdTypeName
	}

	// value
	{
		value := raw.Get("value")
		if value.Size() > 0 {
			valueNodeType := value.Get("nodeType").ToString()

			var vdValue ASTNode
			var err error

			switch valueNodeType {
			case "Literal":
				vdValue, err = GetLiteral(gn, value, logger)
			case "BinaryOperation":
				vdValue, err = GetBinaryOperation(gn, value, logger)
			case "FunctionCall":
				vdValue, err = GetFunctionCall(gn, value, logger)
			case "Identifier":
				vdValue, err = GetIdentifier(gn, value, logger)
			case "MemberAccess":
				vdValue, err = GetMemberAccess(gn, value, logger)
			case "UnaryOperation":
				vdValue, err = GetUnaryOperation(gn, value, logger)
			case "TupleExpression":
				vdValue, err = GetTupleExpression(gn, value, logger)
			default:
				logger.Warnf("Unknown value nodeType [%s] for VariableDeclaration [src:%s]", valueNodeType, vd.Src)
			}

			if err != nil {
				return nil, err
			}
			if vdValue != nil {
				vd.value = vdValue
			}
		}

	}

	gn.AddASTNode(vd)

	return vd, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (vd *VariableDeclaration) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	if vd.typeName != nil {
		switch typeName := vd.typeName.(type) {
		case *Mapping:
			typeName.TraverseFunctionCall(ncp, gn, opt, logger)
		}
	}

	if vd.value != nil {
		switch value := vd.value.(type) {
		case *BinaryOperation:
			value.TraverseFunctionCall(ncp, gn, opt, logger)
		case *FunctionCall:
			value.TraverseFunctionCall(ncp, gn, opt, logger)
		}
	}
}

func (vd *VariableDeclaration) SetTypeName(typeName ASTNode) {
	vd.typeName = typeName
}
