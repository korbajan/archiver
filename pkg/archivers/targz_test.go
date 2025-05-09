package archiver

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// readTarGzContents reads all files from a tar.gz archive into a map[path]content.
func readTarGzContents(t *testing.T, archivePath string) map[string]string {
	t.Helper()
	f, err := os.Open(archivePath)
	if err != nil {
		t.Fatalf("Failed to open archive: %v", err)
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	contents := make(map[string]string)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Error reading tar entry: %v", err)
		}
		if header.Typeflag != tar.TypeReg {
			continue // skip non-regular files
		}
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, tr); err != nil {
			t.Fatalf("Error reading file content: %v", err)
		}
		contents[header.Name] = buf.String()
	}

	return contents
}

func TestTarGzArchiver_Archive(t *testing.T) {
	sourceDir, expectedFiles := setupTestDir(t)
	archivePath := filepath.Join(t.TempDir(), "test_archive.tar.gz")

	tarGzArchiver := &TarGzArchiver{}
	err := tarGzArchiver.Archive(sourceDir, archivePath)
	if err != nil {
		t.Fatalf("Archive() error = %v", err)
	}

	actualFiles := readTarGzContents(t, archivePath)

	for path, expectedContent := range expectedFiles {
		actualContent, ok := actualFiles[path]
		if !ok {
			t.Errorf("Missing file in archive: %s", path)
			continue
		}
		if actualContent != expectedContent {
			t.Errorf("Content mismatch for %s: got %q, want %q", path, actualContent, expectedContent)
		}
	}

	if len(actualFiles) != len(expectedFiles) {
		t.Errorf("Archive contains unexpected files: got %d files, want %d", len(actualFiles), len(expectedFiles))
	}
}

func TestTarGzArchiver_Unpack(t *testing.T) {
	sourceDir, expectedFiles := setupTestDir(t)
	archivePath := filepath.Join(t.TempDir(), "test_archive.tar.gz")

	tarGzArchiver := &TarGzArchiver{}
	if err := tarGzArchiver.Archive(sourceDir, archivePath); err != nil {
		t.Fatalf("Archive() error = %v", err)
	}

	unpackDir := t.TempDir()
	if err := tarGzArchiver.Unpack(archivePath, unpackDir); err != nil {
		t.Fatalf("Unpack() error = %v", err)
	}

	actualFiles := readDirContents(t, unpackDir)

	for path, expectedContent := range expectedFiles {
		actualContent, ok := actualFiles[path]
		if !ok {
			t.Errorf("Missing unpacked file: %s", path)
			continue
		}
		if actualContent != expectedContent {
			t.Errorf("Content mismatch for unpacked file %s: got %q, want %q", path, actualContent, expectedContent)
		}
	}

	if len(actualFiles) != len(expectedFiles) {
		t.Errorf("Unpacked directory contains unexpected files: got %d files, want %d", len(actualFiles), len(expectedFiles))
	}
}
