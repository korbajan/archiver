package archiver

import (
	"fmt"
)

type CompressionLevel int

const (
	DefaultCompression CompressionLevel = -1
	NoCompression      CompressionLevel = 0
	BestSpeed          CompressionLevel = 1
	BestCompression    CompressionLevel = 9
)

// Archiver interface defines methods for archiving files.
type Archiver interface {
	Archive(sourceDir, destFile string) error
	Unpack(archiveFile, destDir string) error
}

// Factory function to create archivers by format name.
func NewArchiver(format string) (Archiver, error) {
	switch format {
	case "zip":
		return &ZipArchiver{}, nil
	case "tar.gz":
		return &TarGzArchiver{}, nil
	default:
		return nil, fmt.Errorf("unsupported archive format: %s", format)
	}
}

// CompressionLevelSetter is an interface for archivers supporting compression level
type CompressionLevelSetter interface {
	SetCompressionLevel(level int) error
}

// SetCompressionLevel sets compression level if archiver supports it
func SetCompressionLevel(a Archiver, level int) error {
	if level < -1 || level > 9 {
		return fmt.Errorf("invalid level, ignoring")
	}
	if setter, ok := a.(CompressionLevelSetter); ok {
		setter.SetCompressionLevel(level)
		return nil
	}
	return fmt.Errorf("this archiver does not support compression levels, ignoring")
}
