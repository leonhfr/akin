# akin

> synchronizes a repository of markdown files with anki

## Disclaimer

This project uses Readme-Driven Development. The readme describes the planned functionality of the software at the time of release. Currently akin is still in a pre-release state. Therefore, I make no guarantees as to whether the software works nor as to when it will be completed. All APIs and packages may change unexpectedly. 

I try to track progress and planned features in the [roadmap](./docs/ROADMAP.md).

## Motivation

Akin synchronizes a repository of markdown files with [anki](https://github.com/ankitects/anki). If you're here, then you probably already know how spaced repetition is efficient for memorizing. 

I like using Anki. It's open source, it has [AnkiWeb](https://ankiweb.net) which allows you to keep your cards synchronized across multiple machines, and it has [AnkiDroid](https://github.com/ankidroid/Anki-Android) which is a free Anki Android client. Alas, the desktop client ans its UI make it difficult to efficiently create cards. It neither exposes an API nor it has a CLI.

I needed a tool to automate the process of creating cards, as well as give more visibility and flexibility.

Akin does just that. It helps you manage your anki collection by parsing a repository of markdown files, and creates anki cards from it. The repository is the single source of truth.

Why markdown? It's lightweight, portable, shareable, and customizable.

It also promotes sharing of knowledge. As of now, we can download anki decks created by others, but we have very visibility of what they contain. A markdown repository is more transparent.

## Requirements

Akin requires:

- [Anki](https://github.com/ankitects/anki) > v2.1.35
- [AnkiConnect](https://github.com/FooSoft/anki-connect/) > v6

Previous versions might work but are untested.

Internally, akin uses [anki-connect](https://github.com/FooSoft/anki-connect/) to connect to anki, so you will need this extension installed in addition to the normal anki client.

## Installation

Akin currently uses the go toolchain:

```shell
go get github.com/leonhfr/akin
```

This command will download and build Akin. If you have your `$GOPATH` correctly configured, you will be able to invoke the `akin` command like any other.

## Quick start

TODO: make my anki repository public and add it as an example

```sh
akin -p ./path/to/md/repository -n -v
# ensure the modifications are as wished
akin -p ./path/to/md/repository -s -e
# also sync to anki web and exit anki
```

Akin accepts the verbose `-v` to print out exactly what changes it applies to the anki collection, as well as the dry run `-n` flag, to not make any actual changes. The combination `-v -n` is very useful if you want to see exactly what changes would be made without committing to them.

It is your responsibility to have version control on your markdown repository, akin only handles the parsing and export parts.

## Usage

```
Usage:
  akin synchronizes a local repository in Markdown with an Anki collection

  akin [OPTIONS]

Application Options:
      --anki-port=    AnkiConnect port to use (default: 8765)
      --anki-address= AnkiConnect address to use (default: localhost)
  -p, --path=         Path to the local Markdown file or directory (default: .)
  -n, --dry-run       Does not apply the changes
  -v, --verbose       Verbose logs
  -s, --sync          Synchronizes the local Anki collection with AnkiWeb
  -e, --exit          Gracefully exists Anki when done
      --version       Prints current akin version

Help Options:
  -h, --help          Show this help message
```

## License

MIT

## Author

[leonhfr](https://github.com/leonhfr)

## Credits

In particular, this project relies on:
- [anki](https://github.com/ankitects/anki)
- [anki-connect](https://github.com/FooSoft/anki-connect/)
- [AnkiDroid](https://github.com/ankidroid/Anki-Android)
- the [goldmark](https://github.com/yuin/goldmark) markdown parser

During the course of the project, I came across a few interesting talks, articles, and documentation. If this is of interest to you, they are listed in the [research](./docs/RESEARCH.md) document.