# gorm-model-generator

[![License](https://img.shields.io/badge/GitHub-GPL--3.0-informational)](https://www.gnu.org/licenses/gpl-3.0.html)
![GitHub all releases](https://img.shields.io/github/downloads/rgglez/gorm-model-generator/total)
![GitHub issues](https://img.shields.io/github/issues/rgglez/gorm-model-generator)
![GitHub commit activity](https://img.shields.io/github/commit-activity/y/rgglez/gorm-model-generator)
[![Go Report Card](https://goreportcard.com/badge/github.com/rgglez/gorm-model-generator)](https://goreportcard.com/report/github.com/rgglez/gorm-model-generator)
[![GitHub release](https://img.shields.io/github/release/rgglez/gorm-model-generator.svg)](https://github.com/rgglez/gorm-model-generator/releases/)
![GitHub stars](https://img.shields.io/github/stars/rgglez/gorm-model-generator?style=social)
![GitHub forks](https://img.shields.io/github/forks/rgglez/gorm-model-generator?style=social)

`gmg` is a command line tool written in Go which generates [models](https://en.wikipedia.org/wiki/Object%E2%80%93relational_mapping) for [GORM](https://gorm.io/index.html).

## Installation

### Clone the repository

```bash
git clone https://github.com/rgglez/gorm-model-generator.git
```

### Build

As a regular user:

```bash
make build
```

### Install

As root:

```bash
make install
```

installs the command to `/usr/local/bin/gmg`.

### Clean build

```bash
$ make clean
```

## Command line arguments

* `--dsn` the DSN for the connection. Example: `user:pass@tcp(localhost:3306)/dbname`.
* `--type` the type of your database (`mysql`, `postgres`, `sqlite`).
* `--output` the output directory.
* `--tables` the table names (optional, comma-separated list of specific tables)

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
export DATABASE_DSN='user:pass@tcp(localhost:3306)/mydatabase'
gmg --type=mysql --output=models/ --tables=users,posts
```

## License

Copyright (C) 2026 Rodolfo González González.

Licensed under GPL-3.0. Read the [LICENSE](LICENSE) file for more information.


