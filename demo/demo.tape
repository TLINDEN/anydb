# -*-sh-*-

Output demo.gif
Set FontSize 20
Set Width 1200
Set Height 1000
Set Theme { "name": "Whimsy", "black": "#535178", "red": "#ef6487", "green": "#5eca89", "yellow": "#fdd877", "blue": "#65aef7", "magenta": "#aa7ff0", "cyan": "#43c1be", "white": "#ffffff", "brightBlack": "#535178", "brightRed": "#ef6487", "brightGreen": "#5eca89", "brightYellow": "#fdd877", "brightBlue": "#65aef7", "brightMagenta": "#aa7ff0", "brightCyan": "#43c1be", "brightWhite": "#ffffff", "background": "#29283b", "foreground": "#b3b0d6", "selection": "#3d3c58", "cursor": "#b3b0d6" }
Set WindowBar Colorful
Set BorderRadius 10
Set Shell zsh
Set FontFamily "IBM Plex Mono"
Set CursorBlink false
Set PlaybackSpeed 1
Set TypingSpeed .2

Hide
Type `PROMPT=''`
Enter
Type "setopt interactivecomments"
Enter
Type "autoload -U colors && colors"
Enter
Type `PS1="%{$fg[magenta]%}demo> %{$reset_color%}"`
Enter
Type "clear"
Enter
Show

Type "anydb -h"
Enter
Sleep 5s

Type@0 "# Store something"
Enter
Sleep 1s
Type "anydb set kitty meow"
Enter
Sleep 4s

Type `# What's in the store?`
Enter
Sleep 1s
Type "anydb ls"
Enter
Sleep 4s

Type@0 "# Fetch something"
Enter
Sleep 1s
Type "anydb get kitty"
Enter
Sleep 4s


Type@0 "# Unicode also works, of course"
Enter
Sleep 1s
Type "anydb set 猫咪 喵"
Enter
Sleep 2s
Type "anydb get 猫咪"
Enter
Sleep 4s

Type@0 "# Do creative things with anydb list"
Enter
Sleep 1s
Type "anydb set penelope marmalade"
Enter
Type "anydb set christian tacos"
Enter
Type "anydb set muesli muesli"
Enter
Type "anydb list | xargs -n 2 printf '%s loves %s.\n'"
Enter
Sleep 4s

Type@0 "# you can assign tags"
Enter
Sleep 1s
Type "anydb set foo bar -t note,important"
Enter
Sleep 4s

Type@0 "# and filter for them"
Enter
Sleep 1s
Type "anydb list -t important"
Enter
Sleep 4s

Type@0 "# beside tags filtering you can also use regexps for searching"
Enter
Type@0 "# note, by default the list command only searches through keys"
Enter
Sleep 1s
Type "anydb list '[a-z]+\d'"
Enter
Sleep 4s

Type@0 "# do a full text search"
Enter
Sleep 1s
Type "anydb list '[a-z]+\d' -s"
Enter
Sleep 4s

Type@0 "# anydb also supports a wide output"
Enter
Sleep 1s
Type "anydb list -m wide"
Enter
Sleep 4s

Type@0 "# there are shortcuts as well"
Enter
Sleep 1s
Type "anydb ls -l"
Enter
Sleep 2s
Type "anydb /"
Enter
Sleep 4s

Type@0 "# other outputs are possible as well"
Enter
Sleep 1s
Type "anydb list -m json"
Enter
Sleep 4s

Type@0 "# you can backup your database"
Enter
Sleep 1s
Type "anydb export -o backup.json"
Enter
Sleep 4s

Type@0 "# and import it somewhere else"
Enter
Sleep 1s
Type "rm local.db"
Enter
Sleep 1s
Type "anydb ls -l"
Enter
Sleep 1s
Type "anydb import -i backup.json"
Enter
Sleep 1s
Type "anydb ls -l"
Enter
Sleep 4s

Type@0 "# you can encrypt entries. anydb asks for a passphrase"
Enter
Type@0 "# and will do the same when you retrieve the key using the"
Enter
Type@0 "# get command. anydb will ask you interactively for a password"
Enter
Sleep 1s
Type "anydb set address 'Beatstreet 42' -e"
Enter
Type "pass"
Enter
Sleep 4s

Type@0 "# but you can provide it via an environment variable too"
Enter
Sleep 1s
Type "ANYDB_PASSWORD=foo anydb set -e secretkey blahblah"
Enter
Sleep 4s

Type@0 "# using template output mode you can freely design how to print stuff"
Enter
Type@0 "# here, we print the values in CSV format ONLY if they have some tag"
Enter
Sleep 1s
Type `anydb ls -m template -T "{{ if .Tags }}{{ .Key }},{{ .Value }},{{ .Created}}{{ end }}"`
Enter
Sleep 4s

Type@0 "# or, to simulate skate's -k or -v"
Enter
Sleep 1s
Type `anydb ls -m template -T "{{ .Key }}"`
Enter
Sleep 1s
Type `anydb ls -m template -T "{{ .Value }}"`
Enter
Sleep 4s

# Type@0 "# maybe you want to digest the item in a shell script? also"
# Enter
# Type@0 "# note, that both the list and get commands support templates"
# Enter
# Sleep 1s
# Type `eval $(anydb get kitty -m template -T "value='{{ .Value }}'"); echo "value: $value"`
# Enter
# Sleep 4s

Type@0 "# sometimes you need to know some details about the current database"
Enter
Type@0 "# add -d for more details"
Enter
Sleep 1
Type "anydb info"
Enter
Sleep 4s

Type@0 "# it comes with a manpage builtin"
Enter
Sleep 1s
Type "anydb man"
Enter
Sleep 4s
Type "q"
Enter
Sleep 4s

Type@0 "# Try it out yourself: github.com/tlinden/anydb!"
Enter
Sleep 4s
