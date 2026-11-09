package tray

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"time"

	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/assets"
	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/dynamicicon"
	"gabe565.com/nightscout-menu-bar/internal/fetch"
	"gabe565.com/nightscout-menu-bar/internal/oauth"
	"gabe565.com/nightscout-menu-bar/internal/socket"
	"gabe565.com/nightscout-menu-bar/internal/systheme"
	"gabe565.com/nightscout-menu-bar/internal/ticker"
	"gabe565.com/nightscout-menu-bar/internal/tray/items"
	"gabe565.com/nightscout-menu-bar/internal/tray/messages"
	"github.com/skratchdot/open-golang/open"
)

const AboutURL = "https://github.com/gabe565/nightscout-menu-bar"

func New(version string) *Tray {
	t := &Tray{
		config:       config.New(config.WithVersion(version)),
		bus:          make(chan any, 1),
		checkUpdate:  make(chan struct{}, 1),
		settingsCmds: make(map[*exec.Cmd]struct{}),
	}
	if err := t.config.Flags.Parse(os.Args[1:]); err != nil {
		_, _ = io.WriteString(os.Stderr, err.Error()+"\n")
		os.Exit(2)
	}

	if err := t.config.Load(); err != nil {
		t.displayError(err)
	}

	t.oauth = oauth.NewStore(t.config)
	if err := t.oauth.Load(); err != nil {
		t.displayError(err)
	}

	t.socket = socket.New(t.config, t.handleCommand)
	t.ticker = ticker.New(t.config, t.oauth, t.socket, t.bus)

	if t.config.Data().DynamicIcon.Enabled {
		t.dynamicIcon = dynamicicon.New(t.config)
	}

	t.config.AddCallback(func() {
		t.bus <- messages.ReloadConfigMsg{}
	})
	return t
}

type Tray struct {
	config      *config.Config
	oauth       *oauth.Store
	socket      *socket.Socket
	ticker      *ticker.Ticker
	dynamicIcon *dynamicicon.DynamicIcon
	bus         chan any
	items       items.Items

	checkUpdate chan struct{}

	settingsMu   sync.Mutex
	settingsCmds map[*exec.Cmd]struct{}
}

func (t *Tray) Run(ctx context.Context) {
	defer func() {
		slog.Info("Exiting")
		t.Close()
	}()
	if err := t.socket.Listen(); err != nil {
		if errors.Is(err, socket.ErrRunning) {
			slog.Info("Already running, opening settings in the existing instance")
			return
		}
		t.displayError(fmt.Errorf("failed to listen on socket: %w", err))
	}
	t.ticker.Start(ctx)
	t.checkTrayHost(ctx)
	go t.checkUpdates(ctx)
	if err := t.config.Watch(ctx); err != nil {
		t.displayError(err)
	}
	if err := t.config.WatchSecrets(ctx); err != nil {
		t.displayError(err)
	}
	if err := systheme.Watch(ctx, func() { t.bus <- messages.ThemeChangedMsg{} }); err != nil {
		t.displayError(err)
	}
	if err := t.oauth.Watch(ctx, t.ticker.Reload); err != nil {
		t.displayError(err)
	}
	systray.Run(t.onReady(ctx), nil)
}

func (t *Tray) onReady(ctx context.Context) func() {
	return func() {
		systray.SetTemplateIcon(assets.Nightscout, assets.Nightscout)
		if t.dynamicIcon == nil {
			systray.SetTitle(t.config.Data().Title)
		}
		systray.SetTooltip(t.config.Data().Title)

		t.items = items.New(t.config)
		t.items.UpdateIcons()
		go t.eventLoop(ctx)
	}
}

