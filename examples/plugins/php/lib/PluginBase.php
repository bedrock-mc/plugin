<?php

namespace Dragonfly\PluginLib;

use Df\Plugin\EventEnvelope;
use Df\Plugin\CommandSpec;
use Df\Plugin\EventSubscribe;
use Df\Plugin\EventType;
use Df\Plugin\PluginClient;
use Df\Plugin\PluginHello;
use Df\Plugin\PluginToHost;
use Dragonfly\PluginLib\Events\EventContext;
use Dragonfly\PluginLib\Events\Listener;
use Grpc\ChannelCredentials;
use ReflectionClass;
use ReflectionMethod;
use ReflectionNamedType;

abstract class PluginBase {
    protected string $pluginId;
    protected string $serverAddress;

    protected string $name = 'example-php';
    protected string $version = '0.1.0';
    protected string $apiVersion = 'v1';

    /** @var array<int, array<int, callable>> */
    private array $handlers = [];

    /** @var int[] */
    private array $subscriptions = [];

    /** @var PluginClient */
    private PluginClient $client;

    /** @var mixed */
    private $call;

    private StreamSender $sender;

    private bool $running = false;

    /** @var array<int, array{name: string, description: string}> */
    private array $commandSpecs = [];

    public function __construct(?string $pluginId = null, ?string $serverAddress = null) {
        $this->pluginId = $pluginId ?? (getenv('DF_PLUGIN_ID') ?: 'php-plugin');
        $address = $serverAddress ?? (getenv('DF_PLUGIN_SERVER_ADDRESS') ?: $this->getDefaultAddress());
        $this->serverAddress = $this->normalizeServerAddress($address);
    }

    private function getDefaultAddress(): string {
        if (PHP_OS_FAMILY === 'Windows') {
            return 'unix://C:/temp/dragonfly_plugin.sock';
        }
        // PHP gRPC extension format for Unix sockets
        return 'unix:/tmp/dragonfly_plugin.sock';
    }

    private function normalizeServerAddress(string $address): string {
        // Handle bare Unix socket paths: "/tmp/dragonfly_plugin.sock" -> "unix:/tmp/dragonfly_plugin.sock"
        if ($address !== '' && $address[0] === '/') {
            return 'unix:' . $address;
        }
        // Normalize triple-slash form to single-slash: "unix:///path" -> "unix:/path"
        $normalized = preg_replace('#^unix:///#', 'unix:/', $address);
        return $normalized ?? $address;
    }

    // Lifecycle hooks
    public function onEnable(): void {}
    public function onDisable(): void {}

    // Registration APIs
    public function subscribe(array $eventTypes): void {
        $this->subscriptions = array_values(array_unique($eventTypes));
    }

    public function addEventHandler(int $eventType, callable $handler): void {
        if (!isset($this->handlers[$eventType])) {
            $this->handlers[$eventType] = [];
        }
        $this->handlers[$eventType][] = $handler;
    }

    /**
     * Register many handlers at once.
     * Keys must be int EventType values (e.g. EventType::PLAYER_JOIN).
     *
     * Handlers receive (string $eventId, EventEnvelope $event).
     */
    public function registerHandlers(array $map): void {
        foreach ($map as $key => $handler) {
            if (is_int($key)) {
                $this->addEventHandler($key, $handler);
            } else {
                throw new \InvalidArgumentException('Handler map keys must be int EventType values.');
            }
        }
    }

    /**
     * Subscribe to the set of types that have handlers registered.
     */
    public function subscribeToRegisteredHandlers(): void {
        $types = [];
        foreach ($this->handlers as $type => $_) {
            if (is_int($type)) {
                $types[] = $type;
            }
        }
        if (!empty($types)) {
            $this->subscriptions = array_values(array_unique($types));
        }
    }

    public function registerCommand(string $name, string $description): void {
        $this->commandSpecs[] = ['name' => $name, 'description' => $description];
    }

    /**
     * Register a listener object.
     * Public, non-static methods with:
     *  - first parameter typed to a payload class under \Df\Plugin\... ending with "Event"
     *  - optional second parameter typed to HandlerContext
     * are auto-registered. Method names are arbitrary.
     *
     * The handler is invoked as either:
     *  - (TypedPayload $payload)
     *  - (TypedPayload $payload, HandlerContext $ctx)
     *
     * The context auto-ACKs if the handler returns without respond/cancel.
     */
    public function registerListener(object $listener): void {
        if (!$listener instanceof Listener) {
            throw new \InvalidArgumentException('Listener must implement ' . Listener::class);
        }

        $ref = new ReflectionClass($listener);
        foreach ($ref->getMethods(ReflectionMethod::IS_PUBLIC) as $method) {
            if ($method->isStatic() || $method->isConstructor() || $method->isDestructor()) {
                continue;
            }
            $params = $method->getParameters();
            if (count($params) < 1) {
                continue;
            }
            $param = $params[0];
            $type = $param->getType();
            if (!$type instanceof ReflectionNamedType || $type->isBuiltin()) {
                continue;
            }
            $paramClass = $type->getName();
            $binding = $this->resolveEventBinding($paramClass);
            if ($binding === null) {
                continue;
            }

            $eventType = $binding['type'];
            $getter = $binding['getter'];
            $methodName = $method->getName();

            $wantsContext = $method->getNumberOfParameters() >= 2;
            $this->addEventHandler($eventType, function (string $eventId, EventEnvelope $event) use ($listener, $methodName, $getter, $wantsContext): void {
                $payload = $event->{$getter}();
                $ctx = new EventContext($this->pluginId, $eventId, $this->sender, $event->getExpectsResponse());
                try {
                    if ($wantsContext) {
                        $listener->{$methodName}($payload, $ctx);
                    } else {
                        $listener->{$methodName}($payload);
                    }
                } catch (\Throwable $e) {
                    fwrite(STDERR, "[php] listener error: {$e->getMessage()}\n");
                } finally {
                    $ctx->ackIfUnhandled();
                }
            });
        }
    }

