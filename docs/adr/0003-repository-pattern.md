# ADR-0003: Repository Pattern para Acceso a Datos

**Date**: 2026-04-23
**Status**: accepted
**Deciders**: Alan Reynoso Jacuinde

## Context

Actualmente, acceso a datos (persistencia de tiendas, servidores) está mezclado con lógica de negocio:

- `store.go`: `cargarTiendas()` y `guardarTiendas()` usan JSON directamente
- `server.go`: `GestorServidores` con mapa en memoria (no persistente)
- Cambiar de JSON a DB requiere tocar múltiples archivos
- Difícil testear sin archivos reales

## Decision

**Usar Repository Pattern**: interfaces en Domain Layer, implementaciones en Infrastructure.

```go
// domain/repositories.go
type TiendaRepository interface {
    GetAll(ctx context.Context) ([]Tienda, error)
    Save(ctx context.Context, tienda Tienda) error
    Delete(ctx context.Context, id string) error
}

// infrastructure/storage/json_store.go
type JSONTiendaRepository struct { ... }
func (r *JSONTiendaRepository) GetAll(ctx context.Context) ([]Tienda, error) { ... }
```

**Ventaja**: En tests, inyectamos `MockRepository{}` sin archivo real.

## Alternatives Considered

### Alternative 1: DAO
- **Pros**: Familiar en Java
- **Cons**: DAO y Repository casi lo mismo
- **Why not**: Repository es patrón moderno

### Alternative 2: CQRS
- **Pros**: Separación reads/writes, escalable
- **Cons**: Overkill para <10 entidades
- **Why not**: Repository es suficiente ahora

### Alternative 3: Active Record
- **Pros**: Menos boilerplate
- **Cons**: Tienda acoplada a persistencia, imposible testear
- **Why not**: Acoplamiento inaceptable

### Alternative 4: Global Service Locator
- **Pros**: Menos parámetros
- **Cons**: Anti-pattern, hidden dependencies
- **Why not**: DI explícita es preferible

## Consequences

### Positive
- **Testabilidad**: Inyectar mocks sin archivos reales
- **Flexibilidad**: JSON → SQLite sin refactor masivo
- **Explicitness**: Claro cuáles operaciones de datos existen
- **Concurrency**: RWMutex, prepared statements

### Negative
- **Boilerplate**: Interfaces + implementations
- **Parámetro passing**: Repositories en constructores (menos convenient)
- **Overhead**: CRUD simple requiere interface + impl

### Risks
- **Leaky Abstraction**: Repository "parecida" a DB API
  - **Mitigation**: Repository modela dominio, no DB
- **N+1 queries**: Con SQL, fácil cometer N+1
  - **Mitigation**: Documentar optimizaciones

## Rationale

Repository Pattern es estándar de facto en arquitecturas limpias. Desacopla lógica de persistencia, permitiendo testeo y flexibilidad.

Beneficios para Shopify TUI:
1. Tests sin archivos reales
2. Transición JSON → SQLite sin refactor
3. Claridad sobre datos

---

**Related**: [ADR-0002](0002-layered-architecture.md)
