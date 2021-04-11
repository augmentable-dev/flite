[![Go Reference](https://pkg.go.dev/badge/github.com/augmentable-dev/flite.svg)](https://pkg.go.dev/github.com/augmentable-dev/flite)
[![Go Report Card](https://goreportcard.com/badge/github.com/augmentable-dev/flite)](https://goreportcard.com/report/github.com/augmentable-dev/flite)

## flite

`flite` is a SQLite extension and command line utility for working with local data files.
It's meant to work in tandem with built-in functionality such as the [SQLite JSON1 extension](https://www.sqlite.org/json1.html).

### lines

`lines` is an [eponoymous-only virtual table](https://www.sqlite.org/vtab.html#eponymous_only_virtual_tables) (table-valued-function) that reads a file from disk (or stdin if no file is specified) by line.

```sql
SELECT * FROM lines("/path/to/some/file.ndjson")
```

### readfile

`readfile` is a scalar function that returns the contents of a file (path provided as an argument).
If no path is supplied, it reads from stdin.

```sql
SELECT readfile("/path/to/file.json")
``
