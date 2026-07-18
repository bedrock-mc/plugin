package plugin

import (
	pb "github.com/bedrock-mc/plugin/proto/generated/go"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

func (m *Manager) execMethod(id uuid.UUID, method func(pl *player.Player)) {
	if m.srv == nil {
		return
	}
	if handle, ok := m.srv.Player(id); ok {
		task := player.Do(handle, func(_ *world.Tx, pl *player.Player) {
			method(pl)
		})
		task.OnDone(func(err error) {
			if err != nil && m.log != nil {
				m.log.Warn("player action failed", "player_uuid", id, "error", err)
			}
		})
	}
}

func (m *Manager) sendActionResult(p *pluginProcess, result *pb.ActionResult) {
	if p == nil || result == nil || result.CorrelationId == "" {
		return
	}
	p.queue(&pb.HostToPlugin{
		PluginId: p.id,
		Payload:  &pb.HostToPlugin_ActionResult{ActionResult: result},
	})
}

func (m *Manager) sendActionOK(p *pluginProcess, correlationID string) {
	if correlationID == "" {
		return
	}
	m.sendActionResult(p, &pb.ActionResult{CorrelationId: correlationID, Status: &pb.ActionStatus{Ok: true}})
}

func (m *Manager) sendActionError(p *pluginProcess, correlationID, msg string) {
	if correlationID == "" {
		return
	}
	m.sendActionResult(p, &pb.ActionResult{CorrelationId: correlationID, Status: &pb.ActionStatus{Ok: false, Error: &msg}})
}

func (m *Manager) completeWorldTask(p *pluginProcess, correlationID string, task *world.Task) {
	task.OnDone(func(err error) {
		if err != nil {
			m.sendActionError(p, correlationID, err.Error())
			return
		}
		m.sendActionOK(p, correlationID)
	})
}
