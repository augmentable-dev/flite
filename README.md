## jqlite

`jqlite` is a SQLite extension and command line utility for working with various forms of json.
It's meant to work in tandem with the [SQLite JSON1 extension](https://www.sqlite.org/json1.html).

### ndjson

`ndjson` is an [eponoymous-only virtual table](https://www.sqlite.org/vtab.html#eponymous_only_virtual_tables) (table-valued-function) that reads an [ndjson](https://github.com/ndjson/ndjson-spec) file from disk (or stdin if no file is specified)

```sql
SELECT * FROM ndjson("/path/to/some/file.ndjson")
```
