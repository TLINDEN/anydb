package cmd

var manpage = `
anydb
    anydb - a personal key value store

SYNOPSIS
        Usage:
          anydb <command> [options] [flags]
          anydb [command]
    
        Available Commands:
          completion  Generate the autocompletion script for the specified shell
          del         Delete key
          export      Export database to json
          get         Retrieve value for a key
          help        Help about any command
          import      Import database dump
          list        List database contents
          set         Insert key/value pair
    
        Flags:
          -f, --dbfile string   DB file to use (default "/home/scip/.config/anydb/default.db")
          -d, --debug           Enable debugging
          -h, --help            help for anydb
          -v, --version         Print program version
    
        Use "anydb [command] --help" for more information about a command.

DESCRIPTION
    Anydb is a simple to use commandline tool to store anything you'd like,
    even binary files etc. It uses a key/value store (bbolt) in your home
    directory.

LICENSE
    This software is licensed under the GNU GENERAL PUBLIC LICENSE version
    3.

    Copyright (c) 2024 by Thomas von Dein

AUTHORS
    Thomas von Dein tom AT vondein DOT org

`
var usage = `

Usage:
  anydb <command> [options] [flags]
  anydb [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  del         Delete key
  export      Export database to json
  get         Retrieve value for a key
  help        Help about any command
  import      Import database dump
  list        List database contents
  set         Insert key/value pair

Flags:
  -f, --dbfile string   DB file to use (default "/home/scip/.config/anydb/default.db")
  -d, --debug           Enable debugging
  -h, --help            help for anydb
  -v, --version         Print program version

Use "anydb [command] --help" for more information about a command.

`
