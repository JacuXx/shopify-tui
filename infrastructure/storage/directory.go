package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DirectoryService implementa la gestión de directorios
type DirectoryService struct {
	baseDir string
}

// NewDirectoryService crea una nueva instancia
func NewDirectoryService(baseDir string) *DirectoryService {
	return &DirectoryService{baseDir: baseDir}
}

// CreateStoreDirectory crea el directorio para una tienda
func (ds *DirectoryService) CreateStoreDirectory(storeName string) (string, error) {
	if err := os.MkdirAll(ds.baseDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create base directory: %w", err)
	}

	sanitized := sanitizeName(storeName)
	storePath := filepath.Join(ds.baseDir, sanitized)

	if err := os.MkdirAll(storePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create store directory: %w", err)
	}

	return storePath, nil
}

// StoreDirectoryExists verifica si el directorio existe
func (ds *DirectoryService) StoreDirectoryExists(storePath string) bool {
	info, err := os.Stat(storePath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// sanitizeName limpia el nombre de la tienda para usarlo como nombre de carpeta
func sanitizeName(name string) string {
	name = strings.ToLower(name)

	replacer := strings.NewReplacer(
		" ", "-",
		"_", "-",
		".", "-",
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "",
	)
	name = replacer.Replace(name)

	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}

	name = strings.Trim(name, "-")
	return name
}
