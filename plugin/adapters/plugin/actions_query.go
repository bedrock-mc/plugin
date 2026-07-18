package plugin

import (
	"fmt"

	pb "github.com/bedrock-mc/plugin/proto/generated/go"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

func (m *Manager) handleWorldQueryEntities(p *pluginProcess, correlationID string, act *pb.WorldQueryEntitiesAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	entityRefs, err := world.Call(m.ctx, w, func(tx *world.Tx) ([]*pb.EntityRef, error) {
		entities := make([]world.Entity, 0)
		for e := range tx.Entities() {
			entities = append(entities, e)
		}
		return protoEntityRefs(entities), nil
	})
	if err != nil {
		m.sendActionError(p, correlationID, err.Error())
		return
	}
	m.sendActionResult(p, &pb.ActionResult{
		CorrelationId: correlationID,
		Status:        &pb.ActionStatus{Ok: true},
		Result: &pb.ActionResult_WorldEntities{WorldEntities: &pb.WorldEntitiesResult{
			World:    protoWorldRef(w),
			Entities: entityRefs,
		}},
	})
}

func (m *Manager) handleWorldQueryPlayers(p *pluginProcess, correlationID string, act *pb.WorldQueryPlayersAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	playerRefs, err := world.Call(m.ctx, w, func(tx *world.Tx) ([]*pb.EntityRef, error) {
		players := make([]world.Entity, 0)
		for pl := range tx.Players() {
			players = append(players, pl)
		}
		return protoEntityRefs(players), nil
	})
	if err != nil {
		m.sendActionError(p, correlationID, err.Error())
		return
	}
	m.sendActionResult(p, &pb.ActionResult{
		CorrelationId: correlationID,
		Status:        &pb.ActionStatus{Ok: true},
		Result: &pb.ActionResult_WorldPlayers{WorldPlayers: &pb.WorldPlayersResult{
			World:   protoWorldRef(w),
			Players: playerRefs,
		}},
	})
}

func (m *Manager) handleWorldQueryEntitiesWithin(p *pluginProcess, correlationID string, act *pb.WorldQueryEntitiesWithinAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	box, ok := bboxFromProto(act.Box)
	if !ok {
		m.sendActionError(p, correlationID, "invalid bounding box")
		return
	}
	entityRefs, err := world.Call(m.ctx, w, func(tx *world.Tx) ([]*pb.EntityRef, error) {
		entities := make([]world.Entity, 0)
		for e := range tx.EntitiesWithin(box) {
			entities = append(entities, e)
		}
		return protoEntityRefs(entities), nil
	})
	if err != nil {
		m.sendActionError(p, correlationID, err.Error())
		return
	}
	m.sendActionResult(p, &pb.ActionResult{
		CorrelationId: correlationID,
		Status:        &pb.ActionStatus{Ok: true},
		Result: &pb.ActionResult_WorldEntitiesWithin{WorldEntitiesWithin: &pb.WorldEntitiesWithinResult{
			World:    protoWorldRef(w),
			Box:      protoBBox(box),
			Entities: entityRefs,
		}},
	})
}

func (m *Manager) handleWorldQueryDefaultGameMode(p *pluginProcess, correlationID string, act *pb.WorldQueryDefaultGameModeAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	mode := w.DefaultGameMode()
	m.sendActionResult(p, &pb.ActionResult{
		CorrelationId: correlationID,
		Status:        &pb.ActionStatus{Ok: true},
		Result: &pb.ActionResult_WorldDefaultGameMode{WorldDefaultGameMode: &pb.WorldDefaultGameModeResult{
			World:    protoWorldRef(w),
			GameMode: func() pb.GameMode { id, _ := world.GameModeID(mode); return pb.GameMode(id) }(),
		}},
	})
}

func (m *Manager) handleWorldQueryPlayerSpawn(p *pluginProcess, correlationID string, act *pb.WorldQueryPlayerSpawnAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	if act.PlayerUuid == "" {
		m.sendActionError(p, correlationID, "missing player_uuid")
		return
	}
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		m.sendActionError(p, correlationID, "invalid player_uuid")
		return
	}
	spawn := w.PlayerSpawn(id)
	m.sendActionResult(p, &pb.ActionResult{
		CorrelationId: correlationID,
		Status:        &pb.ActionStatus{Ok: true},
		Result: &pb.ActionResult_WorldPlayerSpawn{WorldPlayerSpawn: &pb.WorldPlayerSpawnResult{
			World:      protoWorldRef(w),
			PlayerUuid: act.PlayerUuid,
			Spawn:      protoBlockPos(spawn),
		}},
	})
}

