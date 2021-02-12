package dungeon

import (
	"fmt"
	"github.com/floralbit/dungeon/game/data"
	"github.com/floralbit/dungeon/game/dungeon/gen"
	"github.com/floralbit/dungeon/game/entity"
	"github.com/floralbit/dungeon/game/model"
	"github.com/floralbit/dungeon/game/zone"
	"github.com/google/uuid"
	"log"
	"math/rand"
)

var dungeonEntranceObjectUUID = uuid.MustParse("85ab1aaf-fcb2-4fa2-80e0-3cf54f8cad41")

const (
	goblinLikelihood   = .5 // 50%
	skeletonLikelihood = .5 // 50%

	dungeonDepth   = 3
	stairsDownTile = 210
)

var genTileTypeToTileID = map[gen.TileType][]int{
	gen.TileTypeWall:     {260, 262, 263, 264},
	gen.TileTypeGround:   {243, 244, 245, 246},
	gen.TileTypeHall:     {247},
	gen.TileTypeAir:      {216},
	gen.TileTypeDoor:     {224, 225, 230},
	gen.TileTypeTorch:    {409},
	gen.TileTypeEntrance: {211},
	gen.TileTypeExit:     {stairsDownTile},
}

func BuildDungeon(overworld *zone.Zone) map[uuid.UUID]*zone.Zone {
	d := map[uuid.UUID]*zone.Zone{}

	var priorFloor *zone.Zone
	for i := 0; i < dungeonDepth; i++ {
		floorUUID := uuid.New()
		d[floorUUID] = buildFloor(floorUUID, i, overworld, priorFloor)
		priorFloor = d[floorUUID]
	}

	return d
}

func buildFloor(floorUUID uuid.UUID, depth int, overworld *zone.Zone, priorFloor *zone.Zone) *zone.Zone {
	level, err := gen.BuildLevel()
	if err != nil {
		log.Fatal(err)
	}

	z := &zone.Zone{
		UUID:   floorUUID,
		Name:   fmt.Sprintf("Dungeon Floor %d", depth+1),
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
				if depth == 0 {
					tieToOverworld(z, tileID, x, y, overworld)
				} else {
					tieToPriorFloor(z, tileID, x, y, priorFloor)
				}
			} else if tileType == gen.TileTypeExit {
				if depth < dungeonDepth-1 {
					exitUUID := uuid.New()
					z.WorldObjects[exitUUID] = &model.WorldObject{
						UUID: exitUUID,
						Name: "Stairs Down",
						Tile: tileID,
						X:    x,
						Y:    y,
						Type: model.WorldObjectTypePortal,
					}
				}
				z.Tiles = append(z.Tiles, data.Tiles[243]) // just add floor for now, TODO: add treasure chamber to end
			} else {
				z.Tiles = append(z.Tiles, data.Tiles[tileID])
			}
		}
	}

	spawnMonsters(level, z)
	return z
}

func tieToOverworld(z *zone.Zone, tileID, x, y int, overworld *zone.Zone) {
	dungeonEntrance := overworld.WorldObjects[dungeonEntranceObjectUUID]

	entranceUUID := uuid.New()
	z.WorldObjects[entranceUUID] = &model.WorldObject{
		UUID: entranceUUID,
		Name: "Stairs Up",
		Tile: tileID,
		X:    x,
		Y:    y,
		Type: model.WorldObjectTypePortal,
		WarpTarget: &model.WarpTarget{
			Zone:     overworld,
			ZoneUUID: overworld.UUID,
			X:        dungeonEntrance.X,
			Y:        dungeonEntrance.Y,
		},
	}

	// tie overworld entrance to stairs
	dungeonEntrance.WarpTarget = &model.WarpTarget{
		Zone:     z,
		ZoneUUID: z.UUID,
		X:        x,
		Y:        y,
	}
	z.Tiles = append(z.Tiles, data.Tiles[243]) // just add floor for now, TODO: figure out better solution here
}

func tieToPriorFloor(z *zone.Zone, tileID, x, y int, priorFloor *zone.Zone) {
	var stairsDownPriorFloor *model.WorldObject

	for _, obj := range priorFloor.WorldObjects {
		if obj.Tile == stairsDownTile {
			stairsDownPriorFloor = obj
			break
		}
	}

	entranceUUID := uuid.New()
	z.WorldObjects[entranceUUID] = &model.WorldObject{
		UUID: entranceUUID,
		Name: "Stairs Up",
		Tile: tileID,
		X:    x,
		Y:    y,
		Type: model.WorldObjectTypePortal,
		WarpTarget: &model.WarpTarget{
			Zone:     priorFloor,
			ZoneUUID: priorFloor.UUID,
			X:        stairsDownPriorFloor.X,
			Y:        stairsDownPriorFloor.Y,
		},
	}

	// tie overworld entrance to stairs
	stairsDownPriorFloor.WarpTarget = &model.WarpTarget{
		Zone:     z,
		ZoneUUID: z.UUID,
		X:        x,
		Y:        y,
	}
	z.Tiles = append(z.Tiles, data.Tiles[243]) // just add floor for now, TODO: figure out better solution here
}

func spawnMonsters(l *gen.Level, z *zone.Zone) {
	for x := 0; x < l.Width; x++ {
		for y := 0; y < l.Height; y++ {
			if l.Objects[x][y] != nil && l.Objects[x][y].Type == gen.ObjectTypeMonsterSlot {
				var m *entity.Monster
				for m == nil {
					if rand.Float32() < goblinLikelihood {
						m = entity.NewMonster(entity.MonsterTypeGoblin)
					} else if rand.Float32() < skeletonLikelihood {
						m = entity.NewMonster(entity.MonsterTypeSkeleton)
					}
				}

				m.X = x
				m.Y = y
				z.AddEntity(m)
			}
		}
	}
}
