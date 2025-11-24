<?php

namespace Dragonfly\PluginLib\Events;

trait MutationsTrait {
    abstract public function respondWith(object $mutation): void;
    abstract public function respond(callable $mutator): void;

    public function chat(string $message): void
    {
        $this->respond(function (\Df\Plugin\EventResult $r) use ($message): void {
            $m = new \Df\Plugin\ChatMutation();
            $m->setMessage($message);
            $r->setChat($m);
        });
    }

    public function blockBreak(\Df\Plugin\BlockBreakMutation $m): void { $this->respondWith($m); }
    public function playerFoodLoss(\Df\Plugin\PlayerFoodLossMutation $m): void { $this->respondWith($m); }
    public function playerHeal(\Df\Plugin\PlayerHealMutation $m): void { $this->respondWith($m); }
    public function playerHurt(\Df\Plugin\PlayerHurtMutation $m): void { $this->respondWith($m); }
    public function playerDeath(\Df\Plugin\PlayerDeathMutation $m): void { $this->respondWith($m); }
    public function playerRespawn(\Df\Plugin\PlayerRespawnMutation $m): void { $this->respondWith($m); }
    public function playerAttackEntity(\Df\Plugin\PlayerAttackEntityMutation $m): void { $this->respondWith($m); }
    public function playerExperienceGain(\Df\Plugin\PlayerExperienceGainMutation $m): void { $this->respondWith($m); }
    public function playerLecternPageTurn(\Df\Plugin\PlayerLecternPageTurnMutation $m): void { $this->respondWith($m); }
    public function playerItemPickup(\Df\Plugin\PlayerItemPickupMutation $m): void { $this->respondWith($m); }
    public function playerTransfer(\Df\Plugin\PlayerTransferMutation $m): void { $this->respondWith($m); }
    public function worldExplosion(\Df\Plugin\WorldExplosionMutation $m): void { $this->respondWith($m); }
}
