package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/geistwelt/logging"
	"github.com/geistwelt/taintguard/global"
	"github.com/geistwelt/taintguard/src"
	v04 "github.com/geistwelt/taintguard/src/v0.4"
	v05 "github.com/geistwelt/taintguard/src/v0.5"
	v06 "github.com/geistwelt/taintguard/src/v0.6"
	v08 "github.com/geistwelt/taintguard/src/v0.8"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
)

var opt = logging.Option{
	Module:         "TaintGuard",
	FilterLevel:    logging.DebugLevel,
	Spec:           "%{color}[%{time}] [%{module}] %{location}%{color:reset} => %{message}",
	FormatSelector: "terminal",
	Writer:         os.Stdout,
}
var logger = logging.MustNewLogger(opt)

func main() {
	execute()
}

var (
	rootCmd = &cobra.Command{
		Use:   "tguard",
		Short: "Tguard is a smart contract vulnerability detection and patching tool.",
		Long: `Tguard is an automated tool that detects the existence of vulnerabilities in 
		solidity smart contracts due to delegatecall and can automatically patch contract 
		vulnerabilities using code instrumentating technology.`,
		Run: func(cmd *cobra.Command, args []string) {
			jsonBytes := src.MustReadFile(global.Input)
			if global.Output[len(global.Output)-1] != '\\' {
				global.Output = fmt.Sprintf("%s/", global.Output)
			}
			source := jsoniter.Get(jsonBytes)
			absolutePath := source.Get("absolutePath").ToString()

			var solFileName string
			absolutePaths := strings.Split(absolutePath, "/")
			solFileName = absolutePaths[len(absolutePaths)-1]

			relativePath := "contracts"
			sourceUnit := source.Get("nodeType")
			if sourceUnit.ToString() != "SourceUnit" {
				fmt.Printf("Expected SourceUnit, but got [%s].\n", sourceUnit.ToString())
				os.Exit(1)
			}
			sourceUnitNodes := source.Get("nodes")
			if sourceUnitNodes.Size() < 1 {
				fmt.Println("Invalid source file, there should be more than zero ast node in SourceUnit.")
				os.Exit(1)
			}

			var pragmaDirective jsoniter.Any

			for i := 0; i < sourceUnitNodes.Size(); i++ {
				sourceUnitNode := sourceUnitNodes.Get(i)
				if sourceUnitNode.Size() > 0 {
					sourceUnitNodeNodeType := sourceUnitNode.Get("nodeType").ToString()
					if sourceUnitNodeNodeType == "PragmaDirective" {
						pragmaDirective = sourceUnitNode
					}
				}
			}

			literals := pragmaDirective.Get("literals").ToString()
			literals_list := strings.Split(literals, ",")
			var upper, lower, version float64
			var err error
			if strings.Contains(literals, ">=") || strings.Contains(literals, ">") || strings.Contains(literals, "<=") || strings.Contains(literals, "<") {
				if strings.Contains(literals, ">=") {
					for index, word := range literals_list {
						word = src.Trim(word)
						if word == ">=" {
							lower_string := src.Trim(literals_list[index+1])
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
						word = src.Trim(word)
						if word == ">" {
							lower_string := src.Trim(literals_list[index+1])
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
						word = src.Trim(word)
						if word == "<=" {
							upper_string := src.Trim(literals_list[index+1])
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
						word = src.Trim(word)
						if word == "<" {
							upper_string := src.Trim(literals_list[index+1])
							upper_string = src.Trim(upper_string)
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

			if err = src.EnsureDir(fmt.Sprintf("%s%s", global.Output, relativePath)); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			switch version {
			case 0.7, 0.8:
				node, err := v08.Run(jsonBytes, global.Cg, logger, solFileName, global.Output, global.Variables)
				if err != nil {
					os.Exit(1)
				}
				code := node.SourceCode(false, false, "", logger)
				f, err := os.OpenFile(fmt.Sprintf("%s%s/%s", global.Output, relativePath, solFileName), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
				if err != nil {
					fmt.Printf("Failed to open %s: [%v].\n", fmt.Sprintf("%s%s/%s", global.Output, relativePath, solFileName), err)
					os.Exit(1)
				}

				f.Write([]byte(code))

				f.Close()
			case 0.6:
				node, err := v06.Run(jsonBytes, global.Cg, logger, solFileName, global.Output, global.Variables)
				if err != nil {
					os.Exit(1)
				}
				code := node.SourceCode(false, false, "", logger)
				f, err := os.OpenFile(fmt.Sprintf("%s%s/%s", global.Output, relativePath, solFileName), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
				if err != nil {
					fmt.Printf("Failed to open %s: [%v].\n", fmt.Sprintf("%s%s/%s", global.Output, relativePath, solFileName), err)
					os.Exit(1)
				}

				f.Write([]byte(code))

				f.Close()
			case 0.5:
				node, err := v05.Run(jsonBytes, global.Cg, logger, solFileName, global.Output, global.Variables)
				if err != nil {
					os.Exit(1)
				}
				code := node.SourceCode(false, false, "", logger)
				f, err := os.OpenFile(fmt.Sprintf("%s%s/%s", global.Output, relativePath, solFileName), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
				if err != nil {
					fmt.Printf("Failed to open %s: [%v].\n", fmt.Sprintf("%s%s/%s", global.Output, relativePath, solFileName), err)
					os.Exit(1)
				}

				f.Write([]byte(code))

				f.Close()
			case 0.4:
				node, err := v04.Run(jsonBytes, global.Cg, logger, solFileName, global.Output, global.Variables)
				if err != nil {
					os.Exit(1)
				}
				code := node.SourceCode(false, false, "", logger)
				f, err := os.OpenFile(fmt.Sprintf("%s%s/%s", global.Output, relativePath, solFileName), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
				if err != nil {
					fmt.Printf("Failed to open %s: [%v].\n", fmt.Sprintf("%s%s/%s", global.Output, relativePath, solFileName), err)
					os.Exit(1)
				}

				f.Write([]byte(code))

				f.Close()
			default:
				fmt.Println(version)
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringSliceVar(&global.Variables, "variables", []string{"owner", "_owner", "owner_"}, "Variables to store permission information, default is [owner]")
	rootCmd.PersistentFlags().StringVar(&global.Input, "input", "contracts/v0.8/1.sol_json.ast", "Path to the abstract syntax tree file of the smart contract to be analyzed")
	rootCmd.PersistentFlags().StringVar(&global.Output, "output", "test", "The path to the folder where the analysis results are stored.")
	rootCmd.PersistentFlags().BoolVar(&global.Cg, "call-graph", false, "Whether to generate a function call relationship graph within the contract, default is false.")
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
