# shreddit-go

`shreddit` is a small command-line utility for deleting Reddit data. Deleting a 
Reddit account will not delete comments or submissions - it will only
disassociate your account from them.

You can use `shreddit` to overwrite your comments with text before deleting them
to ensure that the originals are (probably) not preserved.

If you don't want your post history to follow you around forever, you can use 
`shreddit` on a cron job.

If you're deactivating your account, you can run `shreddit` first to ensure your
posts are deleted.

## Overview

`shreddit` is a Go implementation of the popular 
[original Python version](https://github.com/x89/Shreddit) and the subsequent 
[Rust fork](https://github.com/andrewbanchich/shreddit).

I've tried to keep the command-line interface as similar as possible to the
Rust version, but there are some subtle differences. Please see the help output
for the most up-to-date information.

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

## Usage

```bash
shreddit --help
```

### Create Reddit App Credentials

1. Navigate to Reddit -> preferences -> apps (tab) and click `create another app...`
    - Access page directly at https://www.reddit.com/prefs/apps
2. Give the app a name, like 'shreddit`. The name doesn't matter.
3. Select `script`.
4. Set the redirect URL to be `http://localhost:8080`.
5. Click `create app`.

This will provide with a client ID and client secret. The `CLIENT_ID` value used
by `shreddit` is shown under the name of the app you created. The
`CLIENT_SECRET` is shown after clicking `edit`.

> IMPORTANT: TOTP is not supported at this time. If you have 2FA enabled, you
> will need to disable it to use `shreddit`.


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

To run [GoReleaser](https://goreleaser.com/) locally (for debugging):

```bash
goreleaser release --snapshot --clean --verbose --skip=publish
```
