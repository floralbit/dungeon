package game

import (
	"log"
	"math/rand"

	"github.com/floralbit/dungeon/game/gen"
	"github.com/google/uuid"
)

var dungeonFloor1UUID = uuid.MustParse("6a67086c-eb9c-44c1-85b1-a140df7e4272")
var dungeonFloor1 = buildDungeonFloor() // TODO: put this logic into game loop

var genTileTypeToTileID = map[gen.TileType][]int{
	gen.TileTypeWall:     {260, 262, 263, 264},
	gen.TileTypeGround:   {243, 244, 245, 246},
	gen.TileTypeHall:     {247},
	gen.TileTypeAir:      {216},
	gen.TileTypeDoor:     {224, 225, 230},
	gen.TileTypeTorch:    {409},
	gen.TileTypeEntrance: {211},
}

func buildDungeonFloor() *zone {
	level, err := gen.BuildLevel()
	if err != nil {
		log.Fatal(err)
	}

	z := &zone{
		UUID:   dungeonFloor1UUID,
		Name:   "dungeon",
		Width:  level.Width,
		Height: level.Height,

		Entities:     map[uuid.UUID]*entity{},
		WorldObjects: map[uuid.UUID]*worldObject{},
	}

	for y := 0; y < level.Height; y++ {
		for x := 0; x < level.Width; x++ {
			tileType := level.Tiles[x][y].Type

			tileIDOptions := genTileTypeToTileID[tileType]
			tileID := tileIDOptions[rand.Intn(len(tileIDOptions))]

			if tileType == gen.TileTypeEntrance {
				entranceUUID := uuid.New()
				z.WorldObjects[entranceUUID] = &worldObject{
					UUID: entranceUUID,
					Name: "exit",
					Tile: tileID,
					X:    x,
					Y:    y,
					Type: worldObjectTypePlayerSpawn,
					WarpTarget: &warpTarget{
						ZoneUUID: startingZoneUUID, // TODO: when multi-layer dungeon, assign to last layer
						X:        33,
						Y:        13,
					},
				}
				z.Tiles = append(z.Tiles, tiles[216]) // just add air for now, TODO: figure out better solution here
			} else {
				z.Tiles = append(z.Tiles, tiles[tileID])
			}
		}
	}

	// register zone - TODO: this is a hack for now, fix up later
	zones[dungeonFloor1UUID] = z
	return z
}

func (z *zone) findPlayerSpawn() (int, int) {
	for _, obj := range z.WorldObjects {
		if obj.Type == worldObjectTypePlayerSpawn {
			return obj.X, obj.Y
		}
	}

	return 0, 0
}
