package game

import (
	"github.com/floralbit/dungeon/model"
	"github.com/google/uuid"
)

const warriorTileID = 21

var activePlayers = map[uuid.UUID]*player{}

type player struct {
	entityData

	client *model.Client
}

func newPlayer(client *model.Client) *player {
	p := &player{
		entityData: entityData{
			UUID: client.Account.UUID,
			Name: client.Account.Username,
			Tile: warriorTileID,
			Type: entityTypePlayer,
		},

		client: client,
	}
	p.rollStats()

	activePlayers[p.UUID] = p
	return p
}

func (p *player) Move(x, y int) {
	t := p.zone.getTile(x, y)
	if t == nil {
		// edge of map, don't move
		p.Send(newMoveEvent(p.Data(), p.X, p.Y)) // tell them they're stationary
		return
	}

	if t.Solid {
		p.Send(newMoveEvent(p.Data(), p.X, p.Y)) // tell them they're stationary
		return
	}

	objs := p.zone.getWorldObjects(x, y)
	for _, obj := range objs {
		if obj.WarpTarget != nil {
			p.zone.removeEntity(p, false)
			zones[obj.WarpTarget.ZoneUUID].addEntity(p)
			p.X = obj.WarpTarget.X
			p.Y = obj.WarpTarget.Y
			return
		}
		if obj.HealZone != nil {
			if obj.HealZone.Full {
				p.HealFull()
				p.Send(newUpdateEvent(p.Data()))
				p.Send(newMoveEvent(p.Data(), p.X, p.Y)) // tell them they're stationary
				p.Send(newServerMessageEvent("You pray to your gods and are fully healed in their light."))
				return
			}
		}
	}

	for _, e := range p.zone.Entities {
		if e.Data().X == x && e.Data().Y == y {
			p.Send(newMoveEvent(p.Data(), p.X, p.Y)) // tell them they're stationary, because attacking
			p.Attack(e)
			return
		}
	}

	p.X = x
	p.Y = y

	p.zone.send(newMoveEvent(p.Data(), x, y))
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
	zones[startingZoneUUID].addEntity(p)
}

func (p *player) Despawn(becauseDeath bool) {
	p.zone.removeEntity(p, becauseDeath)
	delete(activePlayers, p.UUID)
}

func (p *player) Send(event serverEvent) {
	p.client.In <- event
}

func (p *player) rollStats() {
	p.Stats.Level = 1
	p.Stats.XP = 0

	// use 3d6 for stats
	r := roll{6, 3, 0} // 3d6 + 0
	p.Stats.Strength = r.roll()
	p.Stats.Dexterity = r.roll()
	p.Stats.Constitution = r.roll()
	p.Stats.Intelligence = r.roll()
	p.Stats.Wisdom = r.roll()
	p.Stats.Charisma = r.roll()

	// hit dice for players is a d8, so HP = 8 + CON (1d8 + CON on level)
	p.Stats.MaxHP = 8 + modifier(p.Stats.Constitution)
	if p.Stats.MaxHP <= 0 {
		p.Stats.MaxHP = 1
	}
	p.Stats.HP = p.Stats.MaxHP

	p.Stats.AC = 10 + modifier(p.Stats.Dexterity)
}

func (p *player) Die() {
	p.zone.removeEntity(p, true)
	p.rollStats()             // roll new stats cuz they're dead lol
	p.Spawn(startingZoneUUID) // send em back to the starting zone
	return
}
