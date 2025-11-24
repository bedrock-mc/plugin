<?php

namespace Dragonfly\PluginLib\Server;

use Dragonfly\PluginLib\Actions\Actions;
use Dragonfly\PluginLib\Entity\Player;

/**
 * Server provides access to online players and server-wide operations.
 * 
 * Player registry is automatically populated via PlayerJoin/PlayerQuit events
 * when using PluginBase.
 */
final class Server {
    /** @var array<string, string> uuid => name */
    private array $players = [];

    /** @var array<string, string> lowercase name => uuid */
    private array $nameIndex = [];

    public function __construct(
        private Actions $actions,
    ) {}

    /**
     * Register a player as online. Called automatically on PlayerJoin.
     */
    public function addPlayer(string $uuid, string $name): void {
        $this->players[$uuid] = $name;
        $this->nameIndex[strtolower($name)] = $uuid;
    }

    /**
     * Remove a player from the registry. Called automatically on PlayerQuit.
     */
    public function removePlayer(string $uuid): void {
        if (isset($this->players[$uuid])) {
            $name = $this->players[$uuid];
            unset($this->nameIndex[strtolower($name)]);
            unset($this->players[$uuid]);
        }
    }

    /**
     * Get a player by UUID. Returns null if not online.
     */
    public function getPlayer(string $uuid): ?Player {
        if (!isset($this->players[$uuid])) {
            return null;
        }
        return new Player($uuid, $this->players[$uuid], $this->actions);
    }

    /**
     * Get a player by name (case-insensitive). Returns null if not online.
     */
    public function getPlayerByName(string $name): ?Player {
        $lower = strtolower($name);
        if (!isset($this->nameIndex[$lower])) {
            return null;
        }
        $uuid = $this->nameIndex[$lower];
        return new Player($uuid, $this->players[$uuid], $this->actions);
    }

    /**
     * Get all online players.
     * 
     * @return Player[]
     */
    public function getOnlinePlayers(): array {
        $result = [];
        foreach ($this->players as $uuid => $name) {
            $result[] = new Player($uuid, $name, $this->actions);
        }
        return $result;
    }

    /**
     * Get the count of online players.
     */
    public function getOnlineCount(): int {
        return count($this->players);
    }

    /**
     * Check if a player is online by UUID.
     */
    public function isOnline(string $uuid): bool {
        return isset($this->players[$uuid]);
    }

    /**
     * Check if a player is online by name (case-insensitive).
     */
    public function isOnlineByName(string $name): bool {
        return isset($this->nameIndex[strtolower($name)]);
    }

    /**
     * Broadcast a message to all online players.
     */
    public function broadcastMessage(string $message): void {
        foreach ($this->players as $uuid => $_) {
            $this->actions->chatToUuid($uuid, $message);
        }
    }
}

