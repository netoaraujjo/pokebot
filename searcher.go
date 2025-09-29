package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Ability struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Sprites struct {
	Other map[string]map[string]string `json:"other"`
}

type Pokemon struct {
	Sprites   Sprites                        `json:"sprites"`
	Abilities []map[string]map[string]string `json:"abilities"`
}

type Language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EffectEntry struct {
	EffectValue   string   `json:"effect"`
	LanguageValue Language `json:"language"`
}

type Effects struct {
	EffectEntries []EffectEntry `json:"effect_entries"`
}

const baseURL string = "https://pokeapi.co/api/v2/"

func search(pokemonName string) string {

	data := searchPokemon(pokemonName)

	searchAbilities(data["abilities_url"])

	return data["image"]
}

func searchPokemon(pokemonName string) map[string]string {
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
	// fmt.Println(data.Sprites.Other["official-artwork"]["front_default"])
	fmt.Println(data.Abilities[0]["ability"]["url"])

	result := make(map[string]string)
	result["image"] = data.Sprites.Other["official-artwork"]["front_default"]
	result["abilities_url"] = data.Abilities[0]["ability"]["url"]

	return result
}

func searchAbilities(abilitiesURL string) string {
	req, err := http.NewRequest("GET", abilitiesURL, nil)
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

	effects := Effects{}
	json.Unmarshal(content, &effects)
	for _, l := range effects.EffectEntries {
		if l.LanguageValue.Name == "en" {
			fmt.Println(l.EffectValue)
		}
	}

	return ""
}
