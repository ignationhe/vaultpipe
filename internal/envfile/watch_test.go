package envfile

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func writeTempEnvWatch(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestWatch_DetectsChange(t *testing.T) {
	p := writeTempEnvWatch(t, "KEY=old\n")

	var mu sync.Mutex
	var got map[string]string

	opts := DefaultWatchOptions()
	opts.Interval = 50 * time.Millisecond
	opts.OnChange = func(_ string, kv map[string]string) {
		mu.Lock()
		got = kv
		mu.Unlock()
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Watch(ctx, p, opts) }()

	time.Sleep(120 * time.Millisecond)
	if err := os.WriteFile(p, []byte("KEY=new\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)
	cancel()
	<-done

	mu.Lock()
	defer mu.Unlock()
	if got == nil {
		t.Fatal("OnChange was never called")
	}
	if got["KEY"] != "new" {
		t.Errorf("expected KEY=new, got %q", got["KEY"])
	}
}

func TestWatch_NoCallWhenUnchanged(t *testing.T) {
	p := writeTempEnvWatch(t, "KEY=stable\n")

	callCount := 0
	opts := DefaultWatchOptions()
	opts.Interval = 50 * time.Millisecond
	opts.OnChange = func(_ string, _ map[string]string) { callCount++ }

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()
	Watch(ctx, p, opts) //nolint:errcheck

	if callCount != 0 {
		t.Errorf("expected 0 calls, got %d", callCount)
	}
}

func TestWatch_CallsOnErrorForMissingFile(t *testing.T) {
	opts := DefaultWatchOptions()
	opts.Interval = 50 * time.Millisecond
	errCount := 0
	opts.OnError = func(_ string, _ error) { errCount++ }

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()
	Watch(ctx, "/nonexistent/.env", opts) //nolint:errcheck

	if errCount == 0 {
		t.Error("expected at least one error callback")
	}
}

func TestWatch_CancelStopsLoop(t *testing.T) {
	p := writeTempEnvWatch(t, "A=1\n")

	opts := DefaultWatchOptions()
	opts.Interval = 30 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := Watch(ctx, p, opts)
	if err == nil {
		t.Error("expected context error, got nil")
	}
}
