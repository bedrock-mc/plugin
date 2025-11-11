package plugin

import (
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	pb "github.com/secmc/plugin/proto/generated"
)

func (m *Manager) registerCommands(p *pluginProcess, specs []*pb.CommandSpec) {
	for _, spec := range specs {
		if spec == nil || spec.Name == "" {
			continue
		}
		name := strings.TrimPrefix(spec.Name, "/")

		aliases := make([]string, 0, len(spec.Aliases))
		for _, alias := range spec.Aliases {
			alias = strings.TrimPrefix(alias, "/")
			if alias == "" || alias == name {
				continue
			}
			aliases = append(aliases, alias)
		}

		binding := commandBinding{pluginID: p.id, command: name, descriptor: spec}
		m.mu.Lock()
		m.commands[name] = binding
		for _, alias := range aliases {
			m.commands[alias] = binding
		}
		m.mu.Unlock()

		cmd.Register(cmd.New(name, spec.Description, aliases, pluginCommand{mgr: m, pluginID: p.id, name: name}))
	}
}

type pluginCommand struct {
	mgr      *Manager
	pluginID string
	name     string
}

func (c pluginCommand) Run(src cmd.Source, output *cmd.Output, tx *world.Tx) {
	_, ok := src.(*player.Player)
	if !ok {
		output.Errorf("command only available to players")
		return
	}
	// No-op: PlayerHandler.HandleCommandExecution emits command events
}
