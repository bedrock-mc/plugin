package plugin

import (
	"io"
	"log/slog"
	"testing"
	"time"

	pb "github.com/bedrock-mc/plugin/proto/generated/go"
)

func TestSubscriptionModes(t *testing.T) {
	p := new(pluginProcess)
	p.connected.Store(true)
	p.updateSubscriptions(
		[]pb.EventType{pb.EventType_CHAT, pb.EventType_PLAYER_MOVE},
		[]pb.EventType{pb.EventType_PLAYER_MOVE, pb.EventType_PLAYER_JUMP},
	)

	tests := []struct {
		event pb.EventType
		want  subscriptionMode
	}{
		{event: pb.EventType_CHAT, want: subscriptionBlocking},
		{event: pb.EventType_PLAYER_MOVE, want: subscriptionBlocking},
		{event: pb.EventType_PLAYER_JUMP, want: subscriptionObserve},
		{event: pb.EventType_PLAYER_QUIT, want: subscriptionNone},
	}
	for _, tt := range tests {
		if got := p.SubscriptionMode(tt.event); got != tt.want {
			t.Fatalf("SubscriptionMode(%v) = %v, want %v", tt.event, got, tt.want)
		}
	}
	if !p.CanReceiveEvents() {
		t.Fatal("connected, ready plugin should receive events")
	}
	p.connected.Store(false)
	if p.CanReceiveEvents() {
		t.Fatal("disconnected plugin should not receive events")
	}
}

func TestObserverDispatchDoesNotExpectResponse(t *testing.T) {
	p := &pluginProcess{id: "observer", sendCh: make(chan *pb.HostToPlugin, 1)}
	p.connected.Store(true)
	p.updateSubscriptions(nil, []pb.EventType{pb.EventType_PLAYER_MOVE})
	m := &Manager{
		plugins: map[string]*pluginProcess{"observer": p},
		log:     slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	m.emitCancellable(nil, &pb.EventEnvelope{Type: pb.EventType_PLAYER_MOVE})

	select {
	case msg := <-p.sendCh:
		if msg.GetEvent().GetExpectsResponse() {
			t.Fatal("observer event unexpectedly requires a response")
		}
	case <-time.After(time.Second):
		t.Fatal("observer did not receive event")
	}

	p.connected.Store(false)
	m.emitCancellable(nil, &pb.EventEnvelope{Type: pb.EventType_PLAYER_MOVE})
	select {
	case <-p.sendCh:
		t.Fatal("disconnected observer received event")
	default:
	}
}

func TestAllSubscriptionModes(t *testing.T) {
	p := new(pluginProcess)
	p.updateSubscriptions(nil, []pb.EventType{pb.EventType_EVENT_TYPE_ALL})
	if got := p.SubscriptionMode(pb.EventType_CHAT); got != subscriptionObserve {
		t.Fatalf("observer all mode = %v, want %v", got, subscriptionObserve)
	}

	p.updateSubscriptions([]pb.EventType{pb.EventType_EVENT_TYPE_ALL}, []pb.EventType{pb.EventType_EVENT_TYPE_ALL})
	if got := p.SubscriptionMode(pb.EventType_CHAT); got != subscriptionBlocking {
		t.Fatalf("blocking all precedence = %v, want %v", got, subscriptionBlocking)
	}
}
