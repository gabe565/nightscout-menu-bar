# Nightscout Menu Bar

[![Build](https://github.com/gabe565/nightscout-menu-bar/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/nightscout-menu-bar/actions/workflows/build.yml)

A small application that displays live blood sugar data from Nightscout on your menu bar.
Should be supported by macOS, Linux, and Windows, but only macOS has been tested so far.

![macOS Screenshot](../assets/macos-screenshot.png?raw=true)

## Build

The systray menu is provided by
[getlantern/systray](https://github.com/getlantern/systray). See
[systray's platform notes](https://github.com/getlantern/systray#platform-notes)
for required dependencies.

### macOS

To generate a Mac app, run [hack/build-app.sh](hack/build-app.sh).
An app will be created in the `dist` directory.
