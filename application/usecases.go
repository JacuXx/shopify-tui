package application

import (
	"context"
	"fmt"
	"sync"

	"github.com/JacuXx/shopify-cli/domain"
)

// ThemeDownloader define la interfaz para descargar temas
type ThemeDownloader interface {
	DownloadShopifyTheme(ctx context.Context, storePath, storeURL string) error
	CloneGitRepository(ctx context.Context, gitURL, storePath string) error
}

// DirectoryManager define la interfaz para gestión de directorios
type DirectoryManager interface {
	CreateStoreDirectory(storeName string) (string, error)
	StoreDirectoryExists(storePath string) bool
}

// AddStoreUseCase orquesta la lógica de agregar una tienda
type AddStoreUseCase struct {
	tiendaRepo domain.TiendaRepository
	downloader ThemeDownloader
	dirManager DirectoryManager
	mu         sync.RWMutex
}

// NewAddStoreUseCase crea una instancia del caso de uso
func NewAddStoreUseCase(
	tiendaRepo domain.TiendaRepository,
	downloader ThemeDownloader,
	dirManager DirectoryManager,
) *AddStoreUseCase {
	return &AddStoreUseCase{
		tiendaRepo: tiendaRepo,
		downloader: downloader,
		dirManager: dirManager,
	}
}

// Execute ejecuta el caso de uso
func (u *AddStoreUseCase) Execute(ctx context.Context, cmd AddStoreCommand) (*domain.Tienda, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Validar comando
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	// Crear directorio
	storePath, err := u.dirManager.CreateStoreDirectory(cmd.Nombre)
	if err != nil {
		return nil, fmt.Errorf("failed to create store directory: %w", err)
	}

	// Descargar tema
	method := domain.DownloadMethod(cmd.Metodo)
	if method == domain.MethodGitClone {
		if err := u.downloader.CloneGitRepository(ctx, cmd.GitURL, storePath); err != nil {
			return nil, fmt.Errorf("failed to clone git repository: %w", err)
		}
	} else {
		if err := u.downloader.DownloadShopifyTheme(ctx, storePath, cmd.URL); err != nil {
			return nil, fmt.Errorf("failed to download shopify theme: %w", err)
		}
	}

	// Crear entidad de dominio
	tienda := domain.Tienda{
		Nombre: cmd.Nombre,
		URL:    cmd.URL,
		Ruta:   storePath,
		Metodo: method,
		GitURL: cmd.GitURL,
	}

	// Guardar en repositorio
	if err := u.tiendaRepo.Save(ctx, tienda); err != nil {
		return nil, fmt.Errorf("failed to save store: %w", err)
	}

	return &tienda, nil
}

// ListStoresUseCase recupera todas las tiendas
type ListStoresUseCase struct {
	tiendaRepo domain.TiendaRepository
}

// NewListStoresUseCase crea una instancia
func NewListStoresUseCase(tiendaRepo domain.TiendaRepository) *ListStoresUseCase {
	return &ListStoresUseCase{tiendaRepo: tiendaRepo}
}

// Execute ejecuta el caso de uso
func (u *ListStoresUseCase) Execute(ctx context.Context) ([]domain.Tienda, error) {
	return u.tiendaRepo.GetAll(ctx)
}

// DeleteStoreUseCase elimina una tienda
type DeleteStoreUseCase struct {
	tiendaRepo domain.TiendaRepository
}

// NewDeleteStoreUseCase crea una instancia
func NewDeleteStoreUseCase(tiendaRepo domain.TiendaRepository) *DeleteStoreUseCase {
	return &DeleteStoreUseCase{tiendaRepo: tiendaRepo}
}

// Execute ejecuta el caso de uso
func (u *DeleteStoreUseCase) Execute(ctx context.Context, cmd DeleteStoreCommand) error {
	return u.tiendaRepo.Delete(ctx, cmd.StoreID)
}
