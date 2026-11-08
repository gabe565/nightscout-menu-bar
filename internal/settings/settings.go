package settings

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"

	"gabe565.com/nightscout-menu-bar/internal/assets"
	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/oauth"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

const AppID = "com.gabe565.nightscout-menu-bar.settings"

//go:embed frontend
var frontend embed.FS

func assetHandler(fsys fs.FS, svc *Service) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", application.AssetFileServerFS(fsys))
	mux.HandleFunc("GET /font", svc.ServeFont)
	return mux
}

func Run(version string, args []string) error {
	conf := config.New(config.WithVersion(version))
	if err := conf.Flags.Parse(args); err != nil {
		return err
	}
	if err := conf.Load(); err != nil {
		return err
	}

	store := oauth.NewStore(conf)
	if err := store.Load(); err != nil {
		slog.Error("Failed to load Nocturne credentials", "error", err)
	}

	assetFS, err := fs.Sub(frontend, "frontend")
	if err != nil {
		return err
	}

	svc := NewService(conf, store)
	var window *application.WebviewWindow
	app := application.New(application.Options{
		Name:     "Nightscout Menu Bar Settings",
		Icon:     assets.App,
		Services: []application.Service{application.NewService(svc)},
		ShouldQuit: func() bool {
			if svc.canClose() {
				return true
			}
			application.Get().Event.Emit(EventCloseRequested)
			return false
		},
		Assets: application.AssetOptions{Handler: assetHandler(assetFS, svc)},
		Mac: application.MacOptions{
			// The app bundle is a menu bar agent (LSUIElement); the settings
			// window should still get a Dock icon and come to the front.
			ActivationPolicy: application.ActivationPolicyRegular,
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: AppID,
			OnSecondInstanceLaunch: func(application.SecondInstanceData) {
				if window != nil {
					window.Show()
					window.Restore()
					window.Focus()
				}
			},
		},
	})

	window = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Nightscout Menu Bar Settings",
		Width:            560,
		Height:           520,
		MinWidth:         480,
		MinHeight:        320,
		URL:              "/",
		BackgroundColour: application.NewRGBA(0, 0, 0, 0),
		Mac: application.MacWindow{
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInsetUnified,
			InvisibleTitleBarHeight: 52,
		},
	})
	svc.window = window

	// Let the page ask about unsaved changes
	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		if !svc.canClose() {
			e.Cancel()
			app.Event.Emit(EventCloseRequested)
		}
	})

	return app.Run()
}
