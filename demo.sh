#!/bin/bash

# This script demonstrates the pbar utility in various modes.
# It is designed for Unix-like environments (Linux, macOS, WSL, Git Bash).
# On some Windows Git Bash setups, direct execution of .exe files from the current directory
# might require adjustments (e.g., adding the current directory to PATH or using 'winpty').
# If you encounter 'command not found' errors, try running 'export PATH=$PATH:.' before executing,
# or use 'winpty ./pbar.exe' for individual calls if your terminal supports it.

# Exit immediately if a command exits with a non-zero status.
set -euo pipefail
# set -x # Enable debugging

# --- Configuration ---
PBAR_BIN="./pbar/main.exe" # Path to the pbar executable
DEMO_SPEED=0.05    # Adjust for faster/slower demo (seconds per update)

# --- Helper Functions ---

PBAR_PID="" # Global to store the pbar process ID
PBAR_COPROC_PID="" # Global to store the coproc PID

# Function to clean up background processes and terminal on exit or interrupt
cleanup() {
  echo -e "\nCleaning up..."
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
  echo -e "\n\033[1m--- $1 ---\033[0m"
  sleep 1
}

# --- Main Demo Logic ---

clear # Clear the terminal for a clean demo start


# Conditional Metadata Display
# section_header "Conditional Metadata (All Shown by Default)"
# for i in $(seq 0 5 100); do
#   "$PBAR_BIN" "$i" 100
#   sleep "$DEMO_SPEED"
# done
# sleep 0.5

# section_header "Conditional Metadata (Elapsed Hidden)"
# for i in $(seq 0 5 100); do
#   "$PBAR_BIN" "$i" 100 --show-elapsed=false 
#   sleep "$DEMO_SPEED"
# done
# sleep 0.5

# section_header "Conditional Metadata (Throughput Hidden)"
# for i in $(seq 0 5 100); do
#   "$PBAR_BIN" --id="task1" --show-throughput=false "$i" 100
#   sleep "$DEMO_SPEED"
# done
# sleep 0.5

# section_header "Conditional Metadata (ETA Hidden)"
# for i in $(seq 0 5 100); do
#   "$PBAR_BIN" --show-eta=false "$i" 100
#   sleep "$DEMO_SPEED"
# done
# sleep 0.5

section_header "Conditional Metadata (All Hidden)"
for i in $(seq 0 5 100); do
  "$PBAR_BIN" --show-elapsed=false --show-throughput=false --show-eta=false "$i" 100
  sleep "$DEMO_SPEED"
done
sleep 0.5


# # Default Spinner
# section_header "Default Spinner (Indeterminate)"
# for i in $(seq 0 5 100); do
#   "$PBAR_BIN" --style=spinner "$i" 100
#   sleep "$DEMO_SPEED"
# done
# sleep 0.5

# # Braille Spinner
# section_header "Braille Spinner"
# for i in $(seq 0 5 100); do
#   "$PBAR_BIN" --style=braille-spinner "$i" 100
#   sleep "$DEMO_SPEED"
# done
# sleep 0.5

# # Finished State
# section_header "Finished State"
# for i in $(seq 0 10 100); do
#   "$PBAR_BIN" --finished-message="Done!" "$i" 100
# sleep "$DEMO_SPEED"
# done
# sleep 0.5

# # Message Bar
# section_header "Message Bar"
# for i in $(seq 0 10 100); do
#   "$PBAR_BIN" --message="Processing item $i..." "$i" 100
#   sleep "$DEMO_SPEED"
# done
# sleep 0.5

# # Classic Bar
# section_header "Classic Bar"
# for i in $(seq 0 10 100); do
#   "$PBAR_BIN" "$i" 100
#   sleep "$DEMO_SPEED"
# done
# sleep 0.5

# # Block Bar
# section_header "Block Bar"
# for i in $(seq 0 10 100); do
#   "$PBAR_BIN" --style=block "$i" 100
#   sleep "$DEMO_SPEED"
# done
# sleep 0.5

# # Arrow Bar
# section_header "Arrow Bar"
# for i in $(seq 0 10 100); do
#   "$PBAR_BIN" --style=arrow "$i" 100
#   sleep "$DEMO_SPEED"
# done
# sleep 0.5

# # Braille Bar
# section_header "Braille Bar"
# for i in $(seq 0 10 100); do
#   "$PBAR_BIN" --style=braille "$i" 100
#   sleep "$DEMO_SPEED"
# done
# sleep 0.5

# # Custom Bar
# section_header "Custom Bar (chars='*.')"
# for i in $(seq 0 10 100); do
#   "$PBAR_BIN" --style=custom --chars='*.\' --colorbar=yellow "$i" 100
#   sleep "$DEMO_SPEED"
# done
# sleep 0.5

# # Color Support
# section_header "Color Support (Bar: Green, Text: Yellow)"
# for i in $(seq 0 10 100); do
#   "$PBAR_BIN" --colorbar=green --colortext=yellow "$i" 100
#   sleep "$DEMO_SPEED"
# done
# sleep 0.5

sleep 2 # Allow final states to be displayed

# The trap EXIT will handle killing PBAR_PID and showing the cursor
