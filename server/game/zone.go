package game

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Zone ...
type Zone struct {
	Width, Height int
	Layers        []struct {
		Name string
		Data []int
	}
}

// Zones ...
var Zones = map[string]Zone{
	"town": loadZone(),
}

func loadZone() Zone {
	zoneFile, err := os.Open("../data/zones/town.json")
	if err != nil {
		panic(err)
	}
	defer zoneFile.Close()

	rawData, err := ioutil.ReadAll(zoneFile)
	if err != nil {
		panic(err)
	}

	var zone Zone
	json.Unmarshal(rawData, &zone)

	return zone
}
