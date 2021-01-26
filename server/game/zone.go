package game

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var startingZoneUUID = uuid.MustParse("10f8b073-cbd7-46b7-a6e3-9cbdf68a933f")

var zones = loadZones()

type tile struct {
	ID int `json:"id"`
	Solid bool `json:"solid"`
}

type zone struct {
	UUID uuid.UUID `json:"uuid"`
	Name string `json:"name"`
	Width int `json:"width"`
	Height int `json:"height"`
	Tiles []tile `json:"tiles"`

	Entities map[uuid.UUID]*entity `json:"entities"`
	WorldObjects map[uuid.UUID]*worldObject `json:"world_objects"`
}

// parsedZone is from the tiled format
type parsedZone struct {
	Width, Height int
	Layers []struct {
		Name string
		Data []int
	}
}

func loadZones() map[uuid.UUID]*zone {
	zones := map[uuid.UUID]*zone{}
	files, err := ioutil.ReadDir("../data/zones")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files{
		split := strings.Split(file.Name(), ".")
		name, ext := split[0], split[1]
		if ext == "json" {
			zoneUUID := uuid.MustParse(name)
			zones[zoneUUID] = loadZone(zoneUUID)
		}
	}

	return zones
}

func loadZone(zoneUUID uuid.UUID) *zone {
	zoneFile, err := os.Open(fmt.Sprintf("../data/zones/%s.json", zoneUUID.String()))
	if err != nil {
		log.Fatal(err)
	}
	defer zoneFile.Close()

	rawData, err := ioutil.ReadAll(zoneFile)
	if err != nil {
		log.Fatal(err)
	}

	var rawZone parsedZone
	json.Unmarshal(rawData, &rawZone)

	// TODO: parse additional things like worldObjects, etc
	tiles := []tile{}
	for _, t := range rawZone.Layers[0].Data {
		tiles = append(tiles, tile{
			ID: t,
			Solid: false, // TODO: populate from tilemap data
		})
	}

	return &zone{
		UUID: zoneUUID,
		Name: "TODO", // TODO: pull from tiled file
		Width: rawZone.Width,
		Height: rawZone.Height,
		Tiles: tiles,

		Entities: map[uuid.UUID]*entity{},
		WorldObjects: map[uuid.UUID]*worldObject{},
	}
}

func (z *zone) addEntity(e *entity) {
	z.Entities[e.UUID] = e
	e.zone = z
	z.send(newSpawnEvent(e))

	if e.Type == entityTypePlayer {
		e.send(newZoneLoadEvent(z)) // send player zone data
	}
}

func (z *zone) removeEntity(e *entity) {
	delete(z.Entities, e.UUID)
	z.send(newDespawnEvent(e))
}