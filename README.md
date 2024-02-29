# Nightscout Menu Bar

[![Build](https://github.com/gabe565/nightscout-menu-bar/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/nightscout-menu-bar/actions/workflows/build.yml)

A small application that displays live blood sugar data from Nightscout on your menu bar.

Works on Windows, MacOS, and Linux.

![macOS Screenshot](assets/macos-screenshot.webp?raw=true)

## Install

### Binary

Eventually, binaries will be attached to releases.
For now, binaries can be downloaded from CI build artifacts.
1. Go to [the build action](https://github.com/gabe565/nightscout-menu-bar/actions/workflows/build.yml?query=branch%3Amain+is%3Asuccess).
2. Click the latest build job.
3. Scroll down to "Artifacts".
4. Download the artifact for your operating system/architecture.
5. A zip file will be downloaded containing Nightscout Menu Bar!

### Build

The systray menu is provided by
[fyne.io/systray](https://github.com/fyne-io/systray). See
[systray's platform notes](https://github.com/getlantern/systray#platform-notes)
for required dependencies.

#### macOS

To generate a Mac app, run [hack/build-darwin.sh](hack/build-darwin.sh).
An app will be created in the `dist` directory.

## Contrib

Integrations with external tools are available in the [contrib](contrib) directory.
