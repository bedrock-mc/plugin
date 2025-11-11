package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/df-mc/dragonfly/server"
	"github.com/sandertv/gophertunnel/minecraft"
)

// listenerFunc ...
func listenerFunc(c *server.Config, addr string, protocols []minecraft.Protocol) {
	c.Listeners = []func(conf server.Config) (server.Listener, error){
		func(conf server.Config) (server.Listener, error) {
			cfg := minecraft.ListenConfig{
				MaximumPlayers:         conf.MaxPlayers,
				StatusProvider:         conf.StatusProvider,
				AuthenticationDisabled: conf.AuthDisabled,
				ResourcePacks:          conf.Resources,
				TexturePacksRequired:   conf.ResourcesRequired,

				AcceptedProtocols:   protocols,
				AllowUnknownPackets: true,
				AllowInvalidPackets: true,
			}
			if conf.Log.Enabled(context.Background(), slog.LevelDebug) {
				cfg.ErrorLog = conf.Log.With("net origin", "gophertunnel")
			}
			l, err := cfg.Listen("raknet", addr)
			if err != nil {
				return nil, fmt.Errorf("create minecraft listener: %w", err)
			}
			conf.Log.Info("Listener running.", "addr", l.Addr())
			return listener{l}, nil
		},
	}
}

// single ...
func single[Y any](x Y) []Y {
	return []Y{x}
}
