# pbar - Command-line Progress Bar Tool

Add progress bars to your shell scripts with ease

## Why pbar?

When running long-duration scripts, it's crucial to provide visual feedback to the user. While simple spinners or dots can work, a well-designed progress bar offers a much clearer indication of progress, estimated time remaining, and throughput.

`pbar` simplifies the creation of such progress bars, providing a clean, intuitive, and highly customizable command-line interface.

## Command Reference

`pbar` accepts `current` item number and `total` item number as positional arguments, and uses flags for customization.

### Basic Renderer

Show 25% completion progress, 25th item out of 100.

```bash
pbar 25 100
```

### Task Progress

Integrate `pbar` into your task loop for the animation to run.

```bash
for i in {1..100}; do
  sleep 0.1
  pbar $i 100 --style=block --colorbar=green
done
```

To experience `pbar` in action showcasing various styles, colors, and the powerful parallel mode, run the `demo.sh` script:

```bash
bash demo.sh
```

**Note on Compatibility:** Set the path to the executable in the demo script before running

[![asciicast](https://asciinema.org/a/bFPewvgj4ilGI94NP5oxKWajT.svg)](https://asciinema.org/a/bFPewvgj4ilGI94NP5oxKWajT)

### Advanced Features

- **Metadata Display**: Control the visibility of elapsed time, throughput, and estimated time remaining.
    - **Example (Hide all metadata)**: `pbar 50 100 --show-elapsed=false --show-throughput=false --show-eta=false`
- **Color Support**: Allows users to set colors for the bar, background, and text for a high-impact visual style.
    - **Example**: `pbar 75 100 --colorbar=green --colortext=yellow`
- **Finished State**: Defines a distinct appearance for the bar upon completion (e.g., a checkmark and a solid color) to provide clear visual confirmation.
    - **Example**: On completion, the bar could change to `[✔] Download Complete! 100%`.
- **Indeterminate Mode**: For tasks where the total is unknown, a special mode displays an animated indicator (e.g., a spinner) without a percentage.
  
    - **Example**: `pbar --style=spinner`
- **Parallel Mode**: Supports rendering multiple progress bars simultaneously, each updated via a stream of JSON objects from standard input. This is useful for orchestrating complex, concurrent tasks.
  
    - **Usage**: Activate with the `--parallel` flag. Input is a stream of JSON objects, one per line, each representing an update for a specific bar.
    - **Example Input (JSON per line)**:
      
        ```json
        {"id": "File1.zip", "current": 10, "total": 100, "message": "Downloading File1.zip", "style": "block", "colorbar": "green"}
        {"id": "File2.iso", "current": 5, "total": 80, "message": "Downloading File2.iso", "style": "block", "colorbar": "cyan"}
        {"id": "File3.tar.gz", "current": 2, "total": 60, "message": "Downloading File3.tar.gz", "style": "block", "colorbar": "magenta"}
        ```

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
curl -sSfL https://raw.githubusercontent.com/gregory-chatelier/pbar/main/install.sh | INSTALL_DIR=$HOME/bin sh
```

### From Source

If you have Go installed (Go 1.24+ is required):

```bash
go install github.com/gregory-chatelier/pbar@latest
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
