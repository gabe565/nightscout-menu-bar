package socket

import (
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/nightscout/testproperties"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()
	socket := New(config.New(), nil)
	require.NotNil(t, socket)
	assert.NotNil(t, socket.config)
	assert.Equal(t, Path(), socket.path)
}

func TestPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Setenv("TMP", `C:\Temp`)
		assert.Equal(t, `C:\Temp\nightscout.sock`, Path())
		return
	}

	name := "nightscout-" + strconv.Itoa(os.Getuid()) + ".sock"
	t.Setenv("TMPDIR", "/tmp/test")
	t.Setenv("XDG_RUNTIME_DIR", "")
	assert.Equal(t, "/tmp/test/"+name, Path())

	if runtime.GOOS == "linux" {
		t.Setenv("XDG_RUNTIME_DIR", "/run/user/1000")
		assert.Equal(t, "/run/user/1000/"+name, Path())
	}
}

func newTestSocket(t *testing.T, enabled bool, onCommand func(string)) *Socket {
	t.Helper()
	conf := config.New()
	data := conf.Data()
	data.Socket.Enabled = enabled
	conf = config.New(config.WithData(data))

	// t.TempDir can exceed the socket path length limit on macOS
	dir, err := os.MkdirTemp("", "ns") //nolint:usetesting
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	socket := New(conf, onCommand)
	socket.path = filepath.Join(dir, "nightscout.sock")
	t.Cleanup(func() {
		_ = socket.Close()
	})
	return socket
}

func read(t *testing.T, path string) string {
	t.Helper()
	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: path, Net: "unix"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	b, err := io.ReadAll(conn)
	require.NoError(t, err)
	return string(b)
}

func TestSocket(t *testing.T) {
	t.Parallel()

	socket := newTestSocket(t, true, nil)
	require.NoError(t, socket.Listen())
	socket.Write(testproperties.Properties)

	assert.Equal(t, socket.Format(testproperties.Properties), read(t, socket.path))

	if runtime.GOOS != "windows" {
		info, err := os.Stat(socket.path)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0o700), info.Mode().Perm())
	}

	require.NoError(t, socket.Close())
	assert.NoFileExists(t, socket.path)
}

func TestSocket_Disabled(t *testing.T) {
	t.Parallel()

	socket := newTestSocket(t, false, nil)
	require.NoError(t, socket.Listen())
	socket.Write(testproperties.Properties)

	assert.Equal(t, "\n", read(t, socket.path))
}

func TestSocket_SecondInstance(t *testing.T) {
	t.Parallel()

	commands := make(chan string, 1)
	first := newTestSocket(t, true, func(cmd string) { commands <- cmd })
	require.NoError(t, first.Listen())

	second := newTestSocket(t, true, nil)
	second.path = first.path
	require.ErrorIs(t, second.Listen(), ErrRunning)

	select {
	case cmd := <-commands:
		assert.Equal(t, CommandSettings, cmd)
	case <-time.After(time.Second):
		t.Fatal("command was not received")
	}

	// The second instance must not remove the first's socket
	assert.FileExists(t, first.path)
}

func TestSocket_Stale(t *testing.T) {
	t.Parallel()

	socket := newTestSocket(t, true, nil)

	// Leave a socket file behind with nothing listening
	stale, err := net.ListenUnix("unix", &net.UnixAddr{Name: socket.path, Net: "unix"})
	require.NoError(t, err)
	stale.SetUnlinkOnClose(false)
	require.NoError(t, stale.Close())
	_, err = os.Stat(socket.path)
	require.NoError(t, err)

	require.NoError(t, socket.Listen())
	socket.Write(testproperties.Properties)
	assert.Equal(t, socket.Format(testproperties.Properties), read(t, socket.path))
}

func TestSocket_Format(t *testing.T) {
	t.Parallel()
	conf := config.New()
	socket := New(conf, nil)

	timeAgo := testproperties.Properties.Bgnow.Mills.Relative(conf.Data().Advanced.RoundAge)
	relative := strconv.Itoa(int(time.Since(testproperties.Properties.Bgnow.Mills.Time).Seconds()))
	assert.Equal(t, "123,→,+1,"+timeAgo+","+relative+"\n", socket.Format(testproperties.Properties))
}
