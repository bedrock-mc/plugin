package plugin

import (
	"strings"
	"time"

	pb "github.com/bedrock-mc/plugin/proto/generated/go"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/bossbar"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/dialogue"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/player/hud"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
)

func (m *Manager) handleSendChat(act *pb.SendChatAction) {
	if act.TargetUuid == "" {
		for p := range m.srv.Players(nil) {
			p.Message(act.Message)
		}
		chat.Global.WriteString(act.Message)
		return
	}
	id, err := uuid.Parse(act.TargetUuid)
	if err != nil {
		return
	}

	m.execMethod(id, func(pl *player.Player) {
		pl.Message(act.Message)
	})
}

func (m *Manager) handleTeleport(act *pb.TeleportAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}

	m.execMethod(id, func(pl *player.Player) {
		pos, ok := vec3FromProto(act.Position)
		if ok {
			pl.Teleport(pos)
		}
		rot, ok := vec3FromProto(act.Rotation)
		if ok {
			playerRot := pl.Rotation()
			deltaYaw := rot[1] - playerRot.Yaw()
			deltaPitch := rot[0] - playerRot.Pitch()
			if deltaYaw != 0 || deltaPitch != 0 {
				pl.Move(mgl64.Vec3{}, deltaYaw, deltaPitch)
			}
		}
	})
}

func (m *Manager) handleKick(act *pb.KickAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		pl.Disconnect(act.Reason)
	})
}

func (m *Manager) handleSetGameMode(act *pb.SetGameModeAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	gameMode, ok := world.GameModeByID(int(act.GameMode))
	if !ok {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		pl.SetGameMode(gameMode)
	})
}

func (m *Manager) handleGiveItem(act *pb.GiveItemAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		if stack, ok := convertProtoItemStackValue(act.Item); ok {
			_, _ = pl.Inventory().AddItem(stack)
		}
	})
}

func (m *Manager) handleClearInventory(act *pb.ClearInventoryAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		_ = pl.Inventory().Clear()
		pl.SetHeldItems(item.Stack{}, item.Stack{})
	})
}

func (m *Manager) handleSetHeldItem(act *pb.SetHeldItemAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		main, off := pl.HeldItems()

		if act.Main != nil {
			if s, ok := convertProtoItemStackValue(act.Main); ok {
				main = s
			}
		}
		if act.Offhand != nil {
			if s, ok := convertProtoItemStackValue(act.Offhand); ok {
				off = s
			}
		}
		pl.SetHeldItems(main, off)
	})
}

func (m *Manager) handleSetHealth(act *pb.SetHealthAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		if act.MaxHealth != nil {
			pl.SetMaxHealth(*act.MaxHealth)
		}
		current := pl.Health()
		target := act.Health
		if target > current {
			pl.Heal(target-current, entity.FoodHealingSource{})
		} else if target < current {
			pl.Hurt(current-target, entity.VoidDamageSource{})
		}
	})
}

func (m *Manager) handleSetFood(act *pb.SetFoodAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		pl.SetFood(int(act.Food))
	})
}

func (m *Manager) handleSetExperience(act *pb.SetExperienceAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		if act.Level != nil {
			pl.SetExperienceLevel(int(*act.Level))
		}
		if act.Progress != nil {
			pl.SetExperienceProgress(float64(*act.Progress))
		}
		if act.Amount != nil {
			amt := int(*act.Amount)
			if amt >= 0 {
				_ = pl.AddExperience(amt)
			} else {
				pl.RemoveExperience(-amt)
			}
		}
	})
}

func (m *Manager) handleSetVelocity(act *pb.SetVelocityAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		if v, ok := vec3FromProto(act.Velocity); ok {
			pl.SetVelocity(v)
		}
	})
}

