package game

import (
	"log"
	"math/rand"

	"github.com/floralbit/dungeon/game/gen"
	"github.com/floralbit/dungeon/game/model"
	"github.com/google/uuid"
)

const (
	goblinLikelihood   = .5 // 50%
	skeletonLikelihood = .5 // 50%
)

var genTileTypeToTileID = map[gen.TileType][]int{
	gen.TileTypeWall:     {260, 262, 263, 264},
	gen.TileTypeGround:   {243, 244, 245, 246},
	gen.TileTypeHall:     {247},
	gen.TileTypeAir:      {216},
	gen.TileTypeDoor:     {224, 225, 230},
	gen.TileTypeTorch:    {409},
	gen.TileTypeEntrance: {211},
	gen.TileTypeExit:     {210},
}

func buildDungeonFloor() *zone {
	level, err := gen.BuildLevel()
	if err != nil {
		log.Fatal(err)
	}

	z := &zone{
		UUID:   dungeonFloor1UUID,
		Name:   "Dungeon Floor 1",
		Width:  level.Width,
		Height: level.Height,

		Entities:     map[uuid.UUID]model.Entity{},
		WorldObjects: map[uuid.UUID]*model.WorldObject{},
	}

	for y := 0; y < level.Height; y++ {
		for x := 0; x < level.Width; x++ {
			tileType := level.Tiles[x][y].Type

			tileIDOptions := genTileTypeToTileID[tileType]
			tileID := tileIDOptions[rand.Intn(len(tileIDOptions))]

			if tileType == gen.TileTypeEntrance {
				dungeonEntrance := zones[startingZoneUUID].WorldObjects[dungeonEntranceObjectUUID]

				entranceUUID := uuid.New()
				z.WorldObjects[entranceUUID] = &model.WorldObject{
					UUID: entranceUUID,
					Name: "Dungeon Exit",
					Tile: tileID,
					X:    x,
					Y:    y,
					Type: model.WorldObjectTypePortal,
					WarpTarget: &model.WarpTarget{
						Zone:     zones[startingZoneUUID],
						ZoneUUID: startingZoneUUID, // TODO: when multi-layer dungeon, assign to last layer
						X:        dungeonEntrance.X,
						Y:        dungeonEntrance.Y,
					},
				}
				// tie overworld entrance to stairs
				dungeonEntrance.WarpTarget = &model.WarpTarget{
					Zone:     z,
					ZoneUUID: dungeonFloor1UUID,
					X:        x,
					Y:        y,
				}
				z.Tiles = append(z.Tiles, tiles[216]) // just add air for now, TODO: figure out better solution here
			} else {
				z.Tiles = append(z.Tiles, tiles[tileID])
			}
		}
	}

	// register zone - TODO: this is a hack for now, fix up later
	zones[dungeonFloor1UUID] = z

	spawnMonsters(level, z)

	return z
}

func spawnMonsters(l *gen.Level, z *zone) {
	for x := 0; x < l.Width; x++ {
		for y := 0; y < l.Height; y++ {
			if l.Objects[x][y] != nil && l.Objects[x][y].Type == gen.ObjectTypeMonsterSlot {
				var m *monster
				for m == nil {
					if rand.Float32() < goblinLikelihood {
						m = newMonster(monsterTypeGoblin)
					} else if rand.Float32() < skeletonLikelihood {
						m = newMonster(monsterTypeSkeleton)
					}
				}

				m.X = x
				m.Y = y
				z.AddEntity(m)
			}
		}
	}
}
