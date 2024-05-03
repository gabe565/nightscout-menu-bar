package localfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gabe565/nightscout-menu-bar/internal/config"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout/testproperties"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()
	localfile := New(config.NewDefault())
	require.NotNil(t, localfile)
	assert.NotNil(t, localfile.config)
}

func TestLocalFile(t *testing.T) {
	t.Parallel()

	temp, err := os.MkdirTemp("", "nightscout-")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(temp)
	})

	conf := config.NewDefault()
	conf.LocalFile.Enabled = true
	conf.LocalFile.Path = filepath.Join(temp, "nightscout.csv")

	localfile := New(conf)
	require.NotNil(t, localfile)
	require.Equal(t, conf.LocalFile.Path, localfile.path)

	require.NoError(t, localfile.Write(testproperties.Properties))

	contents, err := os.ReadFile(conf.LocalFile.Path)
	require.NoError(t, err)

	assert.Equal(t, "123,→,+1,1664764258\n", string(contents))

	require.NoError(t, localfile.Cleanup())
	assert.NoFileExists(t, conf.LocalFile.Path)
}
