package plugin

import (
	"fmt"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	pb "github.com/secmc/plugin/proto/generated"
)

func (m *Manager) EmitWorldLiquidFlow(ctx *world.Context, from, into cube.Pos, liquid world.Liquid, replaced world.Block) {
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_WORLD_LIQUID_FLOW,
		Payload: &pb.EventEnvelope_WorldLiquidFlow{
			WorldLiquidFlow: &pb.WorldLiquidFlowEvent{
				World:    protoWorldRef(worldFromContext(ctx)),
				From:     protoBlockPos(from),
				To:       protoBlockPos(into),
				Liquid:   protoLiquidState(liquid),
				Replaced: protoBlockState(replaced),
			},
		},
	}
	m.emitCancellable(ctx, evt)
}

func (m *Manager) EmitWorldLiquidDecay(ctx *world.Context, pos cube.Pos, before, after world.Liquid) {
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_WORLD_LIQUID_DECAY,
		Payload: &pb.EventEnvelope_WorldLiquidDecay{
			WorldLiquidDecay: &pb.WorldLiquidDecayEvent{
				World:    protoWorldRef(worldFromContext(ctx)),
				Position: protoBlockPos(pos),
				Before:   protoLiquidState(before),
				After:    protoLiquidState(after),
			},
		},
	}
	m.emitCancellable(ctx, evt)
}

func (m *Manager) EmitWorldLiquidHarden(ctx *world.Context, pos cube.Pos, liquidHardened, otherLiquid, newBlock world.Block) {
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_WORLD_LIQUID_HARDEN,
		Payload: &pb.EventEnvelope_WorldLiquidHarden{
			WorldLiquidHarden: &pb.WorldLiquidHardenEvent{
				World:          protoWorldRef(worldFromContext(ctx)),
				Position:       protoBlockPos(pos),
				LiquidHardened: protoLiquidOrBlockState(liquidHardened),
				OtherLiquid:    protoLiquidOrBlockState(otherLiquid),
				NewBlock:       protoBlockState(newBlock),
			},
		},
	}
	m.emitCancellable(ctx, evt)
}

func (m *Manager) EmitWorldSound(ctx *world.Context, s world.Sound, pos mgl64.Vec3) {
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_WORLD_SOUND,
		Payload: &pb.EventEnvelope_WorldSound{
			WorldSound: &pb.WorldSoundEvent{
				World:    protoWorldRef(worldFromContext(ctx)),
				Sound:    fmt.Sprintf("%T", s),
				Position: protoVec3(pos),
			},
		},
	}
	m.emitCancellable(ctx, evt)
}

func (m *Manager) EmitWorldFireSpread(ctx *world.Context, from, to cube.Pos) {
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_WORLD_FIRE_SPREAD,
		Payload: &pb.EventEnvelope_WorldFireSpread{
			WorldFireSpread: &pb.WorldFireSpreadEvent{
				World: protoWorldRef(worldFromContext(ctx)),
				From:  protoBlockPos(from),
				To:    protoBlockPos(to),
			},
		},
	}
	m.emitCancellable(ctx, evt)
}

func (m *Manager) EmitWorldBlockBurn(ctx *world.Context, pos cube.Pos) {
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_WORLD_BLOCK_BURN,
		Payload: &pb.EventEnvelope_WorldBlockBurn{
			WorldBlockBurn: &pb.WorldBlockBurnEvent{
				World:    protoWorldRef(worldFromContext(ctx)),
				Position: protoBlockPos(pos),
			},
		},
	}
	m.emitCancellable(ctx, evt)
}

func (m *Manager) EmitWorldCropTrample(ctx *world.Context, pos cube.Pos) {
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_WORLD_CROP_TRAMPLE,
		Payload: &pb.EventEnvelope_WorldCropTrample{
			WorldCropTrample: &pb.WorldCropTrampleEvent{
				World:    protoWorldRef(worldFromContext(ctx)),
				Position: protoBlockPos(pos),
			},
		},
	}
	m.emitCancellable(ctx, evt)
}

func (m *Manager) EmitWorldLeavesDecay(ctx *world.Context, pos cube.Pos) {
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_WORLD_LEAVES_DECAY,
		Payload: &pb.EventEnvelope_WorldLeavesDecay{
			WorldLeavesDecay: &pb.WorldLeavesDecayEvent{
				World:    protoWorldRef(worldFromContext(ctx)),
				Position: protoBlockPos(pos),
			},
		},
	}
	m.emitCancellable(ctx, evt)
}

func (m *Manager) EmitWorldEntitySpawn(tx *world.Tx, e world.Entity) {
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_WORLD_ENTITY_SPAWN,
		Payload: &pb.EventEnvelope_WorldEntitySpawn{
			WorldEntitySpawn: &pb.WorldEntitySpawnEvent{
				World:  protoWorldRef(worldFromTx(tx)),
				Entity: protoEntityRef(e),
			},
		},
	}
	m.broadcastEvent(evt)
}

func (m *Manager) EmitWorldEntityDespawn(tx *world.Tx, e world.Entity) {
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_WORLD_ENTITY_DESPAWN,
		Payload: &pb.EventEnvelope_WorldEntityDespawn{
			WorldEntityDespawn: &pb.WorldEntityDespawnEvent{
				World:  protoWorldRef(worldFromTx(tx)),
				Entity: protoEntityRef(e),
			},
		},
	}
	m.broadcastEvent(evt)
}

func (m *Manager) EmitWorldExplosion(ctx *world.Context, position mgl64.Vec3, entities *[]world.Entity, blocks *[]cube.Pos, itemDropChance *float64, spawnFire *bool) {
	var entityRefs []*pb.EntityRef
	if entities != nil {
		entityRefs = protoEntityRefs(*entities)
	}
	var blockPositions []*pb.BlockPos
	if blocks != nil {
		blockPositions = protoBlockPositions(*blocks)
	}
	dropChance := 0.0
	if itemDropChance != nil {
		dropChance = *itemDropChance
	}
	spawnFireVal := false
	if spawnFire != nil {
		spawnFireVal = *spawnFire
	}
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_WORLD_EXPLOSION,
		Payload: &pb.EventEnvelope_WorldExplosion{
			WorldExplosion: &pb.WorldExplosionEvent{
				World:            protoWorldRef(worldFromContext(ctx)),
				Position:         protoVec3(position),
				AffectedEntities: entityRefs,
				AffectedBlocks:   blockPositions,
				ItemDropChance:   dropChance,
				SpawnFire:        spawnFireVal,
			},
		},
	}
	m.emitCancellable(ctx, evt)
}

func (m *Manager) EmitWorldClose(tx *world.Tx) {
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_WORLD_CLOSE,
		Payload: &pb.EventEnvelope_WorldClose{
			WorldClose: &pb.WorldCloseEvent{
				World: protoWorldRef(worldFromTx(tx)),
			},
		},
	}
	m.broadcastEvent(evt)
}

func worldFromTx(tx *world.Tx) *world.World {
	if tx == nil {
		return nil
	}
	return tx.World()
}
