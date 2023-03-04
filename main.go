package main

import (
	_ "embed"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting Lottery Machine...")

	http.HandleFunc("/gen", GenHandler)
	http.ListenAndServe(":8080", nil)
}
