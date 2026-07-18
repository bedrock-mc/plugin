package dragonflyhost

import (
	"context"
	"fmt"
	"strings"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

func WorldName(w *world.World, fallback string) string {
	if w == nil {
		return fallback
	}
	return w.Name()
}

func WorldDimension(w *world.World) string {
	if w == nil {
		return ""
	}
	return strings.ToLower(fmt.Sprint(w.Dimension()))
}

func WorldReference(w *world.World) *WorldRef {
	if w == nil {
		return nil
	}
	return &WorldRef{
		ID:        fmt.Sprintf("%p", w),
		Name:      w.Name(),
		Dimension: WorldDimension(w),
	}
}

func PlayerWorldName(p *player.Player, fallback string) string {
	if p == nil {
		return fallback
	}
	if tx := p.Tx(); tx != nil {
		if name, ok := txWorldName(tx, fallback); ok {
			return name
		}
	}
	name, err := player.Call(context.Background(), p.H(), func(tx *world.Tx, _ *player.Player) (string, error) {
		return WorldName(tx.World(), fallback), nil
	})
	if err == nil && name != "" {
		return name
	}
	return fallback
}

func PlayerWorldDimension(p *player.Player) string {
	if p == nil {
		return ""
	}
	if tx := p.Tx(); tx != nil {
		if dimension, ok := txWorldDimension(tx); ok {
			return dimension
		}
	}
	dimension, err := player.Call(context.Background(), p.H(), func(tx *world.Tx, _ *player.Player) (string, error) {
		return WorldDimension(tx.World()), nil
	})
	if err == nil && dimension != "" {
		return dimension
	}
	return ""
}

func PlayerWorldReference(p *player.Player) *WorldRef {
	if p == nil {
		return nil
	}
	if tx := p.Tx(); tx != nil {
		if ref, ok := txWorldReference(tx); ok {
			return ref
		}
	}
	ref, err := player.Call(context.Background(), p.H(), func(tx *world.Tx, _ *player.Player) (*WorldRef, error) {
		return WorldReference(tx.World()), nil
	})
	if err == nil && ref != nil {
		return ref
	}
	return nil
}

func txWorldName(tx *world.Tx, fallback string) (name string, ok bool) {
	defer func() {
		if recover() != nil {
			name = ""
			ok = false
		}
	}()
	return WorldName(tx.World(), fallback), true
}

func txWorldDimension(tx *world.Tx) (dimension string, ok bool) {
	defer func() {
		if recover() != nil {
			dimension = ""
			ok = false
		}
	}()
	return WorldDimension(tx.World()), true
}

func txWorldReference(tx *world.Tx) (ref *WorldRef, ok bool) {
	defer func() {
		if recover() != nil {
			ref = nil
			ok = false
		}
	}()
	return WorldReference(tx.World()), true
}
