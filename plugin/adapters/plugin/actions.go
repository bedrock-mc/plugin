package plugin

import (
	pb "github.com/bedrock-mc/plugin/proto/generated/go"
)

func (m *Manager) applyActions(p *pluginProcess, batch *pb.ActionBatch) {
	if batch == nil {
		return
	}
	for _, action := range batch.Actions {
		if action == nil {
			continue
		}
		correlationID := action.GetCorrelationId()
		switch kind := action.Kind.(type) {
		case *pb.Action_SendChat:
			m.handleSendChat(kind.SendChat)
		case *pb.Action_Teleport:
			m.handleTeleport(kind.Teleport)
		case *pb.Action_Kick:
			m.handleKick(kind.Kick)
		case *pb.Action_SetGameMode:
			m.handleSetGameMode(kind.SetGameMode)
		case *pb.Action_GiveItem:
			m.handleGiveItem(kind.GiveItem)
		case *pb.Action_ClearInventory:
			m.handleClearInventory(kind.ClearInventory)
		case *pb.Action_SetHeldItem:
			m.handleSetHeldItem(kind.SetHeldItem)
		case *pb.Action_PlayerSetArmour:
			m.handlePlayerSetArmour(kind.PlayerSetArmour)
		case *pb.Action_SetHealth:
			m.handleSetHealth(kind.SetHealth)
		case *pb.Action_SetFood:
			m.handleSetFood(kind.SetFood)
		case *pb.Action_SetExperience:
			m.handleSetExperience(kind.SetExperience)
		case *pb.Action_SetVelocity:
			m.handleSetVelocity(kind.SetVelocity)
		case *pb.Action_AddEffect:
			m.handleAddEffect(kind.AddEffect)
		case *pb.Action_RemoveEffect:
			m.handleRemoveEffect(kind.RemoveEffect)
		case *pb.Action_SendTitle:
			m.handleSendTitle(kind.SendTitle)
		case *pb.Action_SendPopup:
			m.handleSendPopup(kind.SendPopup)
		case *pb.Action_SendTip:
			m.handleSendTip(kind.SendTip)
		case *pb.Action_PlaySound:
			m.handlePlaySound(kind.PlaySound)
		case *pb.Action_ExecuteCommand:
			m.handleExecuteCommand(kind.ExecuteCommand)
		case *pb.Action_WorldSetDefaultGameMode:
			m.handleWorldSetDefaultGameMode(p, correlationID, kind.WorldSetDefaultGameMode)
		case *pb.Action_WorldSetDifficulty:
			m.handleWorldSetDifficulty(p, correlationID, kind.WorldSetDifficulty)
		case *pb.Action_WorldSetTickRange:
			m.handleWorldSetTickRange(p, correlationID, kind.WorldSetTickRange)
		case *pb.Action_WorldSetBlock:
			m.handleWorldSetBlock(p, correlationID, kind.WorldSetBlock)
		case *pb.Action_WorldPlaySound:
			m.handleWorldPlaySound(p, correlationID, kind.WorldPlaySound)
		case *pb.Action_WorldAddParticle:
			m.handleWorldAddParticle(p, correlationID, kind.WorldAddParticle)
		case *pb.Action_WorldSetTime:
			m.handleWorldSetTime(p, correlationID, kind.WorldSetTime)
		case *pb.Action_WorldStopTime:
			m.handleWorldStopTime(p, correlationID, kind.WorldStopTime)
		case *pb.Action_WorldStartTime:
			m.handleWorldStartTime(p, correlationID, kind.WorldStartTime)
		case *pb.Action_WorldSetSpawn:
			m.handleWorldSetSpawn(p, correlationID, kind.WorldSetSpawn)
		case *pb.Action_WorldQueryEntities:
			m.handleWorldQueryEntities(p, correlationID, kind.WorldQueryEntities)
		case *pb.Action_WorldQueryPlayers:
			m.handleWorldQueryPlayers(p, correlationID, kind.WorldQueryPlayers)
		case *pb.Action_WorldQueryEntitiesWithin:
			m.handleWorldQueryEntitiesWithin(p, correlationID, kind.WorldQueryEntitiesWithin)
		case *pb.Action_WorldQueryDefaultGameMode:
			m.handleWorldQueryDefaultGameMode(p, correlationID, kind.WorldQueryDefaultGameMode)
		case *pb.Action_WorldQueryPlayerSpawn:
			m.handleWorldQueryPlayerSpawn(p, correlationID, kind.WorldQueryPlayerSpawn)
		case *pb.Action_WorldQueryBlock:
			m.handleWorldQueryBlock(p, correlationID, kind.WorldQueryBlock)
		case *pb.Action_WorldQueryBiome:
			m.handleWorldQueryBiome(p, correlationID, kind.WorldQueryBiome)
		case *pb.Action_WorldQueryLight:
			m.handleWorldQueryLight(p, correlationID, kind.WorldQueryLight)
		case *pb.Action_WorldQuerySkyLight:
			m.handleWorldQuerySkyLight(p, correlationID, kind.WorldQuerySkyLight)
		case *pb.Action_WorldQueryTemperature:
			m.handleWorldQueryTemperature(p, correlationID, kind.WorldQueryTemperature)
		case *pb.Action_WorldQueryHighestBlock:
			m.handleWorldQueryHighestBlock(p, correlationID, kind.WorldQueryHighestBlock)
		case *pb.Action_WorldQueryRainingAt:
			m.handleWorldQueryRainingAt(p, correlationID, kind.WorldQueryRainingAt)
		case *pb.Action_WorldQuerySnowingAt:
			m.handleWorldQuerySnowingAt(p, correlationID, kind.WorldQuerySnowingAt)
		case *pb.Action_WorldQueryThunderingAt:
			m.handleWorldQueryThunderingAt(p, correlationID, kind.WorldQueryThunderingAt)
		case *pb.Action_WorldQueryLiquid:
			m.handleWorldQueryLiquid(p, correlationID, kind.WorldQueryLiquid)
		case *pb.Action_WorldSetBiome:
			m.handleWorldSetBiome(p, correlationID, kind.WorldSetBiome)
		case *pb.Action_WorldSetLiquid:
			m.handleWorldSetLiquid(p, correlationID, kind.WorldSetLiquid)
		case *pb.Action_WorldScheduleBlockUpdate:
			m.handleWorldScheduleBlockUpdate(p, correlationID, kind.WorldScheduleBlockUpdate)
		case *pb.Action_WorldBuildStructure:
			m.handleWorldBuildStructure(p, correlationID, kind.WorldBuildStructure)
		case *pb.Action_PlayerStartSprinting:
			m.handlePlayerStartSprinting(kind.PlayerStartSprinting)
		case *pb.Action_PlayerStopSprinting:
			m.handlePlayerStopSprinting(kind.PlayerStopSprinting)
		case *pb.Action_PlayerStartSneaking:
			m.handlePlayerStartSneaking(kind.PlayerStartSneaking)
		case *pb.Action_PlayerStopSneaking:
			m.handlePlayerStopSneaking(kind.PlayerStopSneaking)
		case *pb.Action_PlayerStartSwimming:
			m.handlePlayerStartSwimming(kind.PlayerStartSwimming)
		case *pb.Action_PlayerStopSwimming:
			m.handlePlayerStopSwimming(kind.PlayerStopSwimming)
		case *pb.Action_PlayerStartCrawling:
			m.handlePlayerStartCrawling(kind.PlayerStartCrawling)
		case *pb.Action_PlayerStopCrawling:
			m.handlePlayerStopCrawling(kind.PlayerStopCrawling)
		case *pb.Action_PlayerStartGliding:
			m.handlePlayerStartGliding(kind.PlayerStartGliding)
		case *pb.Action_PlayerStopGliding:
			m.handlePlayerStopGliding(kind.PlayerStopGliding)
		case *pb.Action_PlayerStartFlying:
			m.handlePlayerStartFlying(kind.PlayerStartFlying)
		case *pb.Action_PlayerStopFlying:
			m.handlePlayerStopFlying(kind.PlayerStopFlying)
		case *pb.Action_PlayerSetImmobile:
			m.handlePlayerSetImmobile(kind.PlayerSetImmobile)
		case *pb.Action_PlayerSetMobile:
			m.handlePlayerSetMobile(kind.PlayerSetMobile)
		case *pb.Action_PlayerSetSpeed:
			m.handlePlayerSetSpeed(kind.PlayerSetSpeed)
		case *pb.Action_PlayerSetFlightSpeed:
			m.handlePlayerSetFlightSpeed(kind.PlayerSetFlightSpeed)
		case *pb.Action_PlayerSetVerticalFlightSpeed:
			m.handlePlayerSetVerticalFlightSpeed(kind.PlayerSetVerticalFlightSpeed)
		case *pb.Action_PlayerSetAbsorption:
			m.handlePlayerSetAbsorption(kind.PlayerSetAbsorption)
		case *pb.Action_PlayerSetOnFire:
			m.handlePlayerSetOnFire(kind.PlayerSetOnFire)
		case *pb.Action_PlayerExtinguish:
			m.handlePlayerExtinguish(kind.PlayerExtinguish)
		case *pb.Action_PlayerSetInvisible:
			m.handlePlayerSetInvisible(kind.PlayerSetInvisible)
		case *pb.Action_PlayerSetVisible:
			m.handlePlayerSetVisible(kind.PlayerSetVisible)
		case *pb.Action_PlayerSetScale:
			m.handlePlayerSetScale(kind.PlayerSetScale)
		case *pb.Action_PlayerSetHeldSlot:
			m.handlePlayerSetHeldSlot(kind.PlayerSetHeldSlot)
		case *pb.Action_PlayerSendToast:
			m.handlePlayerSendToast(kind.PlayerSendToast)
		case *pb.Action_PlayerSendJukeboxPopup:
			m.handlePlayerSendJukeboxPopup(kind.PlayerSendJukeboxPopup)
		case *pb.Action_PlayerShowCoordinates:
			m.handlePlayerShowCoordinates(kind.PlayerShowCoordinates)
		case *pb.Action_PlayerHideCoordinates:
			m.handlePlayerHideCoordinates(kind.PlayerHideCoordinates)
		case *pb.Action_PlayerEnableInstantRespawn:
			m.handlePlayerEnableInstantRespawn(kind.PlayerEnableInstantRespawn)
		case *pb.Action_PlayerDisableInstantRespawn:
			m.handlePlayerDisableInstantRespawn(kind.PlayerDisableInstantRespawn)
		case *pb.Action_PlayerSetNameTag:
			m.handlePlayerSetNameTag(kind.PlayerSetNameTag)
		case *pb.Action_PlayerSetScoreTag:
			m.handlePlayerSetScoreTag(kind.PlayerSetScoreTag)
		case *pb.Action_PlayerShowParticle:
			m.handlePlayerShowParticle(kind.PlayerShowParticle)
		case *pb.Action_PlayerSendScoreboard:
			m.handlePlayerSendScoreboard(kind.PlayerSendScoreboard)
		case *pb.Action_PlayerRemoveScoreboard:
			m.handlePlayerRemoveScoreboard(kind.PlayerRemoveScoreboard)
		case *pb.Action_PlayerSendMenuForm:
			m.handlePlayerSendMenuForm(kind.PlayerSendMenuForm)
		case *pb.Action_PlayerSendModalForm:
			m.handlePlayerSendModalForm(kind.PlayerSendModalForm)
		case *pb.Action_PlayerSendDialogue:
			m.handlePlayerSendDialogue(p, correlationID, kind.PlayerSendDialogue)
		case *pb.Action_PlayerRespawn:
			m.handlePlayerRespawn(kind.PlayerRespawn)
		case *pb.Action_PlayerTransfer:
			m.handlePlayerTransferAction(kind.PlayerTransfer)
		case *pb.Action_PlayerKnockBack:
			m.handlePlayerKnockBack(kind.PlayerKnockBack)
		case *pb.Action_PlayerSwingArm:
			m.handlePlayerSwingArm(kind.PlayerSwingArm)
		case *pb.Action_PlayerPunchAir:
			m.handlePlayerPunchAirAction(kind.PlayerPunchAir)
		case *pb.Action_PlayerSendBossBar:
			m.handlePlayerSendBossBar(kind.PlayerSendBossBar)
		case *pb.Action_PlayerRemoveBossBar:
			m.handlePlayerRemoveBossBar(kind.PlayerRemoveBossBar)
		case *pb.Action_PlayerShowHudElement:
			m.handlePlayerShowHudElement(kind.PlayerShowHudElement)
		case *pb.Action_PlayerHideHudElement:
			m.handlePlayerHideHudElement(kind.PlayerHideHudElement)
		case *pb.Action_PlayerCloseDialogue:
			m.handlePlayerCloseDialogue(kind.PlayerCloseDialogue)
		case *pb.Action_PlayerCloseForm:
			m.handlePlayerCloseForm(kind.PlayerCloseForm)
		case *pb.Action_PlayerOpenSign:
			m.handlePlayerOpenSign(kind.PlayerOpenSign)
		case *pb.Action_PlayerEditSign:
			m.handlePlayerEditSign(kind.PlayerEditSign)
		case *pb.Action_PlayerTurnLecternPage:
			m.handlePlayerTurnLecternPage(kind.PlayerTurnLecternPage)
		case *pb.Action_PlayerHidePlayer:
			m.handlePlayerHidePlayer(kind.PlayerHidePlayer)
		case *pb.Action_PlayerShowPlayer:
			m.handlePlayerShowPlayer(kind.PlayerShowPlayer)
		case *pb.Action_PlayerRemoveAllDebugShapes:
			m.handlePlayerRemoveAllDebugShapes(kind.PlayerRemoveAllDebugShapes)
		case *pb.Action_PlayerOpenBlockContainer:
			m.handlePlayerOpenBlockContainer(kind.PlayerOpenBlockContainer)
		case *pb.Action_PlayerDropItem:
			m.handlePlayerDropItem(kind.PlayerDropItem)
		case *pb.Action_PlayerSetItemCooldown:
			m.handlePlayerSetItemCooldown(kind.PlayerSetItemCooldown)
		}
	}
}
