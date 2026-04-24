# ADR-0002: Arquitectura en Capas (Domain-Application-Infrastructure)

**Date**: 2026-04-23
**Status**: accepted
**Deciders**: Alan Reynoso Jacuinde

## Context

El proyecto actual (v1.5.x) mezcla lógica en archivos sin separación clara entre dominio, aplicación e infraestructura:

- `server.go`: gestión de procesos + estado + lógica de logs
- `model.go`: estado de UI + estado de negocio
- `commands.go`: llamadas CLI + persistencia + resultados

Esto genera código acoplado, difícil de testear sin procesos reales, y casi imposible cambiar de JSON a base de datos.

## Decision

**Implementar arquitectura en 4 capas**:

```
Presentation (Bubbletea) → Model, View, Update
        ↓
Application → UseCases, Services, Handlers
        ↓
Domain → Entities, Repositories (interfaces), ValueObjects
        ↓
Infrastructure → Storage (JSON/DB), CLI, Process management
```

**Regla**: Outer layers dependen de inner, **nunca al revés**.

## Alternatives Considered

### Alternative 1: Monolithic
- **Pros**: Rápido inicialmente
- **Cons**: Imposible testear, acoplado, hard to refactor
- **Why not**: El proyecto está creciendo

### Alternative 2: Plugin Architecture
- **Pros**: Plugins independientes
- **Cons**: Overkill, requiere IPC
- **Why not**: TUI es single-purpose

### Alternative 3: Hexagonal
- **Pros**: Desacoplado
- **Cons**: Más boilerplate
- **Why not**: Clean Architecture (4 capas) es más simple y suficiente

## Consequences

### Positive
- **Testabilidad**: Domain sin dependencies → tests sin mocks
- **Reusabilidad**: UseCases desde múltiples UIs
- **Mantenibilidad**: Cada capa tiene responsabilidad clara
- **Flexibilidad**: Cambiar de JSON a DB sin refactor masivo
- **Escalabilidad**: Agregar features sin afectar existentes

### Negative
- **Boilerplate**: Más archivos, más interfaces
- **Curva aprendizaje**: Entender flujo Domain → App → Infra
- **Over-engineering**: Features simples tocan 4 capas

### Risks
- **Shortcuts**: Saltarse capas porque "es rápido"
  - **Mitigation**: Code review, documentación
- **Interdependencias**: Presentation importa Infrastructure directamente
  - **Mitigation**: strict import rules

## Rationale

Después de v1.5.x de crecimiento, arquitectura actual no escala. Capas claras permiten crecimiento sin reescrituras.

La inversión upfront se recupera cuando: agregamos tests, cambiamos storage, reutilizamos lógica.

---

**Related**: [ADR-0001](0001-elm-architecture.md), [ADR-0003](0003-repository-pattern.md)
