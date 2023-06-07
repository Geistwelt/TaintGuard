package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geistwelt/logging"
	jsoniter "github.com/json-iterator/go"
)

type YulForLoop struct {
	body      ASTNode
	condition ASTNode
	NodeType  string `json:"nodeType"`
	post      ASTNode
	pre       ASTNode
	Src       string `json:"src"`
}

func (yfl *YulForLoop) SourceCode(isSc bool, isIndent bool, indent string, logger logging.Logger) string {
	var code string
	if isIndent {
		code = code + indent
	}

	code = code + "for {\n"

	// pre
	{
		if yfl.pre != nil {
			switch pre := yfl.pre.(type) {
			case *YulBlock:
				code = code + pre.SourceCode(false, false, indent, logger)
			default:
				if pre != nil {
					logger.Warnf("Unknown pre nodeType [%s] for YulForLoop [src:%s].", pre.Type(), yfl.Src)
				} else {
					logger.Warnf("Unknown pre nodeType for YulForLoop [src:%s].", yfl.Src)
				}
			}
		}
	}

	code = code + "\n"
	if isIndent {
		code = code + indent
	}
	code = code + "}"

	// condition
	{
		if yfl.condition != nil {
			switch condition := yfl.condition.(type) {
			case *YulFunctionCall:
				code = code + " " + condition.SourceCode(false, false, indent, logger)
			default:
				if condition != nil {
					logger.Warnf("Unknown condition nodeType [%s] for YulForLoop [src:%s].", condition.Type(), yfl.Src)
				} else {
					logger.Warnf("Unknown condition nodeType for YulForLoop [src:%s].", yfl.Src)
				}
			}
		}
	}

	code = code + " {\n"

	// post
	{
		if yfl.post != nil {
			switch post := yfl.post.(type) {
			case *YulBlock:
				code = code + post.SourceCode(false, false, indent, logger)
			default:
				if post != nil {
					logger.Warnf("Unknown post nodeType [%s] for YulForLoop [src:%s].", post.Type(), yfl.Src)
				} else {
					logger.Warnf("Unknown post nodeType for YulForLoop [src:%s].", yfl.Src)
				}
			}
		}
	}

	code = code + "\n"
	if isIndent {
		code = code + indent
	}
	code = code + "} {\n"

	// body
	{
		if yfl.body != nil {
			switch body := yfl.body.(type) {
			case *YulBlock:
				code = code + body.SourceCode(false, false, indent, logger)
			default:
				if body != nil {
					logger.Warnf("Unknown body nodeType [%s] for YulForLoop [src:%s].", body.Type(), yfl.Src)
				} else {
					logger.Warnf("Unknown body nodeType for YulForLoop [src:%s].", yfl.Src)
				}
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

func (yfl *YulForLoop) Type() string {
	return yfl.NodeType
}

func (yfl *YulForLoop) Nodes() []ASTNode {
	return nil
}

func (yfl *YulForLoop) NodeID() int {
	return -1
}

func GetYulForLoop(gn *GlobalNodes, raw jsoniter.Any, logger logging.Logger) (*YulForLoop, error) {
	yfl := new(YulForLoop)
	if err := json.Unmarshal([]byte(raw.ToString()), yfl); err != nil {
		logger.Errorf("Failed to unmarshal YulForLoop: [%v].", err)
		return nil, fmt.Errorf("failed to unmarshal YulForLoop: [%v]", err)
	}

	// body
	{
		body := raw.Get("body")
		if body.Size() > 0 {
			bodyNodeType := body.Get("nodeType").ToString()
			var yflBody ASTNode
			var err error

			switch bodyNodeType {
			case "YulBlock":
				yflBody, err = GetYulBlock(gn, body, logger)
			default:
				logger.Warnf("Unknown body nodeType [%s] for YulForLoop [src:%s].", bodyNodeType, yfl.Src)
			}

			if err != nil {
				return nil, err
			}

			if yflBody != nil {
				yfl.body = yflBody
			}
		}
	}

	// condition
	{
		condition := raw.Get("condition")
		if condition.Size() > 0 {
			conditionNodeType := condition.Get("nodeType").ToString()
			var yflCondition ASTNode
			var err error

			switch conditionNodeType {
			case "YulFunctionCall":
				yflCondition, err = GetYulFunctionCall(gn, condition, logger)
			default:
				logger.Warnf("Unknown condition nodeType [%s] for YulForLoop [src:%s].", conditionNodeType, yfl.Src)
			}

			if err != nil {
				return nil, err
			}

			if yflCondition != nil {
				yfl.condition = yflCondition
			}
		}
	}

	// post
	{
		post := raw.Get("post")
		if post.Size() > 0 {
			postNodeType := post.Get("nodeType").ToString()
			var yflPost ASTNode
			var err error

			switch postNodeType {
			case "YulBlock":
				yflPost, err = GetYulBlock(gn, post, logger)
			default:
				logger.Warnf("Unknown post nodeType [%s] for YulForLoop [src:%s].", postNodeType, yfl.Src)
			}

			if err != nil {
				return nil, err
			}

			if yflPost != nil {
				yfl.post = yflPost
			}
		}
	}

	// pre
	{
		pre := raw.Get("pre")
		if pre.Size() > 0 {
			preNodeType := pre.Get("nodeType").ToString()
			var yflPre ASTNode
			var err error

			switch preNodeType {
			case "YulBlock":
				yflPre, err = GetYulBlock(gn, pre, logger)
			default:
				logger.Warnf("Unknown pre nodeType [%s] for YulForLoop [src:%s].", preNodeType, yfl.Src)
			}

			if err != nil {
				return nil, err
			}

			if yflPre != nil {
				yfl.pre = yflPre
			}
		}
	}

	return yfl, nil
}
