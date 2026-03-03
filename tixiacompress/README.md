# tixiacompress

Simple helper package for compressing data with xz and encoding the result as Base64.

## Overview

- Compress raw bytes to Base64 (`CompressToBase64`)
- Compress file content to Base64 (`CompressFileToBase64`)
- Decompress Base64+xz content to bytes (`DecompressFromBase64`)
- Decompress Base64+xz content directly to a file (`DecompressFileFromBase64`)

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/tixiacompress
```

## Usage

### Compress bytes to Base64

```go
package main

import (
	"fmt"
	"log"

	"github.com/TixiaOTA/tixia-pkg/tixiacompress"
)

func main() {
	encoded, err := tixiacompress.CompressToBase64([]byte("hello world"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(encoded)
}
```

### Compress file to Base64

```go
encoded, err := tixiacompress.CompressFileToBase64("/path/to/input.txt")
if err != nil {
	log.Fatal(err)
}
```

### Decompress Base64 to bytes

```go
raw, err := tixiacompress.DecompressFromBase64(encoded)
if err != nil {
	log.Fatal(err)
}

fmt.Println(string(raw))
```

### Decompress Base64 directly to file

```go
err := tixiacompress.DecompressFileFromBase64(encoded, "/path/to/output.txt")
if err != nil {
	log.Fatal(err)
}
```
