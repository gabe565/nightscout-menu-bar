package preferences

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync/atomic"

	"fyne.io/systray"
	"gabe565.com/nightscout-menu-bar/internal/config"
	"gabe565.com/nightscout-menu-bar/internal/fetch"
	"gabe565.com/nightscout-menu-bar/internal/oauth"
	"github.com/ncruces/zenity"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
)

func NewNocturne(conf *config.Config, store *oauth.Store, parent *systray.MenuItem) *Nocturne {
	n := &Nocturne{
		config:   conf,
		store:    store,
		MenuItem: parent.AddSubMenuItem("", ""),
	}
	n.UpdateTitle()
	return n
}

type Nocturne struct {
	config *config.Config
	store  *oauth.Store
	busy   atomic.Bool
	*systray.MenuItem
}

func (n *Nocturne) UpdateTitle() {
	if n.store.Connected(n.config.Data().URL) {
		n.SetTitle("Disconnect from Nocturne")
	} else {
		n.SetTitle("Connect to Nocturne…")
	}
}

func (n *Nocturne) Toggle(ctx context.Context) error {
	if !n.busy.CompareAndSwap(false, true) {
		return nil
	}
	defer n.busy.Store(false)
	defer n.UpdateTitle()

	if n.store.Connected(n.config.Data().URL) {
		return n.disconnect(ctx)
	}
	return n.connect(ctx)
}

func (n *Nocturne) disconnect(ctx context.Context) error {
	if err := zenity.Question(
		"Disconnect from Nocturne? Readings will stop until you reconnect or configure an API token.",
		zenity.Title("Disconnect from Nocturne"),
		zenity.OKLabel("Disconnect"),
	); err != nil {
		if errors.Is(err, zenity.ErrCanceled) {
			return nil
		}
		return err
	}
	return n.store.Disconnect(ctx)
}

func (n *Nocturne) connect(ctx context.Context) error {
	serverURL := n.config.Data().URL
	if serverURL == "" {
		return fetch.ErrNoURL
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	codeCh := make(chan *oauth2.DeviceAuthResponse, 1)
	done := make(chan error, 1)
	go func() {
		done <- n.store.Connect(ctx, serverURL, func(auth *oauth2.DeviceAuthResponse) {
			codeCh <- auth
		})
	}()

	var auth *oauth2.DeviceAuthResponse
	select {
	case auth = <-codeCh:
	case err := <-done:
		return err
	}

	verifyURL := auth.VerificationURIComplete
	if verifyURL == "" {
		verifyURL = auth.VerificationURI
	}
	slog.Info("Waiting for Nocturne approval", "url", verifyURL, "code", auth.UserCode)
	if err := open.Run(verifyURL); err != nil {
		slog.Warn("Failed to open browser", "error", err)
	}

	// The dialog closes itself once the flow finishes
	dialogCtx, closeDialog := context.WithCancel(ctx)
	defer closeDialog()
	finished := make(chan struct{})
	var connectErr error
	go func() {
		connectErr = <-done
		close(finished)
		closeDialog()
	}()

	dialogErr := zenity.Info(
		fmt.Sprintf(
			"Approve Nightscout Menu Bar in your browser and confirm this code:\n\n%s\n\nIf your browser didn't open, visit:\n%s",
			auth.UserCode,
			auth.VerificationURI,
		),
		zenity.Title("Connect to Nocturne"),
		zenity.OKLabel("Cancel"),
		zenity.NoIcon,
		zenity.Context(dialogCtx),
	)

	select {
	case <-finished:
		return connectErr
	default:
		// The user dismissed the dialog before approving
		cancel()
		<-finished
		if dialogErr != nil && !errors.Is(dialogErr, zenity.ErrCanceled) {
			return dialogErr
		}
		return nil
	}
}
