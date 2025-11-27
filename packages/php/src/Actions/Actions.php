<?php

namespace Dragonfly\PluginLib\Actions;

use Df\Plugin\Action;
use Df\Plugin\ActionBatch;
use Df\Plugin\AddEffectAction;
use Df\Plugin\PluginToHost;
use Df\Plugin\ClearInventoryAction;
use Df\Plugin\ExecuteCommandAction;
use Df\Plugin\GiveItemAction;
use Df\Plugin\ItemStack;
use Df\Plugin\PlaySoundAction;
use Df\Plugin\RemoveEffectAction;
use Df\Plugin\SendChatAction;
use Df\Plugin\SendPopupAction;
use Df\Plugin\SendTipAction;
use Df\Plugin\SendTitleAction;
use Df\Plugin\SetExperienceAction;
use Df\Plugin\SetFoodAction;
use Df\Plugin\TeleportAction;
use Df\Plugin\KickAction;
use Df\Plugin\SetGameModeAction;
use Df\Plugin\SetHealthAction;
use Df\Plugin\SetHeldItemAction;
use Df\Plugin\SetVelocityAction;
use Df\Plugin\Vec3;
use Df\Plugin\WorldSetDefaultGameModeAction;
use Df\Plugin\WorldSetDifficultyAction;
use Df\Plugin\WorldSetTickRangeAction;
use Df\Plugin\WorldSetBlockAction;
use Df\Plugin\WorldPlaySoundAction;
use Df\Plugin\WorldAddParticleAction;
use Df\Plugin\WorldQueryEntitiesAction;
use Df\Plugin\WorldQueryPlayersAction;
use Df\Plugin\WorldQueryEntitiesWithinAction;
use Df\Plugin\WorldRef;
use Df\Plugin\BlockPos;
use Df\Plugin\BlockState;
use Df\Plugin\BBox;
use Dragonfly\PluginLib\StreamSender;

final class Actions {
    private ?ActionBatch $activeBatch = null;

    public function __construct(
        private StreamSender $sender,
        private string $pluginId,
    ) {}

    public function startBatch(): void {
        $this->activeBatch = new ActionBatch();
    }

    public function commitBatch(): void {
        if ($this->activeBatch !== null && count($this->activeBatch->getActions()) > 0) {
            $toHost = new PluginToHost();
            $toHost->setPluginId($this->pluginId);
            $toHost->setActions($this->activeBatch);
            $this->sender->enqueue($toHost);
        }
        $this->activeBatch = null;
    }

    private function sendOrBatch(Action $action): void {
        if ($this->activeBatch !== null) {
            $actions = $this->activeBatch->getActions();
            $actions[] = $action;
            $this->activeBatch->setActions($actions);
        } else {
            $this->sendAction($action);
        }
    }

    public function sendActions(array $actions): void {
        $batch = new ActionBatch();
        $batch->setActions($actions);

        $toHost = new PluginToHost();
        $toHost->setPluginId($this->pluginId);
        $toHost->setActions($batch);
        $this->sender->enqueue($toHost);
    }

    private function sendAction(Action $action): void {
        $batch = new ActionBatch();
        $batch->setActions([$action]);

        $toHost = new PluginToHost();
        $toHost->setPluginId($this->pluginId);
        $toHost->setActions($batch);
        $this->sender->enqueue($toHost);
    }

    public function chatToUuid(string $uuid, string $message): void {
        $action = new Action();
        $chat = new SendChatAction();
        $chat->setTargetUuid($uuid);
        $chat->setMessage($message);
        $action->setSendChat($chat);
        $this->sendOrBatch($action);
    }

    public function teleportUuid(string $uuid, ?Vec3 $pos = null, ?Vec3 $rot = null): void {
        $action = new Action();
        $tp = new TeleportAction();
        $tp->setPlayerUuid($uuid);
        if ($pos !== null) $tp->setPosition($pos);
        if ($rot !== null) $tp->setRotation($rot);
        $action->setTeleport($tp);
        $this->sendOrBatch($action);
    }

    public function kickUuid(string $uuid, string $reason): void {
        $action = new Action();
        $k = new KickAction();
        $k->setPlayerUuid($uuid);
        $k->setReason($reason);
        $action->setKick($k);
        $this->sendOrBatch($action);
    }

