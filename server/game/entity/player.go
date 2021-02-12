package entity

import (
	"github.com/floralbit/dungeon/game/event"
	gameModel "github.com/floralbit/dungeon/game/model"
	"github.com/floralbit/dungeon/game/util"
	"github.com/floralbit/dungeon/model"
	"github.com/google/uuid"
)

const (
	warriorTileID         = 21
	playerEnergyThreshold = 3
)

var ActivePlayers = map[uuid.UUID]*player{}

type player struct {
	entityData

	startingZone gameModel.Zone
	client       *model.Client
}

func NewPlayer(client *model.Client) *player {
	p := &player{
		entityData: entityData{
			UUID: client.Account.UUID,
			Name: client.Account.Username,
			Tile: warriorTileID,
			Type: gameModel.EntityTypePlayer,

			EnergyThreshold: playerEnergyThreshold,
		},

		client: client,
	}
	p.rollStats()

	ActivePlayers[p.UUID] = p
	return p
}

func (p *player) GetClient() *model.Client {
	return p.client
}

func (p *player) Act() gameModel.Action {
	a := p.QueuedAction
	p.QueuedAction = nil
	return a
}

func (p *player) Spawn(z gameModel.Zone) {
	for _, obj := range z.GetAllWorldObjects() {
		if obj.Type == gameModel.WorldObjectTypePlayerSpawn {
			p.X = obj.X
			p.Y = obj.Y
			break
		}
	}
	p.startingZone = z
	z.AddEntity(p)
	event.NotifyObservers(event.SpawnEvent{Entity: p})
}

// Despawn is for log off only, not changing zones (TODO: fix, leave vs. despawn)
func (p *player) Despawn() {
	event.NotifyObservers(event.DespawnEvent{Entity: p})
	p.zone.RemoveEntity(p)
	delete(ActivePlayers, p.UUID)
}

func (p *player) Die() {
	event.NotifyObservers(event.DieEvent{Entity: p})
	p.zone.RemoveEntity(p)
	p.rollStats()           // roll new stats cuz they're dead lol
	p.Spawn(p.startingZone) // send em back to the starting zone
	return
}

func (p *player) GainExp(xp int) {
	originalLevel := p.Stats.Level
	p.entityData.GainExp(xp)
	event.NotifyObservers(event.GainXPEvent{Entity: p, LeveledUp: originalLevel != p.Stats.Level})
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
