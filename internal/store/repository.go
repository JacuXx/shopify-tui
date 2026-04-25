package store

import "github.com/JacuXx/shopify-cli/internal/domain"

// Repository define el contrato para persistencia de tiendas.
// Permite cambiar la implementación (JSON, SQLite, etc.) sin tocar la UI.
type Repository interface {
	CargarTiendas() ([]domain.Tienda, error)
	GuardarTiendas(tiendas []domain.Tienda) error
	CrearDirectorioTienda(nombre string) (string, error)
	EliminarTienda(tiendas []domain.Tienda, indice int) []domain.Tienda
	ExisteDirectorio(ruta string) bool
}
