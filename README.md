# Parcel

[![Documentation][godoc-img]][godoc-url]
![License][license-img]
[![Build Status][travis-img]][travis-url]
[![Coverage][coveralls-img]][coveralls-url]

An embeddable resource manager for Golang

## Overview

Parcel is a simple resource manager for Golang that allows embedding artefacts
like SQL and bash scripts. That allows easy release management by deploying
just a single binary rather than many files.

## Installation

```console
$ go get -u github.com/phogolabs/parcel
$ go install github.com/phogolabs/parcel/cmd/parcel
```

## Usage

The best way to use the tool is via `go generate`. In order to embed all
resource in particular directory, you should make it a package that has the
following comment:

```golang
// Package database contains the database artefacts of GOM as embedded resource
package database

//go:generate parcel -r
```

When you run:

```console
$ go generate ./...
```

The tools will create a `resource.go` file that contains
all embedded resource in that directory and its
subdirectories.

You can read the content in the following way:

```golang
file, err := parcel.ResourceManager.Open("your_sub_directory_name/your_file_name")
```
## Command Line Interface

```console
$ parcel -h

NAME:
   parcel - Golang Resource Embedding

USAGE:
   parcel [global options]

VERSION:
   0.1

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --quite, -q                     Disable logging
   --recursive, -r                 Embed the resources recursively
   --resource-dir value, -d value  Path to directory
   --package value, --pkg value    Package name
   --ignore value, -i value        Ignore file name
   --include-docs                  Include API documentation in generated source code
   --help, -h                      show help
   --version, -v                   print the version
```

## Example

You can check working example in our [GOM package](https://github.com/phogolabs/gom/tree/master/example).

## Contributing

We are welcome to any contributions. Just fork the
[project](https://github.com/phogolabs/parcel).

[coveralls-url]: https://coveralls.io/github/phogolabs/parcel
[coveralls-img]: https://coveralls.io/repos/github/phogolabs/parcel/badge.svg?branch=master
[travis-img]: https://travis-ci.org/phogolabs/parcel.svg?branch=master
[travis-url]: https://travis-ci.org/phogolabs/parcel
[parcel-url]: https://github.com/phogolabs/parcel
[godoc-url]: https://godoc.org/github.com/phogolabs/parcel
[godoc-img]: https://godoc.org/github.com/phogolabs/parcel?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[software-license-url]: LICENSE
