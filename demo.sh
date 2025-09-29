#!/bin/bash

# This script demonstrates the pbar utility in various modes.
# It is designed for Unix-like environments (Linux, macOS, WSL, Git Bash).

# Exit immediately if a command exits with a non-zero status.
set -euo pipefail

# --- Configuration ---
PBAR_BIN="./pbar/main.exe" # Path to the pbar executable
# if [ -f "./pbar" ]; then
#     PBAR_BIN="./pbar"
# fi
DEMO_SPEED=0.02    # Adjust for faster/slower demo (seconds per update)

# --- 80s Retro Colors ---
C_BLACK='\033[0;30m'
C_RED='\033[0;31m'
C_GREEN='\033[0;32m'
C_YELLOW='\033[0;33m'
C_BLUE='\033[0;34m'
C_MAGENTA='\033[0;35m'
C_CYAN='\033[0;36m'
C_WHITE='\033[0;37m'
C_BRIGHT_BLACK='\033[0;90m'
C_BRIGHT_RED='\033[0;91m'
C_BRIGHT_GREEN='\033[0;92m'
C_BRIGHT_YELLOW='\033[0;93m'
C_BRIGHT_BLUE='\033[0;94m'
C_BRIGHT_MAGENTA='\033[0;95m'
C_BRIGHT_CYAN='\033[0;96m'
C_BRIGHT_WHITE='\033[0;97m'
C_BOLD='\033[1m'
C_RESET='\033[0m'

# --- Helper Functions ---

# Function to clean up background processes and terminal on exit or interrupt
cleanup() {
  echo -e "\n${C_BOLD}${C_BRIGHT_RED}Cleaning up...${C_RESET}"
  echo -e "\033[?25h" # Show cursor
  echo -e "${C_BOLD}${C_BRIGHT_GREEN}Demo finished.${C_RESET}"
}

# Trap SIGINT (Ctrl+C) and SIGTERM to call the cleanup function
trap cleanup EXIT

# Function to display a section header
section_header() {
  echo -e "\n${C_BOLD}${C_BRIGHT_CYAN}////////////////////////////////////////////////////////////////${C_RESET}"
  echo -e "${C_BOLD}${C_BRIGHT_CYAN}// $1 ${C_RESET}"
  echo -e "${C_BOLD}${C_BRIGHT_CYAN}////////////////////////////////////////////////////////////////${C_RESET}"
  sleep 0.5
}

# Function to run a demo loop
run_demo() {
    local title="$1"
    local cmd="$2"
    section_header "$title"
    for i in $(seq 0 5 100); do
        eval "$cmd"
        sleep "$DEMO_SPEED"
    done
    sleep 0.75
    clear
}

# --- Main Demo Logic ---

clear # Clear the terminal for a clean demo start

# --- Intro ---
echo -e "${C_BOLD}${C_BRIGHT_YELLOW}PBAR demo-ing...${C_RESET}"
echo -e "\n${C_BOLD}${C_BRIGHT_MAGENTA}THE PROGRESS BAR EXPERIENCE${C_RESET}\n"
sleep 1.5

clear

# --- Basic Styles ---
run_demo "Classic Style" "\"$PBAR_BIN\" \$i 100 --style=classic --colorbar=magenta"
run_demo "Block Style"   "\"$PBAR_BIN\" \$i 100 --style=block --width=60 --colorbar=cyan"
run_demo "Arrow Style"   "\"$PBAR_BIN\" \$i 100 --style=arrow --width=30 --colorbar=yellow"
run_demo "Braille Style" "\"$PBAR_BIN\" \$i 100 --style=braille --width=40 --colorbar=green"

# --- Indeterminate Mode ---
run_demo "Spinner"         "\"$PBAR_BIN\" \$i 100 --style=spinner --colortext=magenta"
run_demo "Braille Spinner" "\"$PBAR_BIN\" \$i 100 --style=braille-spinner --colortext=cyan"

# --- Metadata ---
run_demo "No Metadata"  "\"$PBAR_BIN\" \$i 100 --id=meta-none --show-elapsed=false --show-throughput=false --show-eta=false"

# --- Finished State ---
run_demo "Finished Message" "\"$PBAR_BIN\" \$i 100 --finished-message='DOWNLOAD COMPLETE'"

# --- Custom Characters ---
run_demo "Custom Chars: '| '"  "\"$PBAR_BIN\" \$i 100 --style=custom --chars='| ' --colorbar=red"


# --- Parallel Downloads ---
section_header "Parallel Downloads"

# This function generates a stream of JSON objects for pbar.
# Each object represents an update for a specific progress bar.
# In this example, we simulate three download tasks running in parallel.
run_all_steps() {
    # We loop from 0 to 100 with a step of 4 to generate progress updates.
    for i in $(seq 0 4 100); do
        # Task 1: A 100MB file, progresses with the loop.
        echo "{\"id\": \"File1.zip\", \"current\": $i, \"total\": 100, \"style\": \"block\", \"colorbar\": \"green\"}"

        # Task 2: An 80MB file, progresses until it reaches its total.
        if [ "$i" -le 80 ]; then
            echo "{\"id\": \"File2.iso\", \"current\": $i, \"total\": 80, \"style\": \"block\", \"colorbar\": \"cyan\"}"
        fi

        # Task 3: A 60MB file, progresses until it reaches its total.
        if [ "$i" -le 60 ]; then
            echo "{\"id\": \"File3.tar.gz\", \"current\": $i, \"total\": 60, \"style\": \"block\", \"colorbar\": \"magenta\"}"
        fi

        # A short sleep to make the progress visible.
        sleep "$DEMO_SPEED"
    done

    # After the loop, we send a "finished" update for each task.
    echo '{\"id\": \"File1.zip\", \"finished\": true}'
    echo '{\"id\": \"File2.iso\", \"finished\": true}'
    echo '{\"id\": \"File3.tar.gz\", \"finished\": true}'
}

# The main pipeline for the parallel demo:
# 1. `run_all_steps` generates the JSON stream.
# 2. The output of `run_all_steps` is piped to `pbar --parallel`
run_all_steps | "$PBAR_BIN" --parallel


clear
