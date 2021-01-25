package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const startingZone = "town"

// Zone ...
type Zone struct {
	Width, Height int
	Layers        []struct {
		Name string
		Data []int
	}
}

// Zones ...
var Zones = loadZones()

func loadZones() map[string]Zone {
	zones := map[string]Zone{}

	files, err := ioutil.ReadDir("../data/zones")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		split := strings.Split(file.Name(), ".")
		name, ext := split[0], split[1]
		if ext == "json" {
			zones[name] = loadZone(name)
		}
	}

	return zones
}

func loadZone(name string) Zone {
	zoneFile, err := os.Open(fmt.Sprintf("../data/zones/%s.json", name))
	if err != nil {
		log.Fatal(err)
	}
	defer zoneFile.Close()

	rawData, err := ioutil.ReadAll(zoneFile)
	if err != nil {
		log.Fatal(err)
	}

	var zone Zone
	json.Unmarshal(rawData, &zone)

	return zone
}
