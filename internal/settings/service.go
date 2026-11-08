package settings

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gabe565.com/nightscout-menu-bar/internal/autostart"
	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/dynamicicon"
	"gabe565.com/nightscout-menu-bar/internal/fetch"
	"gabe565.com/nightscout-menu-bar/internal/nightscout"
	"gabe565.com/nightscout-menu-bar/internal/oauth"
	"gabe565.com/utils/colorx"
	"gabe565.com/utils/slogx"
	"github.com/flopp/go-findfont"
	"github.com/skratchdot/open-golang/open"
	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/oauth2"
)

var (
	ErrDuration    = errors.New("must be a duration like 30s or 5m")
	ErrFontSize    = errors.New("max font size must be a positive number")
	ErrUnknownUnit = errors.New("unknown units")
)

const (
	EventNocturneCode   = "nocturne:code"
	EventCloseRequested = "settings:close-requested"
)

type Settings struct {
	URL   string `json:"url"`
	Token string `json:"token"`

	Title          string `json:"title"`
	Units          string `json:"units"`
	ShowArrow      bool   `json:"showArrow"`
	ShowDelta      bool   `json:"showDelta"`
	ShowTimeAgo    bool   `json:"showTimeAgo"`
	StaleThreshold string `json:"staleThreshold"`

	DynamicIcon bool    `json:"dynamicIcon"`
	AutoColor   bool    `json:"autoColor"`
	FontColor   string  `json:"fontColor"`
	FontFile    string  `json:"fontFile"`
	MaxFontSize float64 `json:"maxFontSize"`

	StartOnLogin bool   `json:"startOnLogin"`
	Socket       bool   `json:"socket"`
	SocketPath   string `json:"socketPath"`

	FetchDelay       string `json:"fetchDelay"`
	FallbackInterval string `json:"fallbackInterval"`
	RoundAge         bool   `json:"roundAge"`
	LogLevel         string `json:"logLevel"`
}

type Options struct {
	Units     []string `json:"units"`
	LogLevels []string `json:"logLevels"`
	Platform  string   `json:"platform"`
}

type Preview struct {
	Text      string `json:"text"`
	Stale     bool   `json:"stale"`
	Icon      string `json:"icon,omitempty"`
	IconColor string `json:"iconColor,omitempty"`
	Error     string `json:"error,omitempty"`
	Live      bool   `json:"live"`
}

type NocturneCode struct {
	Code string `json:"code"`
	URL  string `json:"url"`
}

func NewService(conf *config.Config, store *oauth.Store) *Service {
	return &Service{conf: conf, store: store, openURL: open.Run}
}

type Service struct {
	conf    *config.Config
	store   *oauth.Store
	openURL func(string) error

	mu   sync.Mutex
	live *nightscout.Properties

	window     *application.WebviewWindow
	guarded    atomic.Bool
	allowClose atomic.Bool
}

// GuardClose is called once the page has loaded, so the window can still be
// closed if it fails to load.
func (s *Service) GuardClose() {
	s.guarded.Store(true)
}

func (s *Service) Close() {
	s.allowClose.Store(true)
	if s.window != nil {
		s.window.Close()
	}
}

func (s *Service) canClose() bool {
	return !s.guarded.Load() || s.allowClose.Load()
}

func (s *Service) Options() Options {
	return Options{Units: config.UnitStrings(), LogLevels: slogx.LevelStrings(), Platform: runtime.GOOS}
}

func (s *Service) Load() Settings {
	data := s.conf.Data()
	startOnLogin, _ := autostart.IsEnabled()
	return Settings{
		URL:   data.URL,
		Token: data.Token,

		Title:          data.Title,
		Units:          data.Units.String(),
		ShowArrow:      !data.LastReading.HideArrow,
		ShowDelta:      !data.LastReading.HideDelta,
		ShowTimeAgo:    !data.LastReading.HideTimeAgo,
		StaleThreshold: formatDuration(data.LastReading.StaleThreshold),

		DynamicIcon: data.DynamicIcon.Enabled,
		AutoColor:   data.DynamicIcon.AutoColor,
		FontColor:   hexColor(data.DynamicIcon.FontColor),
		FontFile:    data.DynamicIcon.FontFile,
		MaxFontSize: data.DynamicIcon.MaxFontSize,

		StartOnLogin: startOnLogin,
		Socket:       data.Socket.Enabled,
		SocketPath:   data.Socket.Path,

		FetchDelay:       formatDuration(data.Advanced.FetchDelay),
		FallbackInterval: formatDuration(data.Advanced.FallbackInterval),
		RoundAge:         data.Advanced.RoundAge,
		LogLevel:         data.Log.Level.String(),
	}
}

