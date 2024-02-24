# Powerlevel10k Segment

Adds a custom segment to the [Powerlevel10k](https://github.com/romkatv/powerlevel10k) zsh theme.

![Demo](./screenshot.webp?raw=true)

## Install

### Prerequisites

#### Enable local file writes
1. Open the nightscout-menu-bar menu.
2. Hover over "Preferences".
3. Check "Write to local file".[^1]

### Install the Powerlevel10k segment

The Nightscout segment can either be downloaded to a separate file, and sourced from `~/.p10k.zsh`, or its contents can be pasted directly into `~/.p10k.zsh`.

#### Install directly into `~/.p10k.zsh`
This method adds the custom segment's code and configuration directly into `~/.p10k.zsh`. It is easier to set up, but may be harder to update in the future.

<details>
   <summary>Click to expand</summary>

1. Copy the contents of the [segment script](nightscout.zsh).
2. Edit `~/.p10k.zsh`.
3. Search for `p10k reload`.
4. Somewhere before this line, paste the segment file contents.
5. Search for `POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS`.
6. Add `nightscout` somewhere in this array, depending on where you would like the widget to be rendered.
7. Open a new shell, or restart your current shell with `exec zsh`.
8. Nightscout data should be rendered as a right-segment!
</details>

#### Install to a separate file
This method places the custom segment's code and configuration in a separate file. It is less standard, but makes it easier to update the segment in the future.

<details>
   <summary>Click to expand</summary>

1. Download the [segment script](nightscout.zsh) to a local file. This file will be sourced during Powerlevel10k initialization.
   1. For example: `~/.p10k/nightscout.zsh`.
2. Edit `~/.p10k.zsh`.
3. Search for `p10k reload`.
4. Somewhere before this line, source the segment file.
   1. For example: `source ~/.p10k/nightscout.zsh`
5. Search for `POWERLEVEL9K_RIGHT_PROMPT_ELEMENTS`.
6. Add `nightscout` somewhere in this array, depending on where you would like the widget to be rendered.
7. Open a new shell, or restart your current shell with `exec zsh`.
8. Nightscout data should be rendered as a right-segment!
</details>

[^1]: If this option is unavailable, update nightscout-menu-bar.
