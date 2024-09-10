package tray

import (
	"context"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/systray"
	"github.com/gabe565/nightscout-menu-bar/internal/app/settings"
	"github.com/gabe565/nightscout-menu-bar/internal/assets"
	"github.com/gabe565/nightscout-menu-bar/internal/dynamicicon"
	"github.com/gabe565/nightscout-menu-bar/internal/fyneutil"
	"github.com/gabe565/nightscout-menu-bar/internal/ticker"
	"github.com/gabe565/nightscout-menu-bar/internal/tray/items"
	"github.com/gabe565/nightscout-menu-bar/internal/tray/messages"
	"github.com/gabe565/nightscout-menu-bar/internal/util"
)

func New(app fyneutil.DesktopApp, version string) *Tray {
	t := &Tray{
		app:   app,
		bus:   make(chan any, 1),
		items: items.New(app),
	}

	t.ticker = ticker.New(app, t.bus, version)

	if app.Preferences().Bool(settings.DynamicIconEnabledKey) {
		t.dynamicIcon = dynamicicon.New(app)
	}

	app.Preferences().AddChangeListener(t.Reload)
	return t
}

type Tray struct {
	app         fyneutil.DesktopApp
	ticker      *ticker.Ticker
	dynamicIcon *dynamicicon.DynamicIcon
	bus         chan any
	menu        *fyne.Menu
	items       items.Items
}

func (t *Tray) Menu() *fyne.Menu {
	if t.menu == nil {
		t.menu = fyne.NewMenu("Nightscout",
			t.items.LastReading,
			fyne.NewMenuItemSeparator(),
			t.items.OpenNightscout,
			t.items.History,
			fyne.NewMenuItemSeparator(),
			t.items.Settings,
			t.items.About,
			t.items.Quit,
		)
	}
	return t.menu
}

func (t *Tray) Run(ctx context.Context) {
	t.ticker.Start(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				t.Close()
				return
			case msg := <-t.bus:
				switch msg := msg.(type) {
				case messages.RenderMessage:
					prefs := t.app.Preferences()

					value := msg.Properties.String(t.app.Preferences())
					slog.Debug("Updating reading", "value", value)
					if t.dynamicIcon == nil {
						systray.SetTitle(value)
					} else {
						if icon, err := t.dynamicIcon.Generate(msg.Properties); err == nil {
							systray.SetTitle("")
							var c util.HexColor
							if err := c.UnmarshalText([]byte(prefs.String(settings.DynamicIconFontColorKey))); err != nil {
								c = util.White()
							}
							switch c {
							case util.White():
								systray.SetTemplateIcon(icon, icon)
							default:
								systray.SetIcon(icon)
							}
						} else {
							t.onError(err)
							systray.SetTitle(value)
							t.app.SetSystemTrayIcon(assets.NightscoutResource)
						}
					}
					t.items.LastReading.Label = value
					t.menu.Refresh()
					systray.SetTooltip(value)

					for i, reading := range msg.Properties.Buckets {
						if i < len(t.items.History.ChildMenu.Items) {
							t.items.History.ChildMenu.Items[i].Label = reading.String(prefs)
						} else {
							entry := fyne.NewMenuItem(reading.String(prefs), nil)
							entry.Disabled = true
							t.items.History.ChildMenu.Items = append(t.items.History.ChildMenu.Items, entry)
						}
					}
					t.items.History.ChildMenu.Refresh()
				case error:
					slog.Error("Displaying error", "error", msg)
					t.items.LastReading.Label = msg.Error()
					t.menu.Refresh()
				}
			}
		}
	}()
}

func (t *Tray) Reload() {
	prefs := t.app.Preferences()
	if prefs.Bool(settings.DynamicIconEnabledKey) {
		t.dynamicIcon = dynamicicon.New(t.app)
	} else if t.dynamicIcon != nil {
		t.dynamicIcon = nil
		t.app.SetSystemTrayIcon(assets.NightscoutResource)
	}
}

func (t *Tray) onError(err error) {
	select {
	case t.bus <- err:
	default:
		slog.Error("Unable to display error due to full bus", "error", err)
	}
}

func (t *Tray) Close() {
	slog.Info("Exiting")
	t.ticker.Close()
	close(t.bus)
}
