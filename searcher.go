package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type AbilityInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Ability struct {
	Ability AbilityInfo `json:"ability"`
}

type Cries struct {
	Latest string `json:"latest"`
}

type Sprites struct {
	Sprite Other `json:"other"`
}

type Other struct {
	Artwork OfficialArtwork `json:"official-artwork"`
}

type OfficialArtwork struct {
	FrontDefault string `json:"front_default"`
}

type PokemonData struct {
	Abilities []Ability `json:"abilities"`
	Sprites   Sprites   `json:"sprites"`
	Cries     Cries     `json:"cries"`
}

// Struct simplificada de um Pokemon
type Pokemon struct {
	Image     string
	Sound     string
	Abilities map[string]string
}

func (pokemon *Pokemon) FormatAbilities() string {
	caption := ""
	for name, effect := range pokemon.Abilities {
		caption = fmt.Sprintf("%s\n\n*%s*\n%s", caption, name, effect)
	}
	return caption
}

// Structs das abilidades
type Language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EffectEntry struct {
	Effect   string   `json:"effect"`
	Language Language `json:"language"`
}

type EffectEntries struct {
	EffectEntries []EffectEntry `json:"effect_entries"`
	Name          string        `json:"name"`
}

const baseURL string = "https://pokeapi.co/api/v2/"

func search(pokemonName string) Pokemon {

	// data := searchPokemon(pokemonName)
	pokemonData := searchPokemon(pokemonName)

	abilities := searchAbilities(pokemonData.Abilities)

	return Pokemon{
		Image:     pokemonData.Sprites.Sprite.Artwork.FrontDefault,
		Sound:     pokemonData.Cries.Latest,
		Abilities: abilities,
	}
}

func searchPokemon(pokemonName string) PokemonData {
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

	pokemon := PokemonData{}
	json.Unmarshal(content, &pokemon)

	return pokemon
}

func searchAbilities(abilities []Ability) map[string]string {

	abilitiesList := make(map[string]string)

	for _, ability := range abilities {
		req, err := http.NewRequest("GET", ability.Ability.URL, nil)
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

		effectEntries := EffectEntries{}
		json.Unmarshal(content, &effectEntries)

		for _, effect := range effectEntries.EffectEntries {
			if effect.Language.Name == "en" {
				abilitiesList[effectEntries.Name] = effect.Effect
			}
		}
	}

	return abilitiesList
}

// func TranslateEffect(effect string) {
// 	jsonBody, err := json.Marshal(map[string]interface{}{
// 		"q":effect,
// 		"source": "en",
// 		"targe": "pr-BR",
// 		"format": "text",
// 		"alternatives": 3,
// 		"api_key": ""
// 	})

// 	req, err := http.NewRequest("POST", "https://libretranslate.com/translate", nil)
// 	if err != nil {
// 		log.Printf("Erro ao traduzir o texto: %s\n", err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// }
