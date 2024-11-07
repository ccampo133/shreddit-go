# shreddit-go

TODO

## Overview

TODO

## Usage

TODO

## Installation

### From Source

```bash
go install github.com/ccampo133/shreddit-go
```

### Release Artifacts

Download the latest release zip for your platform from the
[releases page](https://github.com/ccampo133/shreddit-go/releases) and extract
the binary to a location in your `PATH`.

```bash
unzip shreddit_*.zip
mv shreddit /usr/local/bin
```

The SHA256 checksums for each release are provided in a file named
`shreddit_<version>_checksums.txt`. You can verify the integrity of the
downloaded binary by comparing its checksum to the one in the file. The
checksums are also signed with [my GPG key](https://github.com/ccampo133.gpg),
and you can verify the checksums file, e.g.:

```bash
# Replace <version> with the desired version without the 'v' prefix, e.g. 0.1.0.
# The below commands assume that you are in the same directory as the binary and
# checksums/signature.
sha256sum -c shreddit_<version>_checksums.txt
# First import my GPG key if you haven't already:
# curl https://github.com/ccampo133.gpg | gpg --import
gpg --verify shreddit_<version>_checksums.txt.sig shreddit_<version>_checksums.txt
````

### Docker

You can run the image directly:

```bash
docker run --rm ghcr.io/ccampo133/shreddit:latest
```

Tags for each version of `shreddit` are released, as well as a `latest` tag.

## Development

There is a [`Makefile`](Makefile) with some common development tasks. Please see
the file for more information. It's a pretty standard Go project - there's not
much to it.

To build (requires Go 1.23+):

```bash
make build
```

To run tests:

```bash
make test
```

To build a local Docker image called `shreddit`:
```bash
make docker-build
```
