package plugin

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	pb "github.com/secmc/plugin/proto/generated"
)

func (m *Manager) EmitBlockBreak(ctx *player.Context, p *player.Player, pos cube.Pos, drops *[]item.Stack, xp *int, worldDim string) {
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_PLAYER_BLOCK_BREAK,
		Payload: &pb.EventEnvelope_BlockBreak{
			BlockBreak: &pb.BlockBreakEvent{
				PlayerUuid: p.UUID().String(),
				Name:       p.Name(),
				World:      worldDim,
				Position:   protoBlockPos(pos),
			},
		},
	}
	results := m.emitCancellable(ctx, evt)
	for _, res := range results {
		if res == nil {
			continue
		}
		if bbMut := res.GetBlockBreak(); bbMut != nil {
			if drops != nil {
				*drops = convertProtoDrops(bbMut.Drops)
			}
			if bbMut.Xp != nil && xp != nil {
				*xp = int(*bbMut.Xp)
			}
		}
	}
}
