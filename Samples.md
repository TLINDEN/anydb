for path in $(find . -type f); do file=$(basename $path | sed -e 's/.epub//'); autor=$(dirname $path | sed 's|^./||'); anydb set $file -t $autor $path; done

(echo "TITLE AUTHOR"; anydb ls -m template --template "{{.Key}} {{.Tags}}")|column -t | fzf --header-lines 1 --layout=reverse --info=inline --height 100% --pointer="◆" --separator="─" --scrollbar="│" --preview-window "right:50%" --preview "clear; unzip -p \$(anydb get {1}) *htm* | w3m -T text/html -dump | less"  "$@"
