package tray

import (
	"context"
	"image/color"
	"io"
	"log/slog"
	"os"

	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/assets"
	"gabe565.com/nightscout-menu-bar/internal/autostart"
	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/dynamicicon"
	"gabe565.com/nightscout-menu-bar/internal/fetch"
	"gabe565.com/nightscout-menu-bar/internal/ticker"
	"gabe565.com/nightscout-menu-bar/internal/tray/items"
	"gabe565.com/nightscout-menu-bar/internal/tray/messages"
	"github.com/skratchdot/open-golang/open"
)

const AboutURL = "https://github.com/gabe565/nightscout-menu-bar"

func New(version string) *Tray {
	t := &Tray{
		config: config.New(config.WithVersion(version)),
		bus:    make(chan any, 1),
	}
	if err := t.config.Flags.Parse(os.Args[1:]); err != nil {
		_, _ = io.WriteString(os.Stderr, err.Error()+"\n")
		os.Exit(2)
	}

	if err := t.config.Load(); err != nil {
		t.onError(err)
	}

	t.ticker = ticker.New(t.config, t.bus)

	if t.config.DynamicIcon.Enabled {
		t.dynamicIcon = dynamicicon.New(t.config)
	}

	t.config.AddCallback(func() {
		t.bus <- messages.ReloadConfigMsg{}
	})
	return t
}

type Tray struct {
	config      *config.Config
	ticker      *ticker.Ticker
	dynamicIcon *dynamicicon.DynamicIcon
	bus         chan any
	items       items.Items
}

func (t *Tray) Run(ctx context.Context) {
	t.ticker.Start(ctx)
	if err := t.config.Watch(ctx); err != nil {
		t.onError(err)
	}
	systray.Run(t.onReady(ctx), t.onExit)
}

func (t *Tray) Quit() {
	systray.Quit()
}

func (t *Tray) onReady(ctx context.Context) func() { //nolint:gocyclo
	return func() {
		systray.SetTemplateIcon(assets.Nightscout, assets.Nightscout)
		if t.dynamicIcon == nil {
			systray.SetTitle(t.config.Title)
		}
		systray.SetTooltip(t.config.Title)

		t.items = items.New(t.config)

		for {
			select {
			case <-ctx.Done():
				t.Quit()
			case <-t.items.OpenNightscout.ClickedCh:
				u, err := fetch.BuildURLWithToken(t.config)
				if err != nil {
					t.onError(err)
					return
				}
				slog.Debug("Opening Nightscout", "url", u)
				if err := open.Run(u.String()); err != nil {
					t.onError(err)
				}
			case <-t.items.Preferences.URL.ClickedCh:
				go func() {
					if err := t.items.Preferences.URL.Prompt(); err != nil {
						t.onError(err)
					}
				}()
			case <-t.items.About.ClickedCh:
				if err := open.Run(AboutURL); err != nil {
					t.onError(err)
				}
			case <-t.items.Preferences.Token.ClickedCh:
				go func() {
					if err := t.items.Preferences.Token.Prompt(); err != nil {
						t.onError(err)
					}
				}()
			case <-t.items.Preferences.Units.ClickedCh:
				go func() {
					if err := t.items.Preferences.Units.Prompt(); err != nil {
						t.onError(err)
					}
				}()
			case <-t.items.Preferences.StartOnLogin.ClickedCh:
				if t.items.Preferences.StartOnLogin.Checked() {
					if err := autostart.Disable(); err != nil {
						t.onError(err)
						continue
					}
					t.items.Preferences.StartOnLogin.Uncheck()
				} else {
					if err := autostart.Enable(); err != nil {
						t.onError(err)
						continue
					}
					t.items.Preferences.StartOnLogin.Check()
				}
			case <-t.items.Preferences.Socket.ClickedCh:
				if err := t.items.Preferences.Socket.Toggle(); err != nil {
					t.onError(err)
				}
			case <-t.items.Preferences.DynamicIcon.ClickedCh:
				if err := t.items.Preferences.DynamicIcon.Toggle(); err != nil {
					t.onError(err)
				}
			case <-t.items.Preferences.DynamicIconColor.ClickedCh:
				if err := t.items.Preferences.DynamicIconColor.Choose(); err != nil {
					t.onError(err)
				}
			case <-t.items.Quit.ClickedCh:
				t.Quit()
			case msg := <-t.bus:
				switch msg := msg.(type) {
				case messages.RenderMessage:
					if msg.Type == messages.RenderTypeFetch {
						t.items.Error.Hide()
					}

					value := msg.Properties.String(t.config)
					slog.Debug("Updating reading", "value", value)
					if t.dynamicIcon == nil {
						systray.SetTitle(value)
					} else {
						if icon, err := t.dynamicIcon.Generate(msg.Properties); err == nil {
							systray.SetTitle("")
							if t.config.DynamicIcon.FontColor.Color == color.White {
								systray.SetTemplateIcon(icon, icon)
							} else {
								systray.SetIcon(icon)
							}
						} else {
							t.onError(err)
							systray.SetTitle(value)
							systray.SetTemplateIcon(assets.Nightscout, assets.Nightscout)
						}
					}
					systray.SetTooltip(value)
					t.items.LastReading.SetTitle(value)

					for i, reading := range msg.Properties.Buckets {
						if i < len(t.items.History.Subitems) {
							t.items.History.Subitems[i].SetTitle(reading.String(t.config))
						} else {
							entry := t.items.History.AddSubMenuItem(reading.String(t.config), "")
							entry.Disable()
							t.items.History.Subitems = append(t.items.History.Subitems, entry)
						}
					}
				case error:
					slog.Error("Displaying error", "error", msg)
					t.items.Error.SetTitle(msg.Error())
					t.items.Error.Show()
				case messages.ReloadConfigMsg:
					if t.config.DynamicIcon.Enabled {
						t.dynamicIcon = dynamicicon.New(t.config)
					} else if t.dynamicIcon != nil {
						t.dynamicIcon = nil
						systray.SetTemplateIcon(assets.Nightscout, assets.Nightscout)
					}
				}
			}
		}
	}
}

func (t *Tray) onError(err error) {
	select {
	case t.bus <- err:
	default:
		slog.Error("Unable to display error due to full bus", "error", err)
	}
}

func (t *Tray) onExit() {
	slog.Info("Exiting")
	t.ticker.Close()
	close(t.bus)
}
