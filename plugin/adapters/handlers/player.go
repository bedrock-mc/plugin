package handlers

import (
	"fmt"
	"net"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/secmc/plugin/plugin/ports"
)

type PlayerHandler struct {
	player.NopHandler
	manager ports.EventManager
}

func NewPlayerHandler(manager ports.EventManager) player.Handler {
	return &PlayerHandler{manager: manager}
}

func (h *PlayerHandler) HandleChat(ctx *player.Context, message *string) {
	if h.manager == nil {
		return
	}
	h.manager.EmitChat(ctx, ctx.Val(), message)
}

func (h *PlayerHandler) HandleMove(ctx *player.Context, newPos mgl64.Vec3, newRot cube.Rotation) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerMove(ctx, ctx.Val(), newPos, newRot)
}

func (h *PlayerHandler) HandleJump(p *player.Player) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerJump(p)
}

func (h *PlayerHandler) HandleTeleport(ctx *player.Context, pos mgl64.Vec3) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerTeleport(ctx, ctx.Val(), pos)
}

func (h *PlayerHandler) HandleChangeWorld(p *player.Player, before, after *world.World) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerChangeWorld(p, before, after)
}

func (h *PlayerHandler) HandleToggleSprint(ctx *player.Context, after bool) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerToggleSprint(ctx, ctx.Val(), after)
}

func (h *PlayerHandler) HandleToggleSneak(ctx *player.Context, after bool) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerToggleSneak(ctx, ctx.Val(), after)
}

func (h *PlayerHandler) HandleFoodLoss(ctx *player.Context, from int, to *int) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerFoodLoss(ctx, ctx.Val(), from, to)
}

func (h *PlayerHandler) HandleHeal(ctx *player.Context, health *float64, src world.HealingSource) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerHeal(ctx, ctx.Val(), health, src)
}

func (h *PlayerHandler) HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerHurt(ctx, ctx.Val(), damage, immune, attackImmunity, src)
}

func (h *PlayerHandler) HandleDeath(p *player.Player, src world.DamageSource, keepInv *bool) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerDeath(p, src, keepInv)
}

func (h *PlayerHandler) HandleRespawn(p *player.Player, pos *mgl64.Vec3, w **world.World) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerRespawn(p, pos, w)
}

func (h *PlayerHandler) HandleSkinChange(ctx *player.Context, skin *skin.Skin) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerSkinChange(ctx, ctx.Val(), skin)
}

func (h *PlayerHandler) HandleFireExtinguish(ctx *player.Context, pos cube.Pos) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerFireExtinguish(ctx, ctx.Val(), pos)
}

func (h *PlayerHandler) HandleStartBreak(ctx *player.Context, pos cube.Pos) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerStartBreak(ctx, ctx.Val(), pos)
}

func (h *PlayerHandler) HandleCommandExecution(ctx *player.Context, command cmd.Command, args []string) {
	if h.manager == nil {
		return
	}
	h.manager.EmitCommand(ctx, ctx.Val(), command.Name(), args)
}

func (h *PlayerHandler) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	if h.manager == nil {
		return
	}
	p := ctx.Val()
	worldDim := fmt.Sprint(p.Tx().World().Dimension())
	h.manager.EmitBlockBreak(ctx, p, pos, drops, xp, worldDim)
}

func (h *PlayerHandler) HandleQuit(p *player.Player) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerQuit(p)
}

func (h *PlayerHandler) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerBlockPlace(ctx, ctx.Val(), pos, b)
}

func (h *PlayerHandler) HandleBlockPick(ctx *player.Context, pos cube.Pos, b world.Block) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerBlockPick(ctx, ctx.Val(), pos, b)
}

func (h *PlayerHandler) HandleItemUse(ctx *player.Context) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerItemUse(ctx, ctx.Val())
}

func (h *PlayerHandler) HandleItemUseOnBlock(ctx *player.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	if h.manager == nil {
		return
	}
	p := ctx.Val()
	if p == nil {
		return
	}
	var block world.Block
	if tx := p.Tx(); tx != nil {
		block = tx.Block(pos)
	}
	h.manager.EmitPlayerItemUseOnBlock(ctx, p, pos, face, clickPos, block)
}

func (h *PlayerHandler) HandleItemUseOnEntity(ctx *player.Context, e world.Entity) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerItemUseOnEntity(ctx, ctx.Val(), e)
}

func (h *PlayerHandler) HandleItemRelease(ctx *player.Context, it item.Stack, dur time.Duration) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerItemRelease(ctx, ctx.Val(), it, dur)
}

func (h *PlayerHandler) HandleItemConsume(ctx *player.Context, it item.Stack) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerItemConsume(ctx, ctx.Val(), it)
}

func (h *PlayerHandler) HandleAttackEntity(ctx *player.Context, e world.Entity, force, height *float64, critical *bool) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerAttackEntity(ctx, ctx.Val(), e, force, height, critical)
}

func (h *PlayerHandler) HandleExperienceGain(ctx *player.Context, amount *int) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerExperienceGain(ctx, ctx.Val(), amount)
}

func (h *PlayerHandler) HandlePunchAir(ctx *player.Context) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerPunchAir(ctx, ctx.Val())
}

func (h *PlayerHandler) HandleSignEdit(ctx *player.Context, pos cube.Pos, frontSide bool, oldText, newText string) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerSignEdit(ctx, ctx.Val(), pos, frontSide, oldText, newText)
}

func (h *PlayerHandler) HandleLecternPageTurn(ctx *player.Context, pos cube.Pos, oldPage int, newPage *int) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerLecternPageTurn(ctx, ctx.Val(), pos, oldPage, newPage)
}

func (h *PlayerHandler) HandleItemDamage(ctx *player.Context, it item.Stack, damage int) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerItemDamage(ctx, ctx.Val(), it, damage)
}

func (h *PlayerHandler) HandleItemPickup(ctx *player.Context, it *item.Stack) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerItemPickup(ctx, ctx.Val(), it)
}

func (h *PlayerHandler) HandleHeldSlotChange(ctx *player.Context, from, to int) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerHeldSlotChange(ctx, ctx.Val(), from, to)
}

func (h *PlayerHandler) HandleItemDrop(ctx *player.Context, it item.Stack) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerItemDrop(ctx, ctx.Val(), it)
}

func (h *PlayerHandler) HandleTransfer(ctx *player.Context, addr *net.UDPAddr) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerTransfer(ctx, ctx.Val(), addr)
}

func (h *PlayerHandler) HandleDiagnostics(p *player.Player, d session.Diagnostics) {
	if h.manager == nil {
		return
	}
	h.manager.EmitPlayerDiagnostics(p, d)
}