func (s *Service) Save(settings Settings) error {
	data, err := s.toData(settings)
	if err != nil {
		return err
	}

	if enabled, _ := autostart.IsEnabled(); enabled != settings.StartOnLogin {
		if settings.StartOnLogin {
			err = autostart.Enable()
		} else {
			err = autostart.Disable()
		}
		if err != nil {
			return err
		}
	}

	return s.conf.Write(data)
}

func (s *Service) toData(settings Settings) (config.Data, error) {
	data := s.conf.Data()
	data.URL = strings.TrimSpace(settings.URL)
	data.Token = settings.Token

	data.Title = settings.Title
	if err := data.Units.UnmarshalText([]byte(settings.Units)); err != nil {
		return data, fmt.Errorf("%w: %s", ErrUnknownUnit, settings.Units)
	}
	data.LastReading.HideArrow = !settings.ShowArrow
	data.LastReading.HideDelta = !settings.ShowDelta
	data.LastReading.HideTimeAgo = !settings.ShowTimeAgo

	data.DynamicIcon.Enabled = settings.DynamicIcon
	data.DynamicIcon.AutoColor = settings.AutoColor
	if settings.FontColor != "" {
		if err := data.DynamicIcon.FontColor.UnmarshalText([]byte(settings.FontColor)); err != nil {
			return data, fmt.Errorf("font color: %w", err)
		}
	}
	data.DynamicIcon.FontFile = strings.TrimSpace(settings.FontFile)
	if settings.MaxFontSize <= 0 {
		return data, ErrFontSize
	}
	data.DynamicIcon.MaxFontSize = settings.MaxFontSize

	data.Socket.Enabled = settings.Socket
	data.Socket.Path = settings.SocketPath

	data.Advanced.RoundAge = settings.RoundAge
	if err := data.Log.Level.UnmarshalText([]byte(settings.LogLevel)); err != nil {
		return data, fmt.Errorf("log level: %w", err)
	}

	for _, d := range []struct {
		name  string
		value string
		dst   *config.Duration
	}{
		{"Stale after", settings.StaleThreshold, &data.LastReading.StaleThreshold},
		{"Fetch delay", settings.FetchDelay, &data.Advanced.FetchDelay},
		{"Fallback interval", settings.FallbackInterval, &data.Advanced.FallbackInterval},
	} {
		if err := d.dst.UnmarshalText([]byte(strings.TrimSpace(d.value))); err != nil || d.dst.Duration < 0 {
			return data, fmt.Errorf("%s %w", d.name, ErrDuration)
		}
	}
	return data, nil
}

// Test fetches the latest reading, which is then used for previews.
func (s *Service) Test(ctx context.Context, settings Settings) (string, error) {
	data, err := s.toData(settings)
	if err != nil {
		return "", err
	}
	if data.URL == "" {
		return "", fetch.ErrNoURL
	}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	conf := config.New(config.WithVersion(s.conf.Version), config.WithData(data))
	props, err := fetch.NewFetch(conf, s.store).Do(ctx)
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	s.live = props
	s.mu.Unlock()
	return props.String(data), nil
}

func (s *Service) Preview(settings Settings) Preview {
	data, err := s.toData(settings)
	if err != nil {
		return Preview{Error: err.Error()}
	}

	s.mu.Lock()
	props, live := s.live, s.live != nil
	s.mu.Unlock()
	if props == nil {
		props = sampleProperties()
	}

	p := Preview{
		Text:  props.String(data),
		Stale: time.Since(props.Bgnow.Mills.Time) > data.LastReading.StaleThreshold.Duration,
		Live:  live,
	}
	if data.DynamicIcon.Enabled {
		img, err := dynamicicon.New(config.New(config.WithData(data))).Render(props)
		if err != nil {
			p.Error = "Dynamic icon: " + err.Error()
			return p
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			p.Error = err.Error()
			return p
		}
		p.Icon = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
		p.IconColor = hexColor(colorx.Hex{Color: dynamicicon.Color(data)})
	}
	return p
}

