<?php

namespace Dragonfly\PluginLib\Events;

/**
 * Marker interface for listener classes.
 * Public, non-static methods where the first parameter is a known event payload
 * type (e.g. \Df\Plugin\ChatEvent) are auto-registered. The second parameter
 * may optionally be an EventContext.
 */
interface Listener {}