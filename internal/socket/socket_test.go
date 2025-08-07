package socket

import (
	"io"
	"net"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/nightscout/testproperties"
	"gabe565.com/nightscout-menu-bar/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()
	socket := New(config.New())
	require.NotNil(t, socket)
	assert.NotNil(t, socket.config)
}

func TestSocket(t *testing.T) {
	t.Parallel()

	temp := t.TempDir()

	conf := config.New(config.WithData(config.Data{
		Socket: config.Socket{
			Enabled: true,
			Path:    filepath.Join(temp, "nightscout.sock"),
		},
	}))

	socket := New(conf)
	require.NotNil(t, socket)
	t.Cleanup(func() {
		_ = socket.Close()
	})

	socket.Write(testproperties.Properties)

	path := util.ResolvePath(conf.Data().Socket.Path)

	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: path, Net: "unix"})
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = conn.Close()
	})

	b, err := io.ReadAll(conn)
	require.NoError(t, err)

	assert.Equal(t, socket.Format(testproperties.Properties), string(b))

	require.NoError(t, socket.Close())
	assert.NoFileExists(t, path)
}

func TestSocket_Format(t *testing.T) {
	conf := config.New()
	socket := New(conf)
	t.Cleanup(func() {
		_ = socket.Close()
	})

	timeAgo := testproperties.Properties.Bgnow.Mills.Relative(conf.Data().Advanced.RoundAge)
	relative := strconv.Itoa(int(time.Since(testproperties.Properties.Bgnow.Mills.Time).Seconds()))
	assert.Equal(t, "123,→,+1,"+timeAgo+","+relative+"\n", socket.Format(testproperties.Properties))
}
