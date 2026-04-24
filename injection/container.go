package injection

import (
	"os"
	"path/filepath"

	"github.com/JacuXx/shopify-cli/application"
	"github.com/JacuXx/shopify-cli/infrastructure/cli"
	"github.com/JacuXx/shopify-cli/infrastructure/icons"
	"github.com/JacuXx/shopify-cli/infrastructure/storage"
)

// Container agrupa todas las dependencias de la aplicación
type Container struct {
	// Domain repositories
	TiendaRepository *storage.JSONTiendaRepository

	// Application use cases
	AddStoreUseCase    *application.AddStoreUseCase
	ListStoresUseCase  *application.ListStoresUseCase
	DeleteStoreUseCase *application.DeleteStoreUseCase

	// Application services
	ServerManager *application.ServerManager

	// Infrastructure services
	DirectoryService *storage.DirectoryService
	ShopifyService   *cli.ShopifyService
}

// NewContainer crea una nueva instancia del contenedor
func NewContainer() (*Container, error) {
	// Obtener directorio base
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	baseDir := filepath.Join(homeDir, ".config", "shopify-tui")
	storesDir := filepath.Join(baseDir, "stores")
	storesFile := filepath.Join(baseDir, "stores.json")

	// Inicializar iconos
	icons.InitIcons()

	// Crear servicios de infraestructura
	tiendaRepo := storage.NewJSONTiendaRepository(storesFile)
	dirService := storage.NewDirectoryService(storesDir)
	shopifyService := cli.NewShopifyService()

	// Crear use cases
	addStoreUseCase := application.NewAddStoreUseCase(tiendaRepo, shopifyService, dirService)
	listStoresUseCase := application.NewListStoresUseCase(tiendaRepo)
	deleteStoreUseCase := application.NewDeleteStoreUseCase(tiendaRepo)

	// Crear servicios de aplicación
	serverManager := application.NewServerManager()

	return &Container{
		TiendaRepository:   tiendaRepo,
		AddStoreUseCase:    addStoreUseCase,
		ListStoresUseCase:  listStoresUseCase,
		DeleteStoreUseCase: deleteStoreUseCase,
		ServerManager:      serverManager,
		DirectoryService:   dirService,
		ShopifyService:     shopifyService,
	}, nil
}
