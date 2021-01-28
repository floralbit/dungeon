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

var gameTileset = loadTiledTileset("../data/tileset.json")
var tiles = convertTileset(gameTileset)

// TODO: support multiple tilesets, other options? or keep tiles pretty specific

type rawTiledMap struct {
	Width      int `json:"width"`
	Height     int `json:"height"`
	Properties []struct {
		Name  string          `json:"name"`
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	} `json:"properties"`
	Layers []struct {
		ID      int    `json:"id"`
		Width   int    `json:"width"`
		Height  int    `json:"height"`
		Data    []int  `json:"data"`
		Name    string `json:"name"`
		X       int    `json:"x"`
		Y       int    `json:"y"`
		Opacity int    `json:"opacity"`
		Visible bool   `json:"visible"`
	} `json:"layers"`
}

type rawTiledTileset struct {
	Columns      int    `json:"columns"`
	Image        string `json:"image"`
	ImageHeight  int    `json:"imageheight"`
	ImageWidth   int    `json:"imagewidth"`
	Margin       int    `json:"margin"`
	Name         string `json:"name"`
	Spacing      int    `json:"spacing"`
	TileCount    int    `json:"tilescount"`
	TiledVersion string `json:"tiledversion"`
	TileHeight   int    `json:"tileheight"`
	Tiles        []struct {
		ID         int `json:"id"`
		Properties []struct {
			Name  string          `json:"name"`
			Type  string          `json:"type"`
			Value json.RawMessage `json:"value"`
		} `json:"properties"`
	} `json:"tiles"`
	TileWidth int     `json:"tilewidth"`
	Type      string  `json:"type"`
	Version   float64 `json:"version"`
}

func loadZones() map[uuid.UUID]*zone {
	zones := map[uuid.UUID]*zone{}
	files, err := ioutil.ReadDir("../data/zones")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		split := strings.Split(file.Name(), ".")
		name, ext := split[0], split[1]
		if ext == "json" {
			zoneUUID := uuid.MustParse(name)
			zones[zoneUUID] = loadTiledMap(zoneUUID)
		}
	}

	return zones
}

func loadTiledMap(mapUUID uuid.UUID) *zone {
	mapFile, err := os.Open(fmt.Sprintf("../data/zones/%s.json", mapUUID.String()))
	if err != nil {
		log.Fatal(err)
	}
	defer mapFile.Close()

	rawData, err := ioutil.ReadAll(mapFile)
	if err != nil {
		log.Fatal(err)
	}

	var mapData rawTiledMap
	err = json.Unmarshal(rawData, &mapData)
	if err != nil {
		log.Fatal(err)
	}

	z := zone{
		UUID:   mapUUID,
		Width:  mapData.Width,
		Height: mapData.Height,
		Tiles:  []tile{},

		Entities:     map[uuid.UUID]*entity{},
		WorldObjects: map[uuid.UUID]*worldObject{},
	}

	for _, property := range mapData.Properties {
		if property.Name == "name" {
			err := json.Unmarshal(property.Value, &z.Name)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, layer := range mapData.Layers {
		if layer.Name == "ground" {
			for _, tileID := range layer.Data {
				z.Tiles = append(z.Tiles, tiles[tileID-1]) // -1 because of air tile (TODO: add air tile to -1 or something)
			}
		}
	}

	// TODO: worldObjects (either create from props, or object layer)

	return &z
}

func loadTiledTileset(path string) *rawTiledTileset {
	tilesetFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer tilesetFile.Close()

	rawData, err := ioutil.ReadAll(tilesetFile)
	if err != nil {
		log.Fatal(err)
	}

	var tilesetData rawTiledTileset
	err = json.Unmarshal(rawData, &tilesetData)
	if err != nil {
		log.Fatal(err)
	}

	return &tilesetData
}

func convertTileset(tileset *rawTiledTileset) map[int]tile {
	res := map[int]tile{}

	for _, t := range tileset.Tiles {
		var solid bool
		for _, prop := range t.Properties {
			if prop.Name == "solid" {
				err := json.Unmarshal(prop.Value, &solid)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		res[t.ID] = tile{
			ID:    t.ID,
			Solid: solid,
		}
	}

	return res
}
