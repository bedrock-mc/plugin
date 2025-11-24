<?php

namespace Dragonfly\PluginLib\Commands;

use Dragonfly\PluginLib\Actions\Actions;
use Dragonfly\PluginLib\Entity\Player;

class CommandSender extends Player {
    public function __construct(
        public string $uuid,
        public string $name,
        Actions $actions,
    ) {
        parent::__construct($uuid, $name, $actions);
    }
}
