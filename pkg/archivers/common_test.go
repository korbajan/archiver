package archiver

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// Helper to create a temporary directory with test files
func setupTestDir(t *testing.T) (string, map[string]string) {
	t.Helper()

	dir := t.TempDir()

	files := map[string]string{
		"file1.txt":           "Hello, World!",
		"subdir/file2.txt":    "This is file 2",
		"subdir/nested/file3": "Nested file content",
	}

	for path, content := range files {
		fullPath := filepath.Join(dir, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
	}

	return dir, files
}

// Helper to read files from a ZIP archive into a map[path]content
func readZipContents(t *testing.T, zipPath string) map[string]string {
	t.Helper()

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("Failed to open zip file: %v", err)
	}
	defer r.Close()

	contents := make(map[string]string)
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("Failed to open file in zip: %v", err)
		}
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, rc)
		rc.Close()
		if err != nil {
			t.Fatalf("Failed to read file in zip: %v", err)
		}
		contents[f.Name] = buf.String()
	}

	return contents
}