    public function setGameModeUuid(string $uuid, int $mode): void {
        $action = new Action();
        $gm = new SetGameModeAction();
        $gm->setPlayerUuid($uuid);
        $gm->setGameMode($mode);
        $action->setSetGameMode($gm);
        $this->sendOrBatch($action);
    }

    public function giveItemUuid(string $uuid, ItemStack $item): void {
        $action = new Action();
        $gi = new GiveItemAction();
        $gi->setPlayerUuid($uuid);
        $gi->setItem($item);
        $action->setGiveItem($gi);
        $this->sendOrBatch($action);
    }

    public function clearInventoryUuid(string $uuid): void {
        $action = new Action();
        $ci = new ClearInventoryAction();
        $ci->setPlayerUuid($uuid);
        $action->setClearInventory($ci);
        $this->sendOrBatch($action);
    }

    public function setHeldItemsUuid(string $uuid, ?ItemStack $main = null, ?ItemStack $off = null): void {
        $action = new Action();
        $hi = new SetHeldItemAction();
        $hi->setPlayerUuid($uuid);
        if ($main !== null) $hi->setMain($main);
        if ($off !== null) $hi->setOffhand($off);
        $action->setSetHeldItem($hi);
        $this->sendOrBatch($action);
    }

    public function setHealthUuid(string $uuid, float $health, ?float $max = null): void {
        $action = new Action();
        $sh = new SetHealthAction();
        $sh->setPlayerUuid($uuid);
        $sh->setHealth($health);
        if ($max !== null) $sh->setMaxHealth($max);
        $action->setSetHealth($sh);
        $this->sendOrBatch($action);
    }

    public function setFoodUuid(string $uuid, int $food): void {
        $action = new Action();
        $sf = new SetFoodAction();
        $sf->setPlayerUuid($uuid);
        $sf->setFood($food);
        $action->setSetFood($sf);
        $this->sendOrBatch($action);
    }

    public function setExperienceUuid(string $uuid, ?int $level = null, ?float $progress = null, ?int $amount = null): void {
        $action = new Action();
        $xp = new SetExperienceAction();
        $xp->setPlayerUuid($uuid);
        if ($level !== null) $xp->setLevel($level);
        if ($progress !== null) $xp->setProgress($progress);
        if ($amount !== null) $xp->setAmount($amount);
        $action->setSetExperience($xp);
        $this->sendOrBatch($action);
    }

    public function setVelocityUuid(string $uuid, Vec3 $vel): void {
        $action = new Action();
        $sv = new SetVelocityAction();
        $sv->setPlayerUuid($uuid);
        $sv->setVelocity($vel);
        $action->setSetVelocity($sv);
        $this->sendOrBatch($action);
    }

    public function addEffectUuid(string $uuid, int $type, int $level, int $durationMs, bool $showParticles = true): void {
        $action = new Action();
        $fx = new AddEffectAction();
        $fx->setPlayerUuid($uuid);
        $fx->setEffectType($type);
        $fx->setLevel($level);
        $fx->setDurationMs($durationMs);
        $fx->setShowParticles($showParticles);
        $action->setAddEffect($fx);
        $this->sendOrBatch($action);
    }

    public function removeEffectUuid(string $uuid, int $type): void {
        $action = new Action();
        $fx = new RemoveEffectAction();
        $fx->setPlayerUuid($uuid);
        $fx->setEffectType($type);
        $action->setRemoveEffect($fx);
        $this->sendOrBatch($action);
    }

    public function sendTitleUuid(string $uuid, string $title, ?string $subtitle = null, ?int $fadeIn = null, ?int $duration = null, ?int $fadeOut = null): void {
        $action = new Action();
        $st = new SendTitleAction();
        $st->setPlayerUuid($uuid);
        $st->setTitle($title);
        if ($subtitle !== null) $st->setSubtitle($subtitle);
        if ($fadeIn !== null) $st->setFadeInMs($fadeIn);
        if ($duration !== null) $st->setDurationMs($duration);
        if ($fadeOut !== null) $st->setFadeOutMs($fadeOut);
        $action->setSendTitle($st);
        $this->sendOrBatch($action);
    }

    public function sendPopupUuid(string $uuid, string $msg): void {
        $action = new Action();
        $pu = new SendPopupAction();
        $pu->setPlayerUuid($uuid);
        $pu->setMessage($msg);
        $action->setSendPopup($pu);
        $this->sendOrBatch($action);
    }

