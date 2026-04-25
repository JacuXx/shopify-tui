# Blueprint: Arquitectura Escalable — shopify-tui

**Objetivo:** Refactorizar shopify-tui de un monolito en package raíz a una arquitectura modular, escalable y testeable, sin romper ninguna funcionalidad existente.

**Modo:** Direct (git + GitHub CLI disponibles, pero migración es local-first)
**Fecha:** 2026-04-25
**Branch base:** main
**Pasos totales:** 8 (pasos 2-5 paralelos)

---

## Resumen de Arquitectura Propuesta

### Estructura de directorios objetivo

```
shopify-tui/
├── main.go                          # Entry point mínimo (sin cambios funcionales)
├── go.mod
├── go.sum
├── plans/
└── internal/
    ├── app/                         # Coordinador raíz (Elm root)
    │   ├── model.go                 # Model raíz slim (delega a sub-states)
    │   ├── update.go                # Update raíz (dispatcher a vistas)
    │   └── view.go                  # View raíz (dispatcher a vistas)
    ├── domain/                      # Tipos de dominio puro (sin dependencias)
    │   ├── store.go                 # Tienda struct, MetodoDescarga enum
    │   └── server.go                # ServidorActivo struct
    ├── views/                       # Una carpeta por vista
    │   ├── menu/
    │   │   ├── model.go             # MenuState struct
    │   │   ├── update.go            # updateMenu()
    │   │   └── view.go              # vistaMenu()
    │   ├── addstore/
    │   │   ├── model.go             # AddStoreState struct
    │   │   ├── update.go            # updateAgregarTienda()
    │   │   └── view.go              # vistaAgregarTienda()
    │   ├── selectstore/
    │   │   ├── model.go             # SelectStoreState struct
    │   │   ├── update.go            # updateSeleccionarTienda(), updateSeleccionarModo(), updateInputGit()
    │   │   └── view.go              # vistaSeleccionarTienda(), etc.
    │   ├── servers/
    │   │   ├── model.go             # ServersState struct
    │   │   ├── update.go            # updateServidores()
    │   │   └── view.go              # vistaServidores()
    │   ├── logs/
    │   │   ├── model.go             # LogsState struct
    │   │   ├── update.go            # updateLogs()
    │   │   └── view.go              # vistaLogs()
    │   └── popup/
    │       ├── model.go             # PopupState struct
    │       ├── update.go            # updatePopup()
    │       └── view.go              # vistaPopup()
    ├── store/                       # Capa de persistencia (Repository Pattern)
    │   ├── repository.go            # Interface: Repository
    │   └── json_repository.go       # Implementación JSON (lógica actual de store.go)
    ├── server/                      # Gestor de servidores (Strategy Pattern)
    │   ├── manager.go               # Interface: Manager
    │   └── process_manager.go       # Implementación con exec (lógica actual de server.go)
    ├── commands/                    # Comandos externos (Command Pattern)
    │   └── commands.go              # tea.Cmd wrappers (lógica actual de commands.go)
    ├── ui/                          # Componentes UI compartidos
    │   ├── styles/
    │   │   └── styles.go            # Variables Lipgloss (extraídas de view.go)
    │   ├── icons/
    │   │   └── icons.go             # Sistema de iconos (lógica actual de icons.go)
    │   └── components/
    │       ├── banner.go            # renderBanner()
    │       └── list.go              # crearLista(), helpers de listas
    └── version/                     # Verificación de actualizaciones
        └── version.go               # Lógica actual de version.go
```

---

## Patrones de Diseño Aplicados

### 1. Elm Architecture (preservado y mejorado)

El framework Bubbletea impone este patrón. Se mejora con sub-states por vista:

```go
// internal/app/model.go
type Model struct {
    // Navegación global
    vistaActual  Vista
    windowWidth  int
    windowHeight int

    // Sub-states por vista (composición)
    menu        menu.State
    addStore    addstore.State
    selectStore selectstore.State
    servers     servers.State
    logs        logs.State
    popup       popup.State

    // Dependencias inyectadas (interfaces)
    storeRepo  store.Repository
    serverMgr  server.Manager

    // Estado global compartido entre vistas
    tiendas     []domain.Tienda
    tiendaActual *domain.Tienda
    updateAvail bool
    newVersion  string
}
```

### 2. Repository Pattern (persistencia)

```go
// internal/store/repository.go
type Repository interface {
    CargarTiendas() ([]domain.Tienda, error)
    GuardarTiendas(tiendas []domain.Tienda) error
    CrearDirectorioTienda(nombre string) (string, error)
    EliminarDirectorioTienda(ruta string) error
}
```

