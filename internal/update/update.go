package update

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

const releaseURL = "https://api.github.com/repos/gabe565/nightscout-menu-bar/releases/latest"

var ErrUnexpectedStatus = errors.New("unexpected status")

type Release struct {
	Version string `json:"tag_name"`
	URL     string `json:"html_url"`
}

type Checker struct {
	client  *http.Client
	url     string
	version string
}

func NewChecker(version string) *Checker {
	return &Checker{
		client:  &http.Client{Timeout: time.Minute},
		url:     releaseURL,
		version: version,
	}
}

// Enabled reports whether the running version can be compared to releases.
func (c *Checker) Enabled() bool {
	return semver.IsValid(c.version)
}

// Check returns the latest release if it's newer than the running version, or nil.
func (c *Checker) Check(ctx context.Context) (*Release, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "nightscout-menu-bar/"+c.version)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("checking for updates: %w: %s", ErrUnexpectedStatus, resp.Status)
	}

	var rel Release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}
	if !semver.IsValid(rel.Version) || semver.Compare(rel.Version, c.version) <= 0 {
		return nil, nil //nolint:nilnil
	}
	return &rel, nil
}

type Method uint8

const (
	MethodManual Method = iota
	MethodHomebrew
	MethodPackage
)

// InstallMethod guesses how the running binary was installed.
func InstallMethod() Method {
	exe, err := os.Executable()
	if err != nil {
		return MethodManual
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	return installMethod(runtime.GOOS, exe, func(path string) bool {
		_, err := os.Stat(path)
		return err == nil
	})
}

func installMethod(goos, exe string, exists func(string) bool) Method {
	switch {
	case strings.Contains(exe, "/Caskroom/"):
		return MethodHomebrew
	case goos == "darwin":
		for _, prefix := range []string{"/opt/homebrew", "/usr/local"} {
			if exists(prefix + "/Caskroom/nightscout-menu-bar") {
				return MethodHomebrew
			}
		}
	case goos == "linux" && exe == "/usr/bin/nightscout-menu-bar":
		return MethodPackage
	}
	return MethodManual
}
