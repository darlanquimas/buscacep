package cep

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

// Struct que define a estrutura dos dados do ViaCep
type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

// Interface que define o método FetchCep
type CepFetcher interface {
	FetchCep(cep string) (*ViaCep, error)
}

// Struct que implementa a interface CepFetcher para buscar dados do ViaCep
type ViaCepFetcher struct{}

// Método para buscar informações de um CEP no ViaCep
func (f *ViaCepFetcher) FetchCep(cep string) (*ViaCep, error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var data ViaCep
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// Função para inicializar o banco de dados SQLite e criar a tabela endereco, se não existir
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS endereco (
		cep TEXT PRIMARY KEY,
		logradouro TEXT,
		complemento TEXT,
		bairro TEXT,
		localidade TEXT,
		uf TEXT,
		ibge TEXT,
		gia TEXT,
		ddd TEXT,
		siafi TEXT
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ProcessCeps(db *sql.DB, fetcher CepFetcher, rawCeps []string) {
	for _, rawCep := range rawCeps {
		re := regexp.MustCompile(`\D`)
		cep := re.ReplaceAllString(rawCep, "")

		if len(cep) != 8 {
			fmt.Fprintf(os.Stderr, "CEP inválido: %s\n", rawCep)
			continue
		}

		var data ViaCep

		// Consulta se o CEP já existe no banco de dados
		row := db.QueryRow(`SELECT cep, logradouro, complemento, bairro, localidade, uf, ibge, gia, ddd, siafi FROM endereco WHERE REPLACE(cep, '-', '') = ?`, cep)
		err := row.Scan(&data.Cep, &data.Logradouro, &data.Complemento, &data.Bairro, &data.Localidade, &data.Uf, &data.Ibge, &data.Gia, &data.Ddd, &data.Siafi)
		if err != nil {
			if err != sql.ErrNoRows {
				fmt.Fprintf(os.Stderr, "Erro ao consultar banco de dados: %v\n", err)
				continue
			}

			// Caso não exista, buscamos na API e inserimos no banco de dados
			dataPtr, err := fetcher.FetchCep(cep)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao buscar CEP %s: %v\n", cep, err)
				continue
			}
			data = *dataPtr

			_, err = db.Exec("INSERT INTO endereco (cep, logradouro, complemento, bairro, localidade, uf, ibge, gia, ddd, siafi) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				data.Cep, data.Logradouro, data.Complemento, data.Bairro, data.Localidade, data.Uf, data.Ibge, data.Gia, data.Ddd, data.Siafi)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao inserir dados no banco de dados para o CEP %s: %v\n", cep, err)
				continue
			}
		}

		fmt.Printf("Dados para o CEP %s: %+v\n", cep, data)
	}
}
