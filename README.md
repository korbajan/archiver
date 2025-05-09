# Archiver-Go

A simple and extensible archiving tool written in Go, supporting ZIP archive creation and extraction.

## Features

- Archive a directory into a ZIP file
- Unpack a ZIP archive into a directory
- Extensible architecture to add more archive formats
- Command-line interface for easy use

## Installation

go install github.com/yourusername/archiver/cmd/archiver@latest

Or clone this repo and build the binary:

git clone https://github.com/korbajan/archiver.git && cd archiver && make build -o archiver ./cmd/archiver

Make sure your `$GOPATH/bin` is in your `PATH`.

## Usage

### Archive a directory

./archiver -src /path/to/source -dest /path/to/archive.zip -format zip


### Unpack an archive

./archiver -src /path/to/archive.zip -dest /path/to/destination -format zip -unpack


## Development

- Run tests:

make test


- Run linter (requires [staticcheck](https://staticcheck.io/)):

make lint


- Build the project:

make build


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Contributing

Contributions are welcome! Please open issues or pull requests for improvements and new features.

---

## Contact

korbajan - [jakub.kuba.stepien@gmail.com](mailto:jakub.kuba.stepien@gmail.com)

Project Link: https://github.com/korbajan/archiver
