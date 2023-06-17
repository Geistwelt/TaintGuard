package cfg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/geistwelt/logging"
	"github.com/geistwelt/taintguard/src"
	"github.com/geistwelt/taintguard/src/v0.6/ast"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

// make call graph
func MakeCG(ncp *ast.NormalCallPath, logger logging.Logger, solfileName string, dirName string) error {
	if err := src.EnsureDir(fmt.Sprintf("%s/%s/%s", dirName, "call-graph", solfileName)); err != nil {
		logger.Error(err.Error())
	}
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		logger.Errorf("Failed to make graph: [%v].", err)
		return fmt.Errorf("failed to make graph: [%v]", err)
	}
	if err = makeCG(graph, ncp, logger); err != nil {
		logger.Errorf("Failed to make call graph for [%s]: [%v].", ncp.Name(), err)
		return fmt.Errorf("failed to make call graph for [%s]: [%v]", ncp.Name(), err)
	}

	f, err := os.OpenFile(fmt.Sprintf("%s/%d.gv", fmt.Sprintf("%s/%s/%s", dirName, "call-graph", solfileName), ncp.ID()), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		logger.Errorf("Failed to open file [%s]: [%v].", fmt.Sprintf("%s/%d.gv", fmt.Sprintf("%s/%s/%s", dirName, "call-graph", solfileName), ncp.ID()), err)
		return fmt.Errorf("failed to open file [%s]: [%v]", fmt.Sprintf("%s/%d.gv", fmt.Sprintf("%s/%s/%s", dirName, "call-graph", solfileName), ncp.ID()), err)
	}

	var buf bytes.Buffer
	if err = g.Render(graph, "dot", &buf); err != nil {
		logger.Errorf("Failed to render call graph: [%v].", err)
		return fmt.Errorf("failed to render call graph: [%v]", err)
	}

	if _, err = f.Write(buf.Bytes()); err != nil {
		logger.Errorf("Failed to save dot file: [%v].", err)
		return fmt.Errorf("failed to save dot file: [%v]", err)
	}

	cmd := exec.Command(`/bin/sh`, `-c`, fmt.Sprintf("dot %s/%d.gv -T png -o %s/%d.png", fmt.Sprintf("%s/%s/%s", dirName, "call-graph", solfileName), ncp.ID(), fmt.Sprintf("%s/%s/%s", dirName, "call-graph", solfileName), ncp.ID()))
	var out bytes.Buffer
	cmd.Stdout = &out
	if err = cmd.Run(); err != nil {
		logger.Errorf("Failed to generate file %d.png.", ncp.ID())
		return fmt.Errorf("failed to generate file %d.png", ncp.ID())
	}

	logger.Infof("Successfully generate call graph for function [%s] => [%d.png | %d.gv].", ncp.Name(), ncp.ID(), ncp.ID())

	return nil
}

func makeCG(graph *cgraph.Graph, ncp *ast.NormalCallPath, logger logging.Logger) error {
	caller, err := graph.CreateNode(ncp.Name())
	if err != nil {
		logger.Errorf("Adding caller node failed: [%v].", err)
		return fmt.Errorf("adding caller node failed: [%v]", err)
	}

	var e error
	for _, callee := range ncp.Callees() {
		calleeNode, err := graph.CreateNode(callee.Name())
		if err != nil {
			logger.Errorf("Adding callee node failed: [%v].", err)
			e = fmt.Errorf("adding callee node failed: [%v]", err)
		}
		callEdge, err := graph.CreateEdge("call", caller, calleeNode)
		if err != nil {
			logger.Errorf("Adding call edge failed: [%v].", err)
			e = fmt.Errorf("adding call edge failed: [%v]", err)
		}
		callEdge.SetLabel(" call")
		if err = makeCG(graph, callee, logger); err != nil {
			logger.Errorf("Handling callee node failed: [%v].", err)
			e = fmt.Errorf("handling callee node failed: [%v]", err)
		}
	}
	return e
}
