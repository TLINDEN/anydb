## A personal key value store

[![Actions](https://github.com/tlinden/anydb/actions/workflows/ci.yaml/badge.svg)](https://github.com/tlinden/anydb/actions)
[![License](https://img.shields.io/badge/license-GPL-blue.svg)](https://github.com/tlinden/anydb/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tlinden/anydb)](https://goreportcard.com/report/github.com/tlinden/anydb)

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
- more features

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
  ```
  sudo make install
  ```
  
- You can also install from source. Issue the following commands in your shell:
  ```
  git clone https://github.com/TLINDEN/anydb.git
  cd anydb
  make
  sudo make install
  ```

If you  do not find a  binary release for your  platform, please don't
hesitate to ask me about it, I'll add it.

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
