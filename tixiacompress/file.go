package tixiacompress

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ulikunitz/xz"
)

func CompressToBase64(content []byte) (string, error) {
	if len(content) == 0 {
		return "", fmt.Errorf("content is empty")
	}

	var compressed bytes.Buffer

	writer, err := xz.NewWriter(&compressed)
	if err != nil {
		return "", err
	}

	if _, err = writer.Write(content); err != nil {
		if closeErr := writer.Close(); closeErr != nil {
			return "", fmt.Errorf("write failed: %v, close failed: %w", err, closeErr)
		}
		return "", err
	}

	if err = writer.Close(); err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(compressed.Bytes())

	return encoded, nil
}

func CompressFileToBase64(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return CompressToBase64(content)
}

func DecompressFromBase64(encoded string) ([]byte, error) {
	if len(encoded) == 0 {
		return nil, fmt.Errorf("encoded content is empty")
	}

	compressedContent, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	if len(compressedContent) == 0 {
		return nil, fmt.Errorf("empty compressed content after base64 decoding")
	}

	reader, err := xz.NewReader(bytes.NewReader(compressedContent))
	if err != nil {
		return nil, err
	}

	decompressedContent, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return decompressedContent, nil
}

func DecompressFileFromBase64(encoded, outputPath string) error {
	decompressedContent, err := DecompressFromBase64(encoded)
	if err != nil {
		return err
	}

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err = os.WriteFile(outputPath, decompressedContent, 0o644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
