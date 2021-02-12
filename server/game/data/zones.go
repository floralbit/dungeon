package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/floralbit/dungeon/game/model"
	"github.com/floralbit/dungeon/game/zone"
	"github.com/google/uuid"
)

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

func LoadZones() map[uuid.UUID]*zone.Zone {
	zones := map[uuid.UUID]*zone.Zone{}
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

	for _, z := range zones {
		for _, obj := range z.WorldObjects {
			if obj.WarpTarget != nil {
				obj.WarpTarget.Zone = zones[obj.WarpTarget.ZoneUUID] // tie warp targets to zones via UUIDs
			}
		}
	}

	return zones
}

func loadTiledMap(mapUUID uuid.UUID) *zone.Zone {
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

	z := zone.Zone{
		UUID:   mapUUID,
		Width:  mapData.Width,
		Height: mapData.Height,
		Tiles:  []model.Tile{},

		Entities:     map[uuid.UUID]model.Entity{},
		WorldObjects: map[uuid.UUID]*model.WorldObject{},
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
				z.Tiles = append(z.Tiles, Tiles[tileID-1]) // -1 because of air tile (TODO: add air tile to -1 or something)
			}
		}
		if layer.Name == "world_objects" {
			for _, obj := range layer.Objects {
				var UUID uuid.UUID

				var hasWarpTarget, hasFullHeal bool
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
					if prop.Name == "fullHeal" {
						err := json.Unmarshal(prop.Value, &hasFullHeal)
						if err != nil {
							log.Fatal(err)
						}
					}
				}

				z.WorldObjects[UUID] = &model.WorldObject{
					UUID: UUID,
					Name: obj.Name,
					Tile: obj.TileID - 1,
					X:    obj.X / obj.Width,
					Y:    (obj.Y / obj.Height) - 1, // minus 1 because tiled objects start at the bottom left, tiles are top level (why the hell)

					Type: model.WorldObjectType(obj.Type),
				}
				if hasWarpTarget {
					z.WorldObjects[UUID].WarpTarget = &model.WarpTarget{
						ZoneUUID: warpTargetUUID,
						X:        warpTargetX,
						Y:        warpTargetY,
					}
				}
				if hasFullHeal {
					z.WorldObjects[UUID].HealZone = &model.HealZone{
						Full: true,
					}
				}
			}
		}
	}

	// TODO: worldObjects (either create from props, or object layer)

	return &z
}
