package envfile

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchOptions configures the file watcher behaviour.
type WatchOptions struct {
	// Interval between polling checks.
	Interval time.Duration
	// OnChange is called whenever the file content changes.
	OnChange func(path string, current map[string]string)
	// OnError is called on read or parse errors.
	OnError func(path string, err error)
}

// DefaultWatchOptions returns sensible defaults.
func DefaultWatchOptions() WatchOptions {
	return WatchOptions{
		Interval: 2 * time.Second,
		OnChange: func(_ string, _ map[string]string) {},
		OnError:  func(_ string, _ error) {},
	}
}

// Watch polls path at opts.Interval and fires opts.OnChange when the file
// content has changed. It blocks until ctx is cancelled.
func Watch(ctx context.Context, path string, opts WatchOptions) error {
	if opts.Interval <= 0 {
		opts.Interval = DefaultWatchOptions().Interval
	}

	lastHash, _ := fileHash(path)

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			h, err := fileHash(path)
			if err != nil {
				opts.OnError(path, err)
				continue
			}
			if h == lastHash {
				continue
			}
			lastHash = h
			kv, err := Parse(path)
			if err != nil {
				opts.OnError(path, err)
				continue
			}
			opts.OnChange(path, kv)
		}
	}
}

func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
