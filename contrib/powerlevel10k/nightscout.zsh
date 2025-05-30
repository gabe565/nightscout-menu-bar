#!/usr/bin/env zsh

#################################[ nightscout: blood sugar ]#################################
# Nightscout state file. Typically does not need to be changed.
typeset -g NIGHTSCOUT_SOCKET="$TMPDIR/nightscout.sock"

# Nightscout styling will be chosen if the reading is below a given value.
typeset -g NIGHTSCOUT_THRESHOLD_OLD_MINS=10
typeset -g NIGHTSCOUT_THRESHOLD_URGENT_LOW=55
typeset -g NIGHTSCOUT_THRESHOLD_LOW=80
typeset -g NIGHTSCOUT_THRESHOLD_IN_RANGE=160
typeset -g NIGHTSCOUT_THRESHOLD_HIGH=260

# Show/hide Nightscout parts.
typeset -g NIGHTSCOUT_SHOW_ARROW=true
typeset -g NIGHTSCOUT_SHOW_DELTA=true
typeset -g NIGHTSCOUT_SHOW_TIMESTAMP=true

# Nightscout colors.
# Urgent low styling.
typeset -g POWERLEVEL9K_NIGHTSCOUT_URGENT_LOW_BACKGROUND=1
typeset -g POWERLEVEL9K_NIGHTSCOUT_URGENT_LOW_FOREGROUND=7
# Low styling.
typeset -g POWERLEVEL9K_NIGHTSCOUT_LOW_BACKGROUND=1
typeset -g POWERLEVEL9K_NIGHTSCOUT_LOW_FOREGROUND=7
# In range styling.
typeset -g POWERLEVEL9K_NIGHTSCOUT_IN_RANGE_BACKGROUND=2
typeset -g POWERLEVEL9K_NIGHTSCOUT_IN_RANGE_FOREGROUND=0
# High styling.
typeset -g POWERLEVEL9K_NIGHTSCOUT_HIGH_BACKGROUND=3
typeset -g POWERLEVEL9K_NIGHTSCOUT_HIGH_FOREGROUND=0
# Urgent high styling.
typeset -g POWERLEVEL9K_NIGHTSCOUT_URGENT_HIGH_BACKGROUND=1
typeset -g POWERLEVEL9K_NIGHTSCOUT_URGENT_HIGH_FOREGROUND=7
# Old reading styling.
typeset -g POWERLEVEL9K_NIGHTSCOUT_OLD_BACKGROUND=243
typeset -g POWERLEVEL9K_NIGHTSCOUT_OLD_FOREGROUND=0
# Custom icon.
# typeset -g POWERLEVEL9K_NIGHTSCOUT_VISUAL_IDENTIFIER_EXPANSION='⭐'

zmodload -F zsh/net/socket b:zsocket

typeset -g NIGHTSCOUT_THRESHOLD_OLD_SECS=$(( NIGHTSCOUT_THRESHOLD_OLD_MINS * 60 ))

# Creates segment with Nightscout blood sugar data.
#
# Example output: 120 → -1 [1m]
function prompt_nightscout() {
  emulate -L zsh
  setopt err_return

  if [[ -S "$NIGHTSCOUT_SOCKET" ]]; then
    # Read socket into local variables.
    typeset REPLY bgnow arrow delta timeago relative
    zsocket "$NIGHTSCOUT_SOCKET" 2>/dev/null
    {
      IFS=, read -t 0.25 -r bgnow arrow delta timeago relative <&$REPLY
    } always {
      exec {REPLY}>&-
    }

    # State file is invalid. Segment will be hidden.
    if [[ -z "$bgnow" || -z "$relative" ]]; then
      p10k segment -c ''
      return
    fi

    # Choose current state for styling.
    if (( relative > NIGHTSCOUT_THRESHOLD_OLD_SECS )); then
      typeset state=OLD
    elif (( bgnow <= NIGHTSCOUT_THRESHOLD_URGENT_LOW )); then
      typeset state=URGENT_LOW
    elif (( bgnow < NIGHTSCOUT_THRESHOLD_LOW )); then
      typeset state=LOW
    elif (( bgnow < NIGHTSCOUT_THRESHOLD_IN_RANGE )); then
      typeset state=IN_RANGE
    elif (( bgnow < NIGHTSCOUT_THRESHOLD_HIGH )); then
      typeset state=HIGH
    else
      typeset state=URGENT_HIGH
    fi

    # Generate text
    typeset text="$bgnow"
    [[ "$NIGHTSCOUT_SHOW_ARROW" == true ]] && text+=" $arrow"
    [[ "$NIGHTSCOUT_SHOW_DELTA" == true && -n "$delta" ]] && text+=" $delta"
    [[ "$NIGHTSCOUT_SHOW_TIMESTAMP" == true ]] && text+=" [${timeago}]"

    # Write segment.
    p10k segment -s "$state" -i $'\UF058C' -t "$text"
  else
    # State file does not exist. Segment will be hidden.
    p10k segment -c ''
  fi
}