func (m *Manager) handleWorldQueryBlock(p *pluginProcess, correlationID string, act *pb.WorldQueryBlockAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	if act.Position == nil {
		m.sendActionError(p, correlationID, "missing position")
		return
	}
	pos := cube.Pos{int(act.Position.X), int(act.Position.Y), int(act.Position.Z)}
	block, err := world.Call(m.ctx, w, func(tx *world.Tx) (world.Block, error) {
		return tx.Block(pos), nil
	})
	if err != nil {
		m.sendActionError(p, correlationID, err.Error())
		return
	}
	m.sendActionResult(p, &pb.ActionResult{
		CorrelationId: correlationID,
		Status:        &pb.ActionStatus{Ok: true},
		Result: &pb.ActionResult_WorldBlock{WorldBlock: &pb.WorldBlockResult{
			World:    protoWorldRef(w),
			Position: act.Position,
			Block:    protoBlockState(block),
		}},
	})
}

func (m *Manager) handleWorldQueryBiome(p *pluginProcess, correlationID string, act *pb.WorldQueryBiomeAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	if act.Position == nil {
		m.sendActionError(p, correlationID, "missing position")
		return
	}
	pos := cube.Pos{int(act.Position.X), int(act.Position.Y), int(act.Position.Z)}
	biome, err := world.Call(m.ctx, w, func(tx *world.Tx) (world.Biome, error) {
		return tx.Biome(pos), nil
	})
	if err != nil {
		m.sendActionError(p, correlationID, err.Error())
		return
	}
	biomeID := fmt.Sprintf("%d", biome.EncodeBiome())
	m.sendActionResult(p, &pb.ActionResult{
		CorrelationId: correlationID,
		Status:        &pb.ActionStatus{Ok: true},
		Result: &pb.ActionResult_WorldBiome{WorldBiome: &pb.WorldBiomeResult{
			World:    protoWorldRef(w),
			Position: act.Position,
			BiomeId:  biomeID,
		}},
	})
}

func (m *Manager) handleWorldQueryLight(p *pluginProcess, correlationID string, act *pb.WorldQueryLightAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	if act.Position == nil {
		m.sendActionError(p, correlationID, "missing position")
		return
	}
	pos := cube.Pos{int(act.Position.X), int(act.Position.Y), int(act.Position.Z)}
	lightLevel, err := world.Call(m.ctx, w, func(tx *world.Tx) (uint8, error) {
		return tx.Light(pos), nil
	})
	if err != nil {
		m.sendActionError(p, correlationID, err.Error())
		return
	}
	m.sendActionResult(p, &pb.ActionResult{
		CorrelationId: correlationID,
		Status:        &pb.ActionStatus{Ok: true},
		Result: &pb.ActionResult_WorldLight{WorldLight: &pb.WorldLightResult{
			World:      protoWorldRef(w),
			Position:   act.Position,
			LightLevel: int32(lightLevel),
		}},
	})
}

func (m *Manager) handleWorldQuerySkyLight(p *pluginProcess, correlationID string, act *pb.WorldQuerySkyLightAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	if act.Position == nil {
		m.sendActionError(p, correlationID, "missing position")
		return
	}
	pos := cube.Pos{int(act.Position.X), int(act.Position.Y), int(act.Position.Z)}
	skyLightLevel, err := world.Call(m.ctx, w, func(tx *world.Tx) (uint8, error) {
		return tx.SkyLight(pos), nil
	})
	if err != nil {
		m.sendActionError(p, correlationID, err.Error())
		return
	}
	m.sendActionResult(p, &pb.ActionResult{
		CorrelationId: correlationID,
		Status:        &pb.ActionStatus{Ok: true},
		Result: &pb.ActionResult_WorldSkyLight{WorldSkyLight: &pb.WorldSkyLightResult{
			World:         protoWorldRef(w),
			Position:      act.Position,
			SkyLightLevel: int32(skyLightLevel),
		}},
	})
}

func (m *Manager) handleWorldQueryTemperature(p *pluginProcess, correlationID string, act *pb.WorldQueryTemperatureAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	if act.Position == nil {
		m.sendActionError(p, correlationID, "missing position")
		return
	}
	pos := cube.Pos{int(act.Position.X), int(act.Position.Y), int(act.Position.Z)}
	temperature, err := world.Call(m.ctx, w, func(tx *world.Tx) (float64, error) {
		return tx.Temperature(pos), nil
	})
	if err != nil {
		m.sendActionError(p, correlationID, err.Error())
		return
	}
	m.sendActionResult(p, &pb.ActionResult{
		CorrelationId: correlationID,
		Status:        &pb.ActionStatus{Ok: true},
		Result: &pb.ActionResult_WorldTemperature{WorldTemperature: &pb.WorldTemperatureResult{
			World:       protoWorldRef(w),
			Position:    act.Position,
			Temperature: temperature,
		}},
	})
}

