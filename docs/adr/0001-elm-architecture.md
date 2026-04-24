# ADR-0001: Usar Elm Architecture con Bubbletea

**Date**: 2026-04-23
**Status**: accepted
**Deciders**: Alan Reynoso Jacuinde

## Context

Shopify TUI es una aplicación interactiva en terminal que requiere gestión de estado complejo, respuesta rápida a eventos de usuario, facilidad de debugging y código predecible.

Las aplicaciones TUI pueden volverse desordenadas sin una arquitectura clara.

## Decision

**Continuamos usando Elm Architecture** implementada por Bubbletea.

Ciclo predecible:
```
User Input → Update() → new Model → View() → Render
    ↑                                            │
    └────────────────────────────────────────────┘
```

**Regla clave**: Model siempre es inmutable en cada ciclo.

## Alternatives Considered

### Alternative 1: Event-Driven
- **Pros**: Flexible
- **Cons**: State mutation difícil de debuggear, orden no garantizado
- **Why not**: TUI necesita orden de eventos predecible

### Alternative 2: MVC tradicional
- **Pros**: Familiar
- **Cons**: Difuso quién actualiza Model, bidireccional
- **Why not**: Necesitamos pasos discretos y predecibles

### Alternative 3: MVVM
- **Pros**: Separación clara
- **Cons**: Overkill, watchers implícitos
- **Why not**: Elm Architecture es más explícito

## Consequences

### Positive
- **Predictabilidad**: Código secuencial y determinístico
- **Debugging fácil**: Reproducir cualquier estado con secuencia de msgs
- **Testabilidad**: Update es función pura
- **Escalabilidad**: Escala bien a apps TUI complejas

### Negative
- **Curva de aprendizaje**: Conceptualmente diferente
- **Boilerplate**: Requiere definir tipos de msg
- **No mutable state**: Go permite mutabilidad; requiere disciplina

### Risks
- **Performance con Model grande**: cada Update() crea copia
  - **Mitigation**: Usar índices para logs
- **Asincronía compleja**: múltiples comandos pueden generar msgs en orden inesperado
  - **Mitigation**: Application Layer (próxima ADR)

## Rationale

Elm Architecture es la opción más **robusta y predecible** para TUI. Bubbletea es su implementación Go estándar.

---

**Related**: [ADR-0002](0002-layered-architecture.md)
