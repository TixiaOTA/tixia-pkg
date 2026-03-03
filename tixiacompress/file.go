package tixiacompress

import (
	"bytes"
	"encoding/base64"
	"io"
	"os"

	"github.com/ulikunitz/xz"
)

func CompressToBase64(content []byte) (string, error) {
	var compressed bytes.Buffer

	writer, err := xz.NewWriter(&compressed)
	if err != nil {
		return "", err
	}

	if _, err = writer.Write(content); err != nil {
		_ = writer.Close()
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
	compressedContent, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
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

	if err = os.WriteFile(outputPath, decompressedContent, 0o644); err != nil {
		return err
	}

	return nil
}