func (m *Manager) handleWorldQueryHighestBlock(p *pluginProcess, correlationID string, act *pb.WorldQueryHighestBlockAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	y, err := world.Call(m.ctx, w, func(tx *world.Tx) (int, error) {
		return tx.HighestBlock(int(act.X), int(act.Z)), nil
	})
	if err != nil {
		m.sendActionError(p, correlationID, err.Error())
		return
	}
	m.sendActionResult(p, &pb.ActionResult{
		CorrelationId: correlationID,
		Status:        &pb.ActionStatus{Ok: true},
		Result: &pb.ActionResult_WorldHighestBlock{WorldHighestBlock: &pb.WorldHighestBlockResult{
			World: protoWorldRef(w),
			X:     act.X,
			Z:     act.Z,
			Y:     int32(y),
		}},
	})
}

func (m *Manager) handleWorldQueryRainingAt(p *pluginProcess, correlationID string, act *pb.WorldQueryRainingAtAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	if act.Position == nil {
		m.sendActionError(p, correlationID, "missing position")
		return
	}
	pos := cube.Pos{int(act.Position.X), int(act.Position.Y), int(act.Position.Z)}
	raining, err := world.Call(m.ctx, w, func(tx *world.Tx) (bool, error) {
		return tx.RainingAt(pos), nil
	})
	if err != nil {
		m.sendActionError(p, correlationID, err.Error())
		return
	}
	m.sendActionResult(p, &pb.ActionResult{
		CorrelationId: correlationID,
		Status:        &pb.ActionStatus{Ok: true},
		Result: &pb.ActionResult_WorldRainingAt{WorldRainingAt: &pb.WorldRainingAtResult{
			World:    protoWorldRef(w),
			Position: act.Position,
			Raining:  raining,
		}},
	})
}

func (m *Manager) handleWorldQuerySnowingAt(p *pluginProcess, correlationID string, act *pb.WorldQuerySnowingAtAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	if act.Position == nil {
		m.sendActionError(p, correlationID, "missing position")
		return
	}
	pos := cube.Pos{int(act.Position.X), int(act.Position.Y), int(act.Position.Z)}
	snowing, err := world.Call(m.ctx, w, func(tx *world.Tx) (bool, error) {
		return tx.SnowingAt(pos), nil
	})
	if err != nil {
		m.sendActionError(p, correlationID, err.Error())
		return
	}
	m.sendActionResult(p, &pb.ActionResult{
		CorrelationId: correlationID,
		Status:        &pb.ActionStatus{Ok: true},
		Result: &pb.ActionResult_WorldSnowingAt{WorldSnowingAt: &pb.WorldSnowingAtResult{
			World:    protoWorldRef(w),
			Position: act.Position,
			Snowing:  snowing,
		}},
	})
}

func (m *Manager) handleWorldQueryThunderingAt(p *pluginProcess, correlationID string, act *pb.WorldQueryThunderingAtAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	if act.Position == nil {
		m.sendActionError(p, correlationID, "missing position")
		return
	}
	pos := cube.Pos{int(act.Position.X), int(act.Position.Y), int(act.Position.Z)}
	thundering, err := world.Call(m.ctx, w, func(tx *world.Tx) (bool, error) {
		return tx.ThunderingAt(pos), nil
	})
	if err != nil {
		m.sendActionError(p, correlationID, err.Error())
		return
	}
	m.sendActionResult(p, &pb.ActionResult{
		CorrelationId: correlationID,
		Status:        &pb.ActionStatus{Ok: true},
		Result: &pb.ActionResult_WorldThunderingAt{WorldThunderingAt: &pb.WorldThunderingAtResult{
			World:      protoWorldRef(w),
			Position:   act.Position,
			Thundering: thundering,
		}},
	})
}

func (m *Manager) handleWorldQueryLiquid(p *pluginProcess, correlationID string, act *pb.WorldQueryLiquidAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	if act.Position == nil {
		m.sendActionError(p, correlationID, "missing position")
		return
	}
	pos := cube.Pos{int(act.Position.X), int(act.Position.Y), int(act.Position.Z)}
	liquidState, err := world.Call(m.ctx, w, func(tx *world.Tx) (*pb.LiquidState, error) {
		if liq, ok := tx.Liquid(pos); ok {
			return protoLiquidState(liq), nil
		}
		return nil, nil
	})
	if err != nil {
		m.sendActionError(p, correlationID, err.Error())
		return
	}
	m.sendActionResult(p, &pb.ActionResult{
		CorrelationId: correlationID,
		Status:        &pb.ActionStatus{Ok: true},
		Result: &pb.ActionResult_WorldLiquid{WorldLiquid: &pb.WorldLiquidResult{
			World:    protoWorldRef(w),
			Position: act.Position,
			Liquid:   liquidState,
		}},
	})
}

// World mutation handlers