Beneficios:
- Cambiar a SQLite/YAML en el futuro sin tocar la UI
- Testeable con mock repository

### 3. Strategy Pattern (gestor de servidores)

```go
// internal/server/manager.go
type Manager interface {
    IniciarServidor(tienda domain.Tienda, onLog func(string)) error
    DetenerServidor(nombre string) error
    DetenerTodos()
    ObtenerServidoresActivos() []domain.Tienda
    TieneServidorActivo(nombre string) bool
    ObtenerLogs(nombre string) []string
    ObtenerPuerto(nombre string) int
}
```

Beneficios:
- Elimina el singleton global `gestorGlobal`
- Testeable con mock manager

### 4. Command Pattern (preservado, relocalizado)

Las funciones que retornan `tea.Cmd` en `commands.go` ya implementan este patrón correctamente. Solo se mueven al package `internal/commands/`.

### 5. Sub-State Pattern (vistas)

Cada vista encapsula su propio estado:

```go
// internal/views/logs/model.go
package logs

type State struct {
    Scroll        int
    MaxScroll     int
    Seleccionando bool
    SelInicio     int
    SelFin        int
}

func NewState() State { return State{} }
```

```go
// internal/views/logs/update.go
func Update(s State, msg tea.Msg, logs []string) (State, tea.Cmd) { ... }

// internal/views/logs/view.go
func View(s State, logs []string, width, height int) string { ... }
```

---

## Invariantes del Sistema

Estas condiciones DEBEN mantenerse verdaderas después de cada paso:

1. `go build ./...` compila sin errores
2. La app arranca y muestra el menú principal
3. Se pueden agregar tiendas (ShopifyPull y GitClone)
4. Se pueden iniciar servidores de desarrollo
5. Los logs se muestran en tiempo real
6. La persistencia (stores.json) funciona igual
7. El popup y navegación funcionan igual

---

## Grafo de Dependencias

```
Step 1 (dominio + estructura)  ← BLOQUEANTE
    ├── Step 2 (store/repo)   ─┐
    ├── Step 3 (server/mgr)   ─┤ (paralelos entre sí)
    ├── Step 4 (commands/ver) ─┤
    └── Step 5 (ui shared)   ─┘
                                └── Step 6 (vistas)
                                          └── Step 7 (app root slim)
                                                    └── Step 8 (cleanup)
```

---

## STEP 1: Estructura de packages y tipos de dominio

**Branch:** `refactor/step-1-domain-packages`
**Es bloqueante:** todos los demás pasos dependen de este

### Archivos fuente relevantes
- `model.go` líneas 1-50: `Vista enum`, `MetodoDescarga enum`, `Tienda struct`
- `server.go` líneas 1-30: `ServidorActivo struct`
- `icons.go`: completo

### Tareas
- [ ] Crear estructura de directorios bajo `internal/`
- [ ] Crear `internal/domain/store.go`: mover `Tienda`, `MetodoDescarga`
- [ ] Crear `internal/domain/server.go`: mover `ServidorActivo` (sin mutex — es implementación)
- [ ] Crear `internal/ui/icons/icons.go`: copiar lógica de `icons.go` raíz
- [ ] Crear `internal/ui/styles/styles.go`: extraer variables Lipgloss de `view.go`
- [ ] Actualizar `main.go` para importar `internal/ui/icons`
- [ ] Verificar: `go build ./...`

### Exit Criteria
- [ ] `go build ./...` pasa
- [ ] Package `internal/domain` existe con `Tienda` y `MetodoDescarga`
- [ ] App funciona igual (`go run .`)

### Rollback
```bash
git checkout main -- .
```

---

## STEP 2: Repository Pattern para persistencia

**Branch:** `refactor/step-2-store-repository`
**Depende de:** Step 1

### Tareas
- [ ] Crear `internal/store/repository.go` con interface `Repository`
- [ ] Crear `internal/store/json_repository.go` con struct `JSONRepository`
  - Mover toda la lógica de `store.go` raíz
  - Path de datos: `~/.config/shopify-tui/stores.json` (sin cambios)
- [ ] Constructor: `func NewJSONRepository() *JSONRepository`
- [ ] Agregar campo `storeRepo store.Repository` al `Model`
- [ ] Inyectar en `modeloInicial()`: `store.NewJSONRepository()`
- [ ] Reemplazar llamadas directas a `cargarTiendas()` por `m.storeRepo.CargarTiendas()`
- [ ] Verificar: `go build ./...` + test de persistencia manual

