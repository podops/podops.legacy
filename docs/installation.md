# Installation

With the Podops command-line interface (CLI), the `po` command, you can create podcasts and manage them from a terminal. This page describes various methods for installing the PoOps command-line interface on your system. 

## Downloads the binary

Download the latest binary for your system:

* [Linux/amd64](https://storage.googleapis.com/cdn.podops.dev/downloads/cli-0.9.7/po-linux-0.9.7.gz)
* [MacOS/amd64](https://storage.googleapis.com/cdn.podops.dev/downloads/cli-0.9.7/po-mac-0.9.7.gz)
* [Windows](https://storage.googleapis.com/cdn.podops.dev/downloads/cli-0.9.7/po-windows-0.9.7.gz)

Unpack the archive and place the `po` binary in a directory that is on your `PATH`. 

To check your `PATH`, execute the following command:

```shell
$ echo $PATH
```

Make the binary execuatable:

```shell
$ chmod +x po
```

The binary SHOULD already be executable, but hey !

After you install the CLI, it is available using the `po` command:

```shell
po help
```

## Build from source

### Requirements

* [Go 1.15](https://golang.org/dl) or newer

### Clone the repository

```shell
$ git clone https://github.com/podops/podops.git
```

If you don't have git, you can download the source code as a file [archive from GitHub](https://github.com/podops/podops).
Each [release](https://github.com/podops/podops/releases) has a source snapshot.

### Build

```shell
$ cd podops/cmd/cli
$ go build po.go
```

## Build the services

Instructions on how to build the [API service](https://github.com/podops/podops/tree/main/cmd) and the [CDN service](https://github.com/podops/podops/tree/main/cmd) are currently only in README files on GitHub.