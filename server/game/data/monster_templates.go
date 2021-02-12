package data

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var MonsterTemplates = loadMonsterTemplates()

type MonsterTemplate struct {
	Name string `json:"name"`
	Tile int    `json:"tile"`

	EnergyThreshold int     `json:"energy_threshold"`
	AgroDistance    float64 `json:"agro_distance"`

	Level int `json:"level"`

	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

func loadMonsterTemplates() map[string]MonsterTemplate {
	res := map[string]MonsterTemplate{}

	monsterTemplateFile, err := os.Open("../data/monsters/monsters.json")
	if err != nil {
		log.Fatal(err)
	}
	defer monsterTemplateFile.Close()

	rawData, err := ioutil.ReadAll(monsterTemplateFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(rawData, &res)
	if err != nil {
		log.Fatal(err)
	}

	return res
}
