package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/geistwelt/logging"
	"github.com/geistwelt/taintguard/src"
	v05 "github.com/geistwelt/taintguard/src/v0.5"
	v08 "github.com/geistwelt/taintguard/src/v0.8"
	jsoniter "github.com/json-iterator/go"
)

var contract_ast_json string = "contracts/v0.6/1.sol_json.ast"
var opt = logging.Option{
	Module:         "TaintGuard",
	FilterLevel:    logging.DebugLevel,
	Spec:           "%{color}[%{time}] [%{module}] %{location}%{color:reset} => %{message}",
	FormatSelector: "terminal",
	Writer:         os.Stdout,
}
var logger = logging.MustNewLogger(opt)

func main() {

	jsonBytes := src.MustReadFile(contract_ast_json)

	contract := jsoniter.Get(jsonBytes)
	sourceUnit := contract.Get("nodeType")
	if sourceUnit.ToString() != "SourceUnit" {
		fmt.Printf("Expected SourceUnit, but got [%s].\n", sourceUnit.ToString())
		os.Exit(1)
	}
	sourceUnitNodes := contract.Get("nodes")
	if sourceUnitNodes.Size() < 1 {
		fmt.Println("Invalid source file, there should be more than zero ast node in SourceUnit.")
		os.Exit(1)
	}
	// pragmaDirective := sourceUnitNodes.Get(0)
	var pragmaDirective jsoniter.Any
	// if pragmaDirective.Get("nodeType").ToString() != "PragmaDirective" {
	// 	fmt.Printf("Expected PragmaDirective, but got [%s].\n", pragmaDirective.ToString())
	// 	os.Exit(1)
	// }

	for i := 0; i < sourceUnitNodes.Size(); i++ {
		sourceUnitNode := sourceUnitNodes.Get(i)
		if sourceUnitNode.Size() > 0 {
			sourceUnitNodeNodeType := sourceUnitNode.Get("nodeType").ToString()
			if sourceUnitNodeNodeType == "PragmaDirective" {
				pragmaDirective = sourceUnitNode
			}
		}
	}
	// var pragma string = "pragma"
	literals := pragmaDirective.Get("literals").ToString()
	literals_list := strings.Split(literals, ",")
	var upper, lower, version float64
	var err error
	if strings.Contains(literals, ">=") || strings.Contains(literals, ">") || strings.Contains(literals, "<=") || strings.Contains(literals, "<") {
		if strings.Contains(literals, ">=") {
			for index, word := range literals_list {
				word = strings.Trim(word, "\"")
				if word == ">=" {
					lower_string := strings.Trim(literals_list[index+1], "\"")
					lower, err = strconv.ParseFloat(lower_string, 64)
					if err != nil {
						fmt.Printf("1: Failed parse solidity version: [%v], [%s].\n", err, lower_string)
						os.Exit(1)
					}
					upper = lower
				}
			}
		}
		if strings.Contains(literals, ">") {
			for index, word := range literals_list {
				word = strings.Trim(word, "\"")
				if word == ">" {
					lower_string := strings.Trim(literals_list[index+1], "\"")
					lower, err = strconv.ParseFloat(lower_string, 64)
					if err != nil {
						fmt.Printf("2: Failed parse solidity version: [%v], [%s].\n", err, lower_string)
						os.Exit(1)
					}
					lower += 0.1
					upper = lower
				}
			}
		}
		if strings.Contains(literals, "<=") {
			for index, word := range literals_list {
				word = strings.Trim(word, "\"")
				if word == "<=" {
					upper_string := strings.Trim(literals_list[index+1], "\"")
					upper, err = strconv.ParseFloat(upper_string, 64)
					if err != nil {
						fmt.Printf("3: Failed parse solidity version: [%v], [%s].\n", err, upper_string)
						os.Exit(1)
					}
					if lower == 0.0 {
						lower = upper
					}
				}
			}
		}
		if strings.Contains(literals, "<") {
			for index, word := range literals_list {
				word = strings.Trim(word, "\"")
				if word == "<" {
					upper_string := strings.Trim(literals_list[index+1], "\"")
					upper_string = strings.Trim(upper_string, " ")
					upper, err = strconv.ParseFloat(upper_string, 64)
					if err != nil {
						fmt.Printf("4: Failed parse solidity version: [%v], [%s].\n", err, upper_string)
						os.Exit(1)
					}
					upper -= 0.1
					if lower == 0.0 {
						lower = upper
					}
				}
			}
		}
	} else {
		reNum := regexp.MustCompile(`0\.\d*`)
		for _, word := range literals_list {
			if reNum.Match([]byte(word)) {
				word = src.Trim(word)
				version, err = strconv.ParseFloat(word, 64)
				if err != nil {
					fmt.Printf("5: Failed parse solidity version: [%v], [%s].\n", err, word)
					os.Exit(1)
				}
			}
		}
	}

	if version == 0.0 {
		if upper < lower {
			fmt.Printf("6: Failed parse solidity version: [upper(%.1f) < lower(%.1f)].\n", upper, lower)
			os.Exit(1)
		}
		version = lower
	}

	switch version {
	case 0.8:
		node, err := v08.Run(jsonBytes, false, logger)
		if err != nil {
			os.Exit(1)
		}
		code := node.SourceCode(false, false, "", logger)
		f, err := os.OpenFile("test/0.sol", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
		if err != nil {
			fmt.Println("Failed to open file test/0.sol.")
			os.Exit(1)
		}

		f.Write([]byte(code))

		f.Close()
	case 0.6:
		fmt.Println("Success")
	case 0.5:
		node, err := v05.Run(jsonBytes, false, logger)
		if err != nil {
			os.Exit(1)
		}
		code := node.SourceCode(false, false, "", logger)
		f, err := os.OpenFile("test/0.sol", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
		if err != nil {
			fmt.Println("Failed to open file test/0.sol.")
			os.Exit(1)
		}

		f.Write([]byte(code))

		f.Close()
	default:
		fmt.Println(version)
	}
}
