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
          edit        Edit a key
          export      Export database to json
          get         Retrieve value for a key
          help        Help about any command
          import      Import database dump
          info        info
          list        List database contents
          man         show manual page
          serve       run REST API listener
          set         Insert key/value pair
    
        Flags:
          -b, --bucket string   use other bucket (default: data) (default "data")
          -c, --config string   toml config file
          -f, --dbfile string   DB file to use (default "/home/scip/.config/anydb/default.db")
          -d, --debug           Enable debugging
          -h, --help            help for anydb
          -v, --version         Print program version
    
        Use "anydb [command] --help" for more information about a command.

DESCRIPTION
    Anydb is a commandline personal key value store, it is simple to use and
    can be used to store anything you'd like, even binary files etc. It uses
    a key/value store (bbolt) in your home directory.

    The tool provides a number of subcommands to use it, there are global
    options and each subcommand has its own set of options.

GLOBAL OPTIONS
    "-f, --dbfile filename"
        The default location of your databas is
        "$HOME/.config/anydb/default.db". You can change this with the "-f"
        option.

    "-b, --bucket name"
        Data in a bbolt key-value-store are managed in so called buckets.
        These are kind of namespaces, where each key must be unique.
        However, a database may contain more than one bucket.

        By default anydb uses a bucket named "data", but you can change this
        using the option "-b".

        Buckets can be configured to always encrypt values, see ENCRYTPTION.

    "-c, --config filename"
        Under normal circumstances you don't need a configuration file. But
        if you want, you can provide one using the option "-c".

        Anydb looks for a couple of default locations for a config file. You
        only need this option if you want to supply a configuration on a
        non-standard location. See CONFIGURATION for more details.

    "-d, --debug"
        Enable debug output.

    "-h, --help"
        Show the usage of anydb.

    "-v, --version"
        Show the program version.

    All of these options can be used with subcommands as well.

