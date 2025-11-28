package plugin

import (
	"testing"
	"time"

	pb "github.com/secmc/plugin/proto/generated/go"
)

func TestEventBatching(t *testing.T) {
	// Manually construct pluginProcess to isolate batchSendLoop
	p := &pluginProcess{
		id:     "test-plugin",
		sendCh: make(chan *pb.HostToPlugin, 10),
		done:   make(chan struct{}),
	}
	p.connected.Store(true) // Required for queueEvent to work

	p.wg.Add(1)
	go p.batchSendLoop()

	// Queue multiple events rapidly
	count := 5
	for i := 0; i < count; i++ {
		p.queueEvent(&pb.EventEnvelope{
			EventId: "evt-" + string(rune(i)),
		})
	}

	// Wait for batch (ticker is 5ms)
	select {
	case msg := <-p.sendCh:
		batch := msg.GetEvents()
		if batch == nil {
			t.Fatalf("Expected Events payload, got %T", msg.Payload)
		}
		if len(batch.Events) != count {
			t.Errorf("Expected %d batched events, got %d", count, len(batch.Events))
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for batch")
	}

	// Verify buffer is cleared
	p.eventBufferMu.Lock()
	if len(p.eventBuffer) != 0 {
		t.Error("Event buffer should be empty after sending")
	}
	p.eventBufferMu.Unlock()

	// Cleanup
	close(p.done)
	p.wg.Wait()
}
