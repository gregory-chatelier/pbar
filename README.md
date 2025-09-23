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

### Advanced Features

- **Color Support**: Allows users to set colors for the bar, background, and text for a high-impact visual style.
    - **Example**: `pbar 75 100 --color-bar=green --color-text=yellow`
- **Finished State**: Defines a distinct appearance for the bar upon completion (e.g., a checkmark and a solid color) to provide clear visual confirmation.
    - **Example**: On completion, the bar could change to `[✔] Download Complete! 100%`.
- **Indeterminate Mode**: For tasks where the total is unknown, a special mode displays an animated indicator (e.g., a spinner) without a percentage.
    - **Example**: `my_command | pbar --indeterminate`
- **Parallel Mode**: Supports rendering multiple progress bars simultaneously, each updated via a stream of JSON objects from standard input. This is ideal for orchestrating complex, concurrent tasks.
    - **Usage**: Activate with the `--parallel` flag. Input is a stream of JSON objects, one per line, each representing an update for a specific bar.
    - **Example Input (JSON per line)**:
        ```json
        {"id": "task1", "current": 10, "total": 100, "message": "Processing task 1..."}
        {"id": "task2", "current": 25, "total": 50, "style": "block", "colorBar": "blue"}
        {"id": "task1", "current": 20, "total": 100}
        {"id": "task2", "finished": true, "message": "Task 2 complete!"}
        ```

### Panel of Styles

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