SUBCOMMANDS
  completion
    The completion command can be used to setup completion for anydb. Just
    put something like this into your shell's configuration file:

        source <(anydb completion bash)

    If you use another shell, specify it instead of bash, of course.

  set
    The set command is being used to insert or update a key-value pair.

    Usage:

        Usage:
          anydb set <key> [<value> | -r <file>] [-t <tag>] [flags]
    
        Aliases:
          set, add, s, +
    
        Flags:
          -e, --encrypt            encrypt value
          -r, --file string        Filename or - for STDIN
          -h, --help               help for set
          -t, --tags tag,tag,...   tags, multiple allowed

    The standard way to insert a new entry is really simple:

        anydb set key value

    If you don't specify a value, anydb expects you to feed it some data via
    STDIN. For example:

        anydb set key < file

    You might as well specify a file directly using the "-f" option:

        anydb set key -f file

    Values can be encrypted using ChaCha20Poly1305 when you specify the "-e"
    option. Anydb will ask you interactively for a passphrase. You might as
    well provide the passphrase using the environment variable
    "ANYDB_PASSWORD". To encrypt the value, a cryptographically secure key
    will be derived from the passphrase using the ArgonID2 algorithm. Each
    value can be encrypted with another passphrase. So, the database itself
    is not encrypted, just the values.

    You can supply tags by using the option "-t". Multiple tags can be
    provided either by separating them with a comma or by using multiple
    "-t" parameters:

        anydb set key value -t tag1,tag2
        anydb set key value -t tag1 -t tag2

    You can later filter entries by tag or by a combination of tags.

    To edit or modify an entry, just use the set command with the same key,
    the value in the database will be overwritten with the new value. An
    alternative option is the edit command, see below.

  get
    To retrieve the value of a key, use the get subcommand.

    Usage:

        Usage:
          anydb get  <key> [-o <file>] [-m <mode>] [-n -N] [-T <tpl>] [flags]
    
        Aliases:
          get, show, g, .
    
        Flags:
          -h, --help              help for get
          -m, --mode string       output format (simple|wide|json|template) (default 'simple')
          -n, --no-headers        omit headers in tables
          -N, --no-human          do not translate to human readable values
          -o, --output string     output value to file (ignores -m)
          -T, --template string   go template for '-m template'

    In its simplest form you just call the get subcommand with the key you
    want to have the value for. The value is being printed to STDOUT by
    default:

        anydb get key

    If the value is binary content, it will not just being printed. In those
    cases you need to either redirect output into a file or use the option
    "-o" to write to a file:

        anydb get key > file
        anydb get key -o file

    If the value is encrypted, you will be asked for the passphrase to
    decrypt it. If the environment variable "ANYDB_PASSWORD" is set, its
    value will be used instead.

    There are different output modes you can choos from: simple, wide and
    json. The "simple" mode is the default one, it just prints the value as
    is. The "wide" mode prints a tabular output similar to the list
    subcommand, see there for more details. The options "-n" and "-N" have
    the same meaning as in the list command. The "json" mode prints the raw
    JSON representation of the whole database entry. Decryption will only
    take place in "simple" and "json" mode. The "template" mode provides the
    most flexibily, it is detailed in the section TEMPLATES.

  list
    The list subcommand displays a list of all database entries.

    Usage:

        Usage:
          anydb list  [<filter-regex>] [-t <tag>] [-m <mode>] [-n -N] [-T <tpl>] [flags]
    
        Aliases:
          list, /, ls
    
        Flags:
          -h, --help               help for list
          -m, --mode string        output format (table|wide|json|template), wide is a verbose table. (default 'table')
          -n, --no-headers         omit headers in tables
          -N, --no-human           do not translate to human readable values
          -t, --tags stringArray   tags, multiple allowed
          -T, --template string    go template for '-m template'
          -l, --wide-output        output mode: wide

    In its simplest form - without any options - , the list command just
    prints all keys with their values to STDOUT. Values are being truncated
    to maximum of 60 characters, that is, multiline values are not
    completely shown in order to keep the tabular view readable.

    To get more informations about each entry, use the "-o wide" or "-l"
    option. In addition to the key and value also the size, update timestamp
    and tags will be printed. Time and size values are converted into a
    human readable form, you can suppress this behavior with the "-N"
    option. You may omit the headers using the option "-n"

    Sometimes you might want to filter the list of entries. Either because
    your database grew too large or because you're searching for something.
    In that case you have two options: You may supply one or more tags or
    provide a filter regexp. To filter by tag, do:

        anydb list -t tag1
        anydb list -t tag1,tag2
        anydb list -t tag1 -t tag2

    To filter using a regular expression, do:

       anydb list "foo.*bar"

    Regular expressions follow the golang re2 syntax. For more details about
    the syntax, refer to <https://github.com/google/re2/wiki/Syntax>. Please
    note, that this regexp dialect is not PCRE compatible, but supports most
    of its features.

    You can - as with the get command - use other output modes. The default
    mode is "table". The "wide" mode is, as already mentioned, a more
    detailed table. Also supported is "json" mode and "template" mode. For
    details about using templates see TEMPLATES.

  del
    Use the del command to delete database entries.

    Usage:

        Usage:
          anydb del <key> [flags]
    
        Aliases:
          del, d, rm
    
        Flags:
          -h, --help   help for del

    The subcommand del does not provide any further options, it just deletes
    the entry referred to by the given key. No questions are being asked.

  edit
    The edit command makes it easier to modify larger entries.

    Usage:

       Usage:
          anydb edit <key> [flags]
    
        Aliases:
          edit, modify, mod, ed, vi
    
        Flags:
          -h, --help   help for edit

    The subcommand edit does not provide any further options. It works like
    this:

    1. Write the value info a temporary file.
    2. Execute the editor (which one, see below!) with that file.
    3. Now you can edit the file and save+close it when done.
    4. Anydb picks up the file and if the content has changed, puts its
    value into the DB.

    By default anydb executes the "vi" command. You can modify this behavior
    by setting the environment variable "EDITOR" appropriately.

    Please note, that this does not work with binary content!

  export
    Since the bbold database file is not portable across platforms (it is
    bound to the endianess of the CPU it was being created on), you might
    want to create a backup file of your database. You can do this with the
    export subcommand.

    Usage:

        Usage:
          anydb export [-o <json filename>] [flags]
    
        Aliases:
          export, dump, backup
    
        Flags:
          -h, --help            help for export
          -o, --output string   output to file

    The database dump is a JSON representation of the whole database and
    will be printed to STDOUT by default. Redirect it to a file or use the
    "-o" option:

        anydb export > dump.json
        anydb export -o dump.json

    Please note, that encrypted values will not be decrypted. This might
    change in a future version of anydb.

  import
    The import subcommand can be used to restore a database from a JSON
    dump.

    Usage:

        Usage:
          anydb import [<json file>] [flags]
    
        Aliases:
          import, restore
    
        Flags:
          -r, --file string        Filename or - for STDIN
          -h, --help               help for import
          -t, --tags stringArray   tags, multiple allowed

    By default the "import" subcommand reads the JSON contents from STDIN.
    You might pipe the dump into it or use the option "-r":

        anydb import < dump.json
        anydb import -r dump.json
        cat dump.json | anydb import

    If there is already a database, it will be saved by appending a
    timestamp and a new database with the contents of the dump will be
    created.

  serve
    Anydb provides a RESTful API, which you can use to manage the database
    from somewhere else. The API does not provide any authentication or any
    other security measures, so better only use it on localhost.

    Usage:

        Usage:
          anydb serve [-l host:port] [flags]
    
        Flags:
          -h, --help            help for serve
          -l, --listen string   host:port (default "localhost:8787")

    To start the listener, just execute the serve subcommand. You can tweak
    the ip address and tcp port using the "-l" option. The listener will not
    fork and run in the foreground. Logs are being printed to STDOUT as long
    as the listener runs.

    For more details about the API, please see the "REST API" section.

  info
    The info subcommand shows you some information about your current
    database.

    Usage:

        Usage:
          anydb info [flags]
    
        Flags:
          -h, --help       help for info
          -N, --no-human   do not translate to human readable values

    Data being shown are: filename and size, number of keys per bucket. If
    you supply the "-d" option (debug), some bbolt internals are being
    displayed as well.

  man
    The man subcommand shows an unformatted text variant of the manual page
    (which are currently reading).

    Usage:

        Usage:
          anydb man [flags]
    
        Flags:
          -h, --help   help for man

    The manual is being piped into the "more" command, which is being
    expected to exist according to the POSIX standard on all supported unix
    platforms. It might not work on Windows.

