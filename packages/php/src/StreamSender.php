<?php

namespace Dragonfly\PluginLib;

use Df\Plugin\EventResult;
use Df\Plugin\PluginToHost;

class StreamSender {
    private $call;

    public function __construct($call) {
        $this->call = $call;
    }

    public function enqueue(PluginToHost $message): void {
        $this->call->write($message);
    }

    public function sendEventResult(string $pluginId, EventResult $result): void {
        $resp = new PluginToHost();
        $resp->setPluginId($pluginId);
        $resp->setEventResult($result);
        $this->enqueue($resp);
    }

    /**
     * Build and send an EventResult via a mutator closure.
     * The closure receives the EventResult to populate.
     */
    public function respond(string $pluginId, string $eventId, callable $mutator): void {
        $result = new EventResult();
        $result->setEventId($eventId);
        $mutator($result);
        $this->sendEventResult($pluginId, $result);
    }

}
