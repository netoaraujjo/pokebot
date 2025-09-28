package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Sprites struct {
	Other map[string]map[string]string `json:"other"`
}

type Pokemon struct {
	Sprites Sprites `json:"sprites"`
}

func search(pokemonName string) string {
	baseURL := "https://pokeapi.co/api/v2/"

	req, err := http.NewRequest("GET", fmt.Sprintf("%spokemon/%s", baseURL, pokemonName), nil)
	if err != nil {
		log.Printf("Erro ao criar a requisição: %s\n", err)
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao obter a resposta: %s\n", err)
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Erro ao ler o corpo da resposta: %s\n", err)
	}

	data := Pokemon{}
	json.Unmarshal(content, &data)
	fmt.Println(data.Sprites.Other["official-artwork"]["front_default"])

	return data.Sprites.Other["official-artwork"]["front_default"]
}