    public function sendTipUuid(string $uuid, string $msg): void {
        $action = new Action();
        $tp = new SendTipAction();
        $tp->setPlayerUuid($uuid);
        $tp->setMessage($msg);
        $action->setSendTip($tp);
        $this->sendOrBatch($action);
    }

    public function playSoundUuid(string $uuid, int $soundId, ?Vec3 $pos = null, ?float $volume = null, ?float $pitch = null): void {
        $action = new Action();
        $ps = new PlaySoundAction();
        $ps->setPlayerUuid($uuid);
        $ps->setSound($soundId);
        if ($pos !== null) $ps->setPosition($pos);
        if ($volume !== null) $ps->setVolume($volume);
        if ($pitch !== null) $ps->setPitch($pitch);
        $action->setPlaySound($ps);
        $this->sendOrBatch($action);
    }

    public function executeCommandUuid(string $uuid, string $cmd): void {
        $action = new Action();
        $ec = new ExecuteCommandAction();
        $ec->setPlayerUuid($uuid);
        $ec->setCommand($cmd);
        $action->setExecuteCommand($ec);
        $this->sendOrBatch($action);
    }

    public function worldSetDefaultGameMode(WorldRef $world, int $mode): void {
        $action = new Action();
        $w = new WorldSetDefaultGameModeAction();
        $w->setWorld($world);
        $w->setGameMode($mode);
        $action->setWorldSetDefaultGameMode($w);
        $this->sendOrBatch($action);
    }

    public function worldSetDifficulty(WorldRef $world, int $difficulty): void {
        $action = new Action();
        $wd = new WorldSetDifficultyAction();
        $wd->setWorld($world);
        $wd->setDifficulty($difficulty);
        $action->setWorldSetDifficulty($wd);
        $this->sendOrBatch($action);
    }

    public function worldSetTickRange(WorldRef $world, int $range): void {
        $action = new Action();
        $wr = new WorldSetTickRangeAction();
        $wr->setWorld($world);
        $wr->setTickRange($range);
        $action->setWorldSetTickRange($wr);
        $this->sendOrBatch($action);
    }

    public function worldSetBlock(WorldRef $world, BlockPos $pos, ?BlockState $state = null): void {
        $action = new Action();
        $wb = new WorldSetBlockAction();
        $wb->setWorld($world);
        $wb->setPosition($pos);
        if ($state !== null) $wb->setBlock($state);
        $action->setWorldSetBlock($wb);
        $this->sendOrBatch($action);
    }

    public function worldPlaySound(WorldRef $world, int $soundId, Vec3 $pos): void {
        $action = new Action();
        $ws = new WorldPlaySoundAction();
        $ws->setWorld($world);
        $ws->setSound($soundId);
        $ws->setPosition($pos);
        $action->setWorldPlaySound($ws);
        $this->sendOrBatch($action);
    }

    public function worldAddParticle(WorldRef $world, Vec3 $pos, int $particle, ?BlockState $block = null, ?int $face = null): void {
        $action = new Action();
        $wp = new WorldAddParticleAction();
        $wp->setWorld($world);
        $wp->setPosition($pos);
        $wp->setParticle($particle);
        if ($block !== null) $wp->setBlock($block);
        if ($face !== null) $wp->setFace($face);
        $action->setWorldAddParticle($wp);
        $this->sendOrBatch($action);
    }

    public function worldQueryEntities(WorldRef $world, ?string $corr = null): void {
        $action = new Action();
        if ($corr !== null) $action->setCorrelationId($corr);
        $q = new WorldQueryEntitiesAction();
        $q->setWorld($world);
        $action->setWorldQueryEntities($q);
        $this->sendOrBatch($action);
    }

    public function worldQueryPlayers(WorldRef $world, ?string $corr = null): void {
        $action = new Action();
        if ($corr !== null) $action->setCorrelationId($corr);
        $q = new WorldQueryPlayersAction();
        $q->setWorld($world);
        $action->setWorldQueryPlayers($q);
        $this->sendOrBatch($action);
    }

    public function worldQueryEntitiesWithin(WorldRef $world, BBox $box, ?string $corr = null): void {
        $action = new Action();
        if ($corr !== null) $action->setCorrelationId($corr);
        $q = new WorldQueryEntitiesWithinAction();
        $q->setWorld($world);
        $q->setBox($box);
        $action->setWorldQueryEntitiesWithin($q);
        $this->sendOrBatch($action);
    }
}
