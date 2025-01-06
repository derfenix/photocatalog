# Effortless Photo Organizer

[![Go](https://github.com/derfenix/photocatalog/actions/workflows/go.yml/badge.svg)](https://github.com/derfenix/photocatalog/actions/workflows/go.yml)

A simple tool to organize your photos, videos, or other files by copying or hardlinking them into a date-based directory structure like `$ROOT/year/month/day/`.

## TL;DR

I use a smartphone and Syncthing to automatically sync my photos to my PC. However, if I clean up my phone's memory, the synced photos on my PC are deleted as well.
Dumping everything into one folder wasn't an option — finding anything later would be a nightmare. 

To avoid this, I needed a solution to back up and organize my photos without manual effort. So, I built this tool in one evening to solve the problem. It has worked flawlessly for me and might help you too. If you encounter any issues, feel free to open a ticket — I'll do my best to assist.

## Installation

Install the tool via `go`:

```bash
go install github.com/derfenix/photocatalog/v2@latest
```

Optionally, copy the binary to a directory in your system or user's `$PATH` (e.g., `/usr/local/bin`):

```bash
sudo cp ${GOPATH}/bin/photocatalog /usr/local/bin/photocatalog
```

## Organization Modes

The tool supports the following organization modes:

- **copy** — Copies files to the target directory. If the filesystem supports it, uses Copy-on-Write (COW) for efficiency.
- **hardlink** — Creates hardlinks to the source files, saving disk space. Ideal if the source and target are on the same partition, though file permissions remain linked to the original.
- **move** — Moves files from the source to the target directory.
- **symlink** — Creates symbolic links at the target pointing to the source files.

## Supported Formats

- **JPEG and TIFF files** with valid EXIF metadata.
- Files named in the format `yyyymmdd_HHMMSS.ext` (optionally with suffixes after the timestamp) (e.g., `20230101_123456.jpg`). This format is common in Android cameras and other devices.

If a file lacks EXIF data, the tool falls back to parsing the filename.

Currently, the timestamp format is not customizable. Let me know if support for additional formats is required.

## Usage

Arguments
```shell
  -dir-mode string
        Mode bits for directories can be created while syncing (default "0777")
  -file-mode string
        Mode bits for files created while syncing (not applicable for hardlink mode) (default "0644")
  -mode string
        Mode (default "hardlink")
  -overwrite
        Overwrite existing files
  -skip-full-sync
        Skip full sync at startup
  -source string
        Source directory
  -target string
        Target directory
  -watch
        Watch for changes in the source directory (default true)

```

`-skip-full-sync` and `-watch` are not compatible.

`-source` and `-target` are required.


## Examples

### One-Time Run

#### Copy Files
```shell
photocatalog -mode copy -target ./photos/ -source ./sync/photos/
```

#### Create Hardlinks
```shell
photocatalog -mode hardlink -target ./photos/ -source ./sync/photos/
```

### Watch Mode

Enable continuous monitoring of a source directory:

#### Copy Files
```shell
photocatalog -mode copy -target ./photos -watch -source ./sync/photos/
```

#### Create Hardlinks
```shell
photocatalog -mode hardlink -target ./photos/ -watch -source ./sync/photos/
```

## Running as a Service

### Systemd Setup

Install and configure the service:
```shell
sh ./init/install_service.sh systemd
```

This will:

1. Install a systemd unit file.
2. Create a configuration stub at `$HOME/.config/photocatalog`.
3. Open the config file for editing.

Enable and start the service:
```shell
systemctl --user enable --now photocatalog
```

Now, files added to the monitored directory (`MONITOR` in the config) will automatically be organized into the target directory under the corresponding subdirectories.

## FAQ

### Why did you create this tool when awesome tool XXX already exists?
Two reasons:
1. I wanted to.
2. I could.