TEMPLATES
    The get and list commands support a template feature, which is very
    handy to create you own kind of formatting. The template syntax being
    used is the GO template language, refer to
    <https://pkg.go.dev/text/template> for details.

    Each template operates on one or more entries, no loop construct is
    required, the template provided applies to every matching entry
    separatley.

    The following template variables can be used:

    Key - string =item Value - string =item Bin - []byte =item Created -
    time.Time =item Tags - []string =item Encrypted bool

    Prepend a single dot (".") before each variable name.

    Here are some examples how to use the feature:

    Only show the keys of all entries:

        anydb list -m template -T "{{ .Key }}"

    Format the list in a way so that is possible to evaluate it in a shell:

        eval $(anydb get foo -m template -T "key='{{ .Key }}' value='{{ .Value }}' ts='{{ .Created}}'")
        echo "Key: $key, Value: $value"

    Print the values in CSV format ONLY if they have some tag:

        anydb list -m template -T "{{ if .Tags }}{{ .Key }},{{ .Value }},{{ .Created}}{{ end }}"

CONFIGURATION
    Anydb looks at the following location for a configuration file, in that
    order:

    "$HOME/.config/anydb/anydb.toml"
    "$HOME/.anydb.toml"
    "anydb.toml" in the current directory
    or specify one using "-c"
        The configuration format uses the TOML language, refer to
        <https://toml.io/en/> for more details. The key names correspond to
        the commandline options in most cases.

        Configuration follows a certain precedence: the files are tried to
        be read in the given order, followed by commandline options. That
        is, the last configuration file wins, unless the user provides a
        commandline option, then this setting will be taken.

        A complete configuration file might look like this:

            # defaults
            dbfile     = "~/.config/anydb/default.db"
            dbbucket   = "data"
            noheaders  = false
            nohumanize = false
            encrypt    = false
            listen     = "localhost:8787"
    
            # different setups for different buckets
            [buckets.data]
            encrypt = true
    
            [buckets.test]
            encrypt = false

        Under normal circumstances you don't need a configuration file.
        However, if you want to use different buckets, then this might be a
        handy option. Buckets are being configured in ini-style with the
        term "bucket." followed by the bucket name. In the example above we
        enable encryption for the default bucket "data" and disable it for a
        bucket "test". To use different buckets, use the "-b" option.

REST API
    The subcommand serve starts a simple HTTP service, which responds to
    RESTful HTTP requests. The listener responds to all requests with a JSON
    encoded response. The response contains the status and the content - if
    any - of the requested resource.

    The following requests are supported:

    GET /anydb/v1/
        Returns a JSON encoded list of all entries.

    GET /anydb/v1/key
        Returns the JSON encoded entry, if found.

    PUT /anydb/v1/
        Create an entry. Expects a JSON encoded request object in POST data.

    DELETE /anydb/v1/key
        Delete an entry.

    Some curl example calls to the API:

    Post a new key: curl -X PUT localhost:8787/anydb/v1/ \ -H 'Content-Type:
    application/json' \ -d '{"key":"foo","val":"bar"}'

    Retrieve the value:

        curl localhost:8787/anydb/v1/foo

    List all keys:

        curl localhost:8787/anydb/v1/

BUGS
    In order to report a bug, unexpected behavior, feature requests or to
    submit a patch, please open an issue on github:
    <https://github.com/TLINDEN/anydb/issues>.

    Please repeat the failing command with debugging enabled "-d" and
    include the output in the issue.

LIMITATIONS
    The REST API list request doesn't provide any filtering capabilities
    yet.

LICENSE
    This software is licensed under the GNU GENERAL PUBLIC LICENSE version
    3.

    Copyright (c) 2024 by Thomas von Dein

AUTHORS
    Thomas von Dein tom AT vondein DOT org

`
