## A personal key value store

[![Actions](https://github.com/tlinden/anydb/actions/workflows/ci.yaml/badge.svg)](https://github.com/tlinden/anydb/actions)
[![License](https://img.shields.io/badge/license-GPL-blue.svg)](https://github.com/tlinden/anydb/blob/master/LICENSE)
[![Go Coverage](https://github.com/tlinden/anydb/wiki/coverage.svg)](https://raw.githack.com/wiki/tlinden/anydb/coverage.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/tlinden/anydb)](https://goreportcard.com/report/github.com/tlinden/anydb)
[![GitHub release](https://img.shields.io/github/v/release/tlinden/anydb?color=%2300a719)](https://github.com/TLINDEN/anydb/releases/latest)
[![Documentation](https://img.shields.io/badge/manpage-documentation-blue)](https://github.com/TLINDEN/anydb/blob/master/anydb.pod)

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

```shell
# Store something (and sync it to the network)
anydb set kitty meow

# Fetch something (from the local cache)
anydb get kitty

# What’s in the store?
anydb list

# Spaces are fine
anydb set "kitty litter" "smells great"

# You can store binary data, too
anydb set profile-pic < my-cute-pic.jpg
anydb get profile-pic > here-it-is.jpg

# Unicode also works, of course
anydb set 猫咪 喵
anydb get 猫咪

# For more info
anydb --help

# Do creative things with anydb list
anydb set penelope marmalade
anydb set christian tacos
anydb set muesli muesli

anydb list | xargs -n 2 printf '%s loves %s.\n'
```
  
However, there are more features than just that!

```shell
# you can assign tags
anydb set foo bar -t note,important

# and filter for them
anydb list -t important

# beside tags filtering you can also use regexps for searching
anydb list '[a-z]+\d'

# anydb also supports a wide output
anydb list -o wide
KEY     TAGS            SIZE    AGE             VALUE 
blah    important       4 B     7 seconds ago   haha 
foo                     3 B     15 seconds ago  bar
猫咪                    3 B     3 seconds ago   喵

# there are shortcuts as well
anydb ls -l
anydb /

# other outputs are possible as well
anydb list -o json

# you can backup your database
anydb export -o backup.json

# and import it somewhere else
anydb import -r backup.json

# you can encrypt entries. anydb asks for a passphrase
# and will do the same when you retrieve the key using the
# get command. anydb will ask you interactively for a password
anydb set mypassword -e

# but you can provide it via an environment variable too
ANYDB_PASSWORD=foo anydb set -e secretkey blahblah

# too tiresome to add -e every time you add an entry?
# use a per bucket config
cat ~/.config/anydb/anydb.toml
[buckets.data]
encrypt = true
anydb set foo bar # will be encrypted

# speaking of buckets, you can use different buckets
anydb -b test set foo bar

# and speaking of configs, you can place a config file at these places:
# ~/.config/anydb/anydb.toml
# ~/.anydb.toml
# anydb.toml (current directory)
# or specify one using -c <filename>
# look at example.toml

# using template output mode you can freely design how to print stuff
# here, we print the values in CSV format ONLY if they have some tag
anydb ls -m template -T "{{ if .Tags }}{{ .Key }},{{ .Value }},{{ .Created}}{{ end }}"

# or, to simulate skate's -k or -v
anydb ls -m template -T "{{ .Key }}"
anydb ls -m template -T "{{ .Value }}"

# maybe you want to digest the item in a shell script? also
# note, that both the list and get commands support templates
eval $(anydb get foo -m template -T "key='{{ .Key }}' value='{{ .Value }}' ts='{{ .Created}}'")
echo "$key: $value"

# run the restful api server
anydb serve

# post a new key
curl -X PUT localhost:8787/anydb/v1/ \
  -H 'Content-Type: application/json' \
  -d '{"key":"foo","val":"bar"}'

# retrieve it
curl localhost:8787/anydb/v1/foo

# list keys
curl localhost:8787/anydb/v1/

# as you might correctly suspect you can store multi-line values or
# the content of text files. but what to do if you want to change it?
# here's one way:
anydb get contract24 > file.txt && vi file.txt && anydb set contract24 -r file.txt

# annoying. better do this
anydb edit contract24 

# sometimes you need to know some details about the current database
# add -d for more details
anydb info

# it comes with a manpage builtin
anydb man
```

## Installation

There are multiple ways to install **anydb**:

- Go to the [latest release page](https://github.com/tlinden/anydb/releases/latest),
  locate the binary for your operating system and platform.
  
  Download it and put it into some directory within your `$PATH` variable.
  
- The release page also contains a tarball for every supported platform. Unpack it
  to some temporary directory, extract it and execute the following command inside:
  ```shell
  sudo make install
  ```
  
- You can also install from source. Issue the following commands in your shell:
  ```shell
  git clone https://github.com/TLINDEN/anydb.git
  cd anydb
  make
  sudo make install
  ```

- Or, if you have the GO toolkit installed, just install it like this:
  ```shell
  go install github.com/tlinden/anydb@latest
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

A list of available images is  [here](https://github.com/tlinden/anydb/pkgs/container/anydb/versions?filters%5Bversion_type%5D=tagged)


## Documentation

The  documentation  is  provided  as  a unix  man-page.   It  will  be
automatically installed if  you install from source.  However, you can
[read the man-page online](https://github.com/TLINDEN/anydb/blob/master/anydb.pod)

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
https://github.com/TLINDEN/anydb/issues.

## Copyright and license

This software is licensed under the GNU GENERAL PUBLIC LICENSE version 3.

## Authors

T.v.Dein <tom AT vondein DOT org>

## Project homepage

https://github.com/TLINDEN/anydb

## Copyright and License

Licensed under the GNU GENERAL PUBLIC LICENSE version 3.

## Author

T.v.Dein <tom AT vondein DOT org>
