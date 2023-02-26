package main

import (
	"fmt"
	"os"
	"os/exec"
)

func LaunchLotteries(network string) error {
	err := compileContracts()
	if err != nil {
		return fmt.Errorf("compiling contracts: %w", err)
	}

	err = deployContracts(network)
	if err != nil {
		return fmt.Errorf("deploying contracts: %w", err)
	}

	err = resetContracts()
	if err != nil {
		return fmt.Errorf("reseting ./contracts: %w", err)
	}
	return nil
}

func compileContracts() error {
	complileCommandOut, err := exec.Command("truffle", "compile", "--all").Output()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(complileCommandOut))

	return nil
}

func deployContracts(network string) error {
	deployCommandOut, err := exec.Command("truffle", "migrate", "--network", network).Output()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(deployCommandOut))

	return nil
}

func resetContracts() error {
	contracts, err := os.ReadDir("contracts")
	if err != nil {
		return err
	}

	for _, contract := range contracts {
		if contract.Name() == ".gitkeep" {
			continue
		}

		os.Remove(fmt.Sprintf("./contracts/%s", contract.Name()))
	}

	return nil
}
