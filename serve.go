package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Contract string `json:"contract"`
	Abi      string `json:"abi"`
}

func GenHandler(w http.ResponseWriter, r *http.Request) {
	var contract string

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// parse the request body
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Failed to parse request: %v", err)
		return
	}

	//retrieve arguments from the request body
	lotteryName := r.FormValue("arg1")
	lotteryAbbr := r.FormValue("arg2")
	numOfTickets := r.FormValue("arg3")
	daysTilDraw := r.FormValue("arg4")
	daysTilInt, _ := strconv.Atoi(daysTilDraw)

	fmt.Printf("\n\ninitializing creating lottery: %s, amount of tickets: %s, drawing in %s days\n", lotteryName, numOfTickets, drawDate)
	contractLog, abiLog, err := GenerateNewLottery(fmt.Sprintf("%sContract.sol", lotteryName), fmt.Sprintf("%sContract", lotteryName), lotteryName, lotteryAbbr, 1000, uint64(daysTilInt))
	if err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(http.StatusOK)

	if string(contractLog) == "" {
		contract = "Deploy Issue: No Log"
	} else {
		contract = extractLogContract(contractLog)
	}

	response := Response{Contract: contract, Abi: string(abiLog)}
	json.NewEncoder(w).Encode(response)
}

func extractLogContract(logBytes []byte) string {
	log := string(logBytes)
	splitLeft := strings.Split(log, "contract address:")
	splitRight := strings.Split(splitLeft[1], "> block number:")

	contract := strings.ReplaceAll(splitRight[0], " ", "")
	contract = strings.ReplaceAll(contract, "\n", "")

	return contract
}
