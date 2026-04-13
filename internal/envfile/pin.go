package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// PinEntry records a pinned value for a key at a specific version/timestamp.
type PinEntry struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	PinnedAt  time.Time `json:"pinned_at"`
	Comment   string    `json:"comment,omitempty"`
}

// PinFile is the on-disk representation of all pinned keys.
type PinFile struct {
	Pins []PinEntry `json:"pins"`
}

// DefaultPinOptions returns sensible defaults.
func DefaultPinOptions() PinOptions {
	return PinOptions{
		Overwrite: false,
	}
}

// PinOptions controls Pin behaviour.
type PinOptions struct {
	// Overwrite allows an existing pin for the same key to be replaced.
	Overwrite bool
	// Comment is an optional annotation stored alongside the pin.
	Comment string
}

// Pin records the current values of the given keys from env into the pin file
// located at pinPath. Keys not found in env are silently skipped.
func Pin(env map[string]string, keys []string, pinPath string, opts PinOptions) (PinFile, error) {
	pf, err := loadPinFile(pinPath)
	if err != nil {
		return PinFile{}, fmt.Errorf("pin: load: %w", err)
	}

	index := buildPinIndex(pf)

	for _, k := range keys {
		v, ok := env[k]
		if !ok {
			continue
		}
		if _, exists := index[k]; exists && !opts.Overwrite {
			continue
		}
		entry := PinEntry{
			Key:      k,
			Value:    v,
			PinnedAt: time.Now().UTC(),
			Comment:  opts.Comment,
		}
		if idx, exists := index[k]; exists {
			pf.Pins[idx] = entry
		} else {
			pf.Pins = append(pf.Pins, entry)
		}
	}

	if err := savePinFile(pinPath, pf); err != nil {
		return PinFile{}, fmt.Errorf("pin: save: %w", err)
	}
	return pf, nil
}

// ApplyPins overlays pinned values onto env, returning a new map.
func ApplyPins(env map[string]string, pinPath string) (map[string]string, error) {
	pf, err := loadPinFile(pinPath)
	if err != nil {
		return nil, fmt.Errorf("apply pins: %w", err)
	}
	out := copyMapPin(env)
	for _, p := range pf.Pins {
		out[p.Key] = p.Value
	}
	return out, nil
}

func loadPinFile(path string) (PinFile, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return PinFile{}, nil
	}
	if err != nil {
		return PinFile{}, err
	}
	var pf PinFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return PinFile{}, err
	}
	return pf, nil
}

func savePinFile(path string, pf PinFile) error {
	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func buildPinIndex(pf PinFile) map[string]int {
	m := make(map[string]int, len(pf.Pins))
	for i, p := range pf.Pins {
		m[p.Key] = i
	}
	return m
}

func copyMapPin(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}
