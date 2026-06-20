package dragonflyhost

import (
	"strings"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

func Vec3Position(pos mgl64.Vec3, worldName string) Position {
	return Position{X: pos.X(), Y: pos.Y(), Z: pos.Z(), World: worldName}
}

func BlockPosition(pos cube.Pos, worldName string) Position {
	return Position{X: float64(pos.X()), Y: float64(pos.Y()), Z: float64(pos.Z()), World: worldName}
}

func PlayerPosition(p *player.Player, pos mgl64.Vec3, fallbackWorld string) Position {
	return Vec3Position(pos, PlayerWorldName(p, fallbackWorld))
}

func PlayerBlockPosition(p *player.Player, pos cube.Pos, fallbackWorld string) Position {
	return BlockPosition(pos, PlayerWorldName(p, fallbackWorld))
}

func PlayerStateSnapshot(p *player.Player, fallbackWorld string) PlayerState {
	pos := p.Position()
	return PlayerState{
		Position:   PlayerPosition(p, pos, fallbackWorld),
		Health:     p.Health(),
		MaxHealth:  p.MaxHealth(),
		Gamemode:   GameModeName(p.GameMode()),
		XPLevel:    p.ExperienceLevel(),
		XPProgress: p.ExperienceProgress(),
	}
}

func InventorySlots(p *player.Player, size int) []InventorySlot {
	if p == nil || size <= 0 {
		return nil
	}
	inv := p.Inventory()
	slots := make([]InventorySlot, 0, size)
	for slot := 0; slot < size; slot++ {
		stack, err := inv.Item(slot)
		if err != nil {
			continue
		}
		slots = append(slots, InventorySlot{Slot: slot, Item: ItemStackSnapshot(stack)})
	}
	return slots
}

func ItemStackSnapshot(stack item.Stack) ItemStack {
	if stack.Empty() {
		return ItemStack{TypeID: "minecraft:air", Name: "Air", Count: 0}
	}
	name, meta := stack.Item().EncodeItem()
	return ItemStack{TypeID: name, Name: DisplayName(name), Meta: meta, Count: stack.Count()}
}

func DisplayName(typeID string) string {
	name := strings.TrimPrefix(typeID, "minecraft:")
	name = strings.ReplaceAll(name, "_", " ")
	if name == "" {
		return typeID
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

func GameModeName(mode world.GameMode) string {
	id, ok := world.GameModeID(mode)
	if !ok {
		return "survival"
	}
	switch id {
	case 1:
		return "creative"
	case 2:
		return "adventure"
	case 3:
		return "spectator"
	default:
		return "survival"
	}
}
