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

//go:embed templates/1_deploy_contract.js.gotmpl
var tplLotteryDeployerTemplate string

type ContractData struct {
	DaysTilClose uint64
	ContractFull string
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

	err := Generate("lottery", tplLotteryMachineTemplate, e.Data, fmt.Sprintf("./contracts/%s.sol", e.Data.ContractName))
	if err != nil {
		return fmt.Errorf("generating lottery contract: %w", err)
	}

	err = Generate("deployer", tplLotteryDeployerTemplate, e.Data, "./migrations/1_deploy_contract.js")
	if err != nil {
		return fmt.Errorf("generating lottery deployer: %w", err)
	}
	return nil

}

func Generate(name, tpl string, data any, outputFile string) (err error) {
	var w io.Writer

	w, err = os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", outputFile, err)
	}

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

func GenerateNewLottery(contractFull, contractName, lotteryName, abbr string, tickets int, daysTil uint64) ([]byte, []byte, error) {
	lotteryData := &ContractData{
		daysTil,
		contractFull,
		contractName,
		lotteryName,
		abbr,
		tickets,
	}
	lotteryEngine := NewEngine(lotteryData)

	fmt.Println("Contract being created...")
	err := lotteryEngine.GenerateWrapper()
	if err != nil {
		return nil, nil, fmt.Errorf("generating lottery: %w", err)
	}

	err = lotteryEngine.GenerateWrapper()
	if err != nil {
		return nil, nil, fmt.Errorf("generating deployer: %w", err)
	}
	fmt.Println("Now launching lottery...")

	deployLog, abiLog, err := LaunchLotteries("sepolia", contractName)
	if err != nil {
		return nil, nil, fmt.Errorf("deploying all contracts: %w", err)
	}

	fmt.Println("Done")
	return deployLog, abiLog, nil
}