func (m *Manager) handleAddEffect(act *pb.AddEffectAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		t, ok := effect.ByID(int(act.EffectType))
		if !ok {
			return
		}
		var e effect.Effect
		if lt, ok := t.(effect.LastingType); ok {
			d := time.Duration(act.DurationMs) * time.Millisecond
			if d <= 0 {
				e = effect.NewInfinite(lt, int(act.Level))
			} else {
				e = effect.New(lt, int(act.Level), d)
			}
		} else {
			e = effect.NewInstantWithPotency(t, int(act.Level), 1.0)
		}
		if !act.ShowParticles {
			e = e.WithoutParticles()
		}
		pl.AddEffect(e)
	})
}

func (m *Manager) handleRemoveEffect(act *pb.RemoveEffectAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		t, ok := effect.ByID(int(act.EffectType))
		if !ok {
			return
		}
		pl.RemoveEffect(t)
	})
}

func (m *Manager) handleSendTitle(act *pb.SendTitleAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		t := playerTitleFromAction(act)
		pl.SendTitle(t)
	})
}

func (m *Manager) handleSendPopup(act *pb.SendPopupAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		pl.SendPopup(act.Message)
	})
}

func (m *Manager) handleSendTip(act *pb.SendTipAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		pl.SendTip(act.Message)
	})
}

func (m *Manager) handlePlaySound(act *pb.PlaySoundAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		s := soundFromProto(act.Sound)
		pl.PlaySound(s)
	})
}

func (m *Manager) handleExecuteCommand(act *pb.ExecuteCommandAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		cmd := act.Command
		if cmd != "" && !strings.HasPrefix(cmd, "/") {
			cmd = "/" + cmd
		}
		pl.ExecuteCommand(cmd)
	})
}

func (m *Manager) handlePlayerStartSprinting(act *pb.PlayerStartSprintingAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.StartSprinting() })
}
func (m *Manager) handlePlayerStopSprinting(act *pb.PlayerStopSprintingAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.StopSprinting() })
}
func (m *Manager) handlePlayerStartSneaking(act *pb.PlayerStartSneakingAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.StartSneaking() })
}
func (m *Manager) handlePlayerStopSneaking(act *pb.PlayerStopSneakingAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.StopSneaking() })
}
func (m *Manager) handlePlayerStartSwimming(act *pb.PlayerStartSwimmingAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.StartSwimming() })
}
func (m *Manager) handlePlayerStopSwimming(act *pb.PlayerStopSwimmingAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.StopSwimming() })
}
func (m *Manager) handlePlayerStartCrawling(act *pb.PlayerStartCrawlingAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.StartCrawling() })
}
func (m *Manager) handlePlayerStopCrawling(act *pb.PlayerStopCrawlingAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.StopCrawling() })
}
func (m *Manager) handlePlayerStartGliding(act *pb.PlayerStartGlidingAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.StartGliding() })
}
func (m *Manager) handlePlayerStopGliding(act *pb.PlayerStopGlidingAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.StopGliding() })
}
func (m *Manager) handlePlayerStartFlying(act *pb.PlayerStartFlyingAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.StartFlying() })
}
func (m *Manager) handlePlayerStopFlying(act *pb.PlayerStopFlyingAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.StopFlying() })
}

// Player mobility lock
func (m *Manager) handlePlayerSetImmobile(act *pb.PlayerSetImmobileAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.SetImmobile() })
}
func (m *Manager) handlePlayerSetMobile(act *pb.PlayerSetMobileAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.SetMobile() })
}

// Player movement attributes
func (m *Manager) handlePlayerSetSpeed(act *pb.PlayerSetSpeedAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.SetSpeed(act.Speed) })
}
func (m *Manager) handlePlayerSetFlightSpeed(act *pb.PlayerSetFlightSpeedAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.SetFlightSpeed(act.FlightSpeed) })
}
func (m *Manager) handlePlayerSetVerticalFlightSpeed(act *pb.PlayerSetVerticalFlightSpeedAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.SetVerticalFlightSpeed(act.VerticalFlightSpeed) })
}

