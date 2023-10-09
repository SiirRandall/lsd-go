# lsd-go


**lsd-go** is a simple command-line tool written in Go for listing and exploring directories. It provides a user-friendly interface with various options to customize the way directory contents are displayed. This project draws considerable inspiration from lsd-rs, which is an outstanding Rust-based tool. While lsd-rs is a robust and highly recommended option, I chose to develop this tool in Go for my own learning purposes, and not to compete with or replace it.

I want to emphasize that there is absolutely nothing wrong with Rust; in fact, it's a fantastic programming language. My decision is solely based on personal preferences related to the community around Rust. I remain open to contributions and pull requests from anyone interested in this project.

Once again, I'd like to acknowledge the skill and excellence of the lsd-rs creator and their exceptional program.

I would also like to express this is a work in progress. So there maybe bugs and fetures missing. I will fix and add over time.

## Features

- List files and directories in a given directory.
- Customize the output with various flags and options.
- Display directory structures as a tree view.
- Sort files alphabetically and in reverse order.
- Show or hide dotfiles (hidden files).
- Display file and directory details.
- And more!

## Installation

To install **lsd-go**, you can use the Go toolchain:

```bash
go get github.com/SiirRandall/lsd-go
```
This will download and build the latest version of the tool and make it available in your Go workspace.

### Usage

```bash
lsd-go [options] [directory]
```

### Options

- `-a`: Show dotfiles.
- `--no-color`: Disable colored output.
- `--inodes`: Show inodes.
- `--headers`: Show headers.
- `-l`: List files and directories.
- `--alpha`: Sort files alphabetically.
- `--reverse`: Sort files in reverse order.
- `--dirsfirst`: Sort directories first and then files alphabetically.
- `--depth`: Maximum depth for directory traversal (-1 means no limit).
- `--tree`: Show tree view.

### Examples

- List the contents of the current directory:
  ```bash
  lsd-go
  ```
- List the contents of a specific directory:
```bash
lsd-go /path/to/directory
```
- Display a tree view of the directory structure:
```bash
lsd-go --tree
```
- List files and directories with details:
```bash
lsd-go -l
```



