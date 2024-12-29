## Features

- repl
- mime-type => exec app + value
- add waitgroup to db.go funcs
- RestList does not support any params?
- lc() incoming tags+keys

## DB Structure

- put tags into sub bucket see #1
- change structure to:

data bucket
key => {key,value[0:60],isbin:bool}

value bucket
key => value (maybe always use []byte here)

tags bucket
key/tag => tag/key
tag/key => tag

So, list just uses the data bucket, no large contents.
A tag search only looksup matching tags, see #1.
Only a full text search and get would need to dig into the value bucket.

A delete would just delete all keys from all values and then:
lookup in tags bucket for all key/*, then iterate over the values and
remove all tag/key's. Then deleting a key would not leave any residue
behind.

However, maybe change the list command to just list everything and add
an extra find command for fulltext or tag search. Maybe still provide
filter options in list command but only filter for keys.

DONE: most of the above, except the tag stuff. manpage needs update and tests.
