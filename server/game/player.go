package game

import (
	"github.com/floralbit/dungeon/game/event"
	"github.com/floralbit/dungeon/game/util"
	"github.com/floralbit/dungeon/model"
	"github.com/google/uuid"
)

const (
	warriorTileID         = 21
	maxValidMoveDist      = 3
	playerEnergyThreshold = 3
)

var activePlayers = map[uuid.UUID]*player{}

type player struct {
	entityData

	clientQueuedAction action
	client             *model.Client
}

func newPlayer(client *model.Client) *player {
	p := &player{
		entityData: entityData{
			UUID: client.Account.UUID,
			Name: client.Account.Username,
			Tile: warriorTileID,
			Type: entityTypePlayer,

			EnergyThreshold: playerEnergyThreshold,
		},

		client: client,
	}
	p.rollStats()

	activePlayers[p.UUID] = p
	return p
}

func (p *player) GetClient() *model.Client {
	return p.client
}

func (p *player) Act() action {
	a := p.queuedAction
	p.queuedAction = nil
	return a
}

func (p *player) Spawn(zoneUUID uuid.UUID) {
	z := zones[zoneUUID]
	for _, obj := range z.WorldObjects {
		if obj.Type == worldObjectTypePlayerSpawn {
			p.X = obj.X
			p.Y = obj.Y
			break
		}
	}
	zones[zoneUUID].addEntity(p)
	notifyObservers(event.SpawnEvent{Entity: p})
}

// Despawn is for log off only, not changing zones (TODO: fix, leave vs. despawn)
func (p *player) Despawn() {
	notifyObservers(event.DespawnEvent{Entity: p})
	p.zone.removeEntity(p)
	delete(activePlayers, p.UUID)
}

func (p *player) Die() {
	notifyObservers(event.DieEvent{Entity: p})
	p.zone.removeEntity(p)
	p.rollStats()             // roll new stats cuz they're dead lol
	p.Spawn(startingZoneUUID) // send em back to the starting zone
	return
}

func (p *player) GainExp(xp int) {
	originalLevel := p.Stats.Level
	p.entityData.GainExp(xp)
	notifyObservers(event.GainXPEvent{Entity: p, LeveledUp: originalLevel != p.Stats.Level})
}

func (p *player) rollStats() {
	p.Stats.Level = 1
	p.Stats.XP = 0
	p.Stats.XPToNextLevel = util.XPForLevel(2)

	// use 3d6 for stats
	r := util.Roll{6, 3, 0} // 3d6 + 0
	p.Stats.Strength = r.Roll()
	p.Stats.Dexterity = r.Roll()
	p.Stats.Constitution = r.Roll()
	p.Stats.Intelligence = r.Roll()
	p.Stats.Wisdom = r.Roll()
	p.Stats.Charisma = r.Roll()

	// hit dice for players is a d8, so HP = 8 + CON (1d8 + CON on level)
	p.Stats.MaxHP = 8 + util.Modifier(p.Stats.Constitution)
	if p.Stats.MaxHP <= 0 {
		p.Stats.MaxHP = 1
	}
	p.Stats.HP = p.Stats.MaxHP

	p.Stats.AC = 10 + util.Modifier(p.Stats.Dexterity)
}
