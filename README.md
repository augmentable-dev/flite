[![Go Reference](https://pkg.go.dev/badge/github.com/augmentable-dev/flite.svg)](https://pkg.go.dev/github.com/augmentable-dev/flite)
[![Go Report Card](https://goreportcard.com/badge/github.com/augmentable-dev/flite)](https://goreportcard.com/report/github.com/augmentable-dev/flite)
[![codecov](https://codecov.io/gh/augmentable-dev/flite/branch/main/graph/badge.svg?token=ZLISOQV2WV)](https://codecov.io/gh/augmentable-dev/flite)

# flite

`flite` is a SQLite extension and command line utility for working with local data files.
It's meant to work in tandem with built-in functionality such as the [SQLite JSON1 extension](https://www.sqlite.org/json1.html).

## Usage

### SQLite Extension

`flite` can be compiled to a shared library and be loaded as a SQLite [runtime extension](https://sqlite.org/loadext.html).
Run `make` and the shared library will available be at `./build/flite.so`.

### Command Line Interface

`make` will also produce a binary at `./build/flite`.

## lines

`split` is an [eponoymous-only virtual table](https://www.sqlite.org/vtab.html#eponymous_only_virtual_tables) (table-valued-function) that reads a file from disk (or stdin if no file is specified) and splits it into rows by a delimiter (defaults to `\n`).

```sql
SELECT * FROM split("/path/to/some/file.ndjson")
```

## file_read

`file_read` is a scalar function that returns the contents of a file (path provided as an argument).
If no path is supplied, it reads from stdin.

```sql
SELECT file_read("/path/to/file.json")
```

## yaml_to_json

`yaml_to_json` is a scalar function that expects a single argument (a YAML string) and returns it as a JSON string (which can be used in the built-in JSON methods)

```sql
SELECT yaml_to_json("hello: world")
-- {"hello":"world"}
```

## json_to_yaml

`json_to_yaml` is a scalar function that expects a single argument (a JSON string) and returns it as a YAML string.

```sql
SELECT json_to_yaml('{"hello":"world"}')
-- hello: world
```
## cmd

`cmd` is a scalar function that expects a single argument (a bash command string) and it will return the results

```sql
SELECT cmd("echo 'Hello, World'")
-- Hello, World
```

## cmd_table

`cmd_table` is a module that takes in a bash string command and an optional delimiter(default "\n") returning a row for each line

| Column          | Type     |
|-----------------|----------|
| line_no         | INT      |
| contents        | TEXT     |

```sql
SELECT cmd_table("echo 'Hello, World'",' ')
-- 1 , Hello,
-- 2 , World
```