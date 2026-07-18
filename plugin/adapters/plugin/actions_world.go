package plugin

import (
	"fmt"
	"time"

	pb "github.com/bedrock-mc/plugin/proto/generated/go"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

func (m *Manager) handleWorldSetDefaultGameMode(p *pluginProcess, correlationID string, act *pb.WorldSetDefaultGameModeAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	mode, ok := world.GameModeByID(int(act.GameMode))
	if !ok {
		m.sendActionError(p, correlationID, "unknown game mode")
		return
	}
	w.SetDefaultGameMode(mode)
	m.sendActionOK(p, correlationID)
}

func (m *Manager) handleWorldBuildStructure(p *pluginProcess, correlationID string, act *pb.WorldBuildStructureAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	if act.Origin == nil {
		m.sendActionError(p, correlationID, "missing origin")
		return
	}
	if act.Structure == nil {
		m.sendActionError(p, correlationID, "missing structure")
		return
	}
	ps, err := buildProtoStructure(act.Structure)
	if err != nil {
		m.sendActionError(p, correlationID, err.Error())
		return
	}
	origin := cube.Pos{int(act.Origin.X), int(act.Origin.Y), int(act.Origin.Z)}
	m.completeWorldTask(p, correlationID, w.Do(func(tx *world.Tx) {
		tx.BuildStructure(origin, ps)
	}))
}

func (m *Manager) handleWorldSetDifficulty(p *pluginProcess, correlationID string, act *pb.WorldSetDifficultyAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	diff, ok := world.DifficultyByID(int(act.Difficulty))
	if !ok {
		m.sendActionError(p, correlationID, "unknown difficulty")
		return
	}
	w.SetDifficulty(diff)
	m.sendActionOK(p, correlationID)
}

// Player movement toggles
func (m *Manager) handleWorldSetTickRange(p *pluginProcess, correlationID string, act *pb.WorldSetTickRangeAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	if act.TickRange < 0 {
		m.sendActionError(p, correlationID, "tick range must be non-negative")
		return
	}
	w.SetTickRange(int(act.TickRange))
	m.sendActionOK(p, correlationID)
}

func (m *Manager) handleWorldSetBlock(p *pluginProcess, correlationID string, act *pb.WorldSetBlockAction) {
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
	var blk world.Block
	var ok bool
	if act.Block != nil {
		blk, ok = blockFromProto(act.Block)
		if !ok {
			m.sendActionError(p, correlationID, "unknown block")
			return
		}
	}
	m.completeWorldTask(p, correlationID, w.Do(func(tx *world.Tx) {
		tx.SetBlock(pos, blk, nil)
	}))
}

func (m *Manager) handleWorldPlaySound(p *pluginProcess, correlationID string, act *pb.WorldPlaySoundAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	pos, ok := vec3FromProto(act.Position)
	if !ok {
		m.sendActionError(p, correlationID, "invalid position")
		return
	}
	s := soundFromProto(act.Sound)
	m.completeWorldTask(p, correlationID, w.Do(func(tx *world.Tx) {
		tx.PlaySound(pos, s)
	}))
}

func (m *Manager) handleWorldAddParticle(p *pluginProcess, correlationID string, act *pb.WorldAddParticleAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	pos, ok := vec3FromProto(act.Position)
	if !ok {
		m.sendActionError(p, correlationID, "invalid position")
		return
	}
	part, ok := particleFromProto(act)
	if !ok {
		m.sendActionError(p, correlationID, "unknown particle")
		return
	}
	m.completeWorldTask(p, correlationID, w.Do(func(tx *world.Tx) {
		tx.AddParticle(pos, part)
	}))
}

func (m *Manager) handleWorldSetTime(p *pluginProcess, correlationID string, act *pb.WorldSetTimeAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	w.SetTime(int(act.Time))
	m.sendActionOK(p, correlationID)
}

func (m *Manager) handleWorldStopTime(p *pluginProcess, correlationID string, act *pb.WorldStopTimeAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	w.StopTime()
	m.sendActionOK(p, correlationID)
}

func (m *Manager) handleWorldStartTime(p *pluginProcess, correlationID string, act *pb.WorldStartTimeAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	w.StartTime()
	m.sendActionOK(p, correlationID)
}

func (m *Manager) handleWorldSetSpawn(p *pluginProcess, correlationID string, act *pb.WorldSetSpawnAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	if act.Spawn == nil {
		m.sendActionError(p, correlationID, "missing spawn position")
		return
	}
	pos := cube.Pos{int(act.Spawn.X), int(act.Spawn.Y), int(act.Spawn.Z)}
	w.SetSpawn(pos)
	m.sendActionOK(p, correlationID)
}

func (m *Manager) handleWorldSetBiome(p *pluginProcess, correlationID string, act *pb.WorldSetBiomeAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	if act.Position == nil {
		m.sendActionError(p, correlationID, "missing position")
		return
	}
	if act.BiomeId == "" {
		m.sendActionError(p, correlationID, "missing biome_id")
		return
	}
	// Parse biome ID - can be numeric ID or canonical biome name
	var biome world.Biome
	var biomeID int
	if _, err := fmt.Sscanf(act.BiomeId, "%d", &biomeID); err == nil {
		var ok bool
		biome, ok = world.BiomeByID(biomeID)
		if !ok {
			m.sendActionError(p, correlationID, "unknown biome ID")
			return
		}
	} else {
		var ok bool
		biome, ok = world.BiomeByName(act.BiomeId)
		if !ok {
			m.sendActionError(p, correlationID, "unknown biome name")
			return
		}
	}
	pos := cube.Pos{int(act.Position.X), int(act.Position.Y), int(act.Position.Z)}
	m.completeWorldTask(p, correlationID, w.Do(func(tx *world.Tx) {
		tx.SetBiome(pos, biome)
	}))
}

func (m *Manager) handleWorldSetLiquid(p *pluginProcess, correlationID string, act *pb.WorldSetLiquidAction) {
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
	var liquid world.Liquid
	if act.Liquid != nil && act.Liquid.Block != nil {
		if blk, ok := blockFromProto(act.Liquid.Block); ok {
			if liq, ok := blk.(world.Liquid); ok {
				liquid = liq
			} else {
				m.sendActionError(p, correlationID, "block is not a liquid")
				return
			}
		} else {
			m.sendActionError(p, correlationID, "unknown liquid block")
			return
		}
	}
	m.completeWorldTask(p, correlationID, w.Do(func(tx *world.Tx) {
		tx.SetLiquid(pos, liquid)
	}))
}

func (m *Manager) handleWorldScheduleBlockUpdate(p *pluginProcess, correlationID string, act *pb.WorldScheduleBlockUpdateAction) {
	w := m.worldFromRef(act.GetWorld())
	if w == nil {
		m.sendActionError(p, correlationID, "world not found")
		return
	}
	if act.Position == nil {
		m.sendActionError(p, correlationID, "missing position")
		return
	}
	if act.Block == nil {
		m.sendActionError(p, correlationID, "missing block")
		return
	}
	blk, ok := blockFromProto(act.Block)
	if !ok {
		m.sendActionError(p, correlationID, "unknown block")
		return
	}
	pos := cube.Pos{int(act.Position.X), int(act.Position.Y), int(act.Position.Z)}
	delay := time.Duration(act.DelayMs) * time.Millisecond
	m.completeWorldTask(p, correlationID, w.Do(func(tx *world.Tx) {
		tx.ScheduleBlockUpdate(pos, blk, delay)
	}))
}
