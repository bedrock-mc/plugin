package dragonflyhost

import (
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
		return WorldName(tx.World(), fallback)
	}
	name := ""
	if ok := p.H().ExecWorld(func(tx *world.Tx, _ world.Entity) {
		name = WorldName(tx.World(), fallback)
	}); ok && name != "" {
		return name
	}
	return fallback
}

func PlayerWorldDimension(p *player.Player) string {
	if p == nil {
		return ""
	}
	if tx := p.Tx(); tx != nil {
		return WorldDimension(tx.World())
	}
	dimension := ""
	if ok := p.H().ExecWorld(func(tx *world.Tx, _ world.Entity) {
		dimension = WorldDimension(tx.World())
	}); ok && dimension != "" {
		return dimension
	}
	return ""
}

func PlayerWorldReference(p *player.Player) *WorldRef {
	if p == nil {
		return nil
	}
	if tx := p.Tx(); tx != nil {
		return WorldReference(tx.World())
	}
	var ref *WorldRef
	if ok := p.H().ExecWorld(func(tx *world.Tx, _ world.Entity) {
		ref = WorldReference(tx.World())
	}); ok && ref != nil {
		return ref
	}
	return nil
}
