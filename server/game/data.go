package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
)

var gameTileset = loadTiledTileset("../data/tileset.json")
var tiles = convertTileset(gameTileset)
var monsterTemplates = loadMonsterTemplates()

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
		Objects []struct {
			Name       string `json:"name"`
			X          int    `json:"x"`
			Y          int    `json:"y"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
			Type       string `json:"type"`
			TileID     int    `json:"gid"`
			Properties []struct {
				Name  string          `json:"name"`
				Type  string          `json:"type"`
				Value json.RawMessage `json:"value"`
			} `json:"properties"`
		} `json:"objects"`
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

type monsterTemplate struct {
	Name string `json:"name"`
	Tile int    `json:"tile"`

	MoveSpeed    float64 `json:"move_speed"`
	AgroDistance float64 `json:"agro_distance"`

	Level int `json:"level"`

	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
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

		Entities:     map[uuid.UUID]entity{},
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
		if layer.Name == "world_objects" {
			for _, obj := range layer.Objects {
				var UUID uuid.UUID

				var hasWarpTarget bool
				var warpTargetUUID uuid.UUID
				var warpTargetX, warpTargetY int

				for _, prop := range obj.Properties {
					if prop.Name == "UUID" {
						err := json.Unmarshal(prop.Value, &UUID)
						if err != nil {
							log.Fatal(err)
						}
					}
					if prop.Name == "warpTargetUUID" {
						hasWarpTarget = true
						err := json.Unmarshal(prop.Value, &warpTargetUUID)
						if err != nil {
							log.Fatal(err)
						}
					}
					if prop.Name == "warpTargetX" {
						err := json.Unmarshal(prop.Value, &warpTargetX)
						if err != nil {
							log.Fatal(err)
						}
					}
					if prop.Name == "warpTargetY" {
						err := json.Unmarshal(prop.Value, &warpTargetY)
						if err != nil {
							log.Fatal(err)
						}
					}
				}

				z.WorldObjects[UUID] = &worldObject{
					UUID: UUID,
					Name: obj.Name,
					Tile: obj.TileID - 1,
					X:    obj.X / obj.Width,
					Y:    (obj.Y / obj.Height) - 1, // minus 1 because tiled objects start at the bottom left, tiles are top level (why the hell)

					Type: worldObjectType(obj.Type),
				}
				if hasWarpTarget {
					z.WorldObjects[UUID].WarpTarget = &warpTarget{
						ZoneUUID: warpTargetUUID,
						X:        warpTargetX,
						Y:        warpTargetY,
					}
				}
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
		res[t.ID] = tile{
			ID:    t.ID,
			Solid: solid,
			Name:  name,
		}
	}

	return res
}

func loadMonsterTemplates() map[monsterType]monsterTemplate {
	res := map[monsterType]monsterTemplate{}

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