// Player health/status
func (m *Manager) handlePlayerSetAbsorption(act *pb.PlayerSetAbsorptionAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.SetAbsorption(act.Absorption) })
}
func (m *Manager) handlePlayerSetOnFire(act *pb.PlayerSetOnFireAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	d := time.Duration(act.DurationMs) * time.Millisecond
	m.execMethod(id, func(pl *player.Player) { pl.SetOnFire(d) })
}
func (m *Manager) handlePlayerExtinguish(act *pb.PlayerExtinguishAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.Extinguish() })
}
func (m *Manager) handlePlayerSetInvisible(act *pb.PlayerSetInvisibleAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.SetInvisible() })
}
func (m *Manager) handlePlayerSetVisible(act *pb.PlayerSetVisibleAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.SetVisible() })
}

// Player misc attributes
func (m *Manager) handlePlayerSetScale(act *pb.PlayerSetScaleAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.SetScale(act.Scale) })
}
func (m *Manager) handlePlayerSetHeldSlot(act *pb.PlayerSetHeldSlotAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	slot := int(act.Slot)
	m.execMethod(id, func(pl *player.Player) { _ = pl.SetHeldSlot(slot) })
}

// Player UI
func (m *Manager) handlePlayerSendToast(act *pb.PlayerSendToastAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	titleText := act.Title
	message := act.Message
	m.execMethod(id, func(pl *player.Player) { pl.SendToast(titleText, message) })
}
func (m *Manager) handlePlayerSendJukeboxPopup(act *pb.PlayerSendJukeboxPopupAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	msg := act.Message
	m.execMethod(id, func(pl *player.Player) { pl.SendJukeboxPopup(msg) })
}
func (m *Manager) handlePlayerShowCoordinates(act *pb.PlayerShowCoordinatesAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.ShowCoordinates() })
}
func (m *Manager) handlePlayerHideCoordinates(act *pb.PlayerHideCoordinatesAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.HideCoordinates() })
}
func (m *Manager) handlePlayerEnableInstantRespawn(act *pb.PlayerEnableInstantRespawnAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.EnableInstantRespawn() })
}
func (m *Manager) handlePlayerDisableInstantRespawn(act *pb.PlayerDisableInstantRespawnAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.DisableInstantRespawn() })
}
func (m *Manager) handlePlayerSetNameTag(act *pb.PlayerSetNameTagAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	name := act.NameTag
	m.execMethod(id, func(pl *player.Player) { pl.SetNameTag(name) })
}
func (m *Manager) handlePlayerSetScoreTag(act *pb.PlayerSetScoreTagAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	text := act.ScoreTag
	m.execMethod(id, func(pl *player.Player) { pl.SetScoreTag(text) })
}

// Player visuals
func (m *Manager) handlePlayerShowParticle(act *pb.PlayerShowParticleAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	pos, ok := vec3FromProto(act.Position)
	if !ok {
		return
	}
	part, ok := particleFromPlayerAction(act)
	if !ok {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.ShowParticle(pos, part) })
}

// Player lifecycle/control
func (m *Manager) handlePlayerRespawn(act *pb.PlayerRespawnAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { _ = pl.Respawn() })
}
func (m *Manager) handlePlayerTransferAction(act *pb.PlayerTransferAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	addr := parseProtoAddress(act.Address)
	if addr == nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { _ = pl.Transfer(addr.String()) })
}
func (m *Manager) handlePlayerKnockBack(act *pb.PlayerKnockBackAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	src, ok := vec3FromProto(act.Source)
	if !ok {
		return
	}
	force := act.Force
	height := act.Height
	m.execMethod(id, func(pl *player.Player) { pl.KnockBack(src, force, height) })
}
func (m *Manager) handlePlayerSwingArm(act *pb.PlayerSwingArmAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.SwingArm() })
}
func (m *Manager) handlePlayerPunchAirAction(act *pb.PlayerPunchAirAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.PunchAir() })
}

