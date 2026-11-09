# Nightscout Menu Bar

[![Build](https://github.com/gabe565/nightscout-menu-bar/actions/workflows/build.yml/badge.svg)](https://github.com/gabe565/nightscout-menu-bar/actions/workflows/build.yml)

A small application that displays live blood sugar data from Nightscout on your menu bar.

Works on Windows, MacOS, and Linux.

<picture>
  <source media="(prefers-color-scheme: light)" srcset="https://github.com/gabe565/nightscout-menu-bar/assets/7717888/2e9e673a-d69f-43b8-9e82-d7168ffa5766">
  <img width="250" alt="Nightscout Menu Bar Screenshot" src="https://github.com/gabe565/nightscout-menu-bar/assets/7717888/c8c61eca-76ad-4078-ada4-95bad9cb01b9">
</picture>

## Install

### Homebrew (macOS, Linux)

<details>
  <summary>Click to expand</summary>

Install Nightscout Menu Bar from [gabe565/homebrew-tap](https://github.com/gabe565/homebrew-tap):
```shell
brew install gabe565/tap/nightscout-menu-bar
```
</details>

### APT (Ubuntu, Debian)

<details>
  <summary>Click to expand</summary>

1. If you don't have it already, install the `ca-certificates` package
   ```shell
   sudo apt install ca-certificates
   ```

2. Add gabe565 apt repository
   ```
   echo 'deb [trusted=yes] https://apt.gabe565.com /' | sudo tee /etc/apt/sources.list.d/gabe565.list
   ```

3. Update apt repositories
   ```shell
   sudo apt update
   ```

4. Install nightscout-menu-bar
   ```shell
   sudo apt install nightscout-menu-bar
   ```
</details>

### RPM (Fedora, RHEL)

<details>
  <summary>Click to expand</summary>

1. If you don't have it already, install the `ca-certificates` package
   ```shell
   sudo dnf install ca-certificates
   ```

2. Add gabe565 rpm repository to `/etc/yum.repos.d/gabe565.repo`
   ```ini
   [gabe565]
   name=gabe565
   baseurl=https://rpm.gabe565.com
   enabled=1
   gpgcheck=0
   ```

3. Install nightscout-menu-bar
   ```shell
   sudo dnf install nightscout-menu-bar
   ```
</details>

### AUR (Arch Linux)

<details>
  <summary>Click to expand</summary>

Install [nightscout-menu-bar-bin](https://aur.archlinux.org/packages/nightscout-menu-bar-bin) with your [AUR helper](https://wiki.archlinux.org/index.php/AUR_helpers) of choice.
</details>

### Binary

Automated builds are uploaded during the release process. See the [latest release](https://github.com/gabe565/nightscout-menu-bar/releases/latest) for download links.

### Linux requirements

The packages above install these automatically. Linux builds need GTK 4.14+ and WebKitGTK 6.0, which are available on Ubuntu 24.04, Debian 13, Fedora 40, and newer:

- **Debian/Ubuntu:** `sudo apt install libgtk-4-1 libwebkitgtk-6.0-4`
- **Fedora:** `sudo dnf install gtk4 webkitgtk6.0`
- **Arch:** `sudo pacman -S gtk4 webkitgtk-6.0`

GNOME needs the [AppIndicator extension](https://extensions.gnome.org/extension/615/appindicator-support/) to show tray icons. If it's missing, the app offers to install it.

## Usage

After launching Nightscout Menu Bar, open its tray menu and choose "Settings…" to configure the integration.

The settings window contains the following options:
- **Connection:** Nightscout URL (required), API token, and Nocturne connected apps. Connecting to Nocturne shows a code and opens your browser to approve it; no token needs to be copied. When connected, this takes precedence over the API token.
- **Display:** title, units (mg/dL or mmol/L), which parts of the reading appear in the menu bar, and the dynamic icon.
- **General:** start on login, and exposing readings over a local socket (see [`contrib/powerlevel10k`](contrib/powerlevel10k)).
- **Advanced:** fetch timing, age rounding, and log level.

Tokens are stored in the system keyring: Keychain on macOS, Credential Manager on Windows, or the Secret Service on Linux.

## Contrib

Integrations with external tools are available in the [contrib](contrib) directory.

## Development

The systray menu is provided by
[fyne.io/systray](https://github.com/fyne-io/systray). See
[systray's platform notes](https://github.com/getlantern/systray#platform-notes)
for required dependencies.

The settings window is built with [Wails](https://wails.io), which uses the
system webview. On Linux, building it requires GTK 4 and WebKitGTK 6.0
development packages (e.g. `libgtk-4-dev libwebkitgtk-6.0-dev` on Debian/Ubuntu).

#### macOS

To generate a Mac app, run [hack/build-darwin.sh](hack/build-darwin.sh).
An app for each architecture will be created in the `dist` directory.
