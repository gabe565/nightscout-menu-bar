# Nightscout Menu Bar

[![Build](https://github.com/gabe565/nightscout-menu-bar/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/nightscout-menu-bar/actions/workflows/build.yml)

A small application that displays live blood sugar data from Nightscout on your menu bar.
Should be supported by macOS, Linux, and Windows, but only macOS has been tested so far.

![macOS Screenshot](assets/macos-screenshot.webp?raw=true)

## Build

The systray menu is provided by
[fyne.io/systray](https://github.com/fyne-io/systray). See
[systray's platform notes](https://github.com/getlantern/systray#platform-notes)
for required dependencies.

### macOS

To generate a Mac app, run [hack/build-darwin.sh](hack/build-darwin.sh).
An app will be created in the `dist` directory.

## Contrib

Integrations with external tools are available in the [contrib](contrib) directory.
