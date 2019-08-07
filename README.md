# Simple photo cataloguer

Just copy/hardlink photos (or video, or any other files) from one place to
another, separating them in sub-directories like `$ROOT/year/month/day/`.

### TL;DR 

I have a smartphone, I have a Syncthing ~~uugh... SmartThing~~ and all photos
from smartphone nicely synced to my PC without my attention. But I can't just
keep all photos in synced folder: if I'll clean my phone memory - all photos 
from pc will be cleaned too. I need to not forget copy files in another 
place before cleaning phone's memory. Also, I can't  just drop all photos in 
one dir - I will not find anything there later, and a folder with thousands 
photos looks like a bad idea from either side.    
So I create this tool in one evening. All it does - copy (or create hardlinks for)
files from one place to another, creating basic date-aware directories
structure for that files.


## Installing
```bash
go install github.com/derfenix/photocatalog
```
Optionally you could copy created binary from the GO's bin path to 
system or user $PATH, e.g. /usr/local/bin/.
```bash
sudo cp ${GOPATH}/bin/photocatalog /usr/local/bin/photocatalog
```

## Supported formats
At this moment supported jpeg files with filled exif data or any other 
files but with names matching pattern `yyymmdd_HHMMSS.ext`. Such 
names format applied by android's camera software (I guess all cams 
use this format, fix me if I'm wrong).

There is no support for changing names format without modifying  source code 
at this time.

## Usage
### One-shot 
#### Copy files (make a COW if fs supports it)
```bash
photocalog -mode copy -target ./photos/ ./sync/photos/*
```

#### Create hardlinks (only withing one disk partition)
```bash
photocalog -mode hardlink -target ./photos/ ./sync/photos/*
```
or 
```bash
photocalog -target ./photos/ ./sync/photos/*
```

### Monitor
#### Copy files (make a COW if fs supports it)
```bash
photocalog -mode copy -target ./photos/ ./sync/photos/*
```

#### Create hardlinks (only withing one disk partition)
```bash
photocalog -mode hardlink -target ./photos/ -monitor ./sync/photos/
```
or 
```bash
photocalog -target ./photos/ -monitor ./sync/photos/
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
1. I wanted
2. I can