### Exit Criteria
- [ ] `go build ./...` pasa
- [ ] Agregar tienda, reiniciar app, tienda persiste
- [ ] Interface `Repository` definida en `internal/store/repository.go`

---

## STEP 3: Strategy Pattern para gestor de servidores

**Branch:** `refactor/step-3-server-manager`
**Depende de:** Step 1 | **Paralelo con:** Step 2

### Tareas
- [ ] Crear `internal/server/manager.go` con interface `Manager`
- [ ] Crear `internal/server/process_manager.go` con struct `ProcessManager`
  - Mover lógica de `GestorServidores` de `server.go` raíz
  - Mantener goroutines de lectura de logs
  - Mantener thread-safety con `sync.RWMutex`
- [ ] Constructor: `func NewProcessManager() *ProcessManager`
- [ ] Agregar campo `serverMgr server.Manager` al `Model`
- [ ] Inyectar en `modeloInicial()`: `server.NewProcessManager()`
- [ ] Reemplazar `ObtenerGestor()` en `update.go` por `m.serverMgr`
- [ ] Verificar: `go build ./...` + test de servidor manual

### Exit Criteria
- [ ] `go build ./...` pasa
- [ ] Servidores inician y muestran logs en tiempo real
- [ ] Ninguna referencia a `gestorGlobal` o `ObtenerGestor()` en código activo

---

## STEP 4: Extraer packages auxiliares (commands + version)

**Branch:** `refactor/step-4-aux-packages`
**Depende de:** Step 1 | **Paralelo con:** Steps 2, 3

### Tareas
- [ ] Crear `internal/commands/commands.go`:
  - Copiar `commands.go` raíz, cambiar a `package commands`
  - Exportar funciones: `EjecutarShopifyLogin`, `EjecutarShopifyPull`, `EjecutarGitClone`, `EjecutarAbrirEditor`, `EjecutarAbrirTerminal`, etc.
  - Usar `domain.Tienda` para el tipo de tienda
- [ ] Crear `internal/version/version.go`:
  - Copiar `version.go` raíz, cambiar a `package version`
  - Exportar: `Version`, `NpmPackageName`, `VerificarActualizacion()`
- [ ] Actualizar referencias en `update.go` y `model.go`
- [ ] Verificar: `go build ./...`

### Exit Criteria
- [ ] `go build ./...` pasa
- [ ] Todos los comandos externos funcionan (shopify pull, git clone, etc.)
- [ ] Verificación de actualizaciones funciona

---

## STEP 5: Extraer componentes UI compartidos

**Branch:** `refactor/step-5-ui-components`
**Depende de:** Step 1 | **Paralelo con:** Steps 2, 3, 4

### Tareas
- [ ] Completar `internal/ui/styles/styles.go`: todas las variables `lipgloss.Style` de `view.go`
- [ ] Crear `internal/ui/components/banner.go`:
  - Extraer `renderBanner()` de `view.go`
  - Exportar: `RenderBanner(width int) string`
- [ ] Crear `internal/ui/components/list.go`:
  - Extraer `crearLista()`, `renderMenuConAtajos()` de `model.go`/`view.go`
- [ ] Actualizar referencias en `view.go`
- [ ] Verificar: `go build ./...` + verificar visualmente que UI es idéntica

### Exit Criteria
- [ ] `go build ./...` pasa
- [ ] La UI se ve exactamente igual (banner, estilos, listas)

---

## STEP 6: Dividir vistas en sub-packages

**Branch:** `refactor/step-6-views`
**Depende de:** Steps 2, 3, 4, 5
**Esfuerzo:** El más grande del plan

### Patrón a seguir por cada vista

Cada vista exporta tres elementos:
```go
type State struct { /* campos propios de la vista */ }
func NewState() State { return State{} }
func Update(s State, msg tea.Msg, ...) (State, tea.Cmd)
func View(s State, ...) string
```

### Sub-tareas por vista

**6a. `internal/views/menu/`**
- State: `{ lista list.Model }`
- Extraer: `updateMenu()` → `Update()`, `vistaMenu()` → `View()`

**6b. `internal/views/addstore/`**
- State: `{ inputs []textinput.Model, focusIdx int, metodo domain.MetodoDescarga }`
- Extraer: `updateAgregarTienda()`, `vistaAgregarTienda()`

**6c. `internal/views/selectstore/`**
- State: `{ lista list.Model, inputGit textinput.Model }`
- Extraer: `updateSeleccionarTienda()`, `updateSeleccionarModo()`, `updateInputGit()`

