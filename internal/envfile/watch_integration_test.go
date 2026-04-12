package envfile_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/envfile"
)

func TestWatch_Integration_WriteThenRead(t *testing.T)"=init\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var mu sync.Mutex
	var snapshots []map[string]string

	opts := envfile.DefaultWatchOptions()
	opts.Interval = 40 * time.Millisecond
	opts.OnChange = func(_ string, kv map[string]string) {
		mu.Lock()
		cp := make(map[string]string, len(kv))
		for k, v := range kv {
			cp[k] = v
		}
		snapshots = append(snapshots, cp)
		mu.Unlock()
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		envfile.Watch(ctx, p, opts) //nolint:errcheck
		close(done)
	}()

	time.Sleep(80 * time.Millisecond)
	os.WriteFile(p, []byte("STAGE=alpha\n"), 0o600) //nolint:errcheck
	time.Sleep(80 * time.Millisecond)
	os.WriteFile(p, []byte("STAGE=beta\n"), 0o600) //nolint:errcheck
	time.Sleep(80 * time.Millisecond)
	cancel()
	<-done

	mu.Lock()
	defer mu.Unlock()
	if len(snapshots) < 2 {
		t.Fatalf("expected at least 2 change events, got %d", len(snapshots))
	}
	last := snapshots[len(snapshots)-1]
	if last["STAGE"] != "beta" {
		t.Errorf("last snapshot STAGE = %q, want beta", last["STAGE"])
	}
}
