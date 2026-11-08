package secret

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/util"
	"github.com/zalando/go-keyring"
)

const (
	service = "com.gabe565.nightscout-menu-bar"
	// stateFile changes whenever a secret does, so other processes can watch for changes.
	// When no keyring is available, it also holds the secrets.
	stateFile      = "secrets.json"
	storageKeyring = "keyring"
	storageFile    = "file"
	probeKey       = "probe"
)

var ErrNotFound = errors.New("secret not found")

type Backend interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

// Keyring stores secrets in the system keyring:
// Keychain on macOS, Credential Manager on Windows, or the Secret Service on Linux.
type Keyring struct{}

func (Keyring) Get(key string) (string, error) {
	v, err := keyring.Get(service, key)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", ErrNotFound
	}
	return v, err
}

func (Keyring) Set(key, value string) error {
	return keyring.Set(service, key, value)
}

func (Keyring) Delete(key string) error {
	if err := keyring.Delete(service, key); err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return err
	}
	return nil
}

// Memory stores secrets in memory, for tests.
type Memory struct {
	mu     sync.Mutex
	values map[string]string
}

func (m *Memory) Get(key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.values[key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

func (m *Memory) Set(key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.values == nil {
		m.values = make(map[string]string)
	}
	m.values[key] = value
	return nil
}

func (m *Memory) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.values, key)
	return nil
}

func New(dir string) *Store {
	return NewWithBackend(dir, Keyring{})
}

func NewWithBackend(dir string, backend Backend) *Store {
	return &Store{path: filepath.Join(dir, stateFile), backend: backend}
}

// Store saves secrets with a [Backend], falling back to a file readable only by
// the current user if the backend is unavailable (e.g. no Secret Service on Linux).
type Store struct {
	mu       sync.Mutex
	path     string
	backend  Backend
	checked  bool
	fallback bool
}

type state struct {
	Storage string            `json:"storage"`
	Updated time.Time         `json:"updated"`
	Values  map[string]string `json:"values,omitempty"`
}

// checkLocked switches to the file fallback if the backend doesn't work.
func (s *Store) checkLocked() {
	if s.checked {
		return
	}
	s.checked = true
	if _, err := s.backend.Get(probeKey); err != nil && !errors.Is(err, ErrNotFound) {
		slog.Warn("System keyring is unavailable; storing secrets in a file instead", "path", s.path, "error", err)
		s.fallback = true
	}
}

func (s *Store) readState() (state, error) {
	var st state
	b, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return st, nil
		}
		return st, err
	}
	err = json.Unmarshal(b, &st)
	return st, err
}

// update rewrites the state file so watchers see the change, applying fn to
// the file's values when using the fallback.
func (s *Store) update(fn func(values map[string]string)) error {
	st := state{Storage: storageKeyring, Updated: time.Now()}
	if s.fallback {
		prev, err := s.readState()
		if err != nil {
			return err
		}
		st.Storage = storageFile
		st.Values = prev.Values
		if st.Values == nil {
			st.Values = make(map[string]string)
		}
		fn(st.Values)
	}
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	return util.WriteFileAtomic(s.path, b, 0o600)
}

func (s *Store) Get(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checkLocked()

	if !s.fallback {
		return s.backend.Get(key)
	}
	st, err := s.readState()
	if err != nil {
		return "", err
	}
	v, ok := st.Values[key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

func (s *Store) Set(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checkLocked()

	if !s.fallback {
		if err := s.backend.Set(key, value); err != nil {
			return err
		}
	}
	return s.update(func(values map[string]string) { values[key] = value })
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checkLocked()

	if !s.fallback {
		if err := s.backend.Delete(key); err != nil {
			return err
		}
	}
	return s.update(func(values map[string]string) { delete(values, key) })
}

// Watch calls onChange when another process changes a secret.
func (s *Store) Watch(ctx context.Context, onChange func()) error {
	return util.WatchFile(ctx, s.path, onChange)
}
