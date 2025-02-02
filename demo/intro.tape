# -*-sh-*-

Output intro.gif
Set FontSize 20
Set Width 1000
Set Height 800
Set Theme { "name": "Whimsy", "black": "#535178", "red": "#ef6487", "green": "#5eca89", "yellow": "#fdd877", "blue": "#65aef7", "magenta": "#aa7ff0", "cyan": "#43c1be", "white": "#ffffff", "brightBlack": "#535178", "brightRed": "#ef6487", "brightGreen": "#5eca89", "brightYellow": "#fdd877", "brightBlue": "#65aef7", "brightMagenta": "#aa7ff0", "brightCyan": "#43c1be", "brightWhite": "#ffffff", "background": "#29283b", "foreground": "#b3b0d6", "selection": "#3d3c58", "cursor": "#b3b0d6" }
Set WindowBar Colorful
Set BorderRadius 10
Set Shell zsh
Set FontFamily "IBM Plex Mono"
Set CursorBlink false
Set PlaybackSpeed 1
Set TypingSpeed .05

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

Type "# Store something"
Enter
Sleep 1s
Type "anydb set kitty meow"
Enter
Sleep 3s

Enter
Type `# What's in the store?`
Enter
Sleep 1s
Type "anydb ls"
Enter
Sleep 3s

Enter
Type "# Fetch something"
Enter
Sleep 1s
Type "anydb get kitty"
Enter
Sleep 3s

Enter
Type "# Unicode also works, of course"
Enter
Sleep 1s
Type "anydb set 猫咪 喵"
Enter
Sleep 2s
Type "anydb get 猫咪"
Enter
Sleep 3s

Enter
Type "# Do creative things with anydb list"
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
Sleep 3s

