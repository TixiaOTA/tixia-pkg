package tixiacompress

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCompressToBase64_EmptyContent(t *testing.T) {
	_, err := CompressToBase64([]byte{})
	if err == nil {
		t.Fatal("expected error for empty content, got nil")
	}
	if err.Error() != "content is empty" {
		t.Errorf("expected 'content is empty' error, got: %v", err)
	}
}

func TestDecompressFromBase64_EmptyEncoded(t *testing.T) {
	_, err := DecompressFromBase64("")
	if err == nil {
		t.Fatal("expected error for empty encoded string, got nil")
	}
	if err.Error() != "encoded content is empty" {
		t.Errorf("expected 'encoded content is empty' error, got: %v", err)
	}
}

func TestDecompressFromBase64_EmptyCompressedContent(t *testing.T) {
	// Base64 encoding of empty byte slice is empty string ""
	// which is already validated by empty encoded check.
	// The empty compressed content check happens AFTER base64 decode,
	// so we need to test a case where base64 decodes to empty but
	// is not an empty string itself (which is impossible).
	// This edge case is covered by the empty encoded check above.
}

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

func TestDecompressFileFromBase64_NestedDirectories(t *testing.T) {
	content := []byte("decompress to nested directory")

	encoded, err := CompressToBase64(content)
	if err != nil {
		t.Fatalf("CompressToBase64 returned error: %v", err)
	}

	tmpDir := t.TempDir()
	nestedPath := filepath.Join(tmpDir, "level1", "level2", "level3", "output.txt")

	if err = DecompressFileFromBase64(encoded, nestedPath); err != nil {
		t.Fatalf("DecompressFileFromBase64 returned error: %v", err)
	}

	writtenContent, err := os.ReadFile(nestedPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if !bytes.Equal(content, writtenContent) {
		t.Fatalf("written output mismatch: want %q, got %q", content, writtenContent)
	}
}
