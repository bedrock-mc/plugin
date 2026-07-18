package plugin

import (
	"errors"
	"testing"
	"time"

	pb "github.com/bedrock-mc/plugin/proto/generated/go"
	"github.com/df-mc/dragonfly/server/world"
)

func TestCompleteWorldTask(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantOK  bool
		wantErr string
	}{
		{name: "success", wantOK: true},
		{name: "failure", err: errors.New("world closed"), wantErr: "world closed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pluginProcess{id: "test", sendCh: make(chan *pb.HostToPlugin, 1)}
			p.connected.Store(true)
			m := new(Manager)

			m.completeWorldTask(p, "correlation", world.NewFinishedTask(tt.err))

			select {
			case msg := <-p.sendCh:
				result := msg.GetActionResult()
				if got := result.GetStatus().GetOk(); got != tt.wantOK {
					t.Fatalf("ok = %v, want %v", got, tt.wantOK)
				}
				if got := result.GetStatus().GetError(); got != tt.wantErr {
					t.Fatalf("error = %q, want %q", got, tt.wantErr)
				}
			case <-time.After(time.Second):
				t.Fatal("timed out waiting for task result")
			}
		})
	}
}
