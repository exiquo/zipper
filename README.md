# zipper

zipper is a simple CLI tool written in Go that archives a directory into a zip file.

## Features

- Archive a directory into a `.zip` file
- Keeps the source directory as the root inside the archive
- Validates input paths and prevents invalid usage

## Build
```bash
go build -o zipper
```

## Usage
```bash
zipper --src ./my-folder --out ./archive.zip
```

## Notes
- The source path must exist and be a directory
- The target file must have a .zip extension
- The target file must not be located inside the source directory