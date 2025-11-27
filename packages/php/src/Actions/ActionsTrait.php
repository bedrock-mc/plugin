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

    public function startBatch(): void {
        $this->getActions()->startBatch();
    }

    public function commitBatch(): void {
        $this->getActions()->commitBatch();
    }

    public function sendActions(array $actions): void {
        $this->getActions()->sendActions($actions);
    }

    public function chatToUuid(string $uuid, string $message): void {
        $this->getActions()->chatToUuid($uuid, $message);
    }

    public function teleportUuid(string $uuid, ?Vec3 $pos = null, ?Vec3 $rot = null): void {
        $this->getActions()->teleportUuid($uuid, $pos, $rot);
    }

    public function kickUuid(string $uuid, string $reason): void {
        $this->getActions()->kickUuid($uuid, $reason);
    }

    public function setGameModeUuid(string $uuid, int $mode): void {
        $this->getActions()->setGameModeUuid($uuid, $mode);
    }

    public function giveItemUuid(string $uuid, ItemStack $item): void {
        $this->getActions()->giveItemUuid($uuid, $item);
    }

    public function clearInventoryUuid(string $uuid): void {
        $this->getActions()->clearInventoryUuid($uuid);
    }

    public function setHeldItemsUuid(string $uuid, ?ItemStack $main = null, ?ItemStack $off = null): void {
        $this->getActions()->setHeldItemsUuid($uuid, $main, $off);
    }

    public function setHealthUuid(string $uuid, float $health, ?float $max = null): void {
        $this->getActions()->setHealthUuid($uuid, $health, $max);
    }

    public function setFoodUuid(string $uuid, int $food): void {
        $this->getActions()->setFoodUuid($uuid, $food);
    }

    public function setExperienceUuid(string $uuid, ?int $xp = null, ?float $progress = null, ?int $level = null): void {
        $this->getActions()->setExperienceUuid($uuid, $xp, $progress, $level);
    }

    public function setVelocityUuid(string $uuid, Vec3 $velocity): void {
        $this->getActions()->setVelocityUuid($uuid, $velocity);
    }

    public function addEffectUuid(string $uuid, int $id, int $duration, int $amplifier, bool $show = true): void {
        $this->getActions()->addEffectUuid($uuid, $id, $duration, $amplifier, $show);
    }

    public function removeEffectUuid(string $uuid, int $id): void {
        $this->getActions()->removeEffectUuid($uuid, $id);
    }

    public function sendTitleUuid(
        string $uuid,
        string $title,
        ?string $subtitle = null,
        ?int $fadeIn = null,
        ?int $stay = null,
        ?int $fadeOut = null
    ): void {
        $this->getActions()->sendTitleUuid($uuid, $title, $subtitle, $fadeIn, $stay, $fadeOut);
    }

    public function sendPopupUuid(string $uuid, string $message): void {
        $this->getActions()->sendPopupUuid($uuid, $message);
    }

    public function sendTipUuid(string $uuid, string $message): void {
        $this->getActions()->sendTipUuid($uuid, $message);
    }

    public function playSoundUuid(
        string $uuid,
        int $sound,
        ?Vec3 $pos = null,
        ?float $volume = null,
        ?float $pitch = null
    ): void {
        $this->getActions()->playSoundUuid($uuid, $sound, $pos, $volume, $pitch);
    }

    public function executeCommandUuid(string $uuid, string $command): void {
        $this->getActions()->executeCommandUuid($uuid, $command);
    }

    public function worldSetDefaultGameMode(WorldRef $world, int $mode): void {
        $this->getActions()->worldSetDefaultGameMode($world, $mode);
    }

    public function worldSetDifficulty(WorldRef $world, int $difficulty): void {
        $this->getActions()->worldSetDifficulty($world, $difficulty);
    }

    public function worldSetTickRange(WorldRef $world, int $range): void {
        $this->getActions()->worldSetTickRange($world, $range);
    }

    public function worldSetBlock(WorldRef $world, BlockPos $pos, ?BlockState $state = null): void {
        $this->getActions()->worldSetBlock($world, $pos, $state);
    }

    public function worldPlaySound(WorldRef $world, int $sound, Vec3 $pos): void {
        $this->getActions()->worldPlaySound($world, $sound, $pos);
    }

    public function worldAddParticle(
        WorldRef $world,
        Vec3 $pos,
        int $type,
        ?BlockState $state = null,
        ?int $count = null
    ): void {
        $this->getActions()->worldAddParticle($world, $pos, $type, $state, $count);
    }

    public function worldQueryEntities(WorldRef $world, ?string $selector = null): void {
        $this->getActions()->worldQueryEntities($world, $selector);
    }

    public function worldQueryPlayers(WorldRef $world, ?string $selector = null): void {
        $this->getActions()->worldQueryPlayers($world, $selector);
    }

    public function worldQueryEntitiesWithin(WorldRef $world, BBox $bbox, ?string $selector = null): void {
        $this->getActions()->worldQueryEntitiesWithin($world, $bbox, $selector);
    }
}
