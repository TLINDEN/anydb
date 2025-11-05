## A personal key value store

[![status-badge](https://ci.codeberg.org/api/badges/15517/status.svg)](https://ci.codeberg.org/repos/15517)
[![License](https://img.shields.io/badge/license-GPL-blue.svg)](https://codeberg.org/scip/anydb/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/codeberg.org/scip/anydb)](https://goreportcard.com/report/codeberg.org/scip/anydb)
[![GitHub release](https://img.shields.io/github/v/release/tlinden/anydb?color=%2300a719)](https://codeberg.org/scip/anydb/releases)
[![Documentation](https://img.shields.io/badge/manpage-documentation-blue)](https://codeberg.org/scip/anydb/raw/branch/main/anydb.pod)

> [!CAUTION]
> Between version 0.1.0 and version 0.2.1 deletion of keys did not work. There
> is a unit test to check for this, but this unit test had a bug as well and
> didn't catch it. The bug and the test have been fixed. You are advised to
> upgrade to 0.2.1 and above.

> [!CAUTION]
> Version 0.1.3 introduced a [regression](https://codeberg.org/scip/anydb/issues/19),
> which caused the encryption feature not to work correctly anymore.
> If you are using anydb 0.1.3, you are urgently advised to
> upgrade to 0.2.0


Anydb is a simple to use commandline tool to store anything you'd
like, even binary files etc. It is a re-implementation of
[skate](https://github.com/charmbracelet/skate) for the following
reasons:

- it's just fun to do
- anydb uses [bbolt](https://github.com/etcd-io/bbolt) instead of
  badger. bbolt has a stable file format, which doesn't change
  anymore. The badger file format on the other hand changes very
  often, which is not good for a tool intended to be used for many
  years.
- more features:
  - output table in list mode uses TAB separator
  - better STDIN + pipe support
  - supports JSON output
  - supports more verbose tabular output
  - filtering using regular expressions
  - tagging
  - filtering using tags
  - encryption of entries
  - templates for custom output for maximum flexibility
  - includes a tiny web server, which serves a restful API

And I wrote a very similar [tool](https://www.daemon.de/projects/dbtool/) 24 years ago and wanted to do it again wit go.

**anydb** can do all the things you can do with skate:

![simple demo](https://codeberg.org/scip/anydb/raw/branch/main/demo/intro.gif)

However, there are more features than just that!

![advanced demo](https://codeberg.org/scip/anydb/raw/branch/main/demo/advanced.gif)

## Installation

There are multiple ways to install **anydb**:

- Go to the [latest release page](https://codeberg.org/scip/anydb/releases),
  locate the binary for your operating system and platform.
  
  Download it and put it into some directory within your `$PATH` variable.
  
- The release page also contains a tarball for every supported platform. Unpack it
  to some temporary directory, extract it and execute the following command inside:
  ```shell
  sudo make install
  ```
  
- You can also install from source. Issue the following commands in your shell:
  ```shell
  git clone https://codeberg.org/scip/anydb.git
  cd anydb
  make
  sudo make install
  ```

- Or, if you have the GO toolkit installed, just install it like this:
  ```shell
  go install codeberg.org/scip/anydb@latest
  ```

If you  do not find a  binary release for your  platform, please don't
hesitate to ask me about it, I'll add it.

### Using the docker image

A pre-built docker  image is available, which you can  use to test the
app without  installing it. To download:

```shell
docker pull ghcr.io/tlinden/anydb:latest
```

To execute anydb  inside the image do something like this:

```shell
mkdir mydb
docker run -ti  -v mydb:/db -u `id -u $USER` -e HOME=/db ghcr.io/tlinden/anydb:latest set foo bar
docker run -ti  -v mydb:/db -u `id -u $USER` -e HOME=/db ghcr.io/tlinden/anydb:latest list -o wide
```

Here, we operate in a local  directory `mydb`, which we'll use as HOME
inside  the  docker  container.  anydb  will  store  its  database  in
`mydb/.config/anydb/default.db`.

A list of available images is  [here](https://codeberg.org/scip/anydb/pkgs/container/anydb/versions?filters%5Bversion_type%5D=tagged)


## Documentation

The  documentation  is  provided  as  a unix  man-page.   It  will  be
automatically installed if  you install from source.  However, you can
[read the man-page online](https://codeberg.org/scip/anydb/blob/master/anydb.pod)

Or if you cloned  the repository you can read it  this way (perl needs
to be installed though): `perldoc anydb.pod`.

If you have the binary installed, you  can also read the man page with
this command:

    anydb man

## Getting help

Although I'm happy to hear from anydb users in private email, that's the
best way for me to forget to do something.

In order to report a bug,  unexpected behavior, feature requests or to
submit    a    patch,    please    open   an    issue    on    github:
https://codeberg.org/scip/anydb/issues.

## Copyright and license

This software is licensed under the GNU GENERAL PUBLIC LICENSE version 3.

## Authors

T.v.Dein <tom AT vondein DOT org>

## Project homepage

https://codeberg.org/scip/anydb

## Copyright and License

Licensed under the GNU GENERAL PUBLIC LICENSE version 3.

## Author

T.v.Dein <tom AT vondein DOT org>
