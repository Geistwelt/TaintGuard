package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type TryCatchClause struct {
	block      ASTNode
	ErrorName  string `json:"errorName"`
	ID         int    `json:"id"`
	NodeType   string `json:"nodeType"`
	parameters ASTNode
	Src        string `json:"src"`
}

func (tcc *TryCatchClause) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string

	if tcc.parameters != nil {
		switch parameters := tcc.parameters.(type) {
		case *ParameterList:
			code = code + "(" + parameters.SourceCode(false, false, indent, logger) + ")"
		default:
			if parameters != nil {
				logger.Warnf("Unknown parameters nodeType [%s] for TryCatchClause [src:%s].", parameters.Type(), tcc.Src)
			} else {
				logger.Warnf("Unknown parameters nodeType for TryCatchClause [src:%s].", tcc.Src)
			}
		}
	}

	code = code + " {\n"

	if tcc.block != nil {
		switch block := tcc.block.(type) {
		case *Block:
			code = code + block.SourceCode(false, false, indent, logger)
		default:
			if block != nil {
				logger.Warnf("Unknown block nodeType [%s] for TryCatchClause [src:%s].", block.Type(), tcc.Src)
			} else {
				logger.Warnf("Unknown block nodeType for TryCatchClause [src:%s].", tcc.Src)
			}
		}
	}
	code = code + "\n"
	if isIndent {
		code = code + indent
	}
	code = code + "}"
	return code
}

func (tcc *TryCatchClause) Type() string {
	return tcc.NodeType
}

func (tcc *TryCatchClause) Nodes() []ASTNode {
	return nil
}

func (tcc *TryCatchClause) NodeID() int {
	return tcc.ID
}

func GetTryCatchClause(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*TryCatchClause, error) {
	tcc := new(TryCatchClause)
	if err := json.Unmarshal([]byte(raw.ToString()), tcc); err != nil {
		logger.Errorf("Failed to unmarshal TryCatchClause: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal TryCatchClause: [%v]", err)
	}

	// body
	{
		block := raw.Get("block")
		if block.Size() > 0 {
			blockNodeType := block.Get("nodeType").ToString()
			var tccBlock ASTNode
			var err error

			switch blockNodeType {
			case "Block":
				tccBlock, err = GetBlock(gn, block, logger)
			default:
				logger.Warnf("Unknown block nodeType [%s] for TryCatchClause [src:%s].", blockNodeType, tcc.Src)
			}

			if err != nil {
				return nil, err
			}

			if tccBlock != nil {
				tcc.block = tccBlock
			}
		}
	}

	// parameters
	{
		parameters := raw.Get("parameters")
		if parameters.Size() > 0 {
			parametersNodeType := parameters.Get("nodeType").ToString()
			var tccParameters ASTNode
			var err error

			switch parametersNodeType {
			case "ParameterList":
				tccParameters, err = GetParameterList(gn, parameters, logger)
			default:
				logger.Warnf("Unknown parameters nodeType [%s] for TryCatchClause [src:%s].", parametersNodeType, tcc.Src)
			}

			if err != nil {
				return nil, err
			}

			if tccParameters != nil {
				tcc.parameters = tccParameters
			}
		}
	}

	gn.AddASTNode(tcc)

	return tcc, nil
}

func (tcc *TryCatchClause) TraverseFunctionCall(ncp *NormalCallPath, gn *GlobalNodes, opt *Option, logger logging.Logger) {
	if tcc.block != nil {
		switch block := tcc.block.(type) {
		case *Block:
			block.TraverseFunctionCall(ncp, gn, opt, logger)
		}
	}
}

func (tcc *TryCatchClause) TraverseTaintOwner(opt *Option, logger logging.Logger) {
	if tcc.block != nil {
		switch block := tcc.block.(type) {
		case *Block:
			block.TraverseTaintOwner(opt, logger)
		}
	}
}

func (tcc *TryCatchClause) TraverseDelegatecall(opt *Option, logger logging.Logger) {
	if tcc.block != nil {
		switch block := tcc.block.(type) {
		case *Block:
			block.TraverseDelegatecall(opt, logger)
		}
	}
}