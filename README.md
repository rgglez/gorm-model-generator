# gorm-model-generator

`gmg` is a command line tool written in Go which generates models for GORM.

## Build

```bash
$ make build
```

## Install

```bash
# make install
```

installs the command to `/usr/local/bin/gmg`.

## Clean build

```bash
$ make clean
```

## Command line arguments

* `--dsn` the DSN for the connection. Examplle: `user:pass@tcp(localhost:3306)/dbname`.
* `--type` the type of your database (mysql, postgres, sqlite).
* `--output` the output directory.
* `--tables` the table name (optional, for a specific table)

## Enviroment variables

You can pass the DSN via an enviroment variable instead of command line:

```bash
export DATABASE_DSN='user:pass@tcp(localhost:3306)/mydatabase'
```

## Examples

```bash
gmg --dsn='user:pass@tcp(localhost:3306)/mydatabase' --type=mysql --output=models/
```

```bash
gmg --dsn='user:pass@tcp(localhost:3306)/mydatabase' --type=mysql --output=models/ --tables=users,posts
```

## License

Copyright (C) 2026 Rodolfo González González.

Licensed under GPL-3.0. Read the [LICENSE](LICENSE) file for more information.


