package ast

import "testing"

func TestNCP(t *testing.T) {
	ncp := new(NormalCallPath)
	ncp.SetCaller("1", 1)
	ncp.SetCallee("2", 2)
	ncp.SetCallee("3", 3)
	ncp.SetCallee("4", 4)
	ncp.SetCallee("5", 5)
	ncp.callees[0].SetCallee("6", 6)
	ncp.callees[0].SetCallee("7", 7)
	ncp.callees[0].SetCallee("8", 8)
	ncp.callees[0].SetCallee("9", 9)
	ncp.callees[3].SetCallee("10", 10)
	ncp.callees[3].SetCallee("11", 11)
	ncp.callees[1].SetCallee("12", 12)
	ncp.callees[0].callees[1].SetCallee("13", 13)
	ncp.callees[0].callees[1].SetCallee("14", 14)
	ncp.callees[0].callees[1].SetCallee("15", 15)
	ncp.callees[0].callees[1].SetCallee("16", 16)

	t.Log(Lookup(ncp, 7))
}