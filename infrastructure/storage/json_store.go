package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/JacuXx/shopify-cli/domain"
)

// JSONTiendaRepository implementa TiendaRepository usando JSON
type JSONTiendaRepository struct {
	filePath string
	mu       sync.RWMutex
}

// NewJSONTiendaRepository crea una nueva instancia
func NewJSONTiendaRepository(filePath string) *JSONTiendaRepository {
	return &JSONTiendaRepository{filePath: filePath}
}

// GetAll obtiene todas las tiendas
func (r *JSONTiendaRepository) GetAll(ctx context.Context) ([]domain.Tienda, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []domain.Tienda{}, nil
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var tiendas []domain.Tienda
	if err := json.Unmarshal(data, &tiendas); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return tiendas, nil
}

// GetByID obtiene una tienda por ID
func (r *JSONTiendaRepository) GetByID(ctx context.Context, id string) (*domain.Tienda, error) {
	tiendas, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for i := range tiendas {
		if tiendas[i].ID == id {
			return &tiendas[i], nil
		}
	}

	return nil, fmt.Errorf("tienda not found")
}

// Save guarda una tienda
func (r *JSONTiendaRepository) Save(ctx context.Context, tienda domain.Tienda) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tiendas, err := r.getAllUnsafe(ctx)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Buscar si ya existe
	found := false
	for i, t := range tiendas {
		if t.Nombre == tienda.Nombre {
			tiendas[i] = tienda
			found = true
			break
		}
	}

	if !found {
		tiendas = append(tiendas, tienda)
	}

	// Crear directorio si no existe
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(tiendas, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Delete elimina una tienda
func (r *JSONTiendaRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tiendas, err := r.getAllUnsafe(ctx)
	if err != nil {
		return err
	}

	var updated []domain.Tienda
	for _, t := range tiendas {
		if t.ID != id && t.Nombre != id {
			updated = append(updated, t)
		}
	}

	data, err := json.MarshalIndent(updated, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return os.WriteFile(r.filePath, data, 0644)
}

// getAllUnsafe es una versión sin lock (para uso interno)
func (r *JSONTiendaRepository) getAllUnsafe(ctx context.Context) ([]domain.Tienda, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []domain.Tienda{}, nil
		}
		return nil, err
	}

	var tiendas []domain.Tienda
	if err := json.Unmarshal(data, &tiendas); err != nil {
		return nil, err
	}

	return tiendas, nil
}
