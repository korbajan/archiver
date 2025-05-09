package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/korbajan/archiver/pkg/archivers"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'pack' or 'unpack' subcommands")
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "pack":
		packCmd(os.Args[2:])
	case "unpack":
		unpackCmd(os.Args[2:])
	default:
		fmt.Printf("Unknown subcommand: %s\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println(`Usage:
		archiver pack   <zip|tar.gz> [compression-level] <source> <target>
		archiver unpack <zip|tar.gz> <source> <target>

		Examples:
		archiver pack zip 5 ./myfolder ./archive.zip
		archiver unpack tar.gz ./archive.tar.gz ./extracted
		`)
}

func packCmd(args []string) {
	// pack requires at least 3 args: format, source, target
	// optionally compression-level between format and source
	if len(args) < 3 {
		fmt.Println("pack command requires at least 3 arguments")
		usage()
		os.Exit(1)
	}

	format := args[0]
	var compLevel int = -1 // default compression level
	var src, dest string

	if len(args) == 3 {
		// No compression level specified
		src = args[1]
		dest = args[2]
	} else {
		// Compression level specified
		level, err := strconv.Atoi(args[1])
		if err != nil || level < -1 || level > 9 {
			fmt.Println("Invalid compression level; must be integer between -1 and 9")
			os.Exit(1)
		}
		compLevel = level
		src = args[2]
		dest = args[3]
	}

	arch, err := archiver.NewArchiver(format)
	if err != nil {
		fmt.Printf("Unsupported archive format: %s\n", format)
		os.Exit(1)
	}

	// Use type assertion to set compression level if supported
	err = archiver.SetCompressionLevel(arch, compLevel)
	if err != nil {
		fmt.Println(err)
	}

	err = arch.Archive(src, dest)
	if err != nil {
		fmt.Printf("Error during archiving: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully archived %s to %s using %s format\n", src, dest, format)
}

func unpackCmd(args []string) {
	// unpack requires exactly 3 args: format, source, target
	if len(args) != 3 {
		fmt.Println("unpack command requires exactly 3 arguments")
		usage()
		os.Exit(1)
	}

	format := args[0]
	src := args[1]
	dest := args[2]

	arch, err := archiver.NewArchiver(format)
	if err != nil {
		fmt.Printf("Unsupported archive format: %s\n", format)
		os.Exit(1)
	}

	err = arch.Unpack(src, dest)
	if err != nil {
		fmt.Printf("Error during unpacking: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully unpacked %s to %s using %s format\n", src, dest, format)
}
