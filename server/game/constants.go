package game

import "github.com/google/uuid"

var startingZoneUUID = uuid.MustParse("10f8b073-cbd7-46b7-a6e3-9cbdf68a933f")
var dungeonFloor1UUID = uuid.MustParse("6a67086c-eb9c-44c1-85b1-a140df7e4272")
var overworldSpawnObjectUUID = uuid.MustParse("b4f195f7-644a-4791-8177-c9eb69b10e9e")
var dungeonEntranceObjectUUID = uuid.MustParse("85ab1aaf-fcb2-4fa2-80e0-3cf54f8cad41")

var tiles = convertTileset(loadTiledTileset("../data/tileset.json"))
var monsterTemplates = loadMonsterTemplates()

var zones = loadZones()
var dungeonFloor1 = buildDungeonFloor() // TODO: put this logic into game loop
