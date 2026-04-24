package domain

import "context"

// TiendaRepository define el contrato para persistencia de tiendas
type TiendaRepository interface {
	GetAll(ctx context.Context) ([]Tienda, error)
	GetByID(ctx context.Context, id string) (*Tienda, error)
	Save(ctx context.Context, tienda Tienda) error
	Delete(ctx context.Context, id string) error
}

// ServidorRepository define el contrato para persistencia de servidores activos
type ServidorRepository interface {
	GetAll(ctx context.Context) ([]ServidorActivo, error)
	GetByID(ctx context.Context, id string) (*ServidorActivo, error)
	Save(ctx context.Context, servidor ServidorActivo) error
	Delete(ctx context.Context, id string) error
}
