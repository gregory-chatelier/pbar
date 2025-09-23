#!/bin/bash

# This script demonstrates the pbar utility in various modes.
# It is designed for Unix-like environments (Linux, macOS, WSL, Git Bash).
# On some Windows Git Bash setups, direct execution of .exe files from the current directory
# might require adjustments (e.g., adding the current directory to PATH or using 'winpty').
# If you encounter 'command not found' errors, try running 'export PATH=$PATH:.' before executing,
# or use 'winpty ./pbar.exe' for individual calls if your terminal supports it.

# Exit immediately if a command exits with a non-zero status.
set -euo pipefail
set -x # Enable debugging

# --- Configuration ---
PBAR_BIN="./pbar.exe" # Path to the pbar executable
DEMO_SPEED=0.1    # Adjust for faster/slower demo (seconds per update)

# --- Helper Functions ---

PBAR_PID="" # Global to store the pbar process ID
PBAR_COPROC_PID="" # Global to store the coproc PID

# Function to clean up background processes and terminal on exit or interrupt
cleanup() {
  echo "\nCleaning up..."
  if [ -n "$PBAR_PID" ]; then
    kill "$PBAR_PID" 2>/dev/null || true # Kill pbar if running
    wait "$PBAR_PID" 2>/dev/null || true # Wait for it to terminate
  fi
  # Close coproc file descriptors if they exist
  if [ -n "$PBAR_COPROC_PID" ]; then
    eval "exec ${PBAR_COPROC[0]}<&-"
    eval "exec ${PBAR_COPROC[1]}>&-"
  fi
  echo -e "\033[?25h" # Show cursor
  echo "Demo finished."
}

# Trap SIGINT (Ctrl+C) and SIGTERM to call the cleanup function
trap cleanup EXIT

# Function to send JSON updates to a background pbar process via coproc
send_pbar_json_update() {
  if [ -n "${PBAR_COPROC[1]}" ]; then
    echo "$1" >&"${PBAR_COPROC[1]}"
  else
    echo "Error: pbar coproc not running for parallel update." >&2
  fi
}

# Function to display a section header
section_header() {
  echo -e "\n\033[1m--- $1 ---\033[0m\n"
  sleep 1
}

# --- Main Demo Logic ---

clear # Clear the terminal for a clean demo start

section_header "Single Bar Examples"

# Classic Bar
section_header "Classic Bar"
for i in $(seq 0 10 100); do
  echo "Calling pbar with arguments: $i 100"
  "$PBAR_BIN" "$i" 100
  sleep "$DEMO_SPEED"
done
sleep 0.5

# Block Bar
section_header "Block Bar"
for i in $(seq 0 10 100); do
  echo "Calling pbar with arguments: --style=block $i 100"
  "$PBAR_BIN" --style=block "$i" 100
  sleep "$DEMO_SPEED"
done
sleep 0.5

# Arrow Bar
section_header "Arrow Bar"
for i in $(seq 0 10 100); do
  echo "Calling pbar with arguments: --style=arrow $i 100"
  "$PBAR_BIN" --style=arrow "$i" 100
  sleep "$DEMO_SPEED"
done
sleep 0.5

# Braille Bar
section_header "Braille Bar"
for i in $(seq 0 10 100); do
  echo "Calling pbar with arguments: --style=braille $i 100"
  "$PBAR_BIN" --style=braille "$i" 100
  sleep "$DEMO_SPEED"
done
sleep 0.5

# Custom Bar
section_header "Custom Bar (chars='*.')"
for i in $(seq 0 10 100); do
  echo "Calling pbar with arguments: --style=custom --chars='*.' --colorbar=yellow $i 100"
  "$PBAR_BIN" --style=custom --chars='*.' --colorbar=yellow "$i" 100
  sleep "$DEMO_SPEED"
done
sleep 0.5

# Color Support
section_header "Color Support (Bar: Green, Text: Yellow)"
for i in $(seq 0 10 100); do
  echo "Calling pbar with arguments: --colorbar=green --colortext=yellow $i 100"
  "$PBAR_BIN" --colorbar=green --colortext=yellow "$i" 100
  sleep "$DEMO_SPEED"
done
sleep 0.5

# Finished State
section_header "Finished State"
echo "Calling pbar with arguments: --finished --message=\"Task Complete!\" 100 100"
"$PBAR_BIN" --finished --message="Task Complete!" 100 100
sleep 1

# Indeterminate Mode
section_header "Indeterminate Mode"
echo "Calling pbar with arguments: --indeterminate"
"$PBAR_BIN" --indeterminate & # Run in background
INDETERMINATE_PID=$! # Store PID
sleep 3 # Let it run for a bit
kill "$INDETERMINATE_PID" 2>/dev/null || true # Kill it
wait "$INDETERMINATE_PID" 2>/dev/null || true # Wait for it to terminate
sleep 0.5

section_header "Parallel Bar Example (JSON Input)"

# Start pbar in parallel mode using coproc
coproc PBAR_COPROC { "$PBAR_BIN" --parallel; }
PBAR_PID=$! # Store the PID of the pbar process
PBAR_COPROC_PID=$! # Assign to global variable
sleep 0.5 # Give pbar a moment to start

# Simulate multiple tasks
MAX_UPDATES=20
for i in $(seq 1 "$MAX_UPDATES"); do
  current_A=$((i * 5))
  current_B=$((i * 3))
  current_C=$((i * 7))

  # Task A: Downloading
  if [ "$current_A" -le 100 ]; then
    send_pbar_json_update "{\"id\": \"task_A\", \"current\": $current_A, \"total\": 100, \"message\": \"Downloading file A...\"}"
  elif [ "$current_A" -gt 100 ] && [ "$current_A" -lt 105 ]; then # A little buffer to ensure finished state is sent
    send_pbar_json_update "{\"id\": \"task_A\", \"current\": 100, \"total\": 100, \"finished\": true, \"message\": \"File A downloaded!\"}"
  fi

  # Task B: Processing
  if [ "$current_B" -le 100 ]; then
    send_pbar_json_update "{\"id\": \"task_B\", \"current\": $current_B, \"total\": 100, \"style\": \"block\", \"colorbar\": \"green\", \"message\": \"Processing data B...\"}"
  elif [ "$current_B" -gt 100 ] && [ "$current_B" -lt 105 ]; then
    send_pbar_json_update "{\"id\": \"task_B\", \"current\": 100, \"total\": 100, \"finished\": true, \"style\": \"block\", \"colorbar\": \"green\", \"message\": \"Data B processed!\"}"
  fi

  # Task C: Uploading (slower)
  if [ "$current_C" -le 100 ]; then
    send_pbar_json_update "{\"id\": \"task_C\", \"current\": $current_C, \"total\": 100, \"style\": \"arrow\", \"colorbar\": \"blue\", \"message\": \"Uploading results C...\"}"
  elif [ "$current_C" -gt 100 ] && [ "$current_C" -lt 105 ]; then
    send_pbar_json_update "{\"id\": \"task_C\", \"current\": 100, \"total\": 100, \"finished\": true, \"style\": \"arrow\", \"colorbar\": \"blue\", \"message\": \"Results C uploaded!\"}"
  fi

  sleep "$DEMO_SPEED"
done

sleep 2 # Allow final states to be displayed

# The trap EXIT will handle killing PBAR_PID and showing the cursor
