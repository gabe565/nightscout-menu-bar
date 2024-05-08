# Nightscout Menu Bar

[![Build](https://github.com/gabe565/nightscout-menu-bar/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/nightscout-menu-bar/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gabe565/nightscout-menu-bar)](https://goreportcard.com/report/github.com/gabe565/nightscout-menu-bar)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=gabe565_nightscout-menu-bar&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=gabe565_nightscout-menu-bar)

A small application that displays live blood sugar data from Nightscout on your menu bar.

Works on Windows, MacOS, and Linux.

<picture>
  <source media="(prefers-color-scheme: light)" srcset="https://github.com/gabe565/nightscout-menu-bar/assets/7717888/2e9e673a-d69f-43b8-9e82-d7168ffa5766">
  <img width="250" alt="Nightscout Menu Bar Screenshot" src="https://github.com/gabe565/nightscout-menu-bar/assets/7717888/c8c61eca-76ad-4078-ada4-95bad9cb01b9">
</picture>

## Install

### Brew (macOS)

```shell
brew install gabe565/tap/nightscout-menu-bar --no-quarantine
```

### Binary

Automated builds are uploaded during the release process. See the [latest release](https://github.com/gabe565/nightscout-menu-bar/releases/latest) for download links.

## Usage

After launching Nightscout Menu Bar, you will need to open its tray menu, then hover over "Preferences" to configure the integration.

The preferences menu contains the following options:
- Nightscout URL (required)
- API Token
- Units: mg/dL or mmol/L
- Start on login
- Write to a local file (see [`contrib/powerlevel10k`](contrib/powerlevel10k))

Additional configuration is available in a configuration file, which can be found in the following locations:
- **Windows:** `%AppData%\nightscout-menu-bar\config.toml`
- **macOS:** `~/Library/Application Support/nightscout-menu-bar/config.toml`
- **Linux:** `~/.config/nightscout-menu-bar/config.toml`

## Contrib

Integrations with external tools are available in the [contrib](contrib) directory.

## Development

The systray menu is provided by
[fyne.io/systray](https://github.com/fyne-io/systray). See
[systray's platform notes](https://github.com/getlantern/systray#platform-notes)
for required dependencies.

#### macOS

To generate a Mac app, run [hack/build-darwin.sh](hack/build-darwin.sh).
An app will be created in the `dist` directory.