// Player boss bar
func (m *Manager) handlePlayerSendBossBar(act *pb.PlayerSendBossBarAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		bar := bossbar.New(act.Text)
		if act.HealthPercentage != nil {
			h := float64(*act.HealthPercentage)
			if h < 0 {
				h = 0
			}
			if h > 1 {
				h = 1
			}
			bar = bar.WithHealthPercentage(h)
		}
		if act.Colour != nil {
			bar = bar.WithColour(convertBossBarColour(*act.Colour))
		}
		pl.SendBossBar(bar)
	})
}

func (m *Manager) handlePlayerRemoveBossBar(act *pb.PlayerRemoveBossBarAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.RemoveBossBar() })
}

// Player HUD
func (m *Manager) handlePlayerShowHudElement(act *pb.PlayerShowHudElementAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		if el, ok := convertHudElement(act.Element); ok {
			pl.ShowHudElement(el)
		}
	})
}

func (m *Manager) handlePlayerHideHudElement(act *pb.PlayerHideHudElementAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		if el, ok := convertHudElement(act.Element); ok {
			pl.HideHudElement(el)
		}
	})
}

// UI closers
func (m *Manager) handlePlayerCloseDialogue(act *pb.PlayerCloseDialogueAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.CloseDialogue() })
}

func (m *Manager) handlePlayerCloseForm(act *pb.PlayerCloseFormAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.CloseForm() })
}

// Signs & Lecterns
func (m *Manager) handlePlayerOpenSign(act *pb.PlayerOpenSignAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	if act.Position == nil {
		return
	}
	pos := cube.Pos{int(act.Position.X), int(act.Position.Y), int(act.Position.Z)}
	m.execMethod(id, func(pl *player.Player) { pl.OpenSign(pos, act.FrontSide) })
}

func (m *Manager) handlePlayerEditSign(act *pb.PlayerEditSignAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	if act.Position == nil {
		return
	}
	pos := cube.Pos{int(act.Position.X), int(act.Position.Y), int(act.Position.Z)}
	m.execMethod(id, func(pl *player.Player) { _ = pl.EditSign(pos, act.FrontText, act.BackText) })
}

func (m *Manager) handlePlayerTurnLecternPage(act *pb.PlayerTurnLecternPageAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	if act.Position == nil {
		return
	}
	pos := cube.Pos{int(act.Position.X), int(act.Position.Y), int(act.Position.Z)}
	page := int(act.Page)
	m.execMethod(id, func(pl *player.Player) { _ = pl.TurnLecternPage(pos, page) })
}

// Entity visibility (players)
func (m *Manager) handlePlayerHidePlayer(act *pb.PlayerHidePlayerAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	targetID, err := uuid.Parse(act.TargetUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		for p := range m.srv.Players(nil) {
			if p.UUID() == targetID {
				pl.HideEntity(p)
				break
			}
		}
	})
}

func (m *Manager) handlePlayerShowPlayer(act *pb.PlayerShowPlayerAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	targetID, err := uuid.Parse(act.TargetUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		for p := range m.srv.Players(nil) {
			if p.UUID() == targetID {
				pl.ShowEntity(p)
				break
			}
		}
	})
}

// Debug shapes
func (m *Manager) handlePlayerRemoveAllDebugShapes(act *pb.PlayerRemoveAllDebugShapesAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) { pl.RemoveAllDebugShapes() })
}

// Interaction extras
func (m *Manager) handlePlayerOpenBlockContainer(act *pb.PlayerOpenBlockContainerAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	if act.Position == nil {
		return
	}
	pos := cube.Pos{int(act.Position.X), int(act.Position.Y), int(act.Position.Z)}
	m.execMethod(id, func(pl *player.Player) { pl.OpenBlockContainer(pos, pl.Tx()) })
}

