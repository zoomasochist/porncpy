# Porncpy

Infest your computer with porn :)

## Download

Artifacts are produced for every tag [here](https://github.com/zoomasochist/porncpy/releases)

## Usage

Run from a terminal. On first run it will create a cache
of every path in --root and use that in the future to
speed things up. **if --root is overridden after this point,
provide --refresh too or it will be ignored**.

Flags:
- --refresh: re-build the cache of directories
- --root: specify the root path. porn images will be copied randomly to
subdirectories of this path
- --every: interval between copying images, in milliseconds.
- --count: number of images to copy per interval
- --debug: print additional debug information