# Scrolls CLI

Scrolls is a CLI tool for making, managing and using snippets/scripts in your terminal.
Make a snippet in one of the many supported languages and execute it whenever you need it or simply echo it to stdout.

![Scrolls CLI](https://scrolls.sh/demo.gif)

## Install

`curl -sSfL https://get.scrolls.sh/releases/install.sh | bash`

## Update

`scrolls update`

## Build

I am building from OSX using [GoReleaser](https://goreleaser.com/).

You might need a few dependencies to compile:

- [SQLite3](https://github.com/mattn/go-sqlite3?tab=readme-ov-file#macos)
- [musl-cross](https://github.com/FiloSottile/homebrew-musl-cross)


**TL;DR**

- Install sqlite3 dependency: `brew install sqlite3`
- Install Linux amd64 & arm64 toolchains: `brew install filosottile/musl-cross/musl-cross`
- Install Linux i386 toolchains: `brew install i686-linux-musl`

That is probably enough, otherwise go through the above links or ask your favourite LLM or something.