func (m *Manager) handlePlayerDropItem(act *pb.PlayerDropItemAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		var s item.Stack
		if act.Item != nil {
			if stack, ok := convertProtoItemStackValue(act.Item); ok {
				s = stack
			} else {
				return
			}
		} else {
			held, _ := pl.HeldItems()
			if held.Empty() {
				return
			}
			s = held
		}
		_ = pl.Drop(s)
	})
}

func (m *Manager) handlePlayerSetItemCooldown(act *pb.PlayerSetItemCooldownAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	if act.Item == nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		if stack, ok := convertProtoItemStackValue(act.Item); ok {
			d := time.Duration(act.DurationMs) * time.Millisecond
			pl.SetCooldown(stack.Item(), d)
		}
	})
}

// Converters
func convertBossBarColour(c pb.BossBarColour) bossbar.Colour {
	switch c {
	case pb.BossBarColour_BOSS_BAR_COLOUR_GREY:
		return bossbar.White()
	case pb.BossBarColour_BOSS_BAR_COLOUR_BLUE:
		return bossbar.Blue()
	case pb.BossBarColour_BOSS_BAR_COLOUR_RED:
		return bossbar.Red()
	case pb.BossBarColour_BOSS_BAR_COLOUR_GREEN:
		return bossbar.Green()
	case pb.BossBarColour_BOSS_BAR_COLOUR_YELLOW:
		return bossbar.Yellow()
	case pb.BossBarColour_BOSS_BAR_COLOUR_PURPLE:
		return bossbar.Purple()
	case pb.BossBarColour_BOSS_BAR_COLOUR_WHITE:
		return bossbar.White()
	default:
		return bossbar.Purple()
	}
}

func convertHudElement(e pb.HudElement) (hud.Element, bool) {
	switch e {
	case pb.HudElement_HUD_ELEMENT_PAPER_DOLL:
		return hud.PaperDoll(), true
	case pb.HudElement_HUD_ELEMENT_ARMOUR:
		return hud.Armour(), true
	case pb.HudElement_HUD_ELEMENT_TOOL_TIPS:
		return hud.ToolTips(), true
	case pb.HudElement_HUD_ELEMENT_TOUCH_CONTROLS:
		return hud.TouchControls(), true
	case pb.HudElement_HUD_ELEMENT_CROSSHAIR:
		return hud.Crosshair(), true
	case pb.HudElement_HUD_ELEMENT_HOT_BAR:
		return hud.HotBar(), true
	case pb.HudElement_HUD_ELEMENT_HEALTH:
		return hud.Health(), true
	case pb.HudElement_HUD_ELEMENT_PROGRESS_BAR:
		return hud.ProgressBar(), true
	case pb.HudElement_HUD_ELEMENT_HUNGER:
		return hud.Hunger(), true
	case pb.HudElement_HUD_ELEMENT_AIR_BUBBLES:
		return hud.AirBubbles(), true
	case pb.HudElement_HUD_ELEMENT_HORSE_HEALTH:
		return hud.HorseHealth(), true
	case pb.HudElement_HUD_ELEMENT_STATUS_EFFECTS:
		return hud.StatusEffects(), true
	case pb.HudElement_HUD_ELEMENT_ITEM_TEXT:
		return hud.ItemText(), true
	default:
		return hud.Element{}, false
	}
}

// Local no-op submittables for forms/dialogues.
type formMenuNoop struct{}

func (formMenuNoop) Submit(form.Submitter, form.Button, *world.Tx) {}

type formModalNoop struct {
	Yes form.Button
	No  form.Button
}

func (formModalNoop) Submit(form.Submitter, form.Button, *world.Tx) {}

type dialogueNoop struct{}

func (dialogueNoop) Submit(dialogue.Submitter, dialogue.Button, *world.Tx) {}

func resolveWorldEntity(pl *player.Player, ref *pb.EntityRef) world.Entity {
	if ref == nil || ref.Uuid == "" {
		return nil
	}
	if targetUUID, err := uuid.Parse(ref.Uuid); err == nil {
		for ent := range pl.Tx().Entities() { // TODO: optimize, direct lookup with df private handles map
			if ent.H().UUID() == targetUUID {
				return ent
			}
		}
	}
	return nil
}