**6d. `internal/views/servers/`**
- State: `{ lista list.Model }`
- Extraer: `updateServidores()`, `vistaServidores()`

**6e. `internal/views/logs/`**
- State: `{ Scroll int, MaxScroll int, Seleccionando bool, SelInicio int, SelFin int }`
- Extraer: `updateLogs()`, `vistaLogs()`

**6f. `internal/views/popup/`**
- State: `{ lista list.Model, visible bool }`
- Extraer: `updatePopup()`, `vistaPopup()`

### Resultado en update.go raíz

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch m.vistaActual {
    case VistaMenu:
        m.menu, cmd = menu.Update(m.menu, msg, ...)
    case VistaLogs:
        m.logs, cmd = logs.Update(m.logs, msg, ...)
    // etc.
    }
    return m, cmd
}
```

### Exit Criteria
- [ ] `go build ./...` pasa
- [ ] Todas las vistas funcionan correctamente
- [ ] `update.go` raíz < 150 líneas (solo dispatch)
- [ ] `view.go` raíz < 100 líneas (solo dispatch)

---

## STEP 7: Consolidar Model raíz en internal/app

**Branch:** `refactor/step-7-slim-root-model`
**Depende de:** Step 6

### Tareas
- [ ] Crear `internal/app/model.go` con el Model slim (estructura mostrada arriba)
- [ ] Crear `internal/app/update.go`: dispatcher puro
- [ ] Crear `internal/app/view.go`: dispatcher puro
- [ ] Crear `internal/app/init.go`: función `New()` como constructor
- [ ] Actualizar `main.go`:
  ```go
  import "shopify-tui/internal/app"

  func main() {
      ui.InitIcons()
      p := tea.NewProgram(app.New(), tea.WithAltScreen(), tea.WithMouseAllMotion())
      p.Run()
  }
  ```
- [ ] Verificar: `go build ./...`

### Exit Criteria
- [ ] `go build ./...` pasa
- [ ] `main.go` importa `internal/app`
- [ ] El `Model` del package raíz ya no se usa

---

## STEP 8: Cleanup — eliminar archivos del package raíz

**Branch:** `refactor/step-8-cleanup`
**Depende de:** Step 7

### Archivos a eliminar del raíz (uno a uno con verificación)
- `model.go` → reemplazado por `internal/app/model.go` + sub-states
- `update.go` → reemplazado por `internal/app/update.go` + views
- `view.go` → reemplazado por `internal/app/view.go` + views
- `commands.go` → movido a `internal/commands/`
- `server.go` → movido a `internal/server/`
- `store.go` → movido a `internal/store/`
- `icons.go` → movido a `internal/ui/icons/`
- `version.go` → movido a `internal/version/`

### Test funcional completo (manual)
- [ ] App arranca y muestra menú
- [ ] Agregar tienda con ShopifyPull
- [ ] Agregar tienda con GitClone
- [ ] Iniciar servidor de desarrollo
- [ ] Logs aparecen en tiempo real
- [ ] Detener servidor
- [ ] Reiniciar app → tiendas siguen ahí (persistencia)
- [ ] Popup aparece y acciones funcionan
- [ ] Notificación de actualización funciona

### Exit Criteria
- [ ] `go build ./...` sin warnings
- [ ] `go vet ./...` limpio
- [ ] Solo `main.go`, `go.mod`, `go.sum`, `plans/` en el raíz
- [ ] Ningún archivo supera 300 líneas
- [ ] Test funcional completo pasa

---

## Beneficios Post-Migración

| Aspecto | Antes | Después |
|---------|-------|---------|
| Archivo más grande | `update.go` (912L) | ~150L (dispatcher) |
| Testabilidad | Imposible (singleton global) | Mockeable via interfaces |
| Agregar vista nueva | Tocar 3 archivos grandes | Crear `internal/views/nueva/` |
| Cambiar persistencia | Reescribir `store.go` | Implementar `store.Repository` |
| Separación de concerns | Todo mezclado | Clara por capa y por dominio |

---

## Regla de Oro para la Migración

> Cada branch debe compilar y funcionar antes de hacer merge.
> Nunca mover y eliminar en el mismo commit sin verificar build.

**Estrategia segura:**
1. Crear el nuevo file en `internal/`
2. Copiar el código (no mover todavía)
3. Actualizar referencias al nuevo package
4. Verificar que compila y funciona
5. Solo entonces eliminar el código original

---

*Blueprint generado para shopify-tui — 2026-04-25*
