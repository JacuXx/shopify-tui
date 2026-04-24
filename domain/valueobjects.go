package domain

import (
	"errors"
	"strings"
)

// StoreURL representa una URL de tienda Shopify validada
type StoreURL struct {
	raw string
}

// NewStoreURL crea una URL de tienda validada
func NewStoreURL(storeName string) (StoreURL, error) {
	if storeName = strings.TrimSpace(storeName); storeName == "" {
		return StoreURL{}, errors.New("store name is required")
	}

	// Si ya contiene .myshopify.com, usar tal cual
	if strings.Contains(storeName, ".myshopify.com") {
		return StoreURL{raw: storeName}, nil
	}

	// Si no, agregar el dominio
	return StoreURL{raw: storeName + ".myshopify.com"}, nil
}

func (s StoreURL) String() string {
	return s.raw
}

func (s StoreURL) StoreName() string {
	return strings.TrimSuffix(s.raw, ".myshopify.com")
}

// GitURL representa una URL de repositorio Git validada
type GitURL struct {
	raw string
}

// NewGitURL crea una URL de Git validada
func NewGitURL(url string) (GitURL, error) {
	if url = strings.TrimSpace(url); url == "" {
		return GitURL{}, errors.New("git URL is required")
	}

	if !strings.HasPrefix(url, "git@") && !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return GitURL{}, errors.New("invalid git URL format")
	}

	return GitURL{raw: url}, nil
}

func (g GitURL) String() string {
	return g.raw
}

// StoreName representa un nombre de tienda validado
type StoreName struct {
	raw string
}

// NewStoreName crea un nombre de tienda validado
func NewStoreName(name string) (StoreName, error) {
	if name = strings.TrimSpace(name); name == "" {
		return StoreName{}, errors.New("store name is required")
	}

	if len(name) > 100 {
		return StoreName{}, errors.New("store name too long")
	}

	return StoreName{raw: name}, nil
}

func (s StoreName) String() string {
	return s.raw
}
