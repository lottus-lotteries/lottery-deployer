package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Contract struct {
	ContractName string          `json:"contractName"`
	Abi          json.RawMessage `json:"abi"`
}

func LaunchLotteries(network, contractName string) ([]byte, []byte, error) {
	os.Setenv("CONTRACT_NAME", contractName)

	fmt.Println("Compiling contract...")
	err := compileContracts()
	if err != nil {
		return nil, nil, fmt.Errorf("compiling contracts: %w", err)
	}

	fmt.Println("Deploying contract...")
	deployLog, err := deployContracts(network)
	if err != nil {
		return nil, nil, fmt.Errorf("deploying contracts: %w", err)
	}

	abiLog, err := getAbi(contractName)
	if err != nil {
		return nil, nil, fmt.Errorf("getting abi logs: %w", err)
	}

	fmt.Println("Cleaning up...")
	err = resetContract(contractName)
	if err != nil {
		return nil, nil, fmt.Errorf("reseting ./contracts: %w", err)
	}

	return deployLog, abiLog, nil
}

func compileContracts() error {
	_, err := exec.Command("truffle", "compile", "--all").Output()
	if err != nil {
		return err
	}

	return nil
}

func deployContracts(network string) ([]byte, error) {
	deployCommandOut, err := exec.Command("truffle", "migrate", "--network", network).Output()
	if err != nil {
		return []byte("Deploying Issue"), err
	}

	return deployCommandOut, nil
}

func getAbi(contractName string) ([]byte, error) {
	abiLogBytes, err := os.ReadFile(fmt.Sprintf("build/contracts/%s.json", contractName))
	if err != nil {
		return nil, fmt.Errorf("opening json file: %w", err)
	}

	var contract Contract
	err = json.Unmarshal(abiLogBytes, &contract)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling: %w", err)
	}

	return contract.Abi, nil
}

func resetContract(name string) error {

	// Remove contract
	err := os.Remove(fmt.Sprintf("contracts/%s.sol", name))
	if err != nil {
		return fmt.Errorf("clearing contract: %w", err)
	}

	// Remove abi
	err = os.Remove(fmt.Sprintf("build/contracts/%s.json", name))
	if err != nil {
		return fmt.Errorf("clearing abi: %w", err)
	}

	// Remove deployer
	err = os.Remove("migrations/1_deploy_contract.js")
	if err != nil {
		return fmt.Errorf("clearing deployer: %w", err)
	}
	return nil
}