type Font struct {
	// File is the value saved in the config.
	File string `json:"file"`
	Name string `json:"name"`
}

func (s *Service) Fonts() []Font {
	var buf sfnt.Buffer
	seen := make(map[string]bool)
	var fonts []Font
	for _, path := range findfont.List() {
		file := filepath.Base(path)
		// The icon renderer can't load font collections (.ttc)
		ext := strings.ToLower(filepath.Ext(path))
		if (ext != ".ttf" && ext != ".otf") || seen[file] {
			continue
		}
		seen[file] = true
		name := fontName(&buf, path)
		// Names starting with a dot are private system fonts, hidden from font menus on macOS
		if strings.HasPrefix(name, ".") {
			continue
		}
		fonts = append(fonts, Font{File: file, Name: name})
	}
	slices.SortFunc(fonts, func(a, b Font) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})
	return fonts
}

// fontName reads a font's full name, falling back to its file name.
func fontName(buf *sfnt.Buffer, path string) string {
	fallback := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	f, err := os.Open(path)
	if err != nil {
		return fallback
	}
	defer func() {
		_ = f.Close()
	}()
	font, err := sfnt.ParseReaderAt(f)
	if err != nil {
		return fallback
	}
	name, err := font.Name(buf, sfnt.NameIDFull)
	if err != nil || name == "" {
		return fallback
	}
	return name
}

// ServeFont serves a font for previews, looked up like the dynamic icon's font.
func (s *Service) ServeFont(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")
	switch strings.ToLower(filepath.Ext(name)) {
	case "", ".ttf", ".otf":
	default:
		http.Error(w, "not a font", http.StatusBadRequest)
		return
	}
	b, err := dynamicicon.LoadFont(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", http.DetectContentType(b))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write(b) //nolint:gosec // Font data, served with its detected type
}

func (s *Service) BrowseFont() (string, error) {
	return application.Get().Dialog.OpenFile().
		SetTitle("Dynamic Icon Font").
		CanChooseFiles(true).
		AddFilter("Fonts", "*.ttf;*.otf").
		PromptForSingleSelection()
}

func (s *Service) NocturneConnected(url string) bool {
	return s.store.Connected(strings.TrimSpace(url))
}

func (s *Service) NocturneSupported(ctx context.Context, url string) bool {
	url = strings.TrimSpace(url)
	if url == "" {
		return false
	}
	if s.store.Connected(url) {
		return true
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.store.Supported(ctx, url)
}

// ConnectNocturne blocks until the user approves or denies the device code,
// which is sent with [EventNocturneCode].
func (s *Service) ConnectNocturne(ctx context.Context, url string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return fetch.ErrNoURL
	}
	err := s.store.Connect(ctx, url, func(auth *oauth2.DeviceAuthResponse) {
		verifyURL := auth.VerificationURIComplete
		if verifyURL == "" {
			verifyURL = auth.VerificationURI
		}
		slog.Info("Waiting for Nocturne approval", "url", verifyURL, "code", auth.UserCode)
		if app := application.Get(); app != nil {
			app.Event.Emit(EventNocturneCode, NocturneCode{Code: auth.UserCode, URL: verifyURL})
		}
		if err := s.openURL(verifyURL); err != nil {
			slog.Warn("Failed to open browser", "error", err)
		}
	})
	if err != nil {
		return err
	}

	// Save only the URL so the tray uses the new link
	data := s.conf.Data()
	data.URL = url
	return s.conf.Write(data)
}

func (s *Service) DisconnectNocturne(ctx context.Context) error {
	return s.store.Disconnect(ctx)
}

func (s *Service) OpenURL(url string) error {
	return s.openURL(url)
}

func sampleProperties() *nightscout.Properties {
	mills := nightscout.Mills{Time: time.Now().Add(-3 * time.Minute)}
	return &nightscout.Properties{
		Bgnow: nightscout.Reading{
			Last:  123,
			Mills: mills,
			Sgvs:  []nightscout.SGV{{Direction: "Flat", Mgdl: 123, Mills: mills}},
		},
		Delta: nightscout.Delta{Mgdl: 2},
	}
}

func formatDuration(d config.Duration) string {
	b, _ := d.MarshalText()
	return string(b)
}

// hexColor formats a color for <input type="color">.
func hexColor(c colorx.Hex) string {
	if c.Color == nil {
		return "#ffffff"
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}
