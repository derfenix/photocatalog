# Simple photo cataloguer

Just copy/hardlink photos (or video, or any other files) from one place to
another, separating them in sub-directories like `$ROOT/year/month/day/`.

### TL;DR 

I have a smartphone, I have a Syncthing ~~uugh... SmartThing~~ and all photos
from smartphone nicely synced to my PC without my attention. But I can't just
keep all photos in synced folder: if I'll clean my phone memory — all photos 
from pc will be cleaned too. I need to not forget copy files in another 
place before cleaning phone's memory. Also, I can't just drop all photos in 
one dir — I will not find anything there later, and a folder with thousands 
photos looks like a bad idea from either side.
So I create this tool in one evening. All it does — copy (or create hardlinks for)
files from one place to another, creating basic date-aware directories
structure for that files.

Created for own usage and used for a long time without any troubles. But if you meet some — 
you are welcomed to the issues.

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

