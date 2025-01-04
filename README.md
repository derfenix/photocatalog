[![Go](https://github.com/derfenix/photocatalog/actions/workflows/go.yml/badge.svg)](https://github.com/derfenix/photocatalog/actions/workflows/go.yml)

# Simple photo cataloguer

Just copy/hardlink photos (or video, or any other files) from one place to
another, separating them in sub-directories like `$ROOT/year/month/day/`.

### TL;DR 

I use a smartphone along with Syncthing to seamlessly sync all my photos to my PC without any manual effort. However, there's a catch: I can't keep all my photos in the synced folder indefinitely. If I clear my phone's memory, the photos on my PC get deleted as well. To avoid this, I need to remember to copy the files to another location before cleaning up my phone.

Simply dumping all my photos into one folder isn't a solution either — finding anything later would be a nightmare, and a folder with thousands of unsorted photos is far from ideal.

To address these issues, I created this tool in just one evening. Its primary purpose is to copy (or create hardlinks for) files from one location to another, while organizing them into a simple, date-based directory structure.

This tool was built for personal use and has been serving me well for quite some time without any problems. However, if you encounter any issues, feel free to report them — I’d be happy to help.

## Installing
```bash
go install github.com/derfenix/photocatalog/v2@latest
```
Optionally you could copy created binary from the GO's bin path to 
system or user $PATH, e.g. /usr/local/bin/.
```bash
sudo cp ${GOPATH}/bin/photocatalog /usr/local/bin/photocatalog
```

## Migrating from v0.*

TODO 

## Organization modes

Next organization modes supported:
    
- **copy** — copy files to target root. Make COW (using syscall) if FS supports it.
- **hardlink** — create hardlink to the source file instead of copying. 
The best choice if source and target are in same partition for compatibility
and resource usage, but we can't chmod target files, because of original file mode will 
be changed too. 
- **move** — moves original files to new place.
- **symlink** — create a symlink at the target for the source files. 

## Supported formats
At this moment supported jpeg and tiff files with filled exif data and any other 
files but with names matching pattern `yyymmdd_HHMMSS.ext` with optional suffixes after a timestamp.
Such names format applied by the Android's camera software (I guess all cams 
use this format, fix me if I'm wrong). 

Jpeg/Tiff files without modification date if exif will be fallen back to the name parsing.

No able to change names format without modifying source code for now. Just because 
I have reasons to believe that this format is the most popular for the application use cases.
But let me know if you need different timestamp formats support.

## Usage
### One-shot 
#### Copy files
```bash
photocalog -mode copy -target ./photos/ -source ./sync/photos/
```

#### Create hardlinks
```bash
photocalog -mode hardlink -target ./photos/ -source ./sync/photos/
```
or 
```bash
photocalog -target ./photos/ -source ./sync/photos/*
```

### Watch mode
#### Copy files
```bash
photocalog -mode copy -target ./photos -watch -source ./sync/photos/
```

#### Create hardlinks
```bash
photocalog -mode hardlink -target ./photos/ -watch -source ./sync/photos/
```
or 
```bash
photocalog -target ./photos/ -watch -source ./sync/photos/
```

## Install and run monitor service

### Systemd
```bash
sh ./init/install_service.sh systemd
```
This command will install unit file, create stub for its config and open
editor to allow you edit configuration. Config file stored at 
`$HOME/.config/photocatalog`.

Then enable and start service
```bash
systemctl --user enable --now photocatalog
```
That's all. Now, if any file will be placed in directory, specified as `MONITOR`
in config file, this file will be copied or hardlinked into the target dir
under corresponding sub-dir. 

## FAQ

### Why this tool was created if there is awesome XXX tool?
I had two good reasons:
1. I wish
2. I can