    /**
     * Resolve event type constant and Event getter name from a payload FQCN.
     * Example: \Df\Plugin\PlayerJoinEvent -> ['type' => EventType::PLAYER_JOIN, 'getter' => 'getPlayerJoin']
     *
     * @return array{type:int,getter:string}|null
     */
    private function resolveEventBinding(string $payloadFqcn): ?array {
        if (!str_starts_with($payloadFqcn, 'Df\\Plugin\\') || !str_ends_with($payloadFqcn, 'Event')) {
            return null;
        }
        $short = ($pos = strrpos($payloadFqcn, '\\')) !== false ? substr($payloadFqcn, $pos + 1) : $payloadFqcn;
        $base = substr($short, 0, -strlen('Event'));
        if ($base === '') {
            return null;
        }
        $getter = 'get' . $base;
        $constName = strtoupper(preg_replace('/(?<!^)[A-Z]/', '_$0', $base));
        $constFq = 'Df\\Plugin\\EventType::' . $constName;
        if (!defined($constFq)) {
            return null;
        }
        /** @var int $type */
        $type = constant($constFq);
        return ['type' => $type, 'getter' => $getter];
    }

    // Action helpers moved to StreamSender and HandlerContext

    // Runner
    public function run(): void {
        fwrite(STDOUT, "[php] connecting to {$this->serverAddress}...\n");

        $this->client = new PluginClient($this->serverAddress, [
            'credentials' => ChannelCredentials::createInsecure(),
        ]);
        $this->call = $this->client->EventStream();
        $this->sender = new StreamSender($this->call);
        $this->running = true;

        // Allow plugin to register handlers/subscriptions/commands
        $this->onEnable();

        // Defaults if not set
        if (empty($this->subscriptions)) {
            // Prefer subscriptions matching registered handlers if present.
            $this->subscribeToRegisteredHandlers();
        }

        // Handshake
        fwrite(STDOUT, "[php] connected, sending handshake\n");
        $hello = new PluginToHost();
        $hello->setPluginId($this->pluginId);
        $pluginHello = new PluginHello();
        $pluginHello->setName($this->name);
        $pluginHello->setVersion($this->version);
        $pluginHello->setApiVersion($this->apiVersion);
        if (!empty($this->commandSpecs)) {
            $cmds = [];
            foreach ($this->commandSpecs as $spec) {
                $c = new CommandSpec();
                $c->setName($spec['name']);
                $c->setDescription($spec['description']);
                $cmds[] = $c;
            }
            $pluginHello->setCommands($cmds);
        }
        $hello->setHello($pluginHello);
        $this->sender->enqueue($hello);

        // Subscribe
        $subscribeMsg = new PluginToHost();
        $subscribeMsg->setPluginId($this->pluginId);
        $subscribe = new EventSubscribe();
        $subscribe->setEvents($this->subscriptions);
        $subscribeMsg->setSubscribe($subscribe);
        $this->sender->enqueue($subscribeMsg);

        try {
            while ($this->running) {
                $message = $this->call->read();
                if ($message === null) {
                    $status = $this->call->getStatus();
                    fwrite(STDOUT, "[php] stream closed - status: code=" . $status->code . " details=" . $status->details . "\n");
                    $this->running = false;
                    break;
                }

                if ($message->hasHello()) {
                    $hostHello = $message->getHello();
                    fwrite(STDOUT, "[php] host hello api=" . $hostHello->getApiVersion() . "\n");
                    if ($hostHello->getApiVersion() !== $this->apiVersion) {
                        fwrite(STDOUT, "[php] WARNING: API version mismatch (host={$hostHello->getApiVersion()}, plugin=" . $this->apiVersion . ")\n");
                    }
                    continue;
                }

                if ($message->hasEvent()) {
                    $event = $message->getEvent();
                    $eventId = $event->getEventId();
                    $type = $event->getType();

                    if (isset($this->handlers[$type])) {
                        foreach ($this->handlers[$type] as $handler) {
                            try {
                                $handler($eventId, $event);
                            } catch (\Throwable $e) {
                                fwrite(STDERR, "[php] handler error: {$e->getMessage()}\n");
                            }
                        }
                        continue;
                    }

                    // Default ack when unhandled
                    (new EventContext($this->pluginId, $eventId, $this->sender, $event->getExpectsResponse()))->ackIfUnhandled();
                    continue;
                }

                if ($message->hasShutdown()) {
                    fwrite(STDOUT, "[php] shutdown received\n");
                    $this->running = false;
                    continue;
                }
            }
        } finally {
            try {
                $this->onDisable();
            } catch (\Throwable $e) {
                fwrite(STDERR, "[php] onDisable error: {$e->getMessage()}\n");
            }
            $this->call->writesDone();
            fwrite(STDOUT, "[php] client completed\n");
            fwrite(STDOUT, "[php] connection closing\n");
        }
    }

}


