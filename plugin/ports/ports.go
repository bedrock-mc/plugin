package ports

import (
	"context"
	"net"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	pb "github.com/secmc/plugin/proto/generated/go"
)

type PluginManager interface {
	Start(configPath string) error
	Close()
}

type PluginProcess interface {
	Start(ctx context.Context) error
	Stop()
	HasSubscription(eventType pb.EventType) bool
	Queue(msg *pb.HostToPlugin)
}

type Stream interface {
	Send(data []byte) error
	Recv() ([]byte, error)
	CloseSend() error
	Close() error
}

type EventManager interface {
	EmitPlayerJoin(p *player.Player)
	EmitPlayerQuit(p *player.Player)
	EmitChat(ctx *player.Context, p *player.Player, msg *string)
	EmitCommand(ctx *player.Context, p *player.Player, cmdName string, args []string)
	EmitBlockBreak(ctx *player.Context, p *player.Player, pos cube.Pos, drops *[]item.Stack, xp *int, worldDim string)
	EmitPlayerMove(ctx *player.Context, p *player.Player, newPos mgl64.Vec3, newRot cube.Rotation)
	EmitPlayerJump(p *player.Player)
	EmitPlayerTeleport(ctx *player.Context, p *player.Player, pos mgl64.Vec3)
	EmitPlayerChangeWorld(p *player.Player, before, after *world.World)
	EmitPlayerToggleSprint(ctx *player.Context, p *player.Player, after bool)
	EmitPlayerToggleSneak(ctx *player.Context, p *player.Player, after bool)
	EmitPlayerFoodLoss(ctx *player.Context, p *player.Player, from int, to *int)
	EmitPlayerHeal(ctx *player.Context, p *player.Player, health *float64, src world.HealingSource)
	EmitPlayerHurt(ctx *player.Context, p *player.Player, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource)
	EmitPlayerDeath(p *player.Player, src world.DamageSource, keepInv *bool)
	EmitPlayerRespawn(p *player.Player, pos *mgl64.Vec3, w **world.World)
	EmitPlayerSkinChange(ctx *player.Context, p *player.Player, sk *skin.Skin)
	EmitPlayerFireExtinguish(ctx *player.Context, p *player.Player, pos cube.Pos)
	EmitPlayerStartBreak(ctx *player.Context, p *player.Player, pos cube.Pos)
	EmitPlayerBlockPlace(ctx *player.Context, p *player.Player, pos cube.Pos, b world.Block)
	EmitPlayerBlockPick(ctx *player.Context, p *player.Player, pos cube.Pos, b world.Block)
	EmitPlayerItemUse(ctx *player.Context, p *player.Player)
	EmitPlayerItemUseOnBlock(ctx *player.Context, p *player.Player, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3, b world.Block)
	EmitPlayerItemUseOnEntity(ctx *player.Context, p *player.Player, target world.Entity)
	EmitPlayerItemRelease(ctx *player.Context, p *player.Player, it item.Stack, dur time.Duration)
	EmitPlayerItemConsume(ctx *player.Context, p *player.Player, it item.Stack)
	EmitPlayerAttackEntity(ctx *player.Context, p *player.Player, target world.Entity, force, height *float64, critical *bool)
	EmitPlayerExperienceGain(ctx *player.Context, p *player.Player, amount *int)
	EmitPlayerPunchAir(ctx *player.Context, p *player.Player)
	EmitPlayerSignEdit(ctx *player.Context, p *player.Player, pos cube.Pos, frontSide bool, oldText, newText string)
	EmitPlayerLecternPageTurn(ctx *player.Context, p *player.Player, pos cube.Pos, oldPage int, newPage *int)
	EmitPlayerItemDamage(ctx *player.Context, p *player.Player, it item.Stack, damage int)
	EmitPlayerItemPickup(ctx *player.Context, p *player.Player, it *item.Stack)
	EmitPlayerHeldSlotChange(ctx *player.Context, p *player.Player, from, to int)
	EmitPlayerItemDrop(ctx *player.Context, p *player.Player, it item.Stack)
	EmitPlayerTransfer(ctx *player.Context, p *player.Player, addr *net.UDPAddr)
	EmitPlayerDiagnostics(p *player.Player, d session.Diagnostics)
	EmitWorldLiquidFlow(ctx *world.Context, from, into cube.Pos, liquid world.Liquid, replaced world.Block)
	EmitWorldLiquidDecay(ctx *world.Context, pos cube.Pos, before, after world.Liquid)
	EmitWorldLiquidHarden(ctx *world.Context, pos cube.Pos, liquidHardened, otherLiquid world.Block, newBlock world.Block)
	EmitWorldSound(ctx *world.Context, s world.Sound, pos mgl64.Vec3)
	EmitWorldFireSpread(ctx *world.Context, from, to cube.Pos)
	EmitWorldBlockBurn(ctx *world.Context, pos cube.Pos)
	EmitWorldCropTrample(ctx *world.Context, pos cube.Pos)
	EmitWorldLeavesDecay(ctx *world.Context, pos cube.Pos)
	EmitWorldEntitySpawn(tx *world.Tx, e world.Entity)
	EmitWorldEntityDespawn(tx *world.Tx, e world.Entity)
	EmitWorldExplosion(ctx *world.Context, position mgl64.Vec3, entities *[]world.Entity, blocks *[]cube.Pos, itemDropChance *float64, spawnFire *bool)
	EmitWorldClose(tx *world.Tx)
}

type PlayerHandlerFactory func(manager EventManager) player.Handler

type WorldHandlerFactory func(manager EventManager) world.Handler

type PluginService interface {
	PluginManager
	EventManager
	AttachWorld(w *world.World)
	AttachPlayer(p *player.Player)
}
