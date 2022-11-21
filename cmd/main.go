package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ViaCEP struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

type ApiCep struct {
	Cep        string `json:"code"`
	Logradouro string `json:"address"`
	Bairro     string `json:"district"`
	Localidade string `json:"city"`
	Uf         string `json:"state"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Busca CEP")
		fmt.Println("=========")
		fmt.Println("Uso:")
		fmt.Println("  go run cmd/main.go [cep]")
		fmt.Println("Exemplo:")
		fmt.Println("  go run cmd/main.go 06233-030")
		return
	}

	respostaViaCep := make(chan *ViaCEP)
	respostaApiCep := make(chan *ApiCep)

	go func() {
		cep, err := buscarViaCep(os.Args[1])
		if err != nil {
			fmt.Printf("Erro ao consultar Via Cep: %v", err)
		}
		respostaViaCep <- cep
	}()

	go func() {
		cep, err := buscarApiCep(os.Args[1])
		if err != nil {
			fmt.Printf("Erro ao consultar Api Cep: %v", err)
		}
		respostaApiCep <- cep
	}()

	select {
	case res := <-respostaViaCep:
		fmt.Printf("Reposta Via Cep: %v\n", res)
	case res := <-respostaApiCep:
		fmt.Printf("Reposta Api Cep: %v\n", res)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout :(")
	}
}

func buscarViaCep(cep string) (*ViaCEP, error) {
	req, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	var data ViaCEP
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func buscarApiCep(cep string) (*ApiCep, error) {
	req, err := http.Get("https://cdn.apicep.com/file/apicep/" + cep + ".json")
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	var data ApiCep
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
