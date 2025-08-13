package postgresparser

import (
	"crypto/sha256"
	"fmt"
)

// calculateChecksum calculates SHA256 checksum of script content
func calculateChecksum(content string) string {
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", hash)
}
