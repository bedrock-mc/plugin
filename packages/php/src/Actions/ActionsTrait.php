<?php

namespace Dragonfly\PluginLib\Actions;

use Df\Plugin\Action;
use Df\Plugin\Vec3;
use Df\Plugin\ItemStack;
use Df\Plugin\WorldRef;
use Df\Plugin\BlockPos;
use Df\Plugin\BlockState;
use Df\Plugin\BBox;

trait ActionsTrait {
    abstract protected function getActions(): Actions;

    public function sendAction(Action $action): void {
        $this->getActions()->sendAction($action);
    }

    public function chatToUuid(string $targetUuid, string $message): void {
        $this->getActions()->chatToUuid($targetUuid, $message);
    }

    public function teleportUuid(string $playerUuid, ?Vec3 $position = null, ?Vec3 $rotation = null): void {
        $this->getActions()->teleportUuid($playerUuid, $position, $rotation);
    }

    public function kickUuid(string $playerUuid, string $reason): void {
        $this->getActions()->kickUuid($playerUuid, $reason);
    }

    public function setGameModeUuid(string $playerUuid, int $gameMode): void {
        $this->getActions()->setGameModeUuid($playerUuid, $gameMode);
    }

    public function giveItemUuid(string $playerUuid, ItemStack $item): void {
        $this->getActions()->giveItemUuid($playerUuid, $item);
    }

    public function clearInventoryUuid(string $playerUuid): void {
        $this->getActions()->clearInventoryUuid($playerUuid);
    }

    public function setHeldItemsUuid(string $playerUuid, ?ItemStack $main = null, ?ItemStack $offhand = null): void {
        $this->getActions()->setHeldItemsUuid($playerUuid, $main, $offhand);
    }

    public function setHealthUuid(string $playerUuid, float $health, ?float $maxHealth = null): void {
        $this->getActions()->setHealthUuid($playerUuid, $health, $maxHealth);
    }

    public function setFoodUuid(string $playerUuid, int $food): void {
        $this->getActions()->setFoodUuid($playerUuid, $food);
    }

    public function setExperienceUuid(string $playerUuid, ?int $level = null, ?float $progress = null, ?int $amount = null): void {
        $this->getActions()->setExperienceUuid($playerUuid, $level, $progress, $amount);
    }

    public function setVelocityUuid(string $playerUuid, Vec3 $velocity): void {
        $this->getActions()->setVelocityUuid($playerUuid, $velocity);
    }

    public function addEffectUuid(string $playerUuid, int $effectType, int $level, int $durationMs, bool $showParticles = true): void {
        $this->getActions()->addEffectUuid($playerUuid, $effectType, $level, $durationMs, $showParticles);
    }

    public function removeEffectUuid(string $playerUuid, int $effectType): void {
        $this->getActions()->removeEffectUuid($playerUuid, $effectType);
    }

    public function sendTitleUuid(string $playerUuid, string $title, ?string $subtitle = null, ?int $fadeInMs = null, ?int $durationMs = null, ?int $fadeOutMs = null): void {
        $this->getActions()->sendTitleUuid($playerUuid, $title, $subtitle, $fadeInMs, $durationMs, $fadeOutMs);
    }

    public function sendPopupUuid(string $playerUuid, string $message): void {
        $this->getActions()->sendPopupUuid($playerUuid, $message);
    }

    public function sendTipUuid(string $playerUuid, string $message): void {
        $this->getActions()->sendTipUuid($playerUuid, $message);
    }

    public function playSoundUuid(string $playerUuid, int $sound, ?Vec3 $position = null, ?float $volume = null, ?float $pitch = null): void {
        $this->getActions()->playSoundUuid($playerUuid, $sound, $position, $volume, $pitch);
    }

    public function executeCommandUuid(string $playerUuid, string $command): void {
        $this->getActions()->executeCommandUuid($playerUuid, $command);
    }

    public function worldSetDefaultGameMode(WorldRef $world, int $gameMode): void {
        $this->getActions()->worldSetDefaultGameMode($world, $gameMode);
    }

    public function worldSetDifficulty(WorldRef $world, int $difficulty): void {
        $this->getActions()->worldSetDifficulty($world, $difficulty);
    }

    public function worldSetTickRange(WorldRef $world, int $tickRange): void {
        $this->getActions()->worldSetTickRange($world, $tickRange);
    }

    public function worldSetBlock(WorldRef $world, BlockPos $position, ?BlockState $block = null): void {
        $this->getActions()->worldSetBlock($world, $position, $block);
    }

    public function worldPlaySound(WorldRef $world, int $sound, Vec3 $position): void {
        $this->getActions()->worldPlaySound($world, $sound, $position);
    }

    public function worldAddParticle(WorldRef $world, Vec3 $position, int $particle, ?BlockState $block = null, ?int $face = null): void {
        $this->getActions()->worldAddParticle($world, $position, $particle, $block, $face);
    }

    public function worldQueryEntities(WorldRef $world, ?string $correlationId = null): void {
        $this->getActions()->worldQueryEntities($world, $correlationId);
    }

    public function worldQueryPlayers(WorldRef $world, ?string $correlationId = null): void {
        $this->getActions()->worldQueryPlayers($world, $correlationId);
    }

    public function worldQueryEntitiesWithin(WorldRef $world, BBox $box, ?string $correlationId = null): void {
        $this->getActions()->worldQueryEntitiesWithin($world, $box, $correlationId);
    }
}
