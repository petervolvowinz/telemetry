#!/bin/bash

set -e

# List of main Go program files (one per package)
PROGRAMS=(
  "server/server.go"
  "mqttlistener/mqttlistener.go"
  "client/client.go"
  "coapobserver/coapobserver.go"
)

# Colors for output (red, green, blue)
COLORS=("red" "green" "blue")

# Function to get ANSI color code by color name
get_color_code() {
  case "$1" in
    red) echo "\033[0;31m" ;;
    green) echo "\033[0;32m" ;;
    blue) echo "\033[0;34m" ;;
    *) echo "\033[0m" ;;
  esac
}

# Ensure bin directory exists
mkdir -p bin

i=0
for path in "${PROGRAMS[@]}"; do
  name=$(basename "$path" .go)
  dir=$(dirname "$path")
  output="bin/$name"

  # Build entire package directory, not single file
  if [[ ! -f "$output" || "$dir" -nt "$output" ]]; then
    echo "Building package in $dir → $output"
    go build -o "$output" "./$dir"
  else
    echo "Up to date: $output"
  fi

  color_code=$(get_color_code "${COLORS[$i]}")
  reset="\033[0m"

  # Prepare shell command to run in new Terminal window
  cmd="cd \"$(pwd)\"; echo -e '${color_code}>>> Running $output${reset}'; ./$output; echo Done; bash"

  # Escape backslashes and double quotes for AppleScript string
  cmd_escaped=$(printf '%s' "$cmd" | sed -e 's/\\/\\\\/g' -e 's/"/\\"/g')

  # Use osascript with separate -e calls to avoid multiline quoting issues
  osascript -e "tell application \"Terminal\" to activate" \
            -e "tell application \"Terminal\" to do script \"$cmd_escaped\""

  i=$(( (i + 1) % ${#COLORS[@]} ))
done
