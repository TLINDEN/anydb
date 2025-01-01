## Features

- repl
- mime-type => exec app + value
- add waitgroup to db.go funcs
- RestList does not support any params?

## DB Structure

- put tags into sub bucket see #1

tags bucket
key/tag => tag/key
tag/key => tag

A delete would just delete all keys from all values and then:
lookup in tags bucket for all key/*, then iterate over the values and
remove all tag/key's. Then deleting a key would not leave any residue
behind.
