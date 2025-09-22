╔═══╗╔═══╗╔═══╗╔═╗ ╔╗
║╔═╗║║╔═╗║║╔═╗║║║╚╗║║
║╚══╗║╚═╝║║║ ║║║╔╗╚╝║
╚══╗║║╔══╝║╚═╝║║║╚╗║║
║╚═╝║║║   ║╔═╗║║║ ║║║
╚═══╝╚╝   ╚╝ ╚╝╚╝ ╚═╝

# pbar - Command-line Progress Bar Tool

Add beautiful and informative progress bars to your Bash or Zsh scripts with ease.

## Why pbar?

When running long-duration scripts, it's crucial to provide visual feedback to the user. While simple spinners or dots can work, a well-designed progress bar offers a much clearer indication of progress, estimated time remaining, and throughput.

`pbar` simplifies the creation of such progress bars, providing a clean, intuitive, and highly customizable command-line interface. It adheres to the Unix philosophy: do one thing well, and work seamlessly with other tools via pipes.

## Command Reference

`pbar` accepts `current` and `total` as positional arguments, and uses flags for customization.

### Global Flags

*   **`-width`**: Specifies the width of the progress bar (default: 40).
*   **`-style`**: Sets the style of the progress bar. Available: `classic`, `block`, `spinner`, `arrow`, `braille`, `custom` (default: `classic`).
*   **`-indeterminate`**: Renders an animated spinner without a percentage, for tasks where the total is unknown.
*   **`-colorbar`**: Sets the color for the progress bar itself (e.g., `green`, `red`). Available colors: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`.
*   **`-colortext`**: Sets the color for the text (percentage, ETA, throughput) (e.g., `yellow`, `blue`). Available colors: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`.
*   **`-finished`**: Renders the bar in a completed state (e.g., `[✔] 100%`).
*   **`-quiet`**: Outputs only the percentage, useful for piping into other tools.
*   **`-chars`**: Custom characters for the progress bar (e.g., `#=`). Requires `--style=custom`.
*   **`-version`**: Prints version information and exits.

## Installation

`pbar` provides flexible installation options.

### Quick Install (Recommended)

This single command will download and install `pbar` to a sensible default location for your system.

**User-level Installation (Recommended for most users):**
Installs `pbar` to `$HOME/.local/bin` (Linux/macOS) or a user-specific `bin` directory (Windows).

```bash
curl -sSfL https://raw.githubusercontent.com/gregory-chatelier/pbar/main/install.sh | sh
```

**System-wide Installation (Requires `sudo`):**
Installs `pbar` to `/usr/local/bin` (Linux/macOS).

```bash
sudo curl -sSfL https://raw.githubusercontent.com/gregory-chatelier/pbar/main/install.sh | sh
```

### Custom Installation Directory

You can specify a custom installation directory using the `INSTALL_DIR` environment variable:

```bash
INSTALL_DIR=$HOME/my-tools curl -sSfL https://raw.githubusercontent.com/gregory-chatelier/pbar/main/install.sh | sh
```

### From Source

If you have Go installed (Go 1.24+ is required):

```bash
go install github.com/gregory-chatelier/pbar@latest
```

## Common Usage

### Basic Progress

Show 25% completion out of 100.

```bash
pbar 25 100
```

### Looping Progress

Integrate `pbar` into a Bash loop.

```bash
for i in {1..100}; do
  sleep 0.1
  pbar $i 100 --style=block --colorbar=green
done
```

### Indeterminate Progress

For tasks where the total is unknown.

```bash
my_long_command | pbar --indeterminate --colortext=cyan
```

### Custom Bar Characters

Use your own characters for the progress bar.

```bash
pbar 70 100 --style=custom --chars="█ " --colorbar=magenta
```

### Quiet Mode for Scripting

Get only the percentage for further processing.

```bash
progress_percent=$(pbar 80 100 --quiet)
echo "Current progress: $progress_percent"
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