// Player armour
func (m *Manager) handlePlayerSetArmour(act *pb.PlayerSetArmourAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		if act.Helmet != nil {
			if s, ok := convertProtoItemStackValue(act.Helmet); ok {
				pl.Armour().SetHelmet(s)
			} else {
				pl.Armour().SetHelmet(item.Stack{})
			}
		}
		if act.Chestplate != nil {
			if s, ok := convertProtoItemStackValue(act.Chestplate); ok {
				pl.Armour().SetChestplate(s)
			} else {
				pl.Armour().SetChestplate(item.Stack{})
			}
		}
		if act.Leggings != nil {
			if s, ok := convertProtoItemStackValue(act.Leggings); ok {
				pl.Armour().SetLeggings(s)
			} else {
				pl.Armour().SetLeggings(item.Stack{})
			}
		}
		if act.Boots != nil {
			if s, ok := convertProtoItemStackValue(act.Boots); ok {
				pl.Armour().SetBoots(s)
			} else {
				pl.Armour().SetBoots(item.Stack{})
			}
		}
	})
}

// Player scoreboard
func (m *Manager) handlePlayerSendScoreboard(act *pb.PlayerSendScoreboardAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		sb := scoreboard.New(act.Title)
		if act.Padding != nil && !*act.Padding {
			sb.RemovePadding()
		}
		if act.Descending != nil && *act.Descending {
			sb.SetDescending()
		}
		// Clamp to 15 lines as per Dragonfly's limit and set them deterministically without trailing newlines.
		max := min(len(act.Lines), 15)
		for i := range max {
			sb.Set(i, act.Lines[i])
		}
		pl.SendScoreboard(sb)
	})
}

func (m *Manager) handlePlayerRemoveScoreboard(act *pb.PlayerRemoveScoreboardAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		pl.RemoveScoreboard()
	})
}

// Player forms (show)
func (m *Manager) handlePlayerSendMenuForm(act *pb.PlayerSendMenuFormAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		sub := formMenuNoop{}
		menu := form.NewMenu(sub, act.Title)
		if act.Body != nil {
			menu = menu.WithBody(*act.Body)
		}
		btns := make([]form.Button, len(act.Buttons))
		for i := range act.Buttons {
			btns[i] = form.NewButton(act.Buttons[i], "")
		}
		menu = menu.WithButtons(btns...)
		pl.SendForm(menu)
	})
}

func (m *Manager) handlePlayerSendModalForm(act *pb.PlayerSendModalFormAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		sub := formModalNoop{Yes: form.NewButton(act.YesText, ""), No: form.NewButton(act.NoText, "")}
		modal := form.NewModal(sub, act.Title).WithBody(act.Body)
		pl.SendForm(modal)
	})
}

// Player dialogue (show)
func (m *Manager) handlePlayerSendDialogue(p *pluginProcess, correlationID string, act *pb.PlayerSendDialogueAction) {
	id, err := uuid.Parse(act.PlayerUuid)
	if err != nil {
		m.sendActionError(p, correlationID, "invalid player_uuid")
		return
	}
	m.execMethod(id, func(pl *player.Player) {
		sub := dialogueNoop{}
		d := dialogue.New(sub, act.Title)
		if act.Body != nil {
			d = d.WithBody(*act.Body)
		}
		// Clamp to 6
		max := min(len(act.Buttons), 6)
		btns := make([]dialogue.Button, max)
		for i := range max {
			btns[i] = dialogue.Button{Text: act.Buttons[i]}
		}
		d = d.WithButtons(btns...)

		e := resolveWorldEntity(pl, act.Entity)
		if e == nil {
			m.sendActionError(p, correlationID, "entity not found")
			return
		}
		pl.SendDialogue(d, e)
	})
}
