package settings

import (
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/gabe565/nightscout-menu-bar/internal/autostart"
	"github.com/gabe565/nightscout-menu-bar/internal/fyneutil"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
	"github.com/ncruces/zenity"
)

const (
	URLKey   = "url"
	TokenKey = "token"
	UnitsKey = "units"

	DynamicIconEnabledKey   = "dynamic_icon_enabled"
	DynamicIconFontColorKey = "dynamic_icon_font_color"
	DynamicIconFontPathKey  = "dynamic_icon_font_path"

	LocalEnabledKey  = "local_file_enabled"
	LocalPathKey     = "local_file_path"
	LocalPathDefault = "$TMPDIR/nightscout.csv"

	FetchDelayKey           = "fetch_delay"
	FetchDelayDefault       = 30 * time.Second
	FallbackIntervalKey     = "fallback_interval"
	FallbackIntervalDefault = 30 * time.Second
)

func OpenSettings(app fyne.App) func() {
	var settings fyne.Window
	return func() {
		if settings == nil {
			settings = loadSettings(app)
			settings.SetOnClosed(func() {
				settings = nil
				runtime.GC()
			})
		} else {
			settings.RequestFocus()
		}
	}
}

func loadSettings(app fyne.App) fyne.Window { //nolint:ireturn
	window := app.NewWindow("Nightscout Settings")
	window.SetIcon(assets.NightscoutResource)

	prefs := app.Preferences()

	url := widget.NewEntry()
	url.Text = prefs.String(URLKey)

	token := widget.NewPasswordEntry()
	token.Text = prefs.String(TokenKey)

	units := widget.NewRadioGroup([]string{UnitMgdl.String(), UnitMmol.String()}, nil)
	units.Selected = prefs.StringWithFallback(UnitsKey, UnitMgdl.String())

	start := widget.NewCheck("", nil)
	start.Checked, _ = autostart.IsEnabled()

	writeLocal := widget.NewCheck("", nil)
	writeLocal.Checked = prefs.Bool(LocalEnabledKey)

	localPath := widget.NewEntry()
	localPath.Text = prefs.StringWithFallback(LocalPathKey, LocalPathDefault)

	dynamicIcon := widget.NewCheck("", nil)
	dynamicIcon.Checked = prefs.BoolWithFallback(DynamicIconEnabledKey, DynamicIconDefault)

	var dynamicIconColorValue util.HexColor
	if err := dynamicIconColorValue.UnmarshalText([]byte(prefs.String(DynamicIconFontColorKey))); err != nil {
		dynamicIconColorValue = util.White()
	}

	var dynamicIconColor *widget.Button
	//nolint:gosec
	dynamicIconColor = widget.NewButtonWithIcon("Open Color Picker",
		fyneutil.ColorImage(dynamicIconColorValue),
		func() {
			newColor, err := zenity.SelectColor(
				zenity.Title("Dynamic Icon Color"),
				zenity.Color(dynamicIconColorValue),
				zenity.ShowPalette(),
			)
			if err == nil {
				r, g, b, a := newColor.RGBA()
				dynamicIconColorValue.R = uint8(r)
				dynamicIconColorValue.G = uint8(g)
				dynamicIconColorValue.B = uint8(b)
				dynamicIconColorValue.A = uint8(a)
				dynamicIconColor.SetIcon(fyneutil.ColorImage(dynamicIconColorValue))
			}
			window.RequestFocus()
		},
	)

	dynamicIconFontPath := prefs.String(DynamicIconFontPathKey)
	dynamicIconFont := widget.NewButton("Dynamic Icon Font", func() {
		newPath, err := zenity.SelectFile(
			zenity.Title("Dynamic Icon Font"),
			zenity.Filename(dynamicIconFontPath),
			zenity.FileFilter{
				Name:     "TrueType Font",
				Patterns: []string{".ttf"},
				CaseFold: true,
			},
		)
		if err != nil {
			dynamicIconFontPath = ""
		} else {
			dynamicIconFontPath = newPath
		}
		window.RequestFocus()
	})

	fetchDelay := widget.NewEntry()
	fetchDelay.Text = prefs.StringWithFallback(FetchDelayKey, FetchDelayDefault.String())
	fetchDelay.Validator = fyneutil.ValidateDuration(true)

	fallbackInterval := widget.NewEntry()
	fallbackInterval.Text = prefs.StringWithFallback(FallbackIntervalKey, FallbackIntervalDefault.String())
	fallbackInterval.Validator = fyneutil.ValidateDuration(true)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "URL", Widget: url},
			{Text: "Token", Widget: token},
			{Text: "Unit", Widget: units},
			{Text: "Start on login", Widget: start},
			{Text: "Write local file", Widget: writeLocal},
			{Text: "Local file path", Widget: localPath},
			{Text: "Dynamic icon", Widget: dynamicIcon},
			{Text: "Dynamic icon color", Widget: dynamicIconColor},
			{Text: "Dynamic icon font", Widget: dynamicIconFont},
			{Text: "Fetch delay", Widget: fetchDelay},
			{Text: "Fallback interval", Widget: fallbackInterval},
		},
		OnSubmit: func() {
			prefs.SetString(URLKey, url.Text)
			prefs.SetString(TokenKey, token.Text)
			prefs.SetString(UnitsKey, units.Selected)
			_ = autostart.Set(start.Checked)
			prefs.SetBool(LocalEnabledKey, writeLocal.Checked)
			prefs.SetString(LocalPathKey, localPath.Text)
			prefs.SetBool(DynamicIconEnabledKey, dynamicIcon.Checked)
			prefs.SetString(DynamicIconFontColorKey, dynamicIconColorValue.String())
			prefs.SetString(DynamicIconFontPathKey, dynamicIconFontPath)
			prefs.SetString(FetchDelayKey, fetchDelay.Text)
			prefs.SetString(FallbackIntervalKey, fallbackInterval.Text)
			window.Close()
		},
	}
	form.SubmitText = "Save"
	window.SetContent(container.NewPadded(form))
	window.Resize(fyne.NewSize(500, 0))
	window.Show()
	return window
}
