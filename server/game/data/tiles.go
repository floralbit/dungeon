package data

import (
	"encoding/json"
	"github.com/floralbit/dungeon/game/model"
	"io/ioutil"
	"log"
	"os"
)

var Tiles = convertTileset(loadTiledTileset("../data/tileset.json"))

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

func convertTileset(tileset *rawTiledTileset) map[int]model.Tile {
	res := map[int]model.Tile{}

	for _, t := range tileset.Tiles {
		var solid bool
		var name string
		for _, prop := range t.Properties {
			if prop.Name == "solid" {
				err := json.Unmarshal(prop.Value, &solid)
				if err != nil {
					log.Fatal(err)
				}
			}
			if prop.Name == "name" {
				err := json.Unmarshal(prop.Value, &name)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		res[t.ID] = model.Tile{
			ID:    t.ID,
			Solid: solid,
			Name:  name,
		}
	}

	return res
}
