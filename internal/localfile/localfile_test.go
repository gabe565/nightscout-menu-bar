package localfile

import (
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2/app"
	"github.com/gabe565/nightscout-menu-bar/internal/app/settings"
	"github.com/gabe565/nightscout-menu-bar/internal/nightscout/testproperties"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()
	localfile := New(app.New())
	require.NotNil(t, localfile)
	assert.NotNil(t, localfile.app)
}

func TestLocalFile(t *testing.T) {
	t.Parallel()

	temp, err := os.MkdirTemp("", "nightscout-")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(temp)
	})

	app := app.New()
	app.Preferences().SetBool(settings.LocalEnabledKey, true)
	app.Preferences().SetString(settings.LocalPathKey, filepath.Join(temp, "nightscout.csv"))

	localfile := New(app)
	require.NotNil(t, localfile)
	require.Equal(t, app.Preferences().String(settings.LocalPathKey), localfile.path)

	require.NoError(t, localfile.Write(testproperties.Properties))

	contents, err := os.ReadFile(app.Preferences().String(settings.LocalPathKey))
	require.NoError(t, err)

	assert.Equal(t, "123,→,+1,1664764258\n", string(contents))

	require.NoError(t, localfile.Cleanup())
	assert.NoFileExists(t, app.Preferences().String(settings.LocalPathKey))
}
