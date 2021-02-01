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
	gen.TileTypeWall:   {260, 262, 263, 264},
	gen.TileTypeGround: {243, 244, 245, 246},
	gen.TileTypeHall:   {247},
	gen.TileTypeAir:    {216},
	gen.TileTypeDoor:   {224, 225, 230},
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

	for x := 0; x < level.Width; x++ {
		for y := 0; y < level.Height; y++ {
			tileType := level.Tiles[x][y].Type
			tileIDOptions := genTileTypeToTileID[tileType]
			tileID := tileIDOptions[rand.Intn(len(tileIDOptions))]
			z.Tiles = append(z.Tiles, tiles[tileID])
		}
	}

	// register zone - TODO: this is a hack for now, fix up later
	zones[dungeonFloor1UUID] = z
	return z
}
