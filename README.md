# Nightscout Menu Bar

[![Build](https://github.com/gabe565/nightscout-menu-bar/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/nightscout-menu-bar/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gabe565/nightscout-menu-bar)](https://goreportcard.com/report/github.com/gabe565/nightscout-menu-bar)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=gabe565_nightscout-menu-bar&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=gabe565_nightscout-menu-bar)

A small application that displays live blood sugar data from Nightscout on your menu bar.

Works on Windows, MacOS, and Linux.

![macOS Screenshot](assets/macos-screenshot.webp?raw=true)

## Install

### Brew (macOS)

```shell
brew install gabe565/tap/nightscout-menu-bar --no-quarantine
```

### Binary

Automated builds are uploaded during the release process. See the [latest release](https://github.com/gabe565/nightscout-menu-bar/releases/latest) for download links.

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
