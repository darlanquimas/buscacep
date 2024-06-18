package main

import (
	"BuscaCEPViaCEP/internal/cep"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <CEP1> <CEP2> ... <CEPn>")
		os.Exit(1)
	}

	db, err := cep.InitDB("./viacep.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao inicializar banco de dados: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	fetcher := &cep.ViaCepFetcher{}
	cep.ProcessCeps(db, fetcher, os.Args[1:])
}
