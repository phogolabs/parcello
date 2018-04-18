# Parcel

[![Documentation][godoc-img]][godoc-url]
![License][license-img]
[![Build Status][travis-img]][travis-url]
[![Coverage][codecov-img]][codecov-url]
[![Go Report Card][report-img]][report-url]

*Golang Resource Bundler*

[![Parcel][parcel-img]][parcel-url]

## Overview

Parcel is a simple resource manager for Golang that allows embedding assets
like SQL, bash scripts and images. That allows easy release management by
deploying just a single binary rather than many files.

## Roadmap

Note that we may introduce breaking changes until we reach v1.0.

* Rename the tool in order not to clash with [parcel-bundler](https://github.com/parcel-bundler/parcel)
* Support [http.FileSystem](https://golang.org/pkg/net/http/#FileSystem)
* Support embedded COFF resources

## Installation

#### GitHub

```console
$ go get -u github.com/phogolabs/parcel
$ go install github.com/phogolabs/parcel/cmd/parcel
```
#### Homebrew (for Mac OS X)

```console
$ brew tap phogolabs/tap
$ brew install parcel
```

## Usage

You can use the parcel command line interface to bundle the desired resources
recursively:

```console
$ parcel -r -d <resource_dir_source> -b <bundle_dir_destination>
```

However, the best way to use the tool is via `go generate`. In order to embed all
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
subdirectories as `tar.gz` archive which is registered in
[parcel.ResourceManager](https://github.com/phogolabs/parcel/blob/master/common.go#L6).

You can read the content in the following way:

```golang
// Import the package that includes 'resource.go'
import _ "database"

file, err := parcel.Open("your_sub_directory_name/your_file_name")
```

The `parcel` package provides an abstraction of
[FileSystem](https://godoc.org/github.com/phogolabs/parcel#FileSystem)
interface:

```golang
// FileSystem provides primitives to work with the underlying file system
type FileSystem interface {
	// Walk walks the file tree rooted at root, calling walkFn for each file or
	// directory in the tree, including root.
	Walk(dir string, fn filepath.WalkFunc) error
	// OpenFile is the generalized open call; most users will use Open
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
}
```

That is implemented by the following:

- [parcel.Manager](https://godoc.org/github.com/phogolabs/parcel#Manager) which provides an access to the bundled resources.
- [parcel.Dir](https://godoc.org/github.com/phogolabs/parcel#Dir) which provides an access to the underlying file system.

That allows easy replacement of the file system with the bundled resources and
vice versa.

Note that downsides of this resource embedding approach are that your compile
time may increase significantly. We are working on better approach by using
`syso`. Stay tuned.

## Command Line Interface

```console
$ parcel -h

NAME:
   parcel - Golang Resource Bundler

USAGE:
   parcel [global options]

VERSION:
   0.1

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --quiet, -q                     Disable logging
   --recursive, -r                 Embed the resources recursively
   --resource-dir value, -d value  Path to directory (default: ".")
   --bundle-dir value, -b value    Path to bundle directory (default: ".")
   --ignore value, -i value        Ignore file name
   --include-docs                  Include API documentation in generated source code
   --help, -h                      show help
   --version, -v                   print the version
```

## Example

You can check working example in our [OAK package](https://github.com/phogolabs/oak/tree/master/example).

## Contributing

We are welcome to any contributions. Just fork the
[project](https://github.com/phogolabs/parcel).

*logo made by [Good Wave][logo-author-url] [CC 3.0][logo-license]*

[report-img]: https://goreportcard.com/badge/github.com/phogolabs/parcel
[report-url]: https://goreportcard.com/report/github.com/phogolabs/parcel
[logo-author-url]: https://www.flaticon.com/authors/good-ware
[logo-license]: http://creativecommons.org/licenses/by/3.0/
[parcel-url]: https://github.com/phogolabs/parcel
[parcel-img]: doc/img/logo.png
[codecov-url]: https://codecov.io/gh/phogolabs/parcel
[codecov-img]: https://codecov.io/gh/phogolabs/parcel/branch/master/graph/badge.svg
[travis-img]: https://travis-ci.org/phogolabs/parcel.svg?branch=master
[travis-url]: https://travis-ci.org/phogolabs/parcel
[parcel-url]: https://github.com/phogolabs/parcel
[godoc-url]: https://godoc.org/github.com/phogolabs/parcel
[godoc-img]: https://godoc.org/github.com/phogolabs/parcel?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[software-license-url]: LICENSE
