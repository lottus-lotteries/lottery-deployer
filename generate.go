package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"text/template"
)

//go:embed templates/lotteryMachine.sol.gotmpl
var tplLotteryMachineTemplate string

type ContractData struct {
	ContractName string
	LotteryName  string
	Abbr         string
	Tickets      int
}

type Engine struct {
	Data *ContractData
}

func (e *Engine) GetEngine() *Engine {
	return e
}

func NewEngine(data *ContractData) *Engine {
	engine := &Engine{
		Data: data,
	}

	return engine
}

func (e *Engine) GenerateWrapper() error {

	err := Generate("Lottery", tplLotteryMachineTemplate, e.Data, fmt.Sprintf("./contracts/%s.sol", e.Data.ContractName))
	if err != nil {
		return fmt.Errorf("generating lottery contract: %w", err)
	}
	return nil

}

func Generate(name, tpl string, data any, outputFile string) (err error) {
	var w io.Writer

	fmt.Printf("%s\n", "Creating Lottery Machine Smart Contract...")
	w, err = os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", outputFile, err)
	}

	fmt.Printf("%s\n", "Parsing Lottery Information...")

	tmpl, err := template.New(name).Parse(tpl)
	if err != nil {
		return fmt.Errorf("parsing lottery info: %w", err)
	}

	err = tmpl.Execute(
		w,
		data,
	)
	if err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	return nil
}

func GenerateNewLottery(contractName, lotteryName, abbr string, tickets int) error {
	lotteryData := &ContractData{
		contractName,
		lotteryName,
		abbr,
		tickets,
	}
	lotteryEngine := NewEngine(lotteryData)

	err := lotteryEngine.GenerateWrapper()
	if err != nil {
		return fmt.Errorf("generating lottery: %w", err)
	}

	fmt.Println("\n----------CONTRACT CREATED----------")

	err = LaunchLotteries("goerli")
	if err != nil {
		return fmt.Errorf("deploying all contracts: %w", err)
	}

	return nil
}
