package tixiacompress

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCompressToBase64(t *testing.T) {
	input := []byte("simple compression content")

	encoded, err := CompressToBase64(input)
	if err != nil {
		t.Fatalf("CompressToBase64 returned error: %v", err)
	}

	decompressed, err := DecompressFromBase64(encoded)
	if err != nil {
		t.Fatalf("DecompressFromBase64 returned error: %v", err)
	}

	if !bytes.Equal(input, decompressed) {
		t.Fatalf("decompressed output mismatch: want %q, got %q", input, decompressed)
	}
}

func TestCompressFileToBase64(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "input.txt")
	content := []byte("file compression content")

	if err := os.WriteFile(filePath, content, 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	encoded, err := CompressFileToBase64(filePath)
	if err != nil {
		t.Fatalf("CompressFileToBase64 returned error: %v", err)
	}

	decompressed, err := DecompressFromBase64(encoded)
	if err != nil {
		t.Fatalf("DecompressFromBase64 returned error: %v", err)
	}

	if !bytes.Equal(content, decompressed) {
		t.Fatalf("decompressed output mismatch: want %q, got %q", content, decompressed)
	}
}

func TestDecompressFromBase64_InvalidInput(t *testing.T) {
	_, err := DecompressFromBase64("not-base64")
	if err == nil {
		t.Fatal("expected error for invalid base64 input, got nil")
	}
}

func TestDecompressFileFromBase64(t *testing.T) {
	content := []byte("decompress to file content")

	encoded, err := CompressToBase64(content)
	if err != nil {
		t.Fatalf("CompressToBase64 returned error: %v", err)
	}

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "output.txt")

	if err = DecompressFileFromBase64(encoded, outputPath); err != nil {
		t.Fatalf("DecompressFileFromBase64 returned error: %v", err)
	}

	writtenContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if !bytes.Equal(content, writtenContent) {
		t.Fatalf("written output mismatch: want %q, got %q", content, writtenContent)
	}
}
