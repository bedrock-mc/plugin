package plugin

import (
	"strings"

	"github.com/df-mc/dragonfly/server/player"
	pb "github.com/secmc/plugin/proto/generated"
)

func (m *Manager) EmitPlayerJoin(p *player.Player) {
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_PLAYER_JOIN,
		Payload: &pb.EventEnvelope_PlayerJoin{
			PlayerJoin: &pb.PlayerJoinEvent{
				PlayerUuid: p.UUID().String(),
				Name:       p.Name(),
			},
		},
	}
	m.broadcastEvent(evt)
}

func (m *Manager) EmitPlayerQuit(p *player.Player) {
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_PLAYER_QUIT,
		Payload: &pb.EventEnvelope_PlayerQuit{
			PlayerQuit: &pb.PlayerQuitEvent{
				PlayerUuid: p.UUID().String(),
				Name:       p.Name(),
			},
		},
	}
	m.broadcastEvent(evt)
	m.detachPlayer(p)
}

func (m *Manager) EmitChat(ctx *player.Context, p *player.Player, msg *string) {
	if msg == nil {
		return
	}
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_CHAT,
		Payload: &pb.EventEnvelope_Chat{
			Chat: &pb.ChatEvent{
				PlayerUuid: p.UUID().String(),
				Name:       p.Name(),
				Message:    *msg,
			},
		},
	}
	results := m.emitCancellable(ctx, evt)
	for _, res := range results {
		if res == nil {
			continue
		}
		if chatMut := res.GetChat(); chatMut != nil {
			*msg = chatMut.Message
		}
	}
}

func (m *Manager) EmitCommand(ctx *player.Context, p *player.Player, cmdName string, args []string) {
	raw := "/" + cmdName
	if len(args) > 0 {
		raw += " " + strings.Join(args, " ")
	}
	evt := &pb.EventEnvelope{
		EventId: m.generateEventID(),
		Type:    pb.EventType_COMMAND,
		Payload: &pb.EventEnvelope_Command{
			Command: &pb.CommandEvent{
				PlayerUuid: p.UUID().String(),
				Name:       p.Name(),
				Raw:        raw,
				Command:    cmdName,
				Args:       args,
			},
		},
	}
	m.emitCancellable(ctx, evt)
}
