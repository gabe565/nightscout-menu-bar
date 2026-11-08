package secret

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errUnavailable = errors.New("no keyring")

// unavailable is a backend that always fails, like Linux without a Secret Service.
type unavailable struct{}

func (unavailable) Get(string) (string, error) { return "", errUnavailable }
func (unavailable) Set(string, string) error   { return errUnavailable }
func (unavailable) Delete(string) error        { return errUnavailable }

func readState(t *testing.T, dir string) state {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, stateFile))
	require.NoError(t, err)
	var st state
	require.NoError(t, json.Unmarshal(b, &st))
	return st
}

func TestStore_Backend(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	backend := &Memory{}
	s := NewWithBackend(dir, backend)

	_, err := s.Get("token")
	require.ErrorIs(t, err, ErrNotFound)

	require.NoError(t, s.Set("token", "abc"))
	v, err := s.Get("token")
	require.NoError(t, err)
	assert.Equal(t, "abc", v)

	// Secrets stay out of the state file when the backend works
	st := readState(t, dir)
	assert.Equal(t, storageKeyring, st.Storage)
	assert.Empty(t, st.Values)

	require.NoError(t, s.Delete("token"))
	_, err = backend.Get("token")
	require.ErrorIs(t, err, ErrNotFound)
	assert.True(t, readState(t, dir).Updated.After(st.Updated))
}

func TestStore_Fallback(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	s := NewWithBackend(dir, unavailable{})

	require.NoError(t, s.Set("token", "abc"))
	v, err := s.Get("token")
	require.NoError(t, err)
	assert.Equal(t, "abc", v)

	stat, err := os.Stat(filepath.Join(dir, stateFile))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), stat.Mode().Perm())
	st := readState(t, dir)
	assert.Equal(t, storageFile, st.Storage)
	assert.Equal(t, map[string]string{"token": "abc"}, st.Values)

	require.NoError(t, s.Delete("token"))
	_, err = s.Get("token")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestStore_Watch(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	backend := &Memory{}
	s := NewWithBackend(dir, backend)
	other := NewWithBackend(dir, backend)

	changed := make(chan struct{}, 1)
	require.NoError(t, s.Watch(t.Context(), func() { changed <- struct{}{} }))

	require.NoError(t, other.Set("token", "abc"))
	select {
	case <-changed:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for change")
	}
}
