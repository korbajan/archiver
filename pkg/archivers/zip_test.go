package archiver

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestZipArchiver_Archive(t *testing.T) {
	sourceDir, expectedFiles := setupTestDir(t)

	archivePath := filepath.Join(t.TempDir(), "test.zip")

	zipArchiver := &ZipArchiver{}
	err := zipArchiver.Archive(sourceDir, archivePath)
	if err != nil {
		t.Fatalf("Archive() error = %v", err)
	}

	// Read back the zip contents
	actualFiles := readZipContents(t, archivePath)

	// Check all expected files exist with correct content
	for path, expectedContent := range expectedFiles {
		// ZIP uses forward slashes
		zipPath := strings.ReplaceAll(path, string(os.PathSeparator), "/")

		actualContent, ok := actualFiles[zipPath]
		if !ok {
			t.Errorf("Missing file in archive: %s", zipPath)
			continue
		}
		if actualContent != expectedContent {
			t.Errorf("Content mismatch for %s: got %q, want %q", zipPath, actualContent, expectedContent)
		}
	}

	// Check no extra files present
	if len(actualFiles) != len(expectedFiles) {
		t.Errorf("Archive contains unexpected files: got %d files, want %d", len(actualFiles), len(expectedFiles))
	}
}

// readDirContents reads all files under a directory into a map[path]content
func readDirContents(t *testing.T, root string) map[string]string {
	t.Helper()

	files := make(map[string]string)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		// Normalize path separators to slash for comparison
		relPath = strings.ReplaceAll(relPath, string(os.PathSeparator), "/")
		files[relPath] = string(content)
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk directory: %v", err)
	}
	return files
}

func TestZipArchiver_Unpack(t *testing.T) {
	// Setup source directory with files
	sourceDir, expectedFiles := setupTestDir(t)
	archivePath := filepath.Join(t.TempDir(), "test_archive.zip")

	zipArchiver := &ZipArchiver{}

	// First archive the source directory
	err := zipArchiver.Archive(sourceDir, archivePath)
	if err != nil {
		t.Fatalf("Archive() error = %v", err)
	}

	// Create a new temp directory to unpack into
	unpackDir := t.TempDir()

	// Unpack the archive
	err = zipArchiver.Unpack(archivePath, unpackDir)
	if err != nil {
		t.Fatalf("Unpack() error = %v", err)
	}

	// Read contents of unpacked directory
	actualFiles := readDirContents(t, unpackDir)

	// Compare unpacked files with original files
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

	// Check no extra files present
	if len(actualFiles) != len(expectedFiles) {
		t.Errorf("Unpacked directory contains unexpected files: got %d files, want %d", len(actualFiles), len(expectedFiles))
	}
}