func (t *Tray) eventLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			systray.Quit()
			return
		case <-t.items.OpenNightscout.ClickedCh:
			u, err := fetch.BuildURLWithToken(t.config.Data())
			if err != nil {
				t.displayError(err)
				continue
			}
			slog.Debug("Opening Nightscout", "url", u)
			if err := open.Run(u.String()); err != nil {
				t.displayError(err)
			}
		case <-t.items.Settings.ClickedCh:
			if err := t.openSettings(ctx); err != nil {
				t.displayError(err)
			}
		case <-systray.ReopenedCh:
			if err := t.openSettings(ctx); err != nil {
				t.displayError(err)
			}
		case <-t.items.Update.ClickedCh:
			if err := open.Run(t.items.Update.URL); err != nil {
				t.displayError(err)
			}
		case <-t.items.About.ClickedCh:
			if err := open.Run(AboutURL); err != nil {
				t.displayError(err)
			}
		case <-t.items.Quit.ClickedCh:
			systray.Quit()
		case msg := <-t.bus:
			switch msg := msg.(type) {
			case messages.RenderMessage:
				if msg.Type == messages.RenderTypeFetch {
					t.items.Error.Hide()
				}

				conf := t.config.Data()
				shortValue := msg.Properties.String(conf)
				conf.LastReading.HideArrow = false
				conf.LastReading.HideDelta = false
				conf.LastReading.HideTimeAgo = false
				fullValue := msg.Properties.String(conf)
				isOldReading := time.Since(
					msg.Properties.Bgnow.Mills.Time,
				) > conf.LastReading.StaleThreshold.Duration

				slog.Debug("Updating reading", "value", fullValue)
				if t.dynamicIcon == nil {
					systray.SetTitleStrikethrough(shortValue, isOldReading)
				} else {
					if icon, err := t.dynamicIcon.Generate(msg.Properties); err == nil {
						systray.SetTitle("")
						if dynamicicon.Color(t.config.Data()) == color.White {
							systray.SetTemplateIcon(icon, icon)
						} else {
							systray.SetIcon(icon)
						}
					} else {
						t.displayError(err)
						systray.SetTitleStrikethrough(shortValue, isOldReading)
						systray.SetTemplateIcon(assets.Nightscout, assets.Nightscout)
					}
				}
				systray.SetTooltip(fullValue)
				t.items.LastReading.SetTitle(fullValue)

				for i, reading := range msg.Properties.Buckets {
					if i < len(t.items.History.Subitems) {
						t.items.History.Subitems[i].SetTitle(reading.String(t.config.Data()))
					} else {
						entry := t.items.History.AddSubMenuItem(reading.String(t.config.Data()), "")
						entry.Disable()
						t.items.History.Subitems = append(t.items.History.Subitems, entry)
					}
				}
			case error:
				slog.Error("Displaying error", "error", msg)
				t.items.Error.SetTitle(msg.Error())
				t.items.Error.Show()
			case messages.UpdateAvailableMsg:
				t.items.Update.Show(msg.Release)
			case messages.OpenSettingsMsg:
				if err := t.openSettings(ctx); err != nil {
					t.displayError(err)
				}
			case messages.ThemeChangedMsg:
				t.items.UpdateIcons()
				t.ticker.Rerender()
			case messages.ReloadConfigMsg:
				if t.config.Data().Updates.Check {
					select {
					case t.checkUpdate <- struct{}{}:
					default:
					}
				} else {
					t.items.Update.Hide()
				}
				if t.config.Data().DynamicIcon.Enabled {
					t.dynamicIcon = dynamicicon.New(t.config)
				} else if t.dynamicIcon != nil {
					t.dynamicIcon = nil
					systray.SetTemplateIcon(assets.Nightscout, assets.Nightscout)
				}
			}
		}
	}
}

func (t *Tray) displayError(err error) {
	select {
	case t.bus <- err:
	default:
		slog.Error("Unable to display error due to full bus", "error", err)
	}
}

// handleCommand runs commands sent over the socket by another instance.
func (t *Tray) handleCommand(cmd string) {
	switch cmd {
	case socket.CommandSettings:
		select {
		case t.bus <- messages.OpenSettingsMsg{}:
		default:
			slog.Error("Unable to open settings due to full bus")
		}
	default:
		slog.Warn("Unknown socket command", "command", cmd)
	}
}

func (t *Tray) Close() {
	t.closeSettings()
	if t.ticker != nil {
		t.ticker.Close()
	}
	if t.socket != nil {
		if err := t.socket.Close(); err != nil {
			slog.Error("Failed to close socket", "error", err)
		}
	}
	if t.bus != nil {
		close(t.bus)
	}
}
