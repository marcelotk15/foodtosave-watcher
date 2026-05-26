package watcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type CachedGondola struct {
	GondolaID         string    `json:"gondola_id"`
	Quantity          int       `json:"quantity"`
	BagDescription    string    `json:"bag_description"`
	BagCategory       string    `json:"bag_category"`
	BagPrice          float64   `json:"bag_price"`
	AvailabilityEndAt time.Time `json:"availability_end_at"`
}

type Cache struct {
	mu   sync.Mutex
	data map[string]map[string]CachedGondola
	path string
}

// NewCache loads the previous state from disk (path).
// If the file does not exist or is corrupted, it starts with an empty cache and logs a warning.
func NewCache(path string) *Cache {
	c := &Cache{
		path: path,
		data: make(map[string]map[string]CachedGondola),
	}
	if err := c.load(); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			slog.Warn("cache do disco ignorado; iniciando vazio", "path", path, "err", err)
		}
	} else {
		slog.Info("cache carregado do disco", "path", path)
	}
	return c
}

// Snapshot returns a copy of the gondola map for the merchant (never nil).
func (c *Cache) Snapshot(merchantID string) map[string]CachedGondola {
	c.mu.Lock()
	defer c.mu.Unlock()
	m := c.data[merchantID]
	if len(m) == 0 {
		return map[string]CachedGondola{}
	}
	out := make(map[string]CachedGondola, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// Replace completely replaces the merchant's state with the current snapshot and persists to disk.
func (c *Cache) Replace(merchantID string, snapshot map[string]CachedGondola) error {
	c.mu.Lock()
	m := make(map[string]CachedGondola, len(snapshot))
	for k, v := range snapshot {
		m[k] = v
	}
	c.data[merchantID] = m
	dataCopy := c.copyDataLocked()
	c.mu.Unlock()

	return c.persistData(dataCopy)
}

// copyDataLocked makes a deep copy of data (must be called with the lock obtained).
func (c *Cache) copyDataLocked() map[string]map[string]CachedGondola {
	out := make(map[string]map[string]CachedGondola, len(c.data))
	for mid, gondolas := range c.data {
		g := make(map[string]CachedGondola, len(gondolas))
		for gid, cg := range gondolas {
			g[gid] = cg
		}
		out[mid] = g
	}
	return out
}

// load reads the JSON file from disk and populates data.
func (c *Cache) load() error {
	raw, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}
	var data map[string]map[string]CachedGondola
	if err := json.Unmarshal(raw, &data); err != nil {
		return fmt.Errorf("json inválido em %s: %w", c.path, err)
	}
	if data != nil {
		c.data = data
	}
	return nil
}

// persistData writes data to disk atomically (temp + rename).
func (c *Cache) persistData(data map[string]map[string]CachedGondola) error {
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("serializar cache: %w", err)
	}

	dir := filepath.Dir(c.path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("criar diretório do cache %s: %w", dir, err)
		}
	}

	tmp := c.path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return fmt.Errorf("escrever cache temporário: %w", err)
	}
	if err := os.Rename(tmp, c.path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("renomear cache: %w", err)
	}
	return nil
}
